package dto

type ItemResponse struct {
	ID                  string   `json:"id"`
	Code                string   `json:"code"`
	ItemName            string   `json:"item_name,omitempty"`
	ItemUomName         string   `json:"item_uom_name,omitempty"`
	UnitPrice           string   `json:"unit_price,omitempty"`
	Description         string   `json:"description,omitempty"`
	OrderMinQty         string   `json:"order_min_qty,omitempty"`
	DecimalEnabled      string   `json:"decimal_enabled,omitempty"`
	ImageUrl            string   `json:"image_url,omitempty"`
	ItemCategoryName    string   `json:"item_category_name,omitempty"`
	ItemCategoryID      string   `json:"item_category_id,omitempty"`
	ImagesUrlArr        []string `json:"image_url_arr,omitempty"`
	ItemCategoryNameArr []string `json:"item_category_name_arr,omitempty"`
}

type RequestGetItemList struct {
	Platform string       `json:"platform" valid:"required"`
	Offset   int32        `json:"offset"`
	Limit    int32        `json:"limit"`
	Data     dataItemList `json:"data" valid:"required"`
}

type RequestGetPrivateItemList struct {
	Platform string              `json:"platform" valid:"required"`
	Offset   int32               `json:"offset"`
	Limit    int32               `json:"limit"`
	Data     dataItemListPrivate `json:"data" valid:"required"`

	Session *SessionDataCustomer
}

type ItemDetailPrivateRequest struct {
	Platform string                `json:"platform" valid:"required"`
	Data     dataItemPrivateDetail `json:"data" valid:"required"`

	Session *SessionDataCustomer
}

type ItemDetailRequest struct {
	Platform string         `json:"platform" valid:"required"`
	Data     dataItemDetail `json:"data" valid:"required"`
}

type dataItemListPrivate struct {
	AddressID      string `json:"address_id" valid:"required"`
	ItemCategoryID string `json:"item_category_id"`
	Search         string `json:"search"`
	IsSkuDisc      int8   `json:"is_sku_discount"`
}

type dataItemList struct {
	AdmDivisionID  string `json:"adm_division_id" valid:"required"`
	ItemCategoryID string `json:"item_category_id"`
	Search         string `json:"search"`
}

type dataItemPrivateDetail struct {
	ItemID    string `json:"item_id" valid:"required"`
	AddressID string `json:"address_id" valid:"required"`
}

type dataItemDetail struct {
	ItemID        string `json:"item_id" valid:"required"`
	AdmDivisionID string `json:"adm_division_id" valid:"required"`
}

type RequestGetPrivateItemByListID struct {
	Platform string                 `json:"platform" valid:"required"`
	Data     DataGetPrivateByListID `json:"data" valid:"required"`
	Session  *SessionDataCustomer
}

type DataGetPrivateByListID struct {
	AddressID string `json:"address_id" valid:"required"`
	ItemsID   string `json:"items_id" valid:"required"`
}

type RequestGetFinishedItems struct {
	Offset   int64  `json:"offset"`
	Limit    int64  `json:"limit"`
	Platform string `json:"platform" valid:"required"`
	Data     struct {
		AddressID string `json:"address_id" valid:"required"`
	} `json:"data" valid:"required"`

	Session *SessionDataCustomer
}

type LastFinTransItemResponse struct {
	ItemId      string `json:"item_id"`
	ItemName    string `json:"item_name,omitempty"`
	ItemUomName string `json:"item_uom_name,omitempty"`
	UnitPrice   string `json:"unit_price,omitempty"`
	ShadowPrice string `json:"shadow_price,omitempty"`
	OrderQty    string `json:"order_qty,omitempty"`
	Subtotal    string `json:"subtotal,omitempty"`
	Weight      string `json:"weight,omitempty"`
	Note        string `json:"note,omitempty"`
	ImageUrl    string `json:"image_url,omitempty"`
}
