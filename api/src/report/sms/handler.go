// Copyright 2020 PT. Eden Pangan Indonesia. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

package sms

import (
	"fmt"
	"strings"

	"git.edenfarm.id/cuxs/common"
	"git.edenfarm.id/cuxs/cuxs"
	"github.com/labstack/echo/v4"

	"git.edenfarm.id/project-version2/api/datastore/repository"
	"git.edenfarm.id/project-version2/api/src/auth"
	"git.edenfarm.id/project-version2/datamodel/model"
)

type Handler struct{}

// URLMapping : function to map url with it's handler and add authorization validation
func (h *Handler) URLMapping(r *echo.Group) {
	r.GET("/purchase-order", h.purchaseOrder, auth.Authorized("src_rpt_2_dl"))
	r.GET("/purchase-order-item", h.purchaseOrderItem, auth.Authorized("src_rpt_3_dl"))
	r.GET("/purchase-invoice", h.purchaseInvoice, auth.Authorized("fin_rpt_4_dl"))
	r.GET("/purchase-payment", h.purchasePayment, auth.Authorized("fin_rpt_6_dl"))
	r.GET("/purchase-invoice-item", h.purchaseInvoiceItem, auth.Authorized("fin_rpt_5_dl"))
	r.GET("/cogs", h.cogs, auth.Authorized("fin_rpt_8_dl"))
	r.GET("/price-comparison", h.priceComparison, auth.Authorized("src_rpt_5_dl"))
	r.GET("/inbound-time", h.inbound, auth.Authorized("src_rpt_6_dl"))
	r.GET("/field-purchase-order-item", h.fieldPurchaser, auth.Authorized("src_rpt_7_dl"))
}

// purchaseOrder : function to get purchase order report
func (h *Handler) purchaseOrder(c echo.Context) (e error) {
	ctx := c.(*cuxs.Context)
	isExport := ctx.QueryParam("export") == "1"

	area, _ := common.Decrypt(ctx.QueryParam("area"))
	warehouse, _ := common.Decrypt(ctx.QueryParam("warehouse"))
	etaDateStr := ctx.QueryParam("eta_date")
	etaDateArr := strings.Split(etaDateStr, "|")
	orderDateStr := ctx.QueryParam("order_date")
	orderDateArr := strings.Split(orderDateStr, "|")
	supplier, _ := common.Decrypt(ctx.QueryParam("supplier"))
	grStatus := ctx.QueryParam("goods_receipt")

	cond := make(map[string]interface{})

	if area != 0 {
		cond["w.area_id = "] = area
	}

	if warehouse != 0 {
		cond["w.id = "] = warehouse
	}

	if etaDateStr != "" {
		cond["po.eta_date between "] = etaDateArr
	}

	if orderDateStr != "" {
		cond["po.recognition_date between "] = orderDateArr
	}

	if supplier != 0 {
		cond["s.id = "] = supplier
	}

	if !(grStatus == "" || grStatus == "0") {
		// if parameter is 1 then return gr that is finished, else return other
		cond["po.has_finished_gr = "] = grStatus
	}

	data, e := getPurchaseOrder(cond)
	if e == nil {
		if isExport {
			var file string
			mArea, e := repository.ValidArea(area)
			if file, e = GetPurchaseOrderXls(etaDateStr, data, mArea); e == nil {
				ctx.Files(file)
			}
		} else {
			ctx.Data(data)
		}
	}

	return ctx.Serve(e)
}

// purchaseOrderItem : function to get purchase order item report
func (h *Handler) purchaseOrderItem(c echo.Context) (e error) {
	ctx := c.(*cuxs.Context)
	isExport := ctx.QueryParam("export") == "1"

	area, _ := common.Decrypt(ctx.QueryParam("area"))
	warehouse, _ := common.Decrypt(ctx.QueryParam("warehouse"))
	etaDateStr := ctx.QueryParam("eta_date")
	etaDateArr := strings.Split(etaDateStr, "|")
	orderDateStr := ctx.QueryParam("order_date")
	orderDateArr := strings.Split(orderDateStr, "|")
	supplier, _ := common.Decrypt(ctx.QueryParam("supplier"))

	cond := make(map[string]interface{})

	if area != 0 {
		cond["a.id = "] = area
	}

	if warehouse != 0 {
		cond["w.id = "] = warehouse
	}

	if etaDateStr != "" {
		cond["po.eta_date between "] = etaDateArr
	}

	if orderDateStr != "" {
		cond["po.recognition_date between "] = orderDateArr
	}

	if supplier != 0 {
		cond["s.id = "] = supplier
	}

	data, e := getPurchaseOrderItem(cond)
	if e == nil {
		if isExport {
			var file string
			mArea, e := repository.ValidArea(area)
			if file, e = GetPurchaseOrderItemXls(etaDateStr, data, mArea); e == nil {
				ctx.Files(file)
			}
		} else {
			ctx.Data(data)
		}
	}

	return ctx.Serve(e)
}

