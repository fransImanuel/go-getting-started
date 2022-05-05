package main

import (
	"database/sql"
	"log"
	"net/http"
	"os"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/heroku/go-getting-started/Handler"
	_ "github.com/heroku/x/hmetrics/onload"
	_ "github.com/lib/pq"
)

func main() {
	// os.Setenv("PORT", "3001")
	port := os.Getenv("PORT")

	if port == "" {
		log.Fatal("$PORT must be set")
	}

	tStr := os.Getenv("REPEAT")
	repeat, err := strconv.Atoi(tStr)
	if err != nil {
		log.Printf("Error converting $REPEAT to an int: %q - Using default\n", err)
		repeat = 5
	}

	// db_url := "postgres://lxepsqkvqrpjwd:a12bfeec5b497d7bfe23f0bc6b2a62026e670c85453cad6ffb1c2194af2a5180@ec2-18-210-64-223.compute-1.amazonaws.com:5432/d6hcskiftuk3ai?sslmode=disable"
	// db, err := sql.Open("postgres", db_url)
	db, err := sql.Open("postgres", os.Getenv("DATABASE_URL"))
	if err != nil {
		log.Fatalf("Error opening database: %q", err)
	}

	router := gin.New()
	router.Use(gin.Logger())
	router.LoadHTMLGlob("templates/*.tmpl.html")
	router.Static("/static", "static")

	router.GET("/", func(c *gin.Context) {
		c.HTML(http.StatusOK, "index.tmpl.html", nil)
	})

	router.GET("/repeat", Handler.RepeatHandler(repeat))

	router.GET("/db", Handler.DbFunc(db))

	router.GET("/test", Handler.TestCall(db))

	router.Run(":" + port)
}
