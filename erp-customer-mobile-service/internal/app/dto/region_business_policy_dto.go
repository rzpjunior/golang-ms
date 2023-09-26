package dto

import "git.edenfarm.id/project-version3/erp-services/erp-customer-mobile-service/internal/app/model"

type RegionBusinessPolRequest struct {
	Platform string                   `json:"platform" valid:"required"`
	Data     DataGetRegionBusinessPol `json:"data" valid:"required"`
}

type DataGetRegionBusinessPol struct {
	RegionID       string `json:"region_id" valid:"required"`
	CustomerTypeID string `json:"customer_type_id" valid:"required"`
}

type RegionBusinessPolResponse struct {
	RegionBusinessPolicy *model.RegionBusinessPolicy
	Region               *RegionResponse `json:"region,omitempty"`
	//BusinessType *model.BusinessType

}
