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
)

// createRequest : struct to hold price set request data
type createRequest struct {
	Code            string `json:"-"`
	JobsID          string `json:"jobs_id"`
	SalesOrderID    string `json:"sales_order_id" valid:"required"`
	WarehouseID     string `json:"warehouse_id" valid:"required"`
	WrtID           string `json:"wrt_id" valid:"required"`
	RecognitionDate string `json:"recognition_date" valid:"required"`
	ShippingAddress string `json:"shipping_address" valid:"required"`
	Note            string `json:"note"`

	Token       string `json:"token"`
	TypeRequest string `json:"type"`

	OrderDate         time.Time `json:"-"`
	TotalWeight       float64   `json:"-"`
	TotalWeightDirect float64   `json:"-"`

	SalesOrder         *model.SalesOrder `json:"-"`
	Warehouse          *model.Warehouse  `json:"-"`
	Wrt                *model.Wrt        `json:"-"`
	Voucher            *model.Voucher    `json:"-"`
	Stock              *model.Stock      `json:"-"`
	DeliveryOrderItems []*itemRequest    `json:"delivery_order_items" valid:"required"`
	Session            *auth.SessionData `json:"-"`

	// FOR CREATE INVOICE ; IF INVOICE TERM = DIRECT_INVOICE
	DiscountAmount         float64 `json:"-"`
	TotalPrice             float64 `json:"-"`
	TotalCharge            float64 `json:"-"`
	DocumentCodeInvoice    string  `json:"-"`
	Discount               float32 `json:"-"`
	NoteInvoice            string  `json:"-"`
	AdjustmentAmount       float64 `json:"-"`
	AdjustmentNote         string  `json:"-"`
	TotalSkuDiscAmount     float64 `json:"-"`
	CreditLimitBefore      float64 `json:"-"`
	CreditLimitAfter       float64 `json:"-"`
	IsCreateCreditLimitLog int64   `json:"-"`

	//RedisLock
	Client   *goredislib.Client
	Mutex    *redsync.Mutex
	UseRedis bool
}

type itemRequest struct {
	ID          string  `json:"id"`
	ProductID   string  `json:"product_id"`
	ProductCode string  `json:"product_code"`
	ProductName string  `json:"product_name"`
	Uom         string  `json:"uom"`
	OrderQty    float64 `json:"order_qty"`
	DeliverQty  float64 `json:"deliver_qty"`
	Note        string  `json:"note" valid:"lte:255"`
	Weight      float64 `json:"weight"`

	DeliveryOrderItem *model.DeliveryOrderItem `json:"-"`
	Product           *model.Product           `json:"-"`
	SalesInvoiceItem  *model.SalesInvoiceItem  `json:"-"`

	// INVOICE
	TotalPrice         float64 `json:"-"`
	DiscountAmount     float64 `json:"-"`
	DiscountPercentage float32 `json:"-"`
	UnitPrice          float64 `json:"-"`
	Subtotal           float64 `json:"-"`
	SkuDiscAmount      float64 `json:"-"`
}

