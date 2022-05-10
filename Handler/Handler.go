package Handler

import (
	"bytes"
	"database/sql"
	"fmt"
	"net/http"
	"sort"
	"strconv"
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
			if err := rows.Scan(&guitar.Guitar_ID, &guitar.Brand, &guitar.Guitar_Name, &guitar.Price,
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
		var Input Model.RequestGuitar		
		var results []Model.Result

		Input = Model.RequestGuitar{
			Back:        c.Query("Back"),
			Side:        c.Query("Side"),
			Neck:        c.Query("Neck"),
			Guitarsize:  c.Query("GuitarSize"),
			Brand:       c.Query("Brand"),
			BottomPrice: c.Query("bottomPrice"),
			UpperPice:   c.Query("upperPrice"),
			Page:        c.Query("Page"),
		} 

		//------
		// First query to get data based on query params
		//------
		q :=`
			select g."Id", b."Rank" as "Brand_Id" , g."Name", g."Price", w1."Rank" as "Back", w2."Rank" as "Side", w3."Rank" as "Neck", s."Rank" as "GuitarSize", g."Description", g."Image" 
			from guitars g
			join woods w1 on (g."Back" = w1."Wood_Id")
			join woods w2 on (g."Side" = w2."Wood_Id")
			join woods w3 on (g."Neck" = w3."Wood_Id")
			join sizes s on (g."GuitarSize" = s."Size_Id")
			join brands b on (g."Brand_Id" = b."Brand_Id")
		`
		cond:= `
			where w1."Wood_Id" = $1 AND --back
			w2."Wood_Id" = $2 AND --side
			w3."Wood_Id" = $3 AND --neck
			s."Size_Id" = $4 AND --guitarsize
			b."Brand_Id" = $5 AND --brand
			(g."Price" >= $6 AND g."Price" <= $7) --Price
		`
		queryLimit := `
			ORDER BY g."Id"
			offset $8 rows fetch next 10 rows only;`
		page, err := strconv.Atoi(Input.Page)
		if err != nil {
			c.JSON(502, Model.Response{
				Message: "Error",
				Error_Message: err,
			} )
			return
		}

		offset := pagination(page)
		rows, err := db.Query(q+cond+queryLimit,
			Input.Back ,Input.Side ,Input.Neck, Input.Guitarsize ,Input.Brand,Input.BottomPrice,Input.UpperPice,offset)
		if err != nil {
			c.JSON(502, Model.Response{
				Message: "Error",
				Error_Message: err,
			} )
			return
		}
		defer rows.Close()

		for rows.Next() {
			if err := rows.Scan(&guitar.Guitar_ID, &guitar.Brand, &guitar.Guitar_Name, &guitar.Price,
					&guitar.Back, &guitar.Side, &guitar.Neck, &guitar.GuitarSize, &guitar.Description, &guitar.Image); err != nil {
					c.JSON(502, Model.Response{
						Message: "Error",
						Error_Message: err,
					} )
					return
				}
			guitars = append(guitars, guitar)
		}

		//IF THE RESULT IS NULL/NOT FOUND
		if len(guitars) == 0 {
			if !rows.Next(){
				cond = `
				where w1."Wood_Id" = $1 OR --back
				w2."Wood_Id" = $2 OR --side
				w3."Wood_Id" = $3 OR --neck
				s."Size_Id" = $4 OR --guitarsize
				b."Brand_Id" = $5 OR --brand
				(g."Price" >= $6 OR g."Price" <= $7) --Price
			`
				rows, err = db.Query(q+cond+queryLimit,
					Input.Back ,Input.Side ,Input.Neck, Input.Guitarsize ,Input.Brand,Input.BottomPrice,Input.UpperPice,offset)
				if err != nil {
					c.JSON(502, Model.Response{
						Message: "Error",
						Error_Message: err,
					} )
					return
				}
			}

			for rows.Next() {
				if err := rows.Scan(&guitar.Guitar_ID, &guitar.Brand, &guitar.Guitar_Name, &guitar.Price,
						&guitar.Back, &guitar.Side, &guitar.Neck, &guitar.GuitarSize, &guitar.Description, &guitar.Image); err != nil {
						c.JSON(502, Model.Response{
							Message: "Error",
							Error_Message: err,
						} )
						return
					}
				guitars = append(guitars, guitar)
			}
		}
		
		//------
		// Second query to count the data for pagination purposes
		//------
		var count int
		q =`
			select count(g."Id")
			from guitars g
			join woods w1 on (g."Back" = w1."Wood_Id")
			join woods w2 on (g."Side" = w2."Wood_Id")
			join woods w3 on (g."Neck" = w3."Wood_Id")
			join sizes s on (g."GuitarSize" = s."Size_Id")
			join brands b on (g."Brand_Id" = b."Brand_Id")
		`
		rows2, err := db.Query(q+cond,Input.Back ,Input.Side ,Input.Neck, Input.Guitarsize ,Input.Brand,Input.BottomPrice,Input.UpperPice)
		if err != nil {
			c.String(http.StatusInternalServerError,
			fmt.Sprintf("Error rows2: %q", err))
			return
		}

		defer rows2.Close()
		for rows2.Next() {
			if err := rows2.Scan(&count); err != nil {
				c.JSON(502, Model.Response{
					Message: "Error",
					Error_Message: err,
				} )
				return
			}
		}


		results = SAW(guitars)
		//reset guitars & guitar to nil and replace to sorted guitar rating
		guitar = Model.Guitars{}
		guitars = []Model.Guitars{}
		
		//------
		// Third requery data based on sorted rating by SAW method
		//------
		for _, r := range results{
			q =`
				select g."Id", b."Rank" as "Brand_Id" , g."Name", g."Price", w1."Name" as "Back", w2."Name" as "Side", w3."Name" as "Neck", s."Size" as "GuitarSize", g."Description", g."Image" 
				from guitars g
				join woods w1 on (g."Back" = w1."Wood_Id")
				join woods w2 on (g."Side" = w2."Wood_Id")
				join woods w3 on (g."Neck" = w3."Wood_Id")
				join sizes s on (g."GuitarSize" = s."Size_Id")
				join brands b on (g."Brand_Id" = b."Brand_Id")
				where  g."Id" = $1
			`
			rows3, err := db.Query(q,r.Guitar_ID)
			if err != nil {
				c.JSON(502, Model.Response{
					Message: "Error",
					Error_Message: err,
				} )
				return
			}

			defer rows3.Close()
			for rows3.Next() {
				if err := rows3.Scan(&guitar.Guitar_ID, &guitar.Brand, &guitar.Guitar_Name, &guitar.Price,
					&guitar.Back, &guitar.Side, &guitar.Neck, &guitar.GuitarSize, &guitar.Description, &guitar.Image); err != nil {
					c.JSON(502, Model.Response{
						Message: "Error",
						Error_Message: err,
					} )
					return
				}
				guitars = append(guitars, guitar)
			}
		}

		c.JSON(200, Model.Response{
			Message: "Success",
			Data: guitars,
			Total_Data: count,
		} )
	}
}

