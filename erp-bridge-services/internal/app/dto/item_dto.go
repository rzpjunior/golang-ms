package dto

import "time"

type ItemResponse struct {
	ID                      int64     `json:"id"`
	Code                    string    `json:"code"`
	UomID                   int64     `json:"uom_id"`
	ClassID                 int64     `json:"class_id"`
	ItemCategoryID          int64     `json:"item_category"`
	Description             string    `json:"description"`
	UnitWeightConversion    float64   `json:"unit_weight_conversion"`
	OrderMinQty             float64   `json:"order_min_qty"`
	OrderMaxQty             float64   `json:"order_max_qty"`
	ItemType                string    `json:"item_type"`
	Packability             string    `json:"packability"`
	Capitalize              string    `json:"capitalize"`
	Note                    string    `json:"note"`
	ExcludeArchetype        string    `json:"exclude_archetype"`
	MaxDayDeliveryDate      int8      `json:"max_day_delivery_date"`
	FragileGoods            string    `json:"fragile_goods"`
	Taxable                 string    `json:"taxable"`
	OrderChannelRestriction string    `json:"order_channel_restriction"`
	Status                  int8      `json:"status"`
	StatusConvert           string    `json:"status_convert"`
	CreatedAt               time.Time `json:"created_at"`
	UpdatedAt               time.Time `json:"updated_at"`

	Uom *UomResponse `json:"uom"`
}
