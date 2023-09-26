// Copyright 2020 PT. Eden Pangan Indonesia. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

package packing

import (
	"encoding/json"
	"reflect"
	"strings"
	"time"

	"git.edenfarm.id/cuxs/mongodb"
	"git.edenfarm.id/cuxs/orm"
	"git.edenfarm.id/project-version2/api/util"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo/options"

	"git.edenfarm.id/cuxs/common/now"
	"git.edenfarm.id/project-version2/api/datastore/repository"
	"git.edenfarm.id/project-version2/api/src/auth"
	"git.edenfarm.id/project-version2/datamodel/model"

	"git.edenfarm.id/cuxs/common"
	"git.edenfarm.id/cuxs/cuxs"
	"github.com/labstack/echo/v4"
)

type Handler struct{}

// URLMapping : function to map url with it's handler and add authorization validation
func (h *Handler) URLMapping(r *echo.Group) {
	r.GET("", h.read, auth.Authorized("pc_rdl"))
	r.GET("/:id", h.detail, auth.Authorized("pc_rdd"))
	r.GET("/filter", h.readFilter, auth.Authorized("filter_rdl"))
	r.GET("/item_assign", h.itemAssign, auth.Authorized("pc_rdd"))
	r.GET("/get_soi", h.getSoi, auth.Authorized("pc_rdd"))
	r.POST("", h.create, auth.Authorized("pc_crt"))
	r.PUT("/upload/:id", h.uploadActualPackingPacking, auth.Authorized("pc_upl"))
	r.GET("/template/:id", h.template, auth.Authorized("pc_dl"))
	r.PUT("/confirm/:id", h.confirm, auth.Authorized("pc_cnf"))
	r.PUT("/cancel/:id", h.cancel, auth.Authorized("pc_can"))
	r.PUT("/item_assign/:id", h.itemAssignUpdate, auth.Authorized("pc_asg_pac"))

	// packing recommendation
	r.GET("/recommendation", h.read, auth.Authorized("pc_rdl"))
	r.GET("/recommendation/:id", h.detail, auth.Authorized("pc_rdd"))
	r.GET("/recommendation/export/:id", h.exportForm, auth.Authorized("pc_exp"))
	r.GET("/recommendation/pack", h.ListPack, auth.Authorized("pc_rdd"))
	r.PUT("/recommendation/print/:id", h.receivedPrint, auth.Authorized("pc_prt"))
	r.PUT("/recommendation/detail/:id", h.detailPack, auth.Authorized("pc_rdd"))
	r.PUT("/recommendation/update/:id", h.updatePackingPack, auth.Authorized("pc_upd"))
	r.PUT("/recommendation/dispose/:id", h.disposePackingPack, auth.Authorized("pc_del"))
	r.POST("/recommendation/generate", h.generatePacking, auth.Authorized("pc_crt"))

	//PACKING ORDER
	r.PUT("/assign_qty/:id", h.assignQuantityPacking, auth.AuthorizedMobile())
}

// read : function to get requested data based on parameters
func (h *Handler) read(c echo.Context) (e error) {
	ctx := c.(*cuxs.Context)
	rq := ctx.RequestQuery()

	var data []*model.PackingOrder
	var total int64

	if data, total, e = repository.GetPackingOrders(rq); e == nil {
		ctx.Data(data, total)
	}

	return ctx.Serve(e)
}

// detail : function to get detailed data by id
func (h *Handler) detail(c echo.Context) (e error) {
	ctx := c.(*cuxs.Context)

	var id int64
	if id, e = ctx.Decrypt("id"); e == nil {
		if ctx.ResponseData, e = repository.GetPackingOrderDetailPack("id", id); e != nil {
			e = echo.ErrNotFound
		}
	}

	return ctx.Serve(e)
}

