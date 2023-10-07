package tests

import (
	"database/sql"
	"github.com/go-testfixtures/testfixtures/v3"
	"os"
)

func NewFixture(db *sql.DB, fileNames ...string) (*testfixtures.Loader, error) {
	fixturePaths := make([]string, len(fileNames))
	dir, err := os.Getwd()
	if err != nil {
		return nil, err
	}
	for i := range fileNames {
		fixturePaths[i] = dir + "/" + fileNames[i]
	}
	fixtures, err := testfixtures.New(
		testfixtures.Database(db),                      // You database connection
		testfixtures.Dialect("postgres"),               // Available: "postgresql", "timescaledb", "mysql", "mariadb", "sqlite" and "sqlserver"
		testfixtures.FilesMultiTables(fixturePaths...), // The directory containing the YAML files
	)
	if err != nil {
		return nil, err
	}
	if err := fixtures.Load(); err != nil {
		return nil, err
	}
	return fixtures, nil
}
