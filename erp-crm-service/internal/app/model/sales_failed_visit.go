package model

import (
	"git.edenfarm.id/edenlabs/edenlabs/orm"
)

type SalesFailedVisit struct {
	ID                    int64   `orm:"column(id)" json:"-"`
	SalesAssignmentItemID int64   `orm:"column(sales_assignment_item_id)" json:"sales_assignment_item_id"`
	FailedStatus          int64   `orm:"column(failed_status)" json:"failed_status"`
	DescriptionFailed     *string `orm:"column(description_failed)" json:"description_failed"`
	FailedImage           string  `orm:"column(failed_image)" json:"failed_image"`
}

func init() {
	orm.RegisterModel(new(SalesFailedVisit))
}

func (m *SalesFailedVisit) TableName() string {
	return "sales_failed_visit"
}
