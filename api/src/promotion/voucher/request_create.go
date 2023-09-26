// Copyright 2020 PT. Eden Pangan Indonesia. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

package voucher

import (
	"encoding/json"
	"sort"
	"strconv"
	"strings"
	"time"

	"git.edenfarm.id/cuxs/common"
	"git.edenfarm.id/cuxs/validation"
	"git.edenfarm.id/project-version2/api/datastore/repository"
	"git.edenfarm.id/project-version2/api/src/auth"
	"git.edenfarm.id/project-version2/api/util"
	"git.edenfarm.id/project-version2/datamodel/model"
)

// createRequest : struct to hold request data
type createRequest struct {
	Code                   string   `json:"-"`
	RedeemCode             string   `json:"redeem_code" valid:"required"`
	Name                   string   `json:"name" valid:"required"`
	Type                   int8     `json:"type" valid:"required"`
	AreaID                 string   `json:"area_id" valid:"required"`
	BusinessTypeID         string   `json:"business_type_id" valid:"required"`
	ArcheTypeID            string   `json:"archetype_id" valid:"required"`
	CustomerTag            []string `json:"customer_tag"`
	OverallQuota           string   `json:"overall_quota" valid:"required"`
	UserQuota              string   `json:"user_quota" valid:"required"`
	MinOrder               string   `json:"min_order" valid:"required"`
	DiscAmount             string   `json:"disc_amount" valid:"required"`
	StartTimestampStr      string   `json:"start_timestamp" valid:"required"`
	EndTimestampStr        string   `json:"end_timestamp" valid:"required"`
	Note                   string   `json:"note"`
	CustomerTagStr         string   `json:"-"`
	UserQuotaInt           int64    `json:"-"`
	OverallQuotaInt        int64    `json:"-"`
	MinOrderFloat          float64  `json:"-"`
	DiscAmountFloat        float64  `json:"-"`
	MembershipLevelID      string   `json:"membership_level_id"`
	MembershipCheckpointID string   `json:"membership_checkpoint_id"`

	ChannelVoucher []string       `json:"channel_voucher"`
	ImageUrl       string         `json:"image_url"`
	TermConditions string         `json:"term_conditions"`
	MerchantID     string         `json:"merchant_id"`
	VoucherItem    []*voucherItem `json:"voucher_item"`
	IsMobile       bool

	Merchant             *model.Merchant
	Area                 *model.Area
	ArcheType            *model.Archetype
	StartTimestamp       time.Time
	EndTimestamp         time.Time
	Session              *auth.SessionData
	HasVoucherItem       int8
	MembershipLevel      *model.MembershipLevel
	MembershipCheckpoint *model.MembershipCheckpoint
}

type voucherItem struct {
	ProductID  string  `json:"product_id"`
	MinQtyDisc float64 `json:"min_qty_disc"`

	Product *model.Product `json:"-"`
}

