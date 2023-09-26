package dto

import "time"

type SalesOrderResponse struct {
	ID            int64     `json:"id,omitempty"`
	Code          string    `json:"code,omitempty"`
	DocNumber     string    `json:"doc_number,omitempty"`
	AddressID     int64     `json:"address_id,omitempty"`
	CustomerID    int64     `json:"customer_id,omitempty"`
	SalespersonID int64     `json:"salesperson_id,omitempty"`
	WrtID         int64     `json:"wrt_id,omitempty"`
	PaymentTermID int64     `json:"payment_term_id,omitempty"`
	OrderTypeID   int64     `json:"order_type_id,omitempty"`
	SiteID        int64     `json:"site_id,omitempty"`
	Application   int8      `json:"application,omitempty"`
	Status        int8      `json:"status,omitempty"`
	StatusConvert string    `json:"status_convert,omitempty"`
	OrderDate     time.Time `json:"order_date,omitempty"`
	Total         float64   `json:"total,omitempty"`
	CreatedDate   time.Time `json:"created_date,omitempty"`
	ModifiedDate  time.Time `json:"modified_date,omitempty"`
	FinishedDate  time.Time `json:"finished_date,omitempty"`
	CreatedAt     time.Time `json:"created_at,omitempty"`
	UpdatedAt     time.Time `json:"updated_at,omitempty"`
}

type CreateSalesOrderGPRequest struct {
	Interid            string                                  `protobuf:"bytes,1,opt,name=interid,proto3" json:"interid"`
	Sopnumbe           string                                  `protobuf:"bytes,2,opt,name=sopnumbe,proto3" json:"sopnumbe"`
	Docid              string                                  `protobuf:"bytes,3,opt,name=docid,proto3" json:"docid"`
	Docdate            string                                  `protobuf:"bytes,4,opt,name=docdate,proto3" json:"docdate"`
	Custnmbr           string                                  `protobuf:"bytes,5,opt,name=custnmbr,proto3" json:"custnmbr"`
	Custname           string                                  `protobuf:"bytes,6,opt,name=custname,proto3" json:"custname"`
	Prstadcd           string                                  `protobuf:"bytes,7,opt,name=prstadcd,proto3" json:"prstadcd"`
	Curncyid           string                                  `protobuf:"bytes,8,opt,name=curncyid,proto3" json:"curncyid"`
	Subtotal           float64                                 `protobuf:"varint,9,opt,name=subtotal,proto3" json:"subtotal"`
	Trdisamt           float64                                 `protobuf:"varint,10,opt,name=trdisamt,proto3" json:"trdisamt"`
	Freight            float64                                 `protobuf:"varint,11,opt,name=freight,proto3" json:"freight"`
	Miscamnt           float64                                 `protobuf:"varint,12,opt,name=miscamnt,proto3" json:"miscamnt"`
	Taxamnt            float64                                 `protobuf:"varint,13,opt,name=taxamnt,proto3" json:"taxamnt"`
	Docamnt            float64                                 `protobuf:"varint,14,opt,name=docamnt,proto3" json:"docamnt"`
	GnlRequestShipDate string                                  `protobuf:"bytes,15,opt,name=gnl_request_ship_date,json=gnlRequestShipDate,proto3" json:"gnl_request_ship_date"`
	GnlRegion          string                                  `protobuf:"bytes,16,opt,name=gnl_region,json=gnlRegion,proto3" json:"gnl_region"`
	GnlWrtId           string                                  `protobuf:"bytes,17,opt,name=gnl_wrt_id,json=gnlWrtId,proto3" json:"gnl_wrt_id"`
	GnlArchetypeId     string                                  `protobuf:"bytes,18,opt,name=gnl_archetype_id,json=gnlArchetypeId,proto3" json:"gnl_archetype_id"`
	GnlOrderChannel    string                                  `protobuf:"bytes,19,opt,name=gnl_order_channel,json=gnlOrderChannel,proto3" json:"gnl_order_channel"`
	GnlSoCodeApps      string                                  `protobuf:"bytes,20,opt,name=gnl_so_code_apps,json=gnlSoCodeApps,proto3" json:"gnl_so_code_apps"`
	GnlTotalweight     float64                                 `protobuf:"varint,21,opt,name=gnl_totalweight,json=gnlTotalweight,proto3" json:"gnl_totalweight"`
	Userid             string                                  `protobuf:"bytes,22,opt,name=userid,proto3" json:"userid"`
	Detailitems        []*CreateSalesOrderGPRequest_DetailItem `protobuf:"bytes,23,rep,name=detailitems,proto3" json:"detailitems"`
	Locncode           string                                  `protobuf:"bytes,24,opt,name=locncode,proto3" json:"locncode"`
	Shipmthd           string                                  `protobuf:"bytes,25,opt,name=shipmthd,proto3" json:"shipmthd"`
	Pymtrmid           string                                  `protobuf:"bytes,28,opt,name=pymtrmid,proto3" json:"pymtrmid,omitempty"`
	Note               string                                  `protobuf:"bytes,29,opt,name=note,proto3" json:"note"`
	VoucherApply       []*VoucherApplyRequest                  `json:"voucher_apply,omitempty"`
}

type CreateSalesOrderGPRequest_DetailItem struct {
	Sopnumbe   string  `protobuf:"bytes,1,opt,name=sopnumbe,proto3" json:"sopnumbe"`
	Itemnmbr   string  `protobuf:"bytes,2,opt,name=itemnmbr,proto3" json:"itemnmbr"`
	Itemdesc   string  `protobuf:"bytes,3,opt,name=itemdesc,proto3" json:"itemdesc"`
	Locncode   string  `protobuf:"bytes,4,opt,name=locncode,proto3" json:"locncode"`
	Uofm       string  `protobuf:"bytes,5,opt,name=uofm,proto3" json:"uofm"`
	Pricelvl   string  `protobuf:"bytes,6,opt,name=pricelvl,proto3" json:"pricelvl"`
	Quantity   float64 `protobuf:"varint,7,opt,name=quantity,proto3" json:"quantity"`
	Unitprce   float64 `protobuf:"varint,8,opt,name=unitprce,proto3" json:"unitprce"`
	Xtndprce   float64 `protobuf:"varint,9,opt,name=xtndprce,proto3" json:"xtndprce"`
	GnL_Weight float64 `protobuf:"varint,10,opt,name=gnL_Weight,json=gnLWeight,proto3" json:"gnL_Weight"`
}

type VoucherApplyRequest struct {
	GnlVoucherType int8    `json:"gnl_voucher_type"`
	GnlVoucherID   string  `json:"gnl_voucher_id"`
	Ordocamt       float64 `json:"ordocamt"`
}
