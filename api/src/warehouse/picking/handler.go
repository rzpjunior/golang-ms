// Copyright 2020 PT. Eden Pangan Indonesia. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

package picking

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"strings"
	"sync"
	"time"

	"git.edenfarm.id/cuxs/common"
	"git.edenfarm.id/cuxs/common/now"
	"git.edenfarm.id/cuxs/cuxs"
	"git.edenfarm.id/cuxs/dbredis"
	"git.edenfarm.id/cuxs/env"
	"git.edenfarm.id/cuxs/mongodb"
	"git.edenfarm.id/cuxs/orm"
	"git.edenfarm.id/project-version2/api/datastore/repository"
	"git.edenfarm.id/project-version2/api/src/auth"
	"git.edenfarm.id/project-version2/api/util"
	"git.edenfarm.id/project-version2/api/util/kafka"
	"git.edenfarm.id/project-version2/datamodel/model"
	"github.com/labstack/echo/v4"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Handler struct{}

// URLMapping : function to map url with it's handler and add authorization validation
func (h *Handler) URLMapping(r *echo.Group) {
	r.GET("", h.read, auth.Authorized("pco_rdl"))
	r.GET("/:id", h.detail, auth.Authorized("pco_rdd"))
	r.GET("/list-assign", h.readAssign, auth.Authorized("pco_rdd"))
	r.GET("/items/:id", h.itemByPackingAssign, auth.Authorized("pco_rdd"))
	r.GET("/items", h.readItems, auth.Authorized("pco_rdl"))
	r.POST("", h.create, auth.Authorized("pco_crt"))

	r.GET("/template", h.pickingOrderTemplate, auth.Authorized("pco_dwl"))
	r.POST("/upload-assign", h.uploadAssign, auth.Authorized("pco_upl"))

	// picking assign
	r.GET("/assign/print/:id", h.receivePrint, auth.Authorized("pco_prt"))
	r.GET("/assign/print", h.printLabels, auth.Authorized("pco_prt"))

	r.GET("/assign/list-group", h.getPickingLists, auth.Authorized("pco_gen_pl"))
	r.POST("/assign/list-generate", h.listPickingGenerate, auth.Authorized("pco_gen_pl"))
	r.PUT("/assign/leadpicker/pl/:id", h.assignLeadPicker, auth.Authorized("pco_asg_picker"))
	r.POST("/assign/summary", h.summaryPicking, auth.Authorized("pco_gen_pl"))

	r.GET("/new-assign", h.readNewAssign, auth.AuthorizedMobile())
	r.GET("/assign", h.readAssign, auth.AuthorizedMobile())
	r.GET("/assign/:id", h.assignDetail, auth.AuthorizedMobile())
	r.PUT("/assign/:id", h.updateAssign, auth.AuthorizedMobile())
	r.PUT("/assign/checkout/:id", h.checkoutAssign, auth.AuthorizedMobile())
	r.PUT("/assign/checkin/:id", h.checkinAssign, auth.AuthorizedMobile())
	r.PUT("/assign/need-approval/:id", h.needApproval, auth.AuthorizedMobile())
	r.PUT("/assign/approve/:id", h.approveAssign, auth.AuthorizedMobile())
	r.PUT("/assign/reject/:id", h.rejectAssign, auth.AuthorizedMobile())

	// picker lead
	r.GET("/helper", h.getHelper, auth.AuthorizedMobile())
	r.PUT("/start-routing/:id", h.startAssignRouting, auth.AuthorizedMobile())
	r.PUT("/cancel-routing/:id", h.cancelAssignRouting, auth.AuthorizedMobile())
	r.POST("/assign/generateroute", h.generateRoute, auth.AuthorizedMobile())

	// picker
	r.GET("/assignment", h.readAssignment, auth.AuthorizedPickerMobile())
	r.GET("/getassignment/:id", h.readFirstAssignment, auth.AuthorizedPickerMobile())
	r.PUT("/pickeraction/:id", h.pickerAction, auth.AuthorizedPickerMobile())

	// checker
	r.GET("/checker/item/:id", h.detailPickingOrderItem, auth.AuthorizedMobile())
	r.PUT("/checker/checkin/:id", h.checkinChecker, auth.AuthorizedMobile())
	r.PUT("/checker/approve/:id", h.approveByChecker, auth.AuthorizedMobile())
	r.PUT("/checker/reject/:id", h.rejectByChecker, auth.AuthorizedMobile())
	r.PUT("/checker/scan_update/:id", h.scanUpdate, auth.AuthorizedMobile())

	r.PUT("/assign/update-bulk/pl/:id", h.updateBulkQty, auth.AuthorizedMobile())
	r.POST("/assign/profile", h.profile, auth.AuthorizedMobile())
	r.POST("/assign/list-products", h.listProduct, auth.AuthorizedMobile())
	r.POST("/assign/list", h.listPicking, auth.AuthorizedMobile())
	r.POST("/assign/group-so", h.getSalesOrderGroup, auth.AuthorizedMobile())
	r.POST("/assign/checkin-bulk", h.checkinBulkAssign, auth.AuthorizedMobile())

	// SPV
	r.POST("/assign/monitoring", h.listWRTMonitoring, auth.AuthorizedMobile())
	r.POST("/assign/list-helper", h.listHelper, auth.AuthorizedMobile())
	r.POST("/assign/so-monitoring", h.listSOMonitoring, auth.AuthorizedMobile())

	// produce kafka
	r.POST("/consume_create", h.createPickingWithKafka, auth.Authorized("pco_gen_pl"))

}

