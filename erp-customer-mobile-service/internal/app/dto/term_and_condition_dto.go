package dto

import "git.edenfarm.id/project-version3/erp-services/erp-customer-mobile-service/internal/app/model"

type TermConditionResponse struct {
	ID          string `json:"id"`
	Application string `json:"application"`
	Version     string `json:"version"`
	Title       string `json:"title"`
	TitleValue  string `json:"title_value"`
	Content     string `json:"content"`
}

type RequestAcceptTNC struct {
	Platform string `json:"platform" valid:"required"`
	//Data     dataRequestAcceptTNC `json:"data" valid:"required"`
	Session *SessionDataCustomer
}

type dataRequestAcceptTNC struct {
	Customer  *model.Customer
	ConfigApp *ApplicationConfigResponse
}
