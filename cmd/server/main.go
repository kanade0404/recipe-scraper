package main

import (
	"github.com/labstack/echo/v4"
	"recipe-scraper/internal/config"
	"recipe-scraper/internal/handlers"
	"recipe-scraper/internal/infrastructures/services/artist"
	"recipe-scraper/internal/infrastructures/services/recipe"
	"recipe-scraper/internal/logger"
	"recipe-scraper/internal/usecases"
)

// main アプリケーションのスターティングポイント
func main() {
	e := echo.New()
	e.Logger.Fatal(run(e))
}

// run
/*
アプリケーションの実行関数
ここでアプリケーションの実行に必要なusecaseやserviceを初期化し、echoサーバの初期化とサーバの起動を行います。
*/
func run(e *echo.Echo) error {
	c, err := config.NewConfig()
	if err != nil {
		return err
	}
	recipeService := recipe.NewRecipeService(c.DB)
	artistService := artist.NewArtistService(c.DB)
	recipeUsecase := usecases.NewRecipeUseCase(recipeService, artistService)
	nadiaHandler := handlers.NewNadiaHandler(recipeUsecase)
	api := e.Group("/api")
	v1 := api.Group("/v1")
	// /api/v1/scraping/nadia?url=<scraping_target_url>
	v1.GET("/scraping/nadia", nadiaHandler.HandleScraping)
	v1.GET("/scraping/nadia/pickup", nadiaHandler.ListPickupRecipes)
	// /api/v1/recipe
	v1.GET("/recipe", nadiaHandler.List)
	logger.InitLogger()
	return e.Start(":1323")
}
