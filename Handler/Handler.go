package Handler

import (
	"database/sql"
	"fmt"
	"sort"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"github.com/heroku/go-getting-started/Model"
)


var validate *validator.Validate

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
		var cond string
		var rows *sql.Rows
		

		Input = Model.RequestGuitar{
			Price:        c.Query("Price"),
			WoodWeight:        c.Query("WoodWeight"),
			GuitarSizeWeight:  c.Query("GuitarSizeWeight"),
			BrandWeight:  c.Query("BrandWeight"),
			BrandId:  c.Query("BrandId"),
			UpperPrice:  c.Query("UpperPrice"),
			BottomPrice:  c.Query("BottomPrice"),
		} 

		

		PriceWeight, err := strconv.Atoi(Input.Price)
		if err != nil {
			fmt.Println(err)
			res = Model.Response{
				Message: "Error",
				Error_Message: err,
			} 
			c.JSON(502, res)
			return
		}

		WoodWeight, err := strconv.Atoi(Input.WoodWeight)
		if err != nil {
			fmt.Println(err)
			res = Model.Response{
				Message: "Error",
				Error_Message: err,
			} 
			c.JSON(502, res)
			return
		}

		GuitarSizeWeight, err := strconv.Atoi(Input.GuitarSizeWeight)
		if err != nil {
			fmt.Println(err)
			res = Model.Response{
				Message: "Error",
				Error_Message: err,
			} 
			c.JSON(502, res)
			return
		}
		
		BrandWeight, err := strconv.Atoi(Input.BrandWeight)
		if err != nil {
			fmt.Println(err)
			res = Model.Response{
				Message: "Error",
				Error_Message: err,
			} 
			c.JSON(502, res)
			return
		}

		Weight := Model.GuitarWeight{
			WoodWeight       : WoodWeight,
			GuitarSizeWeight : GuitarSizeWeight,
			BrandWeight      : BrandWeight,
			PriceWeight : PriceWeight,
		}

		// ------
		// First query to get data based on query params
		// ------
		
		q :=`
			select g."Id", b."Rank" as "Brand_Id" , g."Name", g."Price", w1."Rank" as "Back", w2."Rank" as "Side", w3."Rank" as "Neck", s."Rank" as "GuitarSize", g."Description", g."Image", g."WhereToBuy" 
			from guitars g
			join woods w1 on (g."Back" = w1."Wood_Id")
			join woods w2 on (g."Side" = w2."Wood_Id")
			join woods w3 on (g."Neck" = w3."Wood_Id")
			join sizes s on (g."GuitarSize" = s."Size_Id")
			join brands b on (g."Brand_Id" = b."Brand_Id")
		`
		if Input.BottomPrice == "" {
			cond = `where g."Brand_Id" = $1 AND g."Price" >= $2`
			rows, err = db.Query(q+cond,Input.BrandId, Input.UpperPrice)
			if err != nil {
				fmt.Println("err here 103")
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
						fmt.Println("err here 117")
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

		if Input.UpperPrice == "" {
			cond = `where g."Brand_Id" = $1 AND g."Price" <= $2`
			rows, err = db.Query(q+cond,Input.BrandId, Input.BottomPrice)
			if err != nil {
				fmt.Println("err here 103")
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
						fmt.Println("err here 117")
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

		if Input.UpperPrice != "" &&  Input.BottomPrice != ""{
			cond = `where g."Brand_Id" = $1 AND g."Price" >= $2 AND g."Price" <= $3`
			rows, err = db.Query(q+cond,Input.BrandId, Input.BottomPrice, Input.UpperPrice)
			if err != nil {
				fmt.Println("err here 103")
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
						fmt.Println("err here 117")
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
		
		results,err = SAW(guitars,db, Weight)
		if err != nil {
			fmt.Println("err here 315")
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
				fmt.Println("err here 344")
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
					fmt.Println("err here 358")
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
			Total_Data: len(guitars),
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
			select g."Id", b."Rank" as "Brand_Id", b."Name" as "Brand_Name" , g."Name", g."Price", w1."Rank" as "Back", w2."Rank" as "Side", w3."Rank" as "Neck", s."Rank" as "GuitarSize"
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
			if err := rows.Scan(&guitar.Guitar_ID, &guitar.Brand, &guitar.Brand_Name, &guitar.Guitar_Name, &guitar.Price,&guitar.Back_ID, 
					&guitar.Side_ID, &guitar.Neck_ID, &guitar.GuitarSize); err != nil {
					// fmt.Println(err)
					res = Model.Response{
						Message: "Error",
						Error_Message: err,
					} 
					c.JSON(502, res)
					return
				}
			guitars = append(guitars, guitar)
		}

		results,err = SAW(guitars,db,Model.GuitarWeight{})
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
				// fmt.Println(err)
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
					// fmt.Println(err)
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
			Total_Data: len(guitars),
		}
		c.JSON(200, res )

	}
}

func SAW(guitars []Model.Guitars,db *sql.DB, Weight Model.GuitarWeight)([]Model.Result,error){
	var d Model.Divider //d = divider (pembagi)
	var n Model.Divider //Normalization
	var ns []Model.Divider //Normalization
	var result Model.Result
	var results []Model.Result
	var WeightTotal float64
	
	//insert divider based on cost/benefit
	for _, g := range guitars{

		if d.Price >= *g.Price || d.Price == 0  { d.Price = *g.Price }
		if d.Back <= *g.Back_ID || d.Back == 0 { d.Back = *g.Back_ID }
		if d.Side <= *g.Side_ID || d.Side == 0 { d.Side = *g.Side_ID }
		if d.Neck <= *g.Neck_ID || d.Neck == 0 { d.Neck = *g.Neck_ID }
		if d.Size <= *g.GuitarSize || d.Size == 0 { d.Size = *g.GuitarSize }
		if d.Brand <= *g.Brand || d.Brand == 0 { d.Brand = *g.Brand }
	}

	// fmt.Println(d)

	// c = calculate
	for _, c := range guitars{
		// fmt.Printf("Guitar_ID: %v, Brand: %v, Brand_Name: %v, Guitar_Name:%v, Price: %v, Back_ID:%v, Side_ID:%v,Neck_ID:%v, GuitarSize:%v\n",	*c.Guitar_ID,*c.Brand,*c.Brand_Name,*c.Guitar_Name,*c.Price,*c.Back_ID,*c.Side_ID,*c.Neck_ID,*c.GuitarSize)
		n.Guitar_ID = *c.Guitar_ID
		n.Price =  d.Price / *c.Price
		n.Back = *c.Back_ID / d.Back
		n.Side = *c.Side_ID / d.Side
		n.Neck = *c.Neck_ID / d.Neck
		n.Size = *c.GuitarSize / d.Size
		n.Brand = *c.Brand / d.Brand
		// fmt.Printf("%v | %v, %v, %v, %v, %v, %v\n",n.Guitar_ID,n.Price, n.Back, n.Side, n.Neck, n.Size, n.Brand)		
		ns = append(ns,n)
	}

	// get criteria value from database
	var c Model.Criteria
	//cm = criteria map
	cm := make(map[string]float64)

	if Weight == (Model.GuitarWeight{}){
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
			
			if c.Criteria_Name == "Harga"{cm["Harga"] = c.Value / 39}
			if c.Criteria_Name == "Back" {cm["Back"] = c.Value / 39}
			if c.Criteria_Name == "Side" {cm["Side"] = c.Value / 39}
			if c.Criteria_Name == "Neck" {cm["Neck"] = c.Value / 39}
			if c.Criteria_Name == "Merk" {cm["Merk"] = c.Value / 39}
			if c.Criteria_Name == "Size" {cm["Size"] = c.Value / 39}
			
		}
	}else{
		WeightTotal =  float64(Weight.PriceWeight) + float64((Weight.WoodWeight*3)) +float64(Weight.BrandWeight)+float64(Weight.GuitarSizeWeight)
		// fmt.Printf("WeightTotal = %v",WeightTotal)
		cm["Harga"] = float64(Weight.PriceWeight)/WeightTotal
		cm["Back"] = float64(Weight.WoodWeight)/WeightTotal
		cm["Side"] = float64(Weight.WoodWeight)/WeightTotal
		cm["Neck"] = float64(Weight.WoodWeight)/WeightTotal
		cm["Merk"] = float64(Weight.BrandWeight)/WeightTotal
		cm["Size"] = float64(Weight.GuitarSizeWeight)/WeightTotal
	}

	


	// fmt.Println("\n\nFINAL RESULT\n\n")
	// fmt.Println(cm)

	//fr = finalResult
	for _, fr:= range ns{
		result.Guitar_ID = fr.Guitar_ID
		result.Rating = (fr.Price * cm["Harga"]) + (fr.Back * cm["Back"]) + (fr.Side * cm["Side"]) + 
						(fr.Neck * cm["Neck"]) + (fr.Size * cm["Size"]) + (fr.Brand * cm["Merk"])
		
		// fmt.Printf("%v |harga = %v * %v = %v, %v, %v, %v, %v, %v = %v\n",result.Guitar_ID,fr.Price, cm["Harga"],fr.Price * cm["Harga"], fr.Back * cm["Back"], fr.Side * cm["Side"], fr.Neck * cm["Neck"], fr.Size * cm["Size"], fr.Brand * cm["Merk"],result.Rating)
		// fmt.Printf("%v |%v, %v, %v, %v, %v, %v = %v\n",result.Guitar_ID,fr.Price * cm["Harga"], fr.Back * cm["Back"], fr.Side * cm["Side"], fr.Neck * cm["Neck"], fr.Size * cm["Size"], fr.Brand * cm["Merk"],result.Rating)
		// fmt.Println(result)
		results = append(results,result)
	}


	//sorting rating
	sort.SliceStable(results, func(i, j int) bool {
		return results[i].Rating > results[j].Rating
	})

	// for _, z:=range results{
	// 	fmt.Println(z)
	// }

	return results,nil
	// return []Model.Result{},nil
}

func AddGuitar(db *sql.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		var res Model.Response
		var Input Model.AddGuitar	
		
		err := c.BindJSON(&Input)
		if err != nil {
			res = Model.Response{
				Message: "Error",
				Error_Message: err,
			} 
			c.JSON(400, res)
			return
		}
		// fmt.Println(Input)
		
		validate = validator.New()
		err = validate.Struct(Input)
		if err != nil {
			fmt.Println(err)
			res = Model.Response{
				Message: "Error",
				Error_Message: err,
			} 
			c.JSON(400, res)
			return
		}

		q := `
			INSERT INTO "guitars" ("Brand_Id", "Name", "Price", "Back", "Side", "Neck", "GuitarSize", "Description", "Image", "WhereToBuy") 
			VALUES ($1,$2,$3,$4,$5,$6,$7,$8,$9,$10)
		`
		_, err = db.Exec(q,Input.Brand_ID, Input.Guitar_Name, Input.Price, Input.Back_ID, Input.Side_ID, Input.Neck_ID,
			Input.Size_ID, Input.Description, Input.Image, Input.WhereToBuy)
		if err != nil {
			fmt.Println(err)
			res = Model.Response{
				Message: "Error",
				Error_Message: err,
			} 
			c.JSON(502, res)
			return
		}

		res = Model.Response{
			Message: "Success",
		}
		c.JSON(200, res )

	}
}

func Login(db *sql.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		var res Model.Response
		var countUser int
		var Input Model.Login	
		
		err := c.BindJSON(&Input)
		if err != nil {
			res = Model.Response{
				Message: "Error",
				Error_Message: err,
			} 
			c.JSON(400, res)
			return
		}
		// fmt.Println(Input)
		
		validate = validator.New()
		err = validate.Struct(Input)
		if err != nil {
			fmt.Println(err)
			res = Model.Response{
				Message: "Error",
				Error_Message: err,
			} 
			c.JSON(400, res)
			return
		}

		q := `
			select count(a."User_Id")
			from admininfo a
			where a."username" = $1 AND a."password" = $2;
		`
		rows, err := db.Query(q,Input.Username, Input.Password)
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
			if err := rows.Scan(&countUser); err != nil {
				fmt.Println("err here 543")
				fmt.Println(err)
				res = Model.Response{
					Message: "Error",
					Error_Message: err,
				} 
				c.JSON(502, res)
				return
			}
		}

		if countUser < 1 {
			res = Model.Response{
				Message: "user not found",
			} 
			c.JSON(404, res)
			return
		}

		res = Model.Response{
			Message: "user found",
		}
		c.JSON(200, res )

	}
}

