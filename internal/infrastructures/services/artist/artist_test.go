package artist

import (
	"database/sql"
	"github.com/google/go-cmp/cmp"
	"github.com/google/go-cmp/cmp/cmpopts"
	"github.com/ory/dockertest/v3"
	"log"
	"os"
	"recipe-scraper/internal/domains/models"
	"recipe-scraper/internal/infrastructures/services/tests"
	"recipe-scraper/internal/logger"
	"testing"
)

var db *sql.DB

func TestMain(m *testing.M) {
	var (
		pool     *dockertest.Pool
		resource *dockertest.Resource
		err      error
		// tearDownを実行する前にlog.Fatalln(err)を実行するとtearDownが実行されないため
		exitCode = 1
	)
	if pool, resource, err = tests.SetUp(); err != nil {
		logger.Error(err.Error())
		os.Exit(exitCode)
	}
	defer func(pool *dockertest.Pool, resource *dockertest.Resource) {
		err := tests.TearDown(pool, resource)
		if err != nil {
			log.Fatalln(err)
		}
		os.Exit(exitCode)
	}(pool, resource)
	if db, err = tests.ConnectDatabase(pool); err != nil {
		logger.Error(err.Error())
		return
	}
	m.Run()
	exitCode = 0
}

func Test_artistService_Save(t *testing.T) {
	type args struct {
		ID   int
		name string
	}
	testCases := []struct {
		name    string
		args    args
		want    *models.Artist
		wantErr bool
	}{
		{
			name: "新規追加が成功すること",
			args: args{
				ID:   1,
				name: "test",
			},
			want: &models.Artist{
				ID:   1,
				Name: "test",
			},
		},
	}
	for _, tt := range testCases {
		t.Run(tt.name, func(t *testing.T) {
			a := NewArtistService(db)
			got, err := a.Save(tt.args.ID, tt.args.name)
			if (err != nil) != tt.wantErr {
				t.Errorf("Save() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if diff := cmp.Diff(got, tt.want, cmpopts.IgnoreFields(models.Artist{}, "CreatedAt", "UpdatedAt")); diff != "" {
				t.Errorf("(-got+want)\n%v", diff)
			}
		})
	}
}
