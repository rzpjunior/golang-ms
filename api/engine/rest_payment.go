// Copyright 2020 PT. Eden Pangan Indonesia. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

package engine

import (
	payment_channel "git.edenfarm.id/project-version2/api/src/payment/channel"
	"git.edenfarm.id/project-version2/api/src/payment/payment_group"
	"git.edenfarm.id/project-version2/api/src/payment/payment_group_comb"
)

func init() {
	handlers["payment/payment_group"] = &payment_group.Handler{}
	handlers["payment/payment_group_comb"] = &payment_group_comb.Handler{}
	handlers["payment/channel"] = &payment_channel.Handler{}
}
