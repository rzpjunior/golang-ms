package warehouse_coverage

import (
	"git.edenfarm.id/cuxs/orm"
	"git.edenfarm.id/project-version2/api/log"
	"git.edenfarm.id/project-version2/datamodel/model"
)

// Save : function to save data requested into database
func Save(r createRequest) (warehouseCoverage *model.WarehouseCoverage, e error) {
	o := orm.NewOrm()
	o.Begin()

	warehouseCoverage = &model.WarehouseCoverage{
		MainWarehouse: r.MainWarehouseInt,
		Warehouse:     r.Warehouse,
		SubDistrict:   r.SubDistrict,
	}
	if r.WarehouseType.ValueName == "HUB" {
		warehouseCoverage.ParentWarehouseID = r.Warehouse.ParentID
	}

	_, e = o.Insert(warehouseCoverage)
	if e != nil {
		o.Rollback()
		return nil, e
	}

	e = log.AuditLogByUser(r.Session.Staff, warehouseCoverage.ID, "warehouse_coverage", "create", "")
	if e != nil {
		o.Rollback()
		return nil, e
	}

	if r.ChangeMainWarehouse {
		r.AffectedWarehouseCoverage.MainWarehouse = 2

		if _, err := o.Update(r.AffectedWarehouseCoverage, "MainWarehouse"); err != nil {
			o.Rollback()
			return nil, err
		}

		e = log.AuditLogByUser(r.Session.Staff, r.AffectedWarehouseCoverage.ID, "warehouse_coverage", "update", "")
		if e != nil {
			o.Rollback()
			return nil, e
		}
	}

	o.Commit()
	return warehouseCoverage, e
}

// Delete : function to Delete data requested into database
func Delete(r deleteRequest) (warehouseCoverage *model.WarehouseCoverage, e error) {
	o := orm.NewOrm()
	o.Begin()

	_, e = o.Delete(r.WarehouseCoverage)
	if e != nil {
		o.Rollback()
		return nil, e
	}

	e = log.AuditLogByUser(r.Session.Staff, r.WarehouseCoverage.ID, "warehouse_coverage", "delete", "")
	if e != nil {
		o.Rollback()
		return nil, e
	}

	o.Commit()
	return r.WarehouseCoverage, e
}

// UpdateMain : function to update data requested into database
func UpdateMain(r updateMainRequest) (warehouseCoverage *model.WarehouseCoverage, e error) {
	o := orm.NewOrm()
	o.Begin()

	r.WarehouseCoverage.MainWarehouse = 1

	if _, err := o.Update(r.WarehouseCoverage, "MainWarehouse"); err != nil {
		o.Rollback()
		return nil, err
	}

	e = log.AuditLogByUser(r.Session.Staff, r.WarehouseCoverage.ID, "warehouse_coverage", "update", "")
	if e != nil {
		o.Rollback()
		return nil, e
	}

	if r.ChangeMainWarehouse {
		r.AffectedWarehouseCoverage.MainWarehouse = 2

		if _, err := o.Update(r.AffectedWarehouseCoverage, "MainWarehouse"); err != nil {
			o.Rollback()
			return nil, err
		}

		e = log.AuditLogByUser(r.Session.Staff, r.AffectedWarehouseCoverage.ID, "warehouse_coverage", "update", "")
		if e != nil {
			o.Rollback()
			return nil, e
		}
	}

	o.Commit()
	return r.WarehouseCoverage, e
}
