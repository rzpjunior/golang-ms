// Copyright 2020 PT. Eden Pangan Indonesia. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

package warehouse

import (
	"strconv"

	"git.edenfarm.id/cuxs/orm"
	"git.edenfarm.id/project-version2/api/datastore/repository"
	"git.edenfarm.id/project-version2/api/log"
	"git.edenfarm.id/project-version2/api/util"
	"git.edenfarm.id/project-version2/datamodel/model"
)

// Save : function to save data requested into database
func Save(r createRequest) (u *model.Warehouse, e error) {
	o := orm.NewOrm()
	o.Begin()

	r.Code, e = util.GenerateCode(r.Code, "warehouse")

	binInfo := &model.BinInfo{
		Latitude:  &r.PickerStartingLatitude,
		Longitude: &r.PickerStartingLongitude,
		ImageURL:  r.FloorPlanLink,
	}

	_, e = o.Insert(binInfo)
	if e != nil {
		o.Rollback()
		return nil, e
	}

	e = log.AuditLogByUser(r.Session.Staff, binInfo.ID, "bin", "create", "")
	if e != nil {
		o.Rollback()
		return nil, e
	}

	u = &model.Warehouse{
		Code:                   r.Code,
		Name:                   r.Name,
		PicName:                r.PicName,
		PhoneNumber:            r.PhoneNumber,
		AltPhoneNumber:         r.AltPhoneNumber,
		StreetAddress:          r.StreetAddress,
		Latitude:               r.Latitude,
		Longitude:              r.Longitude,
		Note:                   r.Note,
		AuxData:                2,
		Status:                 int8(1),
		LimitOrderPickingList:  4,
		LimitWeightPickingList: 80,
		WarehouseType:          r.Glossary.ValueInt,
		BinInfo:                binInfo,
		Area:                   r.Area,
		SubDistrict:            r.SubDistrict,
	}
	if r.Glossary.ValueName == "HUB" {
		u.HubProcessingTime = r.HubProcessingTime
		u.ParentID = r.ParentWarehouseStruct
	}

	_, e = o.Insert(u)
	if e != nil {
		o.Rollback()
		return nil, e
	}

	e = log.AuditLogByUser(r.Session.Staff, u.ID, "warehouse", "create", "")
	if e != nil {
		o.Rollback()
		return nil, e
	}

	// stock and product config
	if r.CopyWarehouse == true {
		if e = WarehouseCopyStock(r.StockData, u, o); e != nil {
			o.Rollback()
			return nil, e
		}
		if e = ProductConfig(r.StockData, u, o); e != nil {
			o.Rollback()
			return nil, e
		}
	} else {
		if e = WarehouseStock(r.Product, u, o); e != nil {
			o.Rollback()
			return nil, e
		}
	}

	o.Commit()
	return u, e
}

func Archive(r archiveRequest) (u *model.Warehouse, e error) {
	u = &model.Warehouse{
		ID:     r.ID,
		Status: int8(2),
	}

	if e = u.Save("id", "status"); e == nil {
		e = log.AuditLogByUser(r.Session.Staff, u.ID, "warehouse", "archive", "")
	}

	return u, e
}

func Unarchive(r unarchiveRequest) (u *model.Warehouse, e error) {
	u = &model.Warehouse{
		ID:     r.ID,
		Status: int8(1),
	}

	if e = u.Save("id", "status"); e == nil {
		e = log.AuditLogByUser(r.Session.Staff, u.ID, "warehouse", "unarchive", "")
	}

	return u, e
}