// read : function to get requested data based on parameters
func (h *Handler) read(c echo.Context) (e error) {
	ctx := c.(*cuxs.Context)
	rq := ctx.RequestQuery()

	var data []*model.PickingOrder
	var total int64

	if data, total, e = repository.GetPickingOrders(rq); e == nil {
		ctx.Data(data, total)
	}

	return ctx.Serve(e)
}

// detail : function to get detailed data by id
func (h *Handler) detail(c echo.Context) (e error) {
	ctx := c.(*cuxs.Context)

	var id int64
	if id, e = ctx.Decrypt("id"); e == nil {
		if ctx.ResponseData, e = repository.GetPickingOrder("id", id); e != nil {
			e = echo.ErrNotFound
		}
	}

	return ctx.Serve(e)
}

// detailPickingOrderItem : function to get detailed data by id
func (h *Handler) detailPickingOrderItem(c echo.Context) (e error) {
	ctx := c.(*cuxs.Context)
	m := new(model.PickingOrderItem)

	var id int64
	if id, e = ctx.Decrypt("id"); e != nil {
		return ctx.Serve(e)
	}
	if m, e = repository.GetPickingOrderItem("id", id); e != nil {
		e = echo.ErrNotFound
		return ctx.Serve(e)
	}
	if m.Product.Packability != 1 {
		e = errors.New("Product is not packable.")
		return ctx.Serve(e)
	}
	ctx.ResponseData = m

	return ctx.Serve(e)
}

// create : function to create new data based on input
func (h *Handler) create(c echo.Context) (e error) {
	ctx := c.(*cuxs.Context)
	var r createRequest

	if r.Session, e = auth.UserSession(ctx); e == nil {
		if e = ctx.Bind(&r); e == nil {
			ctx.ResponseData, e = Save(r)
		}
	}

	return ctx.Serve(e)
}

// readAssign : function to get requested data based on parameters
func (h *Handler) readAssign(c echo.Context) (e error) {
	ctx := c.(*cuxs.Context)
	rq := ctx.RequestQuery()

	var data []*model.PickingOrderAssign
	var total int64

	if data, total, e = repository.GetPickingOrderAssigns(rq); e == nil {
		ctx.Data(data, total)
	}

	return ctx.Serve(e)
}

// readNewAssign : function to get requested data based on parameters
func (h *Handler) readNewAssign(c echo.Context) (e error) {
	ctx := c.(*cuxs.Context)
	rq := ctx.RequestQuery()

	var data []*model.PickingOrderAssign
	var total int64

	if data, total, e = repository.GetPickingOrderAssigns(rq); e == nil {
		ctx.Data(data, total)
	}

	return ctx.Serve(e)
}

// assignDetail : function to get detailed data by id
func (h *Handler) assignDetail(c echo.Context) (e error) {
	ctx := c.(*cuxs.Context)

	var id int64
	if id, e = ctx.Decrypt("id"); e == nil {
		if ctx.ResponseData, e = repository.GetPickingOrderAssign("id", id); e != nil {
			e = echo.ErrNotFound
		}
	}

	return ctx.Serve(e)
}

// updateAssign : function to update picking order item qty
func (h *Handler) updateAssign(c echo.Context) (e error) {
	ctx := c.(*cuxs.Context)
	var r updateRequestAssign

	if r.Session, e = auth.UserSession(ctx); e == nil {

		if r.ID, e = ctx.Decrypt("id"); e == nil {
			if e = ctx.Bind(&r); e == nil {
				// for update qty update checker
				if r.TypeRequest == "1" {
					ctx.ResponseData, e = UpdateChecker(r)
				} else {
					ctx.ResponseData, e = UpdateAssign(r)
				}
			}
		}
	}

	return ctx.Serve(e)
}

// checkoutAssign : function to checkout picking order assign
func (h *Handler) checkoutAssign(c echo.Context) (e error) {
	ctx := c.(*cuxs.Context)
	var r checkoutRequestAssign

	if r.Session, e = auth.UserSession(ctx); e == nil {

		if r.ID, e = ctx.Decrypt("id"); e == nil {
			if e = ctx.Bind(&r); e == nil {
				ctx.ResponseData, e = CheckoutAssign(r)
			}
		}
	}

	return ctx.Serve(e)
}

