package dto

import (
	"time"

	"git.edenfarm.id/project-version3/erp-services/erp-customer-mobile-service/internal/app/model"
)

type RequestGetPrivateBanner struct {
	Platform string               `json:"platform" valid:"required"`
	Data     dataGetPrivateBanner `json:"data" valid:"required"`

	Session *SessionDataCustomer
}

type dataGetPrivateBanner struct {
	AddressID    string `json:"address_id" valid:"required"`
	DataResponse responseGetBanner
}

type responseGetBanner struct {
	Banner []*model.Banner
}

type RequestGetBanner struct {
	Platform string        `json:"platform" valid:"required"`
	Data     dataGetBanner `json:"data" valid:"required"`
}

type dataGetBanner struct {
	AdmDivisionId string `json:"adm_division_id" valid:"required"`
	DataResponse  responseGetBanner
}

// Banner model for audit_log table.
type ResponseBanner struct {
	ID             string    `json:"-"`
	Code           string    `json:"code"`
	Name           string    `json:"name"`
	ImageUrl       string    `json:"image_url"`
	StartDate      time.Time `json:"start_date"`
	EndDate        time.Time `json:"end_date"`
	NavigationType string    `json:"navigate_type"`
	NavigationUrl  string    `json:"navigate_url"`
	Region         string    `json:"region"`
	Archetype      string    `json:"archetype"`
	Queue          string    `json:"queue"`
	Note           string    `json:"note"`
	Status         string    `json:"status"`

	ItemCategory *ItemCategoryResponse `json:"item_category"`
	Item         *ItemResponse         `json:"item"`
	ItemSection  *ItemSectionResponse  `json:"item_section"`

	// log
	CreatedAt time.Time `json:"created_at"`
	CreatedBy string    `json:"created_by,omitempty"`
}