func AllGuitarForAdmin(db *sql.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		var guitar Model.Guitars
		var guitars []Model.Guitars	
		var res Model.Response

		q :=`
			select g."Id", b."Brand_Id" as "Brand_Id", b."Name" as "Brand_Name" , g."Name", g."Price",
					w1."Wood_Id" as "Back", w1."Name" as "Back_Name", 
					w2."Wood_Id" as "Side", w2."Name" as "Side_Name", 
					w3."Wood_Id" as "Neck", w3."Name" as "Neck_Name", 
					s."Size" as "GuitarSize",
					g."Description",
					g."Image",
					g."WhereToBuy"
			from guitars g	
			join woods w1 on (g."Back" = w1."Wood_Id")
			join woods w2 on (g."Side" = w2."Wood_Id")
			join woods w3 on (g."Neck" = w3."Wood_Id")
			join sizes s on (g."GuitarSize" = s."Size_Id")
			join brands b on (g."Brand_Id" = b."Brand_Id")
			order by g."Id"
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
			if err := rows.Scan(&guitar.Guitar_ID, &guitar.Brand, &guitar.Brand_Name, &guitar.Guitar_Name, &guitar.Price, 
					&guitar.Back_ID, &guitar.Back_Name, 
					&guitar.Side_ID, &guitar.Side_Name,
					&guitar.Neck_ID, &guitar.Neck_Name,
					&guitar.GuitarSize,&guitar.Description,&guitar.Image,&guitar.WhereToBuy); err != nil {
					// fmt.Println(err)
					res = Model.Response{
						Message: "Error",
						Error_Message: err,
					} 
					c.JSON(502, res)
					return
				}
			guitars = append(guitars, guitar)
		}

		res = Model.Response{
			Message: "Success",
			Data: guitars,
			Total_Data: len(guitars),
		}
		c.JSON(200, res )

	}
}

func UpdateGuitar(db *sql.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		var res Model.Response
		var Input Model.AddGuitar	

		guitarID := c.Param("id")
		
		err := c.BindJSON(&Input)
		if err != nil {
			res = Model.Response{
				Message: "Error",
				Error_Message: err,
			} 
			c.JSON(400, res)
			return
		}
		// fmt.Println(Input)
		
		validate = validator.New()
		err = validate.Struct(Input)
		if err != nil {
			fmt.Println(err)
			res = Model.Response{
				Message: "Error",
				Error_Message: err,
			} 
			c.JSON(400, res)
			return
		}

		q := `
			UPDATE guitars 
			SET "Brand_Id" = $1,
				"Name" = $2,
				"Price" = $3,
				"Back" = $4, 
				"Side" = $5, 
				"Neck" = $6, 
				"GuitarSize" = $7, 
				"Description" = $8, 
				"Image" = $9,
				"WhereToBuy" = $10
			WHERE "Id" = $11;
		`
		_, err = db.Exec(q,Input.Brand_ID, Input.Guitar_Name, Input.Price, Input.Back_ID, Input.Side_ID, Input.Neck_ID,
			Input.Size_ID, Input.Description, Input.Image, Input.WhereToBuy, guitarID)
		if err != nil {
			fmt.Println(err)
			res = Model.Response{
				Message: "Error",
				Error_Message: err,
			} 
			c.JSON(502, res)
			return
		}

		res = Model.Response{
			Message: "Success",
		}
		c.JSON(200, res )

	}
}

func DeleteGuitar(db *sql.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		var res Model.Response

		guitarID := c.Param("id")

		q := `
			DELETE FROM "guitars"
			WHERE "Id" = $1;
		`
		_, err := db.Exec(q,guitarID)
		if err != nil {
			fmt.Println(err)
			res = Model.Response{
				Message: "Error",
				Error_Message: err,
			} 
			c.JSON(502, res)
			return
		}

		res = Model.Response{
			Message: "Success",
		}
		c.JSON(200, res )

	}
}
