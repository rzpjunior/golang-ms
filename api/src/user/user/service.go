// Copyright 2020 PT. Eden Pangan Indonesia. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

package user

import (
	"strings"

	"git.edenfarm.id/cuxs/common"
	"git.edenfarm.id/cuxs/cuxs/event"
	"git.edenfarm.id/project-version2/api/service/document_history_log"

	"fmt"

	"git.edenfarm.id/cuxs/dbredis"
	"git.edenfarm.id/cuxs/orm"
	"git.edenfarm.id/project-version2/api/log"
	"git.edenfarm.id/project-version2/api/util"
	"git.edenfarm.id/project-version2/datamodel/model"
)

func Save(r createRequest) (u *model.User, e error) {
	o := orm.NewOrm()
	e = o.Begin()
	r.CodeUser, e = util.GenerateCode(r.CodeUser, "user")
	r.CodeStaff, e = util.GenerateCode(r.CodeStaff, "staff")

	if e == nil {
		u = &model.User{
			Code:        r.CodeUser,
			Email:       r.Email,
			Password:    r.PasswordHash,
			Note:        r.Note,
			Status:      int8(1),
			ForceLogout: int8(2),
		}
		if _, e = o.Insert(u); e == nil {
			s := &model.Staff{
				Code:         r.CodeStaff,
				Role:         r.Role,
				User:         u,
				Area:         r.Area,
				Parent:       r.Parent,
				Warehouse:    r.Warehouse,
				Name:         r.Name,
				DisplayName:  r.DisplayName,
				EmployeeCode: r.EmployeeCode,
				RoleGroup:    r.RoleGroup,
				PhoneNumber:  r.PhoneNumber,
				SalesGroupID: r.SalesGroupInt,
				Status:       int8(1),
			}
			if _, e = o.Insert(s); e == nil {
				for _, p := range r.Permission {
					pm := &model.UserPermission{
						User:            u,
						Permission:      p,
						PermissionValue: p.Value,
					}
					if _, e = o.Insert(pm); e != nil {
						e = o.Rollback()
					}
				}
			} else {
				e = o.Rollback()
			}
		} else {
			e = o.Rollback()
		}
		e = log.AuditLogByUser(r.Session.Staff, u.ID, "user", "create", "")
	}
	if e == nil {
		e = o.Commit()
		return
	} else {
		e = o.Rollback()
	}

	return nil, e
}

func (u *updateRequest) Update() (usr *model.User, e error) {
	o := orm.NewOrm()
	o.Begin()
	usr = &model.User{
		ID:   u.ID,
		Note: u.Note,
	}
	if _, e = o.Update(usr, "Note"); e != nil {
		e = o.Rollback()
	}

	e, auditLog := log.AuditLogByUserReturnAuditLogRow(u.Session.Staff, u.ID, "user", "update", "")
	if e != nil {
		e = o.Rollback()
	}

	staff := &model.Staff{
		ID:           u.Staff.ID,
		Role:         u.Role,
		Area:         u.Area,
		Parent:       u.Parent,
		Warehouse:    u.Warehouse,
		Name:         u.Name,
		DisplayName:  u.DisplayName,
		RoleGroup:    u.RoleGroup,
		PhoneNumber:  u.PhoneNumber,
		SalesGroupID: u.SalesGroupInt,
	}
	if _, e = o.Update(staff, "Role", "Area", "Parent", "Warehouse", "Name", "DisplayName", "RoleGroup", "PhoneNumber", "SalesGroupID"); e != nil {
		e = o.Rollback()
	}

	documentHistoryLogData := &document_history_log.UserDocumentHistoryLog{
		OldStaff:   u.OldStaff,
		NewStaff:   staff,
		OldUser:    u.OldUser,
		NewUser:    usr,
		AuditLogID: auditLog.ID,
		RefID:      u.ID,
		Type:       "user",
	}

	go event.Call("document_history_log::user", documentHistoryLogData)

	if e != nil {
		e = o.Rollback()
	}

	o.Commit()
	return
}

func (u *updateWarehouseAccessRequest) UpdateWarehouseAccess() (usr *model.User, e error) {
	o := orm.NewOrm()
	o.Begin()

	usr = u.User
	u.Staff.WarehouseAccessStr = u.WarehouseStr

	if _, e = o.Update(u.Staff, "WarehouseAccessStr"); e != nil {
		o.Rollback()
		return nil, e
	}

	o.Commit()

	return
}

