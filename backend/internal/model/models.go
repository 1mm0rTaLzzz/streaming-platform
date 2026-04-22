package model

import "time"

type Group struct {
	ID   int    `db:"id" json:"id"`
	Name string `db:"name" json:"name"`
}

type Team struct {
	ID         int    `db:"id" json:"id"`
	ExternalID *int   `db:"external_id" json:"external_id,omitempty"`
	Code       string `db:"code" json:"code"`
	NameEN     string `db:"name_en" json:"name_en"`
	NameRU     string `db:"name_ru" json:"name_ru"`
	FlagURL    string `db:"flag_url" json:"flag_url"`
	GroupID    *int   `db:"group_id" json:"group_id,omitempty"`
}

type Match struct {
	ID          int       `db:"id" json:"id"`
	ExternalID  *int      `db:"external_id" json:"external_id,omitempty"`
	HomeTeamID  *int      `db:"home_team_id" json:"home_team_id,omitempty"`
	AwayTeamID  *int      `db:"away_team_id" json:"away_team_id,omitempty"`
	GroupID     *int      `db:"group_id" json:"group_id,omitempty"`
	Stage       string    `db:"stage" json:"stage"`
	Venue       string    `db:"venue" json:"venue"`
	City        string    `db:"city" json:"city"`
	ScheduledAt time.Time `db:"scheduled_at" json:"scheduled_at"`
	Status      string    `db:"status" json:"status"`
	HomeScore   int       `db:"home_score" json:"home_score"`
	AwayScore   int       `db:"away_score" json:"away_score"`
	Minute      *int      `db:"minute" json:"minute,omitempty"`
	CreatedAt   time.Time `db:"created_at" json:"created_at"`
	UpdatedAt   time.Time `db:"updated_at" json:"updated_at"`
}

type MatchFull struct {
	Match
	HomeTeam *Team    `json:"home_team,omitempty"`
	AwayTeam *Team    `json:"away_team,omitempty"`
	Streams  []Stream `json:"streams,omitempty"`
}

type Stream struct {
	ID             int       `db:"id" json:"id"`
	MatchID        int       `db:"match_id" json:"match_id"`
	URL            string    `db:"url" json:"url"`
	SourceType     string    `db:"source_type" json:"source_type"`
	Label          string    `db:"label" json:"label"`
	LanguageCode   string    `db:"language_code" json:"language_code"`
	Region         string    `db:"region" json:"region"`
	CommentaryType string    `db:"commentary_type" json:"commentary_type"`
	Quality        string    `db:"quality" json:"quality"`
	IsActive       bool      `db:"is_active" json:"is_active"`
	Priority       int       `db:"priority" json:"priority"`
	CreatedAt      time.Time `db:"created_at" json:"created_at"`
}

type Mirror struct {
	ID            int        `db:"id" json:"id"`
	Domain        string     `db:"domain" json:"domain"`
	IsActive      bool       `db:"is_active" json:"is_active"`
	IsPrimary     bool       `db:"is_primary" json:"is_primary"`
	LastCheckedAt *time.Time `db:"last_health_check" json:"last_checked_at,omitempty"`
	HealthStatus  string     `db:"health_status" json:"health_status"`
	Region        string     `db:"region" json:"region"`
	Priority      int        `db:"priority" json:"priority"`
	CreatedAt     time.Time  `db:"created_at" json:"created_at"`
}

type AdminUser struct {
	ID           int       `db:"id" json:"id"`
	Username     string    `db:"username" json:"username"`
	PasswordHash string    `db:"password_hash" json:"-"`
	Role         string    `db:"role" json:"role"`
	CreatedAt    time.Time `db:"created_at" json:"created_at"`
}

type Standing struct {
	Team    Team `json:"team"`
	GroupID int  `json:"group_id"`
	Played  int  `json:"played"`
	Won     int  `json:"won"`
	Drawn   int  `json:"drawn"`
	Lost    int  `json:"lost"`
	GF      int  `json:"gf"`
	GA      int  `json:"ga"`
	GD      int  `json:"gd"`
	Points  int  `json:"points"`
}
