package Model

type Guitars struct {
	Guitar_ID   *int     `json:"Guitar_ID,omitempty"  db:"Id"`
	Brand_ID    *int     `json:"Brand_ID,omitempty"  db:"Brand_Id"`
	Guitar_Name *string  `json:"Guitar_Name,omitempty"  db:"Name"`
	Price       *float64 `json:"Price,omitempty"  db:"Price"`
	Back        *int     `json:"Back,omitempty"  db:"Back"`
	Side        *int     `json:"Side,omitempty"  db:"Side"`
	Neck        *int     `json:"Neck,omitempty"  db:"Neck"`
	GuitarSize  *int     `json:"GuitarSize,omitempty"  db:"GuitarSize"`
	Description *string  `json:"Description,omitempty"  db:"Description"`
	Image       *string  `json:"Image,omitempty"  db:"Image"`
}

type Response struct {
	Message        string      `json:"message,omitempty"`
	Data           interface{} `json:"data,omitempty"`
	Total_Data     interface{} `json:"total_data,omitempty"`
	Error_Key      string      `json:"error_key,omitempty"`
	Error_Message  string      `json:"error_message,omitempty"`
	Secondary_Data interface{} `json:"secondary_data,omitempty"`
}
