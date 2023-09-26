// Copyright 2020 PT. Qasico Teknologi Indonesia. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

package order

import (
	"net/http"

	"git.edenfarm.id/cuxs/env"
	"git.edenfarm.id/cuxs/validation"
	"git.edenfarm.id/project-version2/datamodel/model"
	"git.edenfarm.id/project-version2/api/src/auth"
)

type bulkCommitRequest struct {
	PurchaseOrderIDs []string `json:"data" valid:"required"`
	SuccessCount     int64    `json:"-"`

	PurchaseOrder *model.PurchaseOrder `json:"-"`
	Session       *auth.SessionData    `json:"-"`
}

func (r *bulkCommitRequest) Validate() *validation.Output {
	o := &validation.Output{Valid: true}
	var client = &http.Client{}
	var request *http.Request
	var err error

	baseURL := env.GetString("SERVER_HOST", "")

	for _, v := range r.PurchaseOrderIDs {
		if request, err = http.NewRequest("PUT", "http://"+baseURL+"/v1/purchase/order/commit/"+v, nil); err == nil {
			request.Header.Set("Authorization", "Bearer "+r.Session.Token)
			request.Header.Set("Content-Type", "application/json")
			response, err := client.Do(request)
			if err == nil {
				if response.StatusCode == 200 {
					r.SuccessCount++
				}
			}

			defer response.Body.Close()
		}
	}

	if r.SuccessCount == 0 {
		o.Failure("id.invalid", "No data has been saved successfully")
	}

	return o
}

func (r *bulkCommitRequest) Messages() map[string]string {
	return map[string]string{}
}
