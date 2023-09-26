// Copyright 2020 PT. Eden Pangan Indonesia. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

package model

import (
	"encoding/json"
	"time"

	"git.edenfarm.id/cuxs/common"
	"git.edenfarm.id/cuxs/orm"
)

func init() {
	orm.RegisterModel(new(NotificationCampaign))
}

// NotificationCampaign : struct to hold NotificationCampaign model data for database
type NotificationCampaign struct {
	ID                int64        `orm:"column(id);auto" json:"-"`
	Code              string       `orm:"column(code);" json:"code"`
	CampaignName      string       `orm:"column(campaign_name);" json:"campaign_name"`
	Area              string       `orm:"column(area)" json:"area"`
	AreaArr           []*Area      `orm:"-" json:"area_arr,omitempty"`
	AreaName          string       `orm:"-" json:"area_name"`
	AreaNameArr       []string     `orm:"-" json:"area_name_arr"`
	Archetype         string       `orm:"column(archetype)" json:"archetype"`
	ArchetypeArr      []*Archetype `orm:"-" json:"archetype_arr,omitempty"`
	ArchetypeName     string       `orm:"-" json:"archetype_name"`
	ArchetypeNameArr  []string     `orm:"-" json:"archetype_name_arr"`
	RedirectTo        int8         `orm:"column(redirect_to);" json:"redirect_to"`
	RedirectValue     string       `orm:"column(redirect_value);" json:"redirect_value"`
	RedirectToName    string       `orm:"-" json:"redirect_to_name"`
	RedirectValueName string       `orm:"-" json:"redirect_value_name"`
	Title             string       `orm:"column(title);" json:"title"`
	Message           string       `orm:"column(message);" json:"message"`
	PushNow           int8         `orm:"column(push_now);" json:"push_now"`
	PushNowStatus     bool         `orm:"-" json:"-"`
	ScheduledAt       time.Time    `orm:"column(scheduled_at);" json:"scheduled_at"`
	Status            int8         `orm:"column(status);" json:"status"`
	CreatedAt         time.Time    `orm:"column(created_at);type(timestamp);null" json:"created_at"`
	UpdatedAt         time.Time    `orm:"column(updated_at);type(timestamp);null" json:"updated_at"`
	SuccessSent       int64        `orm:"column(success_sent)" json:"success_sent"`
	FailedSent        int64        `orm:"column(failed_sent)" json:"failed_sent"`
	Open              int64        `orm:"column(open)" json:"open"`
	Conversion        int64        `orm:"-" json:"conversion"`
}

// MarshalJSON : function to customize data struct into json. encrypt all key ids
func (m *NotificationCampaign) MarshalJSON() ([]byte, error) {
	type Alias NotificationCampaign

	return json.Marshal(&struct {
		ID string `json:"id"`
		*Alias
	}{
		ID:    common.Encrypt(m.ID),
		Alias: (*Alias)(m),
	})
}

// Save : function to save data into database. will update when it has valid id, create otherwise
func (m *NotificationCampaign) Save(fields ...string) (err error) {
	o := orm.NewOrm()
	if m.ID > 0 {
		_, err = o.Update(m, fields...)
	} else {
		m.ID, err = o.Insert(m)
	}
	return
}

// Read : function to get data from database
func (m *NotificationCampaign) Read(fields ...string) error {
	o := orm.NewOrm()
	o.Using("read_only")
	return o.Read(m, fields...)
}
