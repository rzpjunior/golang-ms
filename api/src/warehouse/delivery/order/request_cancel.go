// Copyright 2020 PT. Qasico Teknologi Indonesia. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

package order

import (
	"context"
	"time"

	"git.edenfarm.id/cuxs/common"
	"git.edenfarm.id/cuxs/env"
	"git.edenfarm.id/cuxs/orm"
	"git.edenfarm.id/cuxs/validation"
	"git.edenfarm.id/project-version2/api/datastore/repository"
	"git.edenfarm.id/project-version2/api/src/auth"
	"git.edenfarm.id/project-version2/api/util"
	"git.edenfarm.id/project-version2/datamodel/model"
	goredislib "github.com/go-redis/redis/v8"
	"github.com/go-redsync/redsync/v4"
	"github.com/go-redsync/redsync/v4/redis/goredis/v8"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type cancelRequest struct {
	ID     int64              `json:"-"`
	JobsID primitive.ObjectID `json:"jobs_id"`
	Note   string             `json:"note" valid:"required"`

	Token       string `json:"token"`
	TypeRequest string `json:"type"`
	EncryptID   string `json:"ID"`

	DeliveryOrder     *model.DeliveryOrder       `json:"-"`
	DeliveryOrderItem []*model.DeliveryOrderItem `json:"-"`
	Stock             []*model.Stock             `json:"-"`
	StockLog          []*model.StockLog          `json:"-"`
	WasteLog          []*model.WasteLog          `json:"-"`
	Session           *auth.SessionData          `json:"-"`
	//RedisLock
	Client   *goredislib.Client
	Mutex    *redsync.Mutex
	UseRedis bool
}

func (r *cancelRequest) Validate() *validation.Output {
	o := &validation.Output{Valid: true}
	o1 := orm.NewOrm()
	o1.Using("read_only")
	var stockOpname int64
	var stockType *model.Glossary
	var err error
	r.UseRedis = false

	//redsync-lock
	r.Client = goredislib.NewClient(&goredislib.Options{
		Addr:       env.GetString("REDIS_HOST", "127.0.0.1:6379"),
		Password:   env.GetString("REDIS_PASSWORD", "127.0.0.1:6379"),
		MaxRetries: -1,
	})
	ctx := context.Background()
	res, _ := r.Client.Ping(ctx).Result()
	if res != "" {
		r.UseRedis = true
		pool := goredis.NewPool(r.Client) // or, pool := redigo.NewPool(...)

		// Create an instance of redisync to be used to obtain a mutual exclusion
		// lock.
		rs := redsync.New(pool)

		// Obtain a new mutex by using the same name for all instances wanting the
		// same lock.
		warehouseID := common.Encrypt(r.DeliveryOrder.Warehouse.ID)

		mutexname := "delivery_order" + warehouseID

		r.Mutex = rs.NewMutex(mutexname, redsync.WithRetryDelay(100*time.Millisecond), redsync.WithExpiry(60*time.Second), redsync.WithTries(64))

		if err := r.Mutex.Lock(); err != nil {
			o.Failure("id.invalid", "system is busy please try again later.")
			return o
		}
	}

	if r.DeliveryOrder.Status != 1 && r.DeliveryOrder.Status != 5 && r.DeliveryOrder.Status != 6 && r.DeliveryOrder.Status != 7 {
		o.Failure("status.inactive", util.ErrorActive("delivery order"))
		return o
	}

	if err := r.DeliveryOrder.SalesOrder.Read("ID"); err != nil {
		o.Failure("sales_order_id.invalid", util.ErrorInvalidData("sales order"))
		return o
	}

	if err := r.DeliveryOrder.SalesOrder.OrderType.Read("ID"); err != nil {
		o.Failure("order_type_id.invalid", util.ErrorInvalidData("order type"))
		return o
	}

	if r.DeliveryOrder.SalesOrder.OrderType.Name == "Zero Waste" {
		stockType, err = repository.GetGlossaryMultipleValue("table", "all", "attribute", "stock_type", "value_name", "waste stock")
		if err != nil {
			o.Failure("stock_type_id.invalid", util.ErrorInvalidData("stock type"))
			return o
		}
	} else {
		stockType, err = repository.GetGlossaryMultipleValue("table", "all", "attribute", "stock_type", "value_name", "good stock")
		if err != nil {
			o.Failure("stock_type_id.invalid", util.ErrorInvalidData("stock type"))
			return o
		}
	}
	o1.Raw("SELECT count(id) from stock_opname where warehouse_id = ? AND stock_type = ? AND status = 1", r.DeliveryOrder.Warehouse.ID, stockType.ValueInt).QueryRow(&stockOpname)

	if stockOpname > 0 {
		o.Failure("id.invalid", util.ErrorRelated("active", "stock opname", r.DeliveryOrder.Warehouse.Name))

	}

	return o
}

func (r *cancelRequest) Messages() map[string]string {
	return map[string]string{
		"note.required": util.ErrorInputRequired("cancellation note"),
	}
}
