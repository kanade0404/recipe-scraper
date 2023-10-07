package artist

import (
	"database/sql"
	"fmt"
	"recipe-scraper/internal/domains/models"
	"recipe-scraper/internal/domains/repositories"
	"recipe-scraper/internal/logger"
)

type artistService struct {
	DB *sql.DB
}

func (a artistService) Save(ID int, name string) (*models.Artist, error) {
	logger.Info("start save artist")
	defer logger.Info("end save artist")
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
