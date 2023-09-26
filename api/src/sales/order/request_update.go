// Copyright 2020 PT. Eden Pangan Indonesia. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

package order

import (
	"fmt"
	"strconv"
	"strings"
	"time"

	"git.edenfarm.id/cuxs/common"
	"git.edenfarm.id/cuxs/orm"
	"git.edenfarm.id/cuxs/validation"
	"git.edenfarm.id/project-version2/datamodel/model"

	"git.edenfarm.id/project-version2/api/datastore/repository"
	"git.edenfarm.id/project-version2/api/src/auth"
	"git.edenfarm.id/project-version2/api/util"
)

// updateRequest : struct to hold sales order request data
type updateRequest struct {
	ID                     int64     `json:"-"`
	DeliveryDateStr        string    `json:"delivery_date" valid:"required"`
	WrtID                  string    `json:"wrt_id" valid:"required"`
	RecognitionDateStr     string    `json:"order_date" valid:"required"`
	OrderTypeID            string    `json:"order_type_id" valid:"required"`
	SalespersonID          string    `json:"salesperson_id" valid:"required"`
	SalesTermID            string    `json:"term_payment_sls_id" valid:"required"`
	InvoiceTermID          string    `json:"term_invoice_sls_id" valid:"required"`
	PaymentGroupID         string    `json:"payment_group_id" valid:"required"`
	WarehouseID            string    `json:"warehouse_id" valid:"required"`
	RedeemCode             string    `json:"redeem_code"`
	Note                   string    `json:"note"`
	UpdateAll              int8      `json:"update_all"`
	TotalPrice             float64   `json:"-"`
	TotalWeight            float64   `json:"-"`
	TotalCharge            float64   `json:"-"`
	DeliveryFee            float64   `json:"-"`
	SameTagCustomer        string    `json:"-"`
	DeliveryDate           time.Time `json:"-"`
	RecognitionDate        time.Time `json:"-"`
	TotalSkuDiscAmount     float64   `json:"-"`
	CurrentTime            time.Time `json:"-"`
	NotePriceChange        string    `json:"-"`
	CreditLimitBefore      float64   `json:"-"`
	CreditLimitAfter       float64   `json:"-"`
	OldTotalCharge         float64   `json:"-"`
	IsCreateCreditLimitLog int64     `json:"-"`

	Wrt                *model.Wrt                `json:"-"`
	OrderType          *model.OrderType          `json:"-"`
	Salesperson        *model.Staff              `json:"-"`
	SalesTerm          *model.SalesTerm          `json:"-"`
	InvoiceTerm        *model.InvoiceTerm        `json:"-"`
	PaymentGroup       *model.PaymentGroup       `json:"-"`
	Voucher            *model.Voucher            `json:"-"`
	AreaPolicy         *model.AreaPolicy         `json:"-"`
	SalesOrder         *model.SalesOrder         `json:"-"`
	PackingOrder       *model.PackingOrder       `json:"-"`
	Warehouse          *model.Warehouse          `json:"-"`
	VoucherLog         []*model.VoucherLog       `json:"-"`
	DayOff             *model.DayOff             `json:"-"`
	SkuDiscountItems   []*model.SkuDiscountItem  `json:"-"`
	AreaBusinessPolicy *model.AreaBusinessPolicy `json:"-"`

	Products []*salesOrderItem `json:"products" valid:"required"`

	Session *auth.SessionData `json:"-"`
}