// export : function to get excel data
func (h *Handler) exportForm(c echo.Context) (e error) {
	ctx := c.(*cuxs.Context)

	var s *auth.SessionData
	if s, e = auth.UserSession(ctx); e != nil {
		return ctx.Serve(e)
	}
	var id int64
	var m *model.PackingOrder
	if id, e = ctx.Decrypt("id"); e == nil {
		if m, e = repository.GetPackingOrderDetailPack("id", id); e != nil {
			e = echo.ErrNotFound
		}
	}

	var file string
	if file, e = GetDetailedPack(s, m); e == nil {
		ctx.Files(file)
	}
	return ctx.Serve(e)
}

// ListPack : function to get list pack
func (h *Handler) ListPack(c echo.Context) (e error) {
	ctx := c.(*cuxs.Context)
	m := mongodb.NewMongo()
	o := orm.NewOrm()
	o.Using("read_only")

	var p []*model.PackingOrder

	// region filter checking
	var withFilter bool
	product, _ := common.Decrypt(ctx.QueryParam("product"))
	warehouse, _ := common.Decrypt(ctx.QueryParam("warehouse"))
	packingDateStr := ctx.QueryParam("packing_date")
	packingDateArr := strings.Split(packingDateStr, "|")

	// query param for searching
	queryString := ctx.QueryParam("queryString")
	// endregion

	cond := make(map[string]interface{})

	if product != 0 {
		withFilter = true
	}
	if warehouse != 0 {
		cond["po.warehouse_id = "] = warehouse
		withFilter = true
	}

	if packingDateStr != "" {
		cond["po.delivery_date "] = packingDateArr
		withFilter = true
	}
	// endregion

	var res []byte
	var response []*model.BarcodeModel
	var filter interface{}
	response = make([]*model.BarcodeModel, 0)

	if withFilter {
		var where string
		var values []interface{}

		if warehouse != 0 || packingDateStr != "" {
			for k, v := range cond {
				if reflect.TypeOf(v).Kind().String() == "slice" {
					where = where + " " + k + " >= ? and " + k + " <= ? and"
				} else {
					where = where + " " + k + "? and"
				}

				values = append(values, v)
			}

			where = strings.TrimSuffix(where, " and")
			q := "select po.id from packing_order po where " + where

			_, e = o.Raw(q, values).QueryRows(&p)

			var poID []int64

			for _, v := range p {
				poID = append(poID, v.ID)
			}

			if len(poID) == 0 {
				ctx.ResponseData = response
				m.DisconnectMongoClient()
				return ctx.Serve(e)
			}

			if product != 0 {
				filter = bson.D{
					{"packing_order_id", bson.M{"$in": poID}},
					{"product_id", product},
				}
			} else {
				filter = bson.D{
					{"packing_order_id", bson.M{"$in": poID}},
				}
			}

			if res, e = m.GetAllDataWithFilter("Packing_Barcode", filter); e != nil {
				m.DisconnectMongoClient()
				return ctx.Serve(e)
			}

			if len(res) == 0 {
				ctx.ResponseData = response
				m.DisconnectMongoClient()
				return ctx.Serve(e)
			}

			// region convert byte data to json data
			if e = json.Unmarshal(res, &response); e != nil {
				m.DisconnectMongoClient()
				return ctx.Serve(e)
			}
		} else if product != 0 {
			filter = bson.D{
				{"product_id", product},
			}

			if res, e = m.GetAllDataWithFilter("Packing_Barcode", filter); e != nil {
				m.DisconnectMongoClient()
				return ctx.Serve(e)
			}

			if len(res) == 0 {
				ctx.ResponseData = response
				m.DisconnectMongoClient()
				return ctx.Serve(e)
			}

			// region convert byte data to json data
			if e = json.Unmarshal(res, &response); e != nil {
				m.DisconnectMongoClient()
				return ctx.Serve(e)
			}
			// endregion
		}

	} else {
		if res, e = m.GetAllDataWithoutFilter("Packing_Barcode"); e != nil {
			return ctx.Serve(e)
		}

		if len(res) == 0 {
			ctx.ResponseData = response
			m.DisconnectMongoClient()
			return ctx.Serve(e)
		}
		// region convert byte data to json data
		if e = json.Unmarshal(res, &response); e != nil {
			m.DisconnectMongoClient()
			return ctx.Serve(e)
		}
	}

	// region read id obj
	for _, v := range response {
		v.Product, _ = repository.ValidProduct(v.ProductID)
		v.PackingOrder, _ = repository.ValidPackingOrder(v.PackingOrderID)
		v.Product.Uom.Read("ID")
		v.PackingOrder.Warehouse.Read("ID")

	}
	// endregion

	if len(response) == 0 {
		ctx.ResponseData = response
		m.DisconnectMongoClient()
		return ctx.Serve(e)
	}
	// region searching
	var pbTemp *model.BarcodeModel
	var pbs []*model.BarcodeModel
	if queryString != "" {
		for _, v := range response {
			if strings.Contains(v.Code, queryString) || strings.Contains(v.PackingOrder.Code, queryString) ||
				strings.Contains(v.Product.Code, queryString) || strings.Contains(v.Product.Name, queryString) {
				pbTemp = v
				pbs = append(pbs, pbTemp)
			}
		}
		response = pbs
	}
	// endregion
	ctx.ResponseData = response
	m.DisconnectMongoClient()
	return ctx.Serve(e)
}

