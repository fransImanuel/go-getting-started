package Handler

import (
	"bytes"
	"database/sql"
	"fmt"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/heroku/go-getting-started/Model"
)

func DbFunc(db *sql.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		if _, err := db.Exec("CREATE TABLE IF NOT EXISTS ticks (tick timestamp)"); err != nil {
			c.String(http.StatusInternalServerError,
				fmt.Sprintf("Error creating database table: %q", err))
			return
		}

		if _, err := db.Exec("INSERT INTO ticks VALUES (now())"); err != nil {
			c.String(http.StatusInternalServerError,
				fmt.Sprintf("Error incrementing tick: %q", err))
			return
		}

		rows, err := db.Query("SELECT tick FROM ticks")
		if err != nil {
			c.String(http.StatusInternalServerError,
				fmt.Sprintf("Error reading ticks: %q", err))
			return
		}

		defer rows.Close()
		for rows.Next() {
			var tick time.Time
			if err := rows.Scan(&tick); err != nil {
				c.String(http.StatusInternalServerError,
					fmt.Sprintf("Error scanning ticks: %q", err))
				return
			}
			c.String(http.StatusOK, fmt.Sprintf("Read from DB: %s\n", tick.String()))
		}
	}
}

func RepeatHandler(r int) gin.HandlerFunc {

	return func(c *gin.Context) {
		var buffer bytes.Buffer
		for i := 0; i < r; i++ {
			buffer.WriteString("Hello from Go!\n")
		}
		c.String(http.StatusOK, buffer.String())
	}
}

func TestCall(db *sql.DB) gin.HandlerFunc {
	return func(c *gin.Context) {

		var guitar Model.Guitars
		var guitars []Model.Guitars

		rows, err := db.Query(`SELECT "Id","Brand_Id", "Name", "Price", "Back", "Side", "Neck", "GuitarSize", "Description", "Image"  FROM guitars`)
		if err != nil {
			c.String(http.StatusInternalServerError,
				fmt.Sprintf("Error reading ticks: %q", err))
			return
		}

		defer rows.Close()
		for rows.Next() {
			if err := rows.Scan(&guitar.Guitar_ID, &guitar.Brand_ID, &guitar.Guitar_Name, &guitar.Price,
				&guitar.Back, &guitar.Side, &guitar.Neck, &guitar.GuitarSize, &guitar.Description, &guitar.Image); err != nil {
				c.String(http.StatusInternalServerError,
					fmt.Sprintf("Error scanning ticks: %q", err))
				return
			}
			guitars = append(guitars, guitar)
		}
		c.JSON(200, guitars)
		// c.String(http.StatusOK, fmt.Sprintf("Read from DB: %s\n", tick.String()))
	}
}

func GuitarByFilter(db *sql.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		
		var guitar Model.Guitars
		var guitars []Model.Guitars
		// var Response Model.Response

		// Input := struct {
		// 	Back        string `json:"Back,omitempty"`
		// 	Side        string `json:"Side,omitempty"`
		// 	Neck        string `json:"Neck,omitempty"`
		// 	Guitarsize  string `json:"Guitarsize,omitempty"`
		// 	Brand       string `json:"Brand,omitempty"`
		// 	BottomPrice string `json:"bottomPrice,omitempty"`
		// 	UpperPice   string `json:"upperPrice,omitempty"`
		// 	Page        string `json:"Page,omitempt"`
		// }{
		// 	Back:        c.Query("Back"),
		// 	Side:        c.Query("Side"),
		// 	Neck:        c.Query("Neck"),
		// 	Guitarsize:  c.Query("GuitarSize"),
		// 	Brand:       c.Query("Brand"),
		// 	BottomPrice: c.Query("bottomPrice"),
		// 	UpperPice:   c.Query("upperPrice"),
		// 	Page:        c.Query("Page"),
		// }

		// fmt.Println("TEST")
		// fmt.Println(Input)

		// c.JSON(200, Input)

		q :=`
			select g."Id", b."Rank" as "Brand_Id" , g."Name", g."Price", w1."Rank" as "Back", w2."Rank" as "Side", w3."Rank" as "Neck", s."Rank" as "GuitarSize", g."Description", g."Image" 
			from guitars g
			join woods w1 on (g."Back" = w1."Wood_Id")
			join woods w2 on (g."Side" = w2."Wood_Id")
			join woods w3 on (g."Neck" = w3."Wood_Id")
			join sizes s on (g."GuitarSize" = s."Size_Id")
			join brands b on (g."Brand_Id" = b."Brand_Id")
			
		`
		// i, err := strconv.Atoi(Input.Page)
		// if err != nil {
		// 	c.String(http.StatusInternalServerError,
		// 	fmt.Sprintf("Error when convert page: %q", err))
		// 	return
		// }
		// offset := i * 10
		// fmt.Println(q+cond)
		rows, err := db.Query(q)
		if err != nil {
			c.String(http.StatusInternalServerError,
			fmt.Sprintf("Error reading ticks: %q", err))
			return
		}

		defer rows.Close()
		for rows.Next() {
			if err := rows.Scan(&guitar.Guitar_ID, &guitar.Brand_ID, &guitar.Guitar_Name, &guitar.Price,
					&guitar.Back, &guitar.Side, &guitar.Neck, &guitar.GuitarSize, &guitar.Description, &guitar.Image); err != nil {
					c.String(http.StatusInternalServerError,
					fmt.Sprintf("Error scanning ticks: %q", err))
					return
			
				guitars = append(guitars, guitar)
			}
		}
		c.JSON(200,guitars)
	}
}

