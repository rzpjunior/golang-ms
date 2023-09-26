package dto

import (
	"time"
)

type SalesPaymentResponse struct {
	ID               int64     `json:"-"`
	Code             string    `json:"code"`
	Status           int8      `json:"status"`
	RecognitionDate  time.Time `json:"recognition_date"`
	Amount           float64   `json:"amount"`
	BankReceiveNum   string    `json:"bank_receive_num"`
	PaidOff          int8      `json:"paid_off"`
	ImageUrl         string    `json:"image_url"`
	Note             string    `json:"note"`
	CreatedAt        time.Time `json:"created_at"`
	CreatedBy        int64     `json:"created_by"`
	CancellationNote string    `json:"cancellation_note,omitempty"`
	ReceivedDate     time.Time `json:"received_date"`
}

type CreateSalesPaymentGPResponse struct {
	Code      int32  `protobuf:"varint,1,opt,name=code,proto3" json:"code,omitempty"`
	Message   string `protobuf:"bytes,2,opt,name=message,proto3" json:"message,omitempty"`
	Docnumber string `protobuf:"bytes,3,opt,name=docnumbr,proto3" json:"docnumbr,omitempty"`
}

type CreateSalesPaymentGPnonPBDResponse struct {
	Code          int32  `protobuf:"varint,1,opt,name=code,proto3" json:"code,omitempty"`
	Message       string `protobuf:"bytes,2,opt,name=message,proto3" json:"message,omitempty"`
	Paymentnumber string `protobuf:"bytes,3,opt,name=paymentnumber,proto3" json:"paymentnumber,omitempty"`
}

type CreateSalesPaymentGPnonPBDRequest struct {
	Interid  string                                       `protobuf:"bytes,1,opt,name=interid,proto3" json:"interid"`
	Bachnumb string                                       `protobuf:"bytes,2,opt,name=bachnumb,proto3" json:"bachnumb"`
	Docnumbr string                                       `protobuf:"bytes,3,opt,name=docnumbr,proto3" json:"docnumbr"`
	Custnmbr string                                       `protobuf:"bytes,4,opt,name=custnmbr,proto3" json:"custnmbr"`
	Docdate  string                                       `protobuf:"bytes,5,opt,name=docdate,proto3" json:"docdate"`
	Cshrctyp int32                                        `protobuf:"varint,6,opt,name=cshrctyp,proto3" json:"cshrctyp"`
	Curncyid string                                       `protobuf:"bytes,7,opt,name=curncyid,proto3" json:"curncyid"`
	Chekbkid string                                       `protobuf:"bytes,8,opt,name=chekbkid,proto3" json:"chekbkid"`
	Ortrxamt float64                                      `protobuf:"fixed64,9,opt,name=ortrxamt,proto3" json:"ortrxamt"`
	Trxdscrn string                                       `protobuf:"bytes,10,opt,name=trxdscrn,proto3" json:"trxdscrn"`
	ApplyTo  []*CreateSalesPaymentGPnonPBDRequest_ApplyTo `protobuf:"bytes,11,rep,name=apply_to,json=applyTo,proto3" json:"apply_to"`
}
type CreateSalesPaymentGPnonPBDRequest_ApplyTo struct {
	Sopnumbe    string  `protobuf:"bytes,1,opt,name=sopnumbe,proto3" json:"sopnumbe,omitempty"`
	ApplyAmount float64 `protobuf:"fixed64,2,opt,name=apply_amount,json=applyAmount,proto3" json:"apply_amount,omitempty"`
}
