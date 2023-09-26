// Copyright 2020 PT. Eden Pangan Indonesia. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

package repository

import (
	"git.edenfarm.id/project-version2/datamodel/model"

	"git.edenfarm.id/cuxs/orm"
)

// GetDeliveryKoli find a single data delivery koli using field and value condition.
func GetDeliveryKoli(field string, values ...interface{}) (*model.DeliveryKoli, error) {
	m := new(model.DeliveryKoli)
	o := orm.NewOrm()
	o.Using("read_only")
	if err := o.QueryTable(m).Filter(field, values...).RelatedSel().Limit(1).One(m); err != nil {
		return nil, err
	}
	return m, nil
}

// GetDeliveryKolis get all data koli that matched with query request parameters.
func GetDeliveryKolis(rq *orm.RequestQuery) (m []*model.DeliveryKoli, total int64, err error) {
	// make new orm query
	q, _ := rq.QueryReadOnly(new(model.DeliveryKoli))

	// get total data
	if total, err = q.Count(); err != nil || total == 0 {
		return nil, total, err
	}

	// get data requested
	var mx []*model.DeliveryKoli
	if _, err = q.All(&mx, rq.Fields...); err == nil {
		return mx, total, nil
	}

	// return error some thing went wrong
	return nil, total, err
}

func ValidDeliveryKoli(id int64) (deliveryKoli *model.DeliveryKoli, e error) {
	deliveryKoli = &model.DeliveryKoli{ID: id}
	e = deliveryKoli.Read("ID")

	return
}
