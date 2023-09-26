// Copyright 2020 PT. Eden Pangan Indonesia. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.
package engine

import (
	"git.edenfarm.id/project-version2/api/src/sales/order"
	"git.edenfarm.id/project-version2/api/src/sales/order_edn"
	"git.edenfarm.id/project-version2/api/src/sales/payment"
	paymentedn "git.edenfarm.id/project-version2/api/src/sales/payment-edn"
	"git.edenfarm.id/project-version2/api/src/sales/sales_assignment"
	"git.edenfarm.id/project-version2/api/src/sales/sales_assignment/objective"
	"git.edenfarm.id/project-version2/api/src/sales/sales_group"
	"git.edenfarm.id/project-version2/api/src/sales/sales_person"
)

func init() {
	handlers["sales/person"] = &sales_person.Handler{}
	handlers["sales/order"] = &order.Handler{}
	handlers["sales/payment"] = &payment.Handler{}
	handlers["sales/group"] = &sales_group.Handler{}
	handlers["sales/assignment"] = &sales_assignment.Handler{}
	handlers["sales/assignment/objective"] = &objective.Handler{}
	handlers["sales/order-edn"] = &order_edn.Handler{}
	handlers["sales/payment-edn"] = &paymentedn.Handler{}

}
