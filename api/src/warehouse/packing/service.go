// Copyright 2020 PT. Eden Pangan Indonesia. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

package packing

import (
	"encoding/json"
	"fmt"
	"math"
	"os"
	"sort"
	"strconv"
	"time"

	"git.edenfarm.id/cuxs/mongodb"
	"git.edenfarm.id/project-version2/api/src/auth"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo/options"

	"git.edenfarm.id/cuxs/common"
	"git.edenfarm.id/cuxs/env"
	"git.edenfarm.id/cuxs/orm"
	"git.edenfarm.id/project-version2/api/log"
	"git.edenfarm.id/project-version2/api/util"
	"git.edenfarm.id/project-version2/datamodel/model"
	"github.com/tealeg/xlsx"
)

// Save : function to save data requested into database
func Save(r createRequest) (u *model.PackingOrder, e error) {

	//generate codes for document
	codePacking, _ := util.GenerateDocCode("PC", r.Warehouse.Code, "packing_order")
	o := orm.NewOrm()
	o.Begin()

	u = &model.PackingOrder{
		Code:         codePacking,
		Warehouse:    r.Warehouse,
		DeliveryDate: r.DeliveryDateTime,
		Status:       1,
		Note:         r.Note,
	}

	if _, e = o.Insert(u); e != nil {
		o.Rollback()
		return nil, e
	}

	for _, row := range r.PackingOrderItems {

		item := &model.PackingOrderItem{
			PackingOrder: &model.PackingOrder{ID: u.ID},
			Product:      row.Product,
			TotalOrder:   row.TotalOrder,
		}

		if _, e = o.Insert(item); e != nil {
			o.Rollback()
			return nil, e
		}

		for _, rows := range row.Helper {

			idStaff, e := common.Decrypt(rows)

			items := &model.PackingOrderItemAssign{
				PackingOrderItem: &model.PackingOrderItem{ID: item.ID},
				Staff:            &model.Staff{ID: idStaff},
			}
			if _, e = o.Insert(items); e != nil {
				o.Rollback()
				return nil, e
			}
		}
	}

	e = log.AuditLogByUser(r.Session.Staff, u.ID, "packing_order", "create", "")
	o.Commit()

	return u, e
}

// Update : function to save data requested into database
func Update(r updateRequest) (u *model.PackingOrder, e error) {
	o := orm.NewOrm()
	o.Begin()

	for _, row := range r.PackingOrderItems {

		item := &model.PackingOrderItem{
			PackingOrder: &model.PackingOrder{ID: u.ID},
			Product:      row.Product,
			TotalOrder:   row.TotalOrder,
		}

		if _, e = o.Insert(item); e == nil {

			for _, rows := range row.Helper {

				idStaff, e := common.Decrypt(rows)

				items := &model.PackingOrderItemAssign{
					PackingOrderItem: &model.PackingOrderItem{ID: item.ID},
					Staff:            &model.Staff{ID: idStaff},
				}

				if _, e = o.Insert(items); e != nil {
					o.Rollback()
					return nil, e
				}

			}

		} else {
			o.Rollback()
			return nil, e
		}
	}

	e = log.AuditLogByUser(r.Session.Staff, u.ID, "packing_order", "create", "")

	o.Commit()

	return u, e
}

// UpdateOrderItemAssign : function to save data requested into database
func UpdateItemAssign(r updateItemAssignRequest) (u *model.PackingOrder, e error) {

	o := orm.NewOrm()
	o.Begin()

	var isItemCreated bool
	var deleted []int64

	if len(r.HelperDec) > 0 {
		for _, rows := range r.Helper {

			idStaff, e := common.Decrypt(rows)

			items := &model.PackingOrderItemAssign{
				PackingOrderItem: r.Poi,
				Staff:            &model.Staff{ID: idStaff},
			}

			if isItemCreated, items.ID, e = o.ReadOrCreate(items, "PackingOrderItem", "Staff"); e == nil {
				if !isItemCreated {
					items.PackingOrderItem = r.Poi
					items.Staff = &model.Staff{ID: idStaff}
					if _, e = o.Update(items, "PackingOrderItem", "Staff"); e != nil {
						o.Rollback()
						return nil, e
					}
				}
				deleted = append(deleted, idStaff)

			} else {
				o.Rollback()
				return nil, e
			}

		}

		if _, err := o.QueryTable(new(model.PackingOrderItemAssign)).Filter("packing_order_item_id", r.Poi.ID).Exclude("staff_id__in", deleted).Delete(); err != nil {
			o.Rollback()
			return nil, e
		}
	} else {
		if _, err := o.QueryTable(new(model.PackingOrderItemAssign)).Filter("packing_order_item_id", r.Poi.ID).Delete(); err != nil {
			o.Rollback()
			return nil, e
		}
	}

	o.Commit()

	return u, e
}

