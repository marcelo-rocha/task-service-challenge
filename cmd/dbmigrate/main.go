package main

import (
	"database/sql"
	"flag"
	"log"
	"os"
	"time"

	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/mysql"
	_ "github.com/golang-migrate/migrate/v4/source/file"
)

var (
	flagaddr       = flag.String("addr", "mysql://root:secret7@(localhost:3306)/test?multiStatements=true&parseTime=true", "mysql connection")
	flagup         = flag.Bool("up", false, "run up migrations")
	flagdown       = flag.Bool("down", false, "run down migrations")
	flagmigrations = flag.String("migrations", "file://persistence/migration", "migrations path")
)

func main() {
	flag.Parse()

	if *flagup && *flagdown {
		log.Println("error: should use either -up or -down")
		os.Exit(1)
	}

	db, err := Connect(*flagaddr)
	if err != nil {
		log.Println("db connection error", err)
		os.Exit(2)
	}
	db.Close()

	migrator, err := migrate.New(*flagmigrations, *flagaddr)
	if err != nil {
		log.Println("migrate init error", err)
		os.Exit(2)
	}

	if *flagup {
		if err = migrator.Up(); err != nil {
			log.Println("migrate up error", err)
			os.Exit(3)
		}
	} else if *flagdown {
		if err = migrator.Down(); err != nil {
			log.Println("migrate down error", err)
			os.Exit(3)
		}
	}

}

func Connect(url string) (*sql.DB, error) {
	db, err := sql.Open("mysql", url)
	if err != nil {
		return nil, err
	}

	db.SetConnMaxLifetime(time.Minute * 5)
	db.SetMaxOpenConns(10)
	db.SetMaxIdleConns(10)

	return db, nil
}
