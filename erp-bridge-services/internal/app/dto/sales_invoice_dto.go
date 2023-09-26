package dto

import (
	"time"
)

type SalesInvoiceResponse struct {
	ID                int64     `json:"-"`
	Code              string    `json:"code"`
	CodeExt           string    `json:"code_ext"`
	Status            int8      `json:"status"`
	RecognitionDate   time.Time `json:"recognition_date"`
	DueDate           time.Time `json:"due_date"`
	BillingAddress    string    `json:"billing_address"`
	DeliveryFee       float64   `json:"delivery_fee"`
	VouRedeemCode     string    `json:"vou_redeem_code"`
	VouDiscAmount     float64   `json:"vou_disc_amount"`
	PointRedeemAmount float64   `json:"point_redeem_amount"`
	Adjustment        int8      `json:"adjustment"`
	AdjAmount         float64   `json:"adj_amount"`
	AdjNote           string    `json:"adj_note"`
	TotalPrice        float64   `json:"total_price"`
	TotalCharge       float64   `json:"total_charge"`
	DeltaPrint        int64     `json:"delta_print"`
	Note              string    `json:"note"`
	VoucherID         int64     `json:"voucher_id"`
	RemainingAmount   float64   `json:"remaining_amount"`
}

type CreateSalesInvoiceGPRequest struct {
	Interid            string                                      `protobuf:"bytes,1,opt,name=interid,proto3" json:"interid"`
	Orignumb           string                                      `protobuf:"bytes,2,opt,name=orignumb,proto3" json:"orignumb"`
	Sopnumbe           string                                      `protobuf:"bytes,3,opt,name=sopnumbe,proto3" json:"sopnumbe"`
	Docid              string                                      `protobuf:"bytes,4,opt,name=docid,proto3" json:"docid"`
	Docdate            string                                      `protobuf:"bytes,5,opt,name=docdate,proto3" json:"docdate"`
	Custnmbr           string                                      `protobuf:"bytes,6,opt,name=custnmbr,proto3" json:"custnmbr"`
	Custname           string                                      `protobuf:"bytes,7,opt,name=custname,proto3" json:"custname"`
	Prstadcd           string                                      `protobuf:"bytes,8,opt,name=prstadcd,proto3" json:"prstadcd"`
	Curncyid           string                                      `protobuf:"bytes,9,opt,name=curncyid,proto3" json:"curncyid"`
	VoucherApply       []*CreateSalesInvoiceGPRequest_VoucherApply `protobuf:"bytes,10,rep,name=voucher_apply,json=voucherApply,proto3" json:"voucher_apply"`
	Subtotal           float64                                     `protobuf:"fixed64,11,opt,name=subtotal,proto3" json:"subtotal"`
	Trdisamt           float64                                     `protobuf:"fixed64,12,opt,name=trdisamt,proto3" json:"trdisamt"`
	Freight            float64                                     `protobuf:"fixed64,13,opt,name=freight,proto3" json:"freight"`
	Miscamnt           float64                                     `protobuf:"fixed64,14,opt,name=miscamnt,proto3" json:"miscamnt"`
	Taxamnt            float64                                     `protobuf:"fixed64,15,opt,name=taxamnt,proto3" json:"taxamnt"`
	Docamnt            float64                                     `protobuf:"fixed64,16,opt,name=docamnt,proto3" json:"docamnt"`
	AmountReceived     *CreateSalesInvoiceGPRequest_AmountReceived `protobuf:"bytes,17,opt,name=amount_received,json=amountReceived,proto3" json:"amount_received"`
	GnlRequestShipDate string                                      `protobuf:"bytes,18,opt,name=gnl_request_ship_date,json=gnlRequestShipDate,proto3" json:"gnl_request_ship_date"`
	GnlRegion          string                                      `protobuf:"bytes,19,opt,name=gnl_region,json=gnlRegion,proto3" json:"gnl_region"`
	GnlWrtId           string                                      `protobuf:"bytes,20,opt,name=gnl_wrt_id,json=gnlWrtId,proto3" json:"gnl_wrt_id"`
	GnlArchetypeId     string                                      `protobuf:"bytes,21,opt,name=gnl_archetype_id,json=gnlArchetypeId,proto3" json:"gnl_archetype_id"`
	GnlOrderChannel    string                                      `protobuf:"bytes,22,opt,name=gnl_order_channel,json=gnlOrderChannel,proto3" json:"gnl_order_channel"`
	GnlSoCodeApps      string                                      `protobuf:"bytes,23,opt,name=gnl_so_code_apps,json=gnlSoCodeApps,proto3" json:"gnl_so_code_apps"`
	GnlTotalweight     float64                                     `protobuf:"fixed64,24,opt,name=gnl_totalweight,json=gnlTotalweight,proto3" json:"gnl_totalweight"`
	Userid             string                                      `protobuf:"bytes,25,opt,name=userid,proto3" json:"userid"`
	Detailitems        []*CreateSalesInvoiceGPRequest_DetailItem   `protobuf:"bytes,26,rep,name=detailitems,proto3" json:"detailitems"`
	Shipmthd           string                                      `protobuf:"bytes,27,opt,name=shipmthd,proto3" json:"shipmthd"`
	Locncode           string                                      `protobuf:"bytes,28,opt,name=locncode,proto3" json:"locncode"`
	Pymtrmid           string                                      `protobuf:"bytes,29,opt,name=pymtrmid,proto3" json:"pymtrmid"`
}

type CreateSalesInvoiceGPRequest_VoucherApply struct {
	GnlVoucherType int32   `protobuf:"varint,1,opt,name=gnl_voucher_type,json=gnlVoucherType,proto3" json:"gnl_voucher_type"`
	GnlVoucherId   string  `protobuf:"bytes,2,opt,name=gnl_voucher_id,json=gnlVoucherId,proto3" json:"gnl_voucher_id"`
	Ordocamt       float64 `protobuf:"fixed64,3,opt,name=ordocamt,proto3" json:"ordocamt"`
}
type CreateSalesInvoiceGPRequest_DetailItem struct {
	Lnitmseq   int32   `protobuf:"varint,1,opt,name=lnitmseq,proto3" json:"lnitmseq"`
	Itemnmbr   string  `protobuf:"bytes,2,opt,name=itemnmbr,proto3" json:"itemnmbr"`
	Locncode   string  `protobuf:"bytes,3,opt,name=locncode,proto3" json:"locncode"`
	Uofm       string  `protobuf:"bytes,4,opt,name=uofm,proto3" json:"uofm"`
	Pricelvl   string  `protobuf:"bytes,5,opt,name=pricelvl,proto3" json:"pricelvl"`
	Quantity   float64 `protobuf:"varint,6,opt,name=quantity,proto3" json:"quantity"`
	Unitprce   float64 `protobuf:"fixed64,7,opt,name=unitprce,proto3" json:"unitprce"`
	Xtndprce   float64 `protobuf:"fixed64,8,opt,name=xtndprce,proto3" json:"xtndprce"`
	GnL_Weight int32   `protobuf:"varint,9,opt,name=gnL_Weight,json=gnLWeight,proto3" json:"gnL_Weight"`
}

type CreateSalesInvoiceGPRequest_AmountReceived struct {
	Amount   float64 `protobuf:"fixed64,1,opt,name=amount,proto3" json:"amount"`
	Chekbkid string  `protobuf:"bytes,2,opt,name=chekbkid,proto3" json:"chekbkid"`
}
