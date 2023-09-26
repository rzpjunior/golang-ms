// Copyright 2020 PT. Eden Pangan Indonesia. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

package order

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"time"

	"git.edenfarm.id/cuxs/mongodb"
	"git.edenfarm.id/project-version2/api/datastore/repository"
	"git.edenfarm.id/project-version2/api/src/auth"
	"git.edenfarm.id/project-version2/api/util/kafka"
	"git.edenfarm.id/project-version2/datamodel/model"

	"strings"
	"sync"

	"net/http"

	"git.edenfarm.id/cuxs/common"
	"git.edenfarm.id/cuxs/cuxs"
	"git.edenfarm.id/cuxs/env"
	"git.edenfarm.id/project-version2/api/util"
	"github.com/labstack/echo/v4"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Handler struct{}

// URLMapping : function to map url with it's handler and add authorization validation
func (h *Handler) URLMapping(r *echo.Group) {
	r.GET("", h.read, auth.Authorized("do_rdl"))
	r.GET("/:id", h.detail, auth.Authorized("do_rdd"))
	r.GET("/filter", h.readFilter, auth.Authorized("filter_rdl"))
	r.GET("/item", h.readItem, auth.Authorized("do_rdl"))
	r.GET("/item/filter", h.readItemFilter, auth.Authorized("filter_rdl"))
	r.POST("", h.create, auth.Authorized("do_crt"))
	r.PUT("/:id", h.update, auth.Authorized("do_upd"))
	r.PUT("/confirm/:id", h.confirm, auth.Authorized("do_cnf"))
	r.PUT("/cancel/:id", h.cancel, auth.Authorized("do_can"))
	r.GET("/print/:id", h.receivePrint, auth.Authorized("do_prt"))

	//api for courier apps
	r.GET("/courier/:id", h.detailCourier, auth.Authorized("do_rdl"))

	// produce kafka + write mongodb
	r.POST("/produce_create", h.createKafkaMongo, auth.Authorized("do_crt"))
	r.PUT("/produce_update/:id", h.updateKafkaMongo, auth.Authorized("do_upd"))
	r.PUT("/produce_cancel/:id", h.cancelKafkaMongo, auth.Authorized("do_can"))

}

func (h *Handler) receivePrint(c echo.Context) (e error) {
	ctx := c.(*cuxs.Context)
	var do *model.DeliveryOrder
	var id int64
	var session *auth.SessionData
	configs := make(map[string]string)
	req := make(map[string]interface{})
	if session, e = auth.UserSession(ctx); e == nil {
		if id, e = ctx.Decrypt("id"); e == nil {
			if do, e = repository.GetDeliveryOrder("id", id); e != nil {
				e = echo.ErrNotFound
			} else {
				req["do"] = do
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

				file := util.SendPrint(req, "read/do")
				ctx.Files(file)

				// delta print
				do.DeltaPrint = do.DeltaPrint + 1
				do.Save("DeltaPrint")
			}
		}
	}

	return ctx.Serve(e)
}

// read : function to get requested data based on parameters
func (h *Handler) read(c echo.Context) (e error) {
	ctx := c.(*cuxs.Context)
	rq := ctx.RequestQuery()

	var data []*model.DeliveryOrder
	var total int64

	if data, total, e = repository.GetDeliveryOrders(rq); e == nil {
		ctx.Data(data, total)
	}

	return ctx.Serve(e)
}

// detail : function to get detailed data by id
func (h *Handler) detail(c echo.Context) (e error) {
	ctx := c.(*cuxs.Context)

	var id int64
	if id, e = ctx.Decrypt("id"); e == nil {
		if ctx.ResponseData, e = repository.GetDeliveryOrder("id", id); e != nil {
			e = echo.ErrNotFound
		}
	}

	return ctx.Serve(e)
}

