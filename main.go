package main

import (
	"flag"
	"log"

	"gopkg.in/gin-gonic/gin.v1"

	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
)

var (
	db     *sqlx.DB
	router *gin.Engine
)

func main() {
	var seedFlag = flag.Bool("seed", false, "Seed database")
	flag.Parse()

	db = sqlx.MustConnect("postgres", "sslmode=disable")

	if *seedFlag {
		seed()
	}

	log.Fatal(router.Run())
}
