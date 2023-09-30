package models

import (
	"time"
)

type Recipe struct {
	ID                   int       `db:"id" json:"id"`
	Name                 string    `db:"name" json:"name"`
	ArtistID             int       `db:"artist_id" json:"artist_id"`
	CookingTimeInMinutes int       `db:"cooking_time_in_minutes" json:"cooking_time_in_minutes"`
	CreatedAt            time.Time `db:"created_at" json:"created_at"`
	UpdatedAt            time.Time `db:"updated_at" json:"updated_at"`
}
