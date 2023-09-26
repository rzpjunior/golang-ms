// Copyright 2020 PT. Eden Pangan Indonesia. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

package report

import (
	"fmt"
	"strconv"
	"strings"
	"time"

	"git.edenfarm.id/cuxs/common"
	"git.edenfarm.id/cuxs/common/now"
	"git.edenfarm.id/cuxs/cuxs"
	"github.com/labstack/echo/v4"

	"git.edenfarm.id/project-version2/api/datastore/repository"
	"git.edenfarm.id/project-version2/api/src/auth"
	"git.edenfarm.id/project-version2/datamodel/model"
)

type Handler struct{}

// URLMapping : function to map url with it's handler and add authorization validation
func (h *Handler) URLMapping(r *echo.Group) {
	r.GET("/sales-order", h.salesOrder, auth.Authorized("sls_rpt_1_dl"))
	r.GET("/sales-order-item", h.salesOrderItem, auth.Authorized("sls_rpt_2_dl"))
	r.GET("/item-recap", h.reportItemRecap, auth.Authorized("src_rpt_1_dl"))
	r.GET("/sales-invoice", h.salesInvoice, auth.Authorized("fin_rpt_1_dl"))
	r.GET("/sales-payment", h.salesPayment, auth.Authorized("fin_rpt_3_dl"))
	r.GET("/prospective-customer", h.prospectiveCustomer, auth.Authorized("pro_cst_exp"))
	r.GET("/sku-discount", h.skuDiscount, auth.Authorized("sls_rpt_9_dl"))
	r.GET("/sku-discount-item", h.skuDiscountItem, auth.Authorized("sls_rpt_10_dl"))
	r.GET("/sales-order-feedback", h.salesOrderFeedback, auth.Authorized("sls_rpt_11_dl"))
	r.GET("/eden-point", h.edenPoint, auth.Authorized("sls_rpt_12_dl"))
}

// salesOrder : function to get sales order report
func (h *Handler) salesOrder(c echo.Context) (e error) {
	ctx := c.(*cuxs.Context)
	isExport := ctx.QueryParam("export") == "1"

	var note string
	var session *auth.SessionData

	if session, e = auth.UserSession(ctx); e != nil {
		return ctx.Serve(e)
	}

	var date time.Time
	date = now.NewParse(time.RFC3339, ctx.QueryParam("date")).Time
	area, _ := common.Decrypt(ctx.QueryParam("area"))
	warehouse, _ := common.Decrypt(ctx.QueryParam("warehouse"))
	deliveryDateStr := ctx.QueryParam("delivery_date")
	deliveryDateArr := strings.Split(deliveryDateStr, "|")
	orderDateStr := ctx.QueryParam("order_date")
	orderDateArr := strings.Split(orderDateStr, "|")
	merchant, _ := common.Decrypt(ctx.QueryParam("merchant"))
	branch, _ := common.Decrypt(ctx.QueryParam("branch"))

	cond := make(map[string]interface{})

	if area != 0 {
		cond["so.area_id = "] = area
		var getArea *model.Area
		if getArea, e = repository.ValidArea(area); e != nil {
			return ctx.Serve(e)
		}
		note += fmt.Sprintf("Area : %s - %s; ", getArea.Code, getArea.Name)
	}

	if warehouse != 0 {
		cond["so.warehouse_id = "] = warehouse
		var getWarehouse *model.Warehouse
		if getWarehouse, e = repository.ValidWarehouse(warehouse); e != nil {
			return ctx.Serve(e)
		}
		note += fmt.Sprintf("Warehouse : %s - %s; ", getWarehouse.Code, getWarehouse.Name)
	}

	if deliveryDateStr != "" {
		cond["so.delivery_date between "] = deliveryDateArr
		note += fmt.Sprintf("Order Delivery Date : %s - %s; ", deliveryDateArr[0], deliveryDateArr[1])
	}

	if orderDateStr != "" {
		cond["so.recognition_date between "] = orderDateArr
		note += fmt.Sprintf("Order Date : %s - %s; ", orderDateArr[0], orderDateArr[1])
	}

	if merchant != 0 {
		cond["b.merchant_id = "] = merchant
		var getMerchant *model.Merchant
		if getMerchant, e = repository.ValidMerchant(merchant); e != nil {
			return ctx.Serve(e)
		}
		note += fmt.Sprintf("Customer : %s - %s; ", getMerchant.Code, getMerchant.Name)
	}

	if branch != 0 {
		cond["so.branch_id = "] = branch
		var getBranch *model.Branch
		if getBranch, e = repository.ValidBranch(branch); e != nil {
			return ctx.Serve(e)
		}
		note += fmt.Sprintf("Outlet : %s - %s; ", getBranch.Code, getBranch.Name)
	}

	data, e := getSalesOrder(cond)
	if e == nil {
		if isExport {
			var file string
			mArea, e := repository.ValidArea(area)
			note = strings.TrimSuffix(note, "; ")
			if file, e = getSalesOrderXls(date, deliveryDateStr, data, mArea, session.Staff, note); e == nil {
				ctx.Files(file)
			}
		} else {
			ctx.Data(data)
		}
	}

	return ctx.Serve(e)
}

