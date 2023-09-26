package dto

import (
	"time"
)

type CustomerAcquisitionItemResponse struct {
	ID                    int64         `json:"id"`
	CustomerAcquisitionID int64         `json:"customer_acquisition_id"`
	Item                  *ItemResponse `json:"item"`
	IsTop                 int8          `json:"is_top"`
	CreatedAt             time.Time     `json:"created_at"`
	UpdatedAt             time.Time     `json:"updated_at"`
}

type ItemResponse struct {
	ID          string `json:"id"`
	Code        string `json:"code,omitempty"`
	Description string `json:"description,omitempty"`
}
