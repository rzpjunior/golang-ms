package model

import (
	"time"

	"git.edenfarm.id/edenlabs/edenlabs/orm"
)

type SalesAssignmentObjective struct {
	ID         int64     `orm:"column(id)" json:"id"`
	Code       string    `orm:"column(code)" json:"code"`
	Name       string    `orm:"column(name)" json:"name"`
	Objective  string    `orm:"column(objective)" json:"objective"`
	SurveyLink string    `orm:"column(survey_link)" json:"survey_link"`
	Status     int8      `orm:"column(status)" json:"status"`
	CreatedAt  time.Time `orm:"column(created_at)" json:"created_at"`
	CreatedBy  int64     `orm:"column(created_by)" json:"created_by"`
	UpdatedAt  time.Time `orm:"column(updated_at)" json:"updated_at"`
}

func init() {
	orm.RegisterModel(new(SalesAssignmentObjective))
}

func (m *SalesAssignmentObjective) TableName() string {
	return "sales_assignment_objective"
}
