// Copyright 2020 PT. Eden Pangan Indonesia. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

package engine

import (
	"git.edenfarm.id/project-version2/api/src/inventory/category"
	"git.edenfarm.id/project-version2/api/src/inventory/product"
	"git.edenfarm.id/project-version2/api/src/inventory/tag_product"
	"git.edenfarm.id/project-version2/api/src/inventory/uom"
)

func init() {
	handlers["inventory/uom"] = &uom.Handler{}
	handlers["inventory/category"] = &category.Handler{}
	handlers["inventory/product"] = &product.Handler{}
	handlers["inventory/tag_product"] = &tag_product.Handler{}
}
