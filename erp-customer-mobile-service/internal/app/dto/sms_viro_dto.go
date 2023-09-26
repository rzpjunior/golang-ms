package dto

import "git.edenfarm.id/project-version3/erp-services/erp-customer-mobile-service/internal/app/model"

type UpdateRequestSMSViro struct {
	Results     []results          `json:"results"`
	OTPOutGoing *model.OtpOutgoing `json:"-"`
}
type results struct {
	Status    status `json:"status"`
	MessageId string `json:"messageId"`
}
type status struct {
	GroupID int `json:"groupId"`
}
