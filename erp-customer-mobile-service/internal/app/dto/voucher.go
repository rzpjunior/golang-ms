package dto

type VoucherResponse struct {
	ID                     int64   `json:"id"`
	VoucherName            string  `json:"voucher_name,omitempty"`
	RedeemCode             string  `json:"redeem_code"`
	ImageUrl               string  `json:"image_url,omitempty"`
	MinOrder               float64 `json:"min_order,omitempty"`
	EndTime                string  `json:"end_time,omitempty"`
	RemUserQuota           int64   `json:"rem_user_quota,omitempty"`
	TermCondition          string  `json:"term_condition,omitempty"`
	VoucherItem            int64   `json:"voucher_item,omitempty"`
	DiscAmount             float64 `json:"disc_amount,omitempty"`
	Type                   int8    `json:"type,omitempty"`
	MembershipLevelID      int64   `json:"membership_level_id,omitempty"`
	MembershipCheckpointID int64   `json:"membership_checkpoint_id,omitempty"`

	MembershipLevel      *MembershipLevelResponse      `json:"membership_level,omitempty"`
	MembershipCheckpoint *MembershipCheckpointResponse `json:"membership_checkpoint,omitempty"`
}

// Request Get Voucher List
type VoucherRequestGet struct {
	Platform string                 `json:"platform" valid:"required"`
	Data     *VoucherRequestGetList `json:"data" valid:"required"`
	Limit    int64                  `json:"limit"`
	Offset   int64                  `json:"offset"`

	Session *SessionDataCustomer
}

type VoucherRequestGetList struct {
	RegionID             string `json:"region_id" valid:"required"`
	CustomerTypeID       string `json:"customer_type_id" valid:"required"`
	ArchetypeID          string `json:"archetype_id" valid:"required"`
	MembershipLevel      string `json:"membership_level"`
	MembershipCheckpoint string `json:"membership_checkpoint"`
	IsMembershipOnly     bool   `json:"is_membership_only"`
}

// Request Get Voucher Detail
type VoucherRequestGetDetail struct {
	Platform string            `json:"platform" valid:"required"`
	Data     *VoucherGetDetail `json:"data" valid:"required"`

	Session *SessionDataCustomer
}

type VoucherGetDetail struct {
	RedeemCode string `json:"redeem_code" valid:"required"`
}

// Request Apply Voucher
type VoucherRequestApply struct {
	Platform string        `json:"platform" valid:"required"`
	Data     *VoucherApply `json:"data" valid:"required"`

	Session *SessionDataCustomer
}

type VoucherApply struct {
	RedeemCode   string         `json:"redeem_code" valid:"required"`
	TotalPrice   float64        `json:"total_price" valid:"required"`
	TotalCharge  float64        `json:"total_charge" valid:"required"`
	RegionID     string         `json:"region_id" valid:"required"`
	AddressID    string         `json:"address_id" valid:"required"`
	VoucherItems []*VoucherItem `json:"voucher_items"`
}

type VoucherItem struct {
	ItemID   string  `json:"item_id"`
	OrderQty float64 `json:"order_qty"`
}

// Get List Voucher Item
type VoucherRequestGetItemList struct {
	Platform string          `json:"platform" valid:"required"`
	Data     *VoucherGetItem `json:"data" valid:"required"`

	Session *SessionDataCustomer
}

type VoucherGetItem struct {
	RedeemCode string `json:"redeem_code" valid:"required"`
}

type VoucherGetItemResponse struct {
	ItemID     int64   `json:"item_id"`
	ItemName   string  `json:"item_name"`
	MinQtyDisc float64 `json:"min_qty_disc"`
	UomName    string  `json:"uom_name"`
}