// salesOrderItem : function to get sales order report
func (h *Handler) salesOrderItem(c echo.Context) (e error) {
	ctx := c.(*cuxs.Context)
	isExport := ctx.QueryParam("export") == "1"

	var note string
	var session *auth.SessionData

	if session, e = auth.UserSession(ctx); e != nil {
		return ctx.Serve(e)
	}

	var date time.Time
	date = now.NewParse(time.RFC3339, ctx.QueryParam("date")).Time
	area, _ := common.Decrypt(ctx.QueryParam("area"))
	wrt, _ := common.Decrypt(ctx.QueryParam("wrt"))
	warehouse, _ := common.Decrypt(ctx.QueryParam("warehouse"))
	deliveryDateStr := ctx.QueryParam("delivery_date")
	deliveryDateArr := strings.Split(deliveryDateStr, "|")

	cond := make(map[string]interface{})

	if area != 0 {
		cond["so.area_id = "] = area
		var getArea *model.Area
		if getArea, e = repository.ValidArea(area); e != nil {
			return ctx.Serve(e)
		}
		note += fmt.Sprintf("Area : %s - %s; ", getArea.Code, getArea.Name)
	}

	if wrt != 0 {
		cond["so.wrt_id = "] = wrt
		var getWrt *model.Wrt
		if getWrt, e = repository.ValidWrt(wrt); e != nil {
			return ctx.Serve(e)
		}
		note += fmt.Sprintf("WRT : %s - %s; ", getWrt.Code, getWrt.Name)
	}

	if warehouse != 0 {
		cond["so.warehouse_id = "] = warehouse
		var getWarehouse *model.Warehouse
		if getWarehouse, e = repository.ValidWarehouse(warehouse); e != nil {
			return ctx.Serve(e)
		}
		note += fmt.Sprintf("Warehouse : %s - %s; ", getWarehouse.Code, getWarehouse.Name)
	}

	if deliveryDateStr != "" {
		cond["so.delivery_date between "] = deliveryDateArr
		note += fmt.Sprintf("Order Delivery Date : %s - %s; ", deliveryDateArr[0], deliveryDateArr[1])
	}

	data, e := getSalesOrderItem(cond)
	if e == nil {
		if isExport {
			var file string
			mArea, e := repository.ValidArea(area)
			note = strings.TrimSuffix(note, "; ")
			if file, e = getSalesOrderItemXls(date, deliveryDateStr, data, mArea, session.Staff, note); e == nil {
				ctx.Files(file)
			}
		} else {
			ctx.Data(data)
		}
	}

	return ctx.Serve(e)
}

