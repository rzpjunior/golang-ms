// Copyright 2020 PT. Eden Pangan Indonesia. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

package repository

import (
	"encoding/json"
	"fmt"
	"git.edenfarm.id/cuxs/mongodb"
	"git.edenfarm.id/cuxs/orm"
	"git.edenfarm.id/project-version2/datamodel/model"
	"go.mongodb.org/mongo-driver/bson"
	"sort"
	"strconv"
)

// GetPackingOrder find a single data packing order using field and value condition.
func GetPackingOrder(field string, values ...interface{}) (*model.PackingOrder, error) {
	p := new(model.PackingOrder)
	o := orm.NewOrm()
	o.Using("read_only")
	var e error

	if e = o.QueryTable(p).Filter(field, values...).RelatedSel().Limit(1).One(p); e != nil {
		return nil, e
	}

	return p, nil
}

// GetPackingOrderDetailPack: find a single data packing order using field and value condition.
func GetPackingOrderDetailPack(field string, values ...interface{}) (*model.PackingOrder, error) {
	p := new(model.PackingOrder)
	o := orm.NewOrm()
	o.Using("read_only")
	m := mongodb.NewMongo()
	var e error

	if e = o.QueryTable(p).Filter(field, values...).RelatedSel().Limit(1).One(p); e != nil {
		return nil, e
	}

	filter := bson.D{
		{"packing_order_id", p.ID},
		{"status", 1},
	}

	var res []byte
	if res, e = m.GetAllDataWithFilter("Packing_Item", filter); e != nil {
		return nil, e
	}

	// region convert byte data to json data
	var response []*model.ResponseData
	if e = json.Unmarshal(res, &response); e != nil {
		return nil, e
	}
	// endregion

	var pt []float64
	o.Raw("SELECT value_name FROM glossary WHERE `table` = 'packing_order' and `attribute` = 'pack_size'").QueryRows(&pt)

	// region mapping product per pack type
	mapProductPa := make(map[int64]map[float64]*model.PackAdjustment)
	mapProductPct := make(map[int64]*model.ProductPercentage)
	for _, value := range response {
		mapPa := make(map[float64]*model.PackAdjustment, 0)
		for _, v := range pt {
			pa := new(model.PackAdjustment)
			if _, ok := mapPa[v]; ok {
				continue
			}
			pa.PackType = v
			mapPa[v] = pa
		}
		mapProductPa[value.ProductID] = mapPa
	}

	for _, value := range response {

		if val, ok := mapProductPa[value.ProductID][value.PackType]; ok {
			val.ExpectedTotalPack = value.ExpectedTotalPack
			val.ActualTotalPack = value.ActualTotalPack
		}

		// region percentage
		ppct := new(model.ProductPercentage)
		if valPct, ok := mapProductPct[value.ProductID]; ok {
			valPct.ExpectedTotalPack += value.ExpectedTotalPack
			valPct.ActualTotalPack += value.ActualTotalPack
			continue
		}

		ppct.ExpectedTotalPack = value.ExpectedTotalPack
		ppct.ActualTotalPack = value.ActualTotalPack

		mapProductPct[value.ProductID] = ppct
		// endregion
	}
	// endregion
	var prs []*model.PackingRecommendation
	for k, v := range mapProductPa {
		var pct float64
		pr := new(model.PackingRecommendation)
		product := new(model.Product)
		pr.PackingOrderID = p.ID
		pr.ProductID = k
		for _, v2 := range v {
			pr.ProductPack = append(pr.ProductPack, v2)
		}
		product.ID = k
		product.Read("ID")
		product.Uom.Read("ID")
		pr.Product = product

		// region pct calculation
		pct = (mapProductPct[k].ActualTotalPack / mapProductPct[k].ExpectedTotalPack) * 100
		pctNumb := fmt.Sprintf("%.2f", pct)
		pr.TotalProgressPercentage, _ = strconv.ParseFloat(pctNumb, 10)

		// endregion

		prs = append(prs, pr)
	}
	p.PackingRecommendation = prs
	sort.Slice(p.PackingRecommendation, func(i, j int) bool {
		return p.PackingRecommendation[i].Product.ID < p.PackingRecommendation[j].Product.ID
	})

	m.DisconnectMongoClient()
	return p, nil
}

