// Copyright 2020 PT. Eden Pangan Indonesia. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

package engine

import (
	"git.edenfarm.id/project-version2/api/src/dashboard"
	"git.edenfarm.id/project-version2/api/src/dashboard/field_purchaser"
	"git.edenfarm.id/project-version2/api/src/dashboard/fulfillment"
	"git.edenfarm.id/project-version2/api/src/dashboard/widget"
)

func init() {
	handlers["dashboard"] = &dashboard.Handler{}
	handlers["dashboard/fulfillment"] = &fulfillment.Handler{}
	handlers["dashboard/operation"] = &widget.Handler{}
	handlers["dashboard/field_purchaser"] = &field_purchaser.Handler{}
}