func DownloadActualPackingXls(date time.Time, p *model.PackingOrder) (filePath string, err error) {
	var file *xlsx.File
	var sheet *xlsx.Sheet
	var row *xlsx.Row

	dir := env.GetString("EXPORT_DIRECTORY", "")

	filename := fmt.Sprintf("TemplatePackingOrder_%s_%s_%s.xlsx", p.Code, date.Format("2006-01-02"), util.GenerateRandomDoc(5))
	fileDir := fmt.Sprintf("%s/%s", dir, filename)

	file = xlsx.NewFile()
	if sheet, err = file.AddSheet("Sheet1"); err == nil {

		row = sheet.AddRow()
		row.SetHeight(20)
		row.AddCell().Value = "No"
		row.AddCell().Value = "Packing_Order_Item_ID"
		row.AddCell().Value = "Product_ID"
		row.AddCell().Value = "Product_Code"
		row.AddCell().Value = "Product_Name"
		row.AddCell().Value = "UOM"
		row.AddCell().Value = "Packer_ID"
		row.AddCell().Value = "Packer_Code"
		row.AddCell().Value = "Packer_Name"
		row.AddCell().Value = "Total_Order"
		row.AddCell().Value = "Total_Weight"
		row.AddCell().Value = "Total_Packing"

		var idx int
		for _, v := range p.PackingOrderItems {
			if v.Helper != nil {
				for _, h := range v.Helper {
					idx += 1
					row = sheet.AddRow()
					row.AddCell().SetInt(idx)
					row.AddCell().Value = common.Encrypt(v.ID)            // Packing Order Item ID
					row.AddCell().Value = common.Encrypt(v.Product.ID)    // Product ID
					row.AddCell().Value = v.Product.Code                  // Product Code
					row.AddCell().Value = v.Product.Name                  // Product Name
					row.AddCell().Value = v.Product.Uom.Name              // UOM
					row.AddCell().Value = common.Encrypt(h.ID)            // Packer ID
					row.AddCell().Value = h.Code                          // Packer Code
					row.AddCell().Value = h.Name                          // Packer Name
					row.AddCell().Value = strconv.Itoa(int(v.TotalOrder)) // Total Order
					row.AddCell().Value = "0"                             // Total Weight
					row.AddCell().Value = "0"                             // Total Packing
				}
			} else {
				idx += 1
				row = sheet.AddRow()
				row.AddCell().SetInt(idx)
				row.AddCell().Value = common.Encrypt(v.ID)            // Product ID
				row.AddCell().Value = common.Encrypt(v.Product.ID)    // Packing Order Item ID
				row.AddCell().Value = v.Product.Code                  // Product Code
				row.AddCell().Value = v.Product.Name                  // Product Name
				row.AddCell().Value = v.Product.Uom.Name              // UOM
				row.AddCell().Value = ""                              // Packer ID
				row.AddCell().Value = ""                              // Packer Code
				row.AddCell().Value = ""                              // Packer Name
				row.AddCell().Value = strconv.Itoa(int(v.TotalOrder)) // Total Order
				row.AddCell().Value = "0"                             // Total Weight
				row.AddCell().Value = "0"                             // Total Packing
			}

		}
		boldStyle := xlsx.NewStyle()
		boldFont := xlsx.NewFont(10, "Liberation Sans")
		boldFont.Bold = true
		boldStyle.Font = *boldFont
		boldStyle.ApplyFont = true

		//looping to get column range 0-7. making BOLD font header
		for col := 0; col < 28; col++ {
			sheet.Cell(0, col).SetStyle(boldStyle)
		}

		err = file.Save(fileDir)

		filePath, err = util.UploadToS3(filename, fileDir, "application/xlsx")

		// fungsi ini berguna untuk menghapus kembali file yang sudah diupload ke server
		os.Remove(fileDir)

	}
	return
}

