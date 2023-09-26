// Copyright 2020 PT. Eden Pangan Indonesia. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

package engine

import (
	"git.edenfarm.id/project-version2/api/src/campaign/banner"
	"git.edenfarm.id/project-version2/api/src/campaign/edenpoint"
	log "git.edenfarm.id/project-version2/api/src/campaign/edenpoint/eplog"
	"git.edenfarm.id/project-version2/api/src/campaign/membership"
	"git.edenfarm.id/project-version2/api/src/campaign/product_section"
	"git.edenfarm.id/project-version2/api/src/campaign/push_notification"
)

func init() {
	handlers["campaign/push-notification"] = &push_notification.Handler{}
	handlers["campaign/banner"] = &banner.Handler{}
	handlers["campaign/eden-point-log"] = &log.Handler{}
	handlers["campaign/eden-point"] = &edenpoint.Handler{}
	handlers["campaign/product-section"] = &product_section.Handler{}
	handlers["campaign/membership"] = &membership.Handler{}
}
