// Copyright 2020 PT. Eden Pangan Indonesia. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

package engine

import (
	"git.edenfarm.id/project-version2/api/src/customer/agent"
	"git.edenfarm.id/project-version2/api/src/customer/archetype"
	"git.edenfarm.id/project-version2/api/src/customer/branch"
	"git.edenfarm.id/project-version2/api/src/customer/business_type"
	"git.edenfarm.id/project-version2/api/src/customer/customer_acquisition"
	"git.edenfarm.id/project-version2/api/src/customer/distribution_network"
	"git.edenfarm.id/project-version2/api/src/customer/merchant"
	"git.edenfarm.id/project-version2/api/src/customer/prospect_customer"
	"git.edenfarm.id/project-version2/api/src/customer/tag"
)

func init() {
	handlers["customer/tag"] = &tag.Handler{}
	handlers["customer/archetype"] = &archetype.Handler{}
	handlers["customer/business_type"] = &business_type.Handler{}
	handlers["customer/merchant"] = &merchant.Handler{}
	handlers["customer/branch"] = &branch.Handler{}
	handlers["customer/agent"] = &agent.Handler{}
	handlers["customer/prospect_customer"] = &prospect_customer.Handler{}
	handlers["customer/acquisition"] = &customer_acquisition.Handler{}
	handlers["customer/distribution_network"] = &distribution_network.Handler{}
}
