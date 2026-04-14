package handler

import (
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/jmoiron/sqlx"
	"stream_site/backend/internal/model"
	wslib "stream_site/backend/internal/ws"
)

type MatchHandler struct {
	db  *sqlx.DB
	hub *wslib.Hub
}

func NewMatchHandler(db *sqlx.DB) *MatchHandler {
	return &MatchHandler{db: db}
}

// SetHub injects the WebSocket hub so UpdateScore can broadcast score changes.
func (h *MatchHandler) SetHub(hub *wslib.Hub) {
	h.hub = hub
}

func (h *MatchHandler) List(c *gin.Context) {
	status := c.Query("status")
	stage := c.Query("stage")
	dateStr := c.Query("date")

	query := `
		SELECT m.*,
			ht.id as "home_team.id", ht.code as "home_team.code", ht.name_en as "home_team.name_en", ht.name_ru as "home_team.name_ru", ht.flag_url as "home_team.flag_url",
			at.id as "away_team.id", at.code as "away_team.code", at.name_en as "away_team.name_en", at.name_ru as "away_team.name_ru", at.flag_url as "away_team.flag_url"
		FROM matches m
		LEFT JOIN teams ht ON m.home_team_id = ht.id
		LEFT JOIN teams at ON m.away_team_id = at.id
		WHERE 1=1`
	args := []interface{}{}
	argIdx := 1

	if status != "" {
		query += ` AND m.status = $` + strconv.Itoa(argIdx)
		args = append(args, status)
		argIdx++
	}
	if stage != "" {
		query += ` AND m.stage = $` + strconv.Itoa(argIdx)
		args = append(args, stage)
		argIdx++
	}
	if dateStr != "" {
		date, err := time.Parse("2006-01-02", dateStr)
		if err == nil {
			query += ` AND m.scheduled_at >= $` + strconv.Itoa(argIdx) + ` AND m.scheduled_at < $` + strconv.Itoa(argIdx+1)
			args = append(args, date, date.Add(24*time.Hour))
			argIdx += 2
		}
	}
	query += ` ORDER BY m.scheduled_at ASC LIMIT 200`

	rows, err := h.db.Queryx(query, args...)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	defer rows.Close()

	matches, err := scanMatches(rows)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"matches": matches, "total": len(matches)})
}

func (h *MatchHandler) Get(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid id"})
		return
	}

	var match model.Match
	err = h.db.Get(&match, "SELECT * FROM matches WHERE id = $1", id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "match not found"})
		return
	}

	full := model.MatchFull{Match: match}

	if match.HomeTeamID != nil {
		var ht model.Team
		if err := h.db.Get(&ht, "SELECT * FROM teams WHERE id = $1", *match.HomeTeamID); err == nil {
			full.HomeTeam = &ht
		}
	}
	if match.AwayTeamID != nil {
		var at model.Team
		if err := h.db.Get(&at, "SELECT * FROM teams WHERE id = $1", *match.AwayTeamID); err == nil {
			full.AwayTeam = &at
		}
	}

	var streams []model.Stream
	h.db.Select(&streams, "SELECT * FROM streams WHERE match_id = $1 AND is_active = true ORDER BY priority DESC", id)
	full.Streams = streams

	c.JSON(http.StatusOK, full)
}

func (h *MatchHandler) Create(c *gin.Context) {
	var input struct {
		HomeTeamID  *int      `json:"home_team_id"`
		AwayTeamID  *int      `json:"away_team_id"`
		GroupID     *int      `json:"group_id"`
		Stage       string    `json:"stage" binding:"required"`
		Venue       string    `json:"venue"`
		City        string    `json:"city"`
		ScheduledAt time.Time `json:"scheduled_at" binding:"required"`
		Status      string    `json:"status"`
		HomeScore   int       `json:"home_score"`
		AwayScore   int       `json:"away_score"`
		Minute      *int      `json:"minute"`
	}
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	status := input.Status
	if status == "" {
		status = "scheduled"
	}
	query := `INSERT INTO matches (home_team_id, away_team_id, group_id, stage, venue, city, scheduled_at, status, home_score, away_score, minute)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11) RETURNING id`
	var id int
	err := h.db.QueryRow(query,
		input.HomeTeamID, input.AwayTeamID, input.GroupID,
		input.Stage, input.Venue, input.City, input.ScheduledAt, status, input.HomeScore, input.AwayScore, input.Minute,
	).Scan(&id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusCreated, gin.H{"id": id})
}

