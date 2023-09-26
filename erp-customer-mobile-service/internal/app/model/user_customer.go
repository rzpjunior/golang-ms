package model

import (
	"time"

	"git.edenfarm.id/edenlabs/edenlabs/orm"
)

type UserCustomer struct {
	ID            int64     `orm:"column(id);auto" json:"-"`
	Code          string    `orm:"column(code);size(50);null" json:"code,omitempty"`
	CustomerID    int64     `orm:"column(customer_id);" json:"-"`
	CustomerIDGP  string    `orm:"column(customer_id_gp);" json:"-"`
	FirebaseToken string    `orm:"column(firebase_token);size(250);null" json:"firebase_token,omitempty"`
	Verification  int8      `orm:"column(verification)" json:"verification,omitempty"`
	TncAccVersion string    `orm:"column(tnc_acc_version);size(50);null" json:"tnc_acc_version,omitempty"`
	TncAccAt      time.Time `orm:"column(tnc_acc_at);type(timestamp);null" json:"tnc_acc_at"`
	LastLoginAt   time.Time `orm:"column(last_login_at);type(timestamp);null" json:"last_login_at"`
	Note          string    `orm:"column(note);size(250);null" json:"note,omitempty"`
	Status        int8      `orm:"column(status);null" json:"status,omitempty"`
	ForceLogout   int8      `orm:"column(force_logout);null" json:"force_logout,omitempty"`
	LoginToken    string    `orm:"column(login_token);null" json:"login_token,omitempty"`
}

func init() {
	orm.RegisterModel(new(UserCustomer))
}

func (m *UserCustomer) TableName() string {
	return "user_customer"
}
