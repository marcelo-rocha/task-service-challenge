package persistence

import (
	"context"
	"fmt"
	"os"
	"testing"

	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/mysql"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"github.com/ory/dockertest/v3"
	"go.uber.org/zap"

	_ "github.com/golang-migrate/migrate/v4/source/file"
)

var db *Connection
var logger *zap.Logger

const dbName = "test"
const DemoUserId = 2

func TestMain(m *testing.M) {
	logger, _ = zap.NewDevelopment()
	// uses a sensible default on windows (tcp/http) and linux/osx (socket)
	pool, err := dockertest.NewPool("")
	if err != nil {
		logger.Fatal("Could not connect to docker", zap.Error(err))
	}

	// pulls an image, creates a container based on it and runs it
	resource, err := pool.Run("mysql", "8.0", []string{"MYSQL_ROOT_PASSWORD=secret7", "MYSQL_DATABASE=" + dbName})
	if err != nil {
		logger.Fatal("Could not start resource", zap.Error(err))
	}

	// exponential backoff-retry, because the application in the container might not be ready to accept connections yet
	if err := pool.Retry(func() error {
		var err error
		db, err = NewConnection(context.Background(), fmt.Sprintf("root:secret7@(localhost:%s)/%s?multiStatements=true&parseTime=true",
			resource.GetPort("3306/tcp"), dbName))
		if err != nil {
			return err
		}
		return db.Driver.Ping()
	}); err != nil {
		logger.Fatal("Could not connect to database", zap.Error(err))
	}

	if err := runMigrations(); err != nil {
		logger.Fatal("Could not run migrate", zap.Error(err))
	}

	code := m.Run()

	db.Close()

	// You can't defer this because os.Exit doesn't care for defer
	if err := pool.Purge(resource); err != nil {
		logger.Fatal("Could not purge resource", zap.Error(err))
	}

	os.Exit(code)
}

func runMigrations() error {
	wordDir, _ := os.Getwd()
	driver, _ := mysql.WithInstance(db.Driver, &mysql.Config{})
	migrate, err := migrate.NewWithDatabaseInstance(
		fmt.Sprintf("file:///%s/migration", wordDir),
		dbName, driver)
	if err != nil {
		return err
	}

	return migrate.Up()

}
