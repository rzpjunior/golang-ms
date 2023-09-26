// Copyright 2021 PT. Eden Pangan Indonesia. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

package repository

import (
	"git.edenfarm.id/project-version2/datamodel/model"

	"errors"
	"git.edenfarm.id/cuxs/common"
)

var (
	errorInvalidItem = "Invalid item"
)

// ValidSalesInvoiceItem : function to check if id is valid in database
func ValidSalesInvoiceItem(PurchaseInvoiceID string) (i *model.SalesInvoiceItem, e error) {
	i = new(model.SalesInvoiceItem)

	if i.ID, e = common.Decrypt(PurchaseInvoiceID); e == nil {
		if e = i.Read("ID"); e != nil {
			// saat dikirim ke return variable i tidak boleh memiliki nilai
			i = nil
			e = errors.New(errorInvalidItem)
		}
	} else {
		// saat dikirim ke return variable i tidak boleh memiliki nilai
		i = nil
		e = errors.New(errorInvalidItem)
	}

	return
}
