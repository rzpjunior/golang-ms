package supplier_commodity

import (
	"time"

	"git.edenfarm.id/cuxs/orm"
	"git.edenfarm.id/project-version2/datamodel/model"
	"git.edenfarm.id/project-version2/api/log"
	"git.edenfarm.id/project-version2/api/util"
)

// Save : function to save data requested into database
func Save(r createRequest) (supplierCommodity *model.SupplierCommodity, err error) {
	r.Code, _ = util.GenerateCode(r.Code, "supplier_commodity")
	o := orm.NewOrm()

	o.Begin()

	supplierCommodity = &model.SupplierCommodity{
		Name:      r.Name,
		Code:      r.Code,
		Note:      r.Note,
		Status:    1,
		UpdatedBy: r.Session.Staff,
		CreatedAt: time.Now(),
	}

	_, err = o.Insert(supplierCommodity)

	if err != nil {
		o.Rollback()
		return nil, err
	}

	err = log.AuditLogByUser(r.Session.Staff, supplierCommodity.ID, "supplier_commodity", "create", "")

	if err != nil {
		o.Rollback()
		return nil, err
	}

	o.Commit()

	return supplierCommodity, err
}

func Update(u updateRequest) (supplierCommodity *model.SupplierCommodity, err error) {
	o := orm.NewOrm()
	o.Begin()

	u.SupplierCommodity.Name = u.Name
	u.SupplierCommodity.Note = u.Note
	u.SupplierCommodity.UpdatedAt = time.Now()
	u.SupplierCommodity.UpdatedBy = u.Session.Staff

	_, err = o.Update(u.SupplierCommodity, "Name", "Note", "UpdatedAt", "UpdatedBy")

	if err != nil {
		o.Rollback()
		return nil, err
	}

	err = log.AuditLogByUser(u.Session.Staff, u.SupplierCommodity.ID, "supplier_commodity", "update", "")

	if err != nil {
		o.Rollback()
		return nil, err
	}

	o.Commit()

	return u.SupplierCommodity, err
}
