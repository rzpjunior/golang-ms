package dto

import "time"

type ItemCategoryResponse struct {
	ID                int64                      `json:"id"`
	Code              string                     `json:"code"`
	Region            string                     `json:"region,omitempty"`
	RegionID          string                     `json:"region_id,omitempty"`
	Name              string                     `json:"name,omitempty"`
	Status            int8                       `json:"status,omitempty"`
	StatusConvert     string                     `json:"status_convert,omitempty"`
	CreatedAt         time.Time                  `json:"created_at,omitempty"`
	UpdatedAt         time.Time                  `json:"updated_at,omitempty"`
	ItemCategoryImage *ItemCategoryImageResponse `json:"item_category_images"`
	Regions           []*RegionResponse          `json:"regions,omitempty"`
}

type ItemCategoryRequestCreate struct {
	RegionID []string `json:"region_id" valid:"required"`
	Name     string   `json:"name" valid:"required"`
	ImageUrl string   `json:"image_url" valid:"required"`
}

type ItemCategoryRequestUpdate struct {
	RegionID []string `json:"region_id" valid:"required"`
	Name     string   `json:"name" valid:"required"`
	ImageUrl string   `json:"image_url" valid:"required"`
}