// readFilter : function to get requested data based on parameters with filtered permission
func (h *Handler) readFilter(c echo.Context) (e error) {
	ctx := c.(*cuxs.Context)
	rq := ctx.RequestQuery()

	var data []*model.DeliveryOrder
	var total int64

	if data, total, e = repository.GetFilterDeliveryOrders(rq); e == nil {
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
		if !r.UseRedis {
			return ctx.Serve(e)
		}
		if ok, err := r.Mutex.Unlock(); !ok || err != nil {
			if err != nil {
				err = echo.NewHTTPError(http.StatusLocked, err.Error())
			}
			if e != nil && err == nil {
				err = e
			}
			r.Client.Close()
			return ctx.Serve(err)
		}
		r.Client.Close()
		return ctx.Serve(e)
	}
	ctx.ResponseData, e = Save(r)
	if e != nil {
		if !r.UseRedis {
			return ctx.Serve(e)
		}
		if ok, err := r.Mutex.Unlock(); !ok || err != nil {
			if err != nil {
				err = echo.NewHTTPError(http.StatusLocked, err.Error())
			}
			r.Client.Close()
			return ctx.Serve(err)
		}
		r.Client.Close()
		return ctx.Serve(e)
	}
	if !r.UseRedis {
		return ctx.Serve(e)
	}
	if ok, err := r.Mutex.Unlock(); !ok || err != nil {
		if err != nil {
			err = echo.NewHTTPError(http.StatusLocked, err.Error())
		}
		r.Client.Close()
		return ctx.Serve(err)
	}
	r.Client.Close()
	return ctx.Serve(e)
}

// confirm : function to confirm delivery order
func (h *Handler) confirm(c echo.Context) (e error) {
	ctx := c.(*cuxs.Context)
	var r confirmRequest
	r.DeliveryOrder = new(model.DeliveryOrder)
	if r.ID, e = common.Decrypt(ctx.Param("id")); e == nil {
		if r.DeliveryOrder, e = repository.GetDeliveryOrder("id", r.ID); e == nil {
			if r.Session, e = auth.UserSession(ctx); e == nil {
				if e = ctx.Bind(&r); e == nil {
					ctx.ResponseData, e = r.Confirm()
				}
			}
		}
	}
	return ctx.Serve(e)
}

// update : function to update delivery order
func (h *Handler) update(c echo.Context) (e error) {
	ctx := c.(*cuxs.Context)
	var r updateRequest
	r.DeliveryOrder = new(model.DeliveryOrder)
	if r.ID, e = common.Decrypt(ctx.Param("id")); e != nil {
		return ctx.Serve(e)
	}
	if r.DeliveryOrder, e = repository.GetDeliveryOrder("id", r.ID); e != nil {
		return ctx.Serve(e)
	}
	if r.Session, e = auth.UserSession(ctx); e != nil {
		return ctx.Serve(e)
	}
	if e = ctx.Bind(&r); e != nil {
		if !r.UseRedis {
			return ctx.Serve(e)
		}
		if ok, err := r.Mutex.Unlock(); !ok || err != nil {
			if err != nil {
				err = echo.NewHTTPError(http.StatusLocked, err.Error())
			}
			if e != nil && err == nil {
				err = e
			}
			r.Client.Close()
			return ctx.Serve(err)
		}
		r.Client.Close()
		return ctx.Serve(e)

	}

	ctx.ResponseData, e = r.Update()
	if e != nil {
		if !r.UseRedis {
			return ctx.Serve(e)
		}
		if ok, err := r.Mutex.Unlock(); !ok || err != nil {
			if err != nil {
				err = echo.NewHTTPError(http.StatusLocked, err.Error())
			}
			r.Client.Close()
			return ctx.Serve(err)
		}
		r.Client.Close()
		return ctx.Serve(e)
	}
	if !r.UseRedis {
		return ctx.Serve(e)
	}
	if ok, err := r.Mutex.Unlock(); !ok || err != nil {
		if err != nil {
			err = echo.NewHTTPError(http.StatusLocked, err.Error())
		}
		r.Client.Close()
		return ctx.Serve(err)
	}
	r.Client.Close()
	return ctx.Serve(e)
}

