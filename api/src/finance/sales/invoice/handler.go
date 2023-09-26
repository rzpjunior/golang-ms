// Copyright 2021 PT. Eden Pangan Indonesia. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

package invoice

import (
	"strconv"
	"strings"

	"git.edenfarm.id/cuxs/common"
	"git.edenfarm.id/cuxs/cuxs"
	"git.edenfarm.id/project-version2/api/datastore/repository"
	"git.edenfarm.id/project-version2/api/src/auth"
	"git.edenfarm.id/project-version2/api/util"
	"git.edenfarm.id/project-version2/datamodel/model"
	"github.com/labstack/echo/v4"
)

type Handler struct{}

// URLMapping : function to map url with it's handler and add authorization validation
func (h *Handler) URLMapping(r *echo.Group) {
	r.GET("", h.read, auth.Authorized("si_rdl"))
	r.GET("/:id", h.detail, auth.Authorized("si_rdd"))
	r.GET("/filter", h.readFilter, auth.Authorized("filter_rdl"))
	r.POST("", h.create, auth.Authorized("si_crt"))
	r.PUT("/:id", h.update, auth.Authorized("si_upd"))
	r.PUT("/cancel/:id", h.cancel, auth.Authorized("si_can"))
	r.GET("/print/:id", h.receivePrint, auth.Authorized("si_prt"))
	r.GET("/remaining/:id", h.readInvoiceRemainingAmount, auth.Authorized("si_rdl"))
	r.GET("/print-edn/:id", h.receivePrintNota, auth.Authorized("si_edn_prt"))
}

// receivePrint : function to print data based on invoice sales id
func (h *Handler) receivePrint(c echo.Context) (e error) {
	//var r requestGet

	ctx := c.(*cuxs.Context)
	var si *model.SalesInvoice
	var id int64
	var session *auth.SessionData
	configs := make(map[string]string)
	req := make(map[string]interface{})
	if session, e = auth.UserSession(ctx); e == nil {
		if id, e = ctx.Decrypt("id"); e == nil {
			if si, e = repository.GetSalesInvoice("id", id); e != nil {
				e = echo.ErrNotFound
			} else {
				req["si"] = si
				req["session"] = session.Staff.ID + 56
				if config, _, e := repository.GetConfigAppsByAttribute("attribute__icontains", "company"); e == nil {
					for _, v := range config {
						configs[strings.TrimPrefix(v.Attribute, "company_")] = v.Value
					}
					configs["address"] = strings.ReplaceAll(configs["address"], "<br>", "\n")
					req["config"] = configs
				} else {
					e = echo.ErrNotFound
				}

				file := util.SendPrint(req, "read/si")
				ctx.Files(file)

				// delta print
				si.DeltaPrint = si.DeltaPrint + 1
				si.Save("DeltaPrint")
			}
		}
	}

	return ctx.Serve(e)
}

// read : function to get requested data based on parameters
func (h *Handler) read(c echo.Context) (e error) {
	ctx := c.(*cuxs.Context)
	rq := ctx.RequestQuery()

	var data []*model.SalesInvoice
	var total int64
	var merchantArrList []int64
	var merchantJoin []int64
	var deliveryDate []string

	deliveryDateStr := ctx.QueryParam("delivery_date")
	if deliveryDateStr != "" {
		deliveryDateArr := strings.Split(deliveryDateStr, "|")
		deliveryDate = deliveryDateArr
	}

	merchantStr := ctx.QueryParam("merchant")
	if merchantStr != "" {
		merchantArr := strings.Split(merchantStr, "|")

		for _, v := range merchantArr {
			merchant, _ := common.Decrypt(v)

			merchantArrList = append(merchantArrList, merchant)
		}

		merchantJoin = merchantArrList
	}

	if data, total, e = repository.GetSalesInvoices(rq, deliveryDate, merchantJoin...); e == nil {
		ctx.Data(data, total)
	}

	return ctx.Serve(e)
}

// readFilter : function to get requested data based on parameters with filtered permission
func (h *Handler) readFilter(c echo.Context) (e error) {
	ctx := c.(*cuxs.Context)
	rq := ctx.RequestQuery()

	var data []*model.SalesInvoice
	var total int64
	var empty []string
	var empty2 []int64

	if data, total, e = repository.GetSalesInvoices(rq, empty, empty2...); e == nil {
		ctx.Data(data, total)
	}

	return ctx.Serve(e)
}

