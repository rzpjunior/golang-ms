package dto

type CreateCashReceiptRequest struct {
	Interid        string  `protobuf:"bytes,1,opt,name=interid,proto3" json:"interid"`
	Sopnumbe       string  `protobuf:"bytes,2,opt,name=sopnumbe,proto3" json:"sopnumbe"`
	AmountReceived float64 `protobuf:"fixed64,3,opt,name=amount_received,json=amountReceived,proto3" json:"amount_received"`
	Chekbkid       string  `protobuf:"bytes,4,opt,name=chekbkid,proto3" json:"chekbkid"`
	Docdate        string  `protobuf:"bytes,5,opt,name=docdate,proto3" json:"docdate"`
}
