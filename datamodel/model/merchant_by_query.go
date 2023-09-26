package model

import (
	"encoding/json"

	"git.edenfarm.id/cuxs/common"
	"git.edenfarm.id/cuxs/orm"
)

func init() {
	orm.RegisterModel(new(MerchantByQuery))
}

// ItemByQuery : struct to hold price set model data for database
type MerchantByQuery struct {
	ID           	int64  `orm:"column(id);auto" json:"-"`
	MerchantName    string `orm:"column(merchant_name);null" json:"merchant_name"`
	MerchantCode    string `orm:"column(merchant_code);null" json:"merchant_code"`
	PicName			string `orm:"column(pic_name);null" json:"pic_name"`
	PhoneNumber		string `orm:"column(phone_number);null" json:"phone_number"`
	AltPhoneNumber	string `orm:"column(alt_phone_number);null" json:"alt_phone_number"`
	Email			string `orm:"column(email);null" json:"email"`
	BillingAddress	string `orm:"column(billing_address);null" json:"billing_address"`
	Note			string `orm:"column(note);null" json:"note"`
	Status			string `orm:"column(status);null" json:"status"`
	TagCustomer		string `orm:"column(tag_customer);null" json:"tag_customer"`
	TermInvoice    	string `orm:"column(term_invoice_name);null" json:"term_invoice_name"`
	TermPayment 	string `orm:"column(term_payment);null" json:"term_payment"`
	PaymentMethod   string `orm:"column(payment_method);null" json:"payment_method"`
	BusinessType  	string `orm:"column(business_type);null" json:"business_type"`
}

// MarshalJSON : function to customize data struct into json. encrypt all key ids
func (m *MerchantByQuery) MarshalJSON() ([]byte, error) {
	type Alias MerchantByQuery

	return json.Marshal(&struct {
		ID string `json:"id"`
		*Alias
	}{
		ID:    common.Encrypt(m.ID),
		Alias: (*Alias)(m),
	})
}
