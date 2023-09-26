package dto

type CreateVoucherGPRequest struct {
	InterID            string              `json:"interid"`
	GnlVoucherID       string              `json:"gnl_voucher_id"`
	GnlChannel         int32               `json:"gnl_channel"`
	GnlVoucherType     int32               `json:"gnl_voucher_type"`
	GnlVoucherName     string              `json:"gnl_voucher_name"`
	GnlExpenseAccount  string              `json:"gnl_expense_account"`
	GnlVoucherCode     string              `json:"gnl_voucher_code"`
	GnlMinimumOrder    int32               `json:"gnl_minimum_order"`
	GnlDiscountAmount  int32               `json:"gnl_discount_amount"`
	GnlVoucherStatus   int32               `json:"gnl_voucher_status"`
	Inactive           int32               `json:"inactive"`
	Restriction        *Restriction        `json:"restriction"`
	AdvancedProperties *AdvancedProperties `json:"advanced_properties"`
}

type Restriction struct {
	GnlRegion      string `json:"gnl_region"`
	GnlCustTypeID  string `json:"gnl_cust_type_id"`
	GnlArchetypeID string `json:"gnl_archetype_id"`
	DefaultCB      int32  `json:"default_cb"`
}

type AdvancedProperties struct {
	Custnmbr              string `json:"custnmbr"`
	GnlStartPeriod        string `json:"gnl_start_period"`
	GnlEndPeriod          string `json:"gnl_end_period"`
	GnlTotalQuotaCount    int32  `json:"gnl_total_quota_count"`
	GnlTotalQuotaCountPE  int32  `json:"gnl_total_quota_count_pe"`
	GnlRemainingOverallQu int32  `json:"gnl_remaining_overall_qu"`
	GnlMobileVoucher      int32  `json:"gnl_mobile_voucher"`
}
