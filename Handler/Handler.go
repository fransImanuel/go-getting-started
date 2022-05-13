package Handler

import (
	"database/sql"
	"fmt"
	"sort"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/heroku/go-getting-started/Model"
)



func TestCall(db *sql.DB) gin.HandlerFunc {
	return func(c *gin.Context) {

		var res Model.Response
		res = Model.Response{
			Message: "Test Success",
		}
		c.JSON(200, res)
		// c.String(http.StatusOK, fmt.Sprintf("Read from DB: %s\n", tick.String()))
	}
}

func GuitarByFilter(db *sql.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		
		var guitar Model.Guitars
		var guitars []Model.Guitars
		var Input Model.RequestGuitar		
		var results []Model.Result
		var res Model.Response
		var count int

		Input = Model.RequestGuitar{
			Back_ID:        c.Query("Back"),
			Side_ID:        c.Query("Side"),
			Neck_ID:        c.Query("Neck"),
			Guitarsize:  c.Query("GuitarSize"),
			Brand:       c.Query("Brand"),
			UpperPice:   c.Query("upperPrice"),
			Page:        c.Query("Page"),
		} 

		//------
		// First query to get data based on query params
		//------
		q :=`
			select g."Id", b."Rank" as "Brand_Id" , g."Name", g."Price", w1."Rank" as "Back", w2."Rank" as "Side", w3."Rank" as "Neck", s."Rank" as "GuitarSize", g."Description", g."Image" , g."WhereToBuy" 
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
			g."Price" <= $6 --Price
		`
		queryLimit := `
			ORDER BY g."Id"
			offset $7 rows fetch next 10 rows only;`
		page, err := strconv.Atoi(Input.Page)
		if err != nil {
			fmt.Println(err)
			res = Model.Response{
				Message: "Error",
				Error_Message: err,
			} 
			c.JSON(502, res)
			return
		}

		offset := pagination(page)
		rows, err := db.Query(q+cond+queryLimit,
			Input.Back_ID ,Input.Side_ID ,Input.Neck_ID, Input.Guitarsize ,Input.Brand,Input.UpperPice,offset)
		if err != nil {
			fmt.Println(err)
			res = Model.Response{
				Message: "Error",
				Error_Message: err,
			} 
			c.JSON(502, res)
			return
		}
		defer rows.Close()

		for rows.Next() {
			if err := rows.Scan(&guitar.Guitar_ID, &guitar.Brand, &guitar.Guitar_Name, &guitar.Price, &guitar.Back_ID, 
				&guitar.Side_ID, &guitar.Neck_ID, &guitar.GuitarSize, &guitar.Description, &guitar.Image, &guitar.WhereToBuy); err != nil {
					fmt.Println(err)
					res = Model.Response{
						Message: "Error",
						Error_Message: err,
					} 
					c.JSON(502, res)
					return
				}
			guitars = append(guitars, guitar)
		}

		//IF THE RESULT IS NULL/NOT FOUND
		if len(guitars) == 0 {
			if !rows.Next(){
				cond = `
				where (w1."Wood_Id" = $1 OR --back
				w2."Wood_Id" = $2 OR --side
				w3."Wood_Id" = $3) AND --neck
				s."Size_Id" = $4 And --guitarsize
				g."Price" <= $5 --Price
			`
			queryLimit = `
				ORDER BY g."Id"
				offset $6 rows fetch next 10 rows only;`
				rows, err = db.Query(q+cond+queryLimit,
					Input.Back_ID ,Input.Side_ID ,Input.Neck_ID, Input.Guitarsize ,Input.UpperPice,offset)
				if err != nil {
					fmt.Println(err)
					c.JSON(502, Model.Response{
						Message: "Error",
						Error_Message: err,
					} )
					return
				}
			}

			for rows.Next() {
				if err := rows.Scan(&guitar.Guitar_ID, &guitar.Brand, &guitar.Guitar_Name, &guitar.Price,&guitar.Back_ID, 
					&guitar.Side_ID, &guitar.Neck_ID, &guitar.GuitarSize, &guitar.Description, &guitar.Image, &guitar.WhereToBuy); err != nil {
						fmt.Println(err)
						c.JSON(502, Model.Response{
							Message: "Error",
							Error_Message: err,
						} )
						return
					}
				guitars = append(guitars, guitar)
			}

			//------
			// Second query to count the data for pagination purposes
			//------

			q =`
				select count(g."Id")
				from guitars g
				join woods w1 on (g."Back" = w1."Wood_Id")
				join woods w2 on (g."Side" = w2."Wood_Id")
				join woods w3 on (g."Neck" = w3."Wood_Id")
				join sizes s on (g."GuitarSize" = s."Size_Id")
				join brands b on (g."Brand_Id" = b."Brand_Id")
			`
			cond= `
				where w1."Wood_Id" = $1 AND --back
				w2."Wood_Id" = $2 AND --side
				w3."Wood_Id" = $3 AND --neck
				s."Size_Id" = $4 AND --guitarsize
				b."Brand_Id" = $5 AND --brand
				g."Price" <= $6 --Price
			`
			rows2, err := db.Query(q+cond,Input.Back_ID ,Input.Side_ID ,Input.Neck_ID, Input.Guitarsize ,Input.Brand,Input.UpperPice)
			if err != nil {
				fmt.Println(err)
				res = Model.Response{
					Message: "Error",
					Error_Message: err,
				} 
				c.JSON(502, res)
				return
			}

			defer rows2.Close()
			for rows2.Next() {
				if err := rows2.Scan(&count); err != nil {
					fmt.Println(err)
					res = Model.Response{
						Message: "Error",
						Error_Message: err,
					} 
					c.JSON(502, res)
					return
				}
			}
		}else{
			q =`
				select count(g."Id")
				from guitars g
				join woods w1 on (g."Back" = w1."Wood_Id")
				join woods w2 on (g."Side" = w2."Wood_Id")
				join woods w3 on (g."Neck" = w3."Wood_Id")
				join sizes s on (g."GuitarSize" = s."Size_Id")
				join brands b on (g."Brand_Id" = b."Brand_Id")
			`
			rows2, err := db.Query(q+cond,Input.Back_ID ,Input.Side_ID ,Input.Neck_ID, Input.Guitarsize ,Input.Brand,Input.UpperPice)
			if err != nil {
				fmt.Println(err)
				res = Model.Response{
					Message: "Error",
					Error_Message: err,
				} 
				c.JSON(502, res)
				return
			}

			defer rows2.Close()
			for rows2.Next() {
				if err := rows2.Scan(&count); err != nil {
					fmt.Println(err)
					res = Model.Response{
						Message: "Error",
						Error_Message: err,
					} 
					c.JSON(502, res)
					return
				}
			}
		}
		// if len(guitars) == 0 {
		// 	res = Model.Response{
		// 		Message: "Error",
		// 		Error_Message: err,
		// 	} 
		// 	c.JSON(404, res)
		// 	return
		// }
		
		
		
		


		results,err = SAW(guitars,db)
		if err != nil {
			fmt.Println(err)
			res = Model.Response{
				Message: "Error",
				Error_Message: err,
			} 
			c.JSON(502, res)
			return
		}
		//reset guitars & guitar to nil and replace to sorted guitar rating
		guitar = Model.Guitars{}
		guitars = []Model.Guitars{}
		
		//------
		// Third requery data based on sorted rating by SAW method
		//------
		for _, r := range results{
			q =`
				select g."Id", b."Rank" as "Brand_Id", b."Name" as "Brand_Name" , g."Name", g."Price", w1."Name" as "Back", w2."Name" as "Side", w3."Name" as "Neck", s."Size" as "GuitarSize", g."Description", g."Image", g."WhereToBuy" 
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
				fmt.Println(err)
				res = Model.Response{
					Message: "Error",
					Error_Message: err,
				} 
				c.JSON(502, res)
				return
			}

			defer rows3.Close()
			for rows3.Next() {
				if err := rows3.Scan(&guitar.Guitar_ID, &guitar.Brand, &guitar.Brand_Name, &guitar.Guitar_Name, &guitar.Price,&guitar.Back_Name, 
					&guitar.Side_Name, &guitar.Neck_Name, &guitar.GuitarSize, &guitar.Description, &guitar.Image, &guitar.WhereToBuy); err != nil {
					fmt.Println(err)
					res = Model.Response{
						Message: "Error",
						Error_Message: err,
					} 
					c.JSON(502, res)
					return
				}
				guitars = append(guitars, guitar)
			}
		}

		res = Model.Response{
			Message: "Success",
			Data: guitars,
			Total_Data: count,
		}
		c.JSON(200, res )
	}
}

func AllGuitar(db *sql.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		var guitar Model.Guitars
		var guitars []Model.Guitars	
		var results []Model.Result
		var res Model.Response

		q :=`
			select g."Id", b."Rank" as "Brand_Id" , g."Name", g."Price", w1."Rank" as "Back", w2."Rank" as "Side", w3."Rank" as "Neck", s."Rank" as "GuitarSize", g."Description", g."Image" , g."WhereToBuy" 
			from guitars g
			join woods w1 on (g."Back" = w1."Wood_Id")
			join woods w2 on (g."Side" = w2."Wood_Id")
			join woods w3 on (g."Neck" = w3."Wood_Id")
			join sizes s on (g."GuitarSize" = s."Size_Id")
			join brands b on (g."Brand_Id" = b."Brand_Id")
		`

		rows, err := db.Query(q)
		if err != nil {
			res = Model.Response{
				Message: "Error",
				Error_Message: err,
			} 
			c.JSON(502, res)
			return
		}
		defer rows.Close()

		for rows.Next() {
			if err := rows.Scan(&guitar.Guitar_ID, &guitar.Brand, &guitar.Guitar_Name, &guitar.Price,&guitar.Back_ID, &guitar.Side_ID, 
				&guitar.Neck_ID, &guitar.GuitarSize, &guitar.Description, &guitar.Image, &guitar.WhereToBuy); err != nil {
					res = Model.Response{
						Message: "Error",
						Error_Message: err,
					} 
					c.JSON(502, res)
					return
				}
			guitars = append(guitars, guitar)
		}

		results,err = SAW(guitars,db)
		if err != nil {
			res = Model.Response{
				Message: "Error",
				Error_Message: err,
			} 
			c.JSON(502, res)
			return
		}
		//reset guitars & guitar to nil and replace to sorted guitar rating
		guitar = Model.Guitars{}
		guitars = []Model.Guitars{}
		
		//------
		// Third requery data based on sorted rating by SAW method
		//------
		for _, r := range results{
			q =`
				select g."Id", b."Rank" as "Brand_Id" , g."Name", g."Price", w1."Rank" as "Back", w2."Rank" as "Side", w3."Rank" as "Neck", s."Rank" as "GuitarSize", g."Description", g."Image", g."WhereToBuy" 
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
				res = Model.Response{
					Message: "Error",
					Error_Message: err,
				} 
				c.JSON(502, res)
				return
			}

			defer rows3.Close()
			for rows3.Next() {
				if err := rows3.Scan(&guitar.Guitar_ID, &guitar.Brand, &guitar.Guitar_Name, &guitar.Price,&guitar.Back_ID, 
					&guitar.Side_ID, &guitar.Neck_ID, &guitar.GuitarSize, &guitar.Description, &guitar.Image, &guitar.WhereToBuy); err != nil {
					res = Model.Response{
						Message: "Error",
						Error_Message: err,
					} 
					c.JSON(502, res)
					return
				}
				guitars = append(guitars, guitar)
			}
		}

		res = Model.Response{
			Message: "Success",
			Data: guitars,
		}

		c.JSON(200, res )

	}
}