func AllGuitar(db *sql.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		var guitar Model.Guitars
		var guitars []Model.Guitars	
		var results []Model.Result

		q :=`
			select g."Id", b."Rank" as "Brand_Id" , g."Name", g."Price", w1."Rank" as "Back", w2."Rank" as "Side", w3."Rank" as "Neck", s."Rank" as "GuitarSize", g."Description", g."Image" 
			from guitars g
			join woods w1 on (g."Back" = w1."Wood_Id")
			join woods w2 on (g."Side" = w2."Wood_Id")
			join woods w3 on (g."Neck" = w3."Wood_Id")
			join sizes s on (g."GuitarSize" = s."Size_Id")
			join brands b on (g."Brand_Id" = b."Brand_Id")
		`

		rows, err := db.Query(q)
		if err != nil {
			c.JSON(502, Model.Response{
				Message: "Error",
				Error_Message: err,
			} )
			return
		}
		defer rows.Close()

		for rows.Next() {
			if err := rows.Scan(&guitar.Guitar_ID, &guitar.Brand, &guitar.Guitar_Name, &guitar.Price,
					&guitar.Back, &guitar.Side, &guitar.Neck, &guitar.GuitarSize, &guitar.Description, &guitar.Image); err != nil {
					c.JSON(502, Model.Response{
						Message: "Error",
						Error_Message: err,
					} )
					return
				}
			guitars = append(guitars, guitar)
		}

		results = SAW(guitars)
		//reset guitars & guitar to nil and replace to sorted guitar rating
		guitar = Model.Guitars{}
		guitars = []Model.Guitars{}
		
		//------
		// Third requery data based on sorted rating by SAW method
		//------
		for _, r := range results{
			q =`
				select g."Id", b."Rank" as "Brand_Id" , g."Name", g."Price", w1."Rank" as "Back", w2."Rank" as "Side", w3."Rank" as "Neck", s."Rank" as "GuitarSize", g."Description", g."Image" 
				from guitars g
				join woods w1 on (g."Back" = w1."Wood_Id")
				join woods w2 on (g."Side" = w2."Wood_Id")
				join woods w3 on (g."Neck" = w3."Wood_Id")
				join sizes s on (g."GuitarSize" = s."Size_Id")
				join brands b on (g."Brand_Id" = b."Brand_Id")
				where  g."Id" = $1
			`
			rows3, err := db.Query(q,r.Guitar_ID)
			if err != nil {
				c.JSON(502, Model.Response{
					Message: "Error",
					Error_Message: err,
				} )
				return
			}

			defer rows3.Close()
			for rows3.Next() {
				if err := rows3.Scan(&guitar.Guitar_ID, &guitar.Brand, &guitar.Guitar_Name, &guitar.Price,
					&guitar.Back, &guitar.Side, &guitar.Neck, &guitar.GuitarSize, &guitar.Description, &guitar.Image); err != nil {
					c.JSON(502, Model.Response{
						Message: "Error",
						Error_Message: err,
					} )
					return
				}
				guitars = append(guitars, guitar)
			}
		}

		c.JSON(200, Model.Response{
			Message: "Success",
			Data: guitars,
		} )

	}
}

