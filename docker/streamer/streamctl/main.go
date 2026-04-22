package main

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"os/exec"
	"strings"
	"sync"
	"syscall"
	"time"
)

type activeStream struct {
	Key       string      `json:"key"`
	URL       string      `json:"url"`
	StartedAt time.Time   `json:"started_at"`
	Tail      *tailBuffer `json:"-"`
	procs     []*exec.Cmd
}

var (
	mu     sync.Mutex
	active *activeStream
)

type tailBuffer struct {
	mu      sync.Mutex
	lines   []string
	limit   int
	partial string
}

func newTailBuffer(limit int) *tailBuffer {
	return &tailBuffer{limit: limit}
}

func (t *tailBuffer) Write(p []byte) (int, error) {
	t.mu.Lock()
	defer t.mu.Unlock()

	combined := t.partial + string(p)
	parts := strings.Split(combined, "\n")
	t.partial = parts[len(parts)-1]
	for _, line := range parts[:len(parts)-1] {
		line = strings.TrimSpace(line)
		if line == "" {
			continue
		}
		t.lines = append(t.lines, line)
		if len(t.lines) > t.limit {
			t.lines = t.lines[len(t.lines)-t.limit:]
		}
	}
	return len(p), nil
}

func (t *tailBuffer) String() string {
	t.mu.Lock()
	defer t.mu.Unlock()

	lines := append([]string{}, t.lines...)
	if strings.TrimSpace(t.partial) != "" {
		lines = append(lines, strings.TrimSpace(t.partial))
	}
	return strings.Join(lines, "\n")
}

func main() {
	mux := http.NewServeMux()
	mux.HandleFunc("POST /launch", handleLaunch)
	mux.HandleFunc("POST /stop", handleStop)
	mux.HandleFunc("GET /status", handleStatus)
	mux.HandleFunc("POST /debug", handleDebug)
	log.Println("streamctl ready on :8090")
	log.Fatal(http.ListenAndServe(":8090", mux))
}

func writeJSON(w http.ResponseWriter, code int, v any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	json.NewEncoder(w).Encode(v)
}

