// Copyright 2020 PT. Eden Pangan Indonesia. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

package repository

import (
	"git.edenfarm.id/project-version2/datamodel/model"

	"git.edenfarm.id/cuxs/orm"
)

// GetJob find a single data job using field and value condition.
func GetJob(field string, values ...interface{}) (*model.Jobs, error) {
	m := new(model.Jobs)
	o := orm.NewOrm()
	o.Using("read_only")

	if err := o.QueryTable(m).Filter(field, values...).RelatedSel().Limit(1).One(m); err != nil {
		return nil, err
	}
	return m, nil
}

// GetJobs get all jobs from database
func GetJobs(rq *orm.RequestQuery) (m []*model.Jobs, total int64, err error) {
	// make new orm query
	q, _ := rq.QueryReadOnly(new(model.Jobs))

	// get total data
	if total, err = q.Count(); err != nil || total == 0 {
		return nil, total, err
	}

	// get data requested
	var mx []*model.Jobs
	if _, err = q.All(&mx, rq.Fields...); err == nil {
		return mx, total, nil
	}

	// return error some thing went wrong
	return nil, total, err
}
