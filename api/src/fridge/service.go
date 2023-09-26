// //// Copyright 2021 PT. Eden Pangan Indonesia. All rights reserved.
// // Use of this source code is governed by a MIT style
// // license that can be found in the LICENSE file.

package fridge

import (
	"time"

	"git.edenfarm.id/cuxs/orm"
	"git.edenfarm.id/project-version2/api/log"
	"git.edenfarm.id/project-version2/api/util"
	"git.edenfarm.id/project-version2/datamodel/model"
)

// // Save : function to save data requested into database
// func Save(r createRequest) (tempBox *model.BoxFridge, e error) {
// 	docCode := "TES"
// 	var i int
// 	o := orm.NewOrm()
// 	o.Begin()
// 	orSelect := orm.NewOrm()
// 	orSelect.Using("read_only")
// 	if e := orSelect.Raw("Select count(*) from box_fridge "+
// 		"where box_id = ? and fridge_id = ? order by id desc limit 1", r.Box.ID, r.Fridge.ID).QueryRow(&i); e != nil {
// 		return nil, e
// 	}
// 	if i >= 1 {

// 		if _, e := o.Raw("Update box_fridge "+
// 			"set status = ?,position=? where box_id = ?  and fridge_id = ? order by id desc limit 1", r.Status, r.Position, r.Box.ID, r.Fridge.ID).Exec(); e != nil {
// 			o.Rollback()
// 			return nil, e
// 		}
// 		if e := orSelect.Raw("Select * from box_fridge "+
// 			"where box_id = ? and fridge_id = ? order by id desc limit 1", r.Box.ID, r.Fridge.ID).QueryRow(&tempBox); e != nil {
// 			o.Rollback()
// 			return nil, e
// 		}
// 		o.Commit()
// 		return tempBox, e
// 	}

// 	if _, e := o.Raw("Update box_fridge "+
// 		"set status = 2 where box_id = ?   ", r.Box.ID).Exec(); e != nil {
// 		o.Rollback()
// 		return nil, e
// 	}

// 	tempBox = &model.BoxFridge{
// 		Code:          docCode,
// 		Box:           r.Box,
// 		Fridge:        r.Fridge,
// 		Status:        r.Status,
// 		CreatedAt:     time.Now(),
// 		Position:      r.Position,
// 		CreatedBy:     r.Session.Staff.ID,
// 		LastUpdatedAt: time.Now(),
// 		LastUpdatedBy: r.Session.Staff.ID,
// 	}
// 	_, e = o.Insert(tempBox)
// 	if e != nil {
// 		o.Rollback()
// 		return nil, e
// 	}

// 	if e = log.AuditLogByUser(r.Session.Staff, 1, "box_fridge", "create", ""); e != nil {
// 		o.Rollback()
// 		return nil, e
// 	}

// 	o.Commit()

// 	return tempBox, e
// }

// // Save : function to save data requested into database
// func SaveTransaction(r createRequestTransaction) (tempBox *model.FridgeTransactionHeader, e error) {
// 	o := orm.NewOrm()
// 	o.Begin()

// 	orSelect := orm.NewOrm()
// 	orSelect.Using("read_only")

// 	var tempBox2 []*model.FridgeTransactionDetail
// 	tempBox = &model.FridgeTransactionHeader{
// 		Code:          "TES",
// 		TotalPrice:    r.TotalPrice,
// 		TotalCharge:   r.TotalCharge,
// 		CreatedAt:     time.Now(),
// 		CreatedBy:     r.Session.Staff.ID,
// 		LastUpdatedAt: time.Now(),
// 		LastUpdatedBy: r.Session.Staff.ID,
// 		Fridge:        r.Fridge,
// 		Status:        1,
// 		Branch:        r.BranchFridge.Branch,
// 	}
// 	_, e = o.Insert(tempBox)
// 	if e != nil {
// 		o.Rollback()
// 		return nil, e
// 	}

// 	for _, v := range r.ProductBox {
// 		ftDetail := &model.FridgeTransactionDetail{
// 			Code:                    "TES",
// 			TotalPrice:              v.TotalCharge,
// 			TotalCharge:             v.TotalCharge,
// 			TotalWeight:             v.TotalWeight,
// 			CreatedAt:               time.Now(),
// 			CreatedBy:               r.Session.Staff.ID,
// 			LastUpdatedAt:           time.Now(),
// 			LastUpdatedBy:           r.Session.Staff.ID,
// 			FridgeTransactionHeader: tempBox,
// 			Box:                     v.Box,
// 			Product:                 v.ProductBox.Product,
// 			Status:                  1,
// 		}
// 		tempBox2 = append(tempBox2, ftDetail)
// 		if _, e := o.Raw("Update box_fridge "+
// 			"set status = 3 where box_id = ?  and fridge_id = ? order by id desc limit 1", v.Box.ID, r.Fridge.ID).Exec(); e != nil {
// 			return nil, e
// 		}
// 		if _, e := o.Raw("Update product_box "+
// 			"set status = 3 where box_id = ?  order by id desc limit 1", v.Box.ID).Exec(); e != nil {
// 			return nil, e
// 		}
// 	}
// 	_, e = o.InsertMulti(100, &tempBox2)
// 	if e != nil {
// 		o.Rollback()
// 		return nil, e
// 	}

// 	o.Commit()

// 	return tempBox, e
// }

// SessionData struktur data current user logged in.
type SessionData struct {
	UserFridge *model.UserFridge `json:"user_fridge"`
}

// Save : function to save data requested into database
func SaveUser(r createRequestUser) (c *model.UserFridge, err error) {
	o := orm.NewOrm()
	o.Begin()

	r.Code, err = util.GenerateCode(r.Code, "user_fridge", 6)
	if err != nil {
		o.Rollback()
		return nil, err
	}
	c = &model.UserFridge{
		Code:        r.Code,
		Username:    r.Username,
		Password:    r.PasswordHash,
		Status:      1,
		CreatedAt:   time.Now(),
		LastLoginAt: time.Time{},
		Note:        "",
		Branch:      r.Branch,
		Warehouse:   r.Warehouse,
	}

	if _, err = o.Insert(c); err != nil {
		o.Rollback()
		return nil, err
	}

	err = log.AuditLogByUser(r.Session.Staff, c.ID, "fridge", "create", "Create User Fridge")
	if err != nil {
		o.Rollback()
		return nil, err
	}

	o.Commit()

	return c, err
}

func Login(userFridge *model.UserFridge) (sd *model.UserFridge, e error) {
	userFridge.LastLoginAt = time.Now()
	if e := userFridge.Save("LastLoginAt"); e != nil {
		return nil, e
	}
	return userFridge, nil
}
