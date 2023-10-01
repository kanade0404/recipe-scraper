package handlers

import (
	"github.com/labstack/echo/v4"
	"net/http"
	"recipe-scraper/internal/usecases"
)

type NadiaHandler interface {
	HandleScraping(c echo.Context) error
	ListPickupRecipes(c echo.Context) error
	List(c echo.Context) error
}
type nadiaHandler struct {
	recipeUsecase usecases.RecipeUseCase
}

// ListPickupRecipes
/*
おすすめレシピをスクレイピングで取得する
*/
func (h nadiaHandler) ListPickupRecipes(c echo.Context) error {
	recipes, err := h.recipeUsecase.ScrapingPickupRecipes(c.Request().Context())
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}
	return c.JSON(http.StatusOK, recipes)
}

// NewNadiaHandler
/*
ナディアのhandlerを初期化する
初期化する際には依存しているusecaseが必要になる
*/
func NewNadiaHandler(useCase usecases.RecipeUseCase) NadiaHandler {
	return &nadiaHandler{recipeUsecase: useCase}
}

// HandleScraping
/*
ナディアのレシピURLを渡してスクレイピングして取得する
*/
func (h nadiaHandler) HandleScraping(c echo.Context) error {
	url := c.QueryParam("url")
	if url == "" {
		return echo.NewHTTPError(http.StatusBadRequest, "url is required")
	}
	recipe, err := h.recipeUsecase.ScrapingRecipe(c.Request().Context(), url)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}
	return c.JSON(http.StatusOK, recipe)
}

// List
/*
登録済みのナディアのレシピ一覧を取得する
*/
func (h nadiaHandler) List(c echo.Context) error {
	recipes, err := h.recipeUsecase.ListRecipe(c.Request().Context())
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}
	return c.JSON(http.StatusOK, recipes)
}
