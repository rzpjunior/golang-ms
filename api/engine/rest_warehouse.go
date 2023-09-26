// Copyright 2020 PT. Eden Pangan Indonesia. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

package engine

import (
	"git.edenfarm.id/project-version2/api/src/warehouse/bin"
	delivery_order "git.edenfarm.id/project-version2/api/src/warehouse/delivery/order"
	delivery_return "git.edenfarm.id/project-version2/api/src/warehouse/delivery/return"
	goods_receipt "git.edenfarm.id/project-version2/api/src/warehouse/goods/receipt"
	"git.edenfarm.id/project-version2/api/src/warehouse/goods/receipt/transfer_sku"
	goods_transfer "git.edenfarm.id/project-version2/api/src/warehouse/goods/transfer"
	"git.edenfarm.id/project-version2/api/src/warehouse/packing"
	"git.edenfarm.id/project-version2/api/src/warehouse/picking"
	"git.edenfarm.id/project-version2/api/src/warehouse/stock"
	"git.edenfarm.id/project-version2/api/src/warehouse/stock_opname"
	"git.edenfarm.id/project-version2/api/src/warehouse/supplier_return"
	waste_disposal "git.edenfarm.id/project-version2/api/src/warehouse/waste/disposal"
	waste_entry "git.edenfarm.id/project-version2/api/src/warehouse/waste/entry"
)

func init() {
	handlers["warehouse/stock"] = &stock.Handler{}
	handlers["warehouse/stock_opname"] = &stock_opname.Handler{}
	handlers["warehouse/delivery_order"] = &delivery_order.Handler{}
	handlers["warehouse/goods/receipt"] = &goods_receipt.Handler{}
	handlers["warehouse/goods/transfer"] = &goods_transfer.Handler{}
	handlers["warehouse/waste/entry"] = &waste_entry.Handler{}
	handlers["warehouse/waste/disposal"] = &waste_disposal.Handler{}
	handlers["warehouse/delivery_return"] = &delivery_return.Handler{}
	handlers["warehouse/packing_order"] = &packing.Handler{}
	handlers["warehouse/picking_order"] = &picking.Handler{}
	handlers["warehouse/supplier_return"] = &supplier_return.Handler{}
	handlers["warehouse/goods/receipt/:id/transfer_sku"] = &transfer_sku.Handler{}
	handlers["warehouse/bin"] = &bin.Handler{}
}
