package dto

import "time"

type ItemSectionResponse struct {
	ID              int64                `json:"id,omitempty"`
	Code            string               `json:"code,omitempty"`
	Name            string               `json:"name,omitempty"`
	BackgroundImage string               `json:"background_image,omitempty"`
	StartAt         time.Time            `json:"start_at,omitempty"`
	FinishAt        time.Time            `json:"finish_at,omitempty"`
	Regions         []string             `json:"regions,omitempty"`
	RegionNames     []string             `json:"region_names,omitempty"`
	Archetypes      []string             `json:"archetypes,omitempty"`
	ArchetypeNames  []string             `json:"archetype_names,omitempty"`
	ItemID          []int64              `json:"item,omitempty"`
	Items           []*ItemResponse      `json:"items,omitempty"`
	Sequence        int                  `json:"sequence,omitempty"`
	Note            string               `json:"note,omitempty"`
	Status          int8                 `json:"status,omitempty"`
	StatusConvert   string               `json:"status_convert,omitempty"`
	CreatedAt       time.Time            `json:"created_at,omitempty"`
	UpdatedAt       time.Time            `json:"updated_at,omitempty"`
	Type            int8                 `json:"type,omitempty"`
	Region          []*RegionResponse    `json:"region,omitempty"`
	Archetype       []*ArchetypeResponse `json:"archetype,omitempty"`
}

type ItemSectionRequestCreate struct {
	Name            string    `json:"name" valid:"required|lte:30"`
	BackgroundImage string    `json:"background_image"`
	StartAt         time.Time `json:"start_at" valid:"required"`
	FinishAt        time.Time `json:"finish_at" valid:"required"`
	Regions         []string  `json:"regions" valid:"required"`
	Archetypes      []string  `json:"archetypes" valid:"required"`
	Items           []int64   `json:"items" valid:"required"`
	Sequence        int       `json:"sequence" valid:"required|numeric"`
	Note            string    `json:"note" valid:"lte:255"`
	Type            int8      `json:"type,omitempty"`
}

type ItemSectionRequestUpdate struct {
	Name            string    `json:"name" valid:"required|lte:30"`
	BackgroundImage string    `json:"background_image"`
	StartAt         time.Time `json:"start_at" valid:"required"`
	FinishAt        time.Time `json:"finish_at" valid:"required"`
	Regions         []string  `json:"regions" valid:"required"`
	Archetypes      []string  `json:"archetypes" valid:"required"`
	Items           []int64   `json:"items" valid:"required"`
	Sequence        int       `json:"sequence" valid:"required|numeric"`
	Note            string    `json:"note" valid:"lte:255"`
}

type ItemSectionRequestArchive struct {
	Note string `json:"note" valid:"required|lte:255"`
}

type ItemSectionRequestGet struct {
	Offset        int64     `json:"offset"`
	Limit         int64     `json:"limit"`
	RegionID      string    `json:"region_id"`
	ArchetypeID   string    `json:"archetype_id"`
	Status        []int32   `json:"status"`
	Search        string    `json:"search"`
	OrderBy       string    `json:"order_by"`
	Type          int8      `json:"type"`
	CurrentTime   time.Time `json:"current_time"`
	ItemSectionID int64     `json:"item_section_id"`
}