// detailPack : function to get detailed data by id
func (h *Handler) detailPack(c echo.Context) (e error) {
	ctx := c.(*cuxs.Context)
	var r UpdatePackRequest

	if r.Session, e = auth.UserSession(ctx); e != nil {
		return ctx.Serve(e)
	}
	if r.ID, e = ctx.Decrypt("id"); e != nil {
		return ctx.Serve(e)
	}

	if e = ctx.Bind(&r); e != nil {
		return ctx.Serve(e)
	}
	if ctx.ResponseData, e = DetailPack(r); e != nil {
		return ctx.Serve(e)
	}

	return ctx.Serve(e)
}

// receivedPrint : function to print
func (h *Handler) receivedPrint(c echo.Context) (e error) {
	ctx := c.(*cuxs.Context)
	var r PrintPackRequest
	var po *model.ResponseData
	var pr ResponsePrint
	var code string
	req := make(map[string]interface{})

	po = new(model.ResponseData)

	if r.Session, e = auth.UserSession(ctx); e != nil {
		return ctx.Serve(e)
	}
	if r.ID, e = ctx.Decrypt("id"); e != nil {
		return ctx.Serve(e)
	}

	if e = ctx.Bind(&r); e != nil {
		return ctx.Serve(e)
	}

	po = r.ResponseData

	/*
		type print = 1(print insert), 0(re print)
	*/
	m := mongodb.NewMongo()
	if r.TypePrint == 1 {
		code, _ = util.GenerateDocCode("PK", po.Product.Code, "packing_code")
		// insert to mongo
		mongoData := model.BarcodeModel{
			Code:           code,
			PackingOrderID: po.PackingOrderID,
			ProductID:      po.ProductID,
			PackType:       po.PackType,
			WeightScale:    r.WeightScale,
			Status:         1,
			DeltaPrint:     1,
			CreatedAt:      time.Now().Format(("2006-01-02 15:04:05")),
			CreatedBy:      r.Session.Staff.ID,
		}

		m.CreateIndex("Packing_Barcode", "code", false)
		ID, err := m.InsertID("Packing_Barcode")
		if err != nil {
			m.DisconnectMongoClient()
			return ctx.Serve(e)
		}
		mongoData.ID = ID.(int64)

		_, err = m.InsertOneData("Packing_Barcode", mongoData)
		if err != nil {
			m.DisconnectMongoClient()
			return ctx.Serve(e)
		}
		po.CodePrint = code
	} else {
		var rd = new(model.BarcodeModel)
		filter := bson.D{
			{"packing_order_id", r.PackingOrder.ID},
			{"product_id", r.Product.ID},
			{"pack_type", r.PackType},
			{"status", 1},
		}

		opts := &options.FindOneOptions{}
		opts.SetSort(bson.D{{"code", -1}})
		var res []byte
		if res, e = m.GetOneDataWithFilter("Packing_Barcode", filter, opts); e != nil {
			return ctx.Serve(e)
		}

		// region convert byte data to json data
		if e = json.Unmarshal(res, &rd); e != nil {
			return ctx.Serve(e)
		}
		code = rd.Code
		po.CodePrint = code
		// endregion

		filterUpdate := bson.D{
			{"code", rd.Code},
		}
		updatePayload := bson.D{
			{"delta_print", rd.DeltaPrint + 1},
			{"created_at", time.Now().Format(("2006-01-02 15:04:05"))},
			{"created_by", r.Session.Staff.ID},
		}

		if e = m.UpdateOneDataWithFilter("Packing_Barcode", filterUpdate, updatePayload); e != nil {
			m.DisconnectMongoClient()
			return ctx.Serve(e)
		}
	}
	m.DisconnectMongoClient()

	req["pk"] = po
	file := util.SendPrint(req, "read/label_packing")

	pr.LinkPrint = file
	pr.Code = code
	pr.ExpectedTotalPack = r.ResponseData.ExpectedTotalPack
	pr.ActualTotalPack = r.ResponseData.ActualTotalPack
	ctx.ResponseData = pr

	return ctx.Serve(e)
}