// Validate : function to validate uom request data
func (c *createRequest) Validate() *validation.Output {
	o := &validation.Output{Valid: true}
	var e error
	var filter, exclude map[string]interface{}
	var count, stockOpname int64
	var totalChargeDifferences float64 = 0

	orSelect := orm.NewOrm()
	orSelect.Using("read_only")
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

	salesOrderID, _ := common.Decrypt(c.SalesOrderID)
	c.SalesOrder = &model.SalesOrder{ID: salesOrderID}
	if e = c.SalesOrder.Read("ID"); e != nil {
		o.Failure("id.invalid", util.ErrorInvalidData("sales order"))
		return o
	}
	orSelect.Raw("SELECT count(id) from delivery_order where sales_order_id = ? AND status IN (1,2)", salesOrderID).QueryRow(&count)
	if count > 0 {
		o.Failure("id.invalid", util.ErrorCreateDoc("sales order", "delivery order"))
	}
	if c.SalesOrder.Status != 1 && c.SalesOrder.Status != 9 && c.SalesOrder.Status != 12 {
		o.Failure("id.invalid", util.ErrorCreateDoc("sales order", "delivery order"))
	}

	c.SalesOrder.SalesTerm.Read("ID")
	if c.SalesOrder.Status == 1 && c.SalesOrder.SalesTerm.ID == 11 {
		o.Failure("id.invalid", util.ErrorCreateDocStatus("delivery order", " payment", "invalid"))
	}

	warehouseID, _ := common.Decrypt(c.WarehouseID)
	c.Warehouse = &model.Warehouse{ID: warehouseID}
	c.Warehouse.Read("ID")

	wrtID, _ := common.Decrypt(c.WrtID)
	c.Wrt = &model.Wrt{ID: wrtID}
	c.Wrt.Read("ID")

	if c.RecognitionDate != "" {
		if c.OrderDate, e = time.Parse("2006-01-02", c.RecognitionDate); e != nil {
			o.Failure("recognition_date.invalid", util.ErrorInvalidData("delivery date"))
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

	orSelect.Raw("SELECT count(id) from stock_opname where warehouse_id = ? AND stock_type = ? AND status = 1", warehouseID, stockType.ValueInt).QueryRow(&stockOpname)

	if stockOpname > 0 {
		o.Failure("id.invalid", util.ErrorRelated("active", "stock opname", c.Warehouse.Name))

	}

	var duplicated = make(map[string]bool)

	for n, row := range c.DeliveryOrderItems {
		var productID int64

		if row.ProductID != "" {
			if !duplicated[row.ProductID] {
				if row.DeliverQty < 0 {
					o.Failure("qty"+strconv.Itoa(n)+".greater", util.ErrorGreater("product quantity", "0"))
				}

				productID, _ = common.Decrypt(row.ProductID)
				row.Product = &model.Product{ID: productID}

				if e = row.Product.Read("ID"); e != nil {
					o.Failure("product_id"+strconv.Itoa(n)+".invalid", util.ErrorInputRequired("product"))
				} else {
					filter = map[string]interface{}{"product_id": productID, "warehouse_id": warehouseID, "status": 1}
					if _, countStock, err := repository.CheckStockData(filter, exclude); err == nil && countStock == 0 {
						o.Failure("product_id"+strconv.Itoa(n)+".invalid", util.ErrorProductMustAvailable())
					}

					var salesInvoiceItemQty float64
					if row.DeliverQty >= row.OrderQty {
						salesInvoiceItemQty = row.OrderQty
					} else {
						salesInvoiceItemQty = row.DeliverQty
					}

					// CREATE INVOICE
					if c.SalesOrder.InvoiceTerm.ID == 2 {
						orSelect.Raw("SELECT * from sales_order_item where product_id = ? and sales_order_id = ? ", row.Product.ID, c.SalesOrder.ID).QueryRows(&c.SalesOrder.SalesOrderItems)
						for _, i := range c.SalesOrder.SalesOrderItems {

							row.UnitPrice = i.UnitPrice
							row.Subtotal = common.Rounder(salesInvoiceItemQty*row.UnitPrice, 0.5, 2)
							row.Weight = row.DeliverQty * row.Product.UnitWeight
							c.TotalWeightDirect = common.Rounder(c.TotalWeight+row.Weight, 0.5, 2)

							if i.SkuDiscountItem.ID != 0 {
								if e = i.SkuDiscountItem.Read(("ID")); e != nil {
									o.Failure("product_id"+strconv.Itoa(n)+".invalid", util.ErrorInvalidData("sku discount item"))
									continue
								}

								discQty := i.DiscountQty
								if salesInvoiceItemQty < discQty {
									discQty = salesInvoiceItemQty
								}
								row.SkuDiscAmount = discQty * i.UnitPriceDiscount
								row.Subtotal -= row.SkuDiscAmount
							}

							row.TotalPrice = row.Subtotal
						}

						if c.SalesOrder.Voucher != nil {
							c.SalesOrder.Voucher.Read("ID")
							row.DiscountAmount = common.Rounder(c.SalesOrder.VouDiscAmount, 0.5, 2)
						}

						if c.SalesOrder.PointRedeemAmount != 0 {
							row.DiscountAmount = row.DiscountAmount + common.Rounder(c.SalesOrder.PointRedeemAmount, 0.5, 2)
						}

						c.TotalSkuDiscAmount += row.SkuDiscAmount
						c.TotalPrice += row.TotalPrice
						c.TotalCharge = common.Rounder(c.TotalPrice-row.DiscountAmount, 0.5, 2)
					}

					row.Weight = row.DeliverQty * row.Product.UnitWeight
					c.TotalWeight = c.TotalWeight + row.Weight
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

	if len(c.Note) > 250 {
		o.Failure("note", util.ErrorCharLength("note", 250))
	}

	if e = c.SalesOrder.Branch.Read("ID"); e != nil {
		o.Failure("id.invalid", util.ErrorInvalidData("branch"))
		return o
	}

	if e = c.SalesOrder.Branch.Merchant.Read("ID"); e != nil {
		o.Failure("id.invalid", util.ErrorInvalidData("merchant"))
		return o
	}

	if c.SalesOrder.InvoiceTerm.ID == 2 {

		c.CreditLimitBefore = c.SalesOrder.Branch.Merchant.RemainingCreditLimitAmount

		totalChargeDifferences = c.TotalCharge - c.SalesOrder.TotalCharge
		c.CreditLimitAfter = c.CreditLimitBefore - totalChargeDifferences

		if c.TotalCharge < c.SalesOrder.TotalCharge {
			totalChargeDifferences = c.SalesOrder.TotalCharge - c.TotalCharge
			c.CreditLimitAfter = c.CreditLimitBefore + totalChargeDifferences
		}

		if c.CreditLimitAfter < 0 && c.CreditLimitBefore > 0 {
			o.Failure("credit_limit_amount.invalid", util.ErrorCreditLimitExceeded(c.SalesOrder.Branch.Merchant.Name))
		}

		if totalChargeDifferences > 0 {
			c.IsCreateCreditLimitLog = 1
		}
	}
	return o
}

// Messages : function to return error validation messages
func (c *createRequest) Messages() map[string]string {
	messages := map[string]string{
		"sales_order_id.required":   util.ErrorInputRequired("sales order"),
		"warehouse_id.required":     util.ErrorInputRequired("warehouse"),
		"wrt_id.required":           util.ErrorInputRequired("wrt"),
		"recognition_date.required": util.ErrorInputRequired("recognition date"),
		"delivery_date.required":    util.ErrorInputRequired("delivery date"),
		"shipping_address.required": util.ErrorInputRequired("shipping address"),
	}

	return messages
}
