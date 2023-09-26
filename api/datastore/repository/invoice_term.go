// Copyright 2020 PT. Eden Pangan Indonesia. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

package repository

import (
	"git.edenfarm.id/project-version2/datamodel/model"

	"git.edenfarm.id/cuxs/orm"
)

// GetInvoiceTerm find a single data sales term using field and value condition.
func GetInvoiceTerm(field string, values ...interface{}) (*model.InvoiceTerm, error) {
	m := new(model.InvoiceTerm)
	o := orm.NewOrm()
	o.Using("read_only")

	if err := o.QueryTable(m).Filter(field, values...).RelatedSel().Limit(1).One(m); err != nil {
		return nil, err
	}
	return m, nil
}

// GetInvoiceTerms : function to get data from database based on parameters
func GetInvoiceTerms(rq *orm.RequestQuery) (m []*model.InvoiceTerm, total int64, err error) {
	q, _ := rq.QueryReadOnly(new(model.InvoiceTerm))

	if total, err = q.Exclude("status", 3).Count(); err != nil || total == 0 {
		return nil, total, err
	}

	var mx []*model.InvoiceTerm
	if _, err = q.Exclude("status", 3).All(&mx, rq.Fields...); err == nil {
		return mx, total, nil
	}

	return nil, total, err
}

// GetFilterInvoiceTerms : function to get data from database based on parameters with filtered permission
func GetFilterInvoiceTerms(rq *orm.RequestQuery) (m []*model.InvoiceTerm, total int64, err error) {
	q, _ := rq.QueryReadOnly(new(model.InvoiceTerm))

	if total, err = q.Exclude("status", 3).Count(); err != nil || total == 0 {
		return nil, total, err
	}

	var mx []*model.InvoiceTerm
	if _, err = q.Exclude("status", 3).All(&mx, rq.Fields...); err == nil {
		return mx, total, nil
	}

	return nil, total, err
}

// ValidInvoiceTerm : function to check if id is valid in database
func ValidInvoiceTerm(id int64) (invoiceTerm *model.InvoiceTerm, e error) {
	invoiceTerm = &model.InvoiceTerm{ID: id}
	e = invoiceTerm.Read("ID")

	return
}
