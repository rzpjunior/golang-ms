// Copyright 2020 PT. Eden Pangan Indonesia. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

package agent

import (
	"encoding/json"
	"sort"
	"strconv"
	"strings"

	"git.edenfarm.id/project-version2/datamodel/model"

	"git.edenfarm.id/cuxs/common"
	"git.edenfarm.id/cuxs/validation"
	"git.edenfarm.id/project-version2/api/datastore/repository"
	"git.edenfarm.id/project-version2/api/src/auth"
	"git.edenfarm.id/project-version2/api/util"
)

type updateTagRequest struct {
	ID             int64    `json:"-" valid:"required"`
	CustomerTag    []string `json:"tag_customer" valid:"required"`
	CustomerTagStr string   `json:"-"`

	Merchant *model.Merchant

	Session *auth.SessionData `json:"-"`
}

// Validate : function to validate supplier request data
func (c *updateTagRequest) Validate() *validation.Output {
	o := &validation.Output{Valid: true}
	var err error
	var arrCustomerTagInt []int

	confApp := &model.ConfigApp{
		Attribute: "cust_max_tag",
	}
	confApp.Read("Attribute")
	configAppValue, _ := strconv.Atoi(confApp.Value)
	if len(c.CustomerTag) > configAppValue {
		o.Failure("tag_customer.invalid", util.ErrorSelectMax("3", "product tag"))
	} else {
		for _, v := range c.CustomerTag {
			customerId, _ := common.Decrypt(v)

			if customerTag, err := repository.ValidCustomerTag(customerId); err != nil {
				o.Failure("tag_customer.invalid", util.ErrorInvalidData("customer tag"))
			} else {
				if customerTag.Status != int8(1) {
					o.Failure("tag_customer.active", util.ErrorActive("customer tag"))
				}

				arrCustomerTagInt = append(arrCustomerTagInt, int(customerId))
			}
		}

		if len(arrCustomerTagInt) > 0 {
			// sort integer decrypted customer tag id, then convert it into a string with comma separator
			sort.Ints(arrCustomerTagInt)
			customerTagJson, _ := json.Marshal(arrCustomerTagInt)
			c.CustomerTagStr = strings.Trim(string(customerTagJson), "[]")
		} else {
			c.CustomerTagStr = ""
		}
	}

	if c.Merchant, err = repository.ValidMerchant(c.ID); err == nil {
		if c.Merchant.Status != 1 {
			o.Failure("status.active", util.ErrorActive("status"))
		}
	} else {
		o.Failure("id.invalid", util.ErrorInvalidData("agent"))
	}

	return o
}

// Messages : function to return error validation messages
func (c *updateTagRequest) Messages() map[string]string {
	return map[string]string{
		"tag_customer.required": util.ErrorInputRequired("tag customer"),
	}
}