func pagination(page int)(offset int){
	if page <= 1{
		return 0
	}
	return (10 * page) - 10
}

func SAW(guitars []Model.Guitars,db *sql.DB)([]Model.Result,error){
	var d Model.Divider //d = divider (pembagi)
	var n Model.Divider //Normalization
	var ns []Model.Divider //Normalization
	var result Model.Result
	var results []Model.Result
	
	//insert divider based on cost/benefit
	for _, g := range guitars{

		if d.Price >= *g.Price || d.Price == 0  { d.Price = *g.Price }
		if d.Back <= *g.Back_ID || d.Back == 0 { d.Back = *g.Back_ID }
		if d.Side <= *g.Side_ID || d.Side == 0 { d.Side = *g.Side_ID }
		if d.Neck <= *g.Neck_ID || d.Neck == 0 { d.Neck = *g.Neck_ID }
		if d.Size <= *g.GuitarSize || d.Size == 0 { d.Size = *g.GuitarSize }
		if d.Brand <= *g.Brand || d.Brand == 0 { d.Brand = *g.Brand }
	}

	// c = calculate
	for _, c := range guitars{
		n.Guitar_ID = *c.Guitar_ID
		n.Price =  d.Price / *c.Price
		n.Back = *c.Back_ID / d.Back
		n.Side = *c.Side_ID / d.Side
		n.Neck = *c.Neck_ID / d.Neck
		n.Size = *c.GuitarSize / d.Size
		n.Brand = *c.Brand / d.Brand		
		ns = append(ns,n)
	}

	// get criteria value from database
	var c Model.Criteria
	//cm = criteria map
	cm := make(map[string]float64)

	q :=`
		select "Criteria_Name","Value" from criteria; 
	`
	rows, err := db.Query(q)
	if err != nil {
		return []Model.Result{}, err
	}

	defer rows.Close()
	for rows.Next() {
		if err := rows.Scan(&c.Criteria_Name,&c.Value); err != nil {
			return []Model.Result{}, err
		}
		
		if c.Criteria_Name == "Harga"{cm["Harga"] = c.Value}
		if c.Criteria_Name == "Back" {cm["Back"] = c.Value}
		if c.Criteria_Name == "Side" {cm["Side"] = c.Value}
		if c.Criteria_Name == "Neck" {cm["Neck"] = c.Value}
		if c.Criteria_Name == "Merk" {cm["Merk"] = c.Value}
		if c.Criteria_Name == "Size" {cm["Size"] = c.Value}
		
	}

	fmt.Println(cm)
	// for _, criteria := range cs{

	// }

	//fr = finalResult
	for _, fr:= range ns{
		result.Guitar_ID = fr.Guitar_ID
		result.Rating = (fr.Price * cm["Harga"]) + (fr.Back * cm["Back"]) + (fr.Side * cm["Side"]) + (fr.Neck * cm["Neck"]) + (fr.Size * cm["Size"]) + (fr.Brand * cm["Merk"])
		results = append(results,result)
	}

	//sorting rating
	sort.SliceStable(results, func(i, j int) bool {
		return results[i].Rating > results[j].Rating
	})

	// fmt.Println(results)

	return results,nil
}