// Copyright 2020 PT. Eden Pangan Indonesia. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

package engine

import (
	"git.edenfarm.id/project-version2/api/src/user/role"
	"git.edenfarm.id/project-version2/api/src/user/user"
)

func init() {
	handlers["user"] = &user.Handler{}
	handlers["role"] = &role.Handler{}
}
