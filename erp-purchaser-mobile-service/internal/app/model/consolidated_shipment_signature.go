package model

import (
	"time"

	"git.edenfarm.id/edenlabs/edenlabs/orm"
)

type ConsolidatedShipmentSignature struct {
	ID                     int64     `orm:"column(id)" json:"id"`
	ConsolidatedShipmentID int64     `orm:"column(consolidated_shipment_id)" json:"consolidated_shipment_id"`
	JobFunction            string    `orm:"column(job_function)" json:"job_function"`
	Name                   string    `orm:"column(name)" json:"name"`
	SignatureURL           string    `orm:"column(signature_url)" json:"signature_url"`
	CreatedAt              time.Time `orm:"column(created_at)" json:"created_at"`
	CreatedBy              int64     `orm:"column(created_by)" json:"created_by"`
}

func init() {
	orm.RegisterModel(new(ConsolidatedShipmentSignature))
}

func (m *ConsolidatedShipmentSignature) TableName() string {
	return "consolidated_shipment_signature"
}
