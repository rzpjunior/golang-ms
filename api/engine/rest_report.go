// Copyright 2020 PT. Eden Pangan Indonesia. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.
package engine

import (
	"git.edenfarm.id/project-version2/api/src/report"
	cmsReport "git.edenfarm.id/project-version2/api/src/report/cms"
	fridgeReport "git.edenfarm.id/project-version2/api/src/report/fridge"
	rmsReport "git.edenfarm.id/project-version2/api/src/report/rms"
	smsReport "git.edenfarm.id/project-version2/api/src/report/sms"
	wmsReport "git.edenfarm.id/project-version2/api/src/report/wms"
)

func init() {
	handlers["report"] = &report.Handler{}
	handlers["report/sms"] = &smsReport.Handler{}
	handlers["report/cms"] = &cmsReport.Handler{}
	handlers["report/wms"] = &wmsReport.Handler{}
	handlers["report/rms"] = &rmsReport.Handler{}
	handlers["report/fridge"] = &fridgeReport.Handler{}
}
