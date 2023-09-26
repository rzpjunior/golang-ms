// Copyright 2020 PT. Eden Pangan Indonesia. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

package engine

import (
	upload "git.edenfarm.id/project-version2/api/src/upload_file"
)

func init() {
	handlers["upload"] = &upload.Handler{}
}
