// Copyright 2020 PT. Eden Pangan Indonesia. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

package model

import (
	"encoding/json"
	"time"

	"git.edenfarm.id/cuxs/common"
	"git.edenfarm.id/cuxs/orm"
)

func init() {
	orm.RegisterModel(new(PublicPrice2))
}

type PublicPrice2 struct {
	ID                 int64           `orm:"column(id)" json:"id"`
	ScrapedDate        string          `orm:"column(scraped_date)" json:"scraped_date"`
	Price              float64         `orm:"column(price);digits(12);decimals(2)" json:"price"`
	Discount           float64         `orm:"column(discount);digits(12);decimals(2)" json:"discount"`
	PriceAfterDiscount float64         `orm:"column(price_after_discount);digits(12);decimals(2)" json:"price_after_discount"`
	Area               *PublicArea2    `orm:"column(area_id);null;rel(fk)" json:"area"`
	Product            *PublicProduct2 `orm:"column(product_id);null;rel(fk)" json:"product"`
	CreatedAt          time.Time       `orm:"column(created_at);type(datetime)" json:"created_at"`
	LastUpdatedAt      time.Time       `orm:"column(last_updated_at);type(datetime);null" json:"last_updated_at"`
}

// TableName : set table name used by model
func (PublicPrice2) TableName() string {
	return "public_price_2"
}

// MarshalJSON : function to customize data struct into json. encrypt all key ids
func (m *PublicPrice2) MarshalJSON() ([]byte, error) {
	type Alias PublicPrice2

	return json.Marshal(&struct {
		ID string `json:"id"`
		*Alias
	}{
		ID:    common.Encrypt(m.ID),
		Alias: (*Alias)(m),
	})
}

// Save : function to save data into database. will update when it has valid id, create otherwise
func (m *PublicPrice2) Save(fields ...string) (err error) {
	o := orm.NewOrm()
	o.Using("scrape")

	if m.ID > 0 {
		_, err = o.Update(m, fields...)
	} else {
		m.ID, err = o.Insert(m)
	}
	return
}

// Read : function to get data from database
func (m *PublicPrice2) Read(fields ...string) error {
	o := orm.NewOrm()
	o.Using("scrape")

	return o.Read(m, fields...)
}
