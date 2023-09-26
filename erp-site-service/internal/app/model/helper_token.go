package model

import (
	"git.edenfarm.id/edenlabs/edenlabs/orm"
)

type HelperToken struct {
	Id                int64  `orm:"column(id)" json:"id"`
	HelperIdGp        string `orm:"column(helper_id_gp)" json:"helper_id_gp"`
	NotificationToken string `orm:"column(notif_token)" json:"notification_token"`
}

func init() {
	orm.RegisterModel(new(HelperToken))
}

func (m *HelperToken) TableName() string {
	return "helper_token"
}
