// Copyright 2020 PT. Eden Pangan Indonesia. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

package order

import (
	"net/http"
	"strconv"
	"strings"
	"time"

	"git.edenfarm.id/cuxs/cuxs"
	"git.edenfarm.id/cuxs/env"
	"git.edenfarm.id/project-version2/api/datastore/repository"
	"git.edenfarm.id/project-version2/api/src/auth"
	"git.edenfarm.id/project-version2/api/util"
	"git.edenfarm.id/project-version2/datamodel/model"
	goredislib "github.com/go-redis/redis/v8"
	"github.com/go-redsync/redsync/v4"
	"github.com/go-redsync/redsync/v4/redis/goredis/v8"
	"github.com/labstack/echo/v4"
)

type Handler struct{}

// URLMapping : function to map url with it's handler and add authorization validation
func (h *Handler) URLMapping(r *echo.Group) {
	r.GET("", h.read, auth.Authorized("so_rdl"))
	r.GET("/:id", h.detail, auth.Authorized("so_rdd"))
	r.GET("/filter", h.readFilter, auth.Authorized("filter_rdl"))
	r.POST("", h.create, auth.Authorized("so_crt"))
	r.PUT("/:id", h.update, auth.Authorized("so_upd"))
	r.PUT("/cancel/:id", h.cancel, auth.Authorized("so_can"))
	r.GET("/print/:id", h.receivePrint, auth.Authorized("so_prt"))
	r.PUT("/lock/:id", h.lock, auth.Authorized("so_upd"))
	r.PUT("/price/:id", h.updatePriceDeliveredSO, auth.Authorized("so_upd_prc"))
}

func (h *Handler) receivePrint(c echo.Context) (e error) {
	//var r requestGet

	ctx := c.(*cuxs.Context)
	var so *model.SalesOrder
	var id int64
	configs := make(map[string]string)
	req := make(map[string]interface{})

	if id, e = ctx.Decrypt("id"); e == nil {
		if so, e = repository.GetSalesOrder("id", id); e != nil {
			e = echo.ErrNotFound
		} else {
			req["so"] = so
			if config, _, e := repository.GetConfigAppsByAttribute("attribute__icontains", "company"); e == nil {
				for _, v := range config {
					configs[strings.TrimPrefix(v.Attribute, "company_")] = v.Value
				}
				configs["address"] = strings.ReplaceAll(configs["address"], "<br>", "\n")
				req["config"] = configs
			} else {
				e = echo.ErrNotFound
			}

			file := util.SendPrint(req, "read/so")
			ctx.Files(file)
		}
	}

	//return ctx.Serve(e)
	return ctx.Serve(e)
}

// read : function to get requested data based on parameters
func (h *Handler) read(c echo.Context) (e error) {
	ctx := c.(*cuxs.Context)
	rq := ctx.RequestQuery()

	var data []*model.SalesOrder
	var total int64

	if data, total, e = repository.GetSalesOrders(rq); e == nil {
		ctx.Data(data, total)
	}

	return ctx.Serve(e)
}

// readFilter : function to get requested data based on parameters with filtered permission
func (h *Handler) readFilter(c echo.Context) (e error) {
	ctx := c.(*cuxs.Context)
	rq := ctx.RequestQuery()

	var data []*model.SalesOrder
	var total int64

	if data, total, e = repository.GetFilterSalesOrders(rq); e == nil {
		ctx.Data(data, total)
	}

	return ctx.Serve(e)
}

// create : function to create new data based on input
func (h *Handler) create(c echo.Context) (e error) {
	ctx := c.(*cuxs.Context)
	var (
		r         createRequest
		qr        quotaRequest
		mutexname string
		mutexs    []*redsync.Mutex
		rs        *redsync.Redsync
		mutex     *redsync.Mutex
	)

	if r.Session, e = auth.UserSession(ctx); e != nil {
		return ctx.Serve(e)
	}

	//redsync-lock
	client := goredislib.NewClient(&goredislib.Options{
		Addr:     env.GetString("REDIS_HOST", "127.0.0.1:6379"),
		Password: env.GetString("REDIS_PASSWORD", "127.0.0.1:6379"),

		MaxRetries: 10,
	})

	defer client.Close()

	pool := goredis.NewPool(client)

	// Create an instance of redisync to be used to obtain a mutual exclusion
	// lock.
	rs = redsync.New(pool)
	mutexname = "sales_order_" + r.MerchantID
	mutex = rs.NewMutex(mutexname, redsync.WithRetryDelay(100*time.Millisecond), redsync.WithExpiry(60*time.Second))
	if e = mutex.Lock(); e != nil {
		e = echo.NewHTTPError(http.StatusLocked, e.Error())
		return ctx.Serve(e)
	}

	if e = ctx.Bind(&r); e != nil {
		if _, err := mutex.Unlock(); err != nil {
			return ctx.Serve(err)
		}
		return ctx.Serve(e)
	}

	mutexs = append(mutexs, mutex)

	for _, v := range r.SkuDiscountItems {

		// Obtain a new mutex by using the same name for all instances wanting the
		// same lock.
		mutexname = "sku_discount_item_" + strconv.FormatInt(v.ID, 10)
		mutex = rs.NewMutex(mutexname, redsync.WithRetryDelay(100*time.Millisecond), redsync.WithExpiry(60*time.Second))
		if e = mutex.Lock(); e != nil {
			e = echo.NewHTTPError(http.StatusLocked, e.Error())
			return ctx.Serve(e)
		}

		mutexs = append(mutexs, mutex)
	}

	ctx.ResponseData, e = Save(r)
	if e != nil && strings.Contains(e.Error(), "Failed to save quota") {
		qr.OrderChannel = 1
		qr.Branch = r.Branch
		qr.Products = r.Products
		qr.CurrentTime = r.CurrentTime
		if e = ctx.Bind(&qr); e != nil {
			return ctx.Serve(e)
		}
		ctx.ResponseData, e = Save(r)
	}

	for _, v := range mutexs {
		if _, err := v.Unlock(); err != nil {
			return ctx.Serve(err)
		}
	}

	return ctx.Serve(e)
}