// GetSalesOrderItem get all data sales_order_item that matched with query request parameters.
// returning slices of SalesOrderItem, total data without limit and error.
func GetSalesOrderItem(rq *orm.RequestQuery, date string, wh int64) (m []*model.SalesOrderItem, total int64, err error) {
	// make new orm query
	q, _ := rq.QueryReadOnly(new(model.SalesOrderItem))

	// get total data
	if total, err = q.Filter("sales_order_id__warehouse__id", wh).Filter("sales_order_id__delivery_date", date).Filter("product_id__packability", 1).Filter("sales_order_id__status__in", 1, 9, 12).Count(); err != nil || total == 0 {
		return nil, total, err
	}

	// get data requested
	var mx []*model.SalesOrderItem
	if _, err = q.Filter("sales_order_id__warehouse__id", wh).Filter("sales_order_id__delivery_date", date).Filter("product_id__packability", 1).Filter("sales_order_id__status__in", 1, 9, 12).All(&mx, rq.Fields...); err == nil {
		poi := removeDuplicateSoi(mx)
		for _, i := range poi {
			i.Product.Category.Read("ID")
			i.Product.Uom.Read("ID")
		}
		return poi, total, nil
	}

	// return error some thing went wrong
	return nil, total, err
}

func removeDuplicateSoi(poi []*model.SalesOrderItem) []*model.SalesOrderItem {
	// Use map to record duplicates as we find them.
	duplicate := map[int64]bool{}
	var result []*model.SalesOrderItem
	var storageTemp []*model.SalesOrderItem
	for _, v := range poi {
		v.Product.Read("ID")
		if duplicate[v.Product.ID] == true {
			// data duplicate will append to storageTemp
			storageTemp = append(storageTemp, v)

		} else {
			// Record this element as an encountered element.
			duplicate[v.Product.ID] = true
			result = append(result, v)
			// Append to result slice.
		}
	}

	for _, i := range result {
		for _, a := range storageTemp {
			if i.Product.ID == a.Product.ID {
				i.OrderQty += a.OrderQty
			}
		}
	}

	// Return the new slice.
	return result
}

// Update
func UpdateActualPacking(u updateActualPackingRequest) (p *model.PackingOrder, e error) {
	o := orm.NewOrm()
	o.Begin()

	for _, v := range u.PackingOrderItemUpdate {
		if v.Helper != nil {
			if _, e = o.Raw("update packing_order_item_assign poia "+
				"set poia.subtotal_pack = poia.subtotal_pack + ?, poia.subtotal_weight =  poia.subtotal_weight + ? "+
				"where poia.packing_order_item_id = ? and poia.staff_id = ?", v.TotalPack, v.TotalWeight, v.PackingOrderItem.ID, v.Helper.ID).Exec(); e == nil {
			} else {
				o.Rollback()
			}
		}
	}

	if _, e = o.Raw("update packing_order_item poi " +
		"join ( " +
		"select packing_order_item_id, sum(subtotal_weight) sw, sum(subtotal_pack) sp " +
		"from packing_order_item_assign " +
		"group by packing_order_item_id " +
		") x on poi.id = x.packing_order_item_id " +
		"set poi.total_weight = x.sw, poi.total_pack = x.sp").Exec(); e == nil {
	} else {
		o.Rollback()
	}

	o.Commit()
	return u.PackingOrder, nil
}

// Confirm : function to save data requested into database
func Confirm(r confirmRequest) (po *model.PackingOrder, e error) {
	o := orm.NewOrm()
	o.Begin()

	r.PackingOrder.Status = 2
	if _, e = o.Update(r.PackingOrder, "Status"); e != nil {
		o.Rollback()
		return nil, e
	}

	if e = log.AuditLogByUser(r.Session.Staff, r.PackingOrder.ID, "packing_order", "confirm", ""); e != nil {
		o.Rollback()
		return nil, e
	}

	o.Commit()
	return
}

// Cancel : function to save data requested into database
func Cancel(r cancelRequest) (po *model.PackingOrder, e error) {
	o := orm.NewOrm()
	o.Begin()

	r.PackingOrder.Status = 3
	if _, e = o.Update(r.PackingOrder, "Status"); e != nil {
		o.Rollback()
		return nil, e
	}

	if e = log.AuditLogByUser(r.Session.Staff, r.PackingOrder.ID, "packing_order", "cancel", ""); e != nil {
		o.Rollback()
		return nil, e
	}

	o.Commit()
	return
}

