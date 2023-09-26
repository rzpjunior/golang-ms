// Copyright 2020 PT. Eden Pangan Indonesia. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

package engine

import (
	"git.edenfarm.id/project-version2/api/src/koli"
	"git.edenfarm.id/project-version2/api/src/koli/increment"
)

func init() {
	handlers["koli"] = &koli.Handler{}
	handlers["delivery/koli/increment"] = &koli_increment.Handler{}
}
