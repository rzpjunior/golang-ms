package model

import (
	"time"

	"git.edenfarm.id/edenlabs/edenlabs/orm"
)

type Address struct {
	ID               int64     `orm:"column(id)" json:"id"`
	Code             string    `orm:"column(code)" json:"code"`
	CustomerName     string    `orm:"column(customer_name)" json:"customer_name"`
	ArchetypeID      int64     `orm:"column(archetype_id)" json:"archetype_id"`
	AdmDivisionID    int64     `orm:"column(adm_division_id)" json:"adm_division_id"`
	SiteID           int64     `orm:"column(site_id)" json:"site_id"`
	SalespersonID    int64     `orm:"column(salesperson_id)" json:"salesperson_id"`
	TerritoryID      int64     `orm:"column(territory_id)" json:"territory_id"`
	AddressCode      string    `orm:"column(address_code)" json:"address_code"`
	AddressName      string    `orm:"column(address_name)" json:"address_name"`
	ContactPerson    string    `orm:"column(contact_person)" json:"contact_person"`
	City             string    `orm:"column(city)" json:"city"`
	State            string    `orm:"column(state)" json:"state"`
	ZipCode          string    `orm:"column(zip_code)" json:"zip_code"`
	CountryCode      string    `orm:"column(country_code)" json:"country_code"`
	Country          string    `orm:"column(country)" json:"country"`
	Latitude         float64   `orm:"column(latitude)" json:"latitude"`
	Longitude        float64   `orm:"column(longitude)" json:"longitude"`
	UpsZone          string    `orm:"column(ups_zone)" json:"ups_zone"`
	ShippingMethod   string    `orm:"column(shipping_method)" json:"shipping_method"`
	TaxScheduleID    int64     `orm:"column(tax_schedule_id)" json:"tax_schedule_id"`
	PrintPhoneNumber int8      `orm:"column(print_phone_number)" json:"print_phone_number"`
	Phone1           string    `orm:"column(phone_1)" json:"phone_1"`
	Phone2           string    `orm:"column(phone_2)" json:"phone_2"`
	Phone3           string    `orm:"column(phone_3)" json:"phone_3"`
	FaxNumber        string    `orm:"column(fax_number)" json:"fax_number"`
	ShippingAddress  string    `orm:"column(shipping_address)" json:"shipping_address"`
	BcaVa            string    `orm:"column(bca_va)" json:"bca_va"`
	OtherVa          string    `orm:"column(other_va)" json:"other_va"`
	Note             string    `orm:"column(note)" json:"note"`
	Status           int8      `orm:"column(status)" json:"status"`
	CreatedAt        time.Time `orm:"column(created_at)" json:"created_at"`
	UpdatedAt        time.Time `orm:"column(updated_at)" json:"updated_at"`
}

func init() {
	orm.RegisterModel(new(Address))
}

func (m *Address) TableName() string {
	return "address"
}
