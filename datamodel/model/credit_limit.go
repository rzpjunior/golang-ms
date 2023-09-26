package model

import (
	"encoding/json"

	"git.edenfarm.id/cuxs/common"
	"git.edenfarm.id/cuxs/orm"
)

func init() {
	orm.RegisterModel(new(CreditLimit))
}

// CreditLimit : struct to hold model data for database
type CreditLimit struct {
	ID                      int64   `orm:"column(id);auto" json:"-"`
	BusinessTypeCreditLimit int8    `orm:"column(business_type_credit_limit)" json:"business_type_credit_limit"`
	AmountCreditLimit       float64 `orm:"column(amount_credit_limit)" json:"amount_credit_limit"`

	BusinessType *BusinessType `orm:"column(business_type_id);null;rel(fk)" json:"business_type,omitempty"`
	PaymentTerm  *SalesTerm    `orm:"column(term_payment_sls_id);null;rel(fk)" json:"payment_term,omitempty"`
}

// MarshalJSON : function to customize data struct into json. encrypt all key ids
func (m *CreditLimit) MarshalJSON() ([]byte, error) {
	type Alias CreditLimit

	return json.Marshal(&struct {
		ID string `json:"id"`
		*Alias
	}{
		ID:    common.Encrypt(m.ID),
		Alias: (*Alias)(m),
	})
}

// Save : function to save data into database. will update when it has valid id, create otherwise
func (m *CreditLimit) Save(fields ...string) (err error) {
	o := orm.NewOrm()
	if m.ID > 0 {
		_, err = o.Update(m, fields...)
	} else {
		m.ID, err = o.Insert(m)
	}
	return
}

// Read : function to get data from database
func (m *CreditLimit) Read(fields ...string) error {
	o := orm.NewOrm()
	o.Using("read_only")
	return o.Read(m, fields...)
}
