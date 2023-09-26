package dto

type PackingOrderItemResponse struct {
	ID                string                `json:"_id"`
	PackingOrder      *PackingOrderResponse `json:"packing_order"`
	ItemID            string                `json:"item_id"`
	Item              *ItemResponse         `json:"item"`
	PackType          float64               `json:"pack_type"`
	ExpectedTotalPack float64               `json:"expected_total_pack"`
	ActualTotalPack   float64               `json:"actual_total_pack"`
	WeightScale       float64               `json:"weight_scale"`
	Code              string                `json:"code"`
	Status            int8                  `json:"status"`
	StatusConvert     string                `json:"status_convert"`
	Site              *SiteResponse         `json:"site,omitempty"`
}

type PackingOrderItemBarcodeResponse struct {
	LinkPrint         string  `json:"link_print"`
	Code              string  `json:"code"`
	ExpectedTotalPack float64 `json:"expected_total_pack"`
	ActualTotalPack   float64 `json:"actual_total_pack"`
}

type PackingOrderItemRequestUpdate struct {
	ItemID    string  `json:"item_id" valid:"required"`
	PackType  float64 `json:"pack_type" valid:"required"`
	TypePrint int8    `json:"type_print" valid:"required"`
}

type PackingOrderItemRequestPrint struct {
	ItemID      string  `json:"item_id" valid:"required"`
	PackType    float64 `json:"pack_type" valid:"required"`
	TypePrint   int8    `json:"type_print" valid:"required"`
	WeightScale float64 `json:"weight_scale" valid:"required"`
}

type PackingOrderItemRequestDispose struct {
	ItemID   string  `json:"item_id" valid:"required"`
	PackType float64 `json:"pack_type" valid:"required"`
}

type PrintPackingOrderRequest struct {
	Pk PackingOrderItemResponse `json:"pk"`
}

type PrintResponse struct {
	Data string `json:"data"`
}
