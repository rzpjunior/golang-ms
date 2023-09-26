// Copyright 2020 PT. Eden Pangan Indonesia. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

package invoice

import (
	"git.edenfarm.id/cuxs/orm"
	"git.edenfarm.id/cuxs/validation"
	"strconv"
	"strings"

	"git.edenfarm.id/project-version2/datamodel/model"
	"git.edenfarm.id/project-version2/api/datastore/repository"
	"git.edenfarm.id/project-version2/api/src/auth"
	"git.edenfarm.id/project-version2/api/util"
)

// cancelRequest : struct to hold Cancel Purchase Invoice request data
type cancelRequest struct {
	ID               int64   `json:"-"`
	CancellationNote string  `json:"note" valid:"required"`
	TotalCharge      float64 `json:"-"`

	PurchaseInvoice *model.PurchaseInvoice `json:"-"`

	Session *auth.SessionData `json:"-"`
}

// Validate : function to validate cancel purchase invoice request data
func (r *cancelRequest) Validate() *validation.Output {
	o := &validation.Output{Valid: true}
	var count int64
	var err error

	orSelect := orm.NewOrm()
	orSelect.Using("read_only")

	r.PurchaseInvoice = &model.PurchaseInvoice{ID: r.ID}
	if err = r.PurchaseInvoice.Read("ID"); err != nil {
		o.Failure("id.invalid", util.ErrorInvalidData("purchase invoice"))
	}

	if r.PurchaseInvoice.Status != 1 {
		o.Failure("status.inactive", util.ErrorDocStatus("purchase invoice", "active"))
		return o
	}

	filter := map[string]interface{}{"status__in": []int{1, 2}, "purchase_invoice_id": r.PurchaseInvoice.ID}
	exclude := map[string]interface{}{}
	_, count, err = repository.GetDataPurchasePayment(filter, exclude)

	if err != nil || count >= 1 {
		o.Failure("status.inactive", util.ErrorDocStatus("purchase payment", "inactive"))
		return o
	}

	err = r.PurchaseInvoice.PurchaseOrder.Read("ID")
	if err != nil {
		o.Failure("pi.purchase_order.invalid", util.ErrorInvalidData("purchase order"))
	}

	if len(r.CancellationNote) > 250 {
		o.Failure("note", util.ErrorCharLength("note", 250))
	}

	debitNoteIds := strings.Split(r.PurchaseInvoice.DebitNoteIDs, ",")
	var debitNoteTotalPrice float64
	for _, v := range debitNoteIds {
		ID, _ := strconv.Atoi(v)
		dn := &model.DebitNote{ID: int64(ID)}
		dn.Read("ID")
		debitNoteTotalPrice += dn.TotalPrice
	}
	r.TotalCharge = r.PurchaseInvoice.TotalCharge + debitNoteTotalPrice

	return o
}

// Messages : function to return error validation messages
func (c *cancelRequest) Messages() map[string]string {
	messages := map[string]string{
		"note.required": util.ErrorInputRequired("note"),
	}

	return messages
}
