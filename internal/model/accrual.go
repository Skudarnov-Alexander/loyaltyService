package model

type Accrual struct {
	Number  string		`json:"order"`
    Status  string		`json:"status"`
    Accrual float64		`json:"accrual"`
}