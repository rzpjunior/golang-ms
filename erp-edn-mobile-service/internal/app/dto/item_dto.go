package dto

type ItemResponse struct {
	ID                          string                             `json:"id"`
	Code                        string                             `json:"code,omitempty"`
	Name                        string                             `json:"name,omitempty"`
	ClassID                     int64                              `json:"class_id,omitempty"`
	Description                 string                             `json:"description,omitempty"`
	UnitWeightConversion        float64                            `json:"unit_weight_conversion,omitempty"`
	OrderMinQty                 float64                            `json:"order_min_qty,omitempty"`
	OrderMaxQty                 float64                            `json:"order_max_qty,omitempty"`
	ItemType                    string                             `json:"item_type,omitempty"`
	Capitalize                  string                             `json:"capitalize,omitempty"`
	Note                        string                             `json:"note,omitempty"`
	ExcludeArchetypeName        string                             `json:"exclude_archetype_name,omitempty"`
	MaxDayDeliveryDate          int8                               `json:"max_day_delivery_date,omitempty"`
	Packable                    bool                               `json:"packable,omitempty"`
	Fragile                     bool                               `json:"fragile,omitempty"`
	Taxable                     string                             `json:"taxable,omitempty"`
	OrderChannelRestrictionName string                             `json:"order_channel_restriction_name,omitempty"`
	Status                      int8                               `json:"status,omitempty"`
	ItemCategoryName            string                             `json:"item_category_name,omitempty"`
	UnitPrice                   float64                            `json:"unit_price"`
	RegionID                    string                             `json:"region_id"`
	SiteID                      string                             `json:"site_id"`
	CustomerTypeID              string                             `json:"customer_type_id"`
	PriceLevel                  string                             `json:"price_level"`
	Stock                       float64                            `json:"stock"`
	Uom                         *UomResponse                       `json:"uom,omitempty"`
	ItemCategory                []*ItemCategoryResponse            `json:"item_category,omitempty"`
	ItemImages                  []*ItemImageResponse               `json:"item_images,omitempty"`
	ExcludeArchetypes           []*ArchetypeResponse               `json:"exclude_Archetypes,omitempty"`
	OrderChannelRestrictions    []*OrderChannelRestrictionResponse `json:"order_channel_restrictions,omitempty"`
	Class                       *ClassResponse                     `json:"class,omitempty"`
}

type ItemRequestPackable struct {
	Packable bool `json:"packable,omitempty"`
}

type ItemRequestFragile struct {
	Fragile bool `json:"packable,omitempty"`
}

type ItemRequestUpdate struct {
	ItemCategory            []int64  `json:"item_category,omitempty"`
	MaxDayDeliveryDate      int8     `json:"max_day_delivery_date,omitempty"`
	ExcludeArchetype        []int64  `json:"exclude_archetype,omitempty"`
	OrderChannelRestriction []int64  `json:"order_channel_restriction,omitempty"`
	Images                  []string `json:"images" valid:"required"`
	Note                    string   `json:"note"`
}

type OrderChannelRestrictionResponse struct {
	ValueInt    int64  `json:"value_int"`
	Description string `json:"description"`
}

type ItemListResponse struct {
	Data []*ItemResponse `json:"data"`
}

type ItemListRequest struct {
	Limit      int32  `json:"limit"`
	Offset     int32  `json:"offset"`
	Status     int32  `json:"status"`
	Search     string `json:"search"`
	SiteID     string `json:"site"`
	OrderBy    string `json:"order_by"`
	CustomerID string `json:"customer_id"`
}

type ItemDetailRequest struct {
	Id int32 `json:"id"`
}