// checkoutAssign : function to checkout picking order assign
func (h *Handler) checkinAssign(c echo.Context) (e error) {
	ctx := c.(*cuxs.Context)
	var r checkinRequestAssign

	if r.Session, e = auth.UserSession(ctx); e == nil {

		if r.ID, e = ctx.Decrypt("id"); e == nil {
			if e = ctx.Bind(&r); e == nil {
				ctx.ResponseData, e = CheckinAssign(r)
			}
		}
	}

	return ctx.Serve(e)
}

// assignLeadPicker : function to assign lead picker to picking list
func (h *Handler) assignLeadPicker(c echo.Context) (e error) {
	ctx := c.(*cuxs.Context)
	var r assignLeadPickerRequest

	if r.Session, e = auth.UserSession(ctx); e != nil {
		return ctx.Serve(e)
	}

	if r.ID, e = ctx.Decrypt("id"); e != nil {
		return ctx.Serve(e)
	}

	if e = ctx.Bind(&r); e != nil {
		return ctx.Serve(e)
	}

	ctx.ResponseData, e = AssignLeadPicker(r)

	return ctx.Serve(e)
}

// checkinBulkAssign : function to checkin bulk picking order assign
func (h *Handler) checkinBulkAssign(c echo.Context) (e error) {
	ctx := c.(*cuxs.Context)
	var r checkinBulkRequestAssign

	if r.Session, e = auth.UserSession(ctx); e == nil {

		if e = ctx.Bind(&r); e == nil {
			ctx.ResponseData, e = CheckinBulkAssign(r)
		}
	}

	return ctx.Serve(e)
}

// checkinChecker : function to checkout picking order assign
func (h *Handler) checkinChecker(c echo.Context) (e error) {
	ctx := c.(*cuxs.Context)
	var r checkinRequestChecker

	if r.Session, e = auth.UserSession(ctx); e == nil {

		if r.ID, e = ctx.Decrypt("id"); e == nil {
			if e = ctx.Bind(&r); e == nil {
				ctx.ResponseData, e = CheckinChecker(r)
			}
		}
	}

	return ctx.Serve(e)
}

// itemByPackingAssign : function to get detailed data by id
func (h *Handler) itemByPackingAssign(c echo.Context) (e error) {
	ctx := c.(*cuxs.Context)
	var id int64

	if id, e = ctx.Decrypt("id"); e == nil {
		if ctx.ResponseData, e = repository.GetItemByPickingAssignId("id", id); e != nil {
			e = echo.ErrNotFound
		}
	}

	return ctx.Serve(e)
}

// readItems : function to get requested data based on parameters
func (h *Handler) readItems(c echo.Context) (e error) {
	ctx := c.(*cuxs.Context)
	rq := ctx.RequestQuery()

	var data []*model.PickingOrderItem
	var total int64

	if data, total, e = repository.GetPickingOrderItems(rq); e == nil {
		ctx.Data(data, total)
	}

	return ctx.Serve(e)
}

// needApproval : function to request approval from spv
func (h *Handler) needApproval(c echo.Context) (e error) {
	ctx := c.(*cuxs.Context)
	var r updateRequestApprovalAssign

	if r.Session, e = auth.UserSession(ctx); e == nil {

		if r.ID, e = ctx.Decrypt("id"); e == nil {
			if e = ctx.Bind(&r); e == nil {
				ctx.ResponseData, e = RequestApproval(r)
			}
		}
	}

	return ctx.Serve(e)
}

// approveAssign : function to approve picking order assign
func (h *Handler) approveAssign(c echo.Context) (e error) {
	ctx := c.(*cuxs.Context)
	var r approveRequestAssign

	if r.Session, e = auth.UserSession(ctx); e == nil {

		if r.ID, e = ctx.Decrypt("id"); e == nil {
			if e = ctx.Bind(&r); e == nil {
				ctx.ResponseData, e = Approve(r)
			}
		}
	}

	return ctx.Serve(e)
}

// rejectAssign : function to reject picking order assign
func (h *Handler) rejectAssign(c echo.Context) (e error) {
	ctx := c.(*cuxs.Context)
	var r rejectRequestAssign

	if r.Session, e = auth.UserSession(ctx); e == nil {

		if r.ID, e = ctx.Decrypt("id"); e == nil {
			if e = ctx.Bind(&r); e == nil {
				ctx.ResponseData, e = Reject(r)
			}
		}
	}

	return ctx.Serve(e)
}

// function to get Helper for mobile
func (h *Handler) getHelper(c echo.Context) (e error) {
	ctx := c.(*cuxs.Context)

	pickingListIDStr := ctx.QueryParam("picking_list_id")
	pickingListID, e := common.Decrypt(pickingListIDStr)
	if e != nil {
		return ctx.Serve(e)
	}

	staff, _, e := repository.GetHelpers(ctx.RequestQuery())
	if e != nil {
		return ctx.Serve(e)
	}

	ctx.ResponseData, e = StaffUsedInformation(pickingListID, staff)
	return ctx.Serve(e)
}

