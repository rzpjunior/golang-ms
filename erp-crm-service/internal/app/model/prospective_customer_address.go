package model

import (
	"time"

	"git.edenfarm.id/edenlabs/edenlabs/orm"
)

type ProspectiveCustomerAddress struct {
	ID                    int64     `orm:"column(id)" json:"id"`
	ProspectiveCustomerID int64     `orm:"column(prospective_customer_id)" json:"prospective_customer_id"`
	AddressName           string    `orm:"column(address_name)" json:"address_name"`
	AddressType           string    `orm:"column(address_type)" json:"address_type"`
	AdmDivisionIDGP       string    `orm:"column(adm_division_id_gp)" json:"adm_division_id_gp"`
	Address1              string    `orm:"column(address_1)" json:"address_1"`
	Address2              string    `orm:"column(address_2)" json:"address_2"`
	Address3              string    `orm:"column(address_3)" json:"address_3"`
	Latitude              float64   `orm:"column(latitude)" json:"latitude"`
	Longitude             float64   `orm:"column(longitude)" json:"longitude"`
	Note                  string    `orm:"column(note)" json:"note"`
	CreatedAt             time.Time `orm:"column(created_at)" json:"created_at"`
	UpdatedAt             time.Time `orm:"column(updated_at)" json:"updated_at"`
	ReferTo               int8      `orm:"column(refer_to)" json:"refer_to"`

	City                string `orm:"-" json:"city"`
	State               string `orm:"-" json:"state"`
	PostalCode          string `orm:"-" json:"postal_code"`
	IsCreateAddressToGP bool   `orm:"-" json:"is_generate_code"`
}

func init() {
	orm.RegisterModel(new(ProspectiveCustomerAddress))
}

func (m *ProspectiveCustomerAddress) TableName() string {
	return "prospective_customer_address"
}
