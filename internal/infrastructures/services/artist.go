package services

import (
	"database/sql"
	"fmt"
	"log"
	"recipe-scraper/internal/domains/models"
	"recipe-scraper/internal/domains/repositories"
)

type artistService struct {
	DB *sql.DB
}

func (a artistService) Save(ID int, name string) (*models.Artist, error) {
	log.Println("start save artist")
	defer log.Println("end save artist")
	if err := a.DB.Ping(); err != nil {
		return nil, fmt.Errorf("failed to ping database: %w", err)
	}
	sql := `
INSERT INTO artist(id, name) 
VALUES ($1, $2) 
ON CONFLICT (id) 
DO UPDATE SET name = $2 RETURNING id
`
	if err := a.DB.QueryRow(sql, ID, name).Scan(&ID); err != nil {
		return nil, fmt.Errorf("failed to insert artist: %w", err)
	}
	return &models.Artist{ID: ID, Name: name}, nil
}

func NewArtistService(db *sql.DB) repositories.ArtistRepository {
	return &artistService{DB: db}
}
