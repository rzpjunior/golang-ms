// Copyright 2020 PT. Eden Pangan Indonesia. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

package user

import (
	"git.edenfarm.id/cuxs/orm"
)

func ValidUsername(email string, excludeID int64) bool {
	var total int64
	orSelect := orm.NewOrm()
	orSelect.Using("read_only")

	orSelect.Raw("SELECT count(*) FROM user where email = ? and id != ? and is_active = ?", email, excludeID, 1).QueryRow(&total)

	return total == 0
}

func ValidArchive() (tot int64) {
	var total int64
	orSelect := orm.NewOrm()
	orSelect.Using("read_only")

	orSelect.Raw("SELECT count(usergroup_id) FROM user where usergroup_id = ? AND is_active = ?", 1, 1).QueryRow(&total)

	return total
}
