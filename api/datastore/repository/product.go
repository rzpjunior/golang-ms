// Copyright 2020 PT. Eden Pangan Indonesia. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

package repository

import (
	"strings"

	"git.edenfarm.id/project-version2/api/util"
	"git.edenfarm.id/project-version2/datamodel/model"

	"git.edenfarm.id/cuxs/orm"
)

// GetProduct find a single data using field and value condition.
func GetProduct(field string, values ...interface{}) (*model.Product, error) {
	var where string
	var arrValueTagProduct []interface{}
	var excludeArchetypeArr []string
	var qMark string

	m := new(model.Product)
	o := orm.NewOrm()
	o.Using("read_only")

	if err := o.QueryTable(m).Filter(field, values...).RelatedSel().Limit(1).One(m); err != nil {
		return nil, err
	}

	o.LoadRelated(m, "ProductImage", 1)

	//to create dynamic question mark for arrValueTagProduct
	arrTagProduct := strings.Split(m.TagProduct, ",")
	for _, v := range arrTagProduct {
		where = where + "? ,"
		arrValueTagProduct = append(arrValueTagProduct, v)
	}
	where = strings.TrimSuffix(where, " ,")

	if err := o.Raw("select group_concat(name) from tag_product where value in ("+where+") order by id asc", arrValueTagProduct).QueryRow(&m.TagProductStr); err != nil {
		return nil, err
	}
	excludeArchetypeArr = strings.Split(m.ExcludeArchetype, ",")
	if m.ExcludeArchetype != "" {
		qMark = ""
		for _, _ = range excludeArchetypeArr {
			qMark = qMark + "?,"
		}
		qMark = strings.TrimSuffix(qMark, ",")
		o.Raw("select group_concat(name) from archetype where id in ("+qMark+")", excludeArchetypeArr).QueryRow(&m.ExcludeArchetypeStr)
	}

	m.ExcludeArchetype = util.EncIdInStr(m.ExcludeArchetype)

	if _, err := o.QueryTable(new(model.Warehouse)).RelatedSel("area").Filter("id__in", strings.Split(m.WarehouseStoStr, ",")).All(&m.WarehouseSto); err != nil {
		return nil, err
	}
	if _, err := o.QueryTable(new(model.Warehouse)).RelatedSel("area").Filter("id__in", strings.Split(m.WarehousePurStr, ",")).All(&m.WarehousePur); err != nil {
		return nil, err
	}

	if _, err := o.QueryTable(new(model.Warehouse)).RelatedSel("area").Filter("id__in", strings.Split(m.WarehouseSalStr, ",")).All(&m.WarehouseSal); err != nil {
		return nil, err
	}

	m.WarehouseStoStr = util.EncIdInStr(m.WarehouseStoStr)
	m.WarehousePurStr = util.EncIdInStr(m.WarehousePurStr)
	m.WarehouseSalStr = util.EncIdInStr(m.WarehouseSalStr)

	m.WarehouseStoArr = []string{}
	m.WarehousePurArr = []string{}
	m.WarehouseSalArr = []string{}

	if m.WarehouseStoStr != "" {
		m.WarehouseStoArr = strings.Split(m.WarehouseStoStr, ",")
	}

	if m.WarehousePurStr != "" {
		m.WarehousePurArr = strings.Split(m.WarehousePurStr, ",")
	}

	if m.WarehouseSalStr != "" {
		m.WarehouseSalArr = strings.Split(m.WarehouseSalStr, ",")
	}

	// region get grand parent and parent category
	if m.Category.GrandParentID != 0 {
		m.GrandParent = &model.Category{
			ID: m.Category.GrandParentID,
		}
		m.GrandParent.Read("ID")
	}

	if m.Category.ParentID != 0 {
		m.Parent = &model.Category{
			ID: m.Category.ParentID,
		}
		m.Parent.Read("ID")
	}
	// endregion

	if m.OrderChannelRestriction != "" {
		qMark := ""
		orderChannelArr := strings.Split(m.OrderChannelRestriction, ",")
		for range orderChannelArr {
			qMark += "?,"
		}
		qMark = strings.TrimSuffix(qMark, ",")
		_, err := o.Raw("SELECT * FROM glossary WHERE attribute = 'order_channel' AND value_int IN ("+qMark+")", orderChannelArr).QueryRows(&m.OrderChannelsRestriction)
		if err != nil {
			return nil, err
		}
	}
	return m, nil
}

