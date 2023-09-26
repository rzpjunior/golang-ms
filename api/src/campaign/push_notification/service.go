// Copyright 2020 PT. Eden Pangan Indonesia. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

package push_notification

import (
	"strconv"
	"time"

	"git.edenfarm.id/cuxs/common"
	"git.edenfarm.id/cuxs/orm"

	"git.edenfarm.id/project-version2/api/datastore/repository"
	"git.edenfarm.id/project-version2/api/log"
	"git.edenfarm.id/project-version2/api/util"
	"git.edenfarm.id/project-version2/datamodel/model"
)

// Save : function to save data requested into database
func Save(r createRequest) (pushNotification *model.NotificationCampaign, e error) {
	o := orm.NewOrm()
	o.Begin()
	// set charset to utf8mb4 for supporting emoji connection
	o.Raw("SET NAMES 'utf8mb4'").Exec()

	var areas string
	var archetypes string

	for i, a := range r.Areas {
		id := strconv.Itoa(int(a.ID))
		if i != 0 {
			areas += ","
		}
		areas += id
	}

	for i, a := range r.Archetypes {
		id := strconv.Itoa(int(a.ID))
		if i != 0 {
			archetypes += ","
		}
		archetypes += id
	}

	// generate code notification_campaign
	if r.Code, e = util.GenerateCode(r.Code, "notification_campaign"); e != nil {
		o.Rollback()
		return nil, e
	}

	pushNotification = &model.NotificationCampaign{
		Code:          r.Code,
		CampaignName:  r.CampaignName,
		Area:          areas,
		Archetype:     archetypes,
		RedirectTo:    r.RedirectTo,
		RedirectValue: r.RedirectValue,
		Title:         r.Title,
		Message:       r.Message,
		PushNow:       r.PushNow,
		ScheduledAt:   r.ScheduledAt,
		CreatedAt:     time.Now(),
	}

	// checking push_now
	var pushNowGlossary *model.Glossary
	if pushNowGlossary, e = repository.GetGlossaryMultipleValue("table", "notification_campaign", "attribute", "push_now", "value_int", r.PushNow); e != nil {
		o.Rollback()
		return nil, e
	}
	// checking push now
	if pushNowGlossary.ValueName == "yes" {
		pushNotification.Status = 2
	} else {
		pushNotification.Status = 1
	}

	if pushNotification.ID, e = o.Insert(pushNotification); e != nil {
		o.Rollback()
		return nil, e
	}

	if e = log.AuditLogByUser(r.Session.Staff, pushNotification.ID, "notification_campaign", "create", ""); e != nil {
		o.Rollback()
		return nil, e
	}

	o.Commit()
	return pushNotification, e
}

// Update: function to update data from database
func Update(r updateRequest) (pushNotification *model.NotificationCampaign, e error) {
	o := orm.NewOrm()
	o.Begin()

	// set charset to utf8mb4 for supporting emoji connection
	o.Raw("SET NAMES 'utf8mb4'").Exec()

	var areas string
	var archetypes string

	for i, a := range r.Areas {
		id := strconv.Itoa(int(a.ID))
		if i != 0 {
			areas += ","
		}
		areas += id
	}

	for i, a := range r.Archetypes {
		id := strconv.Itoa(int(a.ID))
		if i != 0 {
			archetypes += ","
		}
		archetypes += id
	}

	pushNotification = &model.NotificationCampaign{
		ID:            r.ID,
		CampaignName:  r.CampaignName,
		Area:          areas,
		Archetype:     archetypes,
		RedirectTo:    r.RedirectTo,
		RedirectValue: r.RedirectValue,
		Title:         r.Title,
		Message:       r.Message,
		PushNow:       r.PushNow,
		ScheduledAt:   r.ScheduledAt,
		UpdatedAt:     time.Now(),
	}

	if _, e = o.Update(pushNotification, "campaign_name", "area", "archetype", "redirect_to", "redirect_value", "title", "message", "push_now", "scheduled_at", "updated_at"); e != nil {
		o.Rollback()
		return nil, e
	}

	if e = log.AuditLogByUser(r.Session.Staff, pushNotification.ID, "notification_campaign", "update", ""); e != nil {
		o.Rollback()
		return nil, e
	}

	o.Commit()

	return pushNotification, e
}

// Cancel: function to change data status into cancel
func Cancel(r cancelRequest) (pn *model.NotificationCampaign, e error) {
	o := orm.NewOrm()
	o.Begin()

	// checking status cancel to glossary
	var cancelGlossary *model.Glossary
	cancelGlossary, e = repository.GetGlossaryMultipleValue("table", "notification_campaign", "attribute", "status", "value_name", "cancel")
	if e != nil {
		o.Rollback()
		return nil, e
	}

	// update status to cancel
	r.NotificationCampaign.Status = cancelGlossary.ValueInt
	if _, e = o.Update(r.NotificationCampaign, "Status"); e != nil {
		o.Rollback()
	}

	if e = log.AuditLogByUser(r.Session.Staff, r.NotificationCampaign.ID, "notification_campaign", "cancel", r.Note); e != nil {
		o.Rollback()
		return nil, e
	}

	o.Commit()
	return r.NotificationCampaign, nil
}

// check area is existing
func CheckArea(arr []string, areaID string) (status bool, e error) {
	for _, a := range arr {
		id, e := common.Decrypt(areaID)
		if e != nil {
			return false, e
		}
		a, e := strconv.ParseInt(a, 10, 64)
		if e != nil {
			return false, e
		}
		if a == id {
			return true, nil
		}
	}
	return false, nil
}