func (h *MatchHandler) Update(c *gin.Context) {
	id, _ := strconv.Atoi(c.Param("id"))
	var input struct {
		HomeTeamID  *int      `json:"home_team_id"`
		AwayTeamID  *int      `json:"away_team_id"`
		GroupID     *int      `json:"group_id"`
		Stage       string    `json:"stage" binding:"required"`
		Venue       string    `json:"venue"`
		City        string    `json:"city"`
		ScheduledAt time.Time `json:"scheduled_at" binding:"required"`
		Status      string    `json:"status" binding:"required"`
		HomeScore   int       `json:"home_score"`
		AwayScore   int       `json:"away_score"`
		Minute      *int      `json:"minute"`
	}
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	_, err := h.db.Exec(
		`UPDATE matches SET home_team_id=$1, away_team_id=$2, group_id=$3, stage=$4, venue=$5, city=$6, scheduled_at=$7, status=$8, home_score=$9, away_score=$10, minute=$11, updated_at=NOW() WHERE id=$12`,
		input.HomeTeamID, input.AwayTeamID, input.GroupID, input.Stage, input.Venue, input.City, input.ScheduledAt, input.Status, input.HomeScore, input.AwayScore, input.Minute, id,
	)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"ok": true})
}

// UpdateScore handles PUT /api/admin/matches/:id/score — updates score and broadcasts via WebSocket.
func (h *MatchHandler) UpdateScore(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid id"})
		return
	}
	var req struct {
		HomeScore int    `json:"home_score"`
		AwayScore int    `json:"away_score"`
		Status    string `json:"status"`
		Minute    *int   `json:"minute"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	_, err = h.db.Exec(
		`UPDATE matches SET home_score=$1, away_score=$2, status=$3, minute=$4, updated_at=NOW() WHERE id=$5`,
		req.HomeScore, req.AwayScore, req.Status, req.Minute, id,
	)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	if h.hub != nil {
		minute := 0
		if req.Minute != nil {
			minute = *req.Minute
		}
		h.hub.PublishScore(wslib.ScoreUpdate{
			MatchID:   id,
			HomeScore: req.HomeScore,
			AwayScore: req.AwayScore,
			Minute:    minute,
			Status:    req.Status,
		})
	}
	c.JSON(http.StatusOK, gin.H{"ok": true})
}

func (h *MatchHandler) Delete(c *gin.Context) {
	id, _ := strconv.Atoi(c.Param("id"))
	h.db.Exec("DELETE FROM matches WHERE id = $1", id)
	c.JSON(http.StatusOK, gin.H{"ok": true})
}

// scanMatches is a simplified scanner for the joined query
func scanMatches(rows interface {
	StructScan(interface{}) error
	Next() bool
	Close() error
}) ([]map[string]interface{}, error) {
	var results []map[string]interface{}
	// For simplicity, use a raw query approach with manual struct
	type row struct {
		model.Match
		HomeTeamID2    *int    `db:"home_team.id"`
		HomeTeamCode   *string `db:"home_team.code"`
		HomeTeamNameEN *string `db:"home_team.name_en"`
		HomeTeamNameRU *string `db:"home_team.name_ru"`
		HomeTeamFlag   *string `db:"home_team.flag_url"`
		AwayTeamID2    *int    `db:"away_team.id"`
		AwayTeamCode   *string `db:"away_team.code"`
		AwayTeamNameEN *string `db:"away_team.name_en"`
		AwayTeamNameRU *string `db:"away_team.name_ru"`
		AwayTeamFlag   *string `db:"away_team.flag_url"`
	}
	for rows.Next() {
		var r row
		if err := rows.StructScan(&r); err != nil {
			return nil, err
		}
		m := map[string]interface{}{
			"id": r.Match.ID, "home_team_id": r.Match.HomeTeamID, "away_team_id": r.Match.AwayTeamID, "group_id": r.Match.GroupID,
			"stage": r.Match.Stage, "venue": r.Match.Venue,
			"city": r.Match.City, "scheduled_at": r.Match.ScheduledAt,
			"status": r.Match.Status, "home_score": r.Match.HomeScore,
			"away_score": r.Match.AwayScore, "minute": r.Match.Minute,
		}
		if r.HomeTeamID2 != nil {
			m["home_team"] = map[string]interface{}{
				"id": r.HomeTeamID2, "code": r.HomeTeamCode,
				"name_en": r.HomeTeamNameEN, "name_ru": r.HomeTeamNameRU,
				"flag_url": r.HomeTeamFlag,
			}
		}
		if r.AwayTeamID2 != nil {
			m["away_team"] = map[string]interface{}{
				"id": r.AwayTeamID2, "code": r.AwayTeamCode,
				"name_en": r.AwayTeamNameEN, "name_ru": r.AwayTeamNameRU,
				"flag_url": r.AwayTeamFlag,
			}
		}
		results = append(results, m)
	}
	return results, nil
}
