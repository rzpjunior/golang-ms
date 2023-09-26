// Copyright 2020 PT. Eden Pangan Indonesia. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

package order

import (
	"strconv"

	"git.edenfarm.id/project-version2/api/datastore/repository"
	"git.edenfarm.id/project-version2/api/src/auth"
	"git.edenfarm.id/project-version2/datamodel/model"

	"strings"

	"git.edenfarm.id/cuxs/common"
	"git.edenfarm.id/cuxs/cuxs"
	"git.edenfarm.id/project-version2/api/util"
	"github.com/labstack/echo/v4"
)

type Handler struct{}

// URLMapping : function to map url with it's handler and add authorization validation
func (h *Handler) URLMapping(r *echo.Group) {
	r.GET("", h.read, auth.Authorized("po_rdl"))
	r.GET("/:id", h.detail, auth.Authorized("po_rdd"))
	r.GET("/filter", h.readFilter, auth.Authorized("filter_rdl"))
	r.GET("/print/:id", h.receivePrint, auth.Authorized("po_prt"))
	r.GET("/download-delivery/:id", h.downloadDelivery, auth.Authorized("pd_dl"))
	r.POST("", h.create, auth.Authorized("po_crt"))
	r.POST("/market-purchase/:id", h.addMarketPurchase, auth.Authorized("po_mrk"))
	r.PUT("/:id", h.update, auth.Authorized("po_upd"))
	r.PUT("/cancel/:id", h.cancel, auth.Authorized("po_can"))
	r.PUT("/commit/:id", h.commit, auth.Authorized("po_cmt"))
	r.PUT("/commit", h.bulkCommit, auth.Authorized("po_cmt"))
	r.PUT("/update-product/:id", h.updateProduct, auth.Authorized("po_upd_prd"))
	r.PUT("/lock/:id", h.lock, auth.Authorized("po_upd"))
	r.PUT("/assign/:id", h.assign, auth.AuthorizedFieldPurchaserMobile())
	r.PUT("/count_print/:id", h.countPrint, auth.AuthorizedFieldPurchaserMobile())
	r.POST("/signature", h.signPurchaseOrder, auth.AuthorizedFieldPurchaserMobile())
	r.GET("/filter/consolidated_shipment", h.filterForConsolidatedShipment, auth.AuthorizedFieldPurchaserMobile())
}

// receivePrint : function to print
func (h *Handler) receivePrint(c echo.Context) (e error) {
	ctx := c.(*cuxs.Context)
	var po *model.PurchaseOrder
	var id int64
	configs := make(map[string]string)
	req := make(map[string]interface{})

	if id, e = ctx.Decrypt("id"); e == nil {
		if po, e = repository.GetPurchaseOrder("id", id); e != nil {
			e = echo.ErrNotFound
		} else {
			req["po"] = po
			if config, _, e := repository.GetConfigAppsByAttribute("attribute__icontains", "company"); e == nil {
				for _, v := range config {
					configs[strings.TrimPrefix(v.Attribute, "company_")] = v.Value
				}
				configs["address"] = strings.ReplaceAll(configs["address"], "<br>", "\n")
				req["config"] = configs
			} else {
				e = echo.ErrNotFound
			}

			file := util.SendPrint(req, "read/po")
			ctx.Files(file)
		}
	}

	return ctx.Serve(e)
}

// read : function to get requested data based on parameters
func (h *Handler) read(c echo.Context) (e error) {
	ctx := c.(*cuxs.Context)
	rq := ctx.RequestQuery()

	isInbound := ctx.QueryParam("is_inbound")

	var data []*model.PurchaseOrder
	var total int64

	if data, total, e = repository.GetPurchaseOrders(rq, isInbound); e == nil {
		ctx.Data(data, total)
	}

	return ctx.Serve(e)
}

