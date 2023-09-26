// Copyright 2020 PT. Eden Pangan Indonesia. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

package order_edn

import (
	"fmt"
	"strconv"
	"strings"
	"time"

	"git.edenfarm.id/cuxs/common"
	"git.edenfarm.id/cuxs/orm"
	"git.edenfarm.id/cuxs/validation"
	"git.edenfarm.id/project-version2/datamodel/model"

	re "regexp"

	"git.edenfarm.id/project-version2/api/datastore/repository"
	"git.edenfarm.id/project-version2/api/src/auth"
	"git.edenfarm.id/project-version2/api/util"
)

// createRequest : struct to hold sales order request data
type createRequest struct {
	Code                   string          `json:"-"`
	MerchantID             string          `json:"merchant_id" valid:"required"`
	BranchID               string          `json:"branch_id" valid:"required"`
	AreaID                 string          `json:"area_id" valid:"required"`
	DeliveryDateStr        string          `json:"delivery_date" valid:"required"`
	WrtID                  string          `json:"wrt_id" valid:"required"`
	WarehouseID            string          `json:"warehouse_id" valid:"required"`
	RecognitionDateStr     string          `json:"order_date" valid:"required"`
	OrderTypeID            string          `json:"order_type_id" valid:"required"`
	SalespersonID          string          `json:"salesperson_id" valid:"required"`
	SalesTermID            string          `json:"term_payment_sls_id" valid:"required"`
	InvoiceTermID          string          `json:"term_invoice_sls_id" valid:"required"`
	PaymentGroupID         string          `json:"payment_group_id" valid:"required"`
	ShippingAddress        string          `json:"shipping_address" valid:"required"`
	BillingAddress         string          `json:"billing_address" valid:"required"`
	RedeemCode             string          `json:"redeem_code"`
	Note                   string          `json:"note"`
	TotalPrice             float64         `json:"-"`
	TotalWeight            float64         `json:"-"`
	TotalCharge            float64         `json:"-"`
	SameTagCustomer        string          `json:"-"`
	DeliveryDate           time.Time       `json:"-"`
	RecognitionDate        time.Time       `json:"-"`
	IsCreateMerchantVa     map[string]int8 `json:"-"`
	CurrentTime            time.Time       `json:"-"`
	CreditLimitBefore      float64         `json:"-"`
	CreditLimitAfter       float64         `json:"-"`
	IsCreateCreditLimitLog bool            `json:"-"`
	NotePriceChange        string          `json:"-"`

	Branch       *model.Branch       `json:"-"`
	Wrt          *model.Wrt          `json:"-"`
	Warehouse    *model.Warehouse    `json:"-"`
	OrderType    *model.OrderType    `json:"-"`
	Salesperson  *model.Staff        `json:"-"`
	SalesTerm    *model.SalesTerm    `json:"-"`
	InvoiceTerm  *model.InvoiceTerm  `json:"-"`
	PaymentGroup *model.PaymentGroup `json:"-"`
	Voucher      *model.Voucher      `json:"-"`
	AreaPolicy   *model.AreaPolicy   `json:"-"`
	PackingOrder *model.PackingOrder `json:"-"`
	DayOff       *model.DayOff

	Products []*salesOrderItem `json:"products" valid:"required"`

	Session *auth.SessionData `json:"-"`
}

type salesOrderItem struct {
	ProductID     string  `json:"product_id"`
	Quantity      float64 `json:"qty"`
	UnitPrice     int64   `json:"unit_price"`
	Note          string  `json:"note"`
	ProductPush   int8    `json:"product_push"`
	TaxPercentage float64 `json:"-"`
	TaxableItem   int8    `json:"-"`
	Subtotal      float64 `json:"-"`
	Weight        float64 `json:"-"`
	DefaultPrice  float64 `json:"-"`

	Product *model.Product `json:"-"`
	Price   *model.Price   `json:"-"`
}

