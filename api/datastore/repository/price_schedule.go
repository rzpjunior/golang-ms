package repository

import (
	"git.edenfarm.id/cuxs/orm"
	"git.edenfarm.id/project-version2/datamodel/model"
)

// GetPriceSchedule find a single data price schedule set using field and value condition.
func GetPriceSchedule(field string, values ...interface{}) (*model.PriceSchedule, error) {
	m := new(model.PriceSchedule)
	o := orm.NewOrm()
	o.Using("read_only")

	if err := o.QueryTable(m).Filter(field, values...).RelatedSel().Limit(1).One(m); err != nil {
		return nil, err
	}

	o.LoadRelated(m, "PriceScheduleDumps", 1)

	return m, nil
}

// GetPriceSchedules : function to get data from database based on parameters
func GetPriceSchedules(rq *orm.RequestQuery) (m []*model.PriceSchedule, total int64, err error) {
	q, _ := rq.QueryReadOnly(new(model.PriceSchedule))

	if total, err = q.Exclude("status", 4).Count(); err != nil || total == 0 {
		return nil, total, err
	}

	var mx []*model.PriceSchedule
	if _, err = q.Exclude("status", 4).All(&mx, rq.Fields...); err == nil {
		return mx, total, nil
	}

	return nil, total, err
}
