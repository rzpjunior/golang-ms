package dto

import "time"

type ItemCategoryImageResponse struct {
	ID             int64     `json:"id"`
	ItemCategoryID int64     `json:"item_category_id"`
	ImageUrl       string    `json:"image_url"`
	CreatedAt      time.Time `json:"created_at"`
}

type ItemCategoryImageRequestUpdate struct {
	ImageUrl string `json:"image_url" valid:"required"`
}