func handleLaunch(w http.ResponseWriter, r *http.Request) {
	var req struct {
		URL    string `json:"url"`
		Key    string `json:"key"`
		Width  int    `json:"width"`
		Height int    `json:"height"`
		FPS    int    `json:"fps"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil || req.URL == "" || req.Key == "" {
		writeJSON(w, http.StatusBadRequest, map[string]string{"error": "url and key required"})
		return
	}
	if req.Width == 0 {
		req.Width = 1920
	}
	if req.Height == 0 {
		req.Height = 1080
	}
	if req.FPS == 0 {
		req.FPS = 30
	}

	mu.Lock()
	defer mu.Unlock()
	if active != nil {
		writeJSON(w, http.StatusConflict, map[string]string{"error": "stream already active, stop first"})
		return
	}

	procs, tail, err := startAll(req.URL, req.Key, req.Width, req.Height, req.FPS)
	if err != nil {
		writeJSON(w, http.StatusInternalServerError, map[string]string{"error": err.Error()})
		return
	}
	active = &activeStream{Key: req.Key, URL: req.URL, StartedAt: time.Now(), Tail: tail, procs: procs}
	writeJSON(w, http.StatusOK, map[string]string{"ok": "true", "key": req.Key})
}

func handleStop(w http.ResponseWriter, r *http.Request) {
	var req struct {
		Key string `json:"key"`
	}
	json.NewDecoder(r.Body).Decode(&req)

	mu.Lock()
	defer mu.Unlock()
	if active == nil {
		writeJSON(w, http.StatusOK, map[string]string{"ok": "true", "msg": "nothing active"})
		return
	}
	if req.Key != "" && req.Key != active.Key {
		writeJSON(w, http.StatusNotFound, map[string]string{"error": "key mismatch"})
		return
	}
	stopProcs(active.procs)
	active = nil
	writeJSON(w, http.StatusOK, map[string]string{"ok": "true"})
}

func handleStatus(w http.ResponseWriter, r *http.Request) {
	mu.Lock()
	defer mu.Unlock()
	if active == nil {
		writeJSON(w, http.StatusOK, map[string]any{"active": nil})
		return
	}
	writeJSON(w, http.StatusOK, map[string]any{
		"active": map[string]any{
			"key":         active.Key,
			"url":         active.URL,
			"started_at":  active.StartedAt,
			"uptime_s":    int(time.Since(active.StartedAt).Seconds()),
			"stderr_tail": func() string {
				if active.Tail == nil {
					return ""
				}
				return active.Tail.String()
			}(),
		},
	})
}

func handleDebug(w http.ResponseWriter, r *http.Request) {
	var req struct {
		Enable bool   `json:"enable"`
		URL    string `json:"url"`
	}
	json.NewDecoder(r.Body).Decode(&req)

	mu.Lock()
	defer mu.Unlock()

	if !req.Enable {
		if active != nil {
			stopProcs(active.procs)
			active = nil
		}
		writeJSON(w, http.StatusOK, map[string]string{"ok": "true", "msg": "debug stopped"})
		return
	}
	if active != nil {
		writeJSON(w, http.StatusConflict, map[string]string{"error": "stop active stream first"})
		return
	}
	url := req.URL
	if url == "" {
		url = "https://accounts.google.com"
	}
	procs, err := startDebug(url)
	if err != nil {
		writeJSON(w, http.StatusInternalServerError, map[string]string{"error": err.Error()})
		return
	}
	active = &activeStream{Key: "__debug__", URL: url, StartedAt: time.Now(), procs: procs}
	writeJSON(w, http.StatusOK, map[string]string{
		"ok":  "true",
		"msg": "debug started. CDP on :9222. Run: ssh -L 9222:localhost:9222 your-vps then open chrome://inspect",
	})
}

func clearProfileLocks() {
	for _, f := range []string{"/profile/SingletonLock", "/profile/SingletonCookie", "/profile/SingletonSocket"} {
		os.Remove(f)
	}
}

func startAll(url, key string, w, h, fps int) ([]*exec.Cmd, *tailBuffer, error) {
	clearProfileLocks()
	display := ":99"
	res := fmt.Sprintf("%dx%d", w, h)
	rtmp := fmt.Sprintf("%s/%s", rtmpBase(), key)
	tail := newTailBuffer(40)

	xvfb := spawnProc(nil, nil, "Xvfb", display, "-screen", "0", res+"x24", "-ac", "-nolisten", "tcp")
	time.Sleep(700 * time.Millisecond)

	chrome := spawnProc(
		[]string{"DISPLAY=" + display},
		nil,
		"chromium",
		"--no-sandbox",
		"--disable-gpu",
		"--use-gl=swiftshader",
		"--no-zygote",
		"--disable-dev-shm-usage",
		"--disable-extensions",
		"--disable-audio-output",
		"--mute-audio",
		"--user-data-dir=/profile",
		"--kiosk",
		"--app="+url,
		"--autoplay-policy=no-user-gesture-required",
	)
	time.Sleep(5 * time.Second) // wait for page to load

	ffmpegCmd := spawnProc(
		[]string{"DISPLAY=" + display},
		io.MultiWriter(os.Stderr, tail),
		"ffmpeg", "-hide_banner", "-loglevel", "warning",
		"-f", "x11grab", "-video_size", res, "-framerate", fmt.Sprint(fps), "-i", display,
		"-f", "lavfi", "-i", fmt.Sprintf("anullsrc=sample_rate=44100:channel_layout=stereo"),
		"-c:v", "libx264", "-preset", "veryfast", "-tune", "zerolatency",
		"-b:v", "2500k", "-g", fmt.Sprint(fps*2), "-pix_fmt", "yuv420p",
		"-c:a", "aac", "-b:a", "64k",
		"-shortest", "-f", "flv", rtmp,
	)

	// Verify ffmpeg stays up for at least 2s (catches RTMP connection failures)
	exitCh := make(chan error, 1)
	go func() { exitCh <- ffmpegCmd.Wait() }()
	select {
	case err := <-exitCh:
		stopProcs([]*exec.Cmd{xvfb, chrome})
		return nil, nil, fmt.Errorf("ffmpeg exited early: %v", err)
	case <-time.After(2 * time.Second):
	}

	return []*exec.Cmd{xvfb, chrome, ffmpegCmd}, tail, nil
}

func startDebug(url string) ([]*exec.Cmd, error) {
	clearProfileLocks()
	display := ":99"
	xvfb := spawnProc(nil, nil, "Xvfb", display, "-screen", "0", "1920x1080x24", "-ac", "-nolisten", "tcp")
	time.Sleep(700 * time.Millisecond)
	chrome := spawnProc(
		[]string{"DISPLAY=" + display},
		nil,
		"chromium",
		"--no-sandbox", "--disable-gpu", "--use-gl=swiftshader", "--no-zygote", "--disable-dev-shm-usage",
		"--user-data-dir=/profile",
		"--remote-debugging-port=9222",
		"--remote-debugging-address=0.0.0.0",
		"--remote-allow-origins=*",
		url,
	)
	return []*exec.Cmd{xvfb, chrome}, nil
}

func spawnProc(extraEnv []string, stderr io.Writer, name string, args ...string) *exec.Cmd {
	cmd := exec.Command(name, args...)
	cmd.Stdout = os.Stdout
	if stderr != nil {
		cmd.Stderr = stderr
	} else {
		cmd.Stderr = os.Stderr
	}
	if extraEnv != nil {
		cmd.Env = append(os.Environ(), extraEnv...)
	}
	if err := cmd.Start(); err != nil {
		log.Printf("WARN: failed to start %s: %v", name, err)
	} else {
		log.Printf("started %s (pid %d)", name, cmd.Process.Pid)
	}
	return cmd
}

func stopProcs(procs []*exec.Cmd) {
	for i := len(procs) - 1; i >= 0; i-- {
		cmd := procs[i]
		if cmd == nil || cmd.Process == nil {
			continue
		}
		if err := cmd.Process.Signal(syscall.SIGTERM); err != nil {
			cmd.Process.Kill()
		}
	}
	done := make(chan struct{})
	go func() {
		for _, cmd := range procs {
			if cmd != nil {
				cmd.Wait()
			}
		}
		close(done)
	}()
	select {
	case <-done:
	case <-time.After(4 * time.Second):
		for _, cmd := range procs {
			if cmd != nil && cmd.Process != nil {
				cmd.Process.Kill()
			}
		}
	}
	log.Println("all processes stopped")
}

func rtmpBase() string {
	if v := os.Getenv("RTMP_BASE"); v != "" {
		return v
	}
	return "rtmp://rtmp:1935/live"
}