// cancel : function to cancel delivery order
func (h *Handler) cancel(c echo.Context) (e error) {
	ctx := c.(*cuxs.Context)
	var r cancelRequest
	r.DeliveryOrder = new(model.DeliveryOrder)
	if r.DeliveryOrder.ID, e = common.Decrypt(ctx.Param("id")); e != nil {
		return ctx.Serve(e)
	}
	if r.DeliveryOrder, e = repository.GetDeliveryOrder("id", r.DeliveryOrder.ID); e != nil {
		return ctx.Serve(e)
	}
	if r.Session, e = auth.UserSession(ctx); e != nil {
		return ctx.Serve(e)
	}
	if e = ctx.Bind(&r); e != nil {
		if !r.UseRedis {
			return ctx.Serve(e)
		}
		if ok, err := r.Mutex.Unlock(); !ok || err != nil {
			if err != nil {
				err = echo.NewHTTPError(http.StatusLocked, err.Error())
			}
			if e != nil && err == nil {
				err = e
			}
			r.Client.Close()
			return ctx.Serve(err)
		}
		r.Client.Close()
		return ctx.Serve(e)
	}
	ctx.ResponseData, e = r.Cancel()
	if e != nil {
		if !r.UseRedis {
			return ctx.Serve(e)
		}
		if ok, err := r.Mutex.Unlock(); !ok || err != nil {
			if err != nil {
				err = echo.NewHTTPError(http.StatusLocked, err.Error())
			}
			r.Client.Close()
			return ctx.Serve(err)
		}
		r.Client.Close()
		return ctx.Serve(e)
	}
	if !r.UseRedis {
		return ctx.Serve(e)
	}
	if ok, err := r.Mutex.Unlock(); !ok || err != nil {
		if err != nil {
			err = echo.NewHTTPError(http.StatusLocked, err.Error())
		}
		r.Client.Close()
		return ctx.Serve(err)
	}
	r.Client.Close()
	return ctx.Serve(e)
}

// readItem : function to get requested item data based on parameters
func (h *Handler) readItem(c echo.Context) (e error) {
	ctx := c.(*cuxs.Context)
	rq := ctx.RequestQuery()

	var data []*model.DeliveryOrderItem
	var total int64

	if data, total, e = repository.GetDeliveryOrderItems(rq); e == nil {
		ctx.Data(data, total)
	}

	return ctx.Serve(e)
}

// readItem : function to get requested item data based on parameters with filtered permission
func (h *Handler) readItemFilter(c echo.Context) (e error) {
	ctx := c.(*cuxs.Context)
	rq := ctx.RequestQuery()

	var data []*model.DeliveryOrderItem
	var total int64

	if data, total, e = repository.GetFilterDeliveryOrderItems(rq); e == nil {
		ctx.Data(data, total)
	}

	return ctx.Serve(e)
}

// detailCourier : function to get detailed data by id
func (h *Handler) detailCourier(c echo.Context) (e error) {
	ctx := c.(*cuxs.Context)

	var id int64
	if id, e = ctx.Decrypt("id"); e == nil {
		if ctx.ResponseData, e = repository.GetCourierByDeliveryOrderID(id); e != nil {
			e = echo.ErrNotFound
		}
	}

	return ctx.Serve(e)
}

