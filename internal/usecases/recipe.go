package usecases

import (
	"context"
	"fmt"
	"github.com/chromedp/chromedp"
	"log"
	"net/url"
	"recipe-scraper/internal/domains/models"
	"recipe-scraper/internal/domains/repositories"
	initialize "recipe-scraper/internal/infrastructures/chromedp/initialize"
	"recipe-scraper/internal/infrastructures/chromedp/run"
	"regexp"
	"strconv"
	"strings"
)

type RecipeUseCase interface {
	ScrapingRecipe(ctx context.Context, recipeURL string) (*models.Recipe, error)
	ListRecipe(ctx context.Context) ([]*models.Recipe, error)
}
type recipeUserCase struct {
	recipeRepo repositories.RecipeRepository
	artistRepo repositories.ArtistRepository
}

func (ru *recipeUserCase) ListRecipe(ctx context.Context) ([]*models.Recipe, error) {
	log.Println("start ListRecipe")
	defer log.Println("end ListRecipe")
	recipes, err := ru.recipeRepo.List(nil)
	if err != nil {
		return nil, fmt.Errorf("failed to select recipe: %w", err)
	}
	return recipes, nil
}

func NewRecipeUseCase(rr repositories.RecipeRepository, ar repositories.ArtistRepository) RecipeUseCase {
	return &recipeUserCase{recipeRepo: rr, artistRepo: ar}
}

func (ru *recipeUserCase) ScrapingRecipe(ctx context.Context, recipeURL string) (*models.Recipe, error) {
	log.Println("start ScrapingRecipe")
	defer log.Println("end ScrapingRecipe")
	ctx, timeoutCancel, allocatorCancel, contextCancel := initialize.Nadia(ctx)
	defer timeoutCancel()
	defer allocatorCancel()
	defer contextCancel()
	// URLを正規表現でチェック
	re, err := regexp.Compile("https://oceans-nadia.com/user/([0-9]+)/recipe/([0-9]+)")
	if err != nil {
		return nil, fmt.Errorf("failed to compile regexp: %w", err)
	}
	if isMatch := re.MatchString(recipeURL); !isMatch {
		return nil, fmt.Errorf("invalid recipeURL: %s", recipeURL)
	}
	var (
		recipeName             string
		textCookingTimeMinutes string
		cookingTimeMinutes     int
		recipeID               int
		artistID               int
		artistName             string
	)
	u, err := url.Parse(recipeURL)
	if err != nil {
		return nil, fmt.Errorf("failed to parse url: %w", err)
	}
	paths := strings.Split(u.Path, "/")
	for i := 0; i < len(paths); i++ {
		if paths[i] == "recipe" {
			if recipeID, err = strconv.Atoi(paths[i+1]); err != nil {
				return nil, fmt.Errorf("failed to convert string to int: %w", err)
			}
		}
		if paths[i] == "user" {
			if artistID, err = strconv.Atoi(paths[i+1]); err != nil {
				return nil, fmt.Errorf("failed to convert string to int: %w", err)
			}
		}
	}
	if recipeID == 0 {
		return nil, fmt.Errorf("failed to get recipeID: %w", err)
	}
	if artistID == 0 {
		return nil, fmt.Errorf("failed to get artistID: %w", err)
	}
	if err := run.Nadia(ctx, chromedp.Tasks{
		// urlにアクセスする
		chromedp.Navigate(recipeURL),
		// css selectorの要素が表示されるまで待つ
		chromedp.WaitVisible(`#__next > div > div > div > div:nth-child(1) > div > h1`),
		// css selectorの要素のtextContentを取得する
		chromedp.Text(`#__next > div > div > div > div:nth-child(1) > div > h1`, &recipeName),
		// css selectorの要素を全て取得する（ここではdocument.querySelectorのように1つの要素しか取得しない）
		chromedp.Text(`#__next > div > div > div > div:nth-child(1) > section:nth-child(3) > div:nth-child(1) > ul > li:nth-child(3) > h2 > span:nth-child(2)`, &textCookingTimeMinutes),
		chromedp.Text(`#__next > div > div > div > div > div > div > a > div > p`, &artistName),
	}); err != nil {
		return nil, fmt.Errorf("failed to run chromedp: %w", err)
	}
	recipe := &models.Recipe{ID: recipeID, Name: recipeName, ArtistID: artistID}
	if textCookingTimeMinutes != "" {
		// [0-9]+分→[0-9]+
		re, err := regexp.Compile("[0-9]+")
		if err != nil {
			return nil, fmt.Errorf("failed to compile regexp: %w", err)
		}
		if cookingTimeMinutes, err = strconv.Atoi(re.FindString(textCookingTimeMinutes)); err != nil {
			return nil, fmt.Errorf("failed to convert string to int: %w", err)
		}
	}
	if _, err := ru.artistRepo.Save(artistID, artistName); err != nil {
		return nil, fmt.Errorf("failed to save artist: %w", err)
	}
	if _, err := ru.recipeRepo.Save(recipeID, cookingTimeMinutes, artistID, artistName); err != nil {
		return nil, fmt.Errorf("failed to save recipe: %w", err)
	}
	return recipe, nil
}
