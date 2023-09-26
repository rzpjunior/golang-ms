// Copyright 2020 PT. Eden Pangan Indonesia. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

package repository

import (
	"git.edenfarm.id/project-version2/datamodel/model"
	"strings"

	"git.edenfarm.id/cuxs/orm"
)

// GetStaff find a single data staff using field and value condition.
func GetStaff(field string, values ...interface{}) (*model.Staff, error) {
	m := new(model.Staff)
	o := orm.NewOrm()
	o.Using("read_only")
	var warehouseAccessArr []string
	var qMark string

	if err := o.QueryTable(m).Filter(field, values...).RelatedSel().Limit(1).One(m); err != nil {
		return nil, err
	}

	o.Raw("SELECT name from sales_group where `id` = ? ", m.SalesGroupID).QueryRow(&m.SalesGroupName)

	if _, err := o.QueryTable(new(model.Warehouse)).RelatedSel("area").Filter("id__in", strings.Split(m.WarehouseAccessStr, ",")).All(&m.WarehouseAccess); err != nil {
		return nil, err
	}

	warehouseAccessArr = strings.Split(m.WarehouseAccessStr, ",")
	if m.WarehouseAccessStr != "" {
		qMark = ""
		for _, _ = range warehouseAccessArr {
			qMark = qMark + "?,"
		}
		qMark = strings.TrimSuffix(qMark, ",")
		o.Raw("select group_concat(name) from warehouse where id in ("+qMark+")", warehouseAccessArr).QueryRow(&m.WarehouseAccessStr)
		m.WarehouseAccessStr = strings.Replace(m.WarehouseAccessStr, ",", ", ", -1)
	}

	return m, nil
}

// GetStaffs get all data staff that matched with query request parameters.
// returning slices of User, total data without limit and error.
func GetStaffs(rq *orm.RequestQuery) (m []*model.Staff, total int64, err error) {
	// make new orm query
	q, _ := rq.QueryReadOnly(new(model.Staff))
	o := orm.NewOrm()
	o.Using("read_only")

	// get total data
	if total, err = q.Exclude("status", 3).Count(); err != nil || total == 0 {
		return nil, total, err
	}

	// get data requested
	var mx []*model.Staff
	if _, err = q.Exclude("status", 3).All(&mx, rq.Fields...); err == nil {
		for _, v := range mx {
			if v.SalesGroupID != 0 {
				if err := o.Raw("select name from sales_group where id = ?", v.SalesGroupID).QueryRow(&v.SalesGroupName); err != nil {
					return nil, total, err
				}
			}
		}
		return mx, total, nil
	}

	// return error some thing went wrong
	return nil, total, err
}

func GetHelpers(rq *orm.RequestQuery) (m []*model.Staff, total int64, err error) {
	// make new orm query
	q, _ := rq.QueryReadOnly(new(model.Staff))

	// get total data
	if total, err = q.Exclude("status", 3).Filter("role_group", 2).Count(); err != nil || total == 0 {
		return nil, total, err
	}

	// get data requested
	var mx []*model.Staff
	if _, err = q.Exclude("status", 3).Filter("role_group", 2).All(&mx, rq.Fields...); err == nil {
		return mx, total, nil
	}

	// return error some thing went wrong
	return nil, total, err
}

// GetStaffs get all data staff that matched with query request parameters.
// returning slices of User, total data without limit and error.
func GetFilterStaff(rq *orm.RequestQuery) (m []*model.Staff, total int64, err error) {
	// make new orm query
	q, _ := rq.QueryReadOnly(new(model.Staff))

	// get total data
	if total, err = q.Exclude("status", 3).Count(); err != nil || total == 0 {
		return nil, total, err
	}

	// get data requested
	var mx []*model.Staff
	if _, err = q.Exclude("status", 3).All(&mx, rq.Fields...); err == nil {
		return mx, total, nil
	}

	// return error some thing went wrong
	return nil, total, err
}

// ValidStaff : check if staff id is valid (exist in database)
func ValidStaff(id int64) (staff *model.Staff, e error) {
	staff = &model.Staff{ID: id}
	e = staff.Read("ID")

	return
}

