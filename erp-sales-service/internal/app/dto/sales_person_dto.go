package dto

type SalesPersonResponse struct {
	ID               string `json:"id"`
	Name             string `json:"name"`
	MiddleName       string `json:"middle_name"`
	LastName         string `json:"last_name"`
	EmployeeID       string `json:"employee_id"`
	SalesTerritoryID string `json:"sales_territory_id"`
	Status           int8   `json:"status"`
	ConvertStatus    string `json:"convert_status"`
}

type GetSalesPersonRequest struct {
	Limit            int
	Offset           int
	SalesTerritoryID string
	Status           int8
	Search           string
}
