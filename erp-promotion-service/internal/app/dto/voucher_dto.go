package dto

import "time"

type VoucherResponse struct {
	ID              int64     `json:"id"`
	Code            string    `json:"code"`
	RedeemCode      string    `json:"redeem_code"`
	Name            string    `json:"name"`
	Type            int8      `json:"type"`
	StartTime       time.Time `json:"start_time"`
	EndTime         time.Time `json:"end_time"`
	OverallQuota    int64     `json:"overall_quota"`
	UserQuota       int64     `json:"user_quota"`
	RemOverallQuota int64     `json:"rem_overall_quota"`
	MinOrder        float64   `json:"min_order"`
	DiscAmount      float64   `json:"disc_amount"`
	TermConditions  string    `json:"term_conditions"`
	ImageUrl        string    `json:"image_url"`
	VoidReason      int8      `json:"void_reason"`
	Note            string    `json:"note"`
	Status          int8      `json:"status"`
	StatusConvert   string    `json:"status_convert"`
	VoucherItem     int8      `json:"voucher_item"`
	CreatedAt       time.Time `json:"created_at"`
	RemUserQuota    int64     `json:"rem_user_quota"`

	Archetype            *ArchetypeResponse            `json:"archetype"`
	Region               *RegionResponse               `json:"region"`
	Division             *DivisionResponse             `json:"division"`
	Customer             *CustomerResponse             `json:"customer"`
	MembershipLevel      *MembershipLevelResponse      `json:"membership_level"`
	MembershipCheckpoint *MembershipCheckpointResponse `json:"membership_checkpoint"`
	VoucherItems         []*VoucherItemResponse        `json:"voucher_items"`
}

type VoucherRequestGet struct {
	Search                 string `json:"search"`
	ArchetypeID            string `json:"archetype_id"`
	RegionID               string `json:"region_id"`
	MembershipLevelID      int64  `json:"membership_level_id"`
	MembershipCheckpointID int64  `json:"membership_check_point_id"`
	CustomerID             int64  `json:"customer_id"`
	Type                   int8   `json:"type"`
	Status                 int8   `json:"status"`
	OrderBy                string `json:"order_by"`
	Offset                 int64  `json:"offset"`
	Limit                  int64  `json:"limit"`
	CustomerTypeID         string `json:"customer_type_id"`
}

type VoucherRequestCreate struct {
	RegionID               string                     `json:"region_id" valid:"required"`
	CustomerID             int64                      `json:"customer_id"`
	CustomerTypeID         string                     `json:"customer_type_id" valid:"required"`
	ArchetypeID            string                     `json:"archetype_id" valid:"required"`
	MembershipLevelID      int64                      `json:"membership_level_id"`
	MembershipCheckPointID int64                      `json:"membership_checkpoint_id"`
	DivisionID             int64                      `json:"division_id"`
	RedeemCode             string                     `json:"redeem_code" valid:"required"`
	Name                   string                     `json:"name" valid:"required"`
	Type                   int8                       `json:"type" valid:"required"`
	StartTime              time.Time                  `json:"start_time" valid:"required"`
	EndTime                time.Time                  `json:"end_time" valid:"required"`
	OverallQuota           int64                      `json:"overall_quota" valid:"gt:0"`
	UserQuota              int64                      `json:"user_quota" valid:"gt:0"`
	MinOrder               string                     `json:"min_order" valid:"required"`
	DiscAmount             float64                    `json:"disc_amount" valid:"gt:0"`
	TermConditions         string                     `json:"term_conditions" valid:"required"`
	ImageUrl               string                     `json:"image_url" valid:"required"`
	VoidReason             int8                       `json:"void_reason"`
	Note                   string                     `json:"note"`
	Status                 int8                       `json:"status"`
	CreatedAt              time.Time                  `json:"created_at"`
	VoucherItem            []VoucherItemCreateRequest `json:"voucher_item"`
}
type ArchetypeResponse struct {
	ID             string                `json:"id"`
	Code           string                `json:"code"`
	Description    string                `json:"description"`
	CustomerTypeID string                `json:"customer_type_id"`
	Status         int8                  `json:"status"`
	ConvertStatus  string                `json:"convert_status"`
	CustomerType   *CustomerTypeResponse `json:"customer_type"`
}

