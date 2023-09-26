// Copyright 2020 PT. Eden Pangan Indonesia. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

package model

import "git.edenfarm.id/edenlabs/edenlabs/orm"

func init() {
	orm.RegisterModel(new(CustomerPriceSet))
}

// CustomerPriceSet model for city table.
type CustomerPriceSet struct {
	ID       int64     `orm:"column(id);auto" json:"-"`
	Customer *Customer `orm:"-"  json:"customer"`
	Region   *Region   `orm:"-" json:"region"`
	PriceSet *PriceSet `orm:"-" json:"price_set"`
}