// reportItemRecap : function to get requested data based on parameters
func (h *Handler) reportItemRecap(c echo.Context) (e error) {
	ctx := c.(*cuxs.Context)
	rq := ctx.RequestQuery()
	var backdate time.Time
	var area *model.Area

	isExport := ctx.QueryParam("export") == "1"
	backdate = now.NewParse(time.RFC3339, ctx.QueryParam("date")).Time
	for _, v := range rq.Conditions {
		if val, ok := v["salesorder.area.id"]; ok {
			n, _ := strconv.ParseInt(val, 10, 64)
			area, _ = repository.ValidArea(n)
		}
	}
	data, total, e := GetItemRecap(rq)
	if e == nil {
		if isExport {
			var file string
			if file, e = GetItemRecapXls(backdate, data, area); e == nil {
				ctx.Files(file)
			}
		} else {
			ctx.Data(data, total)
		}
	}

	return ctx.Serve(e)
}

// salesInvoice : function to get sales invoice report
func (h *Handler) salesInvoice(c echo.Context) (e error) {
	ctx := c.(*cuxs.Context)
	var date time.Time
	var note string
	var session *auth.SessionData

	if session, e = auth.UserSession(ctx); e != nil {
		return ctx.Serve(e)
	}

	isExport := ctx.QueryParam("export") == "1"

	date = now.NewParse(time.RFC3339, ctx.QueryParam("date")).Time

	area, _ := common.Decrypt(ctx.QueryParam("area"))
	merchant, _ := common.Decrypt(ctx.QueryParam("merchant"))
	branch, _ := common.Decrypt(ctx.QueryParam("branch"))
	invoiceDateStr := ctx.QueryParam("invoice_date")
	invoiceDateArr := strings.Split(invoiceDateStr, "|")

	invoiceDueDateStr := ctx.QueryParam("invoice_due_date")
	invoiceDueDateArr := strings.Split(invoiceDueDateStr, "|")

	orderDelDateStr := ctx.QueryParam("order_delivery_date")
	orderDeliveryDateArr := strings.Split(orderDelDateStr, "|")

	invoiceStatus := ctx.QueryParam("status")

	cond := make(map[string]interface{})

	if area != 0 {
		cond["so.area_id = "] = area
		var getArea *model.Area
		if getArea, e = repository.ValidArea(area); e != nil {
			return ctx.Serve(e)
		}
		note += fmt.Sprintf("Area : %s - %s; ", getArea.Code, getArea.Name)
	}

	if merchant != 0 {
		cond["b.merchant_id = "] = merchant
		var getMerchant *model.Merchant
		if getMerchant, e = repository.ValidMerchant(merchant); e != nil {
			return ctx.Serve(e)
		}
		note += fmt.Sprintf("Customer : %s - %s; ", getMerchant.Code, getMerchant.Name)
	}

	if branch != 0 {
		cond["so.branch_id = "] = branch
		var getBranch *model.Branch
		if getBranch, e = repository.ValidBranch(branch); e != nil {
			return ctx.Serve(e)
		}
		note += fmt.Sprintf("Branch : %s - %s; ", getBranch.Code, getBranch.Name)
	}

	if invoiceDateStr != "" {
		cond["si.recognition_date between "] = invoiceDateArr
		note += fmt.Sprintf("Invoice Date : %s - %s; ", invoiceDateArr[0], invoiceDateArr[1])
	}

	if invoiceDueDateStr != "" {
		cond["si.due_date between "] = invoiceDueDateArr
		note += fmt.Sprintf("Invoice Due Date : %s - %s; ", invoiceDueDateArr[0], invoiceDueDateArr[1])
	}

	if orderDelDateStr != "" {
		cond["so.delivery_date between "] = orderDeliveryDateArr
		note += fmt.Sprintf("Order Delivery Date : %s - %s; ", orderDeliveryDateArr[0], orderDeliveryDateArr[1])
	}

	if invoiceStatus != "" {
		cond["si.status = "] = invoiceStatus
		var getStatus *model.Glossary
		if getStatus, e = repository.GetGlossaryMultipleValue("table", "sales_invoice", "attribute", "status", "value_int", invoiceStatus); e != nil {
			return ctx.Serve(e)
		}
		note += fmt.Sprintf("Status : %d - %s; ", getStatus.ValueInt, getStatus.ValueName)
	}

	data, e := getSalesInvoice(cond)
	if e == nil {
		if isExport {
			var file string
			mArea, e := repository.ValidArea(area)
			note = strings.TrimSuffix(note, "; ")
			if file, e = getSalesInvoiceXls(date, data, mArea, session.Staff, note); e == nil {
				ctx.Files(file)
			}
		} else {
			ctx.Data(data)
		}
	}

	return ctx.Serve(e)
}

