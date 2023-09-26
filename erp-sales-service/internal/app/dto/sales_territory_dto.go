package dto

import "time"

type SalesTerritoryResponse struct {
	ID            string    `json:"id"`
	Description   string    `json:"description"`
	SalespersonID string    `json:"salesperson_id"`
	CreatedAt     time.Time `json:"created_at"`
	UpdatedAt     time.Time `json:"updated_at"`
}

type GetSalesTerritoryRequest struct {
	Limit  int
	Offset int
	Search string
}
