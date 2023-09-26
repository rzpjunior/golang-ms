// Copyright 2020 PT. Eden Pangan Indonesia. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

package order

import (
	"context"
	"strconv"
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

type updateRequest struct {
	ID              int64              `json:"-"`
	JobsID          primitive.ObjectID `json:"jobs_id"`
	SalesOrderID    string             `json:"sales_order_id" valid:"required"`
	RecognitionDate string             `json:"recognition_date" valid:"required"`
	ShippingAddress string             `json:"shipping_address" valid:"required"`
	WrtID           string             `json:"wrt_id" valid:"required"`
	SalesInvoiceID  string             `json:"sales_invoice_id"`
	TermInvoiceID   string             `json:"term_invoice_sls"`

	Note        string  `json:"note"`
	WarehouseID string  `json:"warehouse_id" valid:"required"`
	TotalWeight float64 `json:"-"`

	Token       string `json:"token"`
	TypeRequest string `json:"type"`
	EncryptID   string `json:"ID"`

	OrderDate time.Time `json:"-"`

	// FOR UPDATE INVOICE
	TotalPrice         float64 `json:"-"`
	TotalCharge        float64 `json:"-"`
	TotalWeightDirect  float64 `json:"-"`
	TotalSkuDiscAmount float64 `json:"-"`

	DeliveryOrderItems     []*itemRequest        `json:"delivery_order_items" valid:"required"`
	Session                *auth.SessionData     `json:"-"`
	RecognitionAt          time.Time             `json:"-"`
	SalesOrder             *model.SalesOrder     `json:"-"`
	DeliveryOrder          *model.DeliveryOrder  `json:"-"`
	Warehouse              *model.Warehouse      `json:"-"`
	SalesInvoice           []*model.SalesInvoice `json:"-"`
	TermInvoice            *model.InvoiceTerm    `json:"-"`
	CreditLimitBefore      float64               `json:"-"`
	CreditLimitAfter       float64               `json:"-"`
	IsCreateCreditLimitLog int64                 `json:"-"`

	CheckSalesInvoiceTermInvoice int8

	//RedisLock
	Client   *goredislib.Client
	Mutex    *redsync.Mutex
	UseRedis bool
}

func (c *updateRequest) Validate() *validation.Output {
	o := &validation.Output{Valid: true}
	o1 := orm.NewOrm()
	o1.Using("read_only")
	var e error
	var filter, exclude map[string]interface{}
	var duplicated = make(map[string]bool)
	var duplicatedoi = make(map[string]bool)
	var deliveryOrderItemID, productID, stockOpname int64

	c.UseRedis = false
	//redsync-lock
	c.Client = goredislib.NewClient(&goredislib.Options{
		Addr:       env.GetString("REDIS_HOST", "127.0.0.1:6379"),
		Password:   env.GetString("REDIS_PASSWORD", "127.0.0.1:6379"),
		MaxRetries: -1,
	})
	ctx := context.Background()
	res, e := c.Client.Ping(ctx).Result()
	if res != "" {
		c.UseRedis = true
		pool := goredis.NewPool(c.Client) // or, pool := redigo.NewPool(...)

		// Create an instance of redisync to be used to obtain a mutual exclusion
		// lock.
		rs := redsync.New(pool)

		// Obtain a new mutex by using the same name for all instances wanting the
		// same lock.

		mutexname := "delivery_order" + c.WarehouseID

		c.Mutex = rs.NewMutex(mutexname, redsync.WithRetryDelay(100*time.Millisecond), redsync.WithExpiry(60*time.Second), redsync.WithTries(64))

		if err := c.Mutex.Lock(); err != nil {
			o.Failure("id.invalid", "system is busy please try again later.")
			return o
		}
	}

	c.DeliveryOrder = &model.DeliveryOrder{ID: c.ID}
	if e = c.DeliveryOrder.Read("ID"); e == nil {
		if c.DeliveryOrder.Status != 1 && c.DeliveryOrder.Status != 5 && c.DeliveryOrder.Status != 6 && c.DeliveryOrder.Status != 7 {
			o.Failure("id.invalid", util.ErrorActive("delivery order"))
			return o
		}
	} else {
		o.Failure("id.invalid", util.ErrorInvalidData("delivery order"))
	}

	salesOrderID, _ := common.Decrypt(c.SalesOrderID)
	c.SalesOrder = &model.SalesOrder{ID: salesOrderID}
	if e = c.SalesOrder.Read("ID"); e != nil {
		o.Failure("id.invalid", util.ErrorInvalidData("sales order"))
	}

	warehouseID, _ := common.Decrypt(c.WarehouseID)
	c.Warehouse = &model.Warehouse{ID: warehouseID}
	c.Warehouse.Read("ID")

	// DIRECT INVOICE SHOULD HAVE AT LEAST 1 INVOICE - [INVOICE DOC CREATED WHEN CREATE SO]
	if c.SalesOrder.InvoiceTerm.ID == 2 {

		filter = map[string]interface{}{"sales_order_id": salesOrderID, "status": 1}
		exclude = map[string]interface{}{}
		c.SalesInvoice, _, e = repository.CheckSalesInvoicesData(filter, exclude)

		if len(c.SalesInvoice) > 0 {
			c.CheckSalesInvoiceTermInvoice = 1
		}

	}

	if err := c.SalesOrder.OrderType.Read("ID"); err != nil {
		o.Failure("order_type_id.invalid", util.ErrorInvalidData("order type"))
		return o
	}

	var stockType *model.Glossary
	if c.SalesOrder.OrderType.Name == "Zero Waste" {
		stockType, e = repository.GetGlossaryMultipleValue("table", "all", "attribute", "stock_type", "value_name", "waste stock")
		if e != nil {
			o.Failure("stock_type_id.invalid", util.ErrorInvalidData("stock type"))
			return o
		}
	} else {
		stockType, e = repository.GetGlossaryMultipleValue("table", "all", "attribute", "stock_type", "value_name", "good stock")
		if e != nil {
			o.Failure("stock_type_id.invalid", util.ErrorInvalidData("stock type"))
			return o
		}
	}

	o1.Raw("SELECT count(id) from stock_opname where warehouse_id = ? AND stock_type = ? AND status = 1", warehouseID, stockType.ValueInt).QueryRow(&stockOpname)

	if stockOpname > 0 {
		o.Failure("id.invalid", util.ErrorRelated("active", "stock opname", c.Warehouse.Name))

	}

	if len(c.Note) > 250 {
		o.Failure("note", util.ErrorCharLength("note", 250))
	}

	if c.RecognitionDate != "" {
		if c.OrderDate, e = time.Parse("2006-01-02", c.RecognitionDate); e != nil {
			o.Failure("recognition_date.invalid", util.ErrorInputRequired("delivery date"))
		}
	}

	for n, row := range c.DeliveryOrderItems {
		if row.ID != "" {

			if row.DeliverQty < 0 {
				o.Failure("qty"+strconv.Itoa(n)+".greater", util.ErrorGreater("product quantity", "0"))
			}
			if !duplicatedoi[row.ID] {
				deliveryOrderItemID, _ = common.Decrypt(row.ID)
				row.DeliveryOrderItem = &model.DeliveryOrderItem{ID: deliveryOrderItemID}

				if e = row.DeliveryOrderItem.Read("ID"); e != nil {
					o.Failure("product_id"+strconv.Itoa(n)+".invalid", util.ErrorInvalidData("product"))
				}

			} else {
				o.Failure("product_id"+strconv.Itoa(n)+".duplicate", util.ErrorDuplicate("product"))
			}

		}

		if row.ProductID != "" {
			if !duplicated[row.ProductID] {
				productID, _ = common.Decrypt(row.ProductID)
				row.Product = &model.Product{ID: productID}

				if e = row.Product.Read("ID"); e != nil {
					o.Failure("product_id"+strconv.Itoa(n)+".invalid", util.ErrorInputRequired("product"))
				} else {
					filter = map[string]interface{}{"product_id": productID, "warehouse_id": c.Warehouse.ID, "status": 1}
					if _, countStock, err := repository.CheckStockData(filter, exclude); err == nil && countStock == 0 {
						o.Failure("product_id"+strconv.Itoa(n)+".invalid", util.ErrorProductMustAvailable())
					}

					var salesInvoiceItemQty float64
					if row.DeliverQty >= row.OrderQty {
						salesInvoiceItemQty = row.OrderQty
					} else {
						salesInvoiceItemQty = row.DeliverQty
					}

					// FOR UPDATE INVOICE
					if c.CheckSalesInvoiceTermInvoice == 1 {
						if e = o1.Raw("SELECT * from sales_invoice_item where product_id = ? and sales_invoice_id = ? ", row.Product.ID, c.SalesInvoice[0].ID).QueryRow(&row.SalesInvoiceItem); e == nil {
							row.UnitPrice = row.SalesInvoiceItem.UnitPrice
							row.TotalPrice = common.Rounder(salesInvoiceItemQty*row.UnitPrice, 0.5, 2)

							if e = row.SalesInvoiceItem.SalesOrderItem.Read("ID"); e == nil {
								if row.SalesInvoiceItem.SalesOrderItem.UnitPriceDiscount > 0 {
									discQty := row.SalesInvoiceItem.SalesOrderItem.DiscountQty
									if salesInvoiceItemQty < discQty {
										discQty = salesInvoiceItemQty
									}
									row.SkuDiscAmount = row.SalesInvoiceItem.SalesOrderItem.UnitPriceDiscount * discQty
									row.TotalPrice -= row.SkuDiscAmount
								}

								row.SalesInvoiceItem.InvoiceQty = salesInvoiceItemQty
								row.SalesInvoiceItem.UnitPrice = row.UnitPrice
								row.SalesInvoiceItem.Subtotal = row.TotalPrice
								row.SalesInvoiceItem.SkuDiscAmount = row.SkuDiscAmount
							}
						}

						row.Subtotal = row.TotalPrice
						row.Weight = row.DeliverQty * row.Product.UnitWeight
						c.TotalWeightDirect = common.Rounder(c.TotalWeight+row.Weight, 0.5, 2)

						if c.SalesOrder.Voucher != nil {
							c.SalesOrder.Voucher.Read("ID")
							if c.SalesOrder.Voucher.Type != 4 {
								row.DiscountAmount = common.Rounder(c.SalesOrder.Voucher.DiscAmount, 0.5, 2)
							}
						}
						if c.SalesOrder.PointRedeemAmount != 0 {
							row.DiscountAmount = row.DiscountAmount + common.Rounder(c.SalesOrder.PointRedeemAmount, 0.5, 2)
						}
						c.TotalSkuDiscAmount += row.SkuDiscAmount
						c.TotalPrice += row.Subtotal
						c.TotalCharge = common.Rounder(c.TotalPrice-row.DiscountAmount, 0.5, 2)
					}

					row.Product.Category.Read("ID")
					row.Weight = row.DeliverQty * row.Product.UnitWeight
					c.TotalWeight += row.Weight

				}
				duplicated[row.ProductID] = true
			} else {
				o.Failure("product_id"+strconv.Itoa(n)+".duplicate", util.ErrorDuplicate("product"))
			}
		} else {
			o.Failure("product_id"+strconv.Itoa(n)+".invalid", util.ErrorInvalidData("product"))
		}

	}

	c.TotalCharge = common.Rounder(c.TotalCharge+c.SalesOrder.DeliveryFee, 0.5, 2)

	if e = c.DeliveryOrder.SalesOrder.Read("ID"); e != nil {
		o.Failure("id.invalid", util.ErrorInvalidData("sales order"))
	}

	if e = c.DeliveryOrder.SalesOrder.Branch.Read("ID"); e != nil {
		o.Failure("id.invalid", util.ErrorInvalidData("branch"))
	}

	if e = c.DeliveryOrder.SalesOrder.Branch.Merchant.Read("ID"); e != nil {
		o.Failure("id.invalid", util.ErrorInvalidData("merchant"))
	}

	return o
}

// Messages : function to return error validation messages
func (c *updateRequest) Messages() map[string]string {
	return map[string]string{
		"sales_order_id.required":   util.ErrorInputRequired("sales order"),
		"warehouse_id.required":     util.ErrorInputRequired("warehouse"),
		"wrt_id.required":           util.ErrorInputRequired("wrt"),
		"recognition_date.required": util.ErrorInputRequired("recognition date"),
		"delivery_date.required":    util.ErrorInputRequired("delivery date"),
		"shipping_address.required": util.ErrorInputRequired("shipping address"),
	}
}