// startAssignRouting : function to start picking routing
func (h *Handler) startAssignRouting(c echo.Context) (e error) {
	ctx := c.(*cuxs.Context)
	var r startRoutingAssignment

	if r.Session, e = auth.UserSession(ctx); e != nil {
		return ctx.Serve(e)
	}

	if r.ID, e = ctx.Decrypt("id"); e != nil {
		return ctx.Serve(e)
	}

	if e = ctx.Bind(&r); e != nil {
		return ctx.Serve(e)
	}

	ctx.ResponseData, e = StartAssignRouting(r)

	return ctx.Serve(e)
}

// cancelAssignRouting : function to cancel picking routing
func (h *Handler) cancelAssignRouting(c echo.Context) (e error) {
	ctx := c.(*cuxs.Context)
	var r cancelRoutingAssignment

	if r.Session, e = auth.UserSession(ctx); e != nil {
		return ctx.Serve(e)
	}

	if r.ID, e = ctx.Decrypt("id"); e != nil {
		return ctx.Serve(e)
	}

	if e = ctx.Bind(&r); e != nil {
		return ctx.Serve(e)
	}

	ctx.ResponseData, e = CancelAssignRouting(r)

	return ctx.Serve(e)
}

// readAssign : function to get requested data based on parameters
func (h *Handler) readAssignment(c echo.Context) (e error) {
	ctx := c.(*cuxs.Context)
	rq := ctx.RequestQuery()

	var data []*model.PickingRoutingStep

	if data, _, e = repository.GetPickingRoutingSteps(rq); e != nil {
		return ctx.Serve(e)
	}

	ctx.Data(data)
	return ctx.Serve(e)
}

// readFirstAssignment : function to get the first picking routing step with status in progress
func (h *Handler) readFirstAssignment(c echo.Context) (e error) {
	ctx := c.(*cuxs.Context)
	var data *model.PickingRoutingStep
	var id int64
	var session *auth.SessionData

	if session, e = auth.UserSession(ctx); e != nil {
		return ctx.Serve(e)
	}

	if id, e = ctx.Decrypt("id"); e != nil {
		return ctx.Serve(e)
	}

	if data, e = repository.GetFirstPickingRoutingStep(session.Staff.ID, id); e != nil {
		return ctx.Serve(e)
	}

	ctx.Data(data)
	return ctx.Serve(e)
}

// pickerAction : function to get the first picking routing step with status in progress
func (h *Handler) pickerAction(c echo.Context) (e error) {
	ctx := c.(*cuxs.Context)
	var r requestPickerAction

	if r.Session, e = auth.UserSession(ctx); e != nil {
		return ctx.Serve(e)
	}

	if r.ID, e = ctx.Decrypt("id"); e != nil {
		return ctx.Serve(e)
	}

	if e = ctx.Bind(&r); e != nil {
		return ctx.Serve(e)
	}

	ctx.ResponseData, e = PickerAction(r)
	return ctx.Serve(e)
}

// approveByChecker : function to approve picking order assign by checker
func (h *Handler) approveByChecker(c echo.Context) (e error) {
	ctx := c.(*cuxs.Context)
	var r approveRequestChecker

	if r.Session, e = auth.UserSession(ctx); e == nil {

		if r.ID, e = ctx.Decrypt("id"); e == nil {
			if e = ctx.Bind(&r); e == nil {
				ctx.ResponseData, e = ApproveByChecker(r)
			}
		}
	}

	return ctx.Serve(e)
}

// rejectByChecker : function to reject picking order assign by checker
func (h *Handler) rejectByChecker(c echo.Context) (e error) {
	ctx := c.(*cuxs.Context)
	var r rejectRequestChecker

	if r.Session, e = auth.UserSession(ctx); e == nil {

		if r.ID, e = ctx.Decrypt("id"); e == nil {
			if e = ctx.Bind(&r); e == nil {
				ctx.ResponseData, e = RejectByChecker(r)
			}
		}
	}

	return ctx.Serve(e)
}

// scanUpdate : function to update check qty from scanning barcode
func (h *Handler) scanUpdate(c echo.Context) (e error) {
	ctx := c.(*cuxs.Context)
	var r scanRequestChecker

	if r.Session, e = auth.UserSession(ctx); e != nil {
		return ctx.Serve(e)
	}
	if r.ID, e = ctx.Decrypt("id"); e != nil {
		return ctx.Serve(e)
	}
	if e = ctx.Bind(&r); e != nil {
		return ctx.Serve(e)
	}
	ctx.ResponseData, e = UpdateCheckQtyByScan(r)

	return ctx.Serve(e)
}

