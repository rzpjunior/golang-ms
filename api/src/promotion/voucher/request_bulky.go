package voucher

import (
	"fmt"
	"strconv"
	"strings"
	"time"

	"git.edenfarm.id/cuxs/orm"
	"git.edenfarm.id/cuxs/validation"
	"git.edenfarm.id/project-version2/api/datastore/repository"
	"git.edenfarm.id/project-version2/api/src/auth"
	"git.edenfarm.id/project-version2/api/util"
	"git.edenfarm.id/project-version2/datamodel/model"
)

type bulkyRequest struct {
	Sheet         []*sheets `json:"data" valid:"required"`
	ErrorCallback string    `json:"error_callback"`
	Session       *auth.SessionData
}

type sheets struct {
	Code                    string  `json:"-"`
	AreaCode                string  `json:"area_code"`
	MerchantCode            string  `json:"merchant_code"`
	ArchetypeCode           string  `json:"archetype_code"`
	CustomerTagCode         string  `json:"customer_tag_code"`
	RedeemCode              string  `json:"redeem_code"`
	VoucherName             string  `json:"voucher_name"`
	VoucherType             int8    `json:"voucher_type" `
	StartTimestampString    string  `json:"start_timestamp" `
	EndTimestampString      string  `json:"end_timestamp" `
	OverallQuota            int64   `json:"overall_quota"`
	UserQuota               int64   `json:"user_quota"`
	DiscountAmount          float64 `json:"disc_amount" `
	MinOrder                float64 `json:"min_order"`
	Note                    string  `json:"note"`
	Area                    *model.Area
	Merchant                *model.Merchant
	Archetype               *model.Archetype
	StartTimestamp          time.Time `json:"-"`
	StartTimestampConverted string
	EndTimestamp            time.Time
	EndTimestampConverted   string
	CustomerTagID           string
	MembershipLevel         string `json:"membership_level"`
	MembershipCheckpoint    string `json:"membership_checkpoint"`

	MembershipLevelModel      *model.MembershipLevel      `json:"-"`
	MembershipCheckpointModel *model.MembershipCheckpoint `json:"-"`
}

