package dto

import (
	"time"
)

type SalesAssignmentItemResponse struct {
	ID                    int64                               `json:"id"`
	SalesAssignmentID     *int64                              `json:"sales_assignment_id"`
	CustomerAcquisitionID int64                               `json:"customer_acquisition_id"`
	AddressID             string                              `json:"address_id"`
	Latitude              float64                             `json:"latitude"`
	Longitude             float64                             `json:"longitude"`
	Task                  int8                                `json:"task"`
	CustomerType          int8                                `json:"customer_type"`
	ObjectiveCodes        string                              `json:"objective_codes"`
	ActualDistance        float64                             `json:"actual_distance"`
	OutOfRoute            int8                                `json:"out_of_route"`
	StartDate             time.Time                           `json:"start_date"`
	EndDate               time.Time                           `json:"end_date"`
	FinishDate            *time.Time                          `json:"finish_date"`
	SubmitDate            time.Time                           `json:"submit_date"`
	TaskImageUrls         []string                            `json:"task_image_urls"`
	TaskAnswer            int8                                `json:"task_answer"`
	Status                int8                                `json:"status"`
	StatusConvert         string                              `json:"status_convert"`
	EffectiveCall         int8                                `json:"effective_call"`
	Address               *AddressResponse                    `json:"address"`
	CustomerAcquisition   *CustomerAcquisitionResponse        `json:"customer_acquisition"`
	SalesPerson           *SalespersonResponse                `json:"salesperson"`
	ObjectiveValues       []*SalesAssignmentObjectiveResponse `json:"objective_values"`
}

type SalesAssignmentItemRequest struct {
	SalesAssignmentID     *int64     `json:"sales_assignment_id"`
	CustomerAcquisitionID int64      `json:"customer_acquisition_id"`
	AddressID             string     `json:"address_id"`
	Latitude              float64    `json:"latitude"`
	Longitude             float64    `json:"longitude"`
	Task                  int8       `json:"task" valid:"required"`
	CustomerType          int8       `json:"customer_type" valid:"required"`
	ObjectiveCodes        string     `json:"objective_codes"`
	ActualDistance        float64    `json:"actual_distance"`
	OutOfRoute            int8       `json:"out_of_route"`
	StartDate             time.Time  `json:"start_date" valid:"required"`
	EndDate               time.Time  `json:"end_date" valid:"required"`
	FinishDate            *time.Time `json:"finish_date"`
	SubmitDate            time.Time  `json:"submit_date"`
	TaskImageUrls         []string   `json:"task_image_urls"`
	TaskAnswer            int8       `json:"task_answer"`
	Status                int8       `json:"status" valid:"required"`
	StatusConvert         string     `json:"status_convert"`
	EffectiveCall         int8       `json:"effective_call"`
	SalesPersonID         string     `json:"salesperson_id" valid:"required"`
}

type AddressResponse struct {
	ID   string `json:"id"`
	Code string `json:"code"`
	Name string `json:"name"`
}

type SalespersonResponse struct {
	ID   string `json:"id"`
	Code string `json:"code"`
	Name string `json:"name"`
}

type UpdateSubmitTaskVisitFURequest struct {
	SalesAssignmentItemResponse
}