// salesPayment : function to get sales payment report
func (h *Handler) salesPayment(c echo.Context) (e error) {
	ctx := c.(*cuxs.Context)
	var date time.Time
	var note string
	var session *auth.SessionData

	if session, e = auth.UserSession(ctx); e != nil {
		return ctx.Serve(e)
	}

	isExport := ctx.QueryParam("export") == "1"

	date = now.NewParse(time.RFC3339, ctx.QueryParam("date")).Time

	area, _ := common.Decrypt(ctx.QueryParam("area"))
	warehouse, _ := common.Decrypt(ctx.QueryParam("warehouse"))
	recognitionDateStr := ctx.QueryParam("recognition_date")
	recognitionDateArr := strings.Split(recognitionDateStr, "|")
	receivedDateStr := ctx.QueryParam("received_date")
	receivedDateArr := strings.Split(receivedDateStr, "|")
	merchant, _ := common.Decrypt(ctx.QueryParam("merchant"))
	branch, _ := common.Decrypt(ctx.QueryParam("branch"))

	cond := make(map[string]interface{})

	if area != 0 {
		cond["so.area_id = "] = area
		var getArea *model.Area
		if getArea, e = repository.ValidArea(area); e != nil {
			return ctx.Serve(e)
		}
		note += fmt.Sprintf("Area : %s - %s; ", getArea.Code, getArea.Name)
	}

	if warehouse != 0 {
		cond["so.warehouse_id = "] = warehouse
		var getWarehouse *model.Warehouse
		if getWarehouse, e = repository.ValidWarehouse(warehouse); e != nil {
			return ctx.Serve(e)
		}
		note += fmt.Sprintf("Warehouse : %s - %s; ", getWarehouse.Code, getWarehouse.Name)
	}

	if merchant != 0 {
		cond["b.merchant_id = "] = merchant
		var getMerchant *model.Merchant
		if getMerchant, e = repository.ValidMerchant(merchant); e != nil {
			return ctx.Serve(e)
		}
		note += fmt.Sprintf("Customer : %s - %s; ", getMerchant.Code, getMerchant.Name)
	}

	if branch != 0 {
		cond["so.branch_id = "] = branch
		var getBranch *model.Branch
		if getBranch, e = repository.ValidBranch(branch); e != nil {
			return ctx.Serve(e)
		}
		note += fmt.Sprintf("Branch : %s - %s; ", getBranch.Code, getBranch.Name)
	}

	if recognitionDateStr != "" {
		cond["sp.recognition_date between "] = recognitionDateArr
		note += fmt.Sprintf("Payment Date : %s - %s; ", recognitionDateArr[0], recognitionDateArr[1])
	}

	if receivedDateStr != "" {
		cond["sp.received_date between "] = receivedDateArr
		note += fmt.Sprintf("Received Date : %s - %s; ", receivedDateArr[0], receivedDateArr[1])
	}

	data, e := getSalesPayment(cond)
	if e == nil {
		if isExport {
			var file string
			mArea, e := repository.ValidArea(area)
			note = strings.TrimSuffix(note, "; ")
			if file, e = getSalesPaymentXls(date, data, mArea, session.Staff, note); e == nil {
				ctx.Files(file)
			}
		} else {
			ctx.Data(data)
		}
	}

	return ctx.Serve(e)
}