// Validate : function to validate sales order request data
func (r *updateRequest) Validate() *validation.Output {
	o := &validation.Output{Valid: true}
	o1 := orm.NewOrm()
	o1.Using("read_only")
	var (
		err                                                                                                     error
		filtersoi, excludesoi, filter, exclude                                                                  map[string]interface{}
		paymentTermID, invoiceTermID, paymentGroupID, wrtID, orderTypeID, salespersonID, productID, warehouseID int64
		totalMatch, countMatch                                                                                  int
		checkProduct                                                                                            *model.ProductPush
		priceChanged                                                                                            string
		isPriceSetChanged                                                                                       bool
		totalChargeDifferences, minOrderQty, minUnitPrice                                                       float64
		warehouseSelfPickUp                                                                                     *model.Warehouse
		wrtType                                                                                                 int8
	)
	productList, discAmount := make(map[int64]string), float64(0)
	r.CurrentTime = time.Now()

	if r.SalesOrder, err = repository.ValidSalesOrder(r.ID); err != nil {
		o.Failure("id.invalid", util.ErrorInvalidData("sales order"))
		return o
	}

	r.OldTotalCharge = r.SalesOrder.TotalCharge

	if r.SalesOrder.Status != 1 {
		o.Failure("id.inactive", util.ErrorActive("sales order"))
		return o
	}

	if r.SalesOrder.HasExtInvoice != 2 {
		o.Failure("id.invalid", util.ErrorStatusDoc("sales order", "updated", "Invoice Xendit"))
		return o
	}

	if r.SalesOrder.LockedBy != 0 {
		Lock, _ := repository.ValidStaff(r.SalesOrder.LockedBy)

		if r.SalesOrder.IsLocked != 2 {
			if r.SalesOrder.LockedBy != r.Session.Staff.ID {
				o.Failure("id.invalid", util.ErrorOrderTypeCantUpdate(Lock.Name))
			}
		}
	}

	if r.SalesOrder.Branch, err = repository.ValidBranch(r.SalesOrder.Branch.ID); err != nil {
		o.Failure("id.invalid", util.ErrorInvalidData("branch"))
		return o
	}

	if r.SalesOrder.Branch.Status != 1 {
		o.Failure("id.inactive", util.ErrorActive("branch"))
		return o
	}

	if r.SalesOrder.Branch.Merchant, err = repository.ValidMerchant(r.SalesOrder.Branch.Merchant.ID); err != nil {
		o.Failure("id.invalid", util.ErrorInvalidData("merchant"))
		return o
	}

	if r.SalesOrder.Branch.Merchant.Status != 1 {
		o.Failure("id.inactive", util.ErrorActive("merchant"))
		return o
	}

	if r.SalesOrder.Branch.Area, err = repository.ValidArea(r.SalesOrder.Branch.Area.ID); err != nil {
		o.Failure("id.invalid", util.ErrorInvalidData("price set"))
		return o
	}

	if r.AreaPolicy, err = repository.GetAreaPolicy("area_id", r.SalesOrder.Branch.Area.ID); err != nil {
		o.Failure("area_id.invalid", util.ErrorInvalidData("area"))
		return o
	}

	if r.AreaBusinessPolicy, err = repository.GetAreaBusinessPolicyDelivery(r.SalesOrder.Branch.Area.ID, r.SalesOrder.Branch.Merchant.BusinessType.ID); err != nil {
		o.Failure("area_business_config_id.invalid", util.ErrorInvalidData("area business config"))
		return o
	}

	r.DeliveryFee = r.AreaBusinessPolicy.DeliveryFee

	if r.SalesOrder.Branch.PriceSet, err = repository.ValidPriceSet(r.SalesOrder.Branch.PriceSet.ID); err != nil {
		o.Failure("id.invalid", util.ErrorInvalidData("price set"))
	}

	if r.SalesOrder.Archetype, err = repository.ValidArchetype(r.SalesOrder.Archetype.ID); err != nil {
		o.Failure("id.invalid", util.ErrorInvalidData("archetype"))
	}

	if warehouseID, err = common.Decrypt(r.WarehouseID); err != nil {
		o.Failure("warehouse_id.invalid", util.ErrorInvalidData("warehouse"))
	}

	if r.Warehouse, err = repository.ValidWarehouse(warehouseID); err != nil {
		o.Failure("warehouse_id.invalid", util.ErrorInvalidData("warehouse"))
	}

	if r.SalesOrder.Voucher != nil && r.SalesOrder.Voucher.ID != 0 {
		if r.SalesOrder.Voucher, err = repository.ValidVoucher(r.SalesOrder.Voucher.ID); err != nil {
			o.Failure("redeem_code.invalid", util.ErrorInvalidData("voucher"))
		}

		filter = map[string]interface{}{
			"voucher_id":     r.SalesOrder.Voucher.ID,
			"merchant_id":    r.SalesOrder.Branch.Merchant.ID,
			"branch_id":      r.SalesOrder.Branch.ID,
			"sales_order_id": r.SalesOrder.ID,
			"status":         int8(1),
		}
		r.VoucherLog, _, err = repository.CheckVoucherLogData(filter, exclude)
	}

	if paymentTermID, err = common.Decrypt(r.SalesTermID); err == nil {
		if r.SalesTerm, err = repository.ValidSalesTerm(paymentTermID); err != nil {
			o.Failure("term_payment_sls_id.invalid", util.ErrorInvalidData("sales term"))
		}
	}

	if invoiceTermID, err = common.Decrypt(r.InvoiceTermID); err == nil {
		if r.InvoiceTerm, err = repository.ValidInvoiceTerm(invoiceTermID); err != nil {
			o.Failure("term_invoice_sls_id.invalid", util.ErrorInvalidData("invoice term"))
		}
	}

	if paymentGroupID, err = common.Decrypt(r.PaymentGroupID); err == nil {
		if r.PaymentGroup, err = repository.ValidPaymentGroup(paymentGroupID); err != nil {
			o.Failure("payment_group.invalid", util.ErrorInvalidData("payment group"))
		}
	}

	if wrtID, err = common.Decrypt(r.WrtID); err == nil {
		if r.Wrt, err = repository.ValidWrt(wrtID); err != nil {
			o.Failure("wrt.invalid", util.ErrorInvalidData("wrt"))
		}
		wrtType = r.Wrt.Type
	}

	if orderTypeID, err = common.Decrypt(r.OrderTypeID); err == nil {
		if r.OrderType, err = repository.ValidOrderType(orderTypeID); err != nil {
			o.Failure("order_type.invalid", util.ErrorInvalidData("order type"))
		}
		minOrderQty = r.OrderType.MinOrderQty
		minUnitPrice = r.OrderType.MinUnitPrice
	}

	if salespersonID, err = common.Decrypt(r.SalespersonID); err == nil {
		if r.Salesperson, err = repository.ValidStaff(salespersonID); err == nil {
			if r.Salesperson.Status != 1 {
				o.Failure("salesperson_id.inactive", util.ErrorActive("salesperson"))
			}
		} else {
			o.Failure("salesperson_id.invalid", util.ErrorInvalidData("salesperson"))
		}
	}

	if len(r.Note) > 250 {
		o.Failure("note", util.ErrorCharLength("note", 250))
	}

	layout := "2006-01-02"
	loc, _ := time.LoadLocation("Asia/Jakarta")
	if r.RecognitionDate, err = time.ParseInLocation(layout, r.RecognitionDateStr, loc); err != nil {
		o.Failure("order_date.invalid", util.ErrorInvalidData("order date"))
	}

	if r.DeliveryDate, err = time.ParseInLocation(layout, r.DeliveryDateStr, loc); err != nil {
		o.Failure("delivery_date.invalid", util.ErrorInvalidData("delivery date"))
	}

	if r.RedeemCode != "" {
		r.Voucher = &model.Voucher{RedeemCode: r.RedeemCode, Status: 1}
		if err = r.Voucher.Read("RedeemCode", "Status"); err != nil {
			o.Failure("redeem_code.invalid", util.ErrorNotFound("voucher"))
			return o
		}

		if _, err = o1.LoadRelated(r.Voucher, "VoucherItems", 0); err != nil {
			o.Failure("redeem_code.invalid", util.ErrorInvalidData("voucher"))
		}
	} else {
		r.Voucher = nil
	}

	// Validation: voucher suits with the following merchant
	if r.Voucher != nil && r.Voucher.MerchantID != 0 {
		if r.Voucher.MerchantID != r.SalesOrder.Branch.Merchant.ID {
			o.Failure("redeem_code.invalid", util.ErrorNotValidFor("voucher", "merchant"))
		}
	}

	// check whether will update item or not
	if r.UpdateAll == 1 {
		for i, v := range r.Products {
			var (
				listOrder               []float64
				totalOrder, maxRemQuota float64
				maxDayDeliveryDate      int64
				maxAvailableDate        time.Time
				dayCount                int8
				salesOrderItem          *model.SalesOrderItem
			)

			// if product_qty below the sales_order_type.min_order_qty
			if v.Quantity < minOrderQty {
				o.Failure("qty"+strconv.Itoa(i)+".equalorgreater", util.ErrorInvalidData("order qty"))
			}

			// if product_price below the sales_order_type.min_unit_price
			if float64(v.UnitPrice) < minUnitPrice {
				o.Failure("unit_price"+strconv.Itoa(i)+".equalorgreater", util.ErrorInvalidData("unit price"))
			}

			if len(v.Note) > 100 {
				o.Failure("note"+strconv.Itoa(i), util.ErrorCharLength("note", 100))
			}

			if v.ProductID == "" {
				o.Failure("product_id"+strconv.Itoa(i)+".invalid", util.ErrorInputRequired("product"))
			} else {
				if productID, err = common.Decrypt(v.ProductID); err == nil {
					if v.Product, err = repository.ValidProduct(productID); err == nil {
						//Check UOM Decimal
						v.Product.Uom.Read("ID")
						if v.Product.Uom.DecimalEnabled == 2 {
							if v.Quantity != float64((int64(v.Quantity))) {
								o.Failure("qty"+strconv.Itoa(i)+".invalid", util.ErrorInvalidData("product quantity"))
							}
						}

						v.TaxableItem = v.Product.Taxable
						v.TaxPercentage = v.Product.TaxPercentage

						v.Price = &model.Price{Product: v.Product, PriceSet: r.SalesOrder.Branch.PriceSet}
						if err = v.Price.Read("Product", "PriceSet"); err != nil {
							o.Failure("price.invalid", util.ErrorInvalidData("price"))
						}

						v.DefaultPrice = v.Price.UnitPrice

						if v.UnitPrice != int64(v.DefaultPrice) {
							priceChanged += fmt.Sprintf("%s: Default Unit Price: Rp %g - Changed Unit Price: Rp %d, ", v.Product.Name, v.DefaultPrice, v.UnitPrice)
							isPriceSetChanged = true
						}

						r.TotalWeight = r.TotalWeight + (v.Quantity * v.Product.UnitWeight)

						if _, exist := productList[productID]; exist {
							o.Failure("product_id"+strconv.Itoa(i)+".duplicate", util.ErrorDuplicate("product"))
						} else {
							productList[productID] = "t"
						}

						filtersoi = map[string]interface{}{"sales_order_id": r.SalesOrder.ID, "product_id": productID}
						if _, countSOI, err := repository.CheckSalesOrderItemData(filtersoi, excludesoi); err == nil && countSOI == 0 {

							// Check for order type name is zero waste, don't check salable for stock
							if orderTypeID == 5 {
								filter = map[string]interface{}{"product_id": productID, "warehouse_id": warehouseID}
							} else {
								filter = map[string]interface{}{"product_id": productID, "warehouse_id": warehouseID, "salable": 1}
							}

							// Check stock data
							if _, countStock, err := repository.CheckStockData(filter, exclude); err == nil && countStock == 0 {
								o.Failure("product_id"+strconv.Itoa(i)+".invalid", util.ErrorProductMustAvailable())
							}

						}

						//Check daily limit
						if v.Product.OrderMaxQty != 0 {
							if v.Quantity > v.Product.OrderMaxQty {
								o.Failure("qty"+strconv.Itoa(i)+".invalid", util.ErrorInvalidData("max qty order"))
							} else {
								//Check product from another so (based on branch & delivery_date) still under maximum order qty
								o1.Raw("SELECT soi.order_qty "+
									"FROM sales_order_item soi "+
									"JOIN sales_order so ON so.id = soi.sales_order_id "+
									"JOIN branch b ON b.id = so.branch_id "+
									"JOIN merchant m ON m.id = b.merchant_id "+
									"WHERE m.id = ? AND soi.product_id = ? AND so.delivery_date = ? AND so.id NOT IN (?) AND so.status NOT IN (3,4)", r.SalesOrder.Branch.Merchant.ID, v.Product.ID, r.DeliveryDate, r.SalesOrder.ID).QueryRows(&listOrder)

								for _, v := range listOrder {
									totalOrder = totalOrder + v
								}

								if (totalOrder + v.Quantity) > v.Product.OrderMaxQty {
									o.Failure("qty"+strconv.Itoa(i)+".invalid", util.ErrorInvalidData("max qty order"))
								}
							}
						}

						if v.Product.MaxDayDeliveryDate > 0 && v.Product.MaxDayDeliveryDate < int64(r.AreaPolicy.MaxDayDeliveryDate) {
							maxDayDeliveryDate = v.Product.MaxDayDeliveryDate
						} else {
							maxDayDeliveryDate = int64(r.AreaPolicy.MaxDayDeliveryDate)
						}

						// get maximum available date based on max day delivery date for each product or max day delivery date on area policy
						for dayCount, maxAvailableDate = 0, r.RecognitionDate; ; maxAvailableDate = maxAvailableDate.AddDate(0, 0, 1) {
							if (int(maxAvailableDate.Weekday()) == 0 && r.AreaPolicy.WeeklyDayOff == 7) || (int(maxAvailableDate.Weekday()) == r.AreaPolicy.WeeklyDayOff) {
								continue
							}

							if err = o1.Raw("SELECT * FROM day_off WHERE off_date = ? LIMIT 1", maxAvailableDate.Format("2006-01-02")).QueryRow(&r.DayOff); err == nil && r.DayOff.ID != 0 {
								continue
							}

							dayCount++
							if int64(dayCount) > maxDayDeliveryDate {
								break
							}
						}

						maxAvailableDate = time.Date(maxAvailableDate.Year(), maxAvailableDate.Month(), maxAvailableDate.Day(), 0, 0, 0, 0, time.Local)
						if maxAvailableDate.Before(r.DeliveryDate) {
							o.Failure("product_id"+strconv.Itoa(i)+".invalid", "Max. delivery date H+"+strconv.Itoa(int(maxDayDeliveryDate))+" for this product")
						}

						v.Subtotal = v.Quantity * float64(v.UnitPrice)
						v.Weight = v.Quantity * v.Product.UnitWeight

						salesOrderItem = &model.SalesOrderItem{SalesOrder: r.SalesOrder, Product: v.Product}

						if e := salesOrderItem.Read("SalesOrder", "Product"); e == nil {
							v.TaxableItem = salesOrderItem.TaxableItem
							v.TaxPercentage = salesOrderItem.TaxPercentage
						}

						if orderTypeID != 5 {
							// start sku discount validation
							if v.SkuDiscountItem, err = repository.GetSkuDiscountData(r.SalesOrder.Branch.Merchant.ID, r.SalesOrder.Branch.PriceSet.ID, v.Product.ID, salesOrderItem.ID, r.SalesOrder.OrderChannel, r.CurrentTime); err == nil && v.SkuDiscountItem != nil {
								// set maximum available qty for discount
								maxRemQuota = float64(v.SkuDiscountItem.RemOverallQuota)
								if float64(v.SkuDiscountItem.RemQuotaPerUser) < maxRemQuota {
									maxRemQuota = float64(v.SkuDiscountItem.RemQuotaPerUser)
								}

								if float64(v.SkuDiscountItem.RemDailyQuotaPerUser) < maxRemQuota {
									maxRemQuota = float64(v.SkuDiscountItem.RemDailyQuotaPerUser)
								}

								if maxRemQuota > 0 && maxRemQuota != v.MaxDiscQty && maxRemQuota < v.DiscQty {
									o.Failure("rem_qty"+strconv.Itoa(i)+".invalid", util.ErrorHasChange("Max discount quota"))
								}

								if v.DiscQty > 0 && v.Quantity >= v.SkuDiscountItem.SkuDiscountItemTiers[0].MinimumQty {
									if v.SkuDiscountItem.RemOverallQuota <= 0 || (v.SkuDiscountItem.IsUseBudget == 1 && v.SkuDiscountItem.RemBudget <= 0) {
										o.Failure("rem_qty"+strconv.Itoa(i)+".invalid", util.ErrorRunOut("Discount quota for this product"))
									}

									if float64(v.SkuDiscountItem.RemDailyQuotaPerUser) < maxRemQuota {
										maxRemQuota = float64(v.SkuDiscountItem.RemDailyQuotaPerUser)
									}

									if maxRemQuota > 0 && maxRemQuota != v.MaxDiscQty && maxRemQuota < v.DiscQty {
										o.Failure("rem_qty"+strconv.Itoa(i)+".invalid", util.ErrorHasChange("Max discount quota"))
									}

									for _, val := range v.SkuDiscountItem.SkuDiscountItemTiers {
										if v.Quantity < val.MinimumQty {
											break
										} else {
											v.UnitPriceDiscount = val.DiscAmount
										}
									}

									v.DiscAmount = v.DiscQty * v.UnitPriceDiscount
									v.Subtotal = v.Subtotal - v.DiscAmount

									v.IsUseSkuDiscount = 1

									r.SkuDiscountItems = append(r.SkuDiscountItems, v.SkuDiscountItem)

									if err = v.SkuDiscountItem.SkuDiscount.Read("ID"); err != nil {
										o.Failure("product_id"+strconv.Itoa(i)+".invalid", util.ErrorInvalidData("product"))
									}
								}
							}
							// end sku discount validation

							v.SkuDiscountItemID = common.Encrypt(v.SkuDiscountItemID)
							if v.SkuDiscountItemID != "0" && (v.SkuDiscountItem == nil || (v.SkuDiscountItem != nil && strconv.Itoa(int(v.SkuDiscountItem.ID)) != v.SkuDiscountItemID)) {
								o.Failure("rem_qty"+strconv.Itoa(i)+".invalid", util.ErrorHasChange("Promo"))
								return o
							}
						}

						r.TotalSkuDiscAmount += v.DiscAmount
						r.TotalPrice = r.TotalPrice + v.Subtotal
					} else {
						o.Failure("product_id"+strconv.Itoa(i)+".invalid", util.ErrorInvalidData("product"))
					}

					//check push product
					if e := o1.Raw("SELECT soi.product_push FROM sales_order_item soi WHERE soi.sales_order_id = ? AND soi.product_id = ?", r.SalesOrder.ID, productID).QueryRow(&v.ProductPush); e != nil {
						v.ProductPush = 2
						if e = o1.Raw("SELECT * FROM product_push pp "+
							"WHERE pp.product_id = ? AND ? >= pp.start_date AND pp.area_id = ? AND pp.archetype_id = ? AND pp.status = 1",
							productID, r.CurrentTime.Format(layout), r.SalesOrder.Area.ID, r.SalesOrder.Archetype.ID).QueryRow(&checkProduct); e != nil {
							continue
						}

						if checkProduct != nil && r.CurrentTime.Unix() >= checkProduct.StartDate.Unix() {
							v.ProductPush = checkProduct.Status
						}
					}

				} else {
					o.Failure("product_id"+strconv.Itoa(i)+".invalid", util.ErrorInvalidData("product"))
				}

				// Validation: If voucher has items
				if r.Voucher != nil && r.Voucher.VoucherItem == 1 {

					o1.Raw("SELECT EXISTS(SELECT id FROM voucher_item where voucher_id = ? AND product_id = ? AND min_qty_disc <= ?) ", r.Voucher.ID, productID, v.Quantity).QueryRow(&countMatch)
					if countMatch > 0 {
						totalMatch++
					}
				}

			}
		}
	} else {
		r.TotalPrice = r.SalesOrder.TotalPrice
		r.TotalCharge = r.SalesOrder.TotalCharge
		r.TotalWeight = r.SalesOrder.TotalWeight
		r.DeliveryFee = r.SalesOrder.DeliveryFee
	}

	// Change Shipping address and delivery fee based on branch, if changed order type from self pickup
	if r.SalesOrder.OrderType.ID == 6 && orderTypeID != 6 {
		r.DeliveryFee = r.AreaBusinessPolicy.DeliveryFee
		admDivision := &model.AdmDivision{SubDistrictId: r.SalesOrder.Branch.SubDistrict.ID}
		if err = admDivision.Read("SubDistrictId"); err == nil {
			r.SalesOrder.ShippingAddress = r.SalesOrder.Branch.ShippingAddress + " " + admDivision.ConcateAddress
		}
	}

	// Validation for Order Type Self Pick up
	if orderTypeID == 6 {

		// Get Default Warehouse Self Pickup
		if warehouseSelfPickUp, err = repository.GetWarehouseSelfPickupByAreaID(r.SalesOrder.Branch.Area.ID); err != nil {
			o.Failure("area_id.invalid", util.ErrorAreaSelfPickUp(r.SalesOrder.Branch.Area.Name))
			return o
		}

		// Check availability warehouse to self pickup
		if warehouseSelfPickUp.ID != warehouseID {
			o.Failure("warehouse_id.invalid", util.ErrorSelfPickUp(warehouseSelfPickUp.Name))
		}
		// Custom Shipping Address based on warehouse self pickup
		r.SalesOrder.ShippingAddress = warehouseSelfPickUp.Name + " Edenfarm, " + warehouseSelfPickUp.StreetAddress

		// if order type is self pickup, delivery fee is 0
		r.DeliveryFee = 0

		// Get Default Wrt Self Pick Up
		wrtSelfPickup := &model.Wrt{Area: r.SalesOrder.Branch.Area, Type: 2, Status: 1}
		if err = wrtSelfPickup.Read("area_id", "type", "status"); err != nil {
			o.Failure("area_id.invalid", util.ErrorAreaSelfPickUp(r.SalesOrder.Branch.Area.Name))
		}

		// Check availability wrt to self pickup
		if wrtType != 2 {
			o.Failure("wrt_id.invalid", util.ErrorSelfPickUp(wrtSelfPickup.Name))
		}

		// only COD be used as payment term for self pick up
		if paymentTermID != 10 {
			o.Failure("term_payment_sls_id.invalid", util.ErrorSelfPickUp("COD"))
		}
	}

	// Validation for not allowed the wrt Self Pickup used in another order type
	if orderTypeID != 6 && wrtType == 2 {
		o.Failure("wrt_id.invalid", util.ErrorOnlyValidFor("wrt", "order type", "self pickup"))
	}

	if isPriceSetChanged {
		priceChanged = strings.TrimSuffix(priceChanged, ", ")
		r.NotePriceChange = fmt.Sprintf("Price Changed | %s", priceChanged)
	}

	o1.Raw("select poa.status from picking_order_assign poa where poa.sales_order_id = ?", r.SalesOrder.ID).QueryRow(&r.SalesOrder.StatusPickingOrderAssign)

	if r.SalesOrder.StatusPickingOrderAssign == 6 {
		o.Failure("sales_order_id.invalid", util.SalesOrderCannotBeUpdated())
	}

	if r.Voucher != nil && r.Voucher.VoucherItem == 1 {
		if r.UpdateAll == 1 && totalMatch != len(r.Voucher.VoucherItems) {
			o.Failure("redeem_code.invalid", util.ErrorNotValidTermConditions())
		}
	}

	if r.TotalPrice >= r.AreaBusinessPolicy.MinOrder {
		r.DeliveryFee = 0
	}

	r.TotalCharge = float64(r.TotalPrice) + r.DeliveryFee

	filter = map[string]interface{}{"term_payment_sls_id": paymentTermID, "term_invoice_sls_id": invoiceTermID, "payment_group_sls_id": paymentGroupID}
	if _, countPaymentGroup, err := repository.CheckPaymentGroupCombData(filter, exclude); err == nil && countPaymentGroup == 0 {
		o.Failure("payment_group_comb_id.invalid", util.ErrorPaymentCombination())
	}

	if paymentGroupID == 1 {
		o.Failure("payment_group_comb_id.invalid", util.ErrorAddDocument("sales order", "merchant payment group", "bayar langsung"))
	}

	if r.RedeemCode != "" {
		if r.Voucher.Status != 1 {
			o.Failure("redeem_code.inactive", util.ErrorActive("voucher"))
			return o
		}

		if r.CurrentTime.Before(r.Voucher.StartTimestamp) {
			o.Failure("redeem_code.invalid", util.ErrorNotInPeriod("voucher"))
			return o
		}

		if r.CurrentTime.After(r.Voucher.EndTimestamp) {
			o.Failure("redeem_code.invalid", util.ErrorOutOfPeriod("voucher"))
			return o
		}

		if r.Voucher.Type == 1 { //type total discount
			if r.Voucher.MinOrder > r.TotalPrice {
				o.Failure("redeem_code.greater", util.ErrorEqualGreater("total order", "minimum order"))
				return o
			}

			if r.Voucher.DiscAmount > r.TotalPrice {
				o.Failure("redeem_code.greater", util.ErrorEqualGreater("total order", "discount amount"))
				return o
			}
		} else if r.Voucher.Type == 2 { // type grand total discount
			if r.Voucher.MinOrder > r.TotalCharge {
				o.Failure("redeem_code.greater", util.ErrorEqualGreater("grand total order", "minimum order"))
				return o
			}

			if r.Voucher.DiscAmount > r.TotalCharge {
				o.Failure("redeem_code.greater", util.ErrorEqualGreater("grand total order", "discount amount"))
				return o
			}
		} else if r.Voucher.Type == 3 { // type delivery discount
			if r.Voucher.MinOrder > r.TotalPrice {
				o.Failure("redeem_code.greater", util.ErrorEqualGreater("total order", "minimum order"))
				return o
			}

			if r.Voucher.DiscAmount > r.DeliveryFee {
				o.Failure("redeem_code.greater", util.ErrorEqualGreater("delivery fee", "discount amount"))
				return o
			}
		}

		/* condition when voucher in sales order is nil or has id 0 meaning current sales order use no voucher
		or voucher in sales order not nil and voucher id not 0 and voucher id not the same as voucher id in sales order
		meaning current sales order use voucher and the new voucher is not the same as the voucher in current sales order
		*/
		if (r.SalesOrder.Voucher == nil || (r.SalesOrder.Voucher != nil && r.SalesOrder.Voucher.ID == 0)) ||
			(r.SalesOrder.Voucher != nil && r.SalesOrder.Voucher.ID != 0 && r.Voucher.ID != r.SalesOrder.Voucher.ID) {
			if r.Voucher.RemOverallQuota < 1 {
				o.Failure("redeem_code.invalid", util.ErrorFullyUsed("voucher"))
				return o
			}

			filter = map[string]interface{}{"merchant_id": r.SalesOrder.Branch.Merchant.ID, "voucher_id": r.Voucher.ID, "status": int8(1)}
			if _, countVoucherLog, err := repository.CheckVoucherLogData(filter, exclude); err == nil && countVoucherLog >= r.Voucher.UserQuota {
				o.Failure("redeem_code.invalid", util.ErrorFullyUsed("voucher"))
				return o
			}

			if r.Voucher.TagCustomer != "" {
				for _, v := range strings.Split(r.SalesOrder.Branch.Merchant.TagCustomer, ",") {
					if strings.Contains(r.Voucher.TagCustomer, v) {
						tagCustomerID, _ := strconv.Atoi(v)
						tagCustomer := &model.TagCustomer{ID: int64(tagCustomerID)}
						tagCustomer.Read("ID")

						r.SameTagCustomer = r.SameTagCustomer + "," + tagCustomer.Name
					}
				}

				r.SameTagCustomer = strings.Trim(r.SameTagCustomer, ",")
				if r.SameTagCustomer == "" {
					o.Failure("redeem_code.invalid", util.ErrorNotValidFor("voucher", "customer tag"))
					return o
				}
			}

			if r.Voucher.Area.ID != 1 && r.SalesOrder.Branch.Area.ID != r.Voucher.Area.ID {
				o.Failure("redeem_code.invalid", util.ErrorNotValidFor("voucher", "area"))
				return o
			}

			r.Voucher.Archetype.Read("ID")
			r.Voucher.Archetype.BusinessType.Read("ID")
			if r.Voucher.Archetype.BusinessType.AuxData != 1 {
				if r.Voucher.Archetype.AuxData != 1 {
					if r.SalesOrder.Branch.Archetype.ID != r.Voucher.Archetype.ID {
						o.Failure("redeem_code.invalid", util.ErrorNotValidFor("voucher", "archetype"))
						return o
					}
				} else {
					r.SalesOrder.Branch.Archetype.Read("ID")
					if r.SalesOrder.Branch.Archetype.BusinessType.ID != r.Voucher.Archetype.BusinessType.ID {
						o.Failure("redeem_code.invalid", util.ErrorNotValidFor("voucher", "business type"))
						return o
					}
				}
			}
		}

		// only has discount amount if type not extra eden point
		if r.Voucher.Type != 4 {
			discAmount = r.Voucher.DiscAmount
		}
	}

	if r.SalesOrder.PointRedeemAmount != 0 {
		discAmount = discAmount + r.SalesOrder.PointRedeemAmount
	}

	r.TotalCharge = r.TotalCharge - discAmount

	if r.TotalCharge < 0 {
		o.Failure("grand_total.invalid", util.ErrorEqualGreater("grand total", "0"))
	}

	r.PackingOrder = &model.PackingOrder{
		Status:       1,
		DeliveryDate: r.DeliveryDate,
		Warehouse:    r.Warehouse,
	}

	if err = r.PackingOrder.Read("Status", "DeliveryDate", "Warehouse"); err != nil {
		r.PackingOrder = nil
	}

	r.CreditLimitBefore = r.SalesOrder.Branch.Merchant.RemainingCreditLimitAmount
	totalChargeDifferences = r.TotalCharge - r.OldTotalCharge
	r.CreditLimitAfter = r.CreditLimitBefore - totalChargeDifferences

	if r.SalesOrder.Branch.Merchant.CreditLimitAmount > 0 || r.CreditLimitBefore < 0 {

		if r.TotalCharge < r.OldTotalCharge {
			totalChargeDifferences = r.OldTotalCharge - r.TotalCharge

			r.CreditLimitAfter = r.CreditLimitBefore + totalChargeDifferences
		}
		if r.CreditLimitAfter < 0 && r.CreditLimitBefore > 0 {
			o.Failure("credit_limit_amount.invalid", util.ErrorCreditLimitExceeded(r.SalesOrder.Branch.Merchant.Name))
		}
		r.IsCreateCreditLimitLog = 1
	}

	return o
}

// Messages : function to return error validation messages
func (r *updateRequest) Messages() map[string]string {
	messages := map[string]string{
		"delivery_date.required":       util.ErrorInputRequired("delivery date"),
		"wrt_id.required":              util.ErrorInputRequired("wrt"),
		"order_date.required":          util.ErrorInputRequired("order date"),
		"order_type_id.required":       util.ErrorInputRequired("order type"),
		"salesperson_id.required":      util.ErrorInputRequired("salesperson"),
		"term_payment_sls_id.required": util.ErrorInputRequired("payment term"),
		"term_invoice_sls_id.required": util.ErrorInputRequired("invoice term"),
		"payment_group_id.required":    util.ErrorInputRequired("payment group"),
		"warehouse_id.required":        util.ErrorInputRequired("warehouse"),
	}

	return messages
}
