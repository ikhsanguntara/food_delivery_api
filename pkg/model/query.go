package model

type QueryPagination struct {
	Page     int `form:"page" binding:"required"`
	PageSize int `form:"page_size" binding:"required"`
}

type QueryGetTransactions struct {
	StartDate string `form:"start_date"`
	EndDate   string `form:"end_date"`
	Customer  string `form:"customer"`
	QueryPagination
}

type MetodologiFilter struct {
	StartDate string `form:"start_date"`
	EndDate   string `form:"end_date"`
	Material  string `form:"material"`
}

type ResponseMetodologi struct {
	Month       string  `form:"month"`
	TotalQty    float64 `form:"total_qty"`
	Predictions float64 `form:"predictions"`
}
