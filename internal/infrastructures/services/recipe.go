package services

import (
	"database/sql"
	"fmt"
	"log"
	"recipe-scraper/internal/domains/models"
	"recipe-scraper/internal/domains/repositories"
	"strings"
)

type recipeService struct {
	DB *sql.DB
}

type recipeDTO struct {
	ID                   int
	Name                 string
	ArtistID             int
	CookingTimeInMinutes int
}

func (r recipeService) Save(id, cookingTimeInMinutes, artistId int, name string) (*models.Recipe, error) {
	log.Printf("start save recipe: id: %d", id)
	defer log.Printf("end save recipe: %d", id)
	if err := r.DB.Ping(); err != nil {
		return nil, fmt.Errorf("failed to ping database: %w", err)
	}
	sql := `
INSERT INTO recipe(id, name, artist_id, cooking_time_in_minutes) 
VALUES ($1, $2, $3, $4) 
ON CONFLICT (id) 
DO UPDATE SET name = $2, artist_id = $3, cooking_time_in_minutes = $4 RETURNING id`
	if err := r.DB.QueryRow(sql, id, name, artistId, cookingTimeInMinutes).Scan(&id); err != nil {
		return nil, fmt.Errorf("failed to insert recipe: %w", err)
	}
	return &models.Recipe{
		ID:                   id,
		Name:                 name,
		CookingTimeInMinutes: cookingTimeInMinutes,
		ArtistID:             artistId,
	}, nil
}

func (r recipeService) List(input *repositories.RecipeListInput) ([]*models.Recipe, error) {
	log.Printf("start list recipe: %+v", input)
	defer log.Printf("end list recipe: %+v", input)
	if err := r.DB.Ping(); err != nil {
		return nil, fmt.Errorf("failed to ping database: %w", err)
	}
	var (
		whereClauses []string
		args         []interface{}
		whereSQL     string
	)
	if input != nil {
		if input.Artist.Name != "" {
			whereClauses = append(whereClauses, "artist.name LIKE $1")
			args = append(args, "%"+input.Artist.Name+"%")
		}
		if input.Artist.ID != "" {
			whereClauses = append(whereClauses, "artist.id = $2")
			args = append(args, input.Artist.ID)
		}
		if input.Recipe.Name != "" {
			whereClauses = append(whereClauses, "recipe.name LIKE $3")
			args = append(args, "%"+input.Recipe.Name+"%")
		}
		if input.Recipe.CookingTimeMax != 0 {
			whereClauses = append(whereClauses, "recipe.cooking_time_in_minutes <= $4")
			args = append(args, input.Recipe.CookingTimeMax)
		}
	}
	if len(whereClauses) > 0 {
		whereSQL = fmt.Sprintf("WHERE %s", strings.Join(whereClauses, " AND "))
	}
	rows, err := r.DB.Query(fmt.Sprintf("SELECT id, name, cooking_time_in_minutes, artist_id FROM recipe %s", whereSQL), args...)
	if err != nil {
		return nil, fmt.Errorf("failed to select recipe: %w", err)
	}
	defer func(rows *sql.Rows) {
		err := rows.Close()
		if err != nil {
			log.Printf("failed to close rows: %v", err)
		}
	}(rows)
	var recipes []*models.Recipe
	for rows.Next() {
		var recipe models.Recipe
		if err := rows.Scan(&recipe.ID, &recipe.Name, &recipe.CookingTimeInMinutes, &recipe.ArtistID); err != nil {
			return nil, fmt.Errorf("failed to scan recipe: %w", err)
		}
		recipes = append(recipes, &recipe)
	}
	return recipes, nil
}

func NewRecipeService(db *sql.DB) repositories.RecipeRepository {
	return &recipeService{db}
}