// Assign Quantity for packing mobile
func AssignActualPacking(u assignActualPackingRequest) (p *model.PackingOrder, e error) {
	o := orm.NewOrm()
	o.Begin()

	log := &model.PackingHelperLog{
		PackingOrderItem: &model.PackingOrderItem{ID: u.PackingOrderItem.ID},
		Helper:           &model.Staff{ID: u.Helper.ID},
		QtyWeight:        u.TotalWeight,
		QtyPack:          u.TotalPack,
		CreatedAt:        time.Now(),
	}
	if _, e = o.Insert(log); e != nil {
		o.Rollback()
	}

	if u.Helper != nil {
		if _, e = o.Raw("update packing_order_item_assign poia "+
			"set poia.subtotal_pack = poia.subtotal_pack + ?, poia.subtotal_weight =  poia.subtotal_weight + ? "+
			"where poia.packing_order_item_id = ? and poia.staff_id = ?", u.TotalPack, u.TotalWeight, u.PackingOrderItem.ID, u.Helper.ID).Exec(); e == nil {
		} else {
			o.Rollback()
		}
	} else {
		o.Rollback()
	}

	if _, e = o.Raw("update packing_order_item poi " +
		"join ( " +
		"select packing_order_item_id, sum(subtotal_weight) sw, sum(subtotal_pack) sp " +
		"from packing_order_item_assign " +
		"group by packing_order_item_id " +
		") x on poi.id = x.packing_order_item_id " +
		"set poi.total_weight = x.sw, poi.total_pack = x.sp").Exec(); e == nil {
	} else {
		o.Rollback()
	}

	o.Commit()
	return u.PackingOrder, nil
}

