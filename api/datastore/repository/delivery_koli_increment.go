// Copyright 2020 PT. Eden Pangan Indonesia. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

package repository

import (
	"git.edenfarm.id/project-version2/datamodel/model"

	"git.edenfarm.id/cuxs/orm"
)

// GetDeliveryKoliIncrement find a single data delivery koli increment using field and value condition.
func GetDeliveryKoliIncrement(field string, values ...interface{}) (*model.DeliveryKoliIncrement, error) {
	m := new(model.DeliveryKoliIncrement)
	o := orm.NewOrm()
	o.Using("read_only")

	if err := o.QueryTable(m).Filter(field, values...).RelatedSel().Limit(1).One(m); err != nil {
		return nil, err
	}
	return m, nil
}

// GetDeliveryKoliIncrements get all data delivery koli increment that matched with query request parameters.
func GetDeliveryKoliIncrements(rq *orm.RequestQuery) (m []*model.DeliveryKoliIncrement, total int64, err error) {
	// make new orm query
	q, _ := rq.QueryReadOnly(new(model.DeliveryKoliIncrement))

	// get total data
	if total, err = q.Count(); err != nil || total == 0 {
		return nil, total, err
	}

	// get data requested
	var mx []*model.DeliveryKoliIncrement
	if _, err = q.All(&mx, rq.Fields...); err != nil {
		return nil, total, err
	}

	for _, v := range mx {
		v.SalesOrder.Read("ID")
	}
	return mx, total, nil
}

func ValidDeliveryKoliIncrement(id int64) (dki *model.DeliveryKoliIncrement, e error) {
	dki = &model.DeliveryKoliIncrement{ID: id}
	e = dki.Read("ID")

	return
}
