package model

import (
	"time"
)

type ItemTransfer struct {
	ID                 int64     `json:"id"`
	Code               string    `json:"code"`
	RequestDate        time.Time `json:"request_date"`
	RecognitionDate    time.Time `json:"recognition_date"`
	EtaDate            time.Time `json:"eta_date"`
	EtaTime            string    `json:"eta_time"`
	AtaDate            time.Time `json:"ata_date"`
	AtaTime            string    `json:"ata_time"`
	AdditionalCost     float64   `json:"additional_cost"`
	AdditionalCostNote string    `json:"additional_cost_note"`
	StockType          int8      `json:"stock_type"`
	TotalCost          float64   `json:"total_cost"`
	TotalCharge        float64   `json:"total_charge"`
	TotalWeight        float64   `json:"total_weight"`
	Note               string    `json:"note"`
	Status             int8      `json:"status"`
	Locked             int8      `json:"locked"`
	LockedBy           int64     `json:"locked_by"`
	TotalSku           int64     `json:"total_sku"`
	HasFinishedGr      int8      `json:"has_finished_gr"`
	SiteOriginID       int64     `json:"site_origin_id"`
	SiteDestinationID  int64     `json:"site_destination_id"`
	UpdatedAt          time.Time `json:"updated_at"`
	UpdatedBy          int64     `json:"updated_by"`

	ItemTransferItem []*ItemTransferItem `json:"item_transfer_item,omitempty"`
}

// TODO: init to db/gp