// pickingOrderTemplate : function to get picking order template data
func (h *Handler) pickingOrderTemplate(c echo.Context) (e error) {
	ctx := c.(*cuxs.Context)
	rq := ctx.RequestQuery()
	var backdate time.Time
	isExport := ctx.QueryParam("export") == "1"
	backdate = now.NewParse(time.RFC3339, ctx.QueryParam("date")).Time

	param := ctx.QueryParams()
	warehouseID, _ := common.Decrypt(param.Get("warehouse_id"))

	warehouse, _ := repository.GetWarehouse("id", warehouseID)

	deliveryDate := ctx.QueryParam("delivery_dates")

	cond := make(map[string]interface{})

	if deliveryDate != "" {
		cond["so.delivery_date = "] = deliveryDate
	}

	if warehouseID != 0 {
		cond["so.warehouse_id = "] = warehouseID
	}

	data, e := getPickingOrder(rq, cond)
	if e == nil {
		if isExport {
			var file string
			if file, e = GetPickingOrderXls(backdate, data, warehouse); e == nil {
				ctx.Files(file)
			}
		} else {
			ctx.Data(data)
		}
	}

	return ctx.Serve(e)
}

// uploadAssign : function to create new data based on input
func (h *Handler) uploadAssign(c echo.Context) (e error) {
	ctx := c.(*cuxs.Context)
	var r uploadAssignRequest

	if r.Session, e = auth.UserSession(ctx); e == nil {
		if e = ctx.Bind(&r); e == nil {
			if len(r.DataCorrection) != 0 {

				var backdate time.Time
				backdate = now.NewParse(time.RFC3339, ctx.QueryParam("date")).Time
				warehouse, _ := repository.GetWarehouse("id", r.Warehouse.ID)

				var file string
				if file, e = RegeneratePickingOrderXls(backdate, r, warehouse); e == nil {
					ctx.Files(file)
				}
			} else {
				ctx.ResponseData, e = UploadAssign(r)
			}
		}
	}

	return ctx.Serve(e)
}

func (h *Handler) receivePrint(c echo.Context) (e error) {
	ctx := c.(*cuxs.Context)
	var poa *model.PickingOrderAssign
	var id int64
	req := make(map[string]interface{})
	if _, e = auth.UserSession(ctx); e == nil {
		if id, e = ctx.Decrypt("id"); e == nil {
			if poa, e = repository.GetPickingOrderAssign("id", id); e != nil {
				e = echo.ErrNotFound
			} else {
				req["pl"] = poa

				file := util.SendPrint(req, "read/picking_label")
				ctx.Files(file)

			}
		}
	}

	return ctx.Serve(e)
}

// printLabels : function to get requested data based on parameters
func (h *Handler) printLabels(c echo.Context) (e error) {
	ctx := c.(*cuxs.Context)
	rq := ctx.RequestQuery()
	typePrint := ctx.QueryParam("type_print")

	var dataSalesInvoice []*model.SalesInvoice

	var session *auth.SessionData
	configs := make(map[string]string)
	req := make(map[string]interface{})

	switch typePrint {
	case "sj":
		if session, e = auth.UserSession(ctx); e != nil {
			return ctx.Serve(e)
		}
		if dataSalesInvoice, _, e = repository.GetSalesInvoiceForPrints(rq); e != nil {
			e = echo.NewHTTPError(http.StatusUnprocessableEntity)
		} else {
			// if invoice record not found, so it's gonna try to find do record
			if len(dataSalesInvoice) == 0 {
				var file string
				if e, file = GetDeliveryOrderRecordForPrint(session, rq); e != nil {
					return ctx.Serve(e)
				}

				ctx.Files(file)
			} else {
				if e = GetDeliveryKoliesBySalesOrderInSalesInvoice(dataSalesInvoice[0]); e != nil {
					e = echo.NewHTTPError(http.StatusUnprocessableEntity)
					return ctx.Serve(e)
				}
				req["si"] = dataSalesInvoice[0]
				req["session"] = session.Staff.ID + 56

				var config []*model.ConfigApp
				if config, _, e = repository.GetConfigAppsByAttribute("attribute__icontains", "company"); e != nil {
					e = echo.NewHTTPError(http.StatusUnprocessableEntity)
					return ctx.Serve(e)
				}

				for _, v := range config {
					configs[strings.TrimPrefix(v.Attribute, "company_")] = v.Value
				}
				configs["address"] = strings.ReplaceAll(configs["address"], "<br>", "\n")
				req["config"] = configs

				file := util.SendPrint(req, "read/si")
				ctx.Files(file)

				// delta print
				dataSalesInvoice[0].DeltaPrint = dataSalesInvoice[0].DeltaPrint + 1
				dataSalesInvoice[0].Save("DeltaPrint")
			}
		}
		return ctx.Serve(e)

	case "label_picking":
		var data []*model.PickingOrderAssign
		var r ResponsePrint

		if data, _, e = repository.GetPickingOrderAssigns(rq); e != nil {
			return ctx.Serve(echo.NewHTTPError(http.StatusUnprocessableEntity))
		}
		if len(data) == 0 {
			return ctx.Serve(echo.NewHTTPError(http.StatusUnprocessableEntity))
		}

		req["pls"] = data[0]
		file := util.SendPrint(req, "read/picking_print")
		r.LinkPrint = file
		r.TotalKoli = data[0].TotalKoli
		ctx.ResponseData = r

		// print label
		e = UpdatePrintLabel(data[0].SalesOrder.ID)

		return ctx.Serve(e)

	case "picking_list":
		var data []*model.SalesOrder
		var pl *model.PickingList

		if data, pl, e = repository.GetSalesOrderbyPickingList(rq); e != nil {
			return ctx.Serve(echo.NewHTTPError(http.StatusUnprocessableEntity))
		}
		if len(data) == 0 {
			return ctx.Serve(echo.NewHTTPError(http.StatusUnprocessableEntity))
		}

		request := &SalesOrderByPickingList{
			PickingListCode: pl.Code,
			SalesOrders:     data,
		}

		req["plso"] = request
		file := util.SendPrint(req, "read/picking_list_print")
		ctx.Files(file)

		return ctx.Serve(e)

	default:
		var r ResponsePrint
		var data []*model.PickingOrderAssign

		if data, _, e = repository.GetPickingOrderAssigns(rq); e != nil {
			return ctx.Serve(echo.NewHTTPError(http.StatusUnprocessableEntity))
		}
		if len(data) == 0 {
			return ctx.Serve(echo.NewHTTPError(http.StatusUnprocessableEntity))
		}

		req["pls"] = data[0]
		file := util.SendPrint(req, "read/picking_print")
		r.LinkPrint = file
		r.TotalKoli = data[0].TotalKoli
		ctx.ResponseData = r

		return ctx.Serve(e)
	}
}