type RegionResponse struct {
	ID          string `json:"id"`
	Code        string `json:"code"`
	Description string `json:"description"`
}

type MembershipLevelResponse struct {
	ID       int64  `json:"id"`
	Code     string `json:"code"`
	Level    int8   `json:"level"`
	Name     string `json:"name"`
	ImageUrl string `json:"image_url"`
	Status   int8   `json:"status"`
}

type MembershipCheckpointResponse struct {
	ID                int64   `json:"id"`
	Checkpoint        int8    `json:"checkpoint"`
	TargetAmount      float64 `json:"target_amount"`
	Status            int8    `json:"status"`
	MembershipLevelID int64   `json:"membership_level_id"`
}

type DivisionResponse struct {
	ID            int64     `json:"id,omitempty"`
	Code          string    `json:"code,omitempty"`
	Description   string    `json:"description,omitempty"`
	Status        int8      `json:"status,omitempty"`
	CreatedAt     time.Time `json:"created_at,omitempty"`
	UpdatedAt     time.Time `json:"updated_at,omitempty"`
	StatusConvert string    `json:"status_convert"`
	Note          string    `json:"note"`
}

type CustomerResponse struct {
	ID   int64  `json:"id"`
	Code string `json:"code"`
	Name string `json:"name"`
}

type CustomerTypeResponse struct {
	ID            string `json:"id"`
	Code          string `json:"code"`
	Description   string `json:"description"`
	CustomerGroup string `json:"customer_group"`
	Status        int8   `json:"status"`
	ConvertStatus string `json:"convert_status"`
}

type VoucherRequestBulky struct {
	Data []*VoucherRequestCreateBulky `json:"data" valid:"required"`
}

type VoucherRequestCreateBulky struct {
	RegionName             string    `json:"region_name"`
	CustomerCode           string    `json:"customer_code"`
	CustomerTypeCode       string    `json:"customer_type_code"`
	ArchetypeCode          string    `json:"archetype_code"`
	RedeemCode             string    `json:"redeem_code"`
	VoucherName            string    `json:"voucher_name"`
	VoucherType            int8      `json:"voucher_type" `
	StartTime              string    `json:"start_time" `
	EndTime                string    `json:"end_time"`
	OverallQuota           int64     `json:"overall_quota"`
	UserQuota              int64     `json:"user_quota"`
	DiscAmount             float64   `json:"disc_amount" `
	MinOrder               float64   `json:"min_order"`
	Note                   string    `json:"note"`
	MembershipLevel        int8      `json:"membership_level"`
	MembershipCheckpoint   int8      `json:"membership_checkpoint"`
	StartTimeActual        time.Time `json:"-"`
	EndTimeActual          time.Time `json:"-"`
	RegionID               string    `json:"-"`
	CustomerID             int64     `json:"-"`
	ArchetypeID            string    `json:"-"`
	DivisionID             int64     `json:"-"`
	MembershipLevelID      int64     `json:"-"`
	MembershipCheckpointID int64     `json:"-"`
	ExpenseAccount         string    `json:"-"`
}

type UomResponse struct {
	ID            string `json:"id"`
	Code          string `json:"code"`
	Description   string `json:"description"`
	Status        int8   `json:"status"`
	StatusConvert string `json:"status_convert"`
}

type VoucherRequestGetMobileVoucherList struct {
	RegionID             string `json:"region_id"`
	CustomerTypeID       string `json:"customer_type_id"`
	CustomerID           int64  `json:"customer_id"`
	ArchetypeID          string `json:"archetype_id"`
	CustomerLevel        int8   `json:"customer_level"`
	MembershipLevel      int8   `json:"membership_level"`
	MembershipCheckpoint int8   `json:"membership_checkpoint"`
	IsMembershipOnly     bool   `json:"membership_only"`
	Offset               int64  `json:"offset"`
	Limit                int64  `json:"limit"`
	Category             int8   `json:"category"`
}

type VoucherRequestGetMobileVoucherDetail struct {
	RedeemCode string `json:"redeem_code"`
	CustomerID int64  `json:"customer_id"`
	Status     int8   `json:"status"`
	Code       string `json:"code"`
}

type VoucherRequestUpdate struct {
	VoucherID       int64 `json:"voucher_id"`
	RemOverallQuota int64 `json:"rem_overall_quota"`
}
