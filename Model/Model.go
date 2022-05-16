package Model

type Guitars struct {
	Guitar_ID   *int     `json:"Guitar_ID,omitempty"  db:"Id"`
	Brand       *float64 `json:"Brand_ID,omitempty"  db:"Brand_Id"`
	Brand_Name  *string  `json:"Brand_Name,omitempty"  db:"Brand_Name"`
	Guitar_Name *string  `json:"Guitar_Name,omitempty"  db:"Name"`
	Price       *float64 `json:"Price,omitempty"  db:"Price"`
	Back_ID     *float64 `json:"Back,omitempty"  db:"Back"`
	Side_ID     *float64 `json:"Side,omitempty"  db:"Side"`
	Neck_ID     *float64 `json:"Neck,omitempty"  db:"Neck"`
	Back_Name   *string  `json:"Back_Name,omitempty"  db:"Back"`
	Side_Name   *string  `json:"Side_Name,omitempty"  db:"Side"`
	Neck_Name   *string  `json:"Neck_Name,omitempty"  db:"Neck"`
	GuitarSize  *float64 `json:"GuitarSize,omitempty"  db:"GuitarSize"`
	Description *string  `json:"Description,omitempty"  db:"Description"`
	Image       *string  `json:"Image,omitempty"  db:"Image"`
	WhereToBuy  *string  `json:"WhereToBuy,omitempty"  db:"WhereToBuy"`
}

type Response struct {
	Message        string      `json:"message,omitempty"`
	Data           interface{} `json:"data,omitempty"`
	Total_Data     interface{} `json:"total_data,omitempty"`
	Error_Key      string      `json:"error_key,omitempty"`
	Error_Message  error       `json:"error_message,omitempty"`
	Secondary_Data interface{} `json:"secondary_data,omitempty"`
	Criteria_Found bool        `json:"criteria_found,omitempty"`
}

type RequestGuitar struct {
	Back_ID     string `json:"Back,omitempty"`
	Side_ID     string `json:"Side,omitempty"`
	Neck_ID     string `json:"Neck,omitempty"`
	Guitarsize  string `json:"Guitarsize,omitempty"`
	Brand       string `json:"Brand,omitempty"`
	BottomPrice string `json:"bottomPrice,omitempty"`
	UpperPice   string `json:"upperPrice,omitempty"`
	Page        string `json:"Page,omitempt"`
}

type Divider struct {
	Guitar_ID int
	Price     float64
	Back      float64
	Side      float64
	Neck      float64
	Size      float64
	Brand     float64
}

type Result struct {
	Guitar_ID int
	Rating    float64
}

type Criteria struct {
	Criteria_Name string
	Value         float64
}

type AddGuitar struct {
	Brand_ID    int     `json:"Brand_ID,omitempty" validate:"numeric,required"`
	Back_ID     int     `json:"Back_ID,omitempty" validate:"numeric,required"`
	Side_ID     int     `json:"Side_ID,omitempty"  validate:"numeric,required"`
	Neck_ID     int     `json:"Neck_ID,omitempty"  validate:"numeric,required"`
	Size_ID     int     `json:"Size_ID,omitempty"  validate:"numeric,required"`
	Price       float64 `json:"Price,omitempty" validate:"gt=0,required"`
	Guitar_Name string  `json:"Guitar_Name,omitempty" validate:"required"`
	Description string  `json:"Description,omitempty"  validate:"required"`
	Image       string  `json:"Image,omitempty" validate:"required"`
	WhereToBuy  string  `json:"WhereToBuy,omitempty" validate:"required"`
}

type Login struct {
	Username string `json:"Username,omitempty" validate:"required" `
	Password string `json:"Password,omitempty" validate:"required" `
}