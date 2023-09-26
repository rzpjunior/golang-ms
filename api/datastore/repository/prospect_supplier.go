package repository

import (
	"git.edenfarm.id/project-version2/datamodel/model"

	"git.edenfarm.id/cuxs/orm"
)

// GetProspectSupplier find a single data price set using field and value condition.
func GetProspectSupplier(field string, values ...interface{}) (*model.ProspectSupplier, error) {
	m := new(model.ProspectSupplier)
	o := orm.NewOrm()
	o.Using("read_only")

	if err := o.QueryTable(m).Filter(field, values...).RelatedSel().Limit(1).One(m); err != nil {
		return nil, err
	}
	m.SubDistrict.District.City.Province.Read()

	return m, nil
}

// GetProspectSuppliers : function to get data from database based on parameters
func GetProspectSuppliers(rq *orm.RequestQuery) (m []*model.ProspectSupplier, total int64, err error) {
	q, _ := rq.QueryReadOnly(new(model.ProspectSupplier))

	if total, err = q.Exclude("reg_status", 4).Count(); err != nil || total == 0 {
		return nil, total, err
	}

	var mx []*model.ProspectSupplier
	if _, err = q.Exclude("reg_status", 4).All(&mx, rq.Fields...); err == nil {
		return mx, total, nil
	}

	return nil, total, err
}

// GetFilterProspectSuppliers : function to get data from database based on parameters with filtered permission
func GetFilterProspectSuppliers(rq *orm.RequestQuery) (m []*model.ProspectSupplier, total int64, err error) {
	q, _ := rq.QueryReadOnly(new(model.ProspectSupplier))

	if total, err = q.Filter("reg_status", 4).Count(); err != nil || total == 0 {
		return nil, total, err
	}

	var mx []*model.ProspectSupplier
	if _, err = q.Filter("reg_status", 4).All(&mx, rq.Fields...); err == nil {
		return mx, total, nil
	}

	return nil, total, err
}

func ValidProspectSupplier(id int64) (prospectsupplier *model.ProspectSupplier, e error) {
	prospectsupplier = &model.ProspectSupplier{ID: id}
	e = prospectsupplier.Read("ID")

	return
}