func GetDeliveryOrderRecordForPrint(session *auth.SessionData, rq *orm.RequestQuery) (e error, file string) {
	configs := make(map[string]string)
	req := make(map[string]interface{})

	var dataDeliveryOrder []*model.DeliveryOrder
	if dataDeliveryOrder, _, e = repository.GetDeliveryOrdersForPrint(rq); e != nil {
		e = echo.NewHTTPError(http.StatusUnprocessableEntity)
		return e, ""
	}

	// if do record not found, so it's gonna return not found
	if len(dataDeliveryOrder) == 0 {
		e = echo.NewHTTPError(http.StatusUnprocessableEntity)
		return e, ""
	}

	if e = GetDeliveryKoliesBySalesOrderInDeliveryOrder(dataDeliveryOrder[0]); e != nil {
		e = echo.NewHTTPError(http.StatusUnprocessableEntity)
		return e, ""
	}
	req["do"] = dataDeliveryOrder[0]
	req["session"] = session.Staff.ID + 56

	var config []*model.ConfigApp
	if config, _, e = repository.GetConfigAppsByAttribute("attribute__icontains", "company"); e != nil {
		e = echo.NewHTTPError(http.StatusUnprocessableEntity)
		return e, ""
	}

	for _, v := range config {
		configs[strings.TrimPrefix(v.Attribute, "company_")] = v.Value
	}
	configs["address"] = strings.ReplaceAll(configs["address"], "<br>", "\n")
	req["config"] = configs

	file = util.SendPrint(req, "read/do")

	// delta print
	dataDeliveryOrder[0].DeltaPrint = dataDeliveryOrder[0].DeltaPrint + 1
	dataDeliveryOrder[0].Save("DeltaPrint")

	return e, file
}

// profile : function to get list task in mobile
func (h *Handler) profile(c echo.Context) (e error) {
	ctx := c.(*cuxs.Context)
	var r profileRequestPicking

	if r.Session, e = auth.UserSession(ctx); e == nil {
		data, _ := getProfile(r)

		ctx.Data(data)
	}

	return ctx.Serve(e)
}

// listWRTMonitoring : function to get list WRT with amount sales order
func (h *Handler) listWRTMonitoring(c echo.Context) (e error) {
	ctx := c.(*cuxs.Context)
	var r wrtMonitoringRequest

	if r.Session, e = auth.UserSession(ctx); e != nil {
		return ctx.Serve(e)
	}
	if e = ctx.Bind(&r); e != nil {
		return ctx.Serve(e)

	}
	data, _ := getWRTMonitorings(r)

	ctx.Data(data)

	return ctx.Serve(e)
}

// listHelper : function to get list staff or helper base on warehouse
func (h *Handler) listHelper(c echo.Context) (e error) {
	ctx := c.(*cuxs.Context)
	var r listHelperRequest

	if r.Session, e = auth.UserSession(ctx); e != nil {
		return ctx.Serve(e)
	}

	if e = ctx.Bind(&r); e != nil {
		return ctx.Serve(e)
	}

	data, _ := getListHelper(r)

	ctx.Data(data)

	return ctx.Serve(e)
}

// listSOMonitoring : function to get list detail from wrt monitoring
func (h *Handler) listSOMonitoring(c echo.Context) (e error) {
	ctx := c.(*cuxs.Context)
	var r wrtMonitoringRequest

	if r.Session, e = auth.UserSession(ctx); e != nil {
		return ctx.Serve(e)
	}
	if e = ctx.Bind(&r); e != nil {
		return ctx.Serve(e)
	}
	data, _ := getListSOMonitoring(r)

	ctx.Data(data)

	return ctx.Serve(e)
}

