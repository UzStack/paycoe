package domain

type Response struct {
	Status bool `json:"status"`
	Data   any  `json:"data"`
}

type Detail struct {
	Detail string `json:"detail"`
}