func pagination(page int)(offset int){
	if page <= 1{
		return 0
	}
	return (10 * page) - 10
}

func SAW(guitars []Model.Guitars)[]Model.Result{
	var d Model.Divider //d = divider (pembagi)
	var n Model.Divider //Normalization
	var ns []Model.Divider //Normalization
	var result Model.Result
	var results []Model.Result
	
	for _, g := range guitars{

		if d.Price >= *g.Price || d.Price == 0  { d.Price = *g.Price }
		if d.Back <= *g.Back || d.Back == 0 { d.Back = *g.Back }
		if d.Side <= *g.Side || d.Side == 0 { d.Side = *g.Side }
		if d.Neck <= *g.Neck || d.Neck == 0 { d.Neck = *g.Neck }
		if d.Size <= *g.GuitarSize || d.Size == 0 { d.Size = *g.GuitarSize }
		if d.Brand <= *g.Brand || d.Brand == 0 { d.Brand = *g.Brand }
	}

	// calculate = c
	for _, c := range guitars{
		n.Guitar_ID = *c.Guitar_ID
		n.Price =  d.Price / *c.Price
		n.Back = *c.Back / d.Back
		n.Side = *c.Side / d.Side
		n.Neck = *c.Neck / d.Neck
		n.Size = *c.GuitarSize / d.Size
		n.Brand = *c.Brand / d.Brand		
		ns = append(ns,n)
	}

	//fr = finalResult
	for _, fr:= range ns{
		result.Guitar_ID = fr.Guitar_ID
		result.Rating = (fr.Price * d.Price) + (fr.Back * d.Back) + (fr.Side * d.Side) + (fr.Neck * d.Neck) + (fr.Size * d.Size) + (fr.Brand * d.Brand)
		results = append(results,result)
	}

	//sorting rating
	sort.SliceStable(results, func(i, j int) bool {
		return results[i].Rating > results[j].Rating
	})

	return results
}