// purchaseInvoice : function to get purchase invoice report
func (h *Handler) purchaseInvoice(c echo.Context) (e error) {
	ctx := c.(*cuxs.Context)

	var note string
	var session *auth.SessionData

	if session, e = auth.UserSession(ctx); e != nil {
		return ctx.Serve(e)
	}

	isExport := ctx.QueryParam("export") == "1"

	area, _ := common.Decrypt(ctx.QueryParam("area"))
	warehouse, _ := common.Decrypt(ctx.QueryParam("warehouse"))
	invoiceDateStr := ctx.QueryParam("invoice_date")
	invoiceDateArr := strings.Split(invoiceDateStr, "|")
	invoiceDueDateStr := ctx.QueryParam("invoice_due_date")
	invoiceDueDateArr := strings.Split(invoiceDueDateStr, "|")
	supplier, _ := common.Decrypt(ctx.QueryParam("supplier"))
	ataDateStr := ctx.QueryParam("ata_date")
	ataDateArr := strings.Split(ataDateStr, "|")

	cond := make(map[string]interface{})

	if area != 0 {
		cond["a.id = "] = area
		var getArea *model.Area
		if getArea, e = repository.ValidArea(area); e != nil {
			return ctx.Serve(e)
		}
		note += fmt.Sprintf("Area : %s - %s; ", getArea.Code, getArea.Name)
	}

	if warehouse != 0 {
		cond["w.id = "] = warehouse
		var getWarehouse *model.Warehouse
		if getWarehouse, e = repository.ValidWarehouse(warehouse); e != nil {
			return ctx.Serve(e)
		}
		note += fmt.Sprintf("Warehouse : %s - %s; ", getWarehouse.Code, getWarehouse.Name)
	}

	if invoiceDateStr != "" {
		cond["pi.recognition_date between "] = invoiceDateArr
		note += fmt.Sprintf("Invoice Date : %s - %s; ", invoiceDateArr[0], invoiceDateArr[1])
	}

	if invoiceDueDateStr != "" {
		cond["pi.due_date between "] = invoiceDueDateArr
		note += fmt.Sprintf("Invoice Due Date : %s - %s; ", invoiceDueDateArr[0], invoiceDueDateArr[1])
	}

	if supplier != 0 {
		cond["s.id = "] = supplier
		var getSupplier *model.Supplier
		if getSupplier, e = repository.ValidSupplier(supplier); e != nil {
			return ctx.Serve(e)
		}
		note += fmt.Sprintf("Supplier : %s - %s; ", getSupplier.Code, getSupplier.Name)
	}

	if ataDateStr != "" {
		cond["gr.ata_date between "] = ataDateArr
		note += fmt.Sprintf("Ata Date : %s - %s; ", ataDateArr[0], ataDateArr[1])
	}

	data, e := getPurchaseInvoice(cond)
	if e == nil {
		if isExport {
			var file string
			mArea, e := repository.ValidArea(area)
			if file, e = GetPurchaseInvoiceXls(invoiceDueDateStr, data, mArea, session.Staff, note); e == nil {
				ctx.Files(file)
			}
		} else {
			ctx.Data(data)
		}
	}

	return ctx.Serve(e)
}

