package model

import (
	"time"

	"git.edenfarm.id/edenlabs/edenlabs/orm"
)

type AdmDivision struct {
	ID            int64     `orm:"column(id)" json:"id"`
	Code          string    `orm:"column(code)" json:"code"`
	ProvinceID    int64     `orm:"column(province_id)" json:"province_id"`
	CityID        int64     `orm:"column(city_id)" json:"city_id"`
	DistrictID    int64     `orm:"column(district_id)" json:"district_id"`
	SubDistrictID int64     `orm:"column(sub_district_id)" json:"sub_district_id"`
	RegionID      int64     `orm:"column(region_id)" json:"region_id"`
	PostalCode    string    `orm:"column(postal_code)" json:"postal_code"`
	Province      string    `orm:"column(province)" json:"province"`
	City          string    `orm:"column(city)" json:"city"`
	District      string    `orm:"column(district)" json:"district"`
	Region        string    `orm:"column(region)" json:"region"`
	Status        int8      `orm:"column(status)" json:"status"`
	CreatedAt     time.Time `orm:"column(created_at)" json:"created_at"`
	UpdatedAt     time.Time `orm:"column(updated_at)" json:"updated_at"`
}

func init() {
	orm.RegisterModel(new(AdmDivision))
}

func (m *AdmDivision) TableName() string {
	return "adm_division"
}
