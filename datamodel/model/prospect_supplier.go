package model

import (
	"encoding/json"

	"git.edenfarm.id/cuxs/common"
	"git.edenfarm.id/cuxs/orm"
)

func init() {
	orm.RegisterModel(new(ProspectSupplier))
}

// PriceSet : struct to hold price set model data for database
type ProspectSupplier struct {
	ID             int64  `orm:"column(id);auto" json:"-"`
	Code           string `orm:"column(code)" json:"code"`
	Name           string `orm:"column(name)" json:"name"`
	PhoneNumber    string `orm:"column(phone_number);size(15);null" json:"phone_number"`
	AltPhoneNumber string `orm:"column(alt_phone_number);size(15);null" json:"alt_phone_number"`
	StreetAddress  string `orm:"column(street_address);size(350);null" json:"street_address"`
	PicName        string `orm:"column(pic_name);size(100);null" json:"pic_name"`
	PicPhoneNumber string `orm:"column(pic_phone_number);size(15);null" json:"pic_phone_number"`
	Commodity      string `orm:"column(commodity)" json:"commodity"`
	TimeConsent    int8   `orm:"column(time_consent)" json:"time_consent"`
	RegStatus      int8   `orm:"column(reg_status)" json:"reg_status"`
	PicAddress     string `orm:"column(pic_address)" json:"pic_address"`

	SubDistrict *SubDistrict `orm:"column(sub_district_id);null;rel(fk)" json:"sub_district,omitempty"`
}

// MarshalJSON : function to customize data struct into json. encrypt all key ids
func (m *ProspectSupplier) MarshalJSON() ([]byte, error) {
	type Alias ProspectSupplier

	return json.Marshal(&struct {
		ID string `json:"id"`
		*Alias
	}{
		ID:    common.Encrypt(m.ID),
		Alias: (*Alias)(m),
	})
}

// Save : function to save data into database. will update when it has valid id, create otherwise
func (m *ProspectSupplier) Save(fields ...string) (err error) {
	o := orm.NewOrm()
	if m.ID > 0 {
		_, err = o.Update(m, fields...)
	} else {
		m.ID, err = o.Insert(m)
	}
	return
}

// Read : function to get data from database
func (m *ProspectSupplier) Read(fields ...string) error {
	o := orm.NewOrm()
	o.Using("read_only")
	return o.Read(m, fields...)
}
