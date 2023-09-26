package supplier_return

import (
	"git.edenfarm.id/cuxs/orm"
	"git.edenfarm.id/cuxs/validation"
	"git.edenfarm.id/project-version2/datamodel/model"
	"git.edenfarm.id/project-version2/api/datastore/repository"
	"git.edenfarm.id/project-version2/api/src/auth"
	"git.edenfarm.id/project-version2/api/util"
)

// cancelRequest : struct to hold supplier return request data
type cancelRequest struct {
	ID int64 `json:"-"`

	CancellationNote string `json:"cancellation_note"`

	SupplierReturn *model.SupplierReturn `json:"-"`
	DebitNote      *model.DebitNote      `json:"-"`
	Session        *auth.SessionData     `json:"-"`
}

// Validate : function to validate uom request data
func (r *cancelRequest) Validate() *validation.Output {
	o := &validation.Output{Valid: true}
	q := orm.NewOrm()
	q.Using("read_only")
	var e error

	// region supplier return definition
	if r.SupplierReturn, e = repository.ValidSupplierReturn(r.ID); e != nil {
		o.Failure("id.invalid", util.ErrorInvalidData("supplier return"))
	}
	// endregion

	// region debit note definition
	r.DebitNote = new(model.DebitNote)
	if e = q.QueryTable(new(model.DebitNote)).Filter("supplier_return_id", r.SupplierReturn.ID).One(r.DebitNote); e != nil {
		o.Failure("debit_note_id.invalid", util.ErrorInvalidData("debit note"))
		return o
	}
	// endregion

	if r.SupplierReturn.Status != 1 {
		o.Failure("supplier_return_id.invalid", util.ErrorActive("supplier return"))
		return o
	}

	if r.DebitNote.Status != 1 {
		o.Failure("debit_note_id.invalid", util.ErrorActive("debit note"))
		return o
	}

	if r.DebitNote.UsedInPurchaseInvoice == 1 {
		o.Failure("debit_note_id.invalid", util.ErrorIsBeingUsed("debit note"))
		return o
	}
	r.SupplierReturn.GoodsReceipt.Read("ID")

	return o
}

// Messages : function to return error validation messages
func (r *cancelRequest) Messages() map[string]string {
	messages := map[string]string{}

	return messages
}
