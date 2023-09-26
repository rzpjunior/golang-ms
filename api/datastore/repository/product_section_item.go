// Copyright 2020 PT. Eden Pangan Indonesia. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

package repository

import (
	"git.edenfarm.id/cuxs/orm"
	"git.edenfarm.id/project-version2/datamodel/model"
)

// GetProductSectionItem : function to get datas product for product section
func GetProductSectionItem(rq *orm.RequestQuery, category string) (m []*model.ProductSectionItem, total int64, err error) {
	q, _ := rq.QueryReadOnly(new(model.Product))
	var products []*model.Product
	o := orm.NewOrm()
	o.Using("read_only")

	cond := q.GetCond()
	if category != "" {
		cond1 := orm.NewCondition()
		cond1 = cond1.And("tag_product__icontains", ","+category+",").Or("tag_product__istartswith", category+",").Or("tag_product__iendswith", ","+category).Or("tag_product", category)

		cond = cond.AndCond(cond1)
	}

	q = q.SetCond(cond)

	if _, err = q.Filter("status", 1).All(&products, rq.Fields...); err == nil {
		for _, v := range products {
			m = append(m, &model.ProductSectionItem{
				ID:   v.ID,
				Code: v.Code,
				Name: v.Name,
			})
		}

		return m, total, nil
	}

	return m, total, err
}
