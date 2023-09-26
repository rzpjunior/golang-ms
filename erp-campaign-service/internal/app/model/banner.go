package model

import (
	"time"

	"git.edenfarm.id/edenlabs/edenlabs/orm"
)

type Banner struct {
	ID            int64     `orm:"column(id)" json:"-"`
	Regions       string    `orm:"column(regions)" json:"regions"`
	Archetypes    string    `orm:"column(archetypes)" json:"archetypes"`
	Name          string    `orm:"column(name)" json:"name"`
	Code          string    `orm:"column(code)" json:"code"`
	Queue         int       `orm:"column(queue)" json:"queue"`
	RedirectTo    int8      `orm:"column(redirect_to)" json:"redirect_to"`
	RedirectValue string    `orm:"column(redirect_value)" json:"redirect_value"`
	ImageUrl      string    `orm:"column(image_url)" json:"image_url"`
	StartAt       time.Time `orm:"column(start_at)" json:"start_at"`
	FinishAt      time.Time `orm:"column(finish_at)" json:"finish_at"`
	Note          string    `orm:"column(note)" json:"note"`
	Status        int8      `orm:"column(status)" json:"status"`
	CreatedAt     time.Time `orm:"column(created_at)" json:"created_at"`
	UpdatedAt     time.Time `orm:"column(updated_at)" json:"updated_at"`
}

func init() {
	orm.RegisterModel(new(Banner))
}

func (m *Banner) TableName() string {
	return "banner"
}
