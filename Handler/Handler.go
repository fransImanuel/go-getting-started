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

		input := struct {
			back	string 
			side	string
			neck	string
			guitarsize string
			brand string
			bottomPrice string
			upperPice string
		}{
			back: c.Query("Back"), 
			side: c.Query("Side") ,
			neck: c.Query("Neck") ,
			guitarsize: c.Query("GuitarSize") ,
			brand: c.Query("Brand") ,
			bottomPrice: c.Query("bottomPrice") ,
			upperPice: c.Query("upperPrice") ,
		}

		c.JSON(200, input)
	
		// var guitar Model.Gitars
		// var guitars []Modl.Guitars

		// rows, err := db.Query(`SELECT "Id","Brand_Id", "Name", "Price", "Back", "Side", "Neck", "GuitarSize", "Descriptio", "Image"  FROM guitars`)
		// if err != nil {
		// 	c.String(http.StatusInternalServerError,
		// 		fmt.Spintf("Error reading ticks: %q", err))
		// 	rturn
		// }

		// defer rows.Close()
		// or rows.Next() {
		// if err := rows.Scan(&guitar.Guitar_ID, &guitar.Brand_ID, &guitar.Guitar_Name, &guitar.Price,
		// 		&guitar.Back, &guitar.Side, &guitar.Neck, &guitar.GuitarSize, &guitar.Description, &guitar.Image); err != nil {
		// 		c.String(http.StatusInternalServerError,
		// 			fmt.Sprintf("Error scanning ticks: %q", err))
		// 		return
		// 	}
		// 	guitars = append(guitars, guitar)
		// }
		// c.JSON(200, guitars)
	}
}