// readFilter : function to get requested data based on parameters with filtered permission
func (h *Handler) readFilter(c echo.Context) (e error) {
	ctx := c.(*cuxs.Context)
	rq := ctx.RequestQuery()

	var data []*model.PackingOrder
	var total int64

	if data, total, e = repository.GetFilterPackingOrders(rq); e == nil {
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

// itemAssign : function to get detailed data by id
func (h *Handler) itemAssign(c echo.Context) (e error) {
	ctx := c.(*cuxs.Context)
	rq := ctx.RequestQuery()

	var data []*model.PackingOrderItemAssign
	var total int64

	if data, total, e = repository.GetPackingOrderItemAssign(rq); e == nil {
		ctx.Data(data, total)
	}

	return ctx.Serve(e)
}

// itemAssign : function to get detailed data by id
func (h *Handler) getSoi(c echo.Context) (e error) {
	ctx := c.(*cuxs.Context)
	rq := ctx.RequestQuery()

	deliveryDate := ctx.QueryParam("delivery_date")
	warehouse, _ := common.Decrypt(ctx.QueryParam("warehouse"))

	data, total, e := GetSalesOrderItem(rq, deliveryDate, warehouse)

	ctx.Data(data, total)

	return ctx.Serve(e)
}

// update : function to update packing order item assign
func (h *Handler) itemAssignUpdate(c echo.Context) (e error) {
	ctx := c.(*cuxs.Context)
	var r updateItemAssignRequest

	if r.Session, e = auth.UserSession(ctx); e == nil {

		if r.ID, e = ctx.Decrypt("id"); e == nil {
			if e = ctx.Bind(&r); e == nil {
				ctx.ResponseData, e = UpdateItemAssign(r)
			}
		}
	}

	return ctx.Serve(e)
}

// generatePacking : function to generate packing recommendation
func (h *Handler) generatePacking(c echo.Context) (e error) {
	ctx := c.(*cuxs.Context)
	var r generatePackingRequest

	if r.Session, e = auth.UserSession(ctx); e != nil {
		return ctx.Serve(e)
	}

	if e = ctx.Bind(&r); e != nil {
		return ctx.Serve(e)
	}
	if ctx.ResponseData, e = GeneratePacking(r); e != nil {
		return ctx.Serve(e)
	}

	return ctx.Serve(e)
}

// updatePackingPack : function to update packing pack
func (h *Handler) updatePackingPack(c echo.Context) (e error) {
	ctx := c.(*cuxs.Context)
	var r UpdatePackRequest

	if r.Session, e = auth.UserSession(ctx); e != nil {
		return ctx.Serve(e)
	}
	if r.ID, e = ctx.Decrypt("id"); e != nil {
		return ctx.Serve(e)
	}

	if e = ctx.Bind(&r); e != nil {
		return ctx.Serve(e)
	}
	if ctx.ResponseData, e = UpdatePackingPack(r); e != nil {
		return ctx.Serve(e)
	}

	return ctx.Serve(e)
}

// disposePackingPack : function to update packing pack
func (h *Handler) disposePackingPack(c echo.Context) (e error) {
	ctx := c.(*cuxs.Context)
	m := mongodb.NewMongo()
	var r DisposePackRequest

	if r.Session, e = auth.UserSession(ctx); e != nil {
		return ctx.Serve(e)
	}
	if r.ID, e = ctx.Decrypt("id"); e != nil {
		return ctx.Serve(e)
	}

	if e = ctx.Bind(&r); e != nil {
		return ctx.Serve(e)
	}
	if ctx.ResponseData, e = DisposePackingPack(r); e != nil {
		return ctx.Serve(e)
	}

	m.DisconnectMongoClient()
	return ctx.Serve(e)
}

func (h *Handler) template(c echo.Context) (e error) {
	ctx := c.(*cuxs.Context)
	var id int64
	if id, e = ctx.Decrypt("id"); e == nil {
		isExport := ctx.QueryParam("export") == "1"
		backdate := now.NewParse(time.RFC3339, ctx.QueryParam("date")).Time
		data, e := repository.GetPackingOrder("id", id)

		if e == nil {
			if isExport {
				var file string
				if file, e = DownloadActualPackingXls(backdate, data); e == nil {
					ctx.Files(file)
				}
			} else {
				ctx.Data(data, 0)
			}
		} else {
			return e
		}
	}

	return ctx.Serve(e)
}

// uploadActualPackingPacking : function to update packing order item assign
func (h *Handler) uploadActualPackingPacking(c echo.Context) (e error) {
	ctx := c.(*cuxs.Context)
	var r updateActualPackingRequest

	if r.Session, e = auth.UserSession(ctx); e == nil {

		if r.ID, e = ctx.Decrypt("id"); e == nil {
			if e = ctx.Bind(&r); e == nil {
				ctx.ResponseData, e = UpdateActualPacking(r)
			}
		}
	}

	return ctx.Serve(e)
}

// confirm : function to update packing order status into finished
func (h *Handler) confirm(c echo.Context) (e error) {
	ctx := c.(*cuxs.Context)
	var r confirmRequest

	if r.Session, e = auth.UserSession(ctx); e == nil {
		if r.ID, e = ctx.Decrypt("id"); e == nil {
			if e = ctx.Bind(&r); e == nil {
				ctx.ResponseData, e = Confirm(r)
			}
		}
	}

	return ctx.Serve(e)
}

// cancel : function to update packing order status into finished
func (h *Handler) cancel(c echo.Context) (e error) {
	ctx := c.(*cuxs.Context)
	var r cancelRequest

	if r.Session, e = auth.UserSession(ctx); e == nil {
		if r.ID, e = ctx.Decrypt("id"); e == nil {
			if e = ctx.Bind(&r); e == nil {
				ctx.ResponseData, e = Cancel(r)
			}
		}
	}

	return ctx.Serve(e)
}

// uploadActualPackingPacking : function to update packing order item assign
func (h *Handler) assignQuantityPacking(c echo.Context) (e error) {
	ctx := c.(*cuxs.Context)
	var r assignActualPackingRequest

	if r.Session, e = auth.UserSession(ctx); e == nil {

		if r.ID, e = ctx.Decrypt("id"); e == nil {
			if e = ctx.Bind(&r); e == nil {
				ctx.ResponseData, e = AssignActualPacking(r)
			}
		}
	}

	return ctx.Serve(e)
}
