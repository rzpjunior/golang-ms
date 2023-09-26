// Copyright 2020 PT. Eden Pangan Indonesia. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

package product

import (
	"strings"

	"git.edenfarm.id/cuxs/orm"

	"git.edenfarm.id/project-version2/api/log"
	"git.edenfarm.id/project-version2/datamodel/model"
)

// Save : function to save data requested into database
func Save(r createRequest) (product *model.Product, e error) {
	o := orm.NewOrm()
	o.Begin()

	tagProduct := strings.Join(r.TagProduct, ",")

	product = &model.Product{
		Code:            r.Code,
		Name:            r.Name,
		Note:            r.Note,
		Description:     r.Description,
		Status:          int8(1),
		UnitWeight:      r.Weight,
		WarehouseStoStr: r.WarehouseStoStr,
		WarehousePurStr: r.WarehousePurStr,
		TagProduct:      tagProduct,
		Salability:      int8(2),
		Storability:     r.Storability,
		Purchasability:  r.Purchasability,
		Uom:             r.Uom,
		Category:        r.Category,
		UnivProductCode: r.UnivProductCode,
		OrderMinQty:     r.OrderMinQty,
		Packability:     int8(2),
		SparePercentage: r.SparePercentage,
		Taxable:         r.Taxable,
		TaxPercentage:   r.TaxPercentage,
		FragileGoods:    r.FragileGoods,
	}

	if product.ID, e = o.Insert(product); e == nil {
		var mainImage int8 = 1

		for _, v := range r.Images {
			productImage := &model.ProductImage{
				Product:   product,
				ImageUrl:  v,
				MainImage: mainImage,
			}

			if _, productImage.ID, e = o.ReadOrCreate(productImage, "Product", "ImageUrl", "MainImage"); e != nil {
				o.Rollback()
				return nil, e
			}

			mainImage = 2
		}

		for _, v := range r.Warehouse {
			status := 2
			purchasable := 2
			if _, isExist := r.WrhPurExist[v.ID]; isExist {
				purchasable = 1
			}
			if _, isExist := r.WrhStoExist[v.ID]; isExist {
				status = 1
			}
			stock := &model.Stock{
				Product:     product,
				Warehouse:   v,
				Salable:     int8(2),
				Purchasable: int8(purchasable),
				Status:      int8(status),
			}

			if _, e = o.Insert(stock); e != nil {
				o.Rollback()
				return nil, e
			}
		}

		for _, v := range r.PriceSet {
			price := &model.Price{
				Product:     product,
				PriceSet:    v,
				UnitPrice:   float64(0),
				ShadowPrice: float64(0),
			}

			if _, e = o.Insert(price); e != nil {
				o.Rollback()
				return nil, e
			}
		}
	} else {
		o.Rollback()
		return nil, e
	}

	if e = log.AuditLogByUser(r.Session.Staff, product.ID, "product", "create", ""); e != nil {
		o.Rollback()
		return nil, e
	}

	o.Commit()
	return product, e
}

func Archive(r archiveRequest) (product *model.Product, e error) {
	product = &model.Product{
		ID:     r.ID,
		Status: int8(2),
	}

	if e = product.Save("id", "status"); e == nil {
		e = log.AuditLogByUser(r.Session.Staff, product.ID, "product", "archive", "")
	}

	return product, e
}

func Unarchive(r unarchiveRequest) (product *model.Product, e error) {
	product = &model.Product{
		ID:     r.ID,
		Status: int8(1),
	}

	if e = product.Save("id", "status"); e == nil {
		e = log.AuditLogByUser(r.Session.Staff, product.ID, "product", "unarchive", "")
	}

	return product, e
}

