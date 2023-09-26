// Copyright 2020 PT. Eden Pangan Indonesia Indonesia. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

package product_price

import (
	"git.edenfarm.id/cuxs/common"
	"git.edenfarm.id/cuxs/validation"
	"git.edenfarm.id/project-version2/datamodel/model"
	"git.edenfarm.id/project-version2/api/src/auth"
	"git.edenfarm.id/project-version2/api/util"
)

type updateRequest struct {
	UpdateProductPrice []*requestItem `json:"prices" valid:"required"`
	SalabilityStatus   float64        `json:"salability_status" valid:"required"`

	Session *auth.SessionData `json:"-"`
	Price   *model.Price      `json:"-"`
}

type requestItem struct {
	UnitPrice    float64 `json:"unit_price" valid:"required"`
	ProductID    string  `json:"product_id" valid:"required"`
	PriceSetID   string  `json:"price_set_id"`
	CategoryID   string  `json:"category_id"`
	CategoryName string  `json:"category_name"`
	Salability   int     `json:"salability"`

	PriceSet *model.PriceSet
	Product  *model.Product
	Category *model.Category
}

func (c *updateRequest) Validate() *validation.Output {
	o := &validation.Output{Valid: true}
	productCheck := make(map[int64]int8)
	var err error

	for _, v := range c.UpdateProductPrice {
		productId, _ := common.Decrypt(v.ProductID)
		categoryId, _ := common.Decrypt(v.CategoryID)
		v.Product = &model.Product{ID: productId}

		err = v.Product.Read("ID")

		// validation for directory database
		if err != nil {
			o.Failure("id.invalid", util.ErrorMustExistInDirectory("product"))
			return o
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

		// validation category based on selected price
		if v.CategoryID != "" {
			if v.Product.Category.ID != categoryId {
				o.Failure("id.invalid", util.ErrorMustBeSame("category name", "selected product"))
				return o
			}
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

		// validation for unit price gte 0
		if v.UnitPrice < 0 {
			o.Failure("id.invalid", util.ErrorEqualGreater("unit price", "0"))
			return o
		}

		// validation for duplicate product id
		if _, isExist := productCheck[productId]; isExist {
			o.Failure("id.invalid", util.ErrorDuplicateID("product"))
			return o
		} else {
			productCheck[productId] = 1
		}

		v.Product.Category.Read("ID")
		if v.CategoryName != v.Product.Category.Name {
			o.Failure("id.invalid", "Category must be same/valid")
		}

		if int8(v.Salability) != v.Product.Salability {
			o.Failure("id.invalid", "Category must be same/valid")
		}

		// validation shadow price > unit price
		//p := &model.Price{Product: v.Product, PriceSet: v.PriceSet}
		//if err = p.Read("Product", "PriceSet"); err == nil {
		//	//pr, _ := strconv.ParseFloat(p.UnitPrice, 10)
		//	if v.UnitPrice > p.ShadowPrice {
		//		o.Failure("id.invalid", util.ErrorGreater("shadow price", "unit price"))
		//		return o
		//	}
		//}

	}
	return o
}

func (c *updateRequest) Messages() map[string]string {
	messages := map[string]string{
		"unit_price.required": util.ErrorInputRequired("Unit Price"),
	}

	return messages
}