// prospectiveCustomer : function to get prospective customer report
func (h *Handler) prospectiveCustomer(c echo.Context) (e error) {
	ctx := c.(*cuxs.Context)
	var session *auth.SessionData

	if session, e = auth.UserSession(ctx); e != nil {
		return ctx.Serve(e)
	}

	var date time.Time

	isExport := ctx.QueryParam("export") == "1"

	date = now.NewParse(time.RFC3339, ctx.QueryParam("date")).Time

	area, _ := common.Decrypt(ctx.QueryParam("area"))
	archetype, _ := common.Decrypt(ctx.QueryParam("archetype"))
	prosCustStatus := ctx.QueryParam("status")
	salesperson, _ := common.Decrypt(ctx.QueryParam("salesperson"))

	cond := make(map[string]interface{})

	if area != 0 {
		cond["ad.area_id = "] = area
	}

	if archetype != 0 {
		cond["pc.archetype_id = "] = archetype
	}

	if prosCustStatus != "" {
		cond["pc.reg_status = "] = prosCustStatus
	}

	if salesperson != 0 {
		cond["pc.salesperson_id = "] = salesperson
	}

	data, e := getProspectiveCustomer(cond)
	if e == nil {
		if isExport {
			var file string
			mArea, e := repository.ValidArea(area)
			if file, e = getProspectiveCustomerXls(date, data, mArea, session.Staff); e == nil {
				ctx.Files(file)
			}
		} else {
			ctx.Data(data)
		}
	}

	return ctx.Serve(e)
}

// skuDiscount : function to get sku discount report
func (h *Handler) skuDiscount(c echo.Context) (e error) {
	ctx := c.(*cuxs.Context)
	isExport := ctx.QueryParam("export") == "1"
	var note string
	var session *auth.SessionData

	if session, e = auth.UserSession(ctx); e != nil {
		return ctx.Serve(e)
	}

	var date time.Time
	date = now.NewParse(time.RFC3339, ctx.QueryParam("date")).Time
	startPeriodStr := ctx.QueryParam("start_period")
	startPeriodArr := strings.Split(startPeriodStr, "|")
	status := ctx.QueryParam("status")

	cond := make(map[string]interface{})

	if startPeriodStr != "" {
		cond["DATE(sd.start_timestamp) between "] = startPeriodArr
		note += fmt.Sprintf("Start Period : %s - %s; ", startPeriodArr[0], startPeriodArr[1])
	}

	if status != "" {
		cond["sd.status = "] = status
		var getStatus *model.Glossary
		if getStatus, e = repository.GetGlossaryMultipleValue("table", "sku_discount", "attribute", "status", "value_int", status); e != nil {
			return ctx.Serve(e)
		}
		note += fmt.Sprintf("Status : %d - %s; ", getStatus.ValueInt, getStatus.ValueName)
	}

	data, e := getSkuDiscount(cond)
	if e != nil {
		return ctx.Serve(e)
	}

	if isExport {
		var file string
		note = strings.TrimSuffix(note, "; ")
		if file, e = getSkuDiscountXls(date, data, session.Staff, note); e != nil {
			return ctx.Serve(e)
		}

		ctx.Files(file)

		return ctx.Serve(e)
	}

	ctx.Data(data)

	return ctx.Serve(e)
}