// Update : function update data in database
func Update(r updateRequest) (product *model.Product, e error) {
	o := orm.NewOrm()
	o.Begin()

	tagProduct := strings.Join(r.TagProduct, ",")
	orderChannelRestriction := strings.Join(r.OrderChannelRestriction, ",")

	product = &model.Product{
		ID:                      r.ID,
		Code:                    r.Code,
		Name:                    r.Name,
		Note:                    r.Note,
		TagProduct:              tagProduct,
		OrderChannelRestriction: orderChannelRestriction,
		Description:             r.Description,
		Category:                r.Category,
		UnivProductCode:         r.UnivProductCode,
		OrderMinQty:             r.OrderMinQty,
		SparePercentage:         r.SparePercentage,
		Taxable:                 r.Taxable,
		TaxPercentage:           r.TaxPercentage,
		OrderMaxQty:             r.OrderMaxQty,
		ExcludeArchetype:        r.ExcludeArchetypeStr,
		MaxDayDeliveryDate:      r.MaxDayDeliveryDate,
	}

	if _, e = o.Update(product, "ID", "Code", "Name", "Note", "TagProduct", "OrderChannelRestriction", "Description", "Category", "UnivProductCode", "OrderMinQty", "SparePercentage", "Taxable", "TaxPercentage", "OrderMaxQty", "ExcludeArchetype", "MaxDayDeliveryDate"); e == nil {
		mainImage := int8(1)
		var keepImagesId []int64

		for _, v := range r.ImagesUrl {
			productImage := &model.ProductImage{
				Product:   product,
				ImageUrl:  v,
				MainImage: mainImage,
			}

			if _, productImage.ID, e = o.ReadOrCreate(productImage, "Product", "ImageUrl", "MainImage"); e != nil {
				o.Rollback()
				return nil, e
			}

			mainImage = 2

			keepImagesId = append(keepImagesId, productImage.ID)
		}

		if _, e := o.QueryTable(new(model.ProductImage)).Filter("product_id", product.ID).Exclude("ID__in", keepImagesId).Delete(); e != nil {
			o.Rollback()
			return nil, e
		}

		for _, v := range r.StorabilityStock {
			stock := &model.Stock{
				Product:   product,
				Warehouse: v.Warehouse,
			}

			if e = stock.Read("Product", "Warehouse"); e == nil {
				stock.SafetyStock = v.WarehouseStock[v.WarehouseId]
				if _, e = o.Update(stock, "ID", "SafetyStock"); e != nil {
					o.Rollback()
					return
				}
			} else {
				o.Rollback()
				return
			}
		}

		e = log.AuditLogByUser(r.Session.Staff, product.ID, "product", "update", "")
	}

	o.Commit()
	return
}

// Salability : function to update salability in database
func Salability(r salabilityRequest) (product *model.Product, e error) {
	o := orm.NewOrm()
	o.Begin()

	product = &model.Product{
		ID:              r.ID,
		Salability:      r.Salability,
		WarehouseSalStr: r.WarehouseStr,
	}

	if _, e = o.Update(product, "Salability", "WarehouseSalStr"); e == nil {
		q := o.QueryTable("stock").Filter("Product", r.ID)

		if r.Salability == 2 {
			if _, e = q.Update(orm.Params{"salable": 2}); e != nil {
				o.Rollback()
				return nil, e
			}
		} else {
			if _, e = q.Filter("Warehouse__in", r.Warehouse).Update(orm.Params{"salable": 1}); e != nil {
				o.Rollback()
				return nil, e
			}

			if _, e = q.Exclude("Warehouse__in", r.Warehouse).Update(orm.Params{"salable": 2}); e != nil {
				o.Rollback()
				return nil, e
			}
		}

		e = log.AuditLogByUser(r.Session.Staff, product.ID, "product", "update salability", "")
	}

	o.Commit()
	return
}

// Purchasability : function to update purchasability in database
func Purchasability(r purchasabilityRequest) (product *model.Product, e error) {
	o := orm.NewOrm()
	o.Begin()

	product = &model.Product{
		ID:              r.ID,
		Purchasability:  r.Purchasability,
		WarehousePurStr: r.WarehouseStr,
	}

	if _, e = o.Update(product, "Purchasability", "WarehousePurStr"); e == nil {
		q := o.QueryTable("stock").Filter("Product", r.ID)

		if r.Purchasability == 2 {
			if _, e = q.Update(orm.Params{"purchasable": 2}); e != nil {
				o.Rollback()
				return nil, e
			}
		} else {
			if _, e = q.Filter("Warehouse__in", r.WarehouseChecked).Update(orm.Params{"purchasable": 1}); e != nil {
				o.Rollback()
				return nil, e
			}

			if _, e = q.Exclude("Warehouse__in", r.WarehouseChecked).Update(orm.Params{"purchasable": 2}); e != nil {
				o.Rollback()
				return nil, e
			}
		}

		e = log.AuditLogByUser(r.Session.Staff, product.ID, "product", "update purchasability", "")
	} else {
		o.Rollback()
		return nil, e
	}

	o.Commit()
	return
}

