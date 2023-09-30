package handlers

import (
	"github.com/labstack/echo/v4"
	"recipe-scraper/internal/usecases"
)

type NadiaHandler interface {
	HandleScraping(c echo.Context) error
	List(c echo.Context) error
}
type nadiaHandler struct {
	recipeUsecase usecases.RecipeUseCase
}

func NewNadiaHandler(useCase usecases.RecipeUseCase) NadiaHandler {
	return &nadiaHandler{recipeUsecase: useCase}
}
func (h nadiaHandler) HandleScraping(c echo.Context) error {
	url := c.QueryParam("url")
	if url == "" {
		return echo.NewHTTPError(400, "url is required")
	}
	recipe, err := h.recipeUsecase.ScrapingRecipe(c.Request().Context(), url)
	if err != nil {
		return echo.NewHTTPError(500, err.Error())
	}
	return c.JSON(200, recipe)
}

func (h nadiaHandler) List(c echo.Context) error {
	recipes, err := h.recipeUsecase.ListRecipe(c.Request().Context())
	if err != nil {
		return echo.NewHTTPError(500, err.Error())
	}
	return c.JSON(200, recipes)
}
