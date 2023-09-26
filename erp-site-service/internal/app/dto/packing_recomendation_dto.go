package dto

type PackingRecommendationResponse struct {
	ID           string                `json:"_id"`
	PackingOrder *PackingOrderResponse `json:"packing_order"`
	ItemID       string                `json:"item_id"`
	Item         *ItemResponse         `json:"item"`
	ItemPack     []*ItemPackResponse   `json:"item_pack"`
}

type ItemResponse struct {
	ID          string             `json:"id"`
	Code        string             `json:"code"`
	Name        string             `json:"name"`
	Note        string             `json:"note"`
	Description string             `json:"description"`
	Status      int8               `json:"status"`
	UnitWeight  float64            `json:"unit_weight"`
	OrderMinQty float64            `json:"order_min_qty"`
	OrderMaxQty float64            `json:"order_max_qty"`
	Uom         *UomResponse       `json:"uom"`
	ItemImage   *ItemImageResponse `json:"item_image"`
}

type ItemPackResponse struct {
	PackType          float64 `json:"pack_type"`
	ExpectedTotalPack float64 `json:"expected_total_pack"`
	ActualTotalPack   float64 `json:"actual_total_pack"`
}

type ItemImageResponse struct {
	ID        int64  `json:"id"`
	ImageUrl  string `json:"image_url"`
	MainImage int8   `json:"main_image"`
}
