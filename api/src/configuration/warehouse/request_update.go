// Copyright 2020 PT. Eden Pangan Indonesia. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

package warehouse

import (
	"regexp"
	"strconv"
	"strings"

	"git.edenfarm.id/cuxs/common"
	"git.edenfarm.id/cuxs/orm"
	"git.edenfarm.id/cuxs/validation"
	"git.edenfarm.id/project-version2/api/datastore/repository"
	"git.edenfarm.id/project-version2/api/src/auth"
	"git.edenfarm.id/project-version2/api/util"
	"git.edenfarm.id/project-version2/datamodel/model"
)

// updateRequest : struct to hold warehouse set request data
type updateRequest struct {
	ID                      int64   `json:"-" valid:"required"`
	AreaId                  string  `json:"area_id" valid:"required"`
	PicName                 string  `json:"pic_name" valid:"required"`
	PhoneNumber             string  `json:"phone_number" valid:"required"`
	AltPhoneNumber          string  `json:"alt_phone_number"`
	StreetAddress           string  `json:"street_address" valid:"required"`
	Latitude                float64 `json:"latitude" valid:"required"`
	Longitude               float64 `json:"longitude" valid:"required"`
	Note                    string  `json:"note"`
	WarehouseType           string  `json:"warehouse_type" valid:"required"`
	ParentWarehouse         string  `json:"parent_warehouse"`
	PickerStartingLatitude  float64 `json:"picker_starting_latitude"`
	PickerStartingLongitude float64 `json:"picker_starting_longitude"`
	FloorPlanLink           string  `json:"floor_plan_link"`
	HubProcessingTimeStr    string  `json:"hub_processing_time"`
	HubProcessingTime       int64   `json:"-"`

	Warehouse             *model.Warehouse `json:"-"`
	ParentWarehouseStruct *model.Warehouse `json:"-"`
	Glossary              *model.Glossary  `json:"-"`
	BinInfoExist          bool             `json:"-"`

	Session *auth.SessionData `json:"-"`
}

// Validate : function to validate warehouse request data
func (c *updateRequest) Validate() *validation.Output {
	o := &validation.Output{Valid: true}
	o1 := orm.NewOrm()
	o1.Using("read_only")
	var e error

	c.Warehouse = &model.Warehouse{ID: c.ID}
	if e = c.Warehouse.Read("ID"); e != nil {
		o.Failure("warehouse_id.invalid", util.ErrorInvalidData("warehouse"))
	}
	if c.Warehouse.BinInfo != nil {
		c.BinInfoExist = true
		if e = c.Warehouse.BinInfo.Read("ID"); e != nil {
			o.Failure("bin_info_id.invalid", util.ErrorInvalidData("bin info"))
		}
	}

	if len(c.PicName) > 30 {
		o.Failure("pic_name.invalid", util.ErrorCharLength("pic name", 30))
	}

	if len(c.PhoneNumber) < 8 {
		o.Failure("phone_number.invalid", util.ErrorCharLength("phone number", 8))
	}

	if len(c.AltPhoneNumber) > 0 && len(c.AltPhoneNumber) < 8 {
		o.Failure("alt_phone_number.invalid", util.ErrorCharLength("alt phone number", 8))
	}

	// phonenumber
	c.PhoneNumber = strings.TrimPrefix(c.PhoneNumber, "0")
	// altphonenumber
	c.AltPhoneNumber = strings.TrimPrefix(c.AltPhoneNumber, "0")

	// latitude and longitude
	var distance float64
	// multiply by 111195 to change unit measurement to m and the distance cannot be outside of the polygon more than 100m
	if e = o1.Raw("SELECT st_distance(`polygon`,POINT(?,  ?))*111195 as distance FROM adm_division_geometry  WHERE `sub_district_id` = ?", c.Longitude, c.Latitude, c.Warehouse.SubDistrict.ID).QueryRow(&distance); e != nil {
		o.Failure("latitude_longitude.invalid", util.ErrorInvalidData("latitude and longitude"))
	}
	if distance > 100 {
		o.Failure("latitude_longitude.invalid", util.InsidePolygon())
	}

	// picker starting latitude and longitude
	if e = o1.Raw("SELECT st_distance(`polygon`,POINT(?,  ?))*111195 as distance FROM adm_division_geometry  WHERE `sub_district_id` = ?", c.PickerStartingLongitude, c.PickerStartingLatitude, c.Warehouse.SubDistrict.ID).QueryRow(&distance); e != nil {
		o.Failure("picking_latitude_longitude.invalid", util.ErrorInvalidData("picking latitude and longitude"))
	}
	if distance > 100 {
		o.Failure("picking_latitude_longitude.invalid", util.InsidePolygon())
	}

	if len(c.Note) > 100 {
		o.Failure("note", util.ErrorCharLength("note", 100))
	}

	// warehouse type and parent validation
	if c.Glossary, e = repository.GetGlossaryMultipleValue("table", "warehouse", "attribute", "warehouse_type", "value_name", c.WarehouseType); e != nil {
		o.Failure("warehouse_type.invalid", util.ErrorInvalidData("warehouse type"))
		return o
	}
	if c.Glossary.ValueName == "HUB" {
		warehouseId, e := common.Decrypt(c.ParentWarehouse)
		if e != nil {
			o.Failure("parent_warehouse_id.invalid", util.ErrorInvalidData("parent warehouse"))
		}
		if c.ParentWarehouseStruct, e = repository.ValidWarehouse(warehouseId); e != nil {
			o.Failure("parent_warehouse_id.invalid", util.ErrorInvalidData("parent warehouse"))
		}

		serviceTimeArr := strings.Split(c.HubProcessingTimeStr, ":")
		hoursCheck := regexp.MustCompile(`^[0-9]*$`).MatchString(serviceTimeArr[0])
		if hoursCheck == false {
			o.Failure("hub_processing_time_hours.invalid", util.ErrorInvalidData("hub processing time hours"))
		}
		minutesCheck := regexp.MustCompile(`^[0-9]*$`).MatchString(serviceTimeArr[1])
		if minutesCheck == false {
			o.Failure("hub_processing_time_minutes.invalid", util.ErrorInvalidData("hub processing time minutes"))
		}

		hours, _ := strconv.Atoi(serviceTimeArr[0])
		minutes, _ := strconv.Atoi(serviceTimeArr[1])
		c.HubProcessingTime = (int64(hours) * 3600) + int64(minutes)*60
		if c.HubProcessingTime <= 0 {
			o.Failure("hub_processing_time.invalid", util.ErrorInvalidData("hub processing time"))
		}
	}

	return o
}

// Messages : function to return error validation messages
func (c *updateRequest) Messages() map[string]string {
	return map[string]string{
		"area_id.required":        util.ErrorInputRequired("area id"),
		"pic_name.required":       util.ErrorInputRequired("pic name"),
		"phone_number.required":   util.ErrorInputRequired("phone number"),
		"street_address.required": util.ErrorInputRequired("address"),
		"latitude.required":       util.ErrorInputRequired("latitude"),
		"longitude.required":      util.ErrorInputRequired("longitude"),
		"warehouse_type.required": util.ErrorInputRequired("warehouse type"),
	}
}
