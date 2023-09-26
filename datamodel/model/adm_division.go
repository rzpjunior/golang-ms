// Copyright 2020 PT. Eden Pangan Indonesia. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

package model

import (
	"encoding/json"

	"git.edenfarm.id/cuxs/common"
	"git.edenfarm.id/cuxs/orm"
)

func init() {
	orm.RegisterModel(new(AdmDivision))
}

// AdmDivision : struct to hold model data for database
type AdmDivision struct {
	AreaId            int64       `orm:"column(area_id)" json:"area_id"`
	AreaName          string      `orm:"column(area_name)" json:"area_name"`
	AreaStatus        int8        `orm:"column(area_status)" json:"area_status"`
	SubDistrictId     int64       `orm:"column(sub_district_id);pk" json:"sub_district_id"`
	SubDistrictName   string      `orm:"column(sub_district_name)" json:"sub_district_name"`
	SubDistrictStatus int8        `orm:"column(sub_district_status)" json:"sub_district_status"`
	DistrictId        int64       `orm:"column(district_id)" json:"district_id"`
	DistrictName      string      `orm:"column(district_name)" json:"district_name"`
	DistrictStatus    int8        `orm:"column(district_status)" json:"district_status"`
	CityId            int64       `orm:"column(city_id)" json:"city_id"`
	CityName          string      `orm:"column(city_name)" json:"city_name"`
	CityStatus        int8        `orm:"column(city_status)" json:"city_status"`
	ProvinceId        int64       `orm:"column(province_id)" json:"province_id"`
	ProvinceName      string      `orm:"column(province_name)" json:"province_name"`
	ProvinceStatus    int8        `orm:"column(province_status)" json:"province_status"`
	CountryId         int64       `orm:"column(country_id)" json:"country_id"`
	CountryName       string      `orm:"column(country_name)" json:"country_name"`
	CountryStatus     int8        `orm:"column(country_status)" json:"country_status"`
	PostalCode        string      `orm:"column(postal_code)" json:"postal_code"`
	ConcateAddress    string      `orm:"column(concate_address)" json:"concate_address"`
	Centroid          Centroid    `orm:"-" json:"centroid"`
	Polygon           [][]float64 `orm:"-" json:"polygon"`
}

type Centroid struct {
	Latitude  float64 `json:"latitude"`
	Longitude float64 `json:"longitude"`
}

type PolygonJson struct {
	Type        string        `json:"type"`
	Coordinates [][][]float64 `json:"coordinates"`
}

// MarshalJSON : function to customize data struct into json. encrypt all key ids
func (m *AdmDivision) MarshalJSON() ([]byte, error) {
	type Alias AdmDivision

	return json.Marshal(&struct {
		AreaId        string `json:"area_id"`
		SubDistrictId string `json:"sub_district_id"`
		DistrictId    string `json:"district_id"`
		CityId        string `json:"city_id"`
		ProvinceId    string `json:"province_id"`
		CountryId     string `json:"country_id"`
		*Alias
	}{
		AreaId:        common.Encrypt(m.AreaId),
		SubDistrictId: common.Encrypt(m.SubDistrictId),
		DistrictId:    common.Encrypt(m.DistrictId),
		CityId:        common.Encrypt(m.CityId),
		ProvinceId:    common.Encrypt(m.ProvinceId),
		CountryId:     common.Encrypt(m.CountryId),
		Alias:         (*Alias)(m),
	})
}

// Read : function to get data from database
func (m *AdmDivision) Read(fields ...string) error {
	o := orm.NewOrm()
	o.Using("read_only")
	return o.Read(m, fields...)
}
