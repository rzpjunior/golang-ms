package dto

import (
	"time"

	"git.edenfarm.id/project-version3/erp-services/erp-customer-mobile-service/internal/app/model"
)

type ArchetypeResponse struct {
	ID               string              `json:"id"`
	Code             string              `json:"code"`
	CustomerTypeID   string              `json:"customer_type_id"`
	Description      string              `json:"description"`
	Status           string              `json:"status"`
	StatusConvert    string              `json:"status_convert"`
	CreatedAt        time.Time           `json:"created_at"`
	UpdatedAt        time.Time           `json:"updated_at"`
	CustomerGroup    string              `orm:"column(customer_group)" json:"customer_group"`
	Name             string              `orm:"column(name)" json:"name"`
	NameID           string              `orm:"column(name_id)" json:"name_id"`
	Abbreviation     string              `orm:"column(abbreviation)" json:"abbreviation"`
	Note             string              `orm:"column(note)" json:"note"`
	AuxData          string              `orm:"column(aux_data)" json:"aux_data"`
	DocRequired      string              `orm:"-" json:"document_required"`
	DocImageRequired []string            `orm:"-" json:"document_image_required"`
	CustomerType     *model.CustomerType `orm:"column(customer_type_id)" json:"customer_type"`
}