func SaveHelper(r createHelperRequest) (u *model.User, e error) {
	o := orm.NewOrm()
	o1 := orm.NewOrm()
	o1.Using("read_only")
	e = o.Begin()
	r.CodeUser, e = util.GenerateCode(r.CodeUser, "user")
	r.CodeStaff, e = util.GenerateCode(r.CodeStaff, "staff")

	if e == nil {
		u = &model.User{
			Code:     r.CodeUser,
			Email:    r.Email,
			Password: r.PasswordHash,
			Status:   int8(1),
		}
		if _, e = o.Insert(u); e == nil {
			s := &model.Staff{
				Code:        r.CodeStaff,
				Role:        r.Role,
				User:        u,
				Area:        r.Area,
				Parent:      r.Parent,
				Warehouse:   r.Warehouse,
				Name:        r.Name,
				DisplayName: r.Name,
				RoleGroup:   int8(2),
				PhoneNumber: r.PhoneNumber,
				Status:      int8(1),
			}
			if _, e = o.Insert(s); e != nil {
				e = o.Rollback()
			} else {
				// this conditions was used to add permission create DO for checker role
				var rps []*model.RolePermission
				var up *model.UserPermission
				var arrUps []*model.UserPermission
				if r.Role.Code == "ROL0049" || r.Role.Code == "ROL0055" {
					o1.Raw("select * from role_permission rp where rp.role_id = ?", r.Role.ID).QueryRows(&rps)
					if len(rps) > 0 {
						for _, v := range rps {
							v.Permission.Read("ID")
							up = &model.UserPermission{
								User:            u,
								Permission:      v.Permission,
								PermissionValue: v.Permission.Value,
							}
							arrUps = append(arrUps, up)

						}
						if _, e = o.InsertMulti(100, &arrUps); e != nil {
							o.Rollback()
						}
					}
				}
			}
		} else {
			e = o.Rollback()
		}
		e = log.AuditLogByUser(r.Session.Staff, u.ID, "user", "create_helper", "")
	}
	if e == nil {
		e = o.Commit()
		return
	} else {
		e = o.Rollback()
	}

	return nil, e
}

func (u *updateHelperRequest) UpdateHelper() (usr *model.User, e error) {
	o := orm.NewOrm()
	o.Begin()

	staff := &model.Staff{
		ID: u.Staff.ID,
	}

	user := &model.User{
		ID: u.Staff.User.ID,
	}

	if u.Staff.Warehouse != u.Warehouse {
		staff.Warehouse = u.Warehouse
		staff.PhoneNumber = u.PhoneNumber
		staff.Name = u.Name
		staff.Role = u.Role
		user.ForceLogout = 1

		if err := common.PasswordHash(user.Password, u.Password); err != nil {
			user.Password = u.PasswordHash
			if u.Password != "" {
				if _, err = o.Update(user, "password"); err != nil {
					o.Rollback()
				}
			}
			if _, err = o.Update(user, "force_logout"); err == nil {
				if _, err = o.Update(staff, "Warehouse", "PhoneNumber", "Name", "Role"); err != nil {
					o.Rollback()
				}
			} else {
				o.Rollback()
			}
		} else {
			if _, err = o.Update(user, "force_logout"); err == nil {
				if _, err = o.Update(u.Staff, "Warehouse", "PhoneNumber", "Name", "Role"); err != nil {
					o.Rollback()
				}
			} else {
				o.Rollback()
			}
		}
	} else if err := common.PasswordHash(user.Password, u.Password); err != nil {
		staff.Warehouse = u.Warehouse
		staff.PhoneNumber = u.PhoneNumber
		staff.Name = u.Name
		staff.Role = u.Role
		user.ForceLogout = 1
		user.Password = u.PasswordHash
		if u.Password != "" {
			if _, err = o.Update(user, "password"); err != nil {
				o.Rollback()
			}
		}
		if _, err = o.Update(user, "force_logout"); err == nil {
			if _, err = o.Update(staff, "Warehouse", "PhoneNumber", "Name", "Role"); err != nil {
				o.Rollback()
			}
		} else {
			o.Rollback()
		}
	} else {
		staff.Warehouse = u.Warehouse
		staff.PhoneNumber = u.PhoneNumber
		staff.Name = u.Name
		staff.Role = u.Role

		if _, err = o.Update(staff, "Warehouse", "PhoneNumber", "Name", "Role"); err != nil {
			o.Rollback()
		}
	}

	e, auditLog := log.AuditLogByUserReturnAuditLogRow(u.Session.Staff, user.ID, "user", "update_helper", "")
	if e != nil {
		e = o.Rollback()
	}

	documentHistoryLogData := &document_history_log.UserDocumentHistoryLog{
		OldStaff:   u.OldStaff,
		NewStaff:   staff,
		OldUser:    u.OldUser,
		NewUser:    user,
		AuditLogID: auditLog.ID,
		RefID:      user.ID,
		Type:       "user",
	}

	go event.Call("document_history_log::user", documentHistoryLogData)

	if e != nil {
		e = o.Rollback()
	}

	o.Commit()
	return
}

