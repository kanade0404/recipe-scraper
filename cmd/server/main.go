package main

import (
	"github.com/labstack/echo/v4"
	"recipe-scraper/internal/config"
	"recipe-scraper/internal/handlers"
	"recipe-scraper/internal/infrastructures/services"
	"recipe-scraper/internal/usecases"
)

func main() {
	e := echo.New()
	e.Logger.Fatal(run(e))
}

func run(e *echo.Echo) error {
	c, err := config.NewConfig()
	if err != nil {
		return err
	}
	recipeService := services.NewRecipeService(c.DB)
	artistService := services.NewArtistService(c.DB)
	recipeUsecase := usecases.NewRecipeUseCase(recipeService, artistService)
	nadiaHandler := handlers.NewNadiaHandler(recipeUsecase)
	api := e.Group("/api")
	v1 := api.Group("/v1")
	// /api/v1/scraping/nadia
	v1.GET("/scraping/nadia", nadiaHandler.HandleScraping)
	v1.GET("/recipe", nadiaHandler.List)
	return e.Start(":1323")
}
