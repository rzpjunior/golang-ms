package farm_stack

import (
	"git.edenfarm.id/cuxs/validation"
	"git.edenfarm.id/project-version2/api/util"
	"time"
)

type requestGet struct {
	PurchaseDate string `json:"purchase_date" valid:"required"`
	Data         []dataCallback
	EtaDate      time.Time
}

type dataCallback struct {
	FarmerName   string  `json:"farmer_name"`
	Latitude     string  `json:"latitude"`
	Longitude    string  `json:"longitude"`
	ProductName  string  `json:"product_name"`
	ProductPrice float64 `json:"product_price"`
	PurchaseDate string  `json:"purchase_date"`
	PurchaseTime string  `json:"purchase_time"`
}

func (c *requestGet) Validate() *validation.Output {
	o := &validation.Output{Valid: true}
	var e error
	layout := "2006-01-02"
	if c.EtaDate, e = time.Parse(layout, c.PurchaseDate); e != nil {
		o.Failure("purchase_date.invalid", util.ErrorInvalidData("purchase date"))
	}

	return o
}

func (c *requestGet) Messages() map[string]string {
	return map[string]string{}
}
