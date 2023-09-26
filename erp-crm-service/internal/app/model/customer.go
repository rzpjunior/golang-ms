package model

import (
	"time"

	"git.edenfarm.id/edenlabs/edenlabs/orm"
)

type Customer struct {
	ID                      int64     `orm:"column(id)" json:"id"`
	CustomerIDGP            string    `orm:"column(customer_id_gp)" json:"customer_id_gp"`
	ProspectiveCustomerID   int64     `orm:"column(prospective_customer_id)" json:"prospective_customer_id"`
	MembershipLevelID       int64     `orm:"column(membership_level_id)" json:"membership_level_id"`
	MembershipCheckpointID  int64     `orm:"column(membership_checkpoint_id)" json:"membership_checkpoint_id"`
	TotalPoint              int64     `orm:"column(total_point)" json:"total_point"`
	ProfileCode             string    `orm:"column(profile_code)" json:"profile_code"`
	Email                   string    `orm:"column(email)" json:"email"`
	Password                string    `orm:"column(password)" json:"password"`
	ReferenceInfo           string    `orm:"column(reference_info)" json:"reference_info"`
	UpgradeStatus           int8      `orm:"column(upgrade_status)" json:"upgrade_status"`
	KtpPhotosUrl            string    `orm:"column(ktp_photos_url)" json:"ktp_photos_url"`
	CustomerPhotosUrl       string    `orm:"column(customer_photos_url)" json:"customer_photos_url"`
	CustomerSelfieUrl       string    `orm:"column(customer_selfie_url)" json:"customer_selfie_url"`
	CreatedAt               time.Time `orm:"column(created_at)" json:"created_at"`
	UpdatedAt               time.Time `orm:"column(updated_at)" json:"updated_at"`
	MembershipRewardID      int64     `orm:"column(membership_reward_id)" json:"membership_reward_id"`
	MembershipRewardAmmount float64   `orm:"column(membership_reward_amount)" json:"membership_reward_amount"`
	ReferrerID              int64     `orm:"column(referrer_id)" json:"referrer_id"`
	ReferrerCode            string    `orm:"column(referrer_code)" json:"referrer_code"`
	ReferralCode            string    `orm:"column(referral_code)" json:"referral_code"`
	Gender                  int       `orm:"column(gender)" json:"gender,omitempty"`
	BirthDate               time.Time `orm:"column(birth_date)" json:"birth_date,omitempty"`
}

func init() {
	orm.RegisterModel(new(Customer))
}

func (m *Customer) TableName() string {
	return "customer"
}
