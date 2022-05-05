package Handler

import (
	"bytes"
	"database/sql"
	"fmt"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
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

		// var guitar Model.Guitars
		// var guitars []Model.Guitars

		rows := db.QueryRow(`SELECT "Id","Brand_Id", "Name", "Price", "Back", "Side", "Neck", "GuitarSize", "Description", "Image"  FROM guitars`)
		// if err != nil {
		// 	c.String(http.StatusInternalServerError,
		// 		fmt.Sprintf("Error reading ticks: %q", err))
		// 	return
		// }

		// defer rows.Close()
		// for rows.Next() {
		// 	if err := rows.Scan(&guitar); err != nil {
		// 		c.String(http.StatusInternalServerError,
		// 			fmt.Sprintf("Error scanning ticks: %q", err))
		// 		return
		// 	}
		// 	guitars = append(guitars, guitar)
		// }
		// c.JSON(200, guitars)
		c.JSON(200, rows)
		// c.String(http.StatusOK, fmt.Sprintf("Read from DB: %s\n", tick.String()))
	}
}
