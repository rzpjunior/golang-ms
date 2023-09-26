// Copyright 2020 PT. Eden Pangan Indonesia. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

package engine

import (
	"git.edenfarm.id/project-version2/api/src/finance/debit_note"
	invoice_term "git.edenfarm.id/project-version2/api/src/finance/invoice/term"
	"git.edenfarm.id/project-version2/api/src/finance/payment/method"
	purchase_invoice "git.edenfarm.id/project-version2/api/src/finance/purchase/invoice"
	"git.edenfarm.id/project-version2/api/src/finance/purchase/payment"
	purchase_term "git.edenfarm.id/project-version2/api/src/finance/purchase/term"
	"git.edenfarm.id/project-version2/api/src/finance/sales/invoice"
	sales_term "git.edenfarm.id/project-version2/api/src/finance/sales/term"
)

func init() {
	handlers["finance/payment/method"] = &method.Handler{}
	handlers["finance/purchase/term"] = &purchase_term.Handler{}
	handlers["finance/sales/term"] = &sales_term.Handler{}
	handlers["finance/purchase/invoice"] = &purchase_invoice.Handler{}
	handlers["finance/purchase/payment"] = &purchase_payment.Handler{}
	handlers["finance/invoice/term"] = &invoice_term.Handler{}
	handlers["finance/sales/invoice"] = &invoice.Handler{}
	handlers["finance/debit_note"] = &debit_note.Handler{}
}