// Validate : function to validate sales order request data
func (r *createRequest) Validate() *validation.Output {
	o := &validation.Output{Valid: true}
	q := orm.NewOrm()
	q.Using("read_only")
	var (
		err                                                                                                                                     error
		filter, exclude                                                                                                                         map[string]interface{}
		branchID, merchantID, warehouseID, paymentTermID, invoiceTermID, paymentGroupID, wrtID, orderTypeID, salespersonID, productID, totalAcc int64
		totalMatch, countMatch                                                                                                                  int
		checkProduct                                                                                                                            *model.ProductPush
		priceChanged                                                                                                                            string
		isPriceSetChanged                                                                                                                       bool
	)
	productList := make(map[int64]string)
	discAmount := float64(0)
	r.CurrentTime = time.Now()

	r.IsCreateMerchantVa = map[string]int8{"bca": 0, "permata": 0}

	if branchID, err = common.Decrypt(r.BranchID); err != nil {
		o.Failure("branch_id.invalid", util.ErrorInvalidData("branch"))
		return o
	}

	if r.Branch, err = repository.ValidBranch(branchID); err != nil {
		o.Failure("branch_id.invalid", util.ErrorInvalidData("branch"))
		return o
	}

	if r.Branch.Status != 1 {
		o.Failure("branch_id.inactive", util.ErrorActive("branch"))
		return o
	}

	if merchantID, err = common.Decrypt(r.MerchantID); err != nil {
		o.Failure("merchant_id.invalid", util.ErrorInvalidData("merchant"))
		return o
	}

	if r.Branch.Merchant, err = repository.ValidMerchant(merchantID); err != nil {
		o.Failure("merchant_id.invalid", util.ErrorInvalidData("merchant"))
		return o
	}

	if merchantID != r.Branch.Merchant.ID {
		o.Failure("branch_id.invalid", util.ErrorInvalidData("branch"))
		return o
	}

	if r.Branch.Merchant.Status != 1 {
		o.Failure("merchant_id.inactive", util.ErrorActive("merchant"))
		return o
	}

	if r.Branch.Merchant.Suspended == 1 {
		o.Failure("customer.suspended", util.ErrorSuspended("customer"))
		return o
	}

	if r.Branch.SubDistrict, err = repository.ValidSubDistrict(r.Branch.SubDistrict.ID); err != nil {
		o.Failure("subdistrict_id.invalid", util.ErrorInvalidData("subdistrict"))
	}

	if r.Branch.Archetype, err = repository.ValidArchetype(r.Branch.Archetype.ID); err != nil {
		o.Failure("archetype_id.invalid", util.ErrorInvalidData("archetype"))
	}

	if r.Branch.PriceSet, err = repository.ValidPriceSet(r.Branch.PriceSet.ID); err != nil {
		o.Failure("priceset_id.invalid", util.ErrorInvalidData("priceset"))
	}

	if r.Branch.Area, err = repository.ValidArea(r.Branch.Area.ID); err == nil {
		if r.AreaPolicy, err = repository.GetAreaPolicy("area_id", r.Branch.Area.ID); err != nil {
			o.Failure("area_id.invalid", util.ErrorInvalidData("area"))
		}
	} else {
		o.Failure("area_id.invalid", util.ErrorInvalidData("area"))
	}

	// validation staff in session
	if err = r.Session.Staff.Role.Read("ID"); err != nil {
		o.Failure("staff.invalid", util.ErrorInvalidData("staff"))
	}
	if r.WarehouseID != "" {
		if warehouseID, err = common.Decrypt(r.WarehouseID); err == nil {
			if r.Warehouse, err = repository.ValidWarehouse(warehouseID); err != nil {
				o.Failure("warehouse_id.invalid", util.ErrorInvalidData("warehouse"))
			}
		} else {
			o.Failure("warehouse_id.invalid", util.ErrorInvalidData("warehouse"))
		}
	} else {
		o.Failure("warehouse_id.required", util.ErrorInputRequired("warehouse"))
	}

	if r.SalesTermID != "" {
		if paymentTermID, err = common.Decrypt(r.SalesTermID); err == nil {
			if r.SalesTerm, err = repository.ValidSalesTerm(paymentTermID); err != nil {
				o.Failure("sales_term_id.invalid", util.ErrorInvalidData("sales term"))
			}
		} else {
			o.Failure("sales_term_id.invalid", util.ErrorInvalidData("sales term"))
		}
	} else {
		o.Failure("sales_term_id.required", util.ErrorInputRequired("sales term"))
	}

	if r.InvoiceTermID != "" {
		if invoiceTermID, err = common.Decrypt(r.InvoiceTermID); err == nil {
			if r.InvoiceTerm, err = repository.ValidInvoiceTerm(invoiceTermID); err != nil {
				o.Failure("invoice_term_id.invalid", util.ErrorInvalidData("invoice term"))
			}
		} else {
			o.Failure("invoice_term_id.invalid", util.ErrorInvalidData("invoice term"))
		}
	} else {
		o.Failure("invoice_term_id.required", util.ErrorInputRequired("invoice term"))
	}

	if r.PaymentGroupID != "" {
		if paymentGroupID, err = common.Decrypt(r.PaymentGroupID); err == nil {
			if r.PaymentGroup, err = repository.ValidPaymentGroup(paymentGroupID); err != nil {
				o.Failure("payment_group_id.invalid", util.ErrorInvalidData("payment group"))
			}
		} else {
			o.Failure("payment_group_id.invalid", util.ErrorInvalidData("payment group"))
		}
	} else {
		o.Failure("payment_group_id.required", util.ErrorInputRequired("payment group"))
	}

	if r.OrderTypeID != "" {
		if orderTypeID, err = common.Decrypt(r.OrderTypeID); err != nil {
			o.Failure("order_type_id.invalid", util.ErrorInvalidData("order type"))
		}
		if orderTypeID != 13 {
			o.Failure("order_type_id.invalid", util.ErrorDocStatus("order_type", "EDN sales"))
		}
		if r.OrderType, err = repository.ValidOrderType(orderTypeID); err != nil {
			o.Failure("order_type_id.invalid", util.ErrorInvalidData("order type"))
		}
	} else {
		o.Failure("order_type_id.required", util.ErrorInputRequired("order type"))
	}

	if r.SalespersonID != "" {
		if salespersonID, err = common.Decrypt(r.SalespersonID); err == nil {
			if r.Salesperson, err = repository.ValidStaff(salespersonID); err == nil {
				if r.Salesperson.Status != 1 {
					o.Failure("salesperson_id.inactive", util.ErrorActive("salesperson"))
				}
			} else {
				o.Failure("salesperson_id.invalid", util.ErrorInvalidData("salesperson"))
			}
		} else {
			o.Failure("salesperson_id.invalid", util.ErrorInvalidData("salesperson"))
		}
	} else {
		o.Failure("salesperson_id.required", util.ErrorInputRequired("salesperson"))
	}

	if wrtID, err = common.Decrypt(r.WrtID); err == nil {
		if r.Wrt, err = repository.ValidWrt(wrtID); err != nil {
			o.Failure("wrt_id.invalid", util.ErrorInvalidData("wrt"))
		}
	} else {
		o.Failure("wrt_id.invalid", util.ErrorInvalidData("wrt"))
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
	} else {
		r.Voucher = nil
	}

	// Validation: voucher suits with the following merchant
	if r.Voucher != nil && r.Voucher.MerchantID != 0 {
		if r.Voucher.MerchantID != r.Branch.Merchant.ID {
			o.Failure("redeem_code.invalid", util.ErrorNotValidFor("voucher", "merchant"))
		}
	}

	for i, v := range r.Products {
		var (
			listOrder          []float64
			totalOrder         float64
			maxDayDeliveryDate int64
			maxAvailableDate   time.Time
			dayCount           int8
		)

		// if product_qty below the sales_order_type.min_order_qty
		if v.Quantity < r.OrderType.MinOrderQty {
			o.Failure("qty"+strconv.Itoa(i)+".equalorgreater", util.ErrorInvalidData("order qty"))
		}

		// if product_price below the sales_order_type.min_unit_price
		if float64(v.UnitPrice) < r.OrderType.MinUnitPrice {
			o.Failure("unit_price"+strconv.Itoa(i)+".equalorgreater", util.ErrorInvalidData("unit price"))
		}

		if v.ProductID == "" {
			o.Failure("product_id"+strconv.Itoa(i)+".invalid", util.ErrorInputRequired("product"))
			continue
		}

		if productID, err = common.Decrypt(v.ProductID); err != nil {
			o.Failure("product_id"+strconv.Itoa(i)+".invalid", util.ErrorInvalidData("product"))
			continue
		}

		if v.Product, err = repository.ValidProduct(productID); err != nil {
			o.Failure("product_id"+strconv.Itoa(i)+".invalid", util.ErrorInvalidData("product"))
			continue
		}

		if _, exist := productList[productID]; exist {
			o.Failure("product_id"+strconv.Itoa(i)+".duplicate", util.ErrorDuplicate("product"))
			return o
		}

		productList[productID] = "t"

		//Check UOM Decimal
		if err = v.Product.Uom.Read("ID"); err == nil {
			if v.Product.Uom.DecimalEnabled == 2 {
				if v.Quantity != float64((int64(v.Quantity))) {
					o.Failure("qty"+strconv.Itoa(i)+".invalid", util.ErrorInvalidData("product quantity"))
				}
			}
		}

		if v.Product.OrderChannelRestriction != "" {
			matched, err := re.MatchString("1", v.Product.OrderChannelRestriction)
			if err != nil {
				o.Failure("product_id"+strconv.Itoa(i)+".order_channel_restriction", util.ErrorInvalidData("order channel restriction"))
			}
			if matched {
				o.Failure("product_id"+strconv.Itoa(i)+".order_channel_restriction", util.ErrorOrderChannelRestriction())
			}
		}

		v.Price = &model.Price{Product: v.Product, PriceSet: r.Branch.PriceSet}
		if err = v.Price.Read("Product", "PriceSet"); err != nil {
			o.Failure("price.invalid", util.ErrorInvalidData("price"))
		}

		v.DefaultPrice = v.Price.UnitPrice

		if v.UnitPrice != int64(v.DefaultPrice) {
			priceChanged += fmt.Sprintf("%s: Default Unit Price: Rp %g - Changed Unit Price: Rp %d, ", v.Product.Name, v.DefaultPrice, v.UnitPrice)
			isPriceSetChanged = true
		}

		r.TotalWeight = r.TotalWeight + (v.Quantity * v.Product.UnitWeight)

		// check salable for stock
		filter = map[string]interface{}{"product_id": productID, "warehouse_id": r.Warehouse.ID, "salable": 1}

		// Check stock data
		if _, countStock, err := repository.CheckStockData(filter, exclude); err == nil && countStock == 0 {
			o.Failure("product_id"+strconv.Itoa(i)+".invalid", util.ErrorProductMustAvailable())
		}

		//Check daily limit
		if v.Product.OrderMaxQty != 0 {
			if v.Quantity > v.Product.OrderMaxQty {
				o.Failure("qty"+strconv.Itoa(i)+".invalid", util.ErrorInvalidData("max qty order"))
			}

			//Check product from another so (based on branch & delivery_date) still under maximum order qty
			q.Raw("SELECT soi.order_qty "+
				"FROM sales_order_item soi "+
				"JOIN sales_order so ON so.id = soi.sales_order_id "+
				"JOIN branch b ON b.id = so.branch_id "+
				"JOIN merchant m ON m.id = b.merchant_id "+
				"WHERE m.id = ? AND soi.product_id = ? AND so.delivery_date = ? AND so.status NOT IN (3,4)", r.Branch.Merchant.ID, v.Product.ID, r.DeliveryDate).QueryRows(&listOrder)

			for _, v := range listOrder {
				totalOrder = totalOrder + v
			}

			if (totalOrder + v.Quantity) > v.Product.OrderMaxQty {
				o.Failure("qty"+strconv.Itoa(i)+".invalid", util.ErrorInvalidData("max qty order"))
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

			if err = q.Raw("SELECT * FROM day_off WHERE status = 1 AND off_date = ? LIMIT 1", maxAvailableDate.Format("2006-01-02")).QueryRow(&r.DayOff); err == nil && r.DayOff.ID != 0 {
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

		r.TotalPrice = r.TotalPrice + v.Subtotal

		//check push product
		v.ProductPush = 2
		if e := q.Raw("SELECT * FROM product_push pp "+
			"WHERE pp.product_id = ? AND ? >= pp.start_date AND pp.area_id = ? AND pp.archetype_id = ? AND pp.status = 1",
			productID, r.CurrentTime.Format(layout), r.Branch.Area.ID, r.Branch.Archetype.ID).QueryRow(&checkProduct); e != nil {
			continue
		}

		if checkProduct != nil && r.CurrentTime.Unix() >= checkProduct.StartDate.Unix() {
			v.ProductPush = checkProduct.Status
		}

		// Validation: If voucher has items
		if r.Voucher != nil && r.Voucher.VoucherItem == 1 {
			q.Raw("SELECT EXISTS(SELECT id FROM voucher_item where voucher_id = ? AND product_id = ? AND min_qty_disc <= ?) ", r.Voucher.ID, productID, v.Quantity).QueryRow(&countMatch)
			if countMatch > 0 {
				totalMatch++
			}
		}
	}

	if isPriceSetChanged {
		priceChanged = strings.TrimSuffix(priceChanged, ", ")
		r.NotePriceChange = fmt.Sprintf("Price Changed | %s", priceChanged)
	}

	if r.Voucher != nil && r.Voucher.VoucherItem == 1 {
		q.LoadRelated(r.Voucher, "VoucherItems", 0)

		if totalMatch != len(r.Voucher.VoucherItems) {
			o.Failure("redeem_code.invalid", util.ErrorNotValidTermConditions())
		}
	}

	r.TotalCharge = float64(r.TotalPrice)

	filter = map[string]interface{}{"term_payment_sls_id": paymentTermID, "term_invoice_sls_id": invoiceTermID, "payment_group_sls_id": paymentGroupID}
	if _, countPaymentGroup, err := repository.CheckPaymentGroupCombData(filter, exclude); err == nil && countPaymentGroup == 0 {
		o.Failure("payment_group_id.invalid", util.ErrorPaymentCombination())
	}

	if paymentGroupID == 1 {
		o.Failure("payment_group_id.invalid", util.ErrorAddDocument("sales order", "merchant payment group", "bayar langsung"))
	}

	if r.TotalCharge < 0 {
		o.Failure("grand_total.invalid", util.ErrorEqualGreater("grand total", "0"))
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
		} else { // Cannot use voucher type delivery discount
			o.Failure("redeem_code.greater", util.ErrorEqualGreater("delivery fee", "discount amount"))
			return o
		}

		if r.Voucher.RemOverallQuota < 1 {
			o.Failure("redeem_code.invalid", util.ErrorFullyUsed("voucher"))
			return o
		}

		filter = map[string]interface{}{"merchant_id": r.Branch.Merchant.ID, "voucher_id": r.Voucher.ID, "status": 1}
		if _, countVoucherLog, err := repository.CheckVoucherLogData(filter, exclude); err == nil && countVoucherLog >= r.Voucher.UserQuota {
			o.Failure("redeem_code.invalid", util.ErrorFullyUsed("voucher"))
			return o
		}

		if r.Voucher.TagCustomer != "" {
			for _, v := range strings.Split(r.Branch.Merchant.TagCustomer, ",") {
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

		if r.Voucher.Area.ID != 1 && r.Branch.Area.ID != r.Voucher.Area.ID {
			o.Failure("redeem_code.invalid", util.ErrorNotValidFor("voucher", "area"))
			return o
		}

		r.Voucher.Archetype.Read("ID")
		r.Voucher.Archetype.BusinessType.Read("ID")
		if r.Voucher.Archetype.BusinessType.AuxData != 1 {
			if r.Voucher.Archetype.AuxData != 1 {
				if r.Branch.Archetype.ID != r.Voucher.Archetype.ID {
					o.Failure("redeem_code.invalid", util.ErrorNotValidFor("voucher", "archetype"))
					return o
				}
			} else {
				r.Branch.Archetype.Read("ID")
				if r.Branch.Archetype.BusinessType.ID != r.Voucher.Archetype.BusinessType.ID {
					o.Failure("redeem_code.invalid", util.ErrorNotValidFor("voucher", "business type"))
					return o
				}
			}
		}

		discAmount = r.Voucher.DiscAmount

	}

	if r.Branch != nil {
		exclude = map[string]interface{}{"account_number": " "}
		filter = map[string]interface{}{"merchant_id": r.Branch.Merchant.ID, "payment_channel_id": 6, "account_number__isnull": ""}
		_, totalAcc, err = repository.CheckMerchantAccNumData(filter, exclude)
		if totalAcc == 0 {
			r.IsCreateMerchantVa["bca"] = 1
		}

		filter = map[string]interface{}{"merchant_id": r.Branch.Merchant.ID, "payment_channel_id": 7, "account_number__isnull": ""}
		_, totalAcc, err = repository.CheckMerchantAccNumData(filter, exclude)
		if totalAcc == 0 {
			r.IsCreateMerchantVa["permata"] = 1
		}
	}

	r.TotalCharge = r.TotalCharge - discAmount

	r.PackingOrder = &model.PackingOrder{
		Status:       1,
		DeliveryDate: r.DeliveryDate,
		Warehouse:    r.Warehouse,
	}

	r.CreditLimitBefore = r.Branch.Merchant.RemainingCreditLimitAmount
	if r.Branch.Merchant.CreditLimitAmount == 0 {
		if r.CreditLimitBefore < 0 {
			o.Failure("credit_limit_amount.invalid", util.ErrorCreditLimitExceeded(r.Branch.Merchant.Name))
		}
	} else {
		r.CreditLimitAfter = r.CreditLimitBefore - r.TotalCharge
		if r.CreditLimitAfter < 0 {
			o.Failure("credit_limit_amount.invalid", util.ErrorCreditLimitExceeded(r.Branch.Merchant.Name))
		}
		r.IsCreateCreditLimitLog = true
	}

	return o
}

// Messages : function to return error validation messages
func (r *createRequest) Messages() map[string]string {
	messages := map[string]string{
		"merchant_id.required":         util.ErrorInputRequired("merchant"),
		"branch_id.required":           util.ErrorInputRequired("branch"),
		"area_id.required":             util.ErrorInputRequired("area"),
		"delivery_date.required":       util.ErrorInputRequired("delivery date"),
		"wrt_id.required":              util.ErrorInputRequired("wrt"),
		"warehouse_id.required":        util.ErrorInputRequired("warehouse"),
		"order_date.required":          util.ErrorInputRequired("order date"),
		"order_type_id.required":       util.ErrorInputRequired("order type"),
		"salesperson_id.required":      util.ErrorInputRequired("salesperson"),
		"term_payment_sls_id.required": util.ErrorInputRequired("payment term"),
		"term_invoice_sls_id.required": util.ErrorInputRequired("term invoice"),
		"payment_group_id.required":    util.ErrorInputRequired("payment group"),
		"shipping_address.required":    util.ErrorInputRequired("shipping address"),
		"billing_address.required":     util.ErrorInputRequired("billing address"),
		"products.required":            util.ErrorInputRequired("products"),
	}

	return messages
}
