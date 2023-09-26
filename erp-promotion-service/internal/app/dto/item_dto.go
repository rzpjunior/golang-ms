package dto

type ItemResponse struct {
	ID                          int64        `json:"id"`
	Code                        string       `json:"code,omitempty"`
	ClassID                     int64        `json:"class_id,omitempty"`
	Description                 string       `json:"description,omitempty"`
	UnitWeightConversion        float64      `json:"unit_weight_conversion,omitempty"`
	OrderMinQty                 float64      `json:"order_min_qty,omitempty"`
	OrderMaxQty                 float64      `json:"order_max_qty,omitempty"`
	ItemType                    string       `json:"item_type,omitempty"`
	Capitalize                  string       `json:"capitalize,omitempty"`
	Note                        string       `json:"note,omitempty"`
	ExcludeArchetypeName        string       `json:"exclude_archetype_name,omitempty"`
	MaxDayDeliveryDate          int8         `json:"max_day_delivery_date,omitempty"`
	Packable                    bool         `json:"packable,omitempty"`
	Fragile                     bool         `json:"fragile,omitempty"`
	Taxable                     string       `json:"taxable,omitempty"`
	OrderChannelRestrictionName string       `json:"order_channel_restriction_name,omitempty"`
	Status                      int8         `json:"status,omitempty"`
	ItemCategoryName            string       `json:"item_category_name,omitempty"`
	Uom                         *UomResponse `json:"uom,omitempty"`
}
