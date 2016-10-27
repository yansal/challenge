package main

import (
	"flag"
	"log"
	"os"

	"github.com/jmoiron/sqlx"
	"gopkg.in/gin-gonic/gin.v1"

	_ "github.com/lib/pq"
)

var (
	db     *sqlx.DB
	router *gin.Engine
)

func main() {
	var seedFlag = flag.Bool("seed", false, "Seed database")
	flag.Parse()

	databaseURL := os.Getenv("DATABASE_URL")
	if databaseURL == "" {
		databaseURL = "sslmode=disable"
	}
	db = sqlx.MustConnect("postgres", databaseURL)

	if *seedFlag {
		create()
	}
	prepare()
	if *seedFlag {
		seed()
	}

	log.Fatal(router.Run())
}
