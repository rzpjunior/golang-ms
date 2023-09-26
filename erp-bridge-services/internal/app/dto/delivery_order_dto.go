package dto

import "time"

type DeliveryOrderResponse struct {
	ID              int64     `json:"id,omitempty"`
	Code            string    `json:"code,omitempty"`
	CustomerID      int64     `json:"customer_id,omitempty"`
	WrtID           int64     `json:"wrt_id,omitempty"`
	SiteID          int64     `json:"site_id,omitempty"`
	Status          int8      `json:"status,omitempty"`
	StatusConvert   string    `json:"status_convert,omitempty"`
	RecognitionDate time.Time `json:"recognition_date,omitempty"`
	CreatedDate     time.Time `json:"created_date,omitempty"`
}

type CreateDeliveryOrderRequest struct {
	Interid      string                      `json:"interid"`
	Docnumber    string                      `json:"docnumbr"`
	Docdate      string                      `json:"docdate"`
	Custnumber   string                      `json:"custnmbr"`
	Custname     string                      `json:"custname"`
	GnlCourierId string                      `json:"gnL_courier_id"`
	DetailOrder  []*DeliveryOrderDetailOrder `json:"detail_order"`
	DetailItem   []*DeliveryOrderDetailItem  `json:"detail_item"`
}

type DeliveryOrderDetailOrder struct {
	IvmCb      int32   `json:"ivm_cb"`
	SopNumber  string  `json:"sopnumbe"`
	QtyOrder   float64 `json:"qtyorder"`
	IvmQtyPack float64 `json:"ivm_qty_pack"`
}

type DeliveryOrderDetailItem struct {
	Lnseqnbr   int32   `json:"lnseqnbr"`
	ItemNumber string  `json:"itemnmbr"`
	Uofm       string  `json:"uofm"`
	QtyOrder   float64 `json:"qtyorder"`
	IvmQtyPack float64 `json:"ivm_qty_pack"`
	Locncode   string  `json:"locncode"`
}
