package xendit_transaction

import (
	"fmt"
	"time"

	"git.edenfarm.id/cuxs/orm"
	"git.edenfarm.id/cuxs/validation"
	"git.edenfarm.id/project-version2/api/datastore/repository"
	"git.edenfarm.id/project-version2/datamodel/model"
)

type invoiceExpiredRequest struct {
	XenditID                string                         `json:"xendit_invoice_id"`
	TransactionDate         string                         `json:"transaction_date"`
	TransactionTime         string                         `json:"transaction_time"`
	Amount                  float64                        `json:"amount"`
	Token                   string                         `json:"token"`
	StatusInvoice           string                         `json:"status_invoice"`
	TransactionDateAt       time.Time                      `json:"-"`
	TransactionTimeAt       time.Time                      `json:"-"`
	ExternalID              string                         `json:"external_id" valid:"required"`
	RecentPoint             float64                        `json:"-"`
	CreditLimitBefore       float64                        `json:"-"`
	CreditLimitAfter        float64                        `json:"-"`
	IsCreateCreditLimitLog  bool                           `json:"-"`
	Bank                    *model.PaymentChannel          `json:"-"`
	SalesInvoiceExternal    *model.SalesInvoiceExternal    `json:"-"`
	MerchantPointLog        []*model.MerchantPointLog      `json:"-"`
	MerchantPointExpiration *model.MerchantPointExpiration `json:"-"`
}

func (c *invoiceExpiredRequest) Validate() *validation.Output {
	o := &validation.Output{Valid: true}
	var e error
	var filter, exclude map[string]interface{}

	key := []byte("joyfuls joy j0y5")
	dec := decrypt(key, c.Token)
	token := fmt.Sprintf("%s", dec)
	if token != "hey please push on thursday, it will make me happy" {
		o.Failure("token", "invalid token")
	}

	if c.TransactionDate != "" {
		if c.TransactionDateAt, e = time.Parse("2006-01-02", c.TransactionDate); e != nil {
			o.Failure("transaction_date.invalid", "invalid date")
		}
	}
	c.SalesInvoiceExternal = &model.SalesInvoiceExternal{XenditInvoiceID: c.XenditID}

	if e = c.SalesInvoiceExternal.Read("XenditInvoiceID"); e != nil {
		o.Failure("invoice_xendit_id.invalid", "no data")
	} else {
		c.SalesInvoiceExternal.SalesOrder.Read("ID")
		c.SalesInvoiceExternal.SalesOrder.Branch.Read("ID")
		c.SalesInvoiceExternal.SalesOrder.Branch.Merchant.Read("ID")
		c.SalesInvoiceExternal.SalesOrder.Branch.Merchant.UserMerchant.Read("ID")
		if c.SalesInvoiceExternal.SalesOrder.PointRedeemID != 0 && c.SalesInvoiceExternal.SalesOrder.PointRedeemAmount != 0 {
			filter = map[string]interface{}{
				"id":             c.SalesInvoiceExternal.SalesOrder.PointRedeemID,
				"merchant_id":    c.SalesInvoiceExternal.SalesOrder.Branch.Merchant.ID,
				"sales_order_id": c.SalesInvoiceExternal.SalesOrder.ID,
				"status":         int8(2),
				"point_value":    c.SalesInvoiceExternal.SalesOrder.PointRedeemAmount,
			}
			c.MerchantPointLog, _, e = repository.CheckMerchantPointLogData(filter, exclude)

			c.MerchantPointExpiration = &model.MerchantPointExpiration{ID: c.SalesInvoiceExternal.SalesOrder.Branch.Merchant.ID}
			c.MerchantPointExpiration.Read("ID")
		}

		or := orm.NewOrm()
		or.Using("read_only")
		or.Raw("SELECT recent_point from merchant_point_log where merchant_id = ? order by id desc limit 1 ", c.SalesInvoiceExternal.SalesOrder.Branch.Merchant.ID).QueryRow(&c.RecentPoint)

		if c.SalesInvoiceExternal.SalesOrder.Branch.Merchant.CreditLimitAmount > 0 {
			c.CreditLimitBefore = c.SalesInvoiceExternal.SalesOrder.Branch.Merchant.RemainingCreditLimitAmount
			c.CreditLimitAfter = c.CreditLimitBefore + c.SalesInvoiceExternal.SalesOrder.TotalCharge
			c.IsCreateCreditLimitLog = true
		}
	}

	return o
}

func (c *invoiceExpiredRequest) Messages() map[string]string {
	return map[string]string{}
}
