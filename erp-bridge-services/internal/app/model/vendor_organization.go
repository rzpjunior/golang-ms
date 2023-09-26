package model

import "git.edenfarm.id/edenlabs/edenlabs/orm"

type VendorOrganization struct {
	ID                     int64  `orm:"column(id)" json:"id"`
	Code                   string `orm:"column(code)" json:"code"`
	VendorClassificationID int64  `orm:"column(vendor_classification_id)" json:"vendor_classification_id"`
	SubDistrictID          int64  `orm:"column(sub_district_id)" json:"sub_district_id"`
	PaymentTermID          int64  `orm:"column(payment_term_id)" json:"payment_term_id"`
	Name                   string `orm:"column(name)" json:"name"`
	Address                string `orm:"column(address)" json:"address"`
	Note                   string `orm:"column(note)" json:"note"`
	Status                 int32  `orm:"column(status)" json:"status"`
}

func init() {
	orm.RegisterModel(new(VendorOrganization))
}

func (m *VendorOrganization) TableName() string {
	return "vendor_organization"
}
