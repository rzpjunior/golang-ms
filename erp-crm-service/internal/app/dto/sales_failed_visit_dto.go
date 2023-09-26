package dto

type SalesFailedVisitResponse struct {
	ID                    int64    `json:"id"`
	SalesAssignmentItemID int64    `json:"sales_asssignment_item_id"`
	FailedStatus          int64    `json:"failed_status"`
	DescriptionFailed     *string  `json:"description_failed"`
	FailedImage           []string `json:"failed_image"`
}

type SalesFailedVisitRequest struct {
	SalesAssignmentItemID int64   `json:"sales_assignment_item_id" valid:"required"`
	FailedStatus          int64   `json:"failed_status" valid:"required"`
	DescriptionFailed     *string `json:"description_failed"`
	FailedImage           string  `json:"failed_image" valid:"required"`
}
