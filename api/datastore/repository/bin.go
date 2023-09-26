// Copyright 2020 PT. Eden Pangan Indonesia. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

package repository

import (
	"git.edenfarm.id/project-version2/datamodel/model"

	"git.edenfarm.id/cuxs/orm"
)

// GetBin : find a single data using field and value condition.
func GetBin(field string, values ...interface{}) (*model.Bin, error) {
	m := new(model.Bin)
	o := orm.NewOrm()
	o.Using("read_only")

	if err := o.QueryTable(m).Filter(field, values...).RelatedSel().Limit(1).One(m); err != nil {
		if err := o.QueryTable(m).Filter(field, values...).RelatedSel("warehouse").Limit(1).One(m); err != nil {
			return nil, err
		} else {
			o.Raw("select * from bin_info where id=?", m.Warehouse.BinInfo.ID).QueryRow(&m.Warehouse.BinInfo)
		}
	}

	return m, nil
}

// GetBins : function to get data from database based on parameters
func GetBins(rq *orm.RequestQuery) (m []*model.Bin, total int64, err error) {
	q, _ := rq.QueryReadOnly(new(model.Bin))
	o := orm.NewOrm()
	o.Using("read_only")

	if total, err = q.Count(); err != nil || total == 0 {
		return nil, total, err
	}

	var mx []*model.Bin
	if _, err = q.All(&mx, rq.Fields...); err == nil {
		for _, v := range mx {
			if v.Product != nil {
				o.Raw("select name from product where id=?", v.Product.ID).QueryRow(&v.ProductName)
			}
		}

		return mx, total, nil
	}

	return nil, total, err
}

func CheckBinData(filter, exclude map[string]interface{}) (bin []*model.Bin, total int64, err error) {
	rq := orm.RequestQuery{}
	o, _ := rq.QueryReadOnly(new(model.Bin))

	for k, v := range filter {
		o = o.Filter(k, v)
	}

	for k, v := range exclude {
		o = o.Exclude(k, v)
	}

	if total, err = o.All(&bin); err == nil {
		return bin, total, nil
	}

	return nil, 0, err
}

func ValidBin(id int64) (bin *model.Bin, e error) {
	bin = &model.Bin{ID: id}
	e = bin.Read("ID")

	return
}