// GeneratePacking: function to generate packing
func GeneratePacking(r generatePackingRequest) (pc *model.PackingOrder, e error) {
	m := mongodb.NewMongo()
	o := orm.NewOrm()
	var packingRecommendation []interface{}
	var packingSalesOrder []interface{}

	if r.PackingOrder != nil {

		filter := bson.D{
			{"packing_order_id", r.PackingOrder.ID},
		}

		var res []byte
		if res, e = m.GetAllDataWithFilter("Packing_Item", filter); e != nil {
			m.DisconnectMongoClient()
			return nil, e
		}

		// region convert byte data to json data
		if e = json.Unmarshal(res, &r.ResponseData); e != nil {
			m.DisconnectMongoClient()
			return nil, e
		}
		// endregion

		/*
			mapExistData -> for storing from data get from sql data
			mapOldExistData -> for storing from data get from mongo data
		*/
		for _, v := range r.ResponseData {
			r.MapOldExistData[v.ProductID][v.PackType] = v.ActualTotalPack
		}
		// region comparing between old data and the new one
		if e = GeneratePackingRecursive(r); e != nil {
			m.DisconnectMongoClient()
			o.Rollback()
			return nil, e
		}

		for _, v := range r.ResponseData {
			if _, ok := r.MapResProductPackType[v.ProductID]; ok {
				r.MapExistData[v.ProductID][v.PackType] = v.ActualTotalPack
			} else {
				/*
					product is no longer exist and need to update status to 3
					update status = 3 where id = ?
					after changed running code:
					delete(MapResProductPackType,v.ProductID)
				*/
				filterUpdateMany := bson.D{
					{"packing_order_id", r.PackingOrder.ID},
					{"product_id", v.ProductID},
				}
				updatePayload := bson.D{{"status", 3}}
				if e = m.UpdateManyDataWithFilter("Packing_Item", filterUpdateMany, updatePayload); e != nil {
					m.DisconnectMongoClient()
					return nil, e
				}

				if e = m.UpdateManyDataWithFilter("Packing_Sales_Order", filterUpdateMany, updatePayload); e != nil {
					m.DisconnectMongoClient()
					return nil, e
				}

				delete(r.MapResProductPackType, v.ProductID)

			}
		}
		//endregion

		for k, v := range r.MapResProductPackType {
			if _, ok := r.MapExistData[k]; ok {
				/*
					update actual pack size
				*/
				for k2, v2 := range v {
					var rd *model.ResponseData
					filter := bson.D{
						{"packing_order_id", r.PackingOrder.ID},
						{"product_id", k},
						{"pack_type", k2},
					}

					var res []byte
					if res, e = m.GetOneDataWithFilter("Packing_Item", filter); e != nil {
						m.DisconnectMongoClient()
						return nil, e
					}

					// region: warning for new Item
					if len(res) == 0 {
						if k2 == -1 {
							packingRecommendationTemp := bson.D{
								{"packing_order_id", r.PackingOrder.ID},
								{"product_id", k},
								{"pack_type", k2},
								{"expected_total_pack", v2},
								{"actual_total_pack", 0},
								{"weight_pack", v2},
								{"status", 1},
							}

							_, e = m.InsertOneData("Packing_Item", packingRecommendationTemp)
							if e != nil {
								m.DisconnectMongoClient()
								return nil, e
							}

							for i, v := range r.MapPackingSalesOrder[k][k2] {

								packingSalesOrderTemp := bson.D{
									{"packing_order_id", r.PackingOrder.ID},
									{"product_id", k},
									{"pack_type", k2},
									{"expected_total_pack", v},
									{"sales_order_id", i},
									{"status", 1},
								}

								_, e = m.InsertOneData("Packing_Sales_Order", packingSalesOrderTemp)
								if e != nil {
									m.DisconnectMongoClient()
									return nil, e
								}

							}

							continue

						} else {
							packingRecommendationTemp := bson.D{
								{"packing_order_id", r.PackingOrder.ID},
								{"product_id", k},
								{"pack_type", k2},
								{"expected_total_pack", r.MapResProductPackType[k][k2]},
								{"actual_total_pack", 0},
								{"weight_pack", r.MapProductPack[k][k2]},
								{"status", 1},
							}

							_, e = m.InsertOneData("Packing_Item", packingRecommendationTemp)
							if e != nil {
								m.DisconnectMongoClient()
								return nil, e
							}

							for i, v := range r.MapPackingSalesOrder[k][k2] {
								packingSalesOrderTemp := bson.D{
									{"packing_order_id", r.PackingOrder.ID},
									{"product_id", k},
									{"pack_type", k2},
									{"expected_total_pack", v},
									{"sales_order_id", i},
									{"status", 1},
								}

								_, e = m.InsertOneData("Packing_Sales_Order", packingSalesOrderTemp)
								if e != nil {
									m.DisconnectMongoClient()
									return nil, e
								}
							}

							continue
						}
						// endregion
					} else {
						// if packing order item already exist, add sales order packing again
						for i, v := range r.MapPackingSalesOrder[k][k2] {

							filter := bson.D{
								{"packing_order_id", r.PackingOrder.ID},
								{"product_id", k},
								{"pack_type", k2},
								{"sales_order_id", i},
							}

							var resSo []byte
							if resSo, e = m.GetOneDataWithFilter("Packing_Sales_Order", filter); e != nil {
								m.DisconnectMongoClient()
								return nil, e
							}

							if len(resSo) == 0 {
								packingSalesOrderTemp := bson.D{
									{"packing_order_id", r.PackingOrder.ID},
									{"product_id", k},
									{"pack_type", k2},
									{"expected_total_pack", v},
									{"sales_order_id", i},
									{"status", 1},
								}

								_, e = m.InsertOneData("Packing_Sales_Order", packingSalesOrderTemp)
								if e != nil {
									m.DisconnectMongoClient()
									return nil, e
								}
							}
						}

						if e = json.Unmarshal(res, &rd); e != nil {
							m.DisconnectMongoClient()
							return nil, e
						}
						// endregion

						rd.ActualTotalPack = r.MapOldExistData[k][k2]
						rd.ExpectedTotalPack = v2
						rd.Status = 1
						if k2 == -1 {
							rd.WeightPack = v2
						}

						if e = m.UpdateOneDataWithFilter("Packing_Item", filter, rd); e != nil {
							m.DisconnectMongoClient()
							return nil, e
						}

					}

				}
			} else {
				/*
					insert new data
				*/
				for k2, v2 := range v {
					if k2 == -1 {
						packingRecommendationTemp := bson.D{
							{"packing_order_id", pc.ID},
							{"product_id", k},
							{"pack_type", k2},
							{"expected_total_pack", v2},
							{"actual_total_pack", 0},
							{"weight_pack", r.MapResProductPackType[k][-1]},
							{"status", 1},
						}

						packingRecommendation = append(packingRecommendation, packingRecommendationTemp)

						for i, v := range r.MapPackingSalesOrder[k][k2] {
							packingSalesOrderTemp := bson.D{
								{"packing_order_id", pc.ID},
								{"product_id", k},
								{"pack_type", k2},
								{"expected_total_pack", v},
								{"sales_order_id", i},
								{"status", 1},
							}

							packingSalesOrder = append(packingSalesOrder, packingSalesOrderTemp)
						}

						continue
					} else {

						packingRecommendationTemp := bson.D{
							{"packing_order_id", r.PackingOrder.ID},
							{"product_id", k},
							{"pack_type", k2},
							{"expected_total_pack", v2},
							{"actual_total_pack", 0},
							{"weight_pack", r.MapProductPack[k][k2]},
							{"status", 1},
						}

						packingRecommendation = append(packingRecommendation, packingRecommendationTemp)

						for i, v := range r.MapPackingSalesOrder[k][k2] {

							packingSalesOrderTemp := bson.D{
								{"packing_order_id", r.PackingOrder.ID},
								{"sales_order_id", i},
								{"product_id", k},
								{"expected_total_pack", v},
								{"pack_type", k2},
								{"status", 1},
							}

							packingSalesOrder = append(packingSalesOrder, packingSalesOrderTemp)
						}
					}
				}
			}
		}

		/*
			product exist but
			pack no longer exist
		*/
		for k, v := range r.MapOldExistData {
			for k2, _ := range v {
				if _, ok := r.MapResProductPackType[k][k2]; !ok {
					filterUpdate := bson.D{
						{"packing_order_id", r.PackingOrder.ID},
						{"product_id", k},
						{"pack_type", k2},
					}

					updatePayload := bson.D{
						{"expected_total_pack", 0},
					}

					if e = m.UpdateOneDataWithFilter("Packing_Item", filterUpdate, updatePayload); e != nil {
						m.DisconnectMongoClient()
						return nil, e
					}

					if e = m.UpdateOneDataWithFilter("Packing_Sales_Order", filterUpdate, updatePayload); e != nil {
						m.DisconnectMongoClient()
						return nil, e
					}
				}
			}
		}

		if len(packingRecommendation) != 0 {
			e = m.InsertManyData("Packing_Item", packingRecommendation)
			if e != nil {
				m.DisconnectMongoClient()
				return nil, e
			}
		}

		if len(packingSalesOrder) != 0 {
			e = m.InsertManyData("Packing_Sales_Order", packingSalesOrder)
			if e != nil {
				m.DisconnectMongoClient()
				return nil, e
			}
		}

		if e = log.AuditLogByUser(r.Session.Staff, r.PackingOrder.ID, "packing_order", "regenetare_packing_order", ""); e != nil {
			m.DisconnectMongoClient()
			o.Rollback()
			return nil, e
		}

		m.DisconnectMongoClient()
		return r.PackingOrder, nil
	} else {
		//generate codes for document
		codePacking, _ := util.GenerateDocCode("PC", r.Warehouse.Code, "packing_order")
		o.Begin()

		pc = &model.PackingOrder{
			Code:         codePacking,
			Area:         r.Warehouse.Area,
			Warehouse:    r.Warehouse,
			DeliveryDate: r.DeliveryDateTime,
			Status:       1,
			Note:         r.Note,
		}

		if _, e = o.Insert(pc); e != nil {
			m.DisconnectMongoClient()
			o.Rollback()
			return nil, e
		}

		if e = GeneratePackingRecursive(r); e != nil {
			m.DisconnectMongoClient()
			o.Rollback()
			return nil, e
		}

		// region insert packing to mongo db
		for k, v := range r.MapResProductPackType {
			for k2, v2 := range v {
				if k2 == -1 {
					packingRecommendationTemp := bson.D{
						{"packing_order_id", pc.ID},
						{"product_id", k},
						{"pack_type", k2},
						{"expected_total_pack", v2},
						{"actual_total_pack", 0},
						{"weight_pack", r.MapResProductPackType[k][-1]},
						{"status", 1},
					}

					packingRecommendation = append(packingRecommendation, packingRecommendationTemp)

					for i, v := range r.MapPackingSalesOrder[k][k2] {

						packingSalesOrderTemp := bson.D{
							{"packing_order_id", pc.ID},
							{"product_id", k},
							{"pack_type", k2},
							{"expected_total_pack", v},
							{"sales_order_id", i},
							{"status", 1},
						}

						packingSalesOrder = append(packingSalesOrder, packingSalesOrderTemp)
					}

					continue
				}

				packingRecommendationTemp := bson.D{
					{"packing_order_id", pc.ID},
					{"product_id", k},
					{"pack_type", k2},
					{"expected_total_pack", v2},
					{"actual_total_pack", 0},
					{"weight_pack", r.MapProductPack[k][k2]},
					{"status", 1},
				}

				packingRecommendation = append(packingRecommendation, packingRecommendationTemp)

				for i, v := range r.MapPackingSalesOrder[k][k2] {
					packingSalesOrderTemp := bson.D{
						{"packing_order_id", pc.ID},
						{"sales_order_id", i},
						{"product_id", k},
						{"expected_total_pack", v},
						{"pack_type", k2},
						{"status", 1},
					}

					packingSalesOrder = append(packingSalesOrder, packingSalesOrderTemp)
				}
			}
		}

		// region if there is no sku list
		if len(packingRecommendation) == 0 {
			m.DisconnectMongoClient()
			o.Commit()
			return
		}

		if len(packingSalesOrder) == 0 {
			m.DisconnectMongoClient()
			o.Commit()
			return
		}
		// endregion

		e = m.InsertManyData("Packing_Item", packingRecommendation)
		if e != nil {
			m.DisconnectMongoClient()
			return nil, e
		}

		e = m.InsertManyData("Packing_Sales_Order", packingSalesOrder)
		if e != nil {
			m.DisconnectMongoClient()
			return nil, e
		}
		// endregion

		if e = log.AuditLogByUser(r.Session.Staff, pc.ID, "packing_order", "generate_packing_order", ""); e != nil {
			m.DisconnectMongoClient()
			o.Rollback()
			return nil, e
		}

		m.DisconnectMongoClient()
		o.Commit()
	}

	return
}

