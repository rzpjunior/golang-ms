package supplier_organization

import (
	"time"

	"git.edenfarm.id/cuxs/orm"
	"git.edenfarm.id/project-version2/api/log"
	"git.edenfarm.id/project-version2/api/util"
	"git.edenfarm.id/project-version2/datamodel/model"
)

// Save : function to save data requested into database
func Save(r createRequest) (supplierOrganization *model.SupplierOrganization, err error) {
	r.Code, _ = util.GenerateCode(r.Code, "supplier_organization")
	o := orm.NewOrm()

	o.Begin()

	supplierOrganization = &model.SupplierOrganization{
		SupplierCommodity: r.SupplierCommodity,
		SupplierBadge:     r.SupplierBadge,
		SupplierType:      r.SupplierType,
		TermPaymentPur:    r.PurchaseTerm,
		SubDistrict:       r.SubDistrict,
		Name:              r.Name,
		Address:           r.Address,
		Code:              r.Code,
		Note:              r.Note,
		Status:            1,
		CreatedBy:         r.Session.Staff,
		CreatedAt:         time.Now(),
	}

	_, err = o.Insert(supplierOrganization)

	if err != nil {
		o.Rollback()
		return nil, err
	}

	err = log.AuditLogByUser(r.Session.Staff, supplierOrganization.ID, "supplier_organization", "create", "")

	if err != nil {
		o.Rollback()
		return nil, err
	}

	o.Commit()

	return supplierOrganization, err
}

func Update(u updateRequest) (supplierOrganization *model.SupplierOrganization, err error) {
	o := orm.NewOrm()
	o.Begin()

	supplierOrganization = &model.SupplierOrganization{
		ID:        u.ID,
		Name:      u.Name,
		Address:   u.Address,
		Note:      u.Note,
		UpdatedAt: time.Now(),
		UpdatedBy: u.Session.Staff,
	}

	err = supplierOrganization.Save("Address", "Name", "Note", "UpdatedAt", "UpdatedBy")

	if err != nil {
		o.Rollback()
		return nil, err
	}

	err = log.AuditLogByUser(u.Session.Staff, u.SupplierOrganization.ID, "supplier_organization", "update", "")

	if err != nil {
		o.Rollback()
		return nil, err
	}

	o.Commit()

	return supplierOrganization, err
}