func Update(r updateRequest) (u *model.Warehouse, e error) {
	o := orm.NewOrm()
	o.Begin()

	var binInfo *model.BinInfo
	if r.BinInfoExist == true {
		binInfo = &model.BinInfo{
			ID:        r.Warehouse.BinInfo.ID,
			Latitude:  &r.PickerStartingLatitude,
			Longitude: &r.PickerStartingLongitude,
			ImageURL:  r.FloorPlanLink,
		}

		_, e = o.Update(binInfo, "Latitude", "Longitude", "ImageURL")
		if e != nil {
			o.Rollback()
			return nil, e
		}

		e = log.AuditLogByUser(r.Session.Staff, binInfo.ID, "bin info", "update", "")
		if e != nil {
			o.Rollback()
			return nil, e
		}
	} else {
		binInfo = &model.BinInfo{
			Latitude:  &r.PickerStartingLatitude,
			Longitude: &r.PickerStartingLongitude,
			ImageURL:  r.FloorPlanLink,
		}

		_, e = o.Insert(binInfo)
		if e != nil {
			o.Rollback()
			return nil, e
		}

		e = log.AuditLogByUser(r.Session.Staff, binInfo.ID, "bin info", "create", "")
		if e != nil {
			o.Rollback()
			return nil, e
		}
	}

	u = &model.Warehouse{
		ID:             r.ID,
		PicName:        r.PicName,
		PhoneNumber:    r.PhoneNumber,
		AltPhoneNumber: r.AltPhoneNumber,
		StreetAddress:  r.StreetAddress,
		Latitude:       r.Latitude,
		Longitude:      r.Longitude,
		Note:           r.Note,
		WarehouseType:  r.Glossary.ValueInt,
		BinInfo:        binInfo,
	}

	if r.Glossary.ValueName == "HUB" {
		u.HubProcessingTime = r.HubProcessingTime
		u.ParentID = r.ParentWarehouseStruct
		_, e = o.Update(u, "PicName", "PhoneNumber", "AltPhoneNumber", "StreetAddress", "Latitude", "Longitude", "Note", "WarehouseType", "BinInfo", "HubProcessingTime", "ParentID")
	} else {
		_, e = o.Update(u, "PicName", "PhoneNumber", "AltPhoneNumber", "StreetAddress", "Latitude", "Longitude", "Note", "WarehouseType", "BinInfo")
	}

	e = log.AuditLogByUser(r.Session.Staff, u.ID, "warehouse", "update", "")
	if e != nil {
		o.Rollback()
		return nil, e
	}

	o.Commit()
	return
}

func UpdateParam(r updateParamRequest) (u *model.Warehouse, e error) {
	o := orm.NewOrm()
	o.Begin()
	u = &model.Warehouse{
		ID:                     r.ID,
		LimitOrderPickingList:  r.LimitSalesOrder,
		LimitWeightPickingList: r.LimitWeight,
	}

	if _, e = o.Update(u, "LimitOrderPickingList", "LimitWeightPickingList"); e != nil {
		o.Rollback()
		return nil, e
	}

	e = log.AuditLogByUser(r.Session.Staff, u.ID, "warehouse", "update", "update param")
	o.Commit()
	return u, e
}

// WarehouseCopyStock : function to imitiate warehouse's stock data into database
func WarehouseCopyStock(stock []*model.Stock, warehouse *model.Warehouse, o orm.Ormer) (e error) {
	var newStockArr []*model.Stock
	for _, v := range stock {
		newStock := &model.Stock{
			Salable:     v.Salable,
			Purchasable: v.Purchasable,
			Status:      v.Status,
			Product:     v.Product,
			Warehouse:   warehouse,
			Bin:         nil,
		}

		newStockArr = append(newStockArr, newStock)
	}

	if _, e := o.InsertMulti(100, &newStockArr); e != nil {
		return e
	}

	return e
}

// ProductConfig : function to save warehouse's product data of salability,storability and purchasability into database
func ProductConfig(stock []*model.Stock, warehouse *model.Warehouse, o orm.Ormer) (e error) {
	for _, v := range stock {
		if v.Status == 1 {
			product, e := repository.GetProductNoDecrypt("ID", v.Product.ID)
			if e != nil {
				return e
			}

			warehouseStr := strconv.Itoa(int(warehouse.ID))

			// storability
			product.WarehouseStoStr += "," + warehouseStr

			// purchasability
			if v.Purchasable == 1 {
				product.WarehousePurStr += "," + warehouseStr
			}

			// salablity
			if v.Salable == 1 {
				product.WarehouseSalStr += "," + warehouseStr
			}

			if _, e = o.Update(product, "WarehouseStoStr", "WarehousePurStr", "WarehouseSalStr"); e != nil {
				return e
			}
		}
	}

	return e
}

// WarehouseStock : function to create warehouse's stock data into database
func WarehouseStock(product []*model.Product, warehouse *model.Warehouse, o orm.Ormer) (e error) {
	var stockArr []*model.Stock

	for _, v := range product {
		stock := &model.Stock{
			Salable:     2,
			Purchasable: 2,
			Status:      2,
			Product:     &model.Product{ID: v.ID},
			Warehouse:   warehouse,
			Bin:         nil,
		}

		stockArr = append(stockArr, stock)
	}

	if _, e := o.InsertMulti(100, &stockArr); e != nil {
		return e
	}

	return e
}
