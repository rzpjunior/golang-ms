package supplier_badge

import (
	"time"

	"git.edenfarm.id/cuxs/orm"
	"git.edenfarm.id/project-version2/datamodel/model"
	"git.edenfarm.id/project-version2/api/log"
	"git.edenfarm.id/project-version2/api/util"
)

// Save : function to save data requested into database
func Save(r createRequest) (supplierBadge *model.SupplierBadge, err error) {
	r.Code, _ = util.GenerateCode(r.Code, "supplier_badge")
	o := orm.NewOrm()

	o.Begin()

	supplierBadge = &model.SupplierBadge{
		Name:      r.Name,
		Code:      r.Code,
		Note:      r.Note,
		Status:    1,
		UpdatedBy: r.Session.Staff,
		CreatedAt: time.Now(),
	}

	_, err = o.Insert(supplierBadge)

	if err != nil {
		o.Rollback()
		return nil, err
	}

	err = log.AuditLogByUser(r.Session.Staff, supplierBadge.ID, "supplier_badge", "create", "")

	if err != nil {
		o.Rollback()
		return nil, err
	}

	o.Commit()

	return supplierBadge, err
}

func Update(u updateRequest) (supplierBadge *model.SupplierBadge, err error) {
	o := orm.NewOrm()
	o.Begin()

	u.SupplierBadge.Name = u.Name
	u.SupplierBadge.Note = u.Note
	u.SupplierBadge.UpdatedAt = time.Now()
	u.SupplierBadge.UpdatedBy = u.Session.Staff

	_, err = o.Update(u.SupplierBadge, "Name", "Note", "UpdatedAt", "UpdatedBy")

	if err != nil {
		o.Rollback()
		return nil, err
	}

	err = log.AuditLogByUser(u.Session.Staff, u.SupplierBadge.ID, "supplier_badge", "update", "")

	if err != nil {
		o.Rollback()
		return nil, err
	}

	o.Commit()

	return u.SupplierBadge, err
}