// listProduct : function to get list task in mobile
func (h *Handler) listProduct(c echo.Context) (e error) {
	ctx := c.(*cuxs.Context)
	var r listProductRequestPicking

	if r.Session, e = auth.UserSession(ctx); e != nil {
		return ctx.Serve(e)
	}
	if e = ctx.Bind(&r); e != nil {
		return ctx.Serve(e)

	}
	data, _ := getListProduct(r)
	ctx.Data(data)

	return ctx.Serve(e)
}

// listPicking : function to get list task in mobile
func (h *Handler) listPicking(c echo.Context) (e error) {
	ctx := c.(*cuxs.Context)
	var r listRequestPicking

	cond := make(map[string]interface{})

	if r.Session, e = auth.UserSession(ctx); e == nil {
		if e = ctx.Bind(&r); e == nil {
			deliveryDateArr := strings.Split(r.DeliveryDate, "|")
			cond["pl.delivery_date "] = deliveryDateArr
			cond["pl.warehouse_id = "] = r.Warehouse.ID
			if r.Query != "" {
				cond["query"] = r.Query
			}
			if r.FilterPickingList != "" {
				cond["filter_picking_list"] = r.FilterPickingList
			}
			if r.FilterPickingStatus != "" {
				cond["filter_picking_status"] = r.FilterPickingStatus
			}
			data, _ := getListPicking(r, cond)

			ctx.Data(data)
		}
	}

	return ctx.Serve(e)
}

// listPickingGenerate : function to generate picking list code
func (h *Handler) listPickingGenerate(c echo.Context) (e error) {
	ctx := c.(*cuxs.Context)
	var (
		r               generateCodePickingRequest
		respPickingList *model.PickingOrder
		d               sync.Mutex
		wg              sync.WaitGroup
	)

	if r.Session, e = auth.UserSession(ctx); e != nil {
		return ctx.Serve(e)
	}

	if e = ctx.Bind(&r); e != nil {
		// delete key that blocks picking list creation for associated warehouse
		dbredis.Redis.DeleteCache("picking_list_" + r.WarehouseID)
		return ctx.Serve(e)
	}

	data, _ := getPickingListGenCode(r)
	r.PickingListFinal = data

	ctx1 := context.Background()

	// lock unlock
	wg.Add(1)
	go func() {
		d.Lock()
		respPickingList, r.PickingListObj, e = InsertPickingListGenerateCode(r)
		d.Unlock()
		defer wg.Done()
	}()

	wg.Wait()

	if e != nil {
		// delete key that blocks picking list creation for associated warehouse
		dbredis.Redis.DeleteCache("picking_list_" + r.WarehouseID)
		return ctx.Serve(e)
	}

	r.TypeRequest = "create"
	jobs := &model.Jobs{
		EndpointUrl:    "/v1/warehouse/picking_order/consume_create",
		Topic:          env.GetString("KAFKA_TOPIC", ""),
		EndpointMethod: "POST",
		ResponseBody:   "[]",
		Status:         1,
		CreatedAt:      time.Now(),
		CreatedBy:      r.Session.Staff.User.ID,
		RetryCount:     0,
	}
	m := mongodb.NewMongo()

	m.CreateIndex("Jobs", "_id", true)
	jobs.ID = primitive.NewObjectID()
	r.JobsID = jobs.ID
	a, _ := json.Marshal(r)
	jobs.RequestBody = string(a)
	_, e = m.InsertOneData("Jobs", jobs)

	if e != nil {
		fmt.Println(e)
		m.DisconnectMongoClient()
	}
	jobsFilter := *jobs

	e = kafka.Produce(ctx1, jobs, jobs.Topic)
	if e != nil {
		jobs.ResponseBody = "{\"error_produce\":\"" + e.Error() + "\"}"
		jobs.Status = 5
		err := m.UpdateOneDataWithFilter("Jobs", jobsFilter, jobs)
		if err != nil {
			e = err
			fmt.Println(e)
			m.DisconnectMongoClient()
		}
	}

	respPickingList.JobsID = r.JobsID

	ctx.ResponseData = respPickingList

	return ctx.Serve(e)
}

