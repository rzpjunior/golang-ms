package model

import (
	"encoding/json"

	"time"

	"git.edenfarm.id/cuxs/orm"
)

func init() {
	orm.RegisterModel(new(OtpOutgoing))
}

// OtpOutgoing : struct to hold uom model data for database
type OtpOutgoing struct {
	ID              int64     `orm:"column(id)" json:"-"` // id not set
	PhoneNumber     string    `orm:"column(phone_number);" json:"phone_number"`
	OTP             string    `orm:"column(otp)" json:"otp"`
	Application     int       `orm:"column(application)" json:"application"`
	UsageType       int       `orm:"column(usage_type)" json:"usage_type"`
	OtpStatus       int       `orm:"column(otp_status)" json:"otp_status"`
	Vendor          int       `orm:"column(vendor)" json:"vendor"`
	VendorMessageID string    `orm:"column(vendor_message_id)" json:"vendor_message_id"`
	MessageType     int       `orm:"column(message_type)" json:"message_type"`
	Message         string    `orm:"column(message)" json:"message"`
	DeliveryStatus  int       `orm:"column(delivery_status)" json:"delivery_status"`
	CreatedAt       time.Time `orm:"column(created_at)" json:"created_at"`
	UpdatedAt       time.Time `orm:"column(updated_at)" json:"updated_at"`
}

// MarshalJSON : function to customize data struct into json. encrypt all key ids
func (m *OtpOutgoing) MarshalJSON() ([]byte, error) {
	type Alias OtpOutgoing

	return json.Marshal(&struct {
		*Alias
	}{})
}
func (m *OtpOutgoing) Save(fields ...string) (err error) {
	o := orm.NewOrm()
	if m.ID > 0 {
		_, err = o.Update(m, fields...)
	} else {
		m.ID, err = o.Insert(m)
	}
	return
}

// Read : function to get data from database
func (m *OtpOutgoing) Read(fields ...string) error {
	o := orm.NewOrm()
	o.Using("read_only")
	return o.Read(m, fields...)
}
