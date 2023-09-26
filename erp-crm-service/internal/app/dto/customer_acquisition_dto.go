package dto

import (
	"time"
)

type CustomerAcquisitionResponse struct {
	ID                       int64                              `json:"id"`
	Code                     string                             `json:"code"`
	Task                     int8                               `json:"task"`
	Name                     string                             `json:"name"`
	PhoneNumber              string                             `json:"phone_number"`
	Latitude                 float64                            `json:"latitude"`
	Longitude                float64                            `json:"longitude"`
	AddressName              string                             `json:"address_name"`
	FoodApp                  int8                               `json:"food_app"`
	PotentialRevenue         float64                            `json:"potential_revenue"`
	TaskImageUrl             string                             `json:"task_image_url"`
	Salesperson              *SalespersonResponse               `json:"salesperson"`
	Territory                *TerritoryResponse                 `json:"territory"`
	FinishDate               time.Time                          `json:"finish_date"`
	SubmitDate               time.Time                          `json:"submit_date"`
	CreatedAt                time.Time                          `json:"created_at"`
	UpdatedAt                time.Time                          `json:"updated_at"`
	Status                   int8                               `json:"status"`
	StatusConvert            string                             `json:"status_convert"`
	CustomerAcquisitionItems []*CustomerAcquisitionItemResponse `json:"customer_acquisition_items,omitempty"`
}

type SubmitTaskCustomerAcqRequest struct {
	SalesPersonID            string                `json:"salesperson_id"`
	CustomerName             string                `json:"customer_name" valid:"required"`
	PhoneNumber              string                `json:"phone_number" valid:"required"`
	AddressDetail            string                `json:"address_detail" valid:"required"`
	FoodApp                  int8                  `json:"food_app" valid:"required"`
	UserLatitude             float64               `json:"user_lats" valid:"required"`
	UserLongitude            float64               `json:"user_longs" valid:"required"`
	PotentialRevenue         float64               `json:"potential_revenue" valid:"required"`
	CustomerAcquisitionPhoto string                `json:"customer_acquisition_photo" valid:"required"`
	Product                  []*CustomerAcqProduct `json:"products" valid:"required"`
}

type CustomerAcqProduct struct {
	Id  int64 `json:"id" valid:"required"`
	Top int8  `json:"is_top"`
}
