// Copyright 2020 PT. Eden Pangan Indonesia. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

package merchant

import (
	"encoding/json"
	"sort"
	"strconv"
	"strings"

	"git.edenfarm.id/cuxs/common"
	"git.edenfarm.id/cuxs/validation"
	"git.edenfarm.id/project-version2/datamodel/model"
	"git.edenfarm.id/project-version2/api/datastore/repository"
	"git.edenfarm.id/project-version2/api/src/auth"
	"git.edenfarm.id/project-version2/api/util"
)

type updateTagRequest struct {
	ID             int64    `json:"-" valid:"required"`
	CustomerTag    []string `json:"customer_tag" valid:"required"`
	CustomerTagStr string   `json:"-"`

	Merchant *model.Merchant

	Session *auth.SessionData `json:"-"`
}

// Validate : function to validate merchant request data
func (c *updateTagRequest) Validate() *validation.Output {
	o := &validation.Output{Valid: true}
	var err error
	var arrCustomerTagInt []int

	configApp, err := repository.GetConfigApp("attribute", "cust_max_tag")
	maxValue, err := strconv.Atoi(configApp.Value)
	if len(c.CustomerTag) > maxValue {
		o.Failure("customer_tag.invalid", util.ErrorSelectMax(configApp.Value, "tag"))
	}

	for _, v := range c.CustomerTag {
		customerId, _ := common.Decrypt(v)

		if customerTag, err := repository.ValidCustomerTag(customerId); err != nil {
			o.Failure("tag_customer.invalid", util.ErrorInvalidData("tag customer"))
		} else {
			if customerTag.Status != int8(1) {
				o.Failure("tag_customer.active", util.ErrorActive("tag customer"))
			}

			arrCustomerTagInt = append(arrCustomerTagInt, int(customerId))
		}
	}

	// sort integer decrypted customer tag id, then convert it into a string with comma separator
	sort.Ints(arrCustomerTagInt)
	custTagJson, _ := json.Marshal(arrCustomerTagInt)
	c.CustomerTagStr = strings.Trim(string(custTagJson), "[]")

	if c.Merchant, err = repository.ValidMerchant(c.ID); err == nil {
		if c.Merchant.Status != 1 {
			o.Failure("status.active", util.ErrorActive("status"))
		}
	} else {
		o.Failure("merchant.invalid", util.ErrorInvalidData("merchant"))
	}

	return o
}

// Messages : function to return error validation messages
func (c *updateTagRequest) Messages() map[string]string {
	return map[string]string{
		"customer_tag.required": util.ErrorInputRequired("customer_tag"),
	}
}
