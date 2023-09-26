// Copyright 2020 PT. Eden Pangan Indonesia. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

package model

import (
	"encoding/json"
	"time"

	"git.edenfarm.id/cuxs/common"
	"git.edenfarm.id/cuxs/orm"
	"git.edenfarm.id/project-version2/datamodel/util"
)

func init() {
	orm.RegisterModel(new(SalesAssignmentItem))
}

// Sales Assignment Item: struct to hold model data for database
type SalesAssignmentItem struct {
	ID               int64                       `orm:"column(id);auto" json:"-"`
	Task             int8                        `orm:"column(task)" json:"-"`
	StartDate        time.Time                   `orm:"column(start_date)" json:"start_date"`
	EndDate          time.Time                   `orm:"column(end_date)" json:"end_date"`
	SubmitDate       time.Time                   `orm:"column(submit_date);type(timestamp);null" json:"submit_date"`
	FinishDate       time.Time                   `orm:"column(finish_date);type(timestamp);null" json:"finish_date"`
	TaskStr          string                      `orm:"-" json:"task"`
	Status           int8                        `orm:"column(status)" json:"status"`
	Latitude         float64                     `orm:"column(latitude)" json:"latitude"`
	Longitude        float64                     `orm:"column(longitude)" json:"longitude"`
	TaskPhoto        string                      `orm:"column(task_photo)" json:"-"`
	CustomerType     int8                        `orm:"column(customer_type)" json:"customer_type"`
	CustomerTypeStr  string                      `orm:"-" json:"customer_type_str"`
	TaskPhotoArr     []string                    `orm:"-" json:"-"`
	TaskPhotoList    []string                    `orm:"-" json:"task_photo_list"`
	AnswerOption     int8                        `orm:"column(answer_option_id)" json:"-"`
	AsnsweOptionStr  string                      `orm:"-" json:"result"`
	ObjectiveCodes   string                      `orm:"column(objective_codes)" json:"-"`
	ObjectiveCodeArr []string                    `orm:"-" json:"-"`
	ObjectiveCode    []*SalesAssignmentObjective `orm:"-" json:"sales_assignment_objective"`
	OutofRoute       int8                        `orm:"column(out_of_route)" json:"out_of_route"`

	PlanVisit               int64   `orm:"-" json:"plan_visit"`
	PlanFollowUp            int64   `orm:"-" json:"plan_follow_up"`
	VisitActual             int64   `orm:"-" json:"visit_actual"`
	FollowUpActual          int64   `orm:"-" json:"follow_up_actual"`
	VisitPercentage         float64 `orm:"-" json:"visit_percentage"`
	FollowUpPercentage      float64 `orm:"-" json:"follow_up_percentage"`
	EffectiveCall           bool    `orm:"-" json:"effective_call"`
	EffectiveCallPercentage float64 `orm:"-" json:"effective_call_percentage"`
	RevenueEffectiveCall    float64 `orm:"-" json:"revenue_effective_call"`
	RevenueTotal            float64 `orm:"-" json:"revenue_total"`
	TotalCA                 int     `orm:"-" json:"total_ca"`

	SalesAssignment     *SalesAssignment     `orm:"column(sales_assignment_id);null;rel(fk)" json:"sales_assignment"`
	Branch              *Branch              `orm:"column(branch_id);null;rel(fk)" json:"branch"`
	SalesPerson         *Staff               `orm:"column(salesperson_id);null;rel(fk)" json:"salesperson"`
	SalesFailedVisit    *SalesFailedVisit    `orm:"-" json:"sales_failed_visit"`
	CustomerAcquisition *CustomerAcquisition `orm:"column(customer_acquisition_id);null;rel(fk)" json:"customer_acquisition"`
}

type TrackerPerformance struct {
	VisitTracker            *VisitFUTracker `json:"visit_tracker"`
	FollowUpTracker         *VisitFUTracker `json:"follow_up_tracker"`
	EffectiveCallPercentage float64         `orm:"-" json:"effective_call_percentage"`
}

type VisitFUTracker struct {
	TotalPlan       int64 `json:"total_plan"`
	TotalFinished   int64 `json:"total_finished"`
	TotalFailed     int64 `json:"total_failed"`
	TotalCancelled  int64 `json:"total_cancelled"`
	TotalOutOfRoute int64 `json:"total_out_of_route"`
}

// MarshalJSON : function to customize data struct into json. encrypt all key ids
func (m *SalesAssignmentItem) MarshalJSON() ([]byte, error) {
	type Alias SalesAssignmentItem

	return json.Marshal(&struct {
		ID            string `json:"id"`
		StatusConvert string `json:"status_convert"`
		*Alias
	}{
		ID:            common.Encrypt(m.ID),
		StatusConvert: util.ConvertStatusDoc(m.Status),
		Alias:         (*Alias)(m),
	})
}

// Save : function to save data into database. will update when it has valid id, create otherwise
func (m *SalesAssignmentItem) Save(fields ...string) (err error) {
	o := orm.NewOrm()
	if m.ID > 0 {
		_, err = o.Update(m, fields...)
	} else {
		m.ID, err = o.Insert(m)
	}
	return
}

// Read : function to get data from database
func (m *SalesAssignmentItem) Read(fields ...string) error {
	o := orm.NewOrm()
	o.Using("read_only")
	return o.Read(m, fields...)
}
