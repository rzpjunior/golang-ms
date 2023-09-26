// Copyright 2020 PT. Eden Pangan Indonesia. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

package report

import (
	"fmt"
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
	r.GET("/payment-gateway", h.paymentGateway, auth.Authorized("fin_rpt_7_dl"))
	r.GET("/main-outlet", h.mainOutlet, auth.Authorized("main_olt_exp"))
	r.GET("/outlet", h.outlet, auth.Authorized("olt_exp"))
	r.GET("/agent", h.agent, auth.Authorized("agt_exp"))
	r.GET("/voucher-log", h.voucherLog, auth.Authorized("sls_rpt_4_dl"))
	r.GET("/submission", h.submission, auth.Authorized("sls_rpt_8_dl"))
}

// paymentGateway : function to get payment gateway report
func (h *Handler) paymentGateway(c echo.Context) (e error) {
	ctx := c.(*cuxs.Context)
	var date time.Time
	var note string
	var session *auth.SessionData

	if session, e = auth.UserSession(ctx); e != nil {
		return ctx.Serve(e)
	}

	isExport := ctx.QueryParam("export") == "1"
	cond := make(map[string]interface{})

	date = now.NewParse(time.RFC3339, ctx.QueryParam("date")).Time

	transactionDateStr := ctx.QueryParam("transaction_date")
	transactionDateArr := strings.Split(transactionDateStr, "|")

	if transactionDateStr != "" {
		cond["tx.transaction_date between "] = transactionDateArr
		note += fmt.Sprintf("Transaction Date : %s - %s; ", transactionDateArr[0], transactionDateArr[1])
	}

	data, e := getPaymentGateway(cond)
	if e == nil {
		if isExport {
			var file string
			if file, e = getPaymentGatewayXls(date, data, session.Staff, note); e == nil {
				ctx.Files(file)
			}
		} else {
			ctx.Data(data)
		}
	}
	return ctx.Serve(e)
}

// mainOutlet : function to get main outlet report
func (h *Handler) mainOutlet(c echo.Context) (e error) {
	ctx := c.(*cuxs.Context)
	var date time.Time

	var note string
	var session *auth.SessionData

	if session, e = auth.UserSession(ctx); e != nil {
		return ctx.Serve(e)
	}

	isExport := ctx.QueryParam("export") == "1"
	cond := make(map[string]interface{})

	date = now.NewParse(time.RFC3339, ctx.QueryParam("date")).Time

	area, _ := common.Decrypt(ctx.QueryParam("area"))
	termInvoiceSls, _ := common.Decrypt(ctx.QueryParam("term_invoice_sls"))
	termPaymentSls, _ := common.Decrypt(ctx.QueryParam("term_payment_sls"))
	businessType, _ := common.Decrypt(ctx.QueryParam("business_type"))

	status := ctx.QueryParam("status")
	suspended := ctx.QueryParam("suspended")

	if area != 0 {
		cond["m.finance_area_id = "] = area
		var getArea *model.Area
		if getArea, e = repository.ValidArea(area); e != nil {
			return ctx.Serve(e)
		}
		note += fmt.Sprintf("Area : %s - %s; ", getArea.Code, getArea.Name)
	}

	if termInvoiceSls != 0 {
		cond["m.term_invoice_sls_id = "] = termInvoiceSls
		var getTermInvoiceSls *model.InvoiceTerm
		if getTermInvoiceSls, e = repository.ValidInvoiceTerm(termInvoiceSls); e != nil {
			return ctx.Serve(e)
		}
		note += fmt.Sprintf("Default Invoice Term : %s - %s; ", getTermInvoiceSls.Code, getTermInvoiceSls.Name)
	}

	if termPaymentSls != 0 {
		cond["m.term_payment_sls_id = "] = termPaymentSls
		var getTermPaymentSls *model.PurchaseTerm
		if getTermPaymentSls, e = repository.ValidPurchaseTerm(termPaymentSls); e != nil {
			return ctx.Serve(e)
		}
		note += fmt.Sprintf("Payment Term : %s - %s; ", getTermPaymentSls.Code, getTermPaymentSls.Name)
	}

	if businessType != 0 {
		cond["m.business_type_id = "] = businessType
		var getBusinessType *model.BusinessType
		if getBusinessType, e = repository.ValidBusinessType(businessType); e != nil {
			return ctx.Serve(e)
		}
		note += fmt.Sprintf("Business Type : %s - %s; ", getBusinessType.Code, getBusinessType.Name)
	}

	if status != "" {
		cond["m.status = "] = status
		var getStatus *model.Glossary
		if getStatus, e = repository.GetGlossaryMultipleValue("table", "merchant", "attribute", "status", "value_int", status); e != nil {
			return ctx.Serve(e)
		}
		note += fmt.Sprintf("Status : %d - %s; ", getStatus.ValueInt, getStatus.ValueName)
	}

	if suspended != "" {
		cond["m.suspended = "] = suspended
		var getStatus *model.Glossary
		if getStatus, e = repository.GetGlossaryMultipleValue("table", "merchant", "attribute", "suspended", "value_int", suspended); e != nil {
			return ctx.Serve(e)
		}
		note += fmt.Sprintf("Suspended : %d - %s; ", getStatus.ValueInt, getStatus.ValueName)
	}

	data, e := getMainOutlet(cond)
	if e == nil {
		if isExport {
			var file string
			note = strings.TrimSuffix(note, "; ")
			if file, e = getMainOutletXls(date, data, session.Staff, note); e == nil {
				ctx.Files(file)
			}
		} else {
			ctx.Data(data)
		}
	}
	return ctx.Serve(e)
}