// detail : function to get detailed data by id
func (h *Handler) detail(c echo.Context) (e error) {
	ctx := c.(*cuxs.Context)

	var id int64
	if id, e = ctx.Decrypt("id"); e == nil {
		if ctx.ResponseData, e = repository.GetSalesOrder("id", id); e != nil {
			e = echo.ErrNotFound
		}
	}

	return ctx.Serve(e)
}

// update : function to unarchive requested data based on parameters
func (h *Handler) update(c echo.Context) (e error) {
	ctx := c.(*cuxs.Context)
	var (
		r         updateRequest
		qr        quotaRequest
		mutexname string
		mutexs    []*redsync.Mutex
		rs        *redsync.Redsync
	)

	if r.Session, e = auth.UserSession(ctx); e != nil {
		return ctx.Serve(e)
	}

	if r.ID, e = ctx.Decrypt("id"); e != nil {
		return ctx.Serve(e)
	}

	if e = ctx.Bind(&r); e != nil {
		return ctx.Serve(e)
	}

	for _, v := range r.SkuDiscountItems {
		var mutex *redsync.Mutex
		//redsync-lock
		client := goredislib.NewClient(&goredislib.Options{
			Addr:     env.GetString("REDIS_HOST", "127.0.0.1:6379"),
			Password: env.GetString("REDIS_PASSWORD", "127.0.0.1:6379"),

			MaxRetries: 10,
		})

		defer client.Close()

		pool := goredis.NewPool(client)

		// Create an instance of redisync to be used to obtain a mutual exclusion
		// lock.
		rs = redsync.New(pool)

		// Obtain a new mutex by using the same name for all instances wanting the
		// same lock.
		mutexname = "sku_discount_item_" + strconv.FormatInt(v.ID, 10)
		mutex = rs.NewMutex(mutexname, redsync.WithRetryDelay(100*time.Millisecond), redsync.WithExpiry(60*time.Second))
		if e = mutex.Lock(); e != nil {
			e = echo.NewHTTPError(http.StatusLocked, e.Error())
			return ctx.Serve(e)
		}

		mutexs = append(mutexs, mutex)
	}

	ctx.ResponseData, e = Update(r)
	if e != nil && strings.Contains(e.Error(), "Failed to save quota") {
		qr.OrderChannel = r.SalesOrder.OrderChannel
		qr.Branch = r.SalesOrder.Branch
		qr.Products = r.Products
		qr.CurrentTime = r.CurrentTime
		if e = ctx.Bind(&qr); e != nil {
			return ctx.Serve(e)
		}

		ctx.ResponseData, e = Update(r)
	}

	for _, v := range mutexs {
		if _, err := v.Unlock(); err != nil {
			return ctx.Serve(err)
		}
	}

	return ctx.Serve(e)
}

// cancel : function to unarchive requested data based on parameters
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

// lock : function to lock SO
func (h *Handler) lock(c echo.Context) (e error) {
	ctx := c.(*cuxs.Context)
	var r lockRequest

	if r.Session, e = auth.UserSession(ctx); e == nil {
		if r.ID, e = ctx.Decrypt("id"); e == nil {
			if e = ctx.Bind(&r); e == nil {
				ctx.ResponseData, e = Lock(r)
			}
		}
	}

	return ctx.Serve(e)
}

// updatePriceDeliveredSO : function to update when SO status delivered
func (h *Handler) updatePriceDeliveredSO(c echo.Context) (e error) {
	ctx := c.(*cuxs.Context)
	var (
		r updatePriceRequestDeliveredSO
	)

	if r.Session, e = auth.UserSession(ctx); e != nil {
		return ctx.Serve(e)
	}

	if r.ID, e = ctx.Decrypt("id"); e != nil {
		return ctx.Serve(e)
	}

	if e = ctx.Bind(&r); e != nil {
		return ctx.Serve(e)
	}

	ctx.ResponseData, e = updatePriceDeliveredSO(r)

	return ctx.Serve(e)
}
