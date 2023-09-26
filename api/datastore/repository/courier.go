package repository

import (
	"git.edenfarm.id/cuxs/orm"
	"git.edenfarm.id/project-version2/datamodel/model"
)

// GetCouriers get all data koli that matched with query request parameters.
func GetCouriers(rq *orm.RequestQuery) (m []*model.Courier, total int64, err error) {
	// make new orm query
	q, _ := rq.QueryReadOnly(new(model.Courier))

	// get total data
	if total, err = q.Filter("status", 1).Filter("couriervendor__status", 1).Count(); err != nil || total == 0 {
		return nil, total, err
	}

	// get data requested
	var mx []*model.Courier
	if _, err = q.Filter("status", 1).Filter("couriervendor__status", 1).All(&mx, rq.Fields...); err == nil {
		return mx, total, nil
	}

	// return error some thing went wrong
	return nil, total, err
}

func ValidCourier(id int64) (courier *model.Courier, e error) {
	courier = &model.Courier{ID: id}
	e = courier.Read("ID")

	return
}