// outlet : function to get outlet report
func (h *Handler) outlet(c echo.Context) (e error) {
	ctx := c.(*cuxs.Context)
	var session *auth.SessionData

	if session, e = auth.UserSession(ctx); e != nil {
		return ctx.Serve(e)
	}

	var date time.Time

	isExport := ctx.QueryParam("export") == "1"
	cond := make(map[string]interface{})

	date = now.NewParse(time.RFC3339, ctx.QueryParam("date")).Time

	area, _ := common.Decrypt(ctx.QueryParam("area"))
	merchant, _ := common.Decrypt(ctx.QueryParam("merchant"))
	archetype, _ := common.Decrypt(ctx.QueryParam("archetype"))
	priceSet, _ := common.Decrypt(ctx.QueryParam("price_set"))
	salesperson, _ := common.Decrypt(ctx.QueryParam("salesperson"))

	status := ctx.QueryParam("status")

	if area != 0 {
		cond["b.area_id = "] = area
	}

	if status != "" {
		cond["b.status = "] = status
	}

	if merchant != 0 {
		cond["b.merchant_id = "] = merchant
	}

	if archetype != 0 {
		cond["b.archetype_id = "] = archetype
	}

	if priceSet != 0 {
		cond["b.price_set_id = "] = priceSet
	}

	if salesperson != 0 {
		cond["b.salesperson_id = "] = salesperson
	}

	data, e := getOutlet(cond)
	if e == nil {
		if isExport {
			var file string
			if file, e = getOutletXls(date, data, session.Staff); e == nil {
				ctx.Files(file)
			}
		} else {
			ctx.Data(data)
		}
	}
	return ctx.Serve(e)
}

// agent : function to get agent report
func (h *Handler) agent(c echo.Context) (e error) {
	ctx := c.(*cuxs.Context)

	var date time.Time
	var note string
	var session *auth.SessionData

	if session, e = auth.UserSession(ctx); e != nil {
		return ctx.Serve(e)
	}

	isExport := ctx.QueryParam("export") == "1"
	cond := make(map[string]interface{})

	date = now.NewParse(time.RFC3339, ctx.QueryParam("date")).Time

	status := ctx.QueryParam("status")
	suspended := ctx.QueryParam("suspended")

	area, _ := common.Decrypt(ctx.QueryParam("area"))
	businessType, _ := common.Decrypt(ctx.QueryParam("business_type"))
	archetype, _ := common.Decrypt(ctx.QueryParam("archetype"))
	paymentGroup, _ := common.Decrypt(ctx.QueryParam("payment_group"))
	warehouse, _ := common.Decrypt(ctx.QueryParam("warehouse"))
	priceSet, _ := common.Decrypt(ctx.QueryParam("price_set"))
	salesPerson, _ := common.Decrypt(ctx.QueryParam("salesperson"))

	if status != "" {
		cond["m.status = "] = status
		var getStatus *model.Glossary
		if getStatus, e = repository.GetGlossaryMultipleValue("table", "merchant", "attribute", "status", "value_int", status); e != nil {
			return ctx.Serve(e)
		}
		note += fmt.Sprintf("Status : %d - %s; ", getStatus.ValueInt, getStatus.ValueName)
	}

	if area != 0 {
		cond["m.finance_area_id = "] = area
		var getArea *model.Area
		if getArea, e = repository.ValidArea(area); e != nil {
			return ctx.Serve(e)
		}
		note += fmt.Sprintf("Area : %s - %s; ", getArea.Code, getArea.Name)
	}

	if businessType != 0 {
		cond["m.business_type_id = "] = businessType
		var getBusinessType *model.BusinessType
		if getBusinessType, e = repository.ValidBusinessType(businessType); e != nil {
			return ctx.Serve(e)
		}
		note += fmt.Sprintf("Business Type : %s - %s; ", getBusinessType.Code, getBusinessType.Name)
	}

	if archetype != 0 {
		cond["b.archetype_id = "] = archetype
		var getArchetype *model.Archetype
		if getArchetype, e = repository.ValidArchetype(archetype); e != nil {
			return ctx.Serve(e)
		}
		note += fmt.Sprintf("Archetype : %s - %s; ", getArchetype.Code, getArchetype.Name)
	}

	if paymentGroup != 0 {
		cond["m.payment_group_sls_id = "] = paymentGroup
		var getPaymentGroup *model.PaymentGroup
		if getPaymentGroup, e = repository.ValidPaymentGroup(paymentGroup); e != nil {
			return ctx.Serve(e)
		}
		note += fmt.Sprintf("Payment Group : %s - %s; ", getPaymentGroup.Code, getPaymentGroup.Name)
	}

	if priceSet != 0 {
		cond["b.price_set_id = "] = priceSet
		var getPriceSet *model.PriceSet
		if getPriceSet, e = repository.ValidPriceSet(priceSet); e != nil {
			return ctx.Serve(e)
		}
		note += fmt.Sprintf("Price Set : %s - %s; ", getPriceSet.Code, getPriceSet.Name)
	}

	if warehouse != 0 {
		cond["b.warehouse_id = "] = warehouse
		var getWarehouse *model.Warehouse
		if getWarehouse, e = repository.ValidWarehouse(warehouse); e != nil {
			return ctx.Serve(e)
		}
		note += fmt.Sprintf("Warehouse : %s - %s; ", getWarehouse.Code, getWarehouse.Name)
	}

	if salesPerson != 0 {
		cond["b.salesperson_id = "] = salesPerson
		var getSalesPerson *model.Staff
		if getSalesPerson, e = repository.ValidStaff(salesPerson); e != nil {
			return ctx.Serve(e)
		}
		note += fmt.Sprintf("Sales Person : %s - %s; ", getSalesPerson.Code, getSalesPerson.Name)
	}

	if suspended != "" {
		cond["m.suspended = "] = suspended
		var getStatus *model.Glossary
		if getStatus, e = repository.GetGlossaryMultipleValue("table", "merchant", "attribute", "suspended", "value_int", suspended); e != nil {
			return ctx.Serve(e)
		}
		note += fmt.Sprintf("Suspended : %d - %s; ", getStatus.ValueInt, getStatus.ValueName)
	}

	data, e := getAgent(cond)
	if e == nil {
		if isExport {
			var file string
			note = strings.TrimSuffix(note, "; ")
			if file, e = getAgentXls(date, data, session.Staff, note); e == nil {
				ctx.Files(file)
			}
		} else {
			ctx.Data(data)
		}
	}
	return ctx.Serve(e)
}

