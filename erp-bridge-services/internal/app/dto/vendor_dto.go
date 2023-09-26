package dto

import (
	"time"

	pb "git.edenfarm.id/project-version3/erp-services/erp-protobuf/gen/proto/bridge_service"
)

type VendorResponse struct {
	ID                     int64     `json:"id"`
	Code                   string    `json:"code"`
	VendorOrganizationID   int64     `json:"vendor_organization_id"`
	VendorClassificationID int64     `json:"vendor_classification_id"`
	SubDistrictID          int64     `json:"sub_district_id"`
	PicName                string    `json:"pic_name"`
	Email                  string    `json:"email"`
	PhoneNumber            string    `json:"phone_number"`
	PaymentTermID          int64     `json:"payment_term_id"`
	PaymentMethodID        int64     `json:"payment_method_id"`
	Rejectable             int32     `json:"rejectable"`
	Returnable             int32     `json:"returnable"`
	Address                string    `json:"address"`
	Note                   string    `json:"note"`
	Status                 int32     `json:"status"`
	Latitude               string    `json:"latitude"`
	Longitude              string    `json:"longitude"`
	CreatedAt              time.Time `json:"created_at"`
	CreatedBy              int64     `json:"created_by"`
}

type CreateVendorGPRequest struct {
	InterID             string `json:"interid"`
	VendorID            string `json:"vendorid"`
	VendName            string `json:"vendname"`
	VendShnm            string `json:"vendshnm"`
	VndChknm            string `json:"vndchknm"`
	VendStts            string `json:"vendstts"`
	VndClsID            string `json:"vndclsid"`
	PrP_Vendor_Org_ID   string `json:"prp_vendor_org_id"`
	PrP_Vendor_CLASF_ID string `json:"prp_vendor_clasf_id"`
	VndCntct            string `json:"vndcntct"`
	AddresS1            string `json:"address1"`
	AddresS2            string `json:"address2"`
	AddresS3            string `json:"address3"`
	City                string `json:"city"`
	State               string `json:"state"`
	Zipcode             string `json:"zipcode"`
	Ccode               string `json:"ccode"`
	Country             string `json:"country"`
	Phnumbr1            string `json:"phnumbr1"`
	Phnumbr2            string `json:"phnumbr2"`
	Phone3              string `json:"phone3"`
	Faxnumbr            string `json:"faxnumbr"`
	Taxschid            string `json:"taxschid"`
	Shipmthd            string `json:"shipmthd"`
	Upszone             string `json:"upszone"`
	Acnmvndr            string `json:"acnmvndr"`
	Vaddcdpr            string `json:"vaddcdpr"`
	Vadcdpad            string `json:"vadcdpad"`
	Vadcdsfr            string `json:"vadcdsfr"`
	Vadcdtro            string `json:"vadcdtro"`
	Vadcd1099           string `json:"vadcd1099"`
	Comment1            string `json:"comment1"`
	Comment2            string `json:"comment2"`
	Prp_payment_method  string `json:"prp_payment_method"`
	Prp_payment_term    string `json:"pymtrmid"`
}

