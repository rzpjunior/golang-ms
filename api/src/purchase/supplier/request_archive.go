package supplier

import (
	"git.edenfarm.id/cuxs/orm"
	"git.edenfarm.id/cuxs/validation"
	"git.edenfarm.id/project-version2/datamodel/model"
	"git.edenfarm.id/project-version2/api/datastore/repository"
	"git.edenfarm.id/project-version2/api/src/auth"
	"git.edenfarm.id/project-version2/api/util"
)

type archiveRequest struct {
	ID      int64             `json:"-" valid:"required"`
	Session *auth.SessionData `json:"-"`

	PurchaseOrder *model.PurchaseOrder
}

func (c *archiveRequest) Validate() *validation.Output {
	o := &validation.Output{Valid: true}
	orSelect := orm.NewOrm()
	orSelect.Using("read_only")

	if supplier, err := repository.ValidSupplier(c.ID); err == nil {
		if supplier.Status != 1 {
			o.Failure("id.active", util.ErrorActive("status"))
		}
	} else {
		o.Failure("id.invalid", util.ErrorInvalidData("supplier"))
	}

	var countPO int8
	orSelect.Raw("select count(*) from purchase_order po where po.supplier_id = ? and po.status in (1,5)", c.ID).QueryRow(&countPO)
	if countPO > 0 {
		o.Failure("id.invalid", util.ErrorRelated("active/draft", "purchase order", "supplier"))
	}

	return o
}

func (c *archiveRequest) Messages() map[string]string {
	return map[string]string{}
}
