// Copyright 2020 PT. Eden Pangan Indonesia. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

package role

import (
	"git.edenfarm.id/cuxs/orm"
	"git.edenfarm.id/project-version2/datamodel/model"
	"git.edenfarm.id/project-version2/api/log"
	"git.edenfarm.id/project-version2/api/util"
	"strings"
)

func Save(r createRequest) (u *model.Role, e error) {
	o := orm.NewOrm()
	e = o.Begin()
	r.CodeRole, e = util.GenerateCode(r.CodeRole, "role")

	if e == nil {
		u = &model.Role{
			Code:     r.CodeRole,
			Name:     r.Name,
			Division: r.Division,
			Note:     r.Note,
			Status:   int8(1),
		}
		if _, e = o.Insert(u); e == nil {
			for _, p := range r.Permission {
				pm := &model.RolePermission{
					Role:       u,
					Permission: p,
				}
				if _, e = o.Insert(pm); e != nil {
					e = o.Rollback()
				}
			}
		} else {
			e = o.Rollback()
		}
		e = log.AuditLogByUser(r.Session.Staff, u.ID, "role", "create", "")
	}

	if e == nil {
		e = o.Commit()
		return
	} else {
		e = o.Rollback()
	}
	return nil, e
}

func (u *updateRequest) Update() (role *model.Role, e error) {
	//STEP 1, join old-permission array with new-permission array from the request
	o := orm.NewOrm()
	o.Begin()

	role = &model.Role{
		ID:   u.ID,
		Name: u.Name,
		Note: u.Note,
	}
	if _, e = o.Update(role, "Name", "Note"); e != nil {
		e = o.Rollback()
	}

	OldPermission := u.OldRolePermissionID
	var z []int64
	for _, row := range u.NewPermission {
		u.OldRolePermissionID = append(u.OldRolePermissionID, row)
	}
	//STEP 2, get unique-value array from the joined old-new array from step 1
	differentPermission := util.GetUniqueValue(u.OldRolePermissionID)

	//STEP 3, join old-permission array with unique-value array from step 2
	for _, id := range differentPermission {
		OldPermission = append(OldPermission, id)
	}

	//STEP 4, get same-value array from the joined old-unique array from step 3, this will result will-deleted array
	willDeletedPermission := util.GetSameValue(OldPermission)

	//STEP 5, delete permission, according to will-deleted array
	if len(willDeletedPermission) > 0 {
		var cat []string
		for i := 0; i < len(willDeletedPermission); i++ {
			cat = append(cat, "?")
		}
		catLength := strings.Join(cat, ",")
		if _, e = o.Raw("DELETE FROM role_permission WHERE role_id = ? AND permission_id IN ("+catLength+")", u.ID, willDeletedPermission).Exec(); e != nil {
			e = o.Rollback()
		}
	}

	z = differentPermission

	//STEP 6, join will-deleted array resulted from step 4, with unique-value array from step 2
	for _, ii := range willDeletedPermission {
		z = append(z, ii)
	}

	//STEP 7, get unique-value array in step 6, this will result new-permission array
	newPermission := util.GetUniqueValue(z)

	//STEP 8, insert permission, according to new-permission array
	if len(newPermission) > 0 {
		for _, id := range newPermission {
			// this func for decrypt
			// for get data permission from id
			p := &model.Permission{ID: id}
			if e = p.Read("ID"); e == nil {
				// add permission
				up := &model.RolePermission{
					Permission: p,
					Role:       &model.Role{ID: u.ID},
				}
				if _, e = o.Insert(up); e != nil {
					e = o.Rollback()
				}
			}
		}
	}
	e = log.AuditLogByUser(u.Session.Staff, u.ID, "role", "update", "")

	if e == nil {
		e = o.Commit()
	} else {
		e = o.Rollback()
	}

	return
}

// Archive : function to update status data into archive
func Archive(r archiveRequest) (u *model.Role, e error) {
	u = &model.Role{
		ID:     r.ID,
		Status: int8(2),
	}

	if e = u.Save("Status"); e == nil {
		e = log.AuditLogByUser(r.Session.Staff, u.ID, "role", "archive", "")
	}

	return u, e
}

// UnArchive : function to update status data into active
func UnArchive(r unarchiveRequest) (u *model.Role, e error) {
	u = &model.Role{
		ID:     r.ID,
		Status: int8(1),
	}

	if e = u.Save("Status"); e == nil {
		e = log.AuditLogByUser(r.Session.Staff, u.ID, "role", "unarchive", "")
	}

	return u, e
}
