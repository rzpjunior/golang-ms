package repository

import (
	"git.edenfarm.id/project-version2/datamodel/model"

	"git.edenfarm.id/cuxs/orm"
)

// GetBoxItem find a single data Product Box using field and value condition.
func GetBoxItem(field string, values ...interface{}) (*model.BoxItem, error) {
	m := new(model.BoxItem)
	o := orm.NewOrm()
	o.Using("read_only")

	if err := o.QueryTable(m).Filter(field, values...).RelatedSel().Limit(1).One(m); err != nil {
		return nil, err
	}
	return m, nil
}

// GetBoxItems : function to get data from database based on parameters
func GetBoxItems(rq *orm.RequestQuery) (m []*model.BoxItem, total int64, err error) {
	q, _ := rq.QueryReadOnly(new(model.BoxItem))

	if total, err = q.Exclude("status", 3).Count(); err != nil || total == 0 {
		return nil, total, err
	}

	var mx []*model.BoxItem
	if _, err = q.Exclude("status", 3).All(&mx, rq.Fields...); err == nil {
		for _, boxItem := range mx {
			boxItem.Box.Read()
			boxItem.Product.Read()
			boxItem.Product.Uom.Read()
		}
		return mx, total, nil
	}

	return nil, total, err
}

// GetBoxItems : function to get data from database based on parameters
func GetBoxFridgeItems(rq *orm.RequestQuery, statusFilter string) (m []*model.ProductFridgeBoxListQuery, total int64, err error) {
	q, _ := rq.QueryReadOnly(new(model.BoxItem))
	o := orm.NewOrm()
	o.Using("read_only")

	var condition string

	if statusFilter == "2" {
		condition = " and bf.status=2"
	}
	if statusFilter == "3" {
		condition = " and bf.status=4"
	}
	if statusFilter == "" {
		condition = ""
	}

	var mx []*model.BoxItem
	//	var tempData *model.ProductFridgeBoxListQuery
	var data []*model.ProductFridgeBoxListQuery
	if _, err = q.Exclude("status", 2).All(&mx, rq.Fields...); err == nil {
		for _, boxItem := range mx {

			var tempData *model.ProductFridgeBoxListQuery
			boxItem.Box.Read()
			boxItem.Product.Read()
			boxItem.Product.Uom.Read()
			if e := o.Raw("select bi.id,pi.image_url as item_image, "+
				"bf.last_seen_at as processed_at,bf.warehouse_id , "+
				"bf.image_url as waste_image,bf.status as box_fridge_status    "+
				"from box_item bi  "+
				"left join box_fridge bf on bi.id =bf.box_item_id   "+
				"left join product_image pi on bi.product_id =pi.product_id   "+
				"where bi.id = ? "+condition, boxItem.ID).QueryRow(&tempData); e != nil {
				if condition == "" {
					return nil, total, e
				} else {
					continue
				}
			}

			if tempData.WarehouseId != 0 {
				Warehouse := &model.Warehouse{ID: tempData.WarehouseId}
				if e := Warehouse.Read("ID"); e != nil {
					//o.Failure("warehouse.id", e.Error())
					return nil, total, e
				}

				tempData.WarehouseName = Warehouse.Name

				branchFridge := &model.BranchFridge{Warehouse: Warehouse}
				if e := branchFridge.Read("Warehouse"); e != nil {
					return nil, total, e
				}
				if branchFridge.ID != 0 {
					if e := branchFridge.Branch.Read("ID"); e != nil {
						return nil, total, e
					}
					tempData.BranchName = branchFridge.Branch.Name
				}
				if tempData.BoxFridgeStatus == 1 {
					tempData.Status = "active"
				}
			}
			if int64(boxItem.Status) == 3 {
				tempData.Status = "finished"
			} else {
				if int64(tempData.BoxFridgeStatus) == 1 {
					tempData.Status = "active"
				} else if int64(tempData.BoxFridgeStatus) == 2 {
					tempData.Status = "sold"
				} else if int64(tempData.BoxFridgeStatus) == 4 {
					tempData.Status = "waste"
				}
				if tempData.WarehouseId == 0 {
					tempData.Status = "new"
				}
			}
			tempData.BoxItemStatus = int64(boxItem.Status)
			tempData.FinishedAt = boxItem.FinishedAt
			tempData.ProductName = boxItem.Product.Name
			tempData.Rfid = boxItem.Box.Rfid
			tempData.TotalWeight = boxItem.TotalWeight
			tempData.Uom = boxItem.Product.Uom.Name

			data = append(data, tempData)
		}
		total = int64(len(data))
		return data, total, nil
	}

	return nil, total, err
}

// GetFilterBoxItems : function to get data from database based on parameters with filtered permission
func GetFilterBoxItems(rq *orm.RequestQuery) (m []*model.BoxItem, total int64, err error) {
	q, _ := rq.QueryReadOnly(new(model.BoxItem))

	if total, err = q.Exclude("status", 3).Count(); err != nil || total == 0 {
		return nil, total, err
	}

	var mx []*model.BoxItem
	if _, err = q.Exclude("status", 3).All(&mx, rq.Fields...); err == nil {
		return mx, total, nil
	}

	return nil, total, err
}

// ValidBoxItem : function to check if id is valid in database
func ValidBoxItem(id int64) (BoxItem *model.BoxItem, e error) {
	BoxItem = &model.BoxItem{ID: id}
	e = BoxItem.Read("ID")

	return
}