// readFilter : function to get requested data based on parameters with filtered permission
func (h *Handler) readFilter(c echo.Context) (e error) {
	ctx := c.(*cuxs.Context)
	rq := ctx.RequestQuery()

	var data []*model.PurchaseOrder
	var total int64

	if data, total, e = repository.GetFilterPurchaseOrders(rq); e == nil {
		ctx.Data(data, total)
	}

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

// detail : function to get detailed data by id
func (h *Handler) detail(c echo.Context) (e error) {
	ctx := c.(*cuxs.Context)

	var id int64
	if id, e = ctx.Decrypt("id"); e == nil {
		if ctx.ResponseData, e = repository.GetPurchaseOrder("id", id); e != nil {
			e = echo.ErrNotFound
		}
	}

	return ctx.Serve(e)
}

// update : function to update requested data based on parameters
func (h *Handler) update(c echo.Context) (e error) {
	ctx := c.(*cuxs.Context)
	var r updateRequest

	if r.Session, e = auth.UserSession(ctx); e == nil {
		if r.ID, e = ctx.Decrypt("id"); e == nil {
			if e = ctx.Bind(&r); e == nil {
				ctx.ResponseData, e = Update(r)
			}
		}
	}

	return ctx.Serve(e)
}

// cancel : function to cancel purchase order
func (h *Handler) cancel(c echo.Context) (e error) {
	ctx := c.(*cuxs.Context)
	var r cancelRequest
	r.PurchaseOrder = new(model.PurchaseOrder)
	if r.PurchaseOrder.ID, e = common.Decrypt(ctx.Param("id")); e == nil {
		if r.PurchaseOrder, e = repository.GetPurchaseOrder("id", r.PurchaseOrder.ID); e == nil {
			if r.Session, e = auth.UserSession(ctx); e == nil {
				if e = ctx.Bind(&r); e == nil {
					ctx.ResponseData, e = r.Cancel()
				}
			}
		}
	}
	return ctx.Serve(e)
}

// commit : function to commit purchase order
func (h *Handler) commit(c echo.Context) (e error) {
	ctx := c.(*cuxs.Context)
	var r commitRequest
	r.PurchaseOrder = new(model.PurchaseOrder)
	if r.PurchaseOrder.ID, e = common.Decrypt(ctx.Param("id")); e == nil {
		if r.PurchaseOrder, e = repository.GetPurchaseOrder("id", r.PurchaseOrder.ID); e == nil {
			if r.Session, e = auth.UserSession(ctx); e == nil {
				if e = ctx.Bind(&r); e == nil {
					ctx.ResponseData, e = r.Commit()
				}
			}
		}
	}
	return ctx.Serve(e)
}

// bulkCommit : function to commit several po at the same time
func (h *Handler) bulkCommit(c echo.Context) (e error) {
	ctx := c.(*cuxs.Context)
	var r bulkCommitRequest

	if r.Session, e = auth.UserSession(ctx); e == nil {

		if e = ctx.Bind(&r); e == nil {
			ctx.ResponseData = strconv.Itoa(int(r.SuccessCount)) + " of " + strconv.Itoa(len(r.PurchaseOrderIDs)) + " data has been saved successfully"
		}
	}

	return ctx.Serve(e)
}

// updateProduct : function to update product only requested data based on parameters
func (h *Handler) updateProduct(c echo.Context) (e error) {
	ctx := c.(*cuxs.Context)
	var r updateProductRequest

	if r.Session, e = auth.UserSession(ctx); e == nil {
		if r.ID, e = ctx.Decrypt("id"); e == nil {
			if e = ctx.Bind(&r); e == nil {
				ctx.ResponseData, e = UpdateProduct(r)
			}
		}
	}

	return ctx.Serve(e)
}

// downloadDelivery : function to download delivery slip
func (h *Handler) downloadDelivery(c echo.Context) (e error) {
	ctx := c.(*cuxs.Context)
	isExport := ctx.QueryParam("export") == "1"

	var (
		id            int64
		purchaseOrder *model.PurchaseOrder
	)

	if id, e = ctx.Decrypt("id"); e == nil {
		if purchaseOrder, e = repository.GetPurchaseOrder("id", id); e == nil {
			if isExport {
				var file string
				if file, e = PrintDeliverySlipXls(purchaseOrder); e == nil {
					ctx.Files(file)
				}
			} else {
				ctx.Data(purchaseOrder)
			}
		}
	}

	return ctx.Serve(e)
}

// addMarketPurchase : function to add market purchase
func (h *Handler) addMarketPurchase(c echo.Context) (e error) {
	ctx := c.(*cuxs.Context)
	var r marketPurchaseRequest

	if r.Session, e = auth.UserSession(ctx); e == nil {
		if r.ID, e = ctx.Decrypt("id"); e == nil {
			if e = ctx.Bind(&r); e == nil {
				ctx.ResponseData, e = AddMarketPurchase(r)
			}
		}
	}

	return ctx.Serve(e)
}

// assign purchase order to field purchaser
func (h *Handler) assign(c echo.Context) (e error) {
	ctx := c.(*cuxs.Context)
	var r assignRequest

	if r.Session, e = auth.UserSession(ctx); e != nil {
		return ctx.Serve(e)
	}

	if r.ID, e = ctx.Decrypt("id"); e != nil {
		return ctx.Serve(e)
	}

	if e = ctx.Bind(&r); e != nil {
		return ctx.Serve(e)
	}

	ctx.ResponseData, e = Assign(r)

	return ctx.Serve(e)
}

// lock : function to lock GR
func (h *Handler) lock(c echo.Context) (e error) {
	ctx := c.(*cuxs.Context)
	var r lockRequest

	if r.Session, e = auth.UserSession(ctx); e != nil {
		return ctx.Serve(e)
	}
	if r.ID, e = ctx.Decrypt("id"); e != nil {
		return ctx.Serve(e)
	}
	if e = ctx.Bind(&r); e != nil {
		return ctx.Serve(e)
	}
	ctx.ResponseData, e = Lock(r)

	return ctx.Serve(e)
}

// count number of print copy of purchase order
func (h *Handler) countPrint(c echo.Context) (e error) {
	ctx := c.(*cuxs.Context)
	var r countPrintRequest

	if r.Session, e = auth.UserSession(ctx); e != nil {
		return ctx.Serve(e)
	}

	if r.ID, e = ctx.Decrypt("id"); e != nil {
		return ctx.Serve(e)
	}
	if e = ctx.Bind(&r); e != nil {
		return ctx.Serve(e)
	}

	ctx.ResponseData, e = CountPrint(r)

	return ctx.Serve(e)
}

// signPurchaseOrder : function to sign purchase order
func (h *Handler) signPurchaseOrder(c echo.Context) (e error) {
	ctx := c.(*cuxs.Context)
	var r signRequest

	if r.Session, e = auth.UserSession(ctx); e != nil {
		return ctx.Serve(e)
	}

	if e = ctx.Bind(&r); e != nil {
		return ctx.Serve(e)
	}

	ctx.ResponseData, e = Sign(r)

	return ctx.Serve(e)
}

// filterForConsolidatedShipment : function to get requested data based on parameters with filtered permission
func (h *Handler) filterForConsolidatedShipment(c echo.Context) (e error) {
	ctx := c.(*cuxs.Context)
	rq := ctx.RequestQuery()

	var data []*model.PurchaseOrder
	var total int64

	if data, total, e = repository.GetFilterPurchaseOrdersForConsolidatedShipment(rq); e != nil {
		return ctx.Serve(e)
	}

	ctx.Data(data, total)

	return ctx.Serve(e)
}
