package main

import (
	"database/sql"
	"fmt"
	"log"
	"time"
	"net/http"
	"os"

	"authnication/data"

	_ "github.com/jackc/pgconn"
	_ "github.com/jackc/pgx/v4"
	_ "github.com/jackc/pgx/v4/stdlib"
)

var counts = 0

const webPort = "80"

type Config struct {
	DB *sql.DB
	Models data.Models
}

func main() {
	log.Println("starting autnication service")

	//dbへ接続
	conn := connectToDB()
	if conn == nil {
		log.Panic("can't connect to postgres")
	}

	// configを設定
	app := Config{
		DB: conn,
		Models: data.New(conn),
	}

	srv := &http.Server {
		Addr: fmt.Sprintf(":%s", webPort),
		Handler: app.routes(),
	}

	err := srv.ListenAndServe()
	if err != nil {
		log.Panic(err)
	}
}

func openDB(dsn string) (*sql.DB, error) {
	db, err := sql.Open("pgx", dsn)
	if err != nil {
		return nil, err
	}

	err = db.Ping()
	if err != nil {
		return nil, err
	}

	return db, nil
}

func connectToDB() *sql.DB {
	dsn := os.Getenv("DSN")

	for {
		connection, err := openDB(dsn)

		if err != nil {
			log.Println("postgress not yet ready ///")
			counts++ 
		} else {
			log.Println("Conected to Postgres")
			return connection
		}

		if counts > 10 {
			log.Println(err)
			return nil
		}

		log.Println("Backing off for two second ...")
		time.Sleep(2 * time.Second)
		continue
	}
}