// GeneratePackingRecursive: function to construct map packing
func GeneratePackingRecursive(r generatePackingRequest) (e error) {
loopSku:
	for _, v := range r.SalesOrderItems {
		for _, v2 := range r.PackType {
			if v.OrderQty >= r.MapProductPack[v.ProductID][v2] {
				v.OrderQty = v.OrderQty - r.MapProductPack[v.ProductID][v2]
				v.OrderQty = math.Round(v.OrderQty*100) / 100

				r.MapResProductPackType[v.ProductID][v2] += 1
				r.MapPackingSalesOrder[v.ProductID][v2][v.SalesOrderID] += 1

				continue loopSku
			}
			if v.OrderQty < r.PackType[len(r.PackType)-1] && v.OrderQty > 0 {
				r.MapResProductPackType[v.ProductID][-1] += v.OrderQty
				r.MapPackingSalesOrder[v.ProductID][-1][v.SalesOrderID] += v.OrderQty

				v.OrderQty = 0
				continue loopSku

			}
		}
	}

	// region for checking the order_qty != 0
	var salesOrderItemTemps []*SalesOrderItemByProduct
	for _, pr := range r.SalesOrderItems {
		if pr.OrderQty > 0 {
			salesOrderItemTemps = append(salesOrderItemTemps, pr)
		}
	}
	r.SalesOrderItems = salesOrderItemTemps
	// endregion

	// region base case for the recursive call
	if len(r.SalesOrderItems) != 0 {
		return GeneratePackingRecursive(r)
	}
	// endregion

	return nil
}

