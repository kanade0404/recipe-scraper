package repositories

import (
	"recipe-scraper/internal/domains/models"
)

type ArtistRepository interface {
	Save(ID int, name string) (*models.Artist, error)
}