// Storability : function to update storability in database
func Storability(r storabilityRequest) (product *model.Product, e error) {
	o := orm.NewOrm()
	o.Begin()

	product = &model.Product{
		ID:              r.ID,
		WarehouseStoStr: r.WarehouseStr,
		Storability:     r.Storability,
		Status:          r.Status,
	}

	if _, e = o.Update(product, "WarehouseStoStr", "Storability", "Status"); e == nil {
		if r.Storability == 2 {
			q := o.QueryTable("stock").Filter("Product", r.ID)

			if _, e = q.Update(orm.Params{"status": 2}); e != nil {
				o.Rollback()
				return nil, e
			}
		} else {
			q := o.QueryTable("stock").Filter("Product", r.ID)

			if _, e = q.Filter("Warehouse__in", r.WarehouseChecked).Update(orm.Params{"status": 1}); e != nil {
				o.Rollback()
				return nil, e
			}

			if _, e = q.Exclude("Warehouse__in", r.WarehouseChecked).Update(orm.Params{"status": 2}); e != nil {
				o.Rollback()
				return nil, e
			}
		}

		e = log.AuditLogByUser(r.Session.Staff, product.ID, "product", "update storability", "")
	} else {
		o.Rollback()
		return nil, e
	}

	o.Commit()
	return
}

//Packable : function to update packablity in database
func Packable(r packableRequest) (product *model.Product, e error) {
	o := orm.NewOrm()
	o.Begin()

	product = &model.Product{
		ID:          r.ID,
		Packability: 1,
	}

	if _, e = o.Update(product, "Packability"); e != nil {
		o.Rollback()
		return nil, e
	}

	e = log.AuditLogByUser(r.Session.Staff, product.ID, "packable_product", "add", "")

	o.Commit()
	return
}

//Unpackable : function to update packablity in database
func Unpackable(r unpackableRequest) (product *model.Product, e error) {
	o := orm.NewOrm()
	o.Begin()

	product = &model.Product{
		ID:          r.ID,
		Packability: 2,
	}

	if _, e = o.Update(product, "Packability"); e != nil {
		o.Rollback()
		return nil, e
	}

	e = log.AuditLogByUser(r.Session.Staff, product.ID, "packable_product", "delete", "")

	o.Commit()
	return
}

//Fragile : function to update fragilegoods in database
func Fragile(r fragileRequest) (product *model.Product, e error) {
	o := orm.NewOrm()
	o.Begin()

	product = &model.Product{
		ID:           r.ID,
		FragileGoods: 1,
	}

	if _, e = o.Update(product, "fragile_goods"); e != nil {
		o.Rollback()
		return nil, e
	}

	e = log.AuditLogByUser(r.Session.Staff, product.ID, "fragile_goods", "add", "")
	if e != nil {
		o.Rollback()
		return nil, e
	}

	o.Commit()
	return
}

//NotFragile : function to update fragilegoods in database
func NotFragile(r notFragileRequest) (product *model.Product, e error) {
	o := orm.NewOrm()
	o.Begin()

	product = &model.Product{
		ID:           r.ID,
		FragileGoods: 2,
	}

	if _, e = o.Update(product, "fragile_goods"); e != nil {
		o.Rollback()
		return nil, e
	}

	e = log.AuditLogByUser(r.Session.Staff, product.ID, "fragile_goods", "delete", "")
	if e != nil {
		o.Rollback()
		return nil, e
	}

	o.Commit()
	return
}