// UpdatePackingPack: function for update amount of packing pack
func UpdatePackingPack(r UpdatePackRequest) (pc *model.ResponseData, e error) {
	m := mongodb.NewMongo()

	filter := bson.D{
		{"packing_order_id", r.PackingOrder.ID},
		{"product_id", r.Product.ID},
		{"pack_type", r.PackType},
	}

	var res []byte
	if res, e = m.GetOneDataWithFilter("Packing_Item", filter); e != nil {
		return nil, e
	}

	// region convert byte data to json data
	if e = json.Unmarshal(res, &pc); e != nil {
		return nil, e
	}
	// endregion

	pc.ActualTotalPack += 1
	if e = m.UpdateOneDataWithFilter("Packing_Item", filter, pc); e != nil {
		return nil, e
	}

	// region return data product & packing
	pc.PackingOrder = r.PackingOrder
	pc.Product = r.Product
	// endregion

	m.DisconnectMongoClient()

	return pc, nil
}

// DisposePackingPack: function for dispose packing pack
func DisposePackingPack(r DisposePackRequest) (pc *model.ResponseData, e error) {
	m := mongodb.NewMongo()

	filter := bson.D{
		{"packing_order_id", r.PackingOrder.ID},
		{"product_id", r.Product.ID},
		{"pack_type", r.PackType},
		{"status", 1},
	}

	var res1 []byte
	if res1, e = m.GetOneDataWithFilter("Packing_Item", filter); e != nil {
		return nil, e
	}

	// region convert byte data to json data
	if e = json.Unmarshal(res1, &pc); e != nil {
		return nil, e
	}
	// endregion

	pc.ActualTotalPack -= 1
	if e = m.UpdateOneDataWithFilter("Packing_Item", filter, pc); e != nil {
		return nil, e
	}

	opts := &options.FindOneOptions{}
	opts.SetSort(bson.D{{"code", -1}})

	var rd = new(model.BarcodeModel)
	var res []byte
	if res, e = m.GetOneDataWithFilter("Packing_Barcode", filter, opts); e != nil {
		return nil, e
	}

	if len(res) == 0 {
		pc.CodePrint = ""
	} else {
		// region convert byte data to json data
		if e = json.Unmarshal(res, &rd); e != nil {

			return nil, e
		}

		filterUpdate := bson.D{
			{"code", rd.Code},
		}
		rd.Status = 3
		rd.DeletedBy = r.Session.Staff.ID
		rd.DeletedAt = time.Now().Format(("2006-01-02 15:04:05"))
		if e = m.UpdateOneDataWithFilter("Packing_Barcode", filterUpdate, rd); e != nil {
			return nil, e
		}

		var rd2 = new(model.BarcodeModel)
		if res, e = m.GetOneDataWithFilter("Packing_Barcode", filter, opts); e != nil {
			return nil, e
		}
		if len(res) == 0 {
			pc.CodePrint = ""
			pc.PackingOrder = r.PackingOrder
			pc.Product = r.Product
			// endregion

			m.DisconnectMongoClient()

			return pc, nil
		}
		// region convert byte data to json data
		if e = json.Unmarshal(res, &rd2); e != nil {
			m.DisconnectMongoClient()
			return nil, e
		}
		pc.CodePrint = rd2.Code
	}

	// region return data product & packing
	pc.PackingOrder = r.PackingOrder
	pc.Product = r.Product
	// endregion

	m.DisconnectMongoClient()

	return pc, nil
}

