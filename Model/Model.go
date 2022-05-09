package Model

type Guitars struct {
	Guitar_ID   *int     `json:"Guitar_ID,omitempty"  db:"Id"`
	Brand       *float64 `json:"Brand_ID,omitempty"  db:"Brand_Id"`
	Guitar_Name *string  `json:"Guitar_Name,omitempty"  db:"Name"`
	Price       *float64 `json:"Price,omitempty"  db:"Price"`
	Back        *float64 `json:"Back,omitempty"  db:"Back"`
	Side        *float64 `json:"Side,omitempty"  db:"Side"`
	Neck        *float64 `json:"Neck,omitempty"  db:"Neck"`
	GuitarSize  *float64 `json:"GuitarSize,omitempty"  db:"GuitarSize"`
	Description *string  `json:"Description,omitempty"  db:"Description"`
	Image       *string  `json:"Image,omitempty"  db:"Image"`
}

type Response struct {
	Message        string      `json:"message,omitempty"`
	Data           interface{} `json:"data,omitempty"`
	Total_Data     interface{} `json:"total_data,omitempty"`
	Error_Key      string      `json:"error_key,omitempty"`
	Error_Message  error      `json:"error_message,omitempty"`
	Secondary_Data interface{} `json:"secondary_data,omitempty"`
}

type RequestGuitar struct {
	Back        string `json:"Back,omitempty"`
	Side        string `json:"Side,omitempty"`
	Neck        string `json:"Neck,omitempty"`
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