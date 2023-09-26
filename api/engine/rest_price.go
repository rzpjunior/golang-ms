// Copyright 2020 PT. Eden Pangan Indonesia. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.
package engine

import (
	"git.edenfarm.id/project-version2/api/src/price/price_set"
	"git.edenfarm.id/project-version2/api/src/price/product_price"
	"git.edenfarm.id/project-version2/api/src/price/schedule"
)

func init() {
	handlers["price/set"] = &price_set.Handler{}
	handlers["price/product_price"] = &product_price.Handler{}
	handlers["price/schedule"] = &schedule.Handler{}
}
