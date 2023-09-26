package dto

import "time"

type BannerResponse struct {
	ID             int64             `json:"id,omitempty"`
	Regions        []string          `json:"regions,omitempty"`
	RegionNames    []string          `json:"region_names,omitempty"`
	Archetypes     []string          `json:"archetypes,omitempty"`
	ArchetypeNames []string          `json:"archetype_names,omitempty"`
	Code           string            `json:"code,omitempty"`
	Name           string            `json:"name,omitempty"`
	Queue          int               `json:"queue,omitempty"`
	Redirect       *RedirectResponse `json:"redirect,omitempty"`
	ImageUrl       string            `json:"image_url,omitempty"`
	StartAt        time.Time         `json:"start_at,omitempty"`
	FinishAt       time.Time         `json:"finish_at,omitempty"`
	Note           string            `json:"note,omitempty"`
	Status         int8              `json:"status,omitempty"`
	CreatedAt      time.Time         `json:"created_at,omitempty"`
	UpdatedAt      time.Time         `json:"updated_at,omitempty"`
	StatusConvert  string            `json:"status_convert,omitempty"`
}

type RedirectResponse struct {
	To        int8        `json:"to,omitempty"`
	Value     interface{} `json:"value,omitempty"`
	Name      string      `json:"name,omitempty"`
	ValueName string      `json:"value_name,omitempty"`
}

type RedirectToItemResponse struct {
	ID   int64  `json:"id"`
	Code string `json:"code"`
	Name string `json:"description"`
}

type RedirectToItemCategoryResponse struct {
	ID   int64  `json:"id"`
	Code string `json:"code"`
	Name string `json:"name"`
}

type RedirectToItemSectionResponse struct {
	ID   int64  `orm:"column(id)" json:"-"`
	Name string `orm:"column(name)" json:"name"`
}

type BannerRequestCreate struct {
	Regions       []string  `json:"regions" valid:"required"`
	Archetypes    []string  `json:"archetypes" valid:"required"`
	Name          string    `json:"name" valid:"required|lte:20"`
	RedirectTo    int       `json:"redirect_to" valid:"required"`
	RedirectValue string    `json:"redirect_value"`
	ImageUrl      string    `json:"image_url" valid:"required"`
	StartAt       time.Time `json:"start_at" valid:"required"`
	FinishAt      time.Time `json:"finish_at" valid:"required"`
	Queue         int       `json:"queue" valid:"required"`
	Note          string    `json:"note"`
}
type BannerRequestArchive struct {
	Note string `json:"note"`
}

type BannerRequestGet struct {
	Offset      int64     `json:"offset"`
	Limit       int64     `json:"limit"`
	RegionID    string    `json:"region_id"`
	ArchetypeID string    `json:"chetype_id"`
	Status      []int32   `json:"status"`
	Search      string    `json:"search"`
	OrderBy     string    `json:"order_by"`
	CurrentTime time.Time `json:"current_time"`
}
