// Copyright 2020 PT. Eden Pangan Indonesia Indonesia. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

package forecast_demand

import (
	"git.edenfarm.id/cuxs/validation"
	"git.edenfarm.id/project-version2/api/src/auth"
	"git.edenfarm.id/project-version2/api/util"
)

type updateRequest struct {
	Data []*forecastRequest `json:"data" valid:"required"`

	Session *auth.SessionData `json:"-"`
}

type forecastRequest struct {
	WarehouseCode string  `json:"warehouse_code" valid:"required"`
	ProductCode   string  `json:"product_code" valid:"required"`
	ForecastDate  string  `json:"forecast_date" valid:"required"`
	ForecastQty   float64 `json:"forecast_qty" valid:"required"`
}

func (r *updateRequest) Validate() *validation.Output {
	o := &validation.Output{Valid: true}

	return o
}

func (r *updateRequest) Messages() map[string]string {
	messages := map[string]string{
		"data.required": util.ErrorInputRequired("data"),
	}

	return messages
}
