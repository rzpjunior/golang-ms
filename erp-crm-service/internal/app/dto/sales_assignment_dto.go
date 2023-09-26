package dto

import "time"

type SalesAssignmentResponse struct {
	ID                  int64                          `json:"id"`
	Code                string                         `json:"code"`
	Territory           *TerritoryResponse             `json:"territory"`
	StartDate           time.Time                      `json:"start_date"`
	EndDate             time.Time                      `json:"end_date"`
	Status              int8                           `json:"status"`
	StatusConvert       string                         `json:"status_convert"`
	SalesAssignmentItem []*SalesAssignmentItemResponse `json:"sales_assignment_item,omitempty"`
}

type TerritoryResponse struct {
	ID          string `json:"id"`
	Code        string `json:"code"`
	Description string `json:"description"`
}

type SalesAssignmentTemplate struct {
	TerritoryCode   string
	TerritoryName   string
	CustomerType    string
	CustomerCode    string
	CustomerName    string
	SubDistrict     string
	District        string
	SalespersonCode string
	SalespersonName string
	Task            string
	VisitDate       string
	ObjectiveCode   string
}

type SalesAssignmentExportResponse struct {
	Url string `json:"url"`
}

type SalesAssignmentImportRequest struct {
	TerritoryCode string               `json:"territory_code" valid:"required"`
	Assignments   []*AssignmentRequest `json:"assignments" valid:"required"`
	StartDate     time.Time            `json:"-"`
	EndDate       time.Time            `json:"-"`
}

type AssignmentRequest struct {
	TerritoryCode         string `json:"territory_code" valid:"required"`
	CustomerCode          string `json:"customer_code" valid:"required"`
	AddressID             string `json:"-"`
	CustomerAcquisitionID int64  `json:"-"`
	SalespersonCode       string `json:"salesperson_code" valid:"required"`
	SalespersonID         string `json:"-"`
	CustomerType          string `json:"customer_type" valid:"required"`
	CustomerTypeValue     int8   `json:"-"`
	Task                  string `json:"task" valid:"required"`
	TaskValue             int8   `json:"-"`
	VisitDate             string `json:"visit_date" valid:"required"`
	ObjectiveCodes        string `json:"objective_codes"`
}

type CheckoutTaskRequest struct {
	Id                  int64 `json:"id" valid:"required"`
	Task                int8  `json:"task" valid:"required"`
	CustomerAcquisition bool  `json:"customer_acquisition"`
}

type BulkCheckoutTaskRequest struct {
	SalesPersonID string `json:"id" valid:"required"`
}
