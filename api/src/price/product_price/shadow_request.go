// Copyright 2020 PT. Eden Pangan Indonesia Indonesia. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

package product_price

import (
	"git.edenfarm.id/cuxs/common"
	"git.edenfarm.id/cuxs/orm"
	"git.edenfarm.id/cuxs/validation"
	"git.edenfarm.id/project-version2/datamodel/model"
	"git.edenfarm.id/project-version2/api/src/auth"
	"git.edenfarm.id/project-version2/api/util"
)

type shadowRequest struct {
	UpdateShadowPrice []*shadowsItem `json:"prices" valid:"required"`
	SalabilityStatus  float64        `json:"salability_status" valid:"required"`

	Session *auth.SessionData `json:"-"`
	Price   *model.Price      `json:"-"`
}

type shadowsItem struct {
	ShadowPrice float64 `json:"shadow_price" valid:"required"`
	ProductID   string  `json:"product_id" valid:"required"`
	PriceSetID  string  `json:"price_set_id"`
	ProductTag  string  `json:"product_tag"`
	Salability  int     `json:"salability"`

	PriceSet *model.PriceSet
	Product  *model.Product
}

func (c *shadowRequest) Validate() *validation.Output {
	o := &validation.Output{Valid: true}
	productCheck := make(map[int64]int8)
	var err error
	orSelect := orm.NewOrm()
	orSelect.Using("read_only")

	for _, v := range c.UpdateShadowPrice {
		productId, _ := common.Decrypt(v.ProductID)
		v.Product = &model.Product{ID: productId}
		err = v.Product.Read("ID")

		// validation for directory database
		if err != nil {
			o.Failure("id.invalid", util.ErrorMustExistInDirectory("product"))
			return o
		}
		// validation for product tag relation of product
		var count int
		if v.ProductTag != "" {
			err = orSelect.Raw("SELECT FIND_IN_SET(?, p.tag_product) AS result FROM product p where p.id = ?", v.ProductTag, v.Product.ID).QueryRow(&count)
			if count < 1 {
				o.Failure("id.invalid", util.ErrorMustBeSame("product tag", "selected product"))
				return o
			}
		}

		if v.PriceSetID == "" {
			o.Failure("price_set.invalid", util.ErrorSelectRequired("price set"))
			o.Failure("id.invalid", "Please select price set and reupload file")
			return o
		} else {
			priceSetId, _ := common.Decrypt(v.PriceSetID)
			v.PriceSet = &model.PriceSet{ID: priceSetId}
			v.PriceSet.Read("ID")
		}

		// validation for salability
		if v.Salability != 0 {

			// salability from filter
			if c.SalabilityStatus == 1 || c.SalabilityStatus == 2 {
				if int(c.SalabilityStatus) != v.Salability {
					o.Failure("id.invalid", util.ErrorMustBeSame("salability", "selected product"))
				}
			}

			//salability from product itself
			if v.Product.Salability != int8(v.Salability) {
				o.Failure("id.invalid", util.ErrorMustBeSame("salability", "selected product"))
				return o
			}
		}

		// validation for shadow price gte 0
		if v.ShadowPrice < 0 {
			o.Failure("id.invalid", util.ErrorEqualGreater("shadow price", "0"))
			return o
		}

		// validation for duplicate product id
		if _, isExist := productCheck[productId]; isExist {
			o.Failure("id.invalid", util.ErrorDuplicateID("product"))
			return o
		} else {
			productCheck[productId] = 1
		}

		// validation shadow price > unit price
		p := &model.Price{Product: v.Product, PriceSet: v.PriceSet}
		if err = p.Read("Product", "PriceSet"); err == nil {
			if v.ShadowPrice != 0 {
				if p.UnitPrice >= v.ShadowPrice {
					o.Failure("id.invalid", util.ErrorGreater("shadow price", "unit price"))
				}
			}
		}
	}
	return o
}

func (c *shadowRequest) Messages() map[string]string {
	messages := map[string]string{
		"shadow_price.required": util.ErrorInputRequired("Shadow Price"),
	}

	return messages
}
