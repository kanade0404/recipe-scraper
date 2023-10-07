package repositories

import (
	"recipe-scraper/internal/domains/models"
)

// RecipeRepository
/*
レシピ情報の読み取りと書き込みを行います。
実体はinternal/infrastructures/servicesにあります。
*/
type RecipeRepository interface {
	// Save レシピ情報の登録をします。すでに存在する場合は更新をします。
	Save(id, cookingTimeInMinutes, artistId int, name string) (*models.Recipe, error)
	// List すでに登録されているレシピ一覧を取得します。
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
