// Copyright 2022 PT. Eden Pangan Indonesia. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

package repository

import (
	"strconv"

	"git.edenfarm.id/cuxs/orm"
	"git.edenfarm.id/project-version2/datamodel/model"
)

// GetProductGroup find a single data price set using field and value condition.
func GetProductGroup(field string, values ...interface{}) (*model.ProductGroup, error) {
	m := new(model.ProductGroup)
	o := orm.NewOrm()
	o.Using("read_only")

	if err := o.QueryTable(m).Filter(field, values...).RelatedSel().Limit(1).One(m); err != nil {
		return nil, err
	}

	return m, nil
}

func GetProductGroupItem(field string, values ...interface{}) (*model.ProductGroupItem, error) {
	m := new(model.ProductGroupItem)
	o := orm.NewOrm()
	o.Using("read_only")

	if  err := o.QueryTable(m).Filter(field, values...).RelatedSel("ProductGroup").Limit(1).One(m); err != nil {
		return nil, err
	}

	return m, nil
}

// GetProductGroups : function to get data from database based on parameters
func GetProductGroups(rq *orm.RequestQuery) (productGroup []*model.ProductGroup, total int64, err error) {
	q, _ := rq.QueryReadOnly(new(model.ProductGroup))

	if total, err = q.Exclude("status", 3).All(&productGroup, rq.Fields...); err == nil {
		return productGroup, total, nil
	}

	return nil, total, err
}

func GetProductGroupTransferSku(rq *orm.RequestQuery) (stocks []*model.Stock, total int64, err error) {
	
	var productGroupID int64
	var listProductIDFiltered []int64
	var productGroupItems []*model.ProductGroupItem
	productGroupItem := new(model.ProductGroupItem)

	o := orm.NewOrm()
	o.Using("read_only")

	for i,v := range rq.Conditions {
		if productGroupIDStr, ok := v["product.product_group_id"]; ok {
			
			productGroupID, err = strconv.ParseInt(productGroupIDStr, 10, 64)

			if err != nil {
				return nil, 0, err
			}

			rq.Conditions[i] = rq.Conditions[len(rq.Conditions)-1]
			rq.Conditions[len(rq.Conditions)-1] = nil
			rq.Conditions = rq.Conditions[:len(rq.Conditions)-1] 

			break

		}

	}
			
	total, err = o.QueryTable(productGroupItem).Filter("product_group_id", productGroupID).RelatedSel("Product").All(&productGroupItems)

	if err != nil {
		return nil, total, err
	}


	if len(productGroupItems) > 0 {

		for _, v := range productGroupItems {
			listProductIDFiltered = append(listProductIDFiltered, v.Product.ID)
		}

		q, _ := rq.QueryReadOnly(new(model.Stock))

		if total, err = q.RelatedSel(2).Filter("warehouse__status", 1).Filter("product_id__in", listProductIDFiltered).All(&stocks, rq.Fields...); err == nil {
			return stocks, total, nil
		}

	}

	return nil, total, err
}

func CheckValidProductByProductGroup(productId int64, productGroupId int64) (pgi *model.ProductGroupItem, valid bool, err error) {
	m := new(model.ProductGroupItem)
	o := orm.NewOrm()
	o.Using("read_only")

	if err := o.QueryTable(m).Filter("product_id", productId).Filter("product_group_id", productGroupId).Limit(1).One(m); err != nil {
		return nil, false, err
	}

	if m == nil {
		return nil, false, err
	}

	return m, true, nil
}