// Copyright 2020 PT. Eden Pangan Indonesia. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

package model

import (
	"time"

	"git.edenfarm.id/edenlabs/edenlabs/orm"
)

func init() {
	orm.RegisterModel(new(Item))
}

// Item : struct to hold model data for database
type Item struct {
	ID                      int64   `orm:"column(id);auto" json:"-"`
	Code                    string  `orm:"column(code)" json:"code"`
	Name                    string  `orm:"column(name)" json:"name"`
	Note                    string  `orm:"column(note)" json:"note"`
	Description             string  `orm:"column(description)" json:"description"`
	Status                  int8    `orm:"column(status)" json:"status"`
	UnitWeight              float64 `orm:"column(unit_weight)" json:"unit_weight"`
	SiteStoStr              string  `orm:"column(site_sto)" json:"-"`
	SitePurStr              string  `orm:"column(site_pur)" json:"-"`
	SiteSalStr              string  `orm:"column(site_sal)" json:"-"`
	ItemCategory            string  `orm:"column(tag_Item)" json:"tag_Item"`
	OrderChannelRestriction string  `orm:"column(order_channel_restriction)" json:"-"`
	Storability             int8    `orm:"column(storability)" json:"storability"`
	Purchasability          int8    `orm:"column(purchasability)" json:"purchasability"`
	Salability              int8    `orm:"column(salability)" json:"salability"`
	UnivItemCode            string  `orm:"column(up_code)" json:"up_code"`
	OrderMinQty             float64 `orm:"column(order_min_qty)" json:"order_min_qty"`
	OrderMaxQty             float64 `orm:"column(order_max_qty)" json:"order_max_qty"`
	ExcludeArchetype        string  `orm:"column(exclude_archetype)" json:"exclude_archetype"`
	MaxDayDeliveryDate      int64   `orm:"column(max_day_delivery_date)" json:"max_day_delivery_date"`
	Taxable                 int8    `orm:"column(taxable)" json:"taxable"`
	TaxPercentage           float64 `orm:"column(tax_percentage)" json:"tax_percentage"`

	ItemCategoryStr string   `orm:"-" json:"tag_Item_str"`
	SiteSto         []*Site  `orm:"-" json:"site_sto"`
	SitePur         []*Site  `orm:"-" json:"site_pur"`
	SiteSal         []*Site  `orm:"-" json:"site_sal"`
	SiteStoArr      []string `orm:"-" json:"site_sto_arr"`
	SitePurArr      []string `orm:"-" json:"site_pur_arr"`
	SiteSalArr      []string `orm:"-" json:"site_sal_arr"`

	Uom       *Uom         `orm:"-" json:"uom"`
	Category  *Category    `orm:"-" json:"category"`
	ItemImage []*ItemImage `orm:"-" json:"Item_image"`
}

type ItemDetail struct {
	ID                  string    `orm:"-" json:"id"`
	ItemID              int64     `orm:"column(id)" json:"-"`
	Name                string    `orm:"column(name)" json:"item_name"`
	UOM                 string    `orm:"column(uom)" json:"item_uom_name"`
	UnitPrice           float64   `orm:"column(unit_price)" json:"unit_price"`
	ShadowPrice         float64   `orm:"column(shadow_price)" json:"shadow_price"`
	ShadowPricePct      float64   `orm:"column(shadow_price_pct)" json:"shadow_price_pct"`
	Description         string    `orm:"column(description)" json:"description"`
	OrderMinQty         float64   `orm:"column(order_min_qty)" json:"order_min_qty"`
	ItemCategory        string    `orm:"column(tag_Item)" json:"tag_Item"`
	DecimalEnabled      int8      `orm:"column(decimal_enabled)" json:"decimal_enabled"`
	SkuDiscEndPeriod    time.Time `orm:"-" json:"-"`
	SkuDiscEndPeriodStr string    `orm:"-" json:"sku_discount_endtimestamp,omitempty"`
	// SkuDiscountItem     *SkuDiscountItem       `orm:"-" json:"sku_discount_item,omitempty"`
	// SkuDiscountItemTier []*SkuDiscountItemTier `orm:"-" json:"sku_discount_item_tier,omitempty"`
	ImagesUrl         []string `orm:"-" json:"image_url"`
	ItemCategorysName []string `orm:"-" json:"Item_tag_name,omitempty"`
}
