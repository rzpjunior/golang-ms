package dto

import (
	"time"
)

type ItemSectionResponse struct {
	ID              string          `json:"id"`
	Code            string          `json:"code"`
	Name            string          `json:"name"`
	Region          string          `json:"region"`
	Archetype       string          `json:"archetype"`
	BackgroundImage string          `json:"background_image"`
	StartAt         time.Time       `json:"start_at"`
	EndAt           time.Time       `json:"end_at"`
	Sequence        string          `json:"sequence"`
	Type            string          `json:"type"`
	CreatedAt       time.Time       `json:"created_at"`
	UpdatedAt       time.Time       `json:"updated_at"`
	Item            []*ItemResponse `json:"item,omitempty"`
}
type RequestGetItemSection struct {
	Platform string             `json:"platform" valid:"required"`
	Data     dataGetItemSection `json:"data" valid:"required"`
}

type dataGetItemSection struct {
	AdmDivisionID         string `json:"adm_division_id" valid:"required"`
	Type                  string `json:"type"`
	ItemSectionItemOffset int64  `json:"item_section_item_offset"`
	ItemSectionItemLimit  int64  `json:"item_section_item_limit" valid:"required"`
}

type RequestGetPrivateItemSection struct {
	Platform string                    `json:"platform" valid:"required"`
	Data     dataGetPrivateItemSection `json:"data" valid:"required"`

	Session *SessionDataCustomer
}

type RequestGetPrivateItemSectionDetail struct {
	Platform string                          `json:"platform" valid:"required"`
	Data     dataGetPrivateItemSectionDetail `json:"data" valid:"required"`

	Session *SessionDataCustomer
}
type dataGetPrivateItemSectionDetail struct {
	AddressID             string `json:"address_id" valid:"required"`
	Type                  string `json:"type"`
	ItemSectionItemOffset int64  `json:"item_section_item_offset" valid:"required"`
	ItemSectionItemLimit  int64  `json:"item_section_item_limit" valid:"required"`
	ItemSectionID         string `json:"item_section_id" valid:"required"`
}

type dataGetPrivateItemSection struct {
	AddressID             string `json:"address_id" valid:"required"`
	Type                  string `json:"type"`
	ItemSectionItemOffset int64  `json:"item_section_item_offset"`
	ItemSectionItemLimit  int64  `json:"item_section_item_limit" valid:"required"`
}

type RequestGetDetailItemSection struct {
	Platform string                   `json:"platform" valid:"required"`
	Data     dataGetDetailItemSection `json:"data" valid:"required"`
}

type dataGetDetailItemSection struct {
	ItemSectionID         string `json:"item_section_id" valid:"required"`
	Search                string `json:"search,omitempty"`
	AdmDivisionID         string `json:"adm_division_id" valid:"required"`
	ItemSectionItemOffset int64  `json:"item_section_item_offset" `
	ItemSectionItemLimit  int64  `json:"item_section_item_limit" valid:"required"`
}