// skuDiscountItem : function to get sku discount item report
func (h *Handler) skuDiscountItem(c echo.Context) (e error) {
	ctx := c.(*cuxs.Context)
	isExport := ctx.QueryParam("export") == "1"
	var note string
	var session *auth.SessionData

	if session, e = auth.UserSession(ctx); e != nil {
		return ctx.Serve(e)
	}

	var date time.Time
	date = now.NewParse(time.RFC3339, ctx.QueryParam("date")).Time
	priceSet, _ := common.Decrypt(ctx.QueryParam("price_set"))
	promotionID, _ := common.Decrypt(ctx.QueryParam("promotion_id"))
	status := ctx.QueryParam("status")

	cond := make(map[string]interface{})

	if priceSet != 0 {
		cond["FIND_IN_SET(?, sd.price_set)"] = priceSet
		var getPriceSet *model.PriceSet
		if getPriceSet, e = repository.ValidPriceSet(priceSet); e != nil {
			return ctx.Serve(e)
		}
		note += fmt.Sprintf("Price Set : %s - %s; ", getPriceSet.Code, getPriceSet.Name)
	}

	if promotionID != 0 {
		cond["sd.id = "] = promotionID
		var getPromotionID *model.SkuDiscount
		if getPromotionID, e = repository.ValidSkuDiscount(promotionID); e != nil {
			return ctx.Serve(e)
		}
		note += fmt.Sprintf("Promotion Id: %s - %s; ", getPromotionID.Code, getPromotionID.Name)
	}

	if status != "" {
		cond["sd.status = "] = status
		var getStatus *model.Glossary
		if getStatus, e = repository.GetGlossaryMultipleValue("table", "sku_discount", "attribute", "status", "value_int", status); e != nil {
			return ctx.Serve(e)
		}
		note += fmt.Sprintf("Status : %d - %s; ", getStatus.ValueInt, getStatus.ValueName)
	}

	data, e := getSkuDiscountItem(cond)
	if e != nil {
		return ctx.Serve(e)
	}

	if isExport {
		var file string
		note = strings.TrimSuffix(note, "; ")
		if file, e = getSkuDiscounItemXls(date, data, session.Staff, note); e != nil {
			return ctx.Serve(e)
		}

		ctx.Files(file)

		return ctx.Serve(e)
	}

	ctx.Data(data)

	return ctx.Serve(e)
}

// salesOrder : function to get sales order report
func (h *Handler) salesOrderFeedback(c echo.Context) (e error) {
	ctx := c.(*cuxs.Context)
	isExport := ctx.QueryParam("export") == "1"
	data, e := getSalesOrderFeedback()
	var date time.Time
	date = now.NewParse(time.RFC3339, ctx.QueryParam("date")).Time

	var session *auth.SessionData
	if session, e = auth.UserSession(ctx); e != nil {
		return ctx.Serve(e)
	}

	if e == nil {
		if isExport {
			var file string
			if file, e = getSalesOrderFeedbackXls(date, data, session.Staff); e == nil {
				ctx.Files(file)
			}
		} else {
			ctx.Data(data)
		}
	}

	return ctx.Serve(e)
}

// edenPoint : function to get EdenPoint report
func (h *Handler) edenPoint(c echo.Context) (e error) {
	ctx := c.(*cuxs.Context)
	isExport := ctx.QueryParam("export") == "1"
	var date time.Time
	date = now.NewParse(time.RFC3339, ctx.QueryParam("date")).Time

	var note string
	var session *auth.SessionData

	if session, e = auth.UserSession(ctx); e != nil {
		return ctx.Serve(e)
	}

	startPeriodStr := ctx.QueryParam("period")
	startPeriodArr := strings.Split(startPeriodStr, "|")
	merchant, _ := common.Decrypt(ctx.QueryParam("merchant"))

	cond := make(map[string]interface{})

	if startPeriodStr != "" {
		cond["DATE(mpl.created_date) between "] = startPeriodArr
		note += fmt.Sprintf("Period : %s - %s; ", startPeriodArr[0], startPeriodArr[1])
	}

	if merchant != 0 {
		cond["m.id = "] = merchant
		var getMerchant *model.Merchant
		if getMerchant, e = repository.ValidMerchant(merchant); e != nil {
			return ctx.Serve(e)
		}
		note += fmt.Sprintf("Customer : %s - %s; ", getMerchant.Code, getMerchant.Name)
	}

	data, e := getEdenPoint(cond)
	if e != nil {
		return ctx.Serve(e)
	}

	if isExport {
		var file string
		note = strings.TrimSuffix(note, "; ")
		if file, e = getEdenPointXls(date, data, session.Staff, note); e != nil {
			return ctx.Serve(e)
		}

		ctx.Files(file)

		return ctx.Serve(e)
	}

	ctx.Data(data)

	return ctx.Serve(e)
}
