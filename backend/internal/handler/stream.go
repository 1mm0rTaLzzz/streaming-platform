package handler

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/jmoiron/sqlx"
	"stream_site/backend/internal/model"
)

type StreamHandler struct {
	db *sqlx.DB
}

func NewStreamHandler(db *sqlx.DB) *StreamHandler {
	return &StreamHandler{db: db}
}

func (h *StreamHandler) List(c *gin.Context) {
	var streams []model.Stream
	if err := h.db.Select(&streams, "SELECT * FROM streams ORDER BY match_id, priority DESC"); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"streams": streams, "total": len(streams)})
}

func (h *StreamHandler) Create(c *gin.Context) {
	var input struct {
		MatchID        int    `json:"match_id" binding:"required"`
		URL            string `json:"url" binding:"required"`
		Label          string `json:"label"`
		LanguageCode   string `json:"language_code"`
		Region         string `json:"region"`
		CommentaryType string `json:"commentary_type"`
		Quality        string `json:"quality"`
		IsActive       *bool  `json:"is_active"`
		Priority       int    `json:"priority"`
	}
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if input.LanguageCode == "" {
		input.LanguageCode = "en"
	}
	if input.Region == "" {
		input.Region = "global"
	}
	if input.CommentaryType == "" {
		input.CommentaryType = "full"
	}
	if input.Quality == "" {
		input.Quality = "720p"
	}
	if input.MatchID <= 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "match_id must be a positive integer"})
		return
	}

	var matchExists bool
	if err := h.db.Get(&matchExists, "SELECT EXISTS(SELECT 1 FROM matches WHERE id = $1)", input.MatchID); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	if !matchExists {
		c.JSON(http.StatusBadRequest, gin.H{"error": "match not found for given match_id"})
		return
	}

	isActive := true
	if input.IsActive != nil {
		isActive = *input.IsActive
	}

	var id int
	err := h.db.QueryRow(`
		INSERT INTO streams (match_id, url, label, language_code, region, commentary_type, quality, is_active, priority)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9) RETURNING id`,
		input.MatchID, input.URL, input.Label, input.LanguageCode, input.Region, input.CommentaryType, input.Quality, isActive, input.Priority,
	).Scan(&id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusCreated, gin.H{"id": id})
}

func (h *StreamHandler) Update(c *gin.Context) {
	id, _ := strconv.Atoi(c.Param("id"))
	var input struct {
		URL            string `json:"url" binding:"required"`
		Label          string `json:"label"`
		LanguageCode   string `json:"language_code"`
		Region         string `json:"region"`
		CommentaryType string `json:"commentary_type"`
		Quality        string `json:"quality"`
		IsActive       *bool  `json:"is_active"`
		Priority       int    `json:"priority"`
	}
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	_, err := h.db.Exec(
		`UPDATE streams SET url=$1, label=$2, language_code=$3, region=$4, commentary_type=$5, quality=$6, is_active=COALESCE($7, is_active), priority=$8 WHERE id=$9`,
		input.URL, input.Label, input.LanguageCode, input.Region, input.CommentaryType, input.Quality, input.IsActive, input.Priority, id,
	)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"ok": true})
}

func (h *StreamHandler) Delete(c *gin.Context) {
	id, _ := strconv.Atoi(c.Param("id"))
	h.db.Exec("DELETE FROM streams WHERE id = $1", id)
	c.JSON(http.StatusOK, gin.H{"ok": true})
}