// DetailPack: function for display spesific sku packing
func DetailPack(r UpdatePackRequest) (rd *model.ResponseData, e error) {
	m := mongodb.NewMongo()

	rd = new(model.ResponseData)
	filter := bson.D{
		{"packing_order_id", r.PackingOrder.ID},
		{"product_id", r.Product.ID},
		{"pack_type", r.PackType},
	}

	var res []byte
	if res, e = m.GetOneDataWithFilter("Packing_Item", filter); e != nil {
		return nil, e
	}

	// region convert byte data to json data
	if e = json.Unmarshal(res, &rd); e != nil {
		return nil, e
	}
	// endregion
	rd.PackingOrder = r.PackingOrder
	rd.Product = r.Product

	m.DisconnectMongoClient()

	return rd, nil
}

func GetDetailedPack(s *auth.SessionData, data *model.PackingOrder) (filePath string, err error) {
	var file *xlsx.File
	var sheet *xlsx.Sheet
	var row *xlsx.Row
	var filename string
	o := orm.NewOrm()
	o.Begin()

	dir := util.ExportDirectory
	if data.Warehouse.Name == "All Warehouse" {
		filename = fmt.Sprintf("PackDetail_%s_%s.xlsx", time.Now().Format("2006-01-02"), util.GenerateRandomDoc(5))
	} else {
		filename = fmt.Sprintf("PackDetail_%s_%s_%s.xlsx", time.Now().Format("2006-01-02"), util.ReplaceSpace(data.Warehouse.Name), util.GenerateRandomDoc(5))
	}

	fileDir := fmt.Sprintf("%s/%s", dir, filename)

	file = xlsx.NewFile()
	if sheet, err = file.AddSheet("Sheet1"); err == nil {
		row = sheet.AddRow()
		row.SetHeight(20)
		row.AddCell().Value = "Packing_Order_Code"
		row.AddCell().Value = "Warehouse"
		row.AddCell().Value = "Packing_Date"
		row.AddCell().Value = "Product_Code"
		row.AddCell().Value = "Product_Name"
		row.AddCell().Value = "UOM"
		row.AddCell().Value = "Pack(0.25)"
		row.AddCell().Value = "Pack(0.5)"
		row.AddCell().Value = "Pack(1)"
		row.AddCell().Value = "Pack(2)"
		row.AddCell().Value = "Pack(5)"
		row.AddCell().Value = "Pack(10)"
		row.AddCell().Value = "Pack(20)"

		for _, v := range data.PackingRecommendation {
			row = sheet.AddRow()
			row.AddCell().Value = data.Code
			row.AddCell().Value = data.Warehouse.Name
			row.AddCell().Value = data.DeliveryDate.Format(("2006-01-02"))
			row.AddCell().Value = v.Product.Code
			row.AddCell().Value = v.Product.Name
			row.AddCell().Value = v.Product.Uom.Name
			sort.Slice(v.ProductPack, func(i, j int) bool {
				return v.ProductPack[i].PackType < v.ProductPack[j].PackType
			})
			for _, v2 := range v.ProductPack {
				row.AddCell().SetInt(int(v2.ExpectedTotalPack))
			}
		}

	}
	err = file.Save(fileDir)
	filePath, err = util.UploadToS3(filename, fileDir, "application/xlsx")

	if err = log.AuditLogByUser(s.Staff, data.ID, "packing_order", "export", ""); err != nil {
		o.Rollback()
		return "", err
	}
	o.Commit()

	// fungsi ini berguna untuk menghapus kembali file yang sudah diupload ke server
	os.Remove(fileDir)
	return
}
