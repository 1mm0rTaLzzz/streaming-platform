package handler

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/jmoiron/sqlx"
	"stream_site/backend/internal/model"
)

type GroupHandler struct {
	db *sqlx.DB
}

func NewGroupHandler(db *sqlx.DB) *GroupHandler {
	return &GroupHandler{db: db}
}

func (h *GroupHandler) List(c *gin.Context) {
	var groups []model.Group
	if err := h.db.Select(&groups, "SELECT * FROM groups ORDER BY name"); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	type GroupWithTeams struct {
		model.Group
		Teams []TeamStanding `json:"teams"`
	}

	result := make([]GroupWithTeams, 0, len(groups))
	for _, g := range groups {
		teams := h.getStandings(g.ID)
		result = append(result, GroupWithTeams{Group: g, Teams: teams})
	}
	c.JSON(http.StatusOK, gin.H{"groups": result})
}

func (h *GroupHandler) Get(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid id"})
		return
	}
	var g model.Group
	if err := h.db.Get(&g, "SELECT * FROM groups WHERE id = $1", id); err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "group not found"})
		return
	}
	teams := h.getStandings(id)
	c.JSON(http.StatusOK, gin.H{"group": g, "standings": teams})
}

type TeamStanding struct {
	model.Team
	Played int `json:"played"`
	Won    int `json:"won"`
	Drawn  int `json:"drawn"`
	Lost   int `json:"lost"`
	GF     int `json:"gf"`
	GA     int `json:"ga"`
	GD     int `json:"gd"`
	Points int `json:"points"`
}

func (h *GroupHandler) getStandings(groupID int) []TeamStanding {
	// Calculate standings from finished matches
	query := `
		SELECT t.*,
			COUNT(m.id) as played,
			SUM(CASE
				WHEN m.home_team_id = t.id AND m.home_score > m.away_score THEN 1
				WHEN m.away_team_id = t.id AND m.away_score > m.home_score THEN 1
				ELSE 0 END) as won,
			SUM(CASE WHEN m.home_score = m.away_score THEN 1 ELSE 0 END) as drawn,
			SUM(CASE
				WHEN m.home_team_id = t.id AND m.home_score < m.away_score THEN 1
				WHEN m.away_team_id = t.id AND m.away_score < m.home_score THEN 1
				ELSE 0 END) as lost,
			SUM(CASE WHEN m.home_team_id = t.id THEN m.home_score WHEN m.away_team_id = t.id THEN m.away_score ELSE 0 END) as gf,
			SUM(CASE WHEN m.home_team_id = t.id THEN m.away_score WHEN m.away_team_id = t.id THEN m.home_score ELSE 0 END) as ga
		FROM teams t
		LEFT JOIN matches m ON (m.home_team_id = t.id OR m.away_team_id = t.id)
			AND m.group_id = $1 AND m.status = 'finished'
		WHERE t.group_id = $1
		GROUP BY t.id
		ORDER BY
			(SUM(CASE WHEN m.home_team_id = t.id AND m.home_score > m.away_score THEN 3
				WHEN m.away_team_id = t.id AND m.away_score > m.home_score THEN 3
				WHEN m.home_score = m.away_score AND m.id IS NOT NULL THEN 1
				ELSE 0 END)) DESC,
			(SUM(CASE WHEN m.home_team_id = t.id THEN m.home_score WHEN m.away_team_id = t.id THEN m.away_score ELSE 0 END) -
			 SUM(CASE WHEN m.home_team_id = t.id THEN m.away_score WHEN m.away_team_id = t.id THEN m.home_score ELSE 0 END)) DESC`

	rows, err := h.db.Queryx(query, groupID)
	if err != nil {
		return nil
	}
	defer rows.Close()

	var standings []TeamStanding
	for rows.Next() {
		var ts struct {
			model.Team
			Played int `db:"played"`
			Won    int `db:"won"`
			Drawn  int `db:"drawn"`
			Lost   int `db:"lost"`
			GF     int `db:"gf"`
			GA     int `db:"ga"`
		}
		if err := rows.StructScan(&ts); err != nil {
			continue
		}
		gd := ts.GF - ts.GA
		points := ts.Won*3 + ts.Drawn
		standings = append(standings, TeamStanding{
			Team: ts.Team, Played: ts.Played, Won: ts.Won,
			Drawn: ts.Drawn, Lost: ts.Lost, GF: ts.GF, GA: ts.GA,
			GD: gd, Points: points,
		})
	}
	return standings
}
