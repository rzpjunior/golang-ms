// Copyright 2022 PT. Eden Pangan Indonesia. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

package engine

import "git.edenfarm.id/project-version2/api/src/warehouse/transfer_sku"


func init() {
	handlers["warehouse/transfer_sku"] = &transfer_sku.Handler{}
}
