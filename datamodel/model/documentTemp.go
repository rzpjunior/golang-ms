// Copyright 2020 PT. Eden Pangan Indonesia. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

package model

import (
	"encoding/json"

	"git.edenfarm.id/cuxs/common"
	"git.edenfarm.id/cuxs/orm"
	"time"
)

func init() {
	orm.RegisterModel(new(DocumentTemp))
}

// DocumentTemp model for documentTemp table.
type DocumentTemp struct {
	ID             int64     `orm:"column(id);auto" json:"-"`
	SalesOrderCode string    `orm:"column(sales_order_code);" json:"sales_order_code"`
	SalesOrderID   int64     `orm:"column(sales_order_id)" json:"sales_order_id"`
	CreatedAt      time.Time `orm:"column(created_at);type(timestamp);null" json:"created_at"`
	FromCronJob    int       `orm:"column(from_cronjob)" json:"from_cronjob"`
}

// MarshalJSON customized data struct when marshaling data
// into JSON format, all Primary key & Foreign key will be encrypted.
func (m *DocumentTemp) MarshalJSON() ([]byte, error) {
	type Alias DocumentTemp

	alias := &struct {
		ID            string `json:"id"`
		StatusConvert string `json:"status_convert"`
		*Alias
	}{
		ID:    common.Encrypt(m.ID),
		Alias: (*Alias)(m),
	}

	return json.Marshal(alias)
}

// Save inserting or updating User struct into user table.
// It will updating if this struct has valid Id
// if not, will inserting a new row to user.
// The field parameter is an field that will be saved, it is
// usefull for partial updating data.
func (m *DocumentTemp) Save(fields ...string) (err error) {
	o := orm.NewOrm()
	if m.ID > 0 {
		_, err = o.Update(m, fields...)
	} else {
		m.ID, err = o.Insert(m)
	}
	return
}
