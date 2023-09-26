package dto

import "time"

type ItemImageResponse struct {
	ID        int64     `json:"id"`
	ItemID    int64     `json:"item_id"`
	ImageUrl  string    `json:"image_url"`
	MainImage int8      `json:"main_image"`
	CreatedAt time.Time `json:"created_at"`
}

type ItemImageRequestUpdate struct {
	Images []ImageRequest `json:"images" valid:"required"`
}

type ImageRequest struct {
	ImageUrl string `json:"image_url" valid:"required"`
}
