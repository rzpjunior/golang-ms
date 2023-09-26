// Copyright 2020 PT. Eden Pangan Indonesia. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

package engine

import (
	"git.edenfarm.id/project-version2/api/src/configuration/application"
	"git.edenfarm.id/project-version2/api/src/configuration/area"
	businessPolicy "git.edenfarm.id/project-version2/api/src/configuration/area/business_policy"
	"git.edenfarm.id/project-version2/api/src/configuration/area/policy"
	dayOff "git.edenfarm.id/project-version2/api/src/configuration/day_off"
	"git.edenfarm.id/project-version2/api/src/configuration/division"
	glossary "git.edenfarm.id/project-version2/api/src/configuration/glossary"
	orderType "git.edenfarm.id/project-version2/api/src/configuration/order_type"
	"git.edenfarm.id/project-version2/api/src/configuration/user/profile"
	"git.edenfarm.id/project-version2/api/src/configuration/warehouse"
	warehouse_coverage "git.edenfarm.id/project-version2/api/src/configuration/warehouse_coverage"
	"git.edenfarm.id/project-version2/api/src/configuration/wrt"
)

func init() {
	handlers["config/area"] = &area.Handler{}
	handlers["config/division"] = &division.Handler{}
	handlers["config/wrt"] = &wrt.Handler{}
	handlers["config/warehouse"] = &warehouse.Handler{}
	handlers["config/area/policy"] = &policy.Handler{}
	handlers["config/app"] = &application.Handler{}
	handlers["config/user/profile"] = &profile.Handler{}
	handlers["config/glossary"] = &glossary.Handler{}
	handlers["config/warehouse/coverage"] = &warehouse_coverage.Handler{}
	handlers["config/order_type"] = &orderType.Handler{}
	handlers["config/day_off"] = &dayOff.Handler{}
	handlers["config/area/business_policy"] = &businessPolicy.Handler{}
}
