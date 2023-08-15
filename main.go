package main

import (
	"database/sql"
	"log"
	"net/http"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/heroku/go-getting-started/Handler"
	"github.com/heroku/go-getting-started/Middleware"
	_ "github.com/heroku/x/hmetrics/onload"
	_ "github.com/lib/pq"
)

func main() {
	//comment this for local testing
	os.Setenv("PORT", "5000")
	//comment this for local testing
	port := os.Getenv("PORT")

	if port == "" {
		log.Fatal("$PORT must be set")
	}

	// comment this for local testing
	// db_url := "db_url_from_herokupostgresql"
	// db, err := sql.Open("postgres", db_url)
	//comment this for local testing
	// db, err := sql.Open("postgres", os.Getenv("DATABASE_URL"))
	connStr := "user=postgres password=123456 sslmode=disable"
	db, err := sql.Open("postgres", connStr)
	if err != nil {
		log.Fatalf("Error opening database: %q", err)
		panic(err)
	}

	router := gin.New()
	router.Use(gin.Logger())
	router.Use(Middleware.CORSMiddleware())
	router.LoadHTMLGlob("templates/*.tmpl.html")
	router.Static("/static", "static")

	router.GET("/", func(c *gin.Context) {
		c.HTML(http.StatusOK, "index.tmpl.html", nil)
	})

	router.GET("/test", Handler.TestCall(db))

	router.GET("/get/guitarbyfilter", Handler.GuitarByFilter(db))

	router.GET("/get/allguitar", Handler.AllGuitar(db))

	router.GET("/get/gitarforadmin", Handler.AllGuitarForAdmin(db))

	router.POST("/addguitar", Handler.AddGuitar(db))

	router.PUT("/updateguitar/:id", Handler.UpdateGuitar(db))

	router.DELETE("/deleteguitar/:id", Handler.DeleteGuitar(db))

	router.POST("/login", Handler.Login(db))

	router.Run(":" + port)

}
