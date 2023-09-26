package dto

type ItemResponse struct {
	ID                          int64   `json:"id"`
	Code                        string  `json:"code,omitempty"`
	ClassID                     string  `json:"class_id,omitempty"`
	Description                 string  `json:"description,omitempty"`
	UnitWeightConversion        float64 `json:"unit_weight_conversion"`
	OrderMinQty                 float64 `json:"order_min_qty"`
	OrderMaxQty                 float64 `json:"order_max_qty"`
	ItemType                    string  `json:"item_type,omitempty"`
	Capitalize                  string  `json:"capitalize,omitempty"`
	Note                        string  `json:"note,omitempty"`
	ExcludeArchetypeName        string  `json:"exclude_archetype_name,omitempty"`
	MaxDayDeliveryDate          int8    `json:"max_day_delivery_date,omitempty"`
	Packable                    bool    `json:"packable"`
	Fragile                     bool    `json:"fragile"`
	Taxable                     string  `json:"taxable,omitempty"`
	OrderChannelRestrictionName string  `json:"order_channel_restriction_name,omitempty"`
	Status                      int8    `json:"status,omitempty"`
	ItemCategoryName            string  `json:"item_category_name,omitempty"`
	Price                       float64 `json:"price,omitempty"`

	Uom                      *UomGPResponse                     `json:"uom"`
	ItemCategory             []*ItemCategoryResponse            `json:"item_category"`
	ItemImages               []*ItemImageResponse               `json:"item_images"`
	ExcludeArchetypes        []*ArchetypeResponse               `json:"exclude_archetypes"`
	OrderChannelRestrictions []*OrderChannelRestrictionResponse `json:"order_channel_restrictions"`
	Class                    *ItemClassResponse                 `json:"class"`
}

type ItemGPResponse struct {
	ID                          int64   `json:"id"`
	Code                        string  `json:"code"`
	ClassID                     string  `json:"class_id,omitempty"`
	Description                 string  `json:"description,omitempty"`
	UnitWeightConversion        float64 `json:"unit_weight_conversion"`
	OrderMinQty                 float64 `json:"order_min_qty"`
	OrderMaxQty                 float64 `json:"order_max_qty"`
	ItemType                    string  `json:"item_type,omitempty"`
	Capitalize                  string  `json:"capitalize,omitempty"`
	Note                        string  `json:"note,omitempty"`
	ExcludeArchetypeName        string  `json:"exclude_archetype_name,omitempty"`
	MaxDayDeliveryDate          int8    `json:"max_day_delivery_date,omitempty"`
	Packable                    bool    `json:"packable"`
	Fragile                     bool    `json:"fragile"`
	Packability                 string  `json:"packability"`
	Taxable                     string  `json:"taxable,omitempty"`
	OrderChannelRestrictionName string  `json:"order_channel_restriction_name,omitempty"`
	Status                      int8    `json:"status,omitempty"`
	ItemCategoryName            string  `json:"item_category_name,omitempty"`
	Price                       float64 `json:"price,omitempty"`
	DecimalEnabled              bool    `json:"decimal_enabled,omitempty"`

	Uom                      *UomGPResponse                     `json:"uom"`
	ItemCategory             []*ItemCategoryResponse            `json:"item_category"`
	ItemImages               []*ItemImageResponse               `json:"item_images"`
	ExcludeArchetypes        []*ArchetypeResponse               `json:"exclude_archetypes"`
	OrderChannelRestrictions []*OrderChannelRestrictionResponse `json:"order_channel_restrictions"`
	Class                    *ItemClassResponse                 `json:"class"`
	ItemPrice                []*ItemPriceResponse               `json:"item_price"`
	ItemSite                 []*ItemSiteResponse                `json:"item_site"`
	ItemPriceTiering         []*ItemPriceTieringResponse        `json:"item_price_tiering"`
}
type ItemRequestGet struct {
	Offset           int
	OffsetQuery      int
	Limit            int
	Status           int
	Search           string
	OrderBy          string
	UomID            string
	ItemCategoryID   int64
	ClassID          string
	CustomerTypeIDGP string
	RegionIDGP       string
	SiteIDGP         string
	Salability       string
	OrderChannel     int32
	ArchetypeIDGP    string
	ItemIdGP         []*string
	PriceLevel       string
	ID               string
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
	ExcludeArchetype        []string `json:"exclude_archetype,omitempty"`
	OrderChannelRestriction []int64  `json:"order_channel_restriction,omitempty"`
	Images                  []string `json:"images" valid:"required"`
	Note                    string   `json:"note"`
}

type OrderChannelRestrictionResponse struct {
	ID        int64  `json:"id"`
	Table     string `json:"table,omitempty"`
	Attribute string `json:"attribute,omitempty"`
	ValueInt  int8   `json:"value_int"`
	ValueName string `json:"value_name"`
	Note      string `json:"note"`
}

type ItemPriceResponse struct {
	Region       string  `protobuf:"bytes,1,opt,name=gnl_region,json=gnlRegion,proto3" json:"gnl_region,omitempty"`
	CustomerType string  `protobuf:"bytes,2,opt,name=gnl_cust_type_id,json=gnlCustTypeId,proto3" json:"gnl_cust_type_id,omitempty"`
	PriceLevel   string  `protobuf:"bytes,3,opt,name=prclevel,proto3" json:"prclevel,omitempty"`
	Price        float64 `protobuf:"fixed64,4,opt,name=price,proto3" json:"price,omitempty"`
}

type ItemSiteResponse struct {
	Region              string  `protobuf:"bytes,1,opt,name=gnl_region,json=gnlRegion,proto3" json:"gnl_region,omitempty"`
	Location            string  `protobuf:"bytes,2,opt,name=locncode,proto3" json:"locncode,omitempty"`
	GnlCbSalability     int32   `protobuf:"varint,3,opt,name=gnl_cb_salability,json=gnlCbSalability,proto3" json:"gnl_cb_salability,omitempty"`
	GnlCbSalabilityDesc string  `protobuf:"bytes,4,opt,name=gnl_cb_salability_desc,json=gnlCbSalabilityDesc,proto3" json:"gnl_cb_salability_desc,omitempty"`
	TotalStock          float64 `protobuf:"fixed64,5,opt,name=total_stock,json=totalStock,proto3" json:"total_stock,omitempty"`
}

type ItemPriceTieringResponse struct {
	Docnumbr          string  `json:"docnumbr"`
	GnlRegion         string  `json:"gnl_region"`
	EffectiveDate     string  `json:"effective_date"`
	GnlMinQty         int32   `json:"gnl_min_qty"`
	GnlDiscountAmount float64 `json:"gnl_discount_amount"`
	GnlQuotaUser      int32   `json:"gnl_quota_user"`
}
