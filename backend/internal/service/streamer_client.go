package service

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"
)

type StreamerClient struct {
	baseURL string
	http    *http.Client
}

func NewStreamerClient(baseURL string) *StreamerClient {
	return &StreamerClient{
		baseURL: baseURL,
		http:    &http.Client{Timeout: 15 * time.Second},
	}
}

type LaunchRequest struct {
	URL    string `json:"url"`
	Key    string `json:"key"`
	Width  int    `json:"width,omitempty"`
	Height int    `json:"height,omitempty"`
	FPS    int    `json:"fps,omitempty"`
}

type StreamStatusResponse struct {
	Active *struct {
		Key        string    `json:"key"`
		URL        string    `json:"url"`
		StartedAt  time.Time `json:"started_at"`
		UptimeS    int       `json:"uptime_s"`
		StderrTail string    `json:"stderr_tail"`
	} `json:"active"`
}

func (c *StreamerClient) Launch(req LaunchRequest) error {
	return c.post("/launch", req)
}

func (c *StreamerClient) Stop(key string) error {
	return c.post("/stop", map[string]string{"key": key})
}

func (c *StreamerClient) Status() (*StreamStatusResponse, error) {
	resp, err := c.http.Get(c.baseURL + "/status")
	if err != nil {
		return nil, fmt.Errorf("streamer unreachable: %w", err)
	}
	defer resp.Body.Close()
	var s StreamStatusResponse
	if err := json.NewDecoder(resp.Body).Decode(&s); err != nil {
		return nil, err
	}
	return &s, nil
}

func (c *StreamerClient) Debug(enable bool, url string) error {
	return c.post("/debug", map[string]any{"enable": enable, "url": url})
}

func (c *StreamerClient) post(path string, body any) error {
	b, _ := json.Marshal(body)
	resp, err := c.http.Post(c.baseURL+path, "application/json", bytes.NewReader(b))
	if err != nil {
		return fmt.Errorf("streamer unreachable: %w", err)
	}
	defer resp.Body.Close()
	if resp.StatusCode >= 300 {
		data, _ := io.ReadAll(resp.Body)
		return fmt.Errorf("streamer %d: %s", resp.StatusCode, data)
	}
	return nil
}
