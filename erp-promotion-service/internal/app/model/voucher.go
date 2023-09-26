package model

import (
	"encoding/json"
	"time"

	"git.edenfarm.id/edenlabs/edenlabs/orm"
)

type Voucher struct {
	ID                     int64     `orm:"column(id)" json:"id"`
	RegionIDGP             string    `orm:"column(region_id_gp)" json:"region_id_gp"`
	CustomerID             int64     `orm:"column(customer_id)" json:"customer"`
	CustomerTypeIDGP       string    `orm:"column(customer_type_id_gp)" json:"customer_type_id_gp"`
	ArchetypeIDGP          string    `orm:"column(archetype_id_gp)" json:"archetype_id_gp"`
	MembershipLevelID      int64     `orm:"column(membership_level_id)" json:"membership_level_id"`
	MembershipCheckPointID int64     `orm:"column(membership_checkpoint_id)" json:"membership_checkpoint_id"`
	DivisionID             int64     `orm:"column(division_id)" json:"division_id"`
	Code                   string    `orm:"column(code)" json:"code"`
	RedeemCode             string    `orm:"column(redeem_code)" json:"redeem_code"`
	Name                   string    `orm:"column(name)" json:"name"`
	Type                   int8      `orm:"column(type)" json:"type"`
	StartTime              time.Time `orm:"column(start_time)" json:"start_time"`
	EndTime                time.Time `orm:"column(end_time)" json:"end_time"`
	OverallQuota           int64     `orm:"column(overall_quota);" json:"overall_quota"`
	UserQuota              int64     `orm:"column(user_quota);" json:"user_quota"`
	RemOverallQuota        int64     `orm:"column(rem_overall_quota);" json:"rem_overall_quota"`
	MinOrder               float64   `orm:"column(min_order);null;digits(20);decimals(2)" json:"min_order"`
	DiscAmount             float64   `orm:"column(disc_amount);null;digits(20);decimals(2)" json:"disc_amount"`
	TermConditions         string    `orm:"column(term_conditions)" json:"term_conditions"`
	ImageUrl               string    `orm:"column(image_url)" json:"image_url"`
	VoidReason             int8      `orm:"column(void_reason);null" json:"void_reason"`
	Note                   string    `orm:"column(note)" json:"note"`
	Status                 int8      `orm:"column(status);null" json:"status"`
	VoucherItem            int8      `orm:"column(voucher_item);null" json:"voucher_item"`
	CreatedAt              time.Time `orm:"column(created_at)" json:"created_at"`
	RemUserQuota           int64     `orm:"-" json:"rem_user_quota"`
	TotalUserUsed          int64     `orm:"-" json:"total_user_used"`
}

func init() {
	orm.RegisterModel(new(Voucher))
}

func (m *Voucher) TableName() string {
	return "voucher"
}

func (m *Voucher) MarshalJSON() ([]byte, error) {
	type Alias Voucher

	return json.Marshal(&struct {
		*Alias
	}{
		Alias: (*Alias)(m),
	})
}
