package model

import (
	"git.edenfarm.id/edenlabs/edenlabs/orm"
)

type PurchaseOrder struct {
	ID                     int64  `orm:"column(id)" json:"id"`
	PurchaseOrderIDGP      string `orm:"column(purchase_order_id_gp)" json:"purchase_order_id_gp"`
	ConsolidatedShipmentID int64  `orm:"column(consolidated_shipment_id)" json:"consolidated_shipment_id"`
	SiteIDGP               string `orm:"column(site_id_gp)" json:"site_id_gp"`
	DeltaPrint             int32  `orm:"column(delta_print)" json:"delta_print"`
	// Latitude               float64   `orm:"column(latitude)" json:"latitude"`
	// Longitude              float64   `orm:"column(longitude)" json:"longitude"`
	// UpdatedAt              time.Time `orm:"column(updated_at)" json:"updated_at"`
	// UpdatedBy              int64     `orm:"column(updated_by)" json:"updated_by"`
	// CreatedAt              time.Time `orm:"column(created_at)" json:"created_at"`
	// CreatedBy              int64     `orm:"column(created_by)" json:"created_by"`
	// CommittedAt            time.Time `orm:"column(committed_at)" json:"committed_at"`
	// CommittedBy            int64     `orm:"column(committed_by)" json:"committed_by"`
	// AssignedTo             int64     `orm:"column(assigned_to)" json:"assigned_to"`
	// AssignedBy             int64     `orm:"column(assigned_by)" json:"assigned_by"`
	// AssignedAt             time.Time `orm:"column(assigned_at)" json:"assigned_at"`
	// Locked                 int32     `orm:"column(locked)" json:"locked"`
	// LockedBy               int64     `orm:"column(locked_by)" json:"locked_by"`
}

func init() {
	orm.RegisterModel(new(PurchaseOrder))
}

func (m *PurchaseOrder) TableName() string {
	return "purchase_order"
}