// Validate : function to validate request data
func (c *createRequest) Validate() *validation.Output {
	o := &validation.Output{Valid: true}
	var (
		err                                                       error
		arrCustomerTagInt                                         []int
		productID, mID, membershipLevelID, membershipCheckpointID int64
	)

	if c.Code, err = util.CheckTable("voucher"); err != nil {
		o.Failure("code.invalid", util.ErrorInvalidData("code voucher"))
	}

	if len(c.RedeemCode) < 5 || len(c.RedeemCode) > 25 {
		o.Failure("redeem_code.invalid", util.ErrorRangeChar("redeem code", "5", "25"))
	}

	UserQuotaInt, _ := strconv.ParseInt(c.UserQuota, 10, 64)
	OverallQuotaInt, _ := strconv.ParseInt(c.OverallQuota, 10, 64)
	MinOrderFloat, _ := strconv.ParseFloat(c.MinOrder, 64)
	DiscAmountFloat, _ := strconv.ParseFloat(c.DiscAmount, 64)

	filter := map[string]interface{}{"redeem_code": c.RedeemCode}
	exclude := map[string]interface{}{"status": int8(2)}
	if _, countVoucher, err := repository.CheckVoucherData(filter, exclude); err != nil {
		o.Failure("redeem_code.invalid", util.ErrorInvalidData("redeem code"))
	} else if countVoucher > 0 {
		o.Failure("redeem_code.duplicate", util.ErrorUnique("redeem code"))
	}

	if c.StartTimestampStr != "" {
		if c.StartTimestamp, err = time.Parse(time.RFC3339, c.StartTimestampStr); err != nil {
			o.Failure("start_timestamp.invalid", util.ErrorInvalidData("start timestamp"))
		}
	}

	if c.EndTimestampStr != "" {
		if c.EndTimestamp, err = time.Parse(time.RFC3339, c.EndTimestampStr); err != nil {
			o.Failure("end_timestamp.invalid", util.ErrorInvalidData("end timestamp"))
		}
	}

	if c.StartTimestamp.Equal(c.EndTimestamp) || c.EndTimestamp.Before(c.StartTimestamp) {
		o.Failure("start_timestamp.invalid", util.ErrorLater("start timestamp", "end timestamp"))
	}

	if areaID, err := common.Decrypt(c.AreaID); err != nil {
		o.Failure("area_id.invalid", util.ErrorInvalidData("area"))
	} else {
		if c.Area, err = repository.ValidArea(areaID); err != nil {
			o.Failure("area_id.invalid", util.ErrorInvalidData("area"))
		} else {
			if c.Area.Status != int8(1) {
				o.Failure("area_id.invalid", util.ErrorActive("area"))
			}
		}
	}

	if archeTypeID, err := common.Decrypt(c.ArcheTypeID); err != nil {
		o.Failure("archetype_id.invalid", util.ErrorInvalidData("archetype"))
	} else {
		if c.ArcheType, err = repository.ValidArchetype(archeTypeID); err != nil {
			o.Failure("archetype_id.invalid", util.ErrorInvalidData("archetype"))
		} else {
			if c.ArcheType.Status != int8(1) {
				o.Failure("archetype_id.invalid", util.ErrorActive("archetype"))
			}
		}
	}

	if MinOrderFloat < 0 {
		o.Failure("min_order.invalid", util.ErrorGreater("min order", "0"))
	}

	if DiscAmountFloat < 1 {
		o.Failure("disc_amount.invalid", util.ErrorEqualGreater("discount amount", "1"))
	}

	if UserQuotaInt < 1 {
		o.Failure("user_quota.invalid", util.ErrorEqualGreater("user quota", "1"))
	}

	if OverallQuotaInt <= 0 {
		o.Failure("overall_quota.invalid", util.ErrorGreater("overall quota", "1"))
	}

	if OverallQuotaInt < UserQuotaInt {
		o.Failure("overall_quota.invalid", util.ErrorEqualGreater("overall quota", "user quota"))
	}

	c.UserQuotaInt = UserQuotaInt
	c.OverallQuotaInt = OverallQuotaInt
	c.MinOrderFloat = MinOrderFloat
	c.DiscAmountFloat = DiscAmountFloat

	configApp, err := repository.GetConfigApp("attribute", "vou_max_tag")
	maxValue, err := strconv.Atoi(configApp.Value)

	if len(c.CustomerTag) > maxValue {
		o.Failure("customer_tag.invalid", util.ErrorSelectMax(configApp.Value, "tag"))
	}

	if c.MerchantID != "" {
		if mID, err = common.Decrypt(c.MerchantID); err == nil {
			c.Merchant = &model.Merchant{ID: mID}
			c.Merchant.Read("ID")
		}
	}

	if len(c.CustomerTag) > 0 {
		for _, v := range c.CustomerTag {
			customerId, _ := common.Decrypt(v)

			if customerTag, err := repository.ValidCustomerTag(customerId); err != nil {
				o.Failure("customer_tag.invalid", util.ErrorInvalidData("customer tag"))
			} else {
				if customerTag.Status != int8(1) {
					o.Failure("customer_tag.active", util.ErrorActive("customer tag"))
				}

				arrCustomerTagInt = append(arrCustomerTagInt, int(customerId))
			}
		}

		// sort integer decrypted customer tag id, then convert it into a string with comma separator
		sort.Ints(arrCustomerTagInt)
		custTagJson, _ := json.Marshal(arrCustomerTagInt)
		c.CustomerTagStr = strings.Trim(string(custTagJson), "[]")
	}

	for _, v := range c.ChannelVoucher {
		if v == "2" || v == "3" {
			c.IsMobile = true
			break
		}
	}

	if c.IsMobile {
		if c.TermConditions == "" {
			o.Failure("term_conditions.required", util.ErrorInputRequired("term and conditions"))
		}
		if c.ImageUrl == "" {
			o.Failure("image_url.required", util.ErrorInputRequired("image"))
		}
	}

	productList := make(map[int64]string)
	if len(c.VoucherItem) > 0 {
		c.HasVoucherItem = 1
		for i, v := range c.VoucherItem {
			if v.MinQtyDisc < 0.01 {
				o.Failure("qty"+strconv.Itoa(i)+".equalorgreater", util.ErrorGreater("minimal qty", "0"))
			}

			if v.ProductID == "" {
				o.Failure("product_id"+strconv.Itoa(i)+".invalid", util.ErrorInputRequired("product"))
			} else {
				if productID, err = common.Decrypt(v.ProductID); err == nil {
					if v.Product, err = repository.ValidProduct(productID); err == nil {

						if _, exist := productList[productID]; exist {
							o.Failure("product_id"+strconv.Itoa(i)+".duplicate", util.ErrorDuplicate("product"))
						} else {
							productList[productID] = "t"
						}
					} else {
						o.Failure("product_id"+strconv.Itoa(i)+".invalid", util.ErrorInvalidData("product"))
					}
				}
			}
		}
	} else {
		c.HasVoucherItem = 2 // if dont have voucher item
	}

	// start check membership request
	if c.MembershipLevelID != "" {
		// decrypt membership level id
		if membershipLevelID, err = common.Decrypt(c.MembershipLevelID); err != nil {
			o.Failure("membership_level_id.invalid", util.ErrorInvalidData("membership level"))
			return o
		}

		// get membership level data
		if c.MembershipLevel, err = repository.ValidMembershipLevel(membershipLevelID); err != nil {
			o.Failure("membership_level_id.invalid", util.ErrorInvalidData("membership level"))
			return o
		}

		// if membership level not empty then check membership checkpoint
		if membershipLevelID > 0 || c.MembershipLevel != nil {
			// check membership checkpoint id
			if c.MembershipCheckpointID == "" {
				o.Failure("membership_checkpoint_id.required", util.ErrorSelectRequired("membership checkpoint"))
				return o
			}

			// decrypt membership checkpoint id
			if membershipCheckpointID, err = common.Decrypt(c.MembershipCheckpointID); err != nil {
				o.Failure("membership_checkpoint_id.invalid", util.ErrorInvalidData("membership checkpoint"))
				return o
			}

			// get membership checkpoint data
			if c.MembershipCheckpoint, err = repository.ValidMembershipCheckpoint(membershipCheckpointID, membershipLevelID); err != nil {
				o.Failure("membership_checkpoint_id.invalid", util.ErrorInvalidData("membership checkpoint"))
				return o
			}
		}
	}
	// end check membership request

	return o
}

// Messages : function to return error validation messages
func (c *createRequest) Messages() map[string]string {
	return map[string]string{
		"redeem_code.required":      util.ErrorInputRequired("redeem code"),
		"name.required":             util.ErrorInputRequired("name"),
		"type.required":             util.ErrorSelectRequired("type"),
		"area_id.required":          util.ErrorSelectRequired("area"),
		"business_type_id.required": util.ErrorSelectRequired("business type"),
		"archetype_id.required":     util.ErrorSelectRequired("archetype"),
		"overall_quota.required":    util.ErrorInputRequired("overall quota"),
		"user_quota.required":       util.ErrorInputRequired("user quota"),
		"min_order.required":        util.ErrorInputRequired("minimum order"),
		"disc_amount.required":      util.ErrorInputRequired("discount amount"),
		"start_timestamp.required":  util.ErrorSelectRequired("start date & start time"),
		"end_timestamp.required":    util.ErrorSelectRequired("end date & end time"),
		"customer_tag.required":     util.ErrorSelectRequired("customer tag"),
	}
}