// GetProducts : function to get data from database based on parameters
func GetProducts(rq *orm.RequestQuery, tagProduct string) (m []*model.Product, total int64, err error) {
	q, _ := rq.QueryReadOnly(new(model.Product))

	cond := q.GetCond()
	if tagProduct != "" {
		cond1 := orm.NewCondition()
		cond1 = cond1.And("tag_product__icontains", ","+tagProduct+",").Or("tag_product__istartswith", tagProduct+",").Or("tag_product__iendswith", ","+tagProduct).Or("tag_product", tagProduct)

		cond = cond.AndCond(cond1)
	}

	q = q.SetCond(cond)

	if total, err = q.Exclude("status", 3).Count(); err != nil || total == 0 {
		return nil, total, err
	}

	if _, err = q.Exclude("status", 3).RelatedSel().All(&m, rq.Fields...); err == nil {
		o := orm.NewOrm()
		o.Using("read_only")
		for _, v := range m {
			var where string
			var arrValueTagProduct []interface{}

			//to create dynamic question mark for arrValueTagProduct
			arrTagProduct := strings.Split(v.TagProduct, ",")
			for _, v := range arrTagProduct {
				where = where + "? ,"
				arrValueTagProduct = append(arrValueTagProduct, v)
			}
			where = strings.TrimSuffix(where, " ,")

			if err := o.Raw("select group_concat(name) from tag_product where value in ("+where+") order by id asc", arrValueTagProduct...).QueryRow(&v.TagProductStr); err != nil {
				return nil, 0, err
			}
		}
		return m, total, nil
	}

	return m, total, err
}

// GetFilterProducts : function to get data from database based on parameters with filtered permission
func GetFilterProducts(rq *orm.RequestQuery) (m []*model.Product, total int64, err error) {
	q, _ := rq.QueryReadOnly(new(model.Product))

	if total, err = q.Exclude("status", 3).Count(); err != nil || total == 0 {
		return nil, 0, err
	}

	if _, err = q.Exclude("status", 3).All(&m, rq.Fields...); err == nil {
		return m, total, nil
	}

	return nil, total, err
}

// ValidProduct : function to check if id is valid in database
func ValidProduct(id int64) (product *model.Product, e error) {
	product = &model.Product{ID: id}
	e = product.Read("ID")

	return
}

// CheckProductData : function to check data based on filter and exclude parameters
func CheckProductData(filter, exclude map[string]interface{}) (p []*model.Product, total int64, err error) {
	rq := orm.RequestQuery{}
	o, _ := rq.QueryReadOnly(new(model.Product))

	for k, v := range filter {
		o = o.Filter(k, v)
	}

	for k, v := range exclude {
		o = o.Exclude(k, v)
	}

	if total, err := o.All(&p); err == nil {
		return p, total, nil
	}

	return nil, 0, err
}

// CheckWarehouseProduct : function to check if product id and warehouse id is in list of column provided
func CheckWarehouseProduct(productID int64, column string, warehousesID ...string) (resultId int64, err error) {
	o := orm.NewOrm()
	o.Using("read_only")

	qb, _ := orm.NewQueryBuilder("mysql")

	var where string

	qb.Select("id").From("product")

	if column != "warehouse_sto" {
		for _, v := range warehousesID {
			where = where + "find_in_set(" + v + ", " + column + ") or "
		}
		where = "and (" + strings.TrimRight(where, " or ") + ")"
	} else {
		for _, v := range warehousesID {
			where = where + "and find_in_set(" + v + ", " + column + ") "
		}
	}

	qb.Where("id = ? " + where).Limit((1))
	q := qb.String()

	if err = o.Raw(q, productID).QueryRow(&resultId); err == nil {
		return resultId, nil
	}

	return 0, err
}

// GetProductNoDecrypt find a single data using field and value condition without decrypt for salable,purchasability and storability.
func GetProductNoDecrypt(field string, values ...interface{}) (*model.Product, error) {
	m := new(model.Product)
	o := orm.NewOrm()
	o.Using("read_only")

	if err := o.QueryTable(m).Filter(field, values...).RelatedSel().Limit(1).One(m); err != nil {
		return nil, err
	}

	return m, nil
}
