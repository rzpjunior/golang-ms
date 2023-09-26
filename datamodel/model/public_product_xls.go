// Copyright 2020 PT. Eden Pangan Indonesia. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

package model

type PublicProductForXls struct {
	ProductName          string `orm:"column(product_name)"`
	UOM                  string `orm:"column(uom)"`
	ProductImages        string `orm:"column(product_images)"`
	DashboardProductName string `orm:"column(dashboard_product_name)"`
}

type ProductMatchingTemplate struct {
	DashboardProductCode string `orm:"column(dashboard_product_code)"`
	DashboardProductName string `orm:"column(dashboard_product_name)"`
	PublicProduct1       string `orm:"column(public_product_1_name)"`
	PublicProduct2       string `orm:"column(public_product_2_name)"`
}
