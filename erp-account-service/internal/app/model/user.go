package model

import (
	"time"

	"git.edenfarm.id/edenlabs/edenlabs/orm"
)

type User struct {
	ID                     int64     `orm:"column(id)" json:"id"`
	Name                   string    `orm:"column(name)" json:"name"`
	Nickname               string    `orm:"column(nickname)" json:"nickname"`
	Email                  string    `orm:"column(email)" json:"email"`
	Password               string    `orm:"column(password)" json:"password"`
	RegionID               int64     `orm:"column(region_id)" json:"region_id"`
	RegionIDGP             string    `orm:"column(region_id_gp)" json:"region_id_gp"`
	ParentID               int64     `orm:"column(parent_id)" json:"parent_id"`
	SiteID                 int64     `orm:"column(site_id)" json:"site_id"`
	SiteIDGP               string    `orm:"column(site_id_gp)" json:"site_id_gp"`
	TerritoryID            int64     `orm:"column(territory_id)" json:"territory_id"`
	TerritoryIDGP          string    `orm:"column(territory_id_gp)" json:"territory_id_gp"`
	EmployeeCode           string    `orm:"column(employee_code)" json:"employee_code"`
	PhoneNumber            string    `orm:"column(phone_number)" json:"phone_number"`
	Status                 int8      `orm:"column(status)" json:"status"`
	Note                   string    `orm:"column(note)" json:"note"`
	ForceLogout            int32     `orm:"column(force_logout)" json:"force_logout"`
	SalesAppLoginToken     string    `orm:"column(salesapp_login_token)" json:"salesapp_login_token"`
	SalesAppNotifToken     string    `orm:"column(salesapp_notif_token)" json:"salesapp_notif_token"`
	PurchaserAppLoginToken string    `orm:"column(purchaser_login_token)" json:"purchaser_login_token"`
	PurchaserAppNotifToken string    `orm:"column(purchaser_notif_token)" json:"purchaser_notif_token"`
	EdnAppLoginToken       string    `orm:"column(edn_app_login_token)" json:"edn_app_login_token"`
	CreatedAt              time.Time `orm:"column(created_at)" json:"created_at"`
	UpdatedAt              time.Time `orm:"column(updated_at)" json:"updated_at"`
}

func init() {
	orm.RegisterModel(new(User))
}

func (m *User) TableName() string {
	return "user"
}