type CreateVendorGPResponse struct {
	Code    int32  `json:"code"`
	Message string `json:"message"`
}
type GetVendorGPList struct {
	PageNumber   int32               `protobuf:"varint,1,opt,name=pageNumber,proto3" json:"pageNumber,omitempty"`
	PageSize     int32               `protobuf:"varint,2,opt,name=pageSize,proto3" json:"pageSize,omitempty"`
	TotalPages   int32               `protobuf:"varint,3,opt,name=totalPages,proto3" json:"totalPages,omitempty"`
	TotalRecords int32               `protobuf:"varint,4,opt,name=totalRecords,proto3" json:"totalRecords,omitempty"`
	Data         []*VendorGPResponse `protobuf:"bytes,5,rep,name=data,proto3" json:"data,omitempty"`
	Succeeded    bool                `protobuf:"varint,6,opt,name=succeeded,proto3" json:"succeeded,omitempty"`
	Errors       []string            `protobuf:"bytes,7,rep,name=errors,proto3" json:"errors,omitempty"`
	Message      string              `protobuf:"bytes,8,opt,name=message,proto3" json:"message,omitempty"`
}
type VendorGPResponse struct {
	VENDORID           string                           `protobuf:"bytes,1,opt,name=VENDORID,proto3" json:"VENDORID,omitempty"`
	VENDNAME           string                           `protobuf:"bytes,2,opt,name=VENDNAME,proto3" json:"VENDNAME,omitempty"`
	VNDCLSID           string                           `protobuf:"bytes,3,opt,name=VNDCLSID,proto3" json:"VNDCLSID,omitempty"`
	VNDCNTCT           string                           `protobuf:"bytes,4,opt,name=VNDCNTCT,proto3" json:"VNDCNTCT,omitempty"`
	ADDRESS1           string                           `protobuf:"bytes,5,opt,name=ADDRESS1,proto3" json:"ADDRESS1,omitempty"`
	ADDRESS2           string                           `protobuf:"bytes,6,opt,name=ADDRESS2,proto3" json:"ADDRESS2,omitempty"`
	ADDRESS3           string                           `protobuf:"bytes,7,opt,name=ADDRESS3,proto3" json:"ADDRESS3,omitempty"`
	CITY               string                           `protobuf:"bytes,8,opt,name=CITY,proto3" json:"CITY,omitempty"`
	STATE              string                           `protobuf:"bytes,9,opt,name=STATE,proto3" json:"STATE,omitempty"`
	ZIPCODE            string                           `protobuf:"bytes,10,opt,name=ZIPCODE,proto3" json:"ZIPCODE,omitempty"`
	COUNTRY            string                           `protobuf:"bytes,11,opt,name=COUNTRY,proto3" json:"COUNTRY,omitempty"`
	PHNUMBR1           string                           `protobuf:"bytes,12,opt,name=PHNUMBR1,proto3" json:"PHNUMBR1,omitempty"`
	PHNUMBR2           string                           `protobuf:"bytes,13,opt,name=PHNUMBR2,proto3" json:"PHNUMBR2,omitempty"`
	PHONE3             string                           `protobuf:"bytes,14,opt,name=PHONE3,proto3" json:"PHONE3,omitempty"`
	FAXNUMBR           string                           `protobuf:"bytes,15,opt,name=FAXNUMBR,proto3" json:"FAXNUMBR,omitempty"`
	UPSZONE            string                           `protobuf:"bytes,16,opt,name=UPSZONE,proto3" json:"UPSZONE,omitempty"`
	SHIPMTHD           string                           `protobuf:"bytes,17,opt,name=SHIPMTHD,proto3" json:"SHIPMTHD,omitempty"`
	TAXSCHID           string                           `protobuf:"bytes,18,opt,name=TAXSCHID,proto3" json:"TAXSCHID,omitempty"`
	ACNMVNDR           string                           `protobuf:"bytes,19,opt,name=ACNMVNDR,proto3" json:"ACNMVNDR,omitempty"`
	TXIDNMBR           string                           `protobuf:"bytes,20,opt,name=TXIDNMBR,proto3" json:"TXIDNMBR,omitempty"`
	VENDSTTS           uint32                           `protobuf:"varint,21,opt,name=VENDSTTS,proto3" json:"VENDSTTS,omitempty"`
	CREATDDT           string                           `protobuf:"bytes,22,opt,name=CREATDDT,proto3" json:"CREATDDT,omitempty"`
	CURNCYID           string                           `protobuf:"bytes,23,opt,name=CURNCYID,proto3" json:"CURNCYID,omitempty"`
	TXRGNNUM           string                           `protobuf:"bytes,24,opt,name=TXRGNNUM,proto3" json:"TXRGNNUM,omitempty"`
	TRDDISCT           uint32                           `protobuf:"varint,25,opt,name=TRDDISCT,proto3" json:"TRDDISCT,omitempty"`
	MINORDER           float64                          `protobuf:"fixed64,26,opt,name=MINORDER,proto3" json:"MINORDER,omitempty"`
	PYMTRMID           *pb.VendorGP_Vpymtrmid           `protobuf:"bytes,27,opt,name=PYMTRMID,proto3" json:"PYMTRMID,omitempty"`
	COMMENT1           string                           `protobuf:"bytes,28,opt,name=COMMENT1,proto3" json:"COMMENT1,omitempty"`
	COMMENT2           string                           `protobuf:"bytes,29,opt,name=COMMENT2,proto3" json:"COMMENT2,omitempty"`
	USERDEF1           string                           `protobuf:"bytes,30,opt,name=USERDEF1,proto3" json:"USERDEF1,omitempty"`
	USERDEF2           string                           `protobuf:"bytes,31,opt,name=USERDEF2,proto3" json:"USERDEF2,omitempty"`
	PYMNTPRI           string                           `protobuf:"bytes,32,opt,name=PYMNTPRI,proto3" json:"PYMNTPRI,omitempty"`
	Organization       *pb.VendorGP_Vorganization       `protobuf:"bytes,33,opt,name=Organization,proto3" json:"Organization,omitempty"`
	Classification     *pb.VendorGP_Vclassification     `protobuf:"bytes,34,opt,name=Classification,proto3" json:"Classification,omitempty"`
	PaymentMethod      *pb.VendorGP_Vpaymentmethod      `protobuf:"bytes,35,opt,name=PaymentMethod,proto3" json:"PaymentMethod,omitempty"`
	LatestGoodsReceipt *pb.VendorGP_Vlatestgoodsreceipt `protobuf:"bytes,36,opt,name=LatestGoodsReceipt,proto3" json:"LatestGoodsReceipt,omitempty"`
	Vaddcdpr           *VendorGP_VVendorAddress         `protobuf:"bytes,37,opt,name=vaddcdpr,proto3" json:"vaddcdpr,omitempty"`
}

type VendorGP_VVendorAddress struct {
	Vendorid              string `protobuf:"bytes,1,opt,name=vendorid,proto3" json:"vendorid,omitempty"`
	Adrscode              string `protobuf:"bytes,2,opt,name=adrscode,proto3" json:"adrscode,omitempty"`
	PrpAdministrativeCode string `protobuf:"bytes,3,opt,name=prp_administrative_code,json=prpAdministrativeCode,proto3" json:"prp_administrative_code,omitempty"`
}
