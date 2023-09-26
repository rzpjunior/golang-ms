package repository

import (
	"git.edenfarm.id/cuxs/orm"
	"git.edenfarm.id/project-version2/datamodel/model"
)

// GetPickingOrderItems : function to get data from database based on parameters
func GetPickingOrderItems(rq *orm.RequestQuery) (m []*model.PickingOrderItem, total int64, err error) {
	q, _ := rq.QueryReadOnly(new(model.PickingOrderItem))

	if total, err = q.Exclude("PickingOrderAssign__Status", 7).Count(); err != nil || total == 0 {
		return nil, total, err
	}

	var mx []*model.PickingOrderItem
	if _, err = q.RelatedSel(2).Exclude("PickingOrderAssign__Status", 7).All(&mx, rq.Fields...); err == nil {
		return mx, total, nil
	}

	return nil, total, err
}

// GetPickingOrderItem find a single data
func GetPickingOrderItem(field string, values ...interface{}) (*model.PickingOrderItem, error) {
	m := new(model.PickingOrderItem)
	o := orm.NewOrm()
	o.Using("read_only")

	if err := o.QueryTable(m).Filter(field, values...).RelatedSel("PickingOrderAssign").RelatedSel("Product").RelatedSel("Product__Uom").Limit(1).One(m); err != nil {
		return nil, err
	}

	o.LoadRelated(m.Product, "ProductImage", 1)
	return m, nil
}

// GetItemIfExistByProductId : get all item data based on picking assign id and product id
func GetItemIfExistByProductId(pickingOrderAssignId int64, productId int64) (isExist bool, err error) {
	o := orm.NewOrm()
	o.Using("read_only")

	var count int
	q := "SELECT EXISTS(select poi.id from picking_order_item poi join picking_order_assign poa on poa.id = poi.picking_order_assign_id where poi.product_id = ? and picking_order_assign_id = ?) as isExists"
	o.Raw(q, productId, pickingOrderAssignId).QueryRow(&count)
	if count != 0 {
		isExist = true
		return isExist, err
	}

	return isExist, err
}

// CheckPickingOrderItemData : function to check PickingOrderItem data based on filter and exclude parameters
func CheckPickingOrderItemData(filter, exclude map[string]interface{}) (PickingOrderItem []*model.PickingOrderItem, total int64, err error) {
	rq := orm.RequestQuery{}
	o, _ := rq.QueryReadOnly(new(model.PickingOrderItem))

	for k, v := range filter {
		o = o.Filter(k, v)
	}

	for k, v := range exclude {
		o = o.Exclude(k, v)
	}

	if total, err = o.All(&PickingOrderItem); err == nil {
		return PickingOrderItem, total, nil
	}

	return nil, 0, err
}