// Validate : function to validate request data
func (c *bulkyRequest) Validate() *validation.Output {
	o := &validation.Output{Valid: true}
	om := orm.NewOrm()
	om.Using("read_only")
	layout := "2006-01-02 15:04:05"
	loc, _ := time.LoadLocation("Asia/Jakarta")
	redeemCodeList := make(map[string]string)
	configApp, _ := repository.GetConfigApp("attribute", "vou_max_tag")
	var (
		errorCallback []string
		e             error
	)
	for index, row := range c.Sheet {
		if row.VoucherName == "" {
			errorCallback = append(errorCallback, util.ErrorInputRequired("voucher name at row "+strconv.Itoa(index+1)))
		}

		if row.VoucherType == 0 {
			errorCallback = append(errorCallback, util.ErrorInputRequired("voucher type at row "+strconv.Itoa(index+1)))
		}

		if row.RedeemCode != "" {
			if _, exist := redeemCodeList[row.RedeemCode]; exist {
				errorCallback = append(errorCallback, errorDuplicateForVoucherBulkyFromRequest("redeem code", strconv.Itoa(index+1)))
			} else {
				redeemCodeList[row.RedeemCode] = "t"
			}

			filter := map[string]interface{}{"redeem_code": row.RedeemCode, "status": int8(1)}
			exclude := map[string]interface{}{}
			if _, countVoucher, err := repository.CheckVoucherData(filter, exclude); err != nil {
				errorCallback = append(errorCallback, util.ErrorInvalidData("redeem code at row "+strconv.Itoa(index+1)))
			} else if countVoucher > 0 {
				errorCallback = append(errorCallback, errorDuplicateForVoucherBulkyFromData("redeem code", strconv.Itoa(index+1)))
			}

			if len(row.RedeemCode) < 5 || len(row.RedeemCode) > 25 {
				errorCallback = append(errorCallback, "redeem code must contains 5-25 at row "+strconv.Itoa(index+1)+".")
			}
		} else {
			errorCallback = append(errorCallback, util.ErrorInputRequired("redeem code at row "+strconv.Itoa(index+1)))
		}

		if row.Code, e = util.CheckTable("voucher"); e != nil {
			o.Failure("code_row"+strconv.Itoa(index+2)+".invalid", util.ErrorInvalidData("code voucher at row "+strconv.Itoa(index+1)))
		}

		if row.AreaCode != "" {
			row.Area = &model.Area{Code: row.AreaCode}
			if e = row.Area.Read("Code"); e != nil {
				errorCallback = append(errorCallback, util.ErrorInvalidData("area code at row "+strconv.Itoa(index+1)))
			}
		} else {
			errorCallback = append(errorCallback, util.ErrorInputRequired("area code at row "+strconv.Itoa(index+1)))
		}

		if row.MerchantCode != "" {
			row.Merchant = &model.Merchant{Code: row.MerchantCode}
			if e = row.Merchant.Read("Code"); e != nil {
				errorCallback = append(errorCallback, util.ErrorInvalidData("merchant code at row "+strconv.Itoa(index+1)))
			}
		}
		if row.ArchetypeCode != "" {
			row.Archetype = &model.Archetype{Code: row.ArchetypeCode}
			if e = row.Archetype.Read("Code"); e != nil {
				errorCallback = append(errorCallback, util.ErrorInvalidData("archetype code at row "+strconv.Itoa(index+1)))
			}
		} else {
			errorCallback = append(errorCallback, util.ErrorInputRequired("archetype code at row "+strconv.Itoa(index+1)))
		}

		if row.StartTimestampString != "" {
			if row.StartTimestamp, e = time.ParseInLocation(layout, row.StartTimestampString, loc); e != nil {
				errorCallback = append(errorCallback, util.ErrorInvalidData("start timestamp at row "+strconv.Itoa(index+1)))
			} else {
				row.StartTimestampConverted = row.StartTimestamp.Format(time.RFC3339)
				row.StartTimestamp, e = time.ParseInLocation(time.RFC3339, row.StartTimestampConverted, loc)
			}
		} else {
			errorCallback = append(errorCallback, util.ErrorInputRequired("start timestamp at row "+strconv.Itoa(index+1)))
		}

		if row.EndTimestampString != "" {
			if row.EndTimestamp, e = time.ParseInLocation(layout, row.EndTimestampString, loc); e != nil {
				errorCallback = append(errorCallback, util.ErrorInvalidData("end timestamp at row "+strconv.Itoa(index+1)))
			} else {
				row.EndTimestampConverted = row.EndTimestamp.Format(time.RFC3339)
				row.EndTimestamp, e = time.ParseInLocation(time.RFC3339, row.EndTimestampConverted, loc)
			}
		} else {
			errorCallback = append(errorCallback, util.ErrorInputRequired("end timestamp at row "+strconv.Itoa(index+1)))
		}

		if row.StartTimestamp.Equal(row.EndTimestamp) || row.EndTimestamp.Before(row.StartTimestamp) {
			errorCallback = append(errorCallback, util.ErrorLater("start timestamp", "end timestamp at row "+strconv.Itoa(index+1)))
		}

		if row.MinOrder < 0 {
			errorCallback = append(errorCallback, util.ErrorGreater("min order", "0 at row "+strconv.Itoa(index+1)))
		}

		if row.DiscountAmount < 1 {
			errorCallback = append(errorCallback, util.ErrorEqualGreater("discount amount", "1 at row "+strconv.Itoa(index+1)))
		}

		if row.UserQuota < 1 {
			errorCallback = append(errorCallback, util.ErrorEqualGreater("user quota", "1 at row "+strconv.Itoa(index+1)))
		}

		if row.OverallQuota <= 0 {
			errorCallback = append(errorCallback, util.ErrorGreater("overall quota", "1 at row "+strconv.Itoa(index+1)))
		}

		if row.OverallQuota < row.UserQuota {
			errorCallback = append(errorCallback, util.ErrorEqualGreater("overall quota", "user quota at row "+strconv.Itoa(index+1)))
		}

		if len(row.CustomerTagCode) > 0 {
			var customerTags []*model.TagCustomer
			maxValue, _ := strconv.Atoi(configApp.Value)
			customerTagArray := strings.Split(row.CustomerTagCode, ",")

			if len(customerTagArray) > maxValue {
				errorCallback = append(errorCallback, util.ErrorSelectMax(configApp.Value, "tag at row "+strconv.Itoa(index+1)))
			}
			customerTagNotDuplicate := util.RemoveDuplicateValuesString(customerTagArray)

			if len(customerTagNotDuplicate) > 0 {
				var cat []string
				var tags []int64
				var total int
				for i := 0; i < len(customerTagNotDuplicate); i++ {
					cat = append(cat, "?")
				}
				catLength := strings.Join(cat, ",")
				om.Raw("SELECT COUNT(id) FROM tag_customer WHERE code IN ("+catLength+")", customerTagNotDuplicate).QueryRow(&total)
				om.Raw("SELECT * FROM tag_customer WHERE code IN ("+catLength+")", customerTagNotDuplicate).QueryRows(&customerTags)
				if total != len(customerTagNotDuplicate) {
					errorCallback = append(errorCallback, util.ErrorInvalidData("tag at row "+strconv.Itoa(index+1)))
				}
				for _, tag := range customerTags {
					tags = append(tags, tag.ID)
				}
				row.CustomerTagID = strings.Trim(strings.Replace(fmt.Sprint(tags), " ", ",", -1), "[]")
				//append(arrCustomerTag, customerTagToInt...) // the dots (customerTagToInt...) mean just like any other looping function  ex: for i := 0; i < len(is); i++ {  fmt.Println(is[i]) }

			}

		}

		// start check membership request
		if row.MembershipLevel != "" {
			// get membership level data
			if row.MembershipLevelModel, e = repository.GetMembershipLevel("level", row.MembershipLevel); e != nil {
				errorCallback = append(errorCallback, util.ErrorInvalidData("membership level at row "+strconv.Itoa(index+1)))
				continue
			}

			// check membership checkpoint id
			if row.MembershipCheckpoint == "" {
				errorCallback = append(errorCallback, util.ErrorSelectRequired("membership checkpoint at row "+strconv.Itoa(index+1)))
				continue
			}

			// check if checkpoint is a number
			if _, e = strconv.Atoi(row.MembershipCheckpoint); e != nil {
				errorCallback = append(errorCallback, util.ErrorInvalidData("membership checkpoint at row "+strconv.Itoa(index+1)))
				continue
			}

			// get membership checkpoint data
			if row.MembershipCheckpointModel, e = repository.GetMembershipCheckpoint("checkpoint", row.MembershipCheckpoint); e != nil {
				errorCallback = append(errorCallback, util.ErrorInvalidData("membership checkpoint at row "+strconv.Itoa(index+1)))
				continue
			}
		}
		// end check membership request
	}

	if len(errorCallback) > 0 {
		c.ErrorCallback = strings.Trim(strings.Replace(fmt.Sprint(errorCallback), ".", ".|", -1), "[]")
		o.Failure("error_callback.invalid", c.ErrorCallback)
	}

	return o
}

// Messages : function to return error validation messages
func (c *bulkyRequest) Messages() map[string]string {
	return map[string]string{}
}

func errorDuplicateForVoucherBulkyFromRequest(fieldName, n string) string {
	return fieldName + " is duplicate in the file at row " + n + ", please enter another " + fieldName + "."
}

func errorDuplicateForVoucherBulkyFromData(fieldName, n string) string {
	return fieldName + " is duplicate with existing data at row" + n + ", please enter another " + fieldName + "."
}
