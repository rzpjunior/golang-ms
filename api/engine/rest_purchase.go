// Copyright 2020 PT. Eden Pangan Indonesia. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

package engine

import (
	"git.edenfarm.id/project-version2/api/src/purchase/cogs"
	"git.edenfarm.id/project-version2/api/src/purchase/forecast_demand"
	"git.edenfarm.id/project-version2/api/src/purchase/order"
	"git.edenfarm.id/project-version2/api/src/purchase/order/consolidated_shipment"
	"git.edenfarm.id/project-version2/api/src/purchase/order/field_purchaser"
	"git.edenfarm.id/project-version2/api/src/purchase/order/purchase_deliver"
	"git.edenfarm.id/project-version2/api/src/purchase/plan"
	"git.edenfarm.id/project-version2/api/src/purchase/prospect_supplier"
	"git.edenfarm.id/project-version2/api/src/purchase/sales_recap"
	"git.edenfarm.id/project-version2/api/src/purchase/stall"
	"git.edenfarm.id/project-version2/api/src/purchase/supplier"
	"git.edenfarm.id/project-version2/api/src/purchase/supplier_badge"
	"git.edenfarm.id/project-version2/api/src/purchase/supplier_commodity"
	"git.edenfarm.id/project-version2/api/src/purchase/supplier_group"
	"git.edenfarm.id/project-version2/api/src/purchase/supplier_organization"
	"git.edenfarm.id/project-version2/api/src/purchase/supplier_type"
)

func init() {
	handlers["purchase/order"] = &order.Handler{}
	handlers["purchase/supplier"] = &supplier.Handler{}
	handlers["purchase/supplier_type"] = &supplier_type.Handler{}
	handlers["purchase/prospect/supplier"] = &prospect_supplier.Handler{}
	handlers["purchase/supplier/badge"] = &supplier_badge.Handler{}
	handlers["purchase/supplier/group"] = &supplier_group.Handler{}
	handlers["purchase/supplier/commodity"] = &supplier_commodity.Handler{}
	handlers["purchase/forecast_demand"] = &forecast_demand.Handler{}
	handlers["purchase/cogs"] = &cogs.Handler{}
	handlers["purchase/sales_recap"] = &sales_recap.Handler{}
	handlers["purchase/stall"] = &stall.Handler{}
	handlers["purchase/supplier/organization"] = &supplier_organization.Handler{}
	handlers["purchase/plan"] = &plan.Handler{}
	handlers["purchase/order/field_purchaser"] = &field_purchaser.Handler{}
	handlers["purchase/order/purchase_deliver"] = &purchase_deliver.Handler{}
	handlers["purchase/order/consolidated_shipment"] = &consolidated_shipment.Handler{}
}
