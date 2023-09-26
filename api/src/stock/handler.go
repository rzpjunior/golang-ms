// Copyright 2020 PT. Eden Pangan Indonesia. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

package stock

import (
	"encoding/json"
	"fmt"
	"net/http"

	"git.edenfarm.id/project-version2/datamodel/model"
	"git.edenfarm.id/project-version2/api/datastore/repository"
	"git.edenfarm.id/project-version2/api/src/auth"

	"git.edenfarm.id/cuxs/common"
	"git.edenfarm.id/cuxs/cuxs"
	"git.edenfarm.id/cuxs/mongodb"
	"github.com/labstack/echo/v4"
)

type Handler struct{}

// URLMapping : function to map url with it's handler and add authorization validation
func (h *Handler) URLMapping(r *echo.Group) {
	r.POST("/consume_create_stock", h.stockCreateWithKafka, auth.Authorized("do_crt"))
	r.PUT("/consume_update_stock/:id", h.stockUpdateWithKafka, auth.Authorized("do_upd"))
	r.PUT("/consume_cancel_stock/:id", h.stockCancelWithKafka, auth.Authorized("do_crt"))
}

// stockCreateWithKafka : function to create stock using kafka and mongodb based on input
func (h *Handler) stockCreateWithKafka(c echo.Context) (e error) {

	ctx := c.(*cuxs.Context)
	var r createRequest
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
		ctx.ResponseData, e = SaveStockDO(r)
		if e != nil {
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

// stockCancelWithKafka : function to cancel stock using kafka and mongodb based on input
func (h *Handler) stockCancelWithKafka(c echo.Context) (e error) {

	ctx := c.(*cuxs.Context)
	var r cancelRequest

	r.DeliveryOrder = new(model.DeliveryOrder)
	if r.DeliveryOrder.ID, e = common.Decrypt(ctx.Param("id")); e != nil {
		return e
	}
	if r.DeliveryOrder, e = repository.GetDeliveryOrder("id", r.DeliveryOrder.ID); e != nil {
		return e
	}

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

	if jobsChanged.Status == 2 {
		if ctx.ResponseData, e = CancelDOStock(r); e == nil {
			jobsChanged.Status = 3
			e = m.UpdateOneDataWithFilter("Jobs", jobs, jobsChanged)
			if e != nil {
				fmt.Println(e)
				m.DisconnectMongoClient()
			}
		}
	} else if jobsChanged.Status == 3 {
		ctx.ResponseData = http.StatusOK
	}
	m.DisconnectMongoClient()

	return ctx.Serve(e)
}

// stockUpdateWithKafka : function to update stock using kafka and mongodb based on input
func (h *Handler) stockUpdateWithKafka(c echo.Context) (e error) {
	ctx := c.(*cuxs.Context)
	var r updateRequest

	r.DeliveryOrder = new(model.DeliveryOrder)
	if r.DeliveryOrder.ID, e = common.Decrypt(ctx.Param("id")); e != nil {
		return e
	}

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

	if jobsChanged.Status == 2 {
		if ctx.ResponseData, e = UpdateStockDO(r); e == nil {
			jobsChanged.Status = 3
			e = m.UpdateOneDataWithFilter("Jobs", jobs, jobsChanged)
			if e != nil {
				fmt.Println(e)
				m.DisconnectMongoClient()
			}
		}
	} else if jobsChanged.Status == 3 {
		ctx.ResponseData = http.StatusOK
	}
	m.DisconnectMongoClient()

	return ctx.Serve(e)
}