// voucherLog : function to get voucher log report
func (h *Handler) voucherLog(c echo.Context) (e error) {
	ctx := c.(*cuxs.Context)
	var date time.Time

	isExport := ctx.QueryParam("export") == "1"
	cond := make(map[string]interface{})

	date = now.NewParse(time.RFC3339, ctx.QueryParam("date")).Time

	voucherID, _ := common.Decrypt(ctx.QueryParam("voucher_id"))
	redeemDateStr := ctx.QueryParam("redeem_date")
	redeemDateArr := strings.Split(redeemDateStr, "|")

	if redeemDateStr != "" {
		cond["DATE(vl.timestamp) between "] = redeemDateArr
	}

	if voucherID != 0 {
		cond["v.id = "] = voucherID
	}

	data, e := getVoucherLog(cond)
	if e == nil {
		if isExport {
			var file string
			if file, e = getVoucherLogXls(date, data); e == nil {
				ctx.Files(file)
			}
		} else {
			ctx.Data(data)
		}
	}
	return ctx.Serve(e)
}

// submission : function to get submission report
func (h *Handler) submission(c echo.Context) (e error) {
	ctx := c.(*cuxs.Context)
	var session *auth.SessionData

	if session, e = auth.UserSession(ctx); e != nil {
		return ctx.Serve(e)
	}

	var date time.Time

	isExport := ctx.QueryParam("export") == "1"
	date = now.NewParse(time.RFC3339, ctx.QueryParam("date")).Time
	condSA := make(map[string]interface{})
	condCA := make(map[string]interface{})
	condOOR := make(map[string]interface{})

	submissionDate := ctx.QueryParam("submission_date")
	submissionDateArr := strings.Split(submissionDate, "|")

	salesGroup, _ := common.Decrypt(ctx.QueryParam("sales_group"))
	salesPerson, _ := common.Decrypt(ctx.QueryParam("sales_person"))

	condSA["DATE(sai.submit_date) BETWEEN "] = submissionDateArr
	condOOR["DATE(sai.submit_date) BETWEEN "] = submissionDateArr
	condCA["DATE(ca.finish_date) BETWEEN "] = submissionDateArr
	if salesGroup != 0 {
		condSA["sa.sales_group_id = "] = salesGroup
		condCA["ca.sales_group_id = "] = salesGroup
		condOOR["s.sales_group_id = "] = salesGroup
	}
	if salesPerson != 0 {
		condSA["sai.salesperson_id = "] = salesPerson
		condCA["ca.salesperson_id = "] = salesPerson
		condOOR["sai.salesperson_id = "] = salesPerson
	}

	data, e := getSubmission(condSA, condCA, condOOR)
	if e == nil {
		if isExport {
			var file string
			if file, e = getSubmissionXls(date, data, session.Staff); e == nil {
				ctx.Files(file)
			}
		} else {
			ctx.Data(data)
		}
	}
	return ctx.Serve(e)
}
