// Copyright 2020 PT. Eden Pangan Indonesia. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

package webhook

import (
	"encoding/json"
	"time"

	"git.edenfarm.id/cuxs/cuxs"
	"git.edenfarm.id/cuxs/orm"
	"github.com/labstack/echo/v4"
)

type Handler struct{}

// URLMapping : function to map url with it's handler and add authorization validation
func (h *Handler) URLMapping(r *echo.Group) {
	r.POST("", h.webhook)
}

func (h *Handler) webhook(c echo.Context) (e error) {
	ctx := c.(*cuxs.Context)
	o := orm.NewOrm()
	var r createRequest
	var merchantID int64

	if e = ctx.Bind(&r); e == nil {
		soCode := r.Data["sales_order_code"].(string)
		o.Using("read_only")
		o.Raw("SELECT merchant_id FROM sales_order RIGHT JOIN branch ON sales_order.branch_id=branch.id WHERE sales_order.code = ?;", soCode).QueryRow(&merchantID)
		body, _ := json.Marshal(r.Data)
		o.Using("default")
		_, e = o.Raw("insert into webhook_temp_request (merchant_id,request_body, receive_at) values (?, ?, ?)", merchantID, string(body), time.Now()).Exec()
	}

	return
}