// GetPackingOrders : function to get data from database based on parameters
func GetPackingOrders(rq *orm.RequestQuery) (m []*model.PackingOrder, total int64, err error) {
	q, _ := rq.QueryReadOnly(new(model.PackingOrder))

	if total, err = q.Count(); err != nil || total == 0 {
		return nil, total, err
	}

	var mx []*model.PackingOrder
	if _, err = q.All(&mx, rq.Fields...); err == nil {
		return mx, total, nil
	}

	return nil, total, err
}

// GetFilterPackingOrders : function to get data from database based on parameters with filtered permission
func GetFilterPackingOrders(rq *orm.RequestQuery) (m []*model.PackingOrder, total int64, err error) {
	q, _ := rq.QueryReadOnly(new(model.PackingOrder))

	if total, err = q.Exclude("status", 3).Count(); err != nil || total == 0 {
		return nil, total, err
	}

	var mx []*model.PackingOrder
	if _, err = q.Exclude("status", 3).All(&mx, rq.Fields...); err == nil {
		return mx, total, nil
	}

	return nil, total, err
}

// ValidPackingOrder : function to check if id is valid in database
func ValidPackingOrder(id int64) (packingOrder *model.PackingOrder, e error) {
	packingOrder = &model.PackingOrder{ID: id}
	e = packingOrder.Read("ID")

	return
}

// ValidPackingOrderItem : function to check if id is valid in database
func ValidPackingOrderItem(id int64) (packingOrderItem *model.PackingOrderItem, e error) {
	packingOrderItem = &model.PackingOrderItem{ID: id}
	e = packingOrderItem.Read("ID")

	return
}

// CheckPackingOrderData : function to check PackingOrder data based on filter and exclude parameters
func CheckPackingOrderData(filter, exclude map[string]interface{}) (PackingOrder []*model.PackingOrder, total int64, err error) {
	rq := orm.RequestQuery{}
	o, _ := rq.QueryReadOnly(new(model.PackingOrder))

	for k, v := range filter {
		o = o.Filter(k, v)
	}

	for k, v := range exclude {
		o = o.Exclude(k, v)
	}

	if total, err = o.All(&PackingOrder); err == nil {
		return PackingOrder, total, nil
	}

	return nil, 0, err
}

// CheckPackingOrderItemData : function to check PackingOrder data based on filter and exclude parameters
func CheckPackingOrderItemData(filter, exclude map[string]interface{}) (PackingOrderItem []*model.PackingOrderItem, total int64, err error) {

	rq := orm.RequestQuery{}
	o, _ := rq.QueryReadOnly(new(model.PackingOrderItem))
	o.RelatedSel(1)

	for k, v := range filter {
		o = o.Filter(k, v)
	}

	for k, v := range exclude {
		o = o.Exclude(k, v)
	}

	if total, err = o.All(&PackingOrderItem); err == nil {
		return PackingOrderItem, total, nil
	}

	return nil, 0, err
}

// CheckPackingOrderItemHelper : function to check every helper in packing order item has total packing (Pack/KG)
func CheckPackingOrderItemHelper(idPackingOrder int64, idStaff ...int64) (packingOrderItemAssign []*model.PackingOrderItemAssign, count int64, e error) {
	rq := orm.RequestQuery{}
	o, _ := rq.QueryReadOnly(new(model.PackingOrderItemAssign))

	if len(idStaff) > 0 {
		o = o.Exclude("staff_id__in", idStaff)
	}
	if count, e := o.Filter("packing_order_item_id", idPackingOrder).Filter("subtotal_weight__gt", 0).All(&packingOrderItemAssign); e == nil {
		return packingOrderItemAssign, count, e
	}

	return nil, 0, e
}
