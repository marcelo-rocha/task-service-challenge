package main

import (
	"database/sql"
	"flag"
	"log"
	"os"
	"time"

	"github.com/avast/retry-go"
	_ "github.com/golang-migrate/migrate/v4/database/mysql"
	"github.com/marcelo-rocha/task-service-challenge/persistence/migration"
)

var (
	flagaddr = flag.String("addr", "root:secret7@(127.0.0.1:3306)/test?multiStatements=true&parseTime=true", "mysql connection")
	flagup   = flag.Bool("up", true, "run up migrations")
	flagdown = flag.Bool("down", false, "run down migrations")
)

func main() {

	flag.Parse()

	if *flagup && *flagdown {
		log.Println("error: should use either -up or -down")
		os.Exit(1)
	}

	var db *sql.DB
	if err := retry.Do(func() error {
		// checking connection
		var e error
		db, e = Connect(*flagaddr)
		if e != nil {
			return e
		}
		return db.Ping()
	}, retry.Delay(time.Second*2), retry.Attempts(5), retry.DelayType(retry.BackOffDelay)); err != nil {
		log.Println("connection failed", err)
		os.Exit(2)
	}

	defer db.Close()

	if *flagup {
		if err := migration.Up(db); err != nil {
			log.Println("migrate up error", err)
			os.Exit(3)
		}
	} else if *flagdown {
		if err := migration.Down(db); err != nil {
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