// createKafkaMongo : function to create new data based on input (produce kafka and write into mongo)
func (h *Handler) createKafkaMongo(c echo.Context) (e error) {
	ctx := c.(*cuxs.Context)
	var (
		r      createRequest
		tempDO *model.DeliveryOrder
		d      sync.Mutex
		wg     sync.WaitGroup
	)

	if r.Session, e = auth.UserSession(ctx); e != nil {
		return ctx.Serve(e)
	}

	if e = ctx.Bind(&r); e != nil {
		if !r.UseRedis {
			return ctx.Serve(e)
		}
		if ok, err := r.Mutex.Unlock(); !ok || err != nil {
			if err != nil {
				err = echo.NewHTTPError(http.StatusLocked, err.Error())
			}
			if e != nil && err == nil {
				err = e
			}
			r.Client.Close()
			return ctx.Serve(err)
		}
		r.Client.Close()
		return ctx.Serve(e)
	}

	// Check from config app whether set use kafka + mongo or not
	checkStatus, e := repository.GetConfigApp("attribute", "kafka_mongo_delivery_order")
	if e != nil {
		if !r.UseRedis {
			return ctx.Serve(e)
		}
		if ok, err := r.Mutex.Unlock(); !ok || err != nil {
			if err != nil {
				err = echo.NewHTTPError(http.StatusLocked, err.Error())
			}
			r.Client.Close()
			return ctx.Serve(err)
		}
		r.Client.Close()
		return ctx.Serve(e)
	}

	if checkStatus.Value == "1" {
		ctx1 := context.Background()

		// lock unlock
		wg.Add(1)
		go func() {
			d.Lock()
			tempDO, e = Save(r)

			r.TypeRequest = "create"
			jobs := &model.Jobs{
				EndpointUrl:    "/v1/stock/consume_create_stock",
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
			_, e := m.InsertOneData("Jobs", jobs)
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
			ctx.ResponseData = tempDO
			d.Unlock()
			defer wg.Done()
		}()

		wg.Wait()

		if e != nil {
			if !r.UseRedis {
				return ctx.Serve(e)
			}
			if ok, err := r.Mutex.Unlock(); !ok || err != nil {
				if err != nil {
					err = echo.NewHTTPError(http.StatusLocked, err.Error())
				}
				r.Client.Close()
				return ctx.Serve(err)
			}
			r.Client.Close()
			return ctx.Serve(e)
		}

	} else {
		ctx.ResponseData, e = Save(r)
		if e != nil {
			if !r.UseRedis {
				return ctx.Serve(e)
			}
			if ok, err := r.Mutex.Unlock(); !ok || err != nil {
				if err != nil {
					err = echo.NewHTTPError(http.StatusLocked, err.Error())
				}
				r.Client.Close()
				return ctx.Serve(err)
			}
			r.Client.Close()
			return ctx.Serve(e)
		}
	}
	if !r.UseRedis {
		return ctx.Serve(e)
	}
	if ok, err := r.Mutex.Unlock(); !ok || err != nil {
		if err != nil {
			err = echo.NewHTTPError(http.StatusLocked, err.Error())
		}
		r.Client.Close()
		return ctx.Serve(err)
	}
	r.Client.Close()
	return ctx.Serve(e)

}

// updateKafkaMongo : function to update data based on input (produce kafka and write into mongo)
func (h *Handler) updateKafkaMongo(c echo.Context) (e error) {
	ctx := c.(*cuxs.Context)
	var (
		r  updateRequest
		d  sync.Mutex
		wg sync.WaitGroup
	)

	ID := ctx.Param("id")
	r.DeliveryOrder = new(model.DeliveryOrder)
	if r.ID, e = common.Decrypt(ID); e != nil {
		return ctx.Serve(e)
	}
	if r.DeliveryOrder, e = repository.GetDeliveryOrder("id", r.ID); e != nil {
		return ctx.Serve(e)
	}
	if r.Session, e = auth.UserSession(ctx); e != nil {
		return ctx.Serve(e)
	}

	if e = ctx.Bind(&r); e != nil {
		if !r.UseRedis {
			return ctx.Serve(e)
		}
		if ok, err := r.Mutex.Unlock(); !ok || err != nil {
			if err != nil {
				err = echo.NewHTTPError(http.StatusLocked, err.Error())
			}
			if e != nil && err == nil {
				err = e
			}
			r.Client.Close()
			return ctx.Serve(err)
		}
		r.Client.Close()
		return ctx.Serve(e)
	}

	// Check from config app whether set use kafka + mongo or not
	checkStatus, e := repository.GetConfigApp("attribute", "kafka_mongo_delivery_order")
	if e != nil {
		if !r.UseRedis {
			return ctx.Serve(e)
		}
		if ok, err := r.Mutex.Unlock(); !ok || err != nil {
			if err != nil {
				err = echo.NewHTTPError(http.StatusLocked, err.Error())
			}
			r.Client.Close()
			return ctx.Serve(err)
		}
		r.Client.Close()
		return ctx.Serve(e)
	}

	if checkStatus.Value == "1" {
		r.EncryptID = ID

		ctx1 := context.Background()

		// update document
		// lock unlock
		wg.Add(1)
		go func() {
			d.Lock()
			r.DeliveryOrder, e = r.Update()

			jobs := &model.Jobs{
				EndpointUrl:    "/v1/stock/consume_update_stock/" + ID,
				Topic:          env.GetString("KAFKA_TOPIC", ""),
				EndpointMethod: "PUT",
				ResponseBody:   "[]",
				Status:         1,
				CreatedAt:      time.Now(),
				CreatedBy:      r.Session.Staff.User.ID,
				RetryCount:     0,
			}

			m := mongodb.NewMongo()
			m.CreateIndex("Jobs", "_id", true)
			jobs.ID = primitive.NewObjectID()

			_, err := m.InsertOneData("Jobs", jobs)
			if err != nil {
				e = err
				fmt.Println(e)
			}

			jobsFilter := *jobs
			r.JobsID = jobs.ID
			b, _ := json.Marshal(r)
			jobs.RequestBody = string(b)
			e = kafka.Produce(ctx1, jobs, jobs.Topic)
			if e != nil {
				jobs.ResponseBody = "{\"error_produce\":\"" + e.Error() + "\"}"
				jobs.Status = 5
				err := m.UpdateOneDataWithFilter("Jobs", jobsFilter, jobs)
				if err != nil {
					e = err
					fmt.Println(err)
					m.DisconnectMongoClient()
				}
				fmt.Println(e)
			}
			ctx.ResponseData = jobs
			d.Unlock()
			defer wg.Done()
		}()

		wg.Wait()

		if e != nil {
			if !r.UseRedis {
				return ctx.Serve(e)
			}
			if ok, err := r.Mutex.Unlock(); !ok || err != nil {
				if err != nil {
					err = echo.NewHTTPError(http.StatusLocked, err.Error())
				}
				r.Client.Close()
				return ctx.Serve(err)
			}
			r.Client.Close()
			return ctx.Serve(e)
		}

	} else {
		ctx.ResponseData, e = r.Update()
		if e != nil {
			if !r.UseRedis {
				return ctx.Serve(e)
			}
			if ok, err := r.Mutex.Unlock(); !ok || err != nil {
				if err != nil {
					err = echo.NewHTTPError(http.StatusLocked, err.Error())
				}
				r.Client.Close()
				return ctx.Serve(err)
			}
			r.Client.Close()
			return ctx.Serve(e)
		}

	}
	if !r.UseRedis {
		return ctx.Serve(e)
	}
	if ok, err := r.Mutex.Unlock(); !ok || err != nil {
		if err != nil {
			err = echo.NewHTTPError(http.StatusLocked, err.Error())
		}
		r.Client.Close()
		return ctx.Serve(err)
	}
	r.Client.Close()
	return ctx.Serve(e)
}

// cancelKafkaMongo : function to cancel data based on selected id (produce kafka and write into mongo)
func (h *Handler) cancelKafkaMongo(c echo.Context) (e error) {
	ctx := c.(*cuxs.Context)
	var (
		r      cancelRequest
		tempDO *model.DeliveryOrder
	)

	if r.Session, e = auth.UserSession(ctx); e != nil {
		return ctx.Serve(e)
	}

	ID := ctx.Param("id")
	r.DeliveryOrder = new(model.DeliveryOrder)
	if r.DeliveryOrder.ID, e = common.Decrypt(ctx.Param("id")); e != nil {
		return ctx.Serve(e)
	}
	if r.DeliveryOrder, e = repository.GetDeliveryOrder("id", r.DeliveryOrder.ID); e != nil {
		return ctx.Serve(e)
	}

	if e = ctx.Bind(&r); e != nil {
		if !r.UseRedis {
			return ctx.Serve(e)
		}
		if ok, err := r.Mutex.Unlock(); !ok || err != nil {
			if err != nil {
				err = echo.NewHTTPError(http.StatusLocked, err.Error())
			}
			if e != nil && err == nil {
				err = e
			}
			r.Client.Close()
			return ctx.Serve(err)
		}
		r.Client.Close()
		return ctx.Serve(e)
	}

	// Check from config app whether set use kafka + mongo or not
	checkStatus, e := repository.GetConfigApp("attribute", "kafka_mongo_delivery_order")
	if e != nil {
		if !r.UseRedis {
			return ctx.Serve(e)
		}
		if ok, err := r.Mutex.Unlock(); !ok || err != nil {
			if err != nil {
				err = echo.NewHTTPError(http.StatusLocked, err.Error())
			}
			r.Client.Close()
			return ctx.Serve(err)
		}
		r.Client.Close()
		return ctx.Serve(e)
	}

	if checkStatus.Value == "1" {
		ctx1 := context.Background()
		m := mongodb.NewMongo()

		jobs := &model.Jobs{
			EndpointUrl:    "/v1/stock/consume_cancel_stock/" + ID,
			Topic:          env.GetString("KAFKA_TOPIC", ""),
			EndpointMethod: "PUT",
			ResponseBody:   "[]",
			Status:         1,
			CreatedAt:      time.Now(),
			CreatedBy:      r.Session.Staff.User.ID,
			RetryCount:     0,
		}

		tempDO, e = r.Cancel()
		if e != nil {
			if !r.UseRedis {
				return ctx.Serve(e)
			}
			if ok, err := r.Mutex.Unlock(); !ok || err != nil {
				if err != nil {
					err = echo.NewHTTPError(http.StatusLocked, err.Error())
				}
				r.Client.Close()
				return ctx.Serve(err)
			}
			r.Client.Close()
			return ctx.Serve(e)

		}
		ctx.ResponseData = tempDO
		r.TypeRequest = "cancel"

		jobsFilterDocument := &model.Jobs{EndpointUrl: jobs.EndpointUrl}

		if ret, err := m.GetCountDataWithFilter("Jobs", jobsFilterDocument); ret >= 1 || err != nil {
			e = errors.New("File Exist in Mongo")
			if !r.UseRedis {
				return ctx.Serve(e)
			}
			if ok, err := r.Mutex.Unlock(); !ok || err != nil {
				if err != nil {
					err = echo.NewHTTPError(http.StatusLocked, err.Error())
				}
				r.Client.Close()
				return ctx.Serve(err)
			}
			r.Client.Close()
			return ctx.Serve(e)
		}

		m.CreateIndex("Jobs", "_id", true)
		jobs.ID = primitive.NewObjectID()
		r.JobsID = jobs.ID
		a, _ := json.Marshal(r)
		jobs.RequestBody = string(a)

		_, err := m.InsertOneData("Jobs", jobs)
		if err != nil {
			e = err
			fmt.Println(e)
		}

		jobsFilter := *jobs

		e = kafka.Produce(ctx1, jobs, jobs.Topic)
		if e != nil {
			jobs.ResponseBody = "{\"error_produce\":\"" + e.Error() + "\"}"
			jobs.Status = 5
			err := m.UpdateOneDataWithFilter("Jobs", jobsFilter, jobs)
			if err != nil {
				e = err
				fmt.Println(err)
				m.DisconnectMongoClient()
			}
			fmt.Println(e)
		}
		ctx.ResponseData = jobs

		if e != nil {
			if !r.UseRedis {
				return ctx.Serve(e)
			}
			if ok, err := r.Mutex.Unlock(); !ok || err != nil {
				if err != nil {
					err = echo.NewHTTPError(http.StatusLocked, err.Error())
				}
				r.Client.Close()
				return ctx.Serve(err)
			}
			r.Client.Close()
			return ctx.Serve(e)
		}

	} else {
		ctx.ResponseData, e = r.Cancel()
		if e != nil {
			if !r.UseRedis {
				return ctx.Serve(e)
			}
			if ok, err := r.Mutex.Unlock(); !ok || err != nil {
				if err != nil {
					err = echo.NewHTTPError(http.StatusLocked, err.Error())
				}
				r.Client.Close()
				return ctx.Serve(err)
			}
			r.Client.Close()
			return ctx.Serve(e)
		}
	}
	if !r.UseRedis {
		return ctx.Serve(e)
	}
	if ok, err := r.Mutex.Unlock(); !ok || err != nil {
		if err != nil {
			err = echo.NewHTTPError(http.StatusLocked, err.Error())
		}
		r.Client.Close()
		return ctx.Serve(err)
	}
	r.Client.Close()

	return ctx.Serve(e)
}
