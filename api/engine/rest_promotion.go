// Copyright 2020 PT. Eden Pangan Indonesia. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.
package engine

import (
	"git.edenfarm.id/project-version2/api/src/promotion/sku_discount"
	"git.edenfarm.id/project-version2/api/src/promotion/voucher"
)

func init() {
	handlers["promotion/voucher"] = &voucher.Handler{}
	handlers["promotion/sku_discount"] = &sku_discount.Handler{}
}
