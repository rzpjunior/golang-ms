package dto

import "time"

type AddressResponse struct {
	ID               int64     `json:"id"`
	Code             string    `json:"code"`
	CustomerName     string    `json:"customer_name"`
	ArchetypeID      int64     `json:"archetype_id"`
	AdmDivisionID    int64     `json:"adm_division_id"`
	SiteID           int64     `json:"site_id"`
	SalespersonID    int64     `json:"salesperson_id"`
	TerritoryID      int64     `json:"territory_id"`
	AddressCode      string    `json:"address_code"`
	AddressName      string    `json:"address_name"`
	ContactPerson    string    `json:"contact_person"`
	City             string    `json:"city"`
	State            string    `json:"state"`
	ZipCode          string    `json:"zip_code"`
	CountryCode      string    `json:"country_code"`
	Country          string    `json:"country"`
	Latitude         float64   `json:"latitude"`
	Longitude        float64   `json:"longitude"`
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
	Status           int8      `json:"status"`
	DistrictId       int64     `json:"district_id"`
	StatusConvert    string    `json:"status_convert"`
	CreatedAt        time.Time `json:"created_at"`
	UpdatedAt        time.Time `json:"updated_at"`
}

type AddressRequestCreate struct {
	// InterID                 string `json:"interid"`
	// Custnmbr                string `json:"custnmbr"`
	// Custname                string `json:"custname"`
	// Adrscode                string `json:"adrscode"`
	// Slprsnid                string `json:"slprsnid"`
	// Shipmthd                string `json:"shipmthd"`
	// Taxschid                string `json:"taxschid"`
	// Cntcprsn                string `json:"cntcprsn"`
	// AddresS1                string `json:"addresS1"`
	// AddresS2                string `json:"addresS2"`
	// AddresS3                string `json:"addresS3"`
	// Country                 string `json:"country"`
	// City                    string `json:"city"`
	// State                   string `json:"state"`
	// Zip                     string `json:"zip"`
	// PhonE1                  string `json:"phonE1"`
	// PhonE2                  string `json:"phonE2"`
	// PhonE3                  string `json:"phonE3"`
	// CCode                   string `json:"cCode"`
	// Locncode                string `json:"locncode"`
	// Salsterr                string `json:"salsterr"`
	// UserdeF1                string `json:"userdeF1"`
	// UserdeF2                string `json:"userdeF2"`
	// ShipToName              string `json:"shipToName"`
	// GnL_Administrative_Code string `json:"gnL_Administrative_Code"`
	// GnL_Archetype_ID        string `json:"gnL_Archetype_ID"`
	// GnL_Longitude           string `json:"gnL_Longitude"`
	// GnL_Latitude            string `json:"gnL_Latitude"`
	// GnL_Address_Note        string `json:"gnL_Address_Note"`
	// Inactive                string `json:"inactive"`
	// Crusrid                 string `json:"crusrid"`
	// Creatddt                string `json:"creatddt"`
	// Mdfusrid                string `json:"mdfusrid"`
	// Modifdt                 string `json:"modifdt"`
	// Type_address            string `json:"type_address"`
	InterID               string  `json:"interid"`
	CustNmbr              string  `json:"custnmbr"`
	AdrsCode              string  `json:"adrscode"`
	CntcPrsn              string  `json:"cntcprsn"`
	Address1              string  `json:"addresS1"`
	Address2              string  `json:"addresS2"`
	Address3              string  `json:"addresS3"`
	City                  string  `json:"city"`
	State                 string  `json:"state"`
	Zip                   string  `json:"zip"`
	CCode                 string  `json:"cCode"`
	Country               string  `json:"country"`
	GnlAdministrativeCode string  `json:"gnL_Administrative_Code"`
	GnlArchetypeID        string  `json:"gnL_Archetype_ID"`
	Upszone               string  `json:"upszone"`
	Shipmthd              string  `json:"shipmthd"`
	Taxschid              string  `json:"taxschid"`
	Locncode              string  `json:"locncode"`
	Slprsnid              string  `json:"slprsnid"`
	Salsterr              string  `json:"salsterr"`
	GnlLongitude          float64 `json:"gnL_Longitude"`
	GnlLatitude           float64 `json:"gnL_Latitude"`
	Userdef1              string  `json:"userdeF1"`
	Userdef2              string  `json:"userdeF2"`
	ShipToName            string  `json:"shipToName"`
	Phone1                string  `json:"phonE1"`
	Phone2                string  `json:"phonE2"`
	Phone3                string  `json:"phonE3"`
	Fax                   string  `json:"fax"`
	GnlAddressNote        string  `json:"gnL_Address_Note"`
	Param                 string  `json:"param"`
}

type UpdateAddressRequest struct {
	InterID               string `json:"interid"`
	CustNmbr              string `json:"custnmbr"`
	AdrsCode              string `json:"adrscode"`
	CntcPrsn              string `json:"cntcprsn"`
	Address1              string `json:"addresS1"`
	Address2              string `json:"addresS2"`
	Address3              string `json:"addresS3"`
	City                  string `json:"city"`
	State                 string `json:"state"`
	Zip                   string `json:"zip"`
	CCode                 string `json:"cCode"`
	Country               string `json:"country"`
	GnlAdministrativeCode string `json:"gnL_Administrative_Code"`
	GnlArchetypeID        string `json:"gnL_Archetype_ID"`
	Upszone               string `json:"upszone"`
	Shipmthd              string `json:"shipmthd"`
	Taxschid              string `json:"taxschid"`
	Locncode              string `json:"locncode"`
	Slprsnid              string `json:"slprsnid"`
	Salsterr              string `json:"salsterr"`
	GnlLongitude          string `json:"gnL_Longitude"`
	GnlLatitude           string `json:"gnL_Latitude"`
	Userdef1              string `json:"userdeF1"`
	Userdef2              string `json:"userdeF2"`
	ShipToName            string `json:"shipToName"`
	Phone1                string `json:"phonE1"`
	Phone2                string `json:"phonE2"`
	Phone3                string `json:"phonE3"`
	Fax                   string `json:"fax"`
	GnlAddressNote        string `json:"gnL_Address_Note"`
	Param                 string `json:"param"`
	Inactive              int32  `json:"inactive"`
}

type SetDefaultAddressRequest struct {
	InterID  string `json:"interid"`
	CustNmbr string `json:"custnmbr"`
	AdrsCode string `json:"adrscode"`
}

type DeleteAddressRequest struct {
	InterID  string `json:"interid"`
	CustNmbr string `json:"custnmbr"`
	AdrsCode string `json:"adrscode"`
}
