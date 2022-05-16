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
	// os.Setenv("PORT", "3002")
	//comment this for local testing
	port := os.Getenv("PORT")

	if port == "" {
		log.Fatal("$PORT must be set")
	}

	

	//comment this for local testing
	// db_url := "postgres://snrqhapoddkcil:d28075e479a43de8d3563ed9bb676e3278b4b4cb27be41af9eb315243f379654@ec2-54-165-184-219.compute-1.amazonaws.com:5432/d9q283dkhak1u0"
	// db, err := sql.Open("postgres", db_url)
	//comment this for local testing
	db, err := sql.Open("postgres", os.Getenv("DATABASE_URL"))
	if err != nil {
		log.Fatalf("Error opening database: %q", err)
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
	
	router.GET("/get/allguitar",Handler.AllGuitar(db))
	
	router.POST("/addguitar",Handler.AddGuitar(db))
	
	router.PUT("/updateguitar/:id",Handler.UpdateGuitar(db))
	
	router.DELETE("/deleteguitar/:id",Handler.DeleteGuitar(db))
	
	router.POST("/login",Handler.Login(db))

	router.Run(":" + port)
}