// detail : function to get detailed data by id
func (h *Handler) detail(c echo.Context) (e error) {
	ctx := c.(*cuxs.Context)

	var id int64
	if id, e = ctx.Decrypt("id"); e == nil {
		if ctx.ResponseData, e = repository.GetSalesInvoice("id", id); e != nil {
			e = echo.ErrNotFound
		}
	}

	return ctx.Serve(e)
}

//create : function to create new data based on input
func (h *Handler) create(c echo.Context) (e error) {
	ctx := c.(*cuxs.Context)
	var r createRequest

	if r.Session, e = auth.UserSession(ctx); e == nil {
		if e = ctx.Bind(&r); e == nil {
			ctx.ResponseData, e = Save(r)
		} else {
			//post error
			errLog := util.ErrorLog{
				ErrorCode:    422,
				Name:         r.Session.Staff.Name,
				Email:        r.Session.Staff.User.Email,
				ErrorMessage: e.Error(),
				Function:     "create_sales_invoice",
				Platform:     "dashboard",
			}
			util.PostToServiceErrorLog(errLog)
		}
	}

	return ctx.Serve(e)
}

// update : function to update requested data based on parameters
func (h *Handler) update(c echo.Context) (e error) {
	ctx := c.(*cuxs.Context)
	var r updateRequest

	if r.ID, e = ctx.Decrypt("id"); e == nil {
		if r.SalesInvoice, e = repository.GetSalesInvoice("id", r.ID); e == nil {
			if r.Session, e = auth.UserSession(ctx); e == nil {
				if e = ctx.Bind(&r); e == nil {
					ctx.ResponseData, e = Update(r)
				}
			}
		}
	}
	return ctx.Serve(e)
}

// unArchive : function to set status of archived data into active
func (h *Handler) cancel(c echo.Context) (e error) {
	ctx := c.(*cuxs.Context)
	var r cancelRequest

	if r.Session, e = auth.UserSession(ctx); e == nil {
		if r.ID, e = ctx.Decrypt("id"); e == nil {
			if r.SalesInvoice, e = repository.GetSalesInvoice("id", r.ID); e == nil {
				if e = ctx.Bind(&r); e == nil {
					ctx.ResponseData, e = Cancel(r)
				}
			}
		}
	}

	return ctx.Serve(e)
}

// readInvoiceRemainingAmount : function to get remaining amount of sales invoice
func (h *Handler) readInvoiceRemainingAmount(c echo.Context) (e error) {
	ctx := c.(*cuxs.Context)

	var id int64

	if id, e = ctx.Decrypt("id"); e != nil {
		return ctx.Serve(e)
	}
	if ctx.ResponseData, e = repository.CheckRemainingSalesInvoiceAmount(id); e != nil {
		return ctx.Serve(e)
	}

	return ctx.Serve(e)
}

// receivePrintNota : function to print data based on invoice sales id with order type EDN Sales
func (h *Handler) receivePrintNota(c echo.Context) (e error) {

	ctx := c.(*cuxs.Context)
	var (
		si     *model.SalesInvoice
		id     int64
		format int
		file   string
	)
	req := make(map[string]interface{})

	if id, e = ctx.Decrypt("id"); e != nil {
		return ctx.Serve(e)
	}

	format, e = strconv.Atoi(ctx.QueryParam("format"))

	if si, e = repository.GetSalesInvoice("id", id); e != nil {
		return ctx.Serve(e)
	}

	if e = si.SalesOrder.Read("ID"); e != nil {
		return ctx.Serve(e)
	}

	if e = si.SalesOrder.Branch.Read("ID"); e != nil {
		return ctx.Serve(e)
	}

	if e = si.SalesOrder.Warehouse.Read("ID"); e != nil {
		return ctx.Serve(e)
	}

	req["company"] = "UD OKI"

	req["si"] = si

	if format == 1 {
		file = util.SendPrint(req, "read/nota_edn")
	} else {
		file = util.SendPrint(req, "read/nota_edn_termal")
	}

	ctx.Files(file)

	// delta print
	si.DeltaPrint = si.DeltaPrint + 1
	si.Save("DeltaPrint")

	return ctx.Serve(e)
}