// purchasePayment : function to get purchase payment report
func (h *Handler) purchasePayment(c echo.Context) (e error) {
	ctx := c.(*cuxs.Context)

	var note string
	var session *auth.SessionData

	if session, e = auth.UserSession(ctx); e != nil {
		return ctx.Serve(e)
	}

	isExport := ctx.QueryParam("export") == "1"

	area, _ := common.Decrypt(ctx.QueryParam("area"))
	paymentDateStr := ctx.QueryParam("payment_date")
	paymentDateArr := strings.Split(paymentDateStr, "|")
	supplier, _ := common.Decrypt(ctx.QueryParam("supplier"))

	cond := make(map[string]interface{})

	if area != 0 {
		cond["ar.id = "] = area
		var getArea *model.Area
		if getArea, e = repository.ValidArea(area); e != nil {
			return ctx.Serve(e)
		}
		note += fmt.Sprintf("Area : %s - %s; ", getArea.Code, getArea.Name)
	}

	if paymentDateStr != "" {
		cond["pp.recognition_date between "] = paymentDateArr
		note += fmt.Sprintf("Payment Date : %s - %s; ", paymentDateArr[0], paymentDateArr[1])

	}

	if supplier != 0 {
		cond["s.id = "] = supplier
		var getSupplier *model.Supplier
		if getSupplier, e = repository.ValidSupplier(supplier); e != nil {
			return ctx.Serve(e)
		}
		note += fmt.Sprintf("Supplier : %s - %s; ", getSupplier.Code, getSupplier.Name)
	}

	data, e := getPurchasePayment(cond)
	if e == nil {
		if isExport {
			var file string
			mArea, e := repository.ValidArea(area)
			if file, e = GetPurchasePaymentXls(paymentDateStr, data, mArea, session.Staff, note); e == nil {
				ctx.Files(file)
			}
		} else {
			ctx.Data(data)
		}
	}

	return ctx.Serve(e)
}

// purchaseInvoiceItem : function to get purchase invoice item report
func (h *Handler) purchaseInvoiceItem(c echo.Context) (e error) {
	ctx := c.(*cuxs.Context)

	var note string
	var session *auth.SessionData

	if session, e = auth.UserSession(ctx); e != nil {
		return ctx.Serve(e)
	}

	isExport := ctx.QueryParam("export") == "1"

	area, _ := common.Decrypt(ctx.QueryParam("area"))
	warehouse, _ := common.Decrypt(ctx.QueryParam("warehouse"))
	etaDateStr := ctx.QueryParam("eta_date")
	etaDateArr := strings.Split(etaDateStr, "|")
	supplier, _ := common.Decrypt(ctx.QueryParam("supplier"))

	cond := make(map[string]interface{})

	if area != 0 {
		cond["a.id = "] = area
		var getArea *model.Area
		if getArea, e = repository.ValidArea(area); e != nil {
			return ctx.Serve(e)
		}
		note += fmt.Sprintf("Area : %s - %s; ", getArea.Code, getArea.Name)
	}

	if warehouse != 0 {
		cond["w.id = "] = warehouse
		var getWarehouse *model.Warehouse
		if getWarehouse, e = repository.ValidWarehouse(warehouse); e != nil {
			return ctx.Serve(e)
		}
		note += fmt.Sprintf("Warehouse : %s - %s; ", getWarehouse.Code, getWarehouse.Name)
	}

	if etaDateStr != "" {
		cond["po.eta_date between "] = etaDateArr
		note += fmt.Sprintf("Eta Date : %s - %s; ", etaDateArr[0], etaDateArr[1])
	}

	if supplier != 0 {
		cond["s.id = "] = supplier
		var getSupplier *model.Supplier
		if getSupplier, e = repository.ValidSupplier(supplier); e != nil {
			return ctx.Serve(e)
		}
		note += fmt.Sprintf("Supplier : %s - %s; ", getSupplier.Code, getSupplier.Name)
	}

	data, e := getPurchaseInvoiceItem(cond)
	if e == nil {
		if isExport {
			var file string
			mArea, e := repository.ValidArea(area)
			if file, e = GetPurchaseInvoiceItemXls(etaDateStr, data, mArea, session.Staff, note); e == nil {
				ctx.Files(file)
			}
		} else {
			ctx.Data(data)
		}
	}

	return ctx.Serve(e)
}

// cogs : function to get cogs report
func (h *Handler) cogs(c echo.Context) (e error) {
	ctx := c.(*cuxs.Context)

	var note string
	var session *auth.SessionData

	if session, e = auth.UserSession(ctx); e != nil {
		return ctx.Serve(e)
	}

	isExport := ctx.QueryParam("export") == "1"

	area, _ := common.Decrypt(ctx.QueryParam("area"))
	etaDateStr := ctx.QueryParam("eta_date")
	etaDateArr := strings.Split(etaDateStr, "|")
	warehouse, _ := common.Decrypt(ctx.QueryParam("warehouse"))

	cond := make(map[string]interface{})

	if area != 0 {
		cond["a.id = "] = area
		var getArea *model.Area
		if getArea, e = repository.ValidArea(area); e != nil {
			return ctx.Serve(e)
		}
		note += fmt.Sprintf("Area : %s - %s; ", getArea.Code, getArea.Name)
	}

	if etaDateStr != "" {
		cond["c.eta_date between "] = etaDateArr
		note += fmt.Sprintf("Eta Date : %s - %s; ", etaDateArr[0], etaDateArr[1])
	}

	if warehouse != 0 {
		cond["w.id = "] = warehouse
		var getWarehouse *model.Warehouse
		if getWarehouse, e = repository.ValidWarehouse(warehouse); e != nil {
			return ctx.Serve(e)
		}
		note += fmt.Sprintf("Warehouse : %s - %s; ", getWarehouse.Code, getWarehouse.Name)
	}

	data, e := getCogs(cond)
	if e == nil {
		if isExport {
			var file string
			mArea, e := repository.ValidArea(area)
			if file, e = GetCogsXls(etaDateStr, data, mArea, session.Staff, note); e == nil {
				ctx.Files(file)
			}
		} else {
			ctx.Data(data)
		}
	}

	return ctx.Serve(e)
}

