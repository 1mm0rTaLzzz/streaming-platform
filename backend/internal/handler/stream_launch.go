package handler

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/jmoiron/sqlx"
	"stream_site/backend/internal/service"
)

type StreamLaunchHandler struct {
	client *service.StreamerClient
	db     *sqlx.DB
}

func NewStreamLaunchHandler(client *service.StreamerClient, db *sqlx.DB) *StreamLaunchHandler {
	return &StreamLaunchHandler{client: client, db: db}
}

func (h *StreamLaunchHandler) Launch(c *gin.Context) {
	var req struct {
		URL     string `json:"url" binding:"required"`
		Key     string `json:"key" binding:"required"`
		MatchID *int   `json:"match_id"`
		Width   int    `json:"width"`
		Height  int    `json:"height"`
		FPS     int    `json:"fps"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	if err := h.client.Launch(service.LaunchRequest{
		URL:    req.URL,
		Key:    req.Key,
		Width:  req.Width,
		Height: req.Height,
		FPS:    req.FPS,
	}); err != nil {
		c.JSON(http.StatusBadGateway, gin.H{"error": err.Error()})
		return
	}
	hlsURL := fmt.Sprintf("/live/%s.m3u8", req.Key)
	if req.MatchID != nil {
		h.db.Exec(`
			INSERT INTO streams (match_id, url, label, language_code, region, commentary_type, quality, is_active, priority)
			VALUES ($1, $2, 'Live auto', 'en', 'global', 'full', '1080p', true, 100)
			ON CONFLICT DO NOTHING`,
			*req.MatchID, hlsURL,
		)
	}
	c.JSON(http.StatusOK, gin.H{"ok": true, "key": req.Key, "hls": hlsURL})
}

func (h *StreamLaunchHandler) Stop(c *gin.Context) {
	var req struct {
		Key string `json:"key"`
	}
	c.ShouldBindJSON(&req)
	if err := h.client.Stop(req.Key); err != nil {
		c.JSON(http.StatusBadGateway, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"ok": true})
}

func (h *StreamLaunchHandler) Status(c *gin.Context) {
	s, err := h.client.Status()
	if err != nil {
		c.JSON(http.StatusBadGateway, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, s)
}

func (h *StreamLaunchHandler) Debug(c *gin.Context) {
	var req struct {
		Enable bool   `json:"enable"`
		URL    string `json:"url"`
	}
	c.ShouldBindJSON(&req)
	if err := h.client.Debug(req.Enable, req.URL); err != nil {
		c.JSON(http.StatusBadGateway, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"ok": true})
}
