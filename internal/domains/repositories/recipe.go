package repositories

import (
	"recipe-scraper/internal/domains/models"
)

type RecipeRepository interface {
	Save(id, cookingTimeInMinutes, artistId int, name string) (*models.Recipe, error)
	List(input *RecipeListInput) ([]*models.Recipe, error)
}
type artistInput struct {
	Name string
	ID   string
}
type recipeInput struct {
	Name           string
	CookingTimeMax int
}
type RecipeListInput struct {
	Artist artistInput
	Recipe recipeInput
}
