package recipe

import (
	"database/sql"
	"fmt"
	"github.com/google/go-cmp/cmp"
	"github.com/google/go-cmp/cmp/cmpopts"
	"github.com/ory/dockertest/v3"
	"os"
	"recipe-scraper/internal/domains/models"
	"recipe-scraper/internal/infrastructures/services/tests"
	"recipe-scraper/internal/logger"
	"testing"
)

var (
	db *sql.DB
)

func TestMain(m *testing.M) {
	var (
		pool     *dockertest.Pool
		resource *dockertest.Resource
		err      error
		exitCode = 1
	)
	if pool, resource, err = tests.SetUp(); err != nil {
		logger.Error(err.Error())
		os.Exit(exitCode)
	}
	defer func(pool *dockertest.Pool, resource *dockertest.Resource) {
		err := tests.TearDown(pool, resource)
		if err != nil {
			logger.Error(err.Error())
		}
		os.Exit(exitCode)
	}(pool, resource)
	if db, err = tests.ConnectDatabase(pool); err != nil {
		logger.Error(err.Error())
		return
	}
	if _, err := tests.NewFixture(db, "fixture.yml"); err != nil {
		logger.Error(fmt.Sprintf("failed to load fixture: %v", err))
		return
	}

	m.Run()
	exitCode = 0
}

func Test_recipeService_Save(t *testing.T) {
	type args struct {
		id                   int
		cookingTimeInMinutes int
		artistId             int
		name                 string
	}
	testCases := []struct {
		name    string
		args    args
		want    *models.Recipe
		wantErr bool
	}{
		{
			name: "新規追加が成功すること",
			args: args{
				id:                   3,
				cookingTimeInMinutes: 1,
				artistId:             1,
				name:                 "Recipe 4",
			},
			want: &models.Recipe{
				ID:                   3,
				Name:                 "Recipe 4",
				CookingTimeInMinutes: 1,
				ArtistID:             1,
			},
		},
		{
			name: "存在しないartistIdのため外部キーエラーになること",
			args: args{
				id:                   4,
				cookingTimeInMinutes: 1,
				artistId:             3,
				name:                 "Recipe 1'",
			},
			wantErr: true,
		},
		{
			name: "存在するidのため更新が成功すること",
			args: args{
				id:                   1,
				cookingTimeInMinutes: 1,
				artistId:             1,
				name:                 "Recipe 1'",
			},
			want: &models.Recipe{
				ID:                   1,
				Name:                 "Recipe 1'",
				CookingTimeInMinutes: 1,
				ArtistID:             1,
			},
		},
	}
	for _, tt := range testCases {
		t.Run(tt.name, func(t *testing.T) {
			r := NewRecipeService(db)
			got, err := r.Save(tt.args.id, tt.args.cookingTimeInMinutes, tt.args.artistId, tt.args.name)
			if (err != nil) != tt.wantErr {
				t.Errorf("Save() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if err != nil && tt.wantErr {
				return
			}
			if diff := cmp.Diff(got, tt.want, cmpopts.IgnoreFields(models.Recipe{}, "CreatedAt", "UpdatedAt")); diff != "" {
				t.Errorf("(-got+want)\n%v", diff)
			}
		})
	}
}
