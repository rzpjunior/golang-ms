package model

import (
	"time"

	"git.edenfarm.id/edenlabs/edenlabs/orm"
)

type SalesAssignmentItem struct {
	ID                    int64      `orm:"column(id)" json:"-"`
	SalesAssignmentID     *int64     `orm:"column(sales_assignment_id)" json:"sales_assignment_id"`
	SalesPersonID         int64      `orm:"column(salesperson_id)" json:"salesperson_id"`
	SalesPersonIDGP       string     `orm:"column(salesperson_id_gp)" json:"salesperson_id_gp"`
	AddressID             int64      `orm:"column(address_id)" json:"address_id"`
	AddressIDGP           string     `orm:"column(address_id_gp)" json:"address_id_gp"`
	CustomerAcquisitionID int64      `orm:"column(customer_acquisition_id)" json:"customer_acquisition_id"`
	Latitude              float64    `orm:"column(latitude)" json:"latitude"`
	Longitude             float64    `orm:"column(longitude)" json:"longitude"`
	Task                  int8       `orm:"column(task)" json:"task"`
	CustomerType          int8       `orm:"column(customer_type)" json:"customer_type"`
	ObjectiveCodes        string     `orm:"column(objective_codes)" json:"objective_codes"`
	ActualDistance        float64    `orm:"column(actual_distance)" json:"actual_distance"`
	OutOfRoute            int8       `orm:"column(out_of_route)" json:"out_of_route"`
	StartDate             time.Time  `orm:"column(start_date)" json:"start_date"`
	EndDate               time.Time  `orm:"column(end_date)" json:"end_date"`
	FinishDate            *time.Time `orm:"column(finish_date)" json:"finish_date"`
	SubmitDate            time.Time  `orm:"column(submit_date)" json:"submit_date"`
	TaskImageUrl          string     `orm:"column(task_image_url)" json:"task_image_url"`
	TaskAnswer            int8       `orm:"column(task_answer)" json:"task_answer"`
	Status                int8       `orm:"column(status)" json:"status"`
	EffectiveCall         int8       `orm:"-" json:"effective_call"`
}

func init() {
	orm.RegisterModel(new(SalesAssignmentItem))
}

func (m *SalesAssignmentItem) TableName() string {
	return "sales_assignment_item"
}
