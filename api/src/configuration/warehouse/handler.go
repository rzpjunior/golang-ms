// Copyright 2020 PT. Eden Pangan Indonesia. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

package warehouse

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"git.edenfarm.id/cuxs/common"
	"git.edenfarm.id/cuxs/cuxs"
	"git.edenfarm.id/cuxs/dbredis"
	"git.edenfarm.id/cuxs/env"
	"git.edenfarm.id/cuxs/mongodb"
	"git.edenfarm.id/project-version2/api/datastore/repository"
	"git.edenfarm.id/project-version2/api/src/auth"
	"git.edenfarm.id/project-version2/api/util/kafka"
	"git.edenfarm.id/project-version2/datamodel/model"
	"github.com/labstack/echo/v4"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

// Handler collection handler for auth.
type Handler struct{}

// URLMapping declare endpoint with handler function.
func (h *Handler) URLMapping(r *echo.Group) {
	r.GET("", h.get, auth.Authorized("wrh_rdl"))
	r.GET("/:id", h.detail, auth.Authorized("wrh_rdd"))
	r.GET("/filter", h.readFilter, auth.Authorized("filter_rdl"))
	r.POST("", h.create, auth.Authorized("wrh_crt"))
	r.PUT("/:id", h.update, auth.Authorized("wrh_upd"))
	r.PUT("/archive/:id", h.archive, auth.Authorized("wrh_arc"))
	r.PUT("/unarchive/:id", h.unarchive, auth.Authorized("wrh_urc"))
	r.PUT("/update/param/:id", h.updateParam, auth.Authorized("wrh_upd_pl_par"))

	r.POST("/produce_create", h.createKafkaMongo, auth.Authorized("wrh_crt"))
	r.GET("/self-pickup", h.getDefaultWarehouseSelfPickUp, auth.Authorized("filter_rdl"))
}

func (h *Handler) get(c echo.Context) (e error) {
	ctx := c.(*cuxs.Context)

	data, total, e := repository.GetWarehouses(ctx.RequestQuery())
	if e == nil {
		ctx.Data(data, total)
	}

	return ctx.Serve(e)
}

func (h *Handler) readFilter(c echo.Context) (e error) {
	ctx := c.(*cuxs.Context)

	areaId := ctx.QueryParam("area")
	warehouseId := ctx.QueryParam("warehouse")

	data, total, e := repository.GetFilterWarehouse(ctx.RequestQuery(), areaId, warehouseId)
	if e == nil {
		ctx.Data(data, total)
	}

	return ctx.Serve(e)
}

// create : function to create new data based on input
func (h *Handler) create(c echo.Context) (e error) {
	ctx := c.(*cuxs.Context)
	var r createRequest

	if r.Session, e = auth.UserSession(ctx); e != nil {
		return ctx.Serve(e)
	}
	if e = ctx.Bind(&r); e != nil {
		return ctx.Serve(e)
	}

	ctx.ResponseData, e = Save(r)

	dbredis.Redis.DeleteCache("warehouse_creation")

	return ctx.Serve(e)
}

// archive : function to set status of active data into archive
func (h *Handler) archive(c echo.Context) (e error) {
	ctx := c.(*cuxs.Context)
	var r archiveRequest

	if r.Session, e = auth.UserSession(ctx); e == nil {
		if r.ID, e = ctx.Decrypt("id"); e == nil {
			if e = ctx.Bind(&r); e == nil {
				ctx.ResponseData, e = Archive(r)
			}
		}
	}

	return ctx.Serve(e)
}

// unarchive : function to set status of archive data into active
func (h *Handler) unarchive(c echo.Context) (e error) {
	ctx := c.(*cuxs.Context)
	var r unarchiveRequest

	if r.Session, e = auth.UserSession(ctx); e == nil {
		if r.ID, e = ctx.Decrypt("id"); e == nil {
			if e = ctx.Bind(&r); e == nil {
				ctx.ResponseData, e = Unarchive(r)
			}
		}
	}

	return ctx.Serve(e)
}

// detail : function to get detailed data by id
func (h *Handler) detail(c echo.Context) (e error) {
	ctx := c.(*cuxs.Context)
	var id int64

	if id, e = ctx.Decrypt("id"); e == nil {
		if ctx.ResponseData, e = repository.GetWarehouse("id", id); e != nil {
			e = echo.ErrNotFound
		}
	}

	return ctx.Serve(e)
}

//create : function to update data based on input
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

//updateParam : function to update data based on input
func (h *Handler) updateParam(c echo.Context) (e error) {
	ctx := c.(*cuxs.Context)
	var r updateParamRequest

	if r.Session, e = auth.UserSession(ctx); e == nil {
		if r.ID, e = ctx.Decrypt("id"); e == nil {
			if e = ctx.Bind(&r); e == nil {
				ctx.ResponseData, e = UpdateParam(r)
			}
		}
	}

	return ctx.Serve(e)
}

// createKafkaMongo : function to create new data based on input (produce kafka and write into mongo)
func (h *Handler) createKafkaMongo(c echo.Context) (e error) {
	ctx := c.(*cuxs.Context)
	ctx1 := context.Background()
	var r createRequest

	if r.Session, e = auth.UserSession(ctx); e != nil {
		return ctx.Serve(e)
	}

	if e = ctx.Bind(&r); e != nil {
		return ctx.Serve(e)
	}

	r.TypeRequest = "create"
	jobs := &model.Jobs{
		EndpointUrl:    "/v1/config/warehouse",
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
	r.JobsID = jobs.ID.Hex()
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
	m.DisconnectMongoClient()

	return ctx.Serve(e)
}

// getDefaultWarehouseSelfPickUp : function to get default warehouse self pick up by area id
func (h *Handler) getDefaultWarehouseSelfPickUp(c echo.Context) (e error) {
	ctx := c.(*cuxs.Context)
	var areaID int64

	if ctx.QueryParam("area") != "" {
		areaID, _ = common.Decrypt(ctx.QueryParam("area"))
	}

	if ctx.ResponseData, e = repository.GetWarehouseSelfPickupByAreaID(areaID); e != nil {
		return ctx.Serve(e)
	}

	return ctx.Serve(e)
}
