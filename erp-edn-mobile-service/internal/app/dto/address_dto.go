package dto

import "time"

type AddressResponse struct {
	ID               int64     `json:"id"`
	Code             string    `json:"code"`
	CustomerName     string    `json:"customer_name"`
	ArchetypeID      string    `json:"archetype_id"`
	AdmDivisionID    string    `json:"adm_division_id"`
	SiteID           string    `json:"site_id"`
	SalespersonID    string    `json:"salesperson_id"`
	TerritoryID      string    `json:"territory_id"`
	AddressCode      string    `json:"address_code"`
	AddressName      string    `json:"address_name"`
	ContactPerson    string    `json:"contact_person"`
	City             string    `json:"city"`
	State            string    `json:"state"`
	ZipCode          string    `json:"zip_code"`
	CountryCode      string    `json:"country_code"`
	Country          string    `json:"country"`
	Latitude         *float64  `json:"latitude"`
	Longitude        *float64  `json:"longitude"`
	UpsZone          string    `json:"ups_zone"`
	ShippingMethod   string    `json:"shipping_method"`
	TaxScheduleID    int64     `json:"tax_schedule_id"`
	PrintPhoneNumber int8      `json:"print_phone_number"`
	Phone1           string    `json:"phone_1"`
	Phone2           string    `json:"phone_2"`
	Phone3           string    `json:"phone_3"`
	FaxNumber        string    `json:"fax_number"`
	ShippingAddress  string    `json:"shipping_address"`
	BcaVa            string    `json:"bca_va"`
	OtherVa          string    `json:"other_va"`
	Note             string    `json:"note"`
	DistrictId       string    `json:"district_id"`
	Status           int8      `json:"status"`
	CreatedAt        time.Time `json:"created_at"`
	UpdatedAt        time.Time `json:"updated_at"`

	Customer *CustomerResponse `json:"customer"`
}

type AddressListResponse struct {
	Data []*AddressResponse `json:"data"`
}

type AddressListRequest struct {
	Limit   int32  `json:"limit"`
	Offset  int32  `json:"offset"`
	Status  int32  `json:"status"`
	Search  string `json:"search"`
	OrderBy string `json:"order_by"`
}

type AddressDetailRequest struct {
	Id string `json:"id"`
}

type AddressGP struct {
	Custnmbr                string  `json:"custnmbr"`
	Custname                string  `json:"custname"`
	Adrscode                string  `json:"adrscode"`
	Slprsnid                string  `json:"slprsnid"`
	Shipmthd                string  `json:"shipmthd"`
	Taxschid                string  `json:"taxschid"`
	Cntcprsn                string  `json:"cntcprsn"`
	AddresS1                string  `json:"addresS1"`
	AddresS2                string  `json:"addresS2"`
	AddresS3                string  `json:"addresS3"`
	Country                 string  `json:"country"`
	City                    string  `json:"city"`
	State                   string  `json:"state"`
	Zip                     string  `json:"zip"`
	PhonE1                  string  `json:"phonE1"`
	PhonE2                  string  `json:"phonE2"`
	PhonE3                  string  `json:"phonE3"`
	CCode                   string  `json:"cCode"`
	Locncode                string  `json:"locncode"`
	Salsterr                string  `json:"salsterr"`
	UserdeF1                string  `json:"userdeF1"`
	UserdeF2                string  `json:"userdeF2"`
	ShipToName              string  `json:"shipToName"`
	GnL_Administrative_Code string  `json:"gnL_Administrative_Code"`
	GnL_Archetype_ID        string  `json:"gnL_Archetype_ID"`
	GnL_Longitude           float64 `json:"gnL_Longitude"`
	GnL_Latitude            float64 `json:"gnL_Latitude"`
	GnL_Address_Note        string  `json:"gnL_Address_Note"`
	Inactive                int32   `json:"inactive"`
	Crusrid                 string  `json:"crusrid"`
	Creatddt                string  `json:"creatddt"`
	Mdfusrid                string  `json:"mdfusrid"`
	Modifdt                 string  `json:"modifdt"`
	TypeAddress             string  `json:"type_address"`

	Customer *CustomerGP `json:"customer"`
}

type GetAddressGPResponse struct {
	PageNumber   int32        `json:"pageNumber"`
	PageSize     int32        `json:"pageSize"`
	TotalPages   int32        `json:"totalPages"`
	TotalRecords int32        `json:"totalRecords"`
	Data         []*AddressGP `json:"data"`
	Succeeded    bool         `json:"succeeded"`
	Errors       []string     `json:"errors"`
	Message      string       `json:"message"`
}

type GetAddressGPListRequest struct {
	Limit   int32  `query:"limit"`
	Offset  int32  `query:"offset"`
	Status  int32  `query:"status"`
	Search  string `query:"search"`
	OrderBy string `query:"orderBy"`
}

type Address struct {
	// ID              int64     `orm:"column(id);auto" json:"-"`
	Code            string    `orm:"column(code)" json:"code,omitempty"`
	Name            string    `orm:"column(name)" json:"name,omitempty"`
	PicName         string    `orm:"column(pic_name)" json:"pic_name,omitempty"`
	PhoneNumber     string    `orm:"column(phone_number)" json:"phone_number,omitempty"`
	AltPhoneNumber  string    `orm:"column(alt_phone_number)" json:"alt_phone_number,omitempty"`
	AddressName     string    `orm:"column(address_name)" json:"address_name"`
	ShippingAddress string    `orm:"column(shipping_address)" json:"shipping_address,omitempty"`
	Latitude        *float64  `orm:"column(latitude)" json:"latitude,omitempty"`
	Longitude       *float64  `orm:"column(longitude)" json:"longitude,omitempty"`
	Note            string    `orm:"column(note)" json:"note,omitempty"`
	MainBranch      int8      `orm:"column(main_branch)" json:"main_branch,omitempty"`
	Status          int8      `orm:"column(status)" json:"status"`
	CreatedAt       time.Time `orm:"column(created_at);type(timestamp);null" json:"created_at"`
	CreatedBy       int64     `orm:"column(created_by)" json:"created_by"`
	LastUpdatedAt   time.Time `orm:"column(last_updated_at);type(timestamp);null" json:"last_updated_at"`
	LastUpdatedBy   int64     `orm:"column(last_updated_by)" json:"last_updated_by"`
	City            string    `json:"city"`
	State           string    `json:"state"`
	ZipCode         string    `json:"zip_code"`
	CountryCode     string    `json:"country_code"`
	Country         string    `json:"country"`
	UpsZone         string    `json:"ups_zone"`

	CustomerID   string `json:"customer_id,omitempty"`
	CustomerName string `json:"customer_name,omitempty"`
	RegionID     string `json:"region_id,omitempty"`
	ArchetypeID  string `json:"archetype_id,omitempty"`
	SiteID       string `json:"site_id,omitempty"`
	// Salesperson *Staff       `orm:"column(salesperson_id);null" json:"salesperson,omitempty"`
	AdmDivisionID string `orm:"column(sub_district_id);null" json:"adm_division_id,omitempty"`
	ContactPerson string `json:"contact_person"`

	StatusConvert string `orm:"-" json:"status_convert"`
}
