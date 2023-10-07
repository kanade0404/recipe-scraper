package tests

import (
	"database/sql"
	"fmt"
	_ "github.com/lib/pq"
	"github.com/ory/dockertest/v3"
	"github.com/ory/dockertest/v3/docker"
	"path/filepath"
	"recipe-scraper/internal/logger"
	"time"
)

const (
	user     = "postgres"
	password = "secret"
	dbName   = "test_recipe"
	port     = "5433"
	dialect  = "postgres"
	dsn      = "postgres://%s:%s@localhost:%s/%s?sslmode=disable"
)

func newPool() (*dockertest.Pool, error) {
	return dockertest.NewPool("")
}

func prepareContainer(pool *dockertest.Pool) (*dockertest.Resource, error) {
	path, err := filepath.Abs("../../../../database")
	if err != nil {
		return nil, err
	}
	resource, err := pool.RunWithOptions(&dockertest.RunOptions{
		Repository: "postgres",
		Tag:        "16",
		Env: []string{
			fmt.Sprintf("POSTGRES_DB=%s", dbName),
			fmt.Sprintf("POSTGRES_USER=%s", user),
			fmt.Sprintf("POSTGRES_PASSWORD=%s", password),
		},
		ExposedPorts: []string{"5432"},
		PortBindings: map[docker.Port][]docker.PortBinding{"5432": {{HostIP: "0.0.0.0", HostPort: port}}},
		Mounts: []string{
			fmt.Sprintf("%s/sql:/docker-entrypoint-initdb.d", path),
			fmt.Sprintf("%s/data/test:/var/lib/postgresql/data", path),
		},
	}, func(config *docker.HostConfig) {
		config.AutoRemove = true
		config.RestartPolicy = docker.RestartPolicy{
			Name: "no",
		}
	})
	if err != nil {
		return nil, err
	}
	return resource, nil
}

func createConnection() (*sql.DB, error) {
	db, err := sql.Open("postgres", fmt.Sprintf(dsn, user, password, port, dbName))
	if err != nil {
		return nil, err
	}
	if err := db.Ping(); err != nil {
		return nil, fmt.Errorf("failed to ping database: %w", err)
	}
	return db, nil
}

func ConnectDatabase(pool *dockertest.Pool) (*sql.DB, error) {
	var (
		db  *sql.DB
		err error
	)
	if err := pool.Retry(func() error {
		db, err = createConnection()
		return err
	}); err != nil {
		return nil, err
	}
	return db, nil
}

func TearDown(pool *dockertest.Pool, resource *dockertest.Resource) error {
	logger.Info("after all...")
	if err := pool.Purge(resource); err != nil {
		return err
	}
	return nil
}

func SetUp() (*dockertest.Pool, *dockertest.Resource, error) {
	logger.Info("before all...")
	pool, err := newPool()
	pool.MaxWait = time.Minute * 1
	if err != nil {
		return nil, nil, err
	}
	resource, err := prepareContainer(pool)
	return pool, resource, err
}
