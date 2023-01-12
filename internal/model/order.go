package model

type Order struct {
	Number     string `json:"number"`
	Status     string `json:"status"`
	Accrual    int64  `json:"accrual"`
	UploadedAt string `json:"uploaded_at"`
}