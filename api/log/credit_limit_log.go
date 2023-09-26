package log

import (
	"time"

	"git.edenfarm.id/project-version2/datamodel/model"
	"github.com/labstack/echo/v4"
)

// this func for make credit limit log whenever there's a operation that cause change on the credit_limit_amount on the merchant ,
// merchant = related merchant that will have an update on their credit_limit_amount
// refID = get from id like sales_order id (when you create sales order)
// types = get from name of module, currently it has only 2 options ( sales_order & sales_invoice )
// note = get from creation or update note, this note is optional.

func CreditLimitLogByMerchant(merchant *model.Merchant, refID int64, types string, creditLimitBefore float64, creditLimitAfter float64, note ...string) (e error) {

	cll := &model.CreditLimitLog{
		Merchant:          merchant.ID,
		RefID:             refID,
		CreditLimitBefore: creditLimitBefore,
		CreditLimitAfter:  creditLimitAfter,
		Type:              types,
		Note:              note[0],
	}

	if note != nil {
		cll.Note = note[0]
	}

	if e = cll.Save(); e != nil {
		e = echo.ErrBadRequest
	}

	return e
}

// CreditLimitLogByStaff : function to insert Credit Limit Log with staff id
func CreditLimitLogByStaff(merchant *model.Merchant, refID int64, types string, creditLimitBefore float64, creditLimitAfter float64, staffID int64, note ...string) (e error) {

	cll := &model.CreditLimitLog{
		Merchant:          merchant.ID,
		RefID:             refID,
		CreditLimitBefore: creditLimitBefore,
		CreditLimitAfter:  creditLimitAfter,
		Type:              types,
		Note:              note[0],
		CreatedAt:         time.Now(),
		CreatedBy:         staffID,
	}

	if note != nil {
		cll.Note = note[0]
	}

	if e = cll.Save(); e != nil {
		e = echo.ErrBadRequest
	}

	return e
}
