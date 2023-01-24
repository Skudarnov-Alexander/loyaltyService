package model

type Good struct {
	Description string		`json:"description"`
	Price		float64		`json:"price"`
}

type Order struct {
	Order   string	`json:"order"`
	Goods   []Good	`json:"goods"`
	Status  string	`json:"status"`
	Accrual float64	`json:"accrual"`
}
