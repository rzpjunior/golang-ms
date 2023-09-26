package model

import (
	"time"

	"git.edenfarm.id/edenlabs/edenlabs/orm"
)

type PurchaseOrder struct {
	ID                     int64     `orm:"column(id)" json:"id"`
	Code                   string    `orm:"column(code)" json:"code"`
	VendorID               int64     `orm:"column(vendor_id)" json:"vendor_id"`
	SiteID                 int64     `orm:"column(site_id)" json:"site_id"`
	TermPaymentPurID       int64     `orm:"column(term_payment_pur_id)" json:"term_payment_pur_id"`
	VendorClassificationID int64     `orm:"column(vendor_classification_id)" json:"vendor_classification_id"`
	PurchasePlanID         int64     `orm:"column(purchase_plan_id)" json:"purchase_plan_id"`
	ConsolidatedShipmentID int64     `orm:"column(consolidated_shipment_id)" json:"consolidated_shipment_id"`
	Status                 int32     `orm:"column(status)" json:"status"`
	RecognitionDate        time.Time `orm:"column(recognition_date)" json:"recognition_date"`
	EtaDate                time.Time `orm:"column(eta_date)" json:"eta_date"`
	SiteAddress            string    `orm:"column(site_address)" json:"site_address"`
	EtaTime                string    `orm:"column(eta_time)" json:"eta_time"`
	TaxPct                 float64   `orm:"column(tax_pct)" json:"tax_pct"`
	DeliveryFee            float64   `orm:"column(delivery_fee)" json:"delivery_fee"`
	TotalPrice             float64   `orm:"column(total_price)" json:"total_price"`
	TaxAmount              float64   `orm:"column(tax_amount)" json:"tax_amount"`
	TotalCharge            float64   `orm:"column(total_charge)" json:"total_charge"`
	TotalInvoice           float64   `orm:"column(total_invoice)" json:"total_invoice"`
	TotalWeight            float64   `orm:"column(total_weight)" json:"total_weight"`
	Note                   string    `orm:"column(note)" json:"note"`
	DeltaPrint             int32     `orm:"column(delta_print)" json:"delta_print"`
	Latitude               float64   `orm:"column(latitude)" json:"latitude"`
	Longitude              float64   `orm:"column(longitude)" json:"longitude"`
	UpdatedAt              time.Time `orm:"column(updated_at)" json:"updated_at"`
	UpdatedBy              int64     `orm:"column(updated_by)" json:"updated_by"`
	CreatedFrom            int32     `orm:"column(created_from)" json:"created_from"`
	HasFinishedGr          int8      `orm:"column(has_finished_gr)" json:"has_finished_gr"`
	CreatedAt              time.Time `orm:"column(created_at)" json:"created_at"`
	CreatedBy              int64     `orm:"column(created_by)" json:"created_by"`
	CommittedAt            time.Time `orm:"column(committed_at)" json:"committed_at"`
	CommittedBy            int64     `orm:"column(committed_by)" json:"committed_by"`
	AssignedTo             int64     `orm:"column(assigned_to)" json:"assigned_to"`
	AssignedBy             int64     `orm:"column(assigned_by)" json:"assigned_by"`
	AssignedAt             time.Time `orm:"column(assigned_at)" json:"assigned_at"`
	Locked                 int32     `orm:"column(locked)" json:"locked"`
	LockedBy               int64     `orm:"column(locked_by)" json:"locked_by"`

	PurchaseOrderItem []*PurchaseOrderItem `orm:"-" json:"purchase_order_item,omitempty"`
}

func init() {
	orm.RegisterModel(new(PurchaseOrder))
}

func (m *PurchaseOrder) TableName() string {
	return "purchase_order"
}
