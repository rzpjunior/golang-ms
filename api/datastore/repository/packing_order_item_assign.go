package repository

import (
	"time"

	"git.edenfarm.id/project-version2/datamodel/model"

	"git.edenfarm.id/cuxs/orm"
)

// GetPackingOrderItemAssign : function to get data from database based on parameters
func GetPackingOrderItemAssign(rq *orm.RequestQuery) (m []*model.PackingOrderItemAssign, total int64, err error) {
	q, _ := rq.QueryReadOnly(new(model.PackingOrderItemAssign))

	if total, err = q.Count(); err != nil || total == 0 {
		return nil, total, err
	}

	var mx []*model.PackingOrderItemAssign
	if _, err = q.All(&mx, rq.Fields...); err == nil {
		return mx, total, nil
	}

	return nil, total, err
}

// GetPackingOrderItemAssign : function to get data from database based on parameters
func GetPackingOrderItemAssignedToPacker(rq *orm.RequestQuery) (m []*model.PackingOrderItemAssign, total int64, err error) {
	var e error

	o := orm.NewOrm()
	o.Using("read_only")

	q, _ := rq.QueryReadOnly(new(model.PackingOrderItemAssign))

	if total, err = q.Filter("packingorderitem__packingorder__deliverydate", time.Now().AddDate(0, 0, 1).Format("2006-01-02")).
		Filter("packingorderitem__packingorder__status", 1).Count(); err != nil || total == 0 {
		return nil, total, err
	}

	var mx []*model.PackingOrderItemAssign
	if _, err = q.RelatedSel(2).Filter("packingorderitem__packingorder__deliverydate", time.Now().AddDate(0, 0, 1).Format("2006-01-02")).
		Filter("packingorderitem__packingorder__status", 1).All(&mx, rq.Fields...); err == nil {
		for _, v := range mx {
			if _, e = o.Raw("SELECT * from product_image where product_id = ?", v.PackingOrderItem.Product.ID).QueryRows(&v.PackingOrderItem.Product.ProductImage); e != nil {

			}
		}
		return mx, total, nil
	}

	return nil, total, err
}
