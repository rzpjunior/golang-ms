package repository

import (
	"git.edenfarm.id/cuxs/orm"
	"git.edenfarm.id/project-version2/datamodel/model"
)

// GetCourierByDeliveryOrderID : get courier transaction by DO ID
func GetCourierByDeliveryOrderID(values ...interface{}) (ct *model.CourierTransaction, err error) {
	m := new(model.CourierTransaction)
	o := orm.NewOrm()
	o.Using("read_only")
	if err := o.QueryTable(m).Filter("delivery_order_id", values).RelatedSel().Limit(1).One(m); err != nil {
		return nil, err
	}
	return m, nil
}
