package model

import (
	"time"

	"git.edenfarm.id/edenlabs/edenlabs/orm"
)

type PushNotification struct {
	ID            int64     `orm:"column(id)" json:"id"`
	Code          string    `orm:"column(code)" json:"code"`
	CampaignName  string    `orm:"column(campaign_name)" json:"name"`
	Regions       string    `orm:"column(regions)" json:"regions"`
	Archetypes    string    `orm:"column(archetypes)" json:"archetypes"`
	RedirectTo    int8      `orm:"column(redirect_to)" json:"redirect_to"`
	RedirectValue string    `orm:"column(redirect_value)" json:"redirect_value"`
	Title         string    `orm:"column(title)" json:"title"`
	Message       string    `orm:"column(message)" json:"message"`
	PushNow       int8      `orm:"column(push_now)" json:"push_now"`
	ScheduledAt   time.Time `orm:"column(scheduled_at)" json:"scheduled_at"`
	Status        int8      `orm:"column(status)" json:"status"`
	SuccessSent   int       `orm:"column(success_sent)" json:"success_sent"`
	FailedSent    int       `orm:"column(failed_sent)" json:"failed_sent"`
	Opened        int       `orm:"column(opened)" json:"opened"`
	CreatedAt     time.Time `orm:"column(created_at)" json:"created_at"`
	CreatedBy     int64     `orm:"column(created_by)" json:"created_by"`
	UpdatedAt     time.Time `orm:"column(updated_at)" json:"updated_at"`
}

func init() {
	orm.RegisterModel(new(PushNotification))
}

func (m *PushNotification) TableName() string {
	return "push_notification"
}
