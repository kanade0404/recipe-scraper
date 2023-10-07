package usecases

import (
	"context"
	"errors"
	"fmt"
	"github.com/chromedp/cdproto/cdp"
	"github.com/chromedp/chromedp"
	"net/url"
	"recipe-scraper/internal/domains/models"
	"recipe-scraper/internal/domains/repositories"
	initialize "recipe-scraper/internal/infrastructures/chromedp/initialize"
	"recipe-scraper/internal/infrastructures/chromedp/run"
	"recipe-scraper/internal/logger"
	"regexp"
	"strconv"
	"strings"
)

type RecipeUseCase interface {
	ScrapingRecipe(ctx context.Context, recipeURL string) (*models.Recipe, error)
	ScrapingPickupRecipes(ctx context.Context) ([]*models.Recipe, error)
	ListRecipe(ctx context.Context) ([]*models.Recipe, error)
}
type recipeUserCase struct {
	recipeRepo repositories.RecipeRepository
	artistRepo repositories.ArtistRepository
}

// ScrapingPickupRecipes
/*
ナディアのおすすめレシピ一覧からレシピ情報を取得して返す
*/
func (ru *recipeUserCase) ScrapingPickupRecipes(ctx context.Context) ([]*models.Recipe, error) {
	logger.Info("start ScrapingPickupRecipes")
	defer logger.Info("end ScrapingPickupRecipes")
	const URL = "https://oceans-nadia.com/pickup-recipes"
	ctx, timeoutCancel, allocatorCancel, contextCancel := initialize.Nadia(ctx)
	defer timeoutCancel()
	defer allocatorCancel()
	defer contextCancel()
	var (
		recipeCards []*cdp.Node
		recipes     []*models.Recipe
	)
	// お勧めレシピのカードコンポーネントの要素を一括取得
	if err := run.Nadia(ctx, chromedp.Tasks{
		chromedp.Navigate(URL),
		chromedp.WaitVisible("#__next > div > div > div:nth-child(2)"),
		// querySelectorAllなのでByQueryAll
		chromedp.Nodes("#__next > div > div > div > div > ul:nth-child(5) > li", &recipeCards, chromedp.ByQueryAll),
	}); err != nil {
		return nil, fmt.Errorf("failed to run chromedp. %w", err)
	}
	for i := range recipeCards {
		var (
			cookingTimeMinutesText string
			cookingTimeMinutes     int
			recipeName             string
			recipeLink             []*cdp.Node
			recipeID               int
			artistID               int
		)
		/*
			chromedp.FromNode(recipeCards[i])でrecipeCards[i]からselectorを辿ってくれる。
			chromedp.ByQueryをつけないとrecipeCards[0]から全部辿ってしまう
		*/
		if err := run.Nadia(ctx, chromedp.Tasks{
			chromedp.Text("div > a > div > div > span", &cookingTimeMinutesText, chromedp.FromNode(recipeCards[i]), chromedp.ByQuery),
			chromedp.Text("div > a > div > h2", &recipeName, chromedp.FromNode(recipeCards[i]), chromedp.ByQuery),
			chromedp.Nodes("div > a", &recipeLink, chromedp.FromNode(recipeCards[i]), chromedp.ByQuery),
		}); err != nil {
			return nil, fmt.Errorf("failed to run chromedp. %w", err)
		}
		// 調理時間を取得する
		re, err := regexp.Compile("[0-9]+")
		if err != nil {
			return nil, fmt.Errorf("failed to compile regexp: %w", err)
		}
		if cookingTimeMinutes, err = strconv.Atoi(re.FindString(cookingTimeMinutesText)); err != nil {
			return nil, fmt.Errorf("failed to convert string to int: %w", err)
		}
		// artistIDをlinkから取得する
		if len(recipeLink) == 0 {
			return nil, errors.New("recipe link not found")
		}
		href, ok := recipeLink[0].Attribute("href")
		if !ok {
			return nil, errors.New("href not found")
		}
		u, err := url.Parse(href)
		if err != nil {
			return nil, fmt.Errorf("failed to parse link. %w", err)
		}
		paths := strings.Split(u.Path, "/")
		for pi := range paths {
			if paths[pi] == "user" {
				aID, err := strconv.Atoi(paths[pi+1])
				if err != nil {
					return nil, fmt.Errorf("failed to parse artistID. %w", err)
				}
				artistID = aID
			}
			if paths[pi] == "recipe" {
				rID, err := strconv.Atoi(paths[pi+1])
				if err != nil {
					return nil, fmt.Errorf("failed to parse recipeID. %w", err)
				}
				recipeID = rID
			}
		}
		recipes = append(recipes, &models.Recipe{
			ID:                   recipeID,
			Name:                 recipeName,
			ArtistID:             artistID,
			CookingTimeInMinutes: cookingTimeMinutes,
		})
	}
	return recipes, nil
}

// ListRecipe
/*
保存済みのナディアのレシピ一覧を返す
*/
func (ru *recipeUserCase) ListRecipe(ctx context.Context) ([]*models.Recipe, error) {
	logger.Info("start ListRecipe")
	defer logger.Info("end ListRecipe")
	recipes, err := ru.recipeRepo.List(nil)
	if err != nil {
		return nil, fmt.Errorf("failed to select recipe: %w", err)
	}
	return recipes, nil
}

// NewRecipeUseCase
/*
ナディアのusecaseを初期化する。
依存先であるrepositoryが必要になる。
*/
func NewRecipeUseCase(rr repositories.RecipeRepository, ar repositories.ArtistRepository) RecipeUseCase {
	return &recipeUserCase{recipeRepo: rr, artistRepo: ar}
}

// ScrapingRecipe
/*
渡されたナディアのレシピのURLからレシピ情報を取得して保存した後に返す
*/
func (ru *recipeUserCase) ScrapingRecipe(ctx context.Context, recipeURL string) (*models.Recipe, error) {
	logger.Info("start ScrapingRecipe")
	defer logger.Info("end ScrapingRecipe")
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
		chromedp.Text(`#__next > div > div > div > div:nth-child(1) > div > h1`, &recipeName, chromedp.ByQuery),
		// css selectorの要素を全て取得する（ここではdocument.querySelectorのように1つの要素しか取得しない）
		chromedp.Text(`#__next > div > div > div > div:nth-child(1) > section:nth-child(3) > div:nth-child(1) > ul > li:nth-child(3) > h2 > span:nth-child(2)`, &textCookingTimeMinutes, chromedp.ByQuery),
		chromedp.Text(`#__next > div > div > div > div > div > div > a > div > p`, &artistName, chromedp.ByQuery),
	}); err != nil {
		return nil, fmt.Errorf("failed to run chromedp: %w", err)
	}
	recipe := &models.Recipe{ID: recipeID, Name: recipeName, ArtistID: artistID}
	if textCookingTimeMinutes != "" {
		// "[0-9]+分"から[0-9]+を抜き出す
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
	if _, err := ru.recipeRepo.Save(recipeID, cookingTimeMinutes, artistID, recipeName); err != nil {
		return nil, fmt.Errorf("failed to save recipe: %w", err)
	}
	return recipe, nil
}
