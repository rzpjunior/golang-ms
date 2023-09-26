// Copyright 2022 PT. Eden Pangan Indonesia. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

package repository

import (
	"git.edenfarm.id/project-version2/datamodel/model"

	"git.edenfarm.id/cuxs/orm"
)

// GetDebitNote find a single data debit note using field and value condition.
func GetDebitNote(field string, values ...interface{}) (*model.DebitNote, error) {
	m := new(model.DebitNote)
	o := orm.NewOrm()
	o.Using("read_only")

	if err := o.QueryTable(m).Filter(field, values...).RelatedSel(2).Limit(1).One(m); err != nil {
		return nil, err
	}
	o.LoadRelated(m, "DebitNoteItems", 2)

	return m, nil
}

// GetDebitNotes : function to get data from database based on parameters
func GetDebitNotes(rq *orm.RequestQuery) (m []*model.DebitNote, total int64, err error) {
	o := orm.NewOrm()
	o.Using("read_only")

	q, _ := rq.QueryReadOnly(new(model.DebitNote))

	if total, err = q.Count(); err != nil || total == 0 {
		return nil, total, err
	}

	var mx []*model.DebitNote
	if _, err = q.All(&mx, rq.Fields...); err != nil {
		return nil, total, err
	}

	return mx, total, nil
}

// ValidDebitNote : function to check if id is valid in database
func ValidDebitNote(id int64) (dn *model.DebitNote, e error) {
	dn = &model.DebitNote{ID: id}
	e = dn.Read("ID")

	return
}
