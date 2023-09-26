package model

import (
	"time"

	"git.edenfarm.id/edenlabs/edenlabs/orm"
)

type Vendor struct {
	ID                     int64     `orm:"column(id)" json:"id"`
	Code                   string    `orm:"column(code)" json:"code"`
	VendorOrganizationID   int64     `orm:"column(vendor_organization_id)" json:"vendor_organization_id"`
	VendorClassificationID int64     `orm:"column(vendor_classification_id)" json:"vendor_classification_id"`
	SubDistrictID          int64     `orm:"column(sub_district_id)" json:"sub_district_id"`
	PicName                string    `orm:"column(pic_name)" json:"pic_name"`
	Email                  string    `orm:"column(email)" json:"email"`
	PhoneNumber            string    `orm:"column(phone_number)" json:"phone_number"`
	PaymentTermID          int64     `orm:"column(payment_term_id)" json:"payment_term_id"`
	Rejectable             int32     `orm:"column(rejectable)" json:"rejectable"`
	Returnable             int32     `orm:"column(returnable)" json:"returnable"`
	Address                string    `orm:"column(address)" json:"address"`
	Note                   string    `orm:"column(note)" json:"note"`
	Status                 int32     `orm:"column(status)" json:"status"`
	Latitude               string    `orm:"column(latitude)" json:"latitude"`
	Longitude              string    `orm:"column(longitude)" json:"longitude"`
	CreatedAt              time.Time `orm:"column(created_at)" json:"created_at"`
	CreatedBy              int64     `orm:"column(created_by)" json:"created_by"`
}

func init() {
	orm.RegisterModel(new(Vendor))
}

func (m *Vendor) TableName() string {
	return "vendor"
}
