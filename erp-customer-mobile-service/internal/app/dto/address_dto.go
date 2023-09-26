package dto

type RequestGetAddressList struct {
	Platform string `json:"platform" valid:"required"`
	//Data     dataGetAddressList `json:"data" valid:"required"`
	Session *SessionDataCustomer
}

type dataGetAddressList struct {
	AddressList []*ListAddressList
}

type ListAddressList struct {
	AddressID     string  `json:"id"`
	ArchetypeID   string  `json:"archetype_id"`
	AddressName   string  `json:"address_name"`
	PicName       string  `json:"pic_name"`
	PhoneNumber   string  `json:"phone_number"`
	AddressType   string  `json:"address_type"`
	Address1      string  `json:"address_1"`
	Address2      string  `json:"address_2"`
	Address3      string  `json:"address_3"`
	AddressNote   string  `json:"address_note"`
	Latitude      float64 `json:"latitude"`
	Longitude     float64 `json:"longitude"`
	AdmDivisionId string  `json:"adm_division_id"`
	RegionID      string  `json:"region_id"`
	Province      string  `json:"province"`
	City          string  `json:"city"`
	District      string  `json:"district"`
	SubDistrict   string  `json:"sub_district"`
}

type CreateAddressRequest struct {
	Platform string            `json:"platform" valid:"required"`
	Data     dataCreateAddress `json:"data" valid:"required"`

	Session *SessionDataCustomer
}

type dataCreateAddress struct {
	AddressName string `json:"address_name" valid:"required"`
	PICName     string `json:"pic_name" valid:"required"`
	PhoneNumber string `json:"phone_number" valid:"required"`
	Address1    string `json:"address_1" valid:"required"`
	Address2    string `json:"address_2"`
	Address3    string `json:"address_3"`
	AddressNote string `json:"address_note"`
	Latitude    string `json:"latitude"`
	Longitude   string `json:"longitude"`
	SubDistrict string `json:"sub_district" valid:"required"`
}

type UpdateAddressRequest struct {
	Platform string            `json:"platform" valid:"required"`
	Data     dataUpdateAddress `json:"data" valid:"required"`

	Session *SessionDataCustomer
}

type dataUpdateAddress struct {
	AddressID   string `json:"address_id" valid:"required"`
	AddressName string `json:"address_name" valid:"required"`
	PICName     string `json:"pic_name" valid:"required"`
	PhoneNumber string `json:"phone_number" valid:"required"`
	Address1    string `json:"address_1" valid:"required"`
	Address2    string `json:"address_2"`
	Address3    string `json:"address_3"`
	AddressNote string `json:"address_note"`
	Latitude    string `json:"latitude"`
	Longitude   string `json:"longitude"`
	SubDistrict string `json:"sub_district" valid:"required"`
}

type SetDefaultAddressRequest struct {
	Platform string         `json:"platform" valid:"required"`
	Data     dataGetDefault `json:"data" valid:"required"`

	Session *SessionDataCustomer
}

type dataGetDefault struct {
	AddressID string `json:"address_id" valid:"required"`
}

type DeleteAddressRequest struct {
	Platform string        `json:"platform" valid:"required"`
	Data     dataGetDelete `json:"data" valid:"required"`

	Session *SessionDataCustomer
}

type dataGetDelete struct {
	AddressID string `json:"address_id" valid:"required"`
}