// GetSupervisor : function to get supervisor from table staff (staff that is a parent)
func GetSupervisor(rq *orm.RequestQuery) (m []*model.Staff, total int64, err error) {
	var exist bool
	var staffTemp []*model.Staff

	q, _ := rq.QueryReadOnly(new(model.Staff))

	total, err = q.RelatedSel("Parent").Filter("parent_id__isnull", false).Filter("parent_id__status", 1).Count()

	if _, err = q.RelatedSel("Parent").Filter("parent_id__isnull", false).Filter("parent_id__status", 1).All(&staffTemp, rq.Fields...); err == nil {
		existID := make(map[int64]int64)
		for _, v := range staffTemp {
			_, exist = existID[v.Parent.ID]
			if !exist {
				existID[v.Parent.ID] = 1
				m = append(m, v.Parent)
			}
		}

		return m, total, nil
	}

	return nil, 0, err
}

// GetSupervisor : function to get supervisor from table staff (staff that is a parent)
func GetSupervisorFilter(rq *orm.RequestQuery) (m []*model.Staff, total int64, err error) {
	var exist bool
	var staffTemp []*model.Staff

	q, _ := rq.QueryReadOnly(new(model.Staff))

	if _, err = q.RelatedSel("Parent").Filter("parent_id__isnull", false).Filter("parent_id__status", 1).All(&staffTemp, rq.Fields...); err == nil {
		existID := make(map[int64]int64)
		for _, v := range staffTemp {
			_, exist = existID[v.Parent.ID]
			if !exist {
				existID[v.Parent.ID] = 1
				m = append(m, v.Parent)
				total++
			}
		}

		return m, total, nil
	}

	return nil, 0, err
}

// CountParentStaff : function to count how many staff that is a child of the requested staff id
func CountParentStaff(staffID int64) (total int64, err error) {
	o := orm.NewOrm()
	o.Using("read_only")

	if countStaff, err := o.QueryTable(new(model.Staff)).Filter("parent_id", staffID).Count(); err == nil {
		return countStaff, nil
	}

	return 0, err
}

// GetFieldPurchaser get all data staff with role field purchaser.
func GetFieldPurchaser(rq *orm.RequestQuery, warehouseID int64) (mx []*model.Staff, total int64, err error) {
	q, _ := rq.QueryReadOnly(new(model.Staff))
	o := orm.NewOrm()
	o.Using("read_only")
	warehouse := new(model.Warehouse)

	o.QueryTable(warehouse).Filter("name", "All Warehouse").Limit(1).One(warehouse)

	if warehouseID == warehouse.ID || warehouseID == 0 {
		if total, err = q.Exclude("status", 3).Filter("Role__Name__in", "Field Purchaser", "Sourcing Admin").Count(); err != nil || total == 0 {
			return nil, total, err
		}

		if _, err = q.Exclude("status", 3).Filter("Role__Name__in", "Field Purchaser", "Sourcing Admin").All(&mx, rq.Fields...); err != nil {
			return nil, total, err
		}
	} else {
		if total, err = q.Exclude("status", 3).Filter("Warehouse__ID__in", warehouseID, warehouse.ID).Filter("Role__Name__in", "Field Purchaser", "Sourcing Admin").Count(); err != nil || total == 0 {
			return nil, total, err
		}

		if _, err = q.Exclude("status", 3).Filter("Warehouse__ID__in", warehouseID, warehouse.ID).Filter("Role__Name__in", "Field Purchaser", "Sourcing Admin").All(&mx, rq.Fields...); err != nil {
			return nil, total, err
		}
	}

	return mx, total, err
}

// CheckStaffData : function to get all staff data based on filter and exclude parameters
func CheckStaffData(filter, exclude map[string]interface{}) (m []*model.Staff, total int64, err error) {
	rq := orm.RequestQuery{}
	o, _ := rq.QueryReadOnly(new(model.Staff))

	for k, v := range filter {
		o = o.Filter(k, v)
	}

	for k, v := range exclude {
		o = o.Exclude(k, v)
	}

	if total, err = o.All(&m); err != nil {
		return nil, 0, err
	}

	return m, total, nil
}