// createPickingWithKafka : function to create new data based on input
func (h *Handler) createPickingWithKafka(c echo.Context) (e error) {
	ctx := c.(*cuxs.Context)
	var r generateCodePickingRequest
	if r.Session, e = auth.UserSession(ctx); e != nil {
		return e
	}
	if e = ctx.Bind(&r); e != nil {
		return e
	}

	jobs := &model.Jobs{ID: r.JobsID}
	jobsChanged := &model.Jobs{}
	m := mongodb.NewMongo()
	ret, err := m.GetOneDataWithFilter("Jobs", jobs)
	if err != nil {
		fmt.Println(err)
		m.DisconnectMongoClient()
	}
	json.Unmarshal(ret, &jobsChanged)
	//read from mongo here
	if jobsChanged.Status == 3 {
		ctx.ResponseData = http.StatusOK
	}

	if jobsChanged.Status == 2 {
		ctx.ResponseData, e = SaveOrderIntoAssign(r, r.PickingOrder)
		if e != nil {
			dbredis.Redis.DeleteCache("picking_list_" + r.WarehouseID)
			return ctx.Serve(e)
		}
		jobsChanged.Status = 3
		e = m.UpdateOneDataWithFilter("Jobs", jobs, jobsChanged)
		if e != nil {
			fmt.Println(e)
			m.DisconnectMongoClient()
		}

	}
	m.DisconnectMongoClient()
	return ctx.Serve(e)
}

// getSalesOrderGroup : function to generate picking list code
func (h *Handler) getSalesOrderGroup(c echo.Context) (e error) {
	ctx := c.(*cuxs.Context)
	var r groupingSalesOrderRequest

	if r.Session, e = auth.UserSession(ctx); e != nil {
		return ctx.Serve(e)
	}
	if e = ctx.Bind(&r); e != nil {
		return ctx.Serve(e)
	}
	data, _ := getSalesOrderGroupByPickingList(r)
	ctx.Data(data)

	return ctx.Serve(e)

}

func GetDeliveryKoliesBySalesOrderInDeliveryOrder(do *model.DeliveryOrder) error {
	var err error
	o := orm.NewOrm()
	o.Using("read_only")
	o.Raw("select * from delivery_koli dk WHERE dk.sales_order_id = ?", do.SalesOrder.ID).QueryRows(&do.DeliveryKoli)

	var totalKoli float64
	if len(do.DeliveryKoli) != 0 {
		for _, v := range do.DeliveryKoli {
			v.Koli.Read("ID")
			totalKoli += v.Quantity
		}

		do.TotalKoli = totalKoli
	}

	return err
}

func GetDeliveryKoliesBySalesOrderInSalesInvoice(si *model.SalesInvoice) error {
	var err error
	o := orm.NewOrm()
	o.Using("read_only")

	o.Raw("select * from delivery_koli dk WHERE dk.sales_order_id = ?", si.SalesOrder.ID).QueryRows(&si.DeliveryKoli)

	var totalKoli float64
	if len(si.DeliveryKoli) != 0 {
		for _, v := range si.DeliveryKoli {
			v.Koli.Read("ID")
			totalKoli += v.Quantity
		}

		si.TotalKoli = totalKoli
	}

	return err
}

// readAssign : function to get requested data based on parameters
func (h *Handler) getPickingLists(c echo.Context) (e error) {
	ctx := c.(*cuxs.Context)
	rq := ctx.RequestQuery()

	var data []*model.PickingList
	var total int64

	if data, total, e = repository.GetPickingLists(rq); e == nil {
		ctx.Data(data, total)
	}

	return ctx.Serve(e)
}

// updateBulkQty : function to update picking order item qty
func (h *Handler) updateBulkQty(c echo.Context) (e error) {
	ctx := c.(*cuxs.Context)
	var r updateBulkQtyRequest

	if r.Session, e = auth.UserSession(ctx); e != nil {
		return ctx.Serve(e)
	}
	if r.ID, e = ctx.Decrypt("id"); e != nil {
		return ctx.Serve(e)
	}
	if e = ctx.Bind(&r); e != nil {
		return ctx.Serve(e)
	}

	ctx.ResponseData, e = UpdateBulkQty(r)

	return ctx.Serve(e)
}

// generateRoute : function to send data to vroom
func (h *Handler) generateRoute(c echo.Context) (e error) {
	ctx := c.(*cuxs.Context)
	var r model.VroomRequest
	var session *auth.SessionData

	if session, e = auth.UserSession(ctx); e != nil {
		return ctx.Serve(e)
	}

	if e = ctx.Bind(&r); e != nil {
		return ctx.Serve(e)
	}

	ctx.ResponseData, e = PostVroom(&r, session)

	return ctx.Serve(e)
}

// summaryPickingRoute : function to show summary of sales orders affected based on the filter given
func (h *Handler) summaryPicking(c echo.Context) (e error) {
	ctx := c.(*cuxs.Context)
	var r summaryRequest
	if r.Session, e = auth.UserSession(ctx); e != nil {
		return ctx.Serve(e)
	}
	if e = ctx.Bind(&r); e != nil {
		return ctx.Serve(e)
	}
	data := pickingRouteSummary{
		TotalSalesOrder:              r.TotalSalesOrder,
		TotalWeight:                  r.TotalWeight,
		HighestSalesOrderWeight:      r.HighestSalesOrderWeight,
		HighestSalesOrderItemsWeight: r.HighestSalesOrderItemsWeight,
	}
	ctx.Data(data)

	return ctx.Serve(e)
}
