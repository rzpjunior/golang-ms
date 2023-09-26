// Copyright 2020 PT. Eden Pangan Indonesia. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

package model

import (
	"time"

	"git.edenfarm.id/edenlabs/edenlabs/orm"
)

// Banner model for audit_log table.
type Banner struct {
	ID             int64     `orm:"column(id);auto" json:"-"`
	Code           string    `orm:"column(code);size(50);null" json:"code"`
	Name           string    `orm:"column(name);size(100);null" json:"name"`
	ImageUrl       string    `orm:"column(image_url);size(300);null" json:"image_url"`
	StartDate      time.Time `orm:"column(start_date);type(timestamp);null" json:"start_date"`
	EndDate        time.Time `orm:"column(end_date);type(timestamp);null" json:"end_date"`
	NavigationType int8      `orm:"column(navigate_type)" json:"navigate_type"`
	NavigationUrl  string    `orm:"column(navigate_url);size(300);null" json:"navigate_url"`
	Region         int64     `orm:"column(region)" json:"region"`
	Archetype      int64     `orm:"column(archetype)" json:"archetype"`
	Queue          int8      `orm:"column(queue)" json:"queue"`
	Note           string    `orm:"column(note);size(250);null" json:"note"`
	Status         int8      `orm:"column(status);null" json:"status"`

	ItemCategory *ItemCategory `orm:"-" json:"item_category"`
	Item         *Item         `orm:"-" json:"item"`

	// log
	CreatedAt time.Time `orm:"column(created_at)" json:"created_at"`
	CreatedBy int64     `orm:"column(created_by)" json:"created_by"`
}

func init() {
	orm.RegisterModel(new(Banner))
}
