// Copyright 2020 PT. Eden Pangan Indonesia. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

package datascraping

import (
	"git.edenfarm.id/cuxs/validation"
	"git.edenfarm.id/project-version2/api/src/auth"
)

// matchingRequest : struct to hold matching product data from dashboard
type matchingRequest struct {
	Data []*itemsData `json:"data"`

	Session *auth.SessionData `json:"-"`
}

type itemsData struct {
	EdenProductCode string `json:"eden_product_code"`
	EdenProductName string `json:"eden_product_name"`
	PublicProduct1  string `json:"public_product_1"`
	PublicProduct2  string `json:"public_product_2"`
}

// Validate : function to validate uom request data
func (c *matchingRequest) Validate() *validation.Output {
	o := &validation.Output{Valid: true}

	return o
}

// Messages : function to return error validation messages
func (c *matchingRequest) Messages() map[string]string {
	messages := map[string]string{}

	return messages
}
