package dto

import "time"

type ConsolidatedShipmentSignatureResponse struct {
	ID                   int64                         `json:"id"`
	ConsolidatedShipment *ConsolidatedShipmentResponse `json:"consolidated_shipment,omitempty"`
	JobFunction          string                        `json:"job_function"`
	Name                 string                        `json:"name"`
	SignatureURL         string                        `json:"signature_url"`
	CreatedAt            time.Time                     `json:"created_at"`
	CreatedBy            int64                         `json:"created_by"`
}

type ConsolidatedShipmentSignatureRequestCreate struct {
	ConsolidatedShipmentID int64  `json:"consolidated_shipment_id" valid:"required"`
	JobFunction            string `json:"job_function" valid:"required|alpha_num_space|lte:100"`
	Name                   string `json:"name" valid:"required|alpha_num_space|lte:100"`
	SignatureURL           string `json:"signature_url" valid:"required"`
}