// ArchiveHelper : function to update status data into archive
func ArchiveHelper(r archiveHelperRequest) (u *model.User, e error) {
	o := orm.NewOrm()
	o.Begin()

	r.Staff.Status = int8(2)
	r.Staff.User.Status = int8(2)

	if _, e = o.Update(r.Staff, "status"); e == nil {
		if _, e = o.Update(r.Staff.User, "status"); e == nil {
			e = log.AuditLogByUser(r.Session.Staff, r.Staff.ID, "helper", "archive_helper", "")
		} else {
			o.Rollback()
		}
	} else {
		o.Rollback()
	}

	o.Commit()
	return u, e
}

// ArchiveHelper : function to update status data into archive
func UnarchiveHelper(r unarchiveHelperRequest) (u *model.User, e error) {
	o := orm.NewOrm()
	o.Begin()

	r.Staff.Status = int8(1)
	r.Staff.User.Status = int8(1)
	r.Staff.User.ForceLogout = int8(2)
	if _, e = o.Update(r.Staff, "status"); e == nil {
		if _, e = o.Update(r.Staff.User, "status", "ForceLogout"); e == nil {
			e = log.AuditLogByUser(r.Session.Staff, r.Staff.ID, "helper", "unarchive_helper", "")
		} else {
			o.Rollback()
		}
	} else {
		o.Rollback()
	}

	o.Commit()
	return u, e
}

func (u *resetPasswordRequest) Reset() (usr *model.User, e error) {
	usr = &model.User{
		ID:       u.ID,
		Password: u.PasswordHash,
	}
	usr.Save("Password")
	return
}

func (u *updatePermissionRequest) UpdatePermission() (usr *model.User, e error) {
	//STEP 1, join old-permission array with new-permission array from the request
	o := orm.NewOrm()
	o.Begin()
	OldPermission := u.OldUserPermissionID
	var z []int64
	for _, row := range u.NewPermission {
		u.OldUserPermissionID = append(u.OldUserPermissionID, row)
	}
	//STEP 2, get unique-value array from the joined old-new array from step 1
	differentPermission := util.GetUniqueValue(u.OldUserPermissionID)

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
		_, e = o.Raw("DELETE FROM user_permission WHERE user_id = ? AND permission_id IN ("+catLength+")", u.ID, willDeletedPermission).Exec()
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
				up := &model.UserPermission{
					Permission:      p,
					User:            &model.User{ID: u.ID},
					PermissionValue: p.Value,
				}
				_, e = o.Insert(up)
			}
		}
	}
	e = log.AuditLogByUser(u.Session.Staff, u.ID, "user", "update_user_permission", "")

	if e == nil {
		e = o.Commit()
		key := fmt.Sprintf("SELECT COUNT(id) FROM user_permission WHERE user_id = %d *", u.ID)
		dbredis.Redis.DeleteCacheWhereLike(key)

	} else {
		e = o.Rollback()
	}

	return
}

// Archive : function to update status data into archive
func Archive(r archiveRequest) (u *model.User, e error) {
	u = &model.User{
		ID:          r.ID,
		Status:      int8(2),
		ForceLogout: int8(1),
	}

	if e = u.Save("Status", "ForceLogout"); e == nil {
		orm.NewOrm().Raw("DELETE FROM user_permission WHERE user_id = ?", r.ID).Exec()
		e = log.AuditLogByUser(r.Session.Staff, u.ID, "user", "archive", "")
	}

	return u, e
}

// UnArchive : function to update status data into active
func UnArchive(r unarchiveRequest) (u *model.User, e error) {
	u = &model.User{
		ID:          r.ID,
		Status:      int8(1),
		ForceLogout: int8(2),
	}

	if e = u.Save("Status", "ForceLogout"); e == nil {
		e = log.AuditLogByUser(r.Session.Staff, u.ID, "user", "unarchive", "")
	}

	return u, e
}

// Delete : function to update status data into deleted
func Delete(r deleteRequest) (u *model.User, e error) {
	u = &model.User{
		ID:     r.ID,
		Status: int8(3),
	}

	if e = u.Save("Status"); e == nil {
		e = log.AuditLogByUser(r.Session.Staff, u.ID, "user", "delete", r.DeletionNote)
	}

	return u, e
}