// priceComparison : function to get price comparison report
func (h *Handler) priceComparison(c echo.Context) (e error) {
	ctx := c.(*cuxs.Context)
	isExport := ctx.QueryParam("export") == "1"

	area, _ := common.Decrypt(ctx.QueryParam("area"))
	surveyDateStr := ctx.QueryParam("survey_date")
	surveyDateArr := strings.Split(surveyDateStr, "|")

	cond := make(map[string]interface{})

	if area != 0 {
		cond["da.id = "] = area
	}

	if surveyDateStr != "" {
		cond["dpri.scraped_date between "] = surveyDateArr
	}

	data, e := getPriceComparison(cond)
	if e == nil {
		if isExport {
			var file string
			mArea, e := repository.ValidArea(area)
			if file, e = GetPriceComparisonXls(surveyDateStr, data, mArea); e == nil {
				ctx.Files(file)
			}
		} else {
			ctx.Data(data)
		}
	}

	return ctx.Serve(e)
}

// inbound : function to get inbound time report
func (h *Handler) inbound(c echo.Context) (e error) {
	ctx := c.(*cuxs.Context)
	isExport := ctx.QueryParam("export") == "1"

	var warehouse int64

	if ctx.QueryParam("warehouse") != "" {
		warehouse, _ = common.Decrypt(ctx.QueryParam("warehouse"))
	}
	etaDateStr := ctx.QueryParam("eta_date")
	etaDateArr := strings.Split(etaDateStr, "|")

	cond := make(map[string]interface{})

	if warehouse != 0 {
		cond["w.id = "] = warehouse
	}

	if etaDateStr != "" {
		cond["po.eta_date between "] = etaDateArr
	}

	data, e := getInbound(cond)
	if e == nil {
		if isExport {
			var file string
			mWarehouse, e := repository.ValidWarehouse(warehouse)
			if file, e = GetInboundXls(etaDateStr, data, mWarehouse); e == nil {
				ctx.Files(file)
			}
		} else {
			ctx.Data(data)
		}
	}

	return ctx.Serve(e)
}

// fieldPurchaser : function to get field purchaser report
func (h *Handler) fieldPurchaser(c echo.Context) (e error) {
	ctx := c.(*cuxs.Context)
	isExport := ctx.QueryParam("export") == "1"

	area, _ := common.Decrypt(ctx.QueryParam("area"))
	warehouse, _ := common.Decrypt(ctx.QueryParam("warehouse"))
	recognitionDateStr := ctx.QueryParam("recognition_date")
	recognitionDateArr := strings.Split(recognitionDateStr, "|")
	supplier, _ := common.Decrypt(ctx.QueryParam("supplier"))
	ppCode := ctx.QueryParam("pp_code")

	cond := make(map[string]interface{})

	if warehouse != 0 {
		cond["w.id = "] = warehouse
	}

	if recognitionDateStr != "" {
		cond["pp.recognition_date between "] = recognitionDateArr
	}

	if supplier != 0 {
		cond["sup.id = "] = supplier
	}

	if ppCode != "" {
		cond["pp.code = "] = ppCode
	}

	data, e := getFieldPurchaser(cond)
	if e != nil {
		return ctx.Serve(&echo.BindingError{})
	}

	if isExport {
		var file string
		modelArea, e := repository.ValidArea(area)
		if e != nil {
			return ctx.Serve(e)
		}
		file, e = GetFieldPurchaserXls(recognitionDateStr, data, modelArea)
		if e != nil {
			return ctx.Serve(e)
		}
		ctx.Files(file)
	} else {
		ctx.Data(data)
	}

	return ctx.Serve(e)
}
