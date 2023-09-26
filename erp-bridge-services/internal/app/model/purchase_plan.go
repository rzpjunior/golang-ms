package model

import (
	"time"

	"git.edenfarm.id/edenlabs/edenlabs/orm"
)

type PurchasePlan struct {
	ID                   int64     `orm:"column(id)" json:"id"`
	Code                 string    `orm:"column(code)" json:"code"`
	VendorOrganizationID int64     `orm:"column(vendor_organization_id)" json:"vendor_organization_id"`
	SiteID               int64     `orm:"column(site_id)" json:"site_id"`
	RecognitionDate      time.Time `orm:"column(recognition_date)" json:"recognition_date"`
	EtaDate              time.Time `orm:"column(eta_date)" json:"eta_date"`
	EtaTime              string    `orm:"column(eta_time)" json:"eta_time"`
	TotalPrice           float64   `orm:"column(total_price)" json:"total_price"`
	TotalWeight          float64   `orm:"column(total_weight)" json:"total_weight"`
	TotalPurchasePlanQty float64   `orm:"column(total_purchase_plan_qty)" json:"total_purchase_plan_qty"`
	TotalPurchaseQty     float64   `orm:"column(total_purchase_qty)" json:"total_purchase_qty"`
	Note                 string    `orm:"column(note)" json:"note"`
	Status               int32     `orm:"column(status)" json:"status"`
	AssignedTo           int64     `orm:"column(assigned_to)" json:"assigned_to"`
	AssignedBy           int64     `orm:"column(assigned_by)" json:"assigned_by"`
	AssignedAt           time.Time `orm:"column(assigned_at)" json:"assigned_at"`
	CreatedAt            time.Time `orm:"column(created_at)" json:"created_at"`
	CreatedBy            int64     `orm:"column(created_by)" json:"created_by"`
}

func init() {
	orm.RegisterModel(new(PurchasePlan))
}

func (m *PurchasePlan) TableName() string {
	return "purchase_plan"
}
