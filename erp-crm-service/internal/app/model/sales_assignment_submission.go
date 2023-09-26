package model

import (
	"git.edenfarm.id/edenlabs/edenlabs/orm"
)

type SalesAssignmentSubmission struct {
	ID int64 `orm:"column(id)" json:"-"`
}

func init() {
	orm.RegisterModel(new(SalesAssignmentSubmission))
}

func (m *SalesAssignmentSubmission) TableName() string {
	return "sales_assignment_submission"
}
