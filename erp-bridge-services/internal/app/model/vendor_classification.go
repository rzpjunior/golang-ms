package model

import "git.edenfarm.id/edenlabs/edenlabs/orm"

type VendorClassification struct {
	ID            int64  `orm:"column(id)" json:"id"`
	CommodityCode string `orm:"column(commodity_code)" json:"commodity_code"`
	CommodityName string `orm:"column(commodity_name)" json:"commodity_name"`
	BadgeCode     string `orm:"column(badge_code)" json:"badge_code"`
	BadgeName     string `orm:"column(badge_name)" json:"badge_name"`
	TypeCode      string `orm:"column(type_code)" json:"type_code"`
	TypeName      string `orm:"column(type_name)" json:"type_name"`
}

func init() {
	orm.RegisterModel(new(VendorClassification))
}

func (m *VendorClassification) TableName() string {
	return "vendor_classification"
}
