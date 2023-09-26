// Copyright 2020 PT. Eden Pangan Indonesia Indonesia. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

package sales_assignment

import (
	"strconv"
	"strings"
	"time"

	"git.edenfarm.id/cuxs/common"
	"git.edenfarm.id/cuxs/validation"
	"git.edenfarm.id/project-version2/api/datastore/repository"
	"git.edenfarm.id/project-version2/api/src/auth"
	"git.edenfarm.id/project-version2/api/util"
	"git.edenfarm.id/project-version2/datamodel/model"
)

type createRequest struct {
	Code               string      `json:"-"`
	AssignmentList     []*taskItem `json:"assignments" valid:"required"`
	SalesGroupSelected string      `json:"sales_group_id" valid:"required"`
	StartDateBatch     time.Time   `json:"-"`
	EndDateBatch       time.Time   `json:"-"`

	Session    *auth.SessionData `json:"-"`
	SalesGroup *model.SalesGroup
}

type taskItem struct {
	CustomerType    string    `json:"customer_type" valid:"required"`
	SalesGroupID    string    `json:"sales_group_id" valid:"required"`
	BranchID        string    `json:"branch_id" valid:"required"`
	SalespersonID   string    `json:"salesperson_id" valid:"required"`
	Task            string    `json:"task" valid:"required"`
	VisitDateStr    string    `json:"visit_date" valid:"required"`
	ObjectiveCode   string    `json:"objective_code"`
	StartDate       time.Time `json:"-"`
	EndDate         time.Time `json:"-"`
	CustomerTypeInt int8      `json:"-"`

	SalesGroup           *model.SalesGroup
	Branch               *model.Branch
	CustomerAcquisition  *model.CustomerAcquisition
	Salesperson          *model.Staff
	GlossaryTask         *model.Glossary
	GlossaryCustomerType *model.Glossary
	ObjectiveCodeModel   *model.SalesAssignmentObjective
}

func (c *createRequest) Validate() *validation.Output {
	o := &validation.Output{Valid: true}
	var err error
	var salesGroupID int64
	layout := "2006-01-02"
	loc, _ := time.LoadLocation("Asia/Jakarta")
	today, _ := time.ParseInLocation(layout, time.Now().Format(layout), loc)

	if c.Code, err = util.CheckTable("sales_assignment"); err != nil {
		o.Failure("code.invalid", util.ErrorInvalidData("code"))
	}

	if salesGroupID, err = common.Decrypt(c.SalesGroupSelected); err != nil {
		o.Failure("sales_group_id.invalid", util.ErrorInvalidData("sales group"))
		return o
	}

	if c.SalesGroup, err = repository.ValidSalesGroup(salesGroupID); err != nil {
		o.Failure("sales_group_id.invalid", util.ErrorInvalidData("sales group"))
		return o
	}

	if c.SalesGroup.Status != 1 {
		o.Failure("sales_group_id.inactive", util.ErrorActive("sales group"))
		return o
	}

	for i, v := range c.AssignmentList {
		salesGroupId, _ := common.Decrypt(v.SalesGroupID)
		outletID, _ := common.Decrypt(v.BranchID)
		salespersonId, _ := common.Decrypt(v.SalespersonID)

		v.SalesGroup = &model.SalesGroup{ID: salesGroupId}
		v.Salesperson = &model.Staff{ID: salespersonId}

		// check glossary if customer type is exist
		v.GlossaryCustomerType = &model.Glossary{
			Table:     "sales_assignment_item",
			Attribute: "customer_type",
			ValueName: v.CustomerType,
		}

		if err = v.GlossaryCustomerType.Read("Table", "Attribute", "ValueName"); err != nil {
			o.Failure("task_"+strconv.Itoa(i)+".invalid", util.ErrorMustExistInDirectory("task"))
			return o
		}

		if v.CustomerType == "Existing Customer" {
			v.Branch = &model.Branch{ID: outletID}
			if err = v.Branch.Read("ID"); err != nil {
				o.Failure("id_"+strconv.Itoa(i)+".invalid", util.ErrorMustExistInDirectory("branch"))
				return o
			}

			if v.Branch.Status != 1 {
				o.Failure("branch_id_"+strconv.Itoa(i)+".inactive", util.ErrorActive("branch"))
				return o
			}

			if v.Branch.Salesperson.ID != v.Salesperson.ID {
				o.Failure("salesperson_id_"+strconv.Itoa(i)+".invalid", util.ErrorMustBeSame("salesperson", "branch salesperson"))
				return o
			}

		} else if v.CustomerType == "Customer Acquisition" {
			v.CustomerAcquisition = &model.CustomerAcquisition{ID: outletID}
			if err = v.CustomerAcquisition.Read("ID"); err != nil {
				o.Failure("id_"+strconv.Itoa(i)+".invalid", util.ErrorMustExistInDirectory("branch"))
				return o
			}

			if v.CustomerAcquisition.Status != 2 {
				o.Failure("branch_id_"+strconv.Itoa(i)+".inactive", util.ErrorActive("branch"))
				return o
			}

			if v.CustomerAcquisition.Salesperson.ID != v.Salesperson.ID {
				o.Failure("salesperson_id_"+strconv.Itoa(i)+".invalid", util.ErrorMustBeSame("salesperson", "branch salesperson"))
				return o
			}
		} else {
			o.Failure("customer_type_"+strconv.Itoa(i)+".invalid", util.ErrorInvalidData("customer type"))
		}

		// validation for directory database
		if err = v.SalesGroup.Read("ID"); err != nil {
			o.Failure("id_"+strconv.Itoa(i)+".invalid", util.ErrorMustExistInDirectory("sales group"))
			return o
		}

		if err = v.Salesperson.Read("ID"); err != nil {
			o.Failure("id_"+strconv.Itoa(i)+".invalid", util.ErrorMustExistInDirectory("salesperson"))
			return o
		}

		if v.SalesGroup.Status != 1 {
			o.Failure("sales_group_id_"+strconv.Itoa(i)+".inactive", util.ErrorActive("sales group"))
			return o
		}

		if v.SalesGroup.ID != c.SalesGroup.ID {
			o.Failure("sales_group_id_"+strconv.Itoa(i)+".invalid", util.ErrorMustBeSame("sales group", "selected sales group"))
			return o
		}

		if v.Salesperson.Status != 1 {
			o.Failure("salesperson_id_"+strconv.Itoa(i)+".inactive", util.ErrorActive("sales person"))
			return o
		}

		if v.StartDate, err = time.ParseInLocation(layout, v.VisitDateStr, loc); err != nil {
			o.Failure("visit_date_"+strconv.Itoa(i)+".invalid", util.ErrorInvalidData("visit date"))
			return o
		}

		if v.EndDate, err = time.ParseInLocation(layout, v.VisitDateStr, loc); err != nil {
			o.Failure("visit_date_"+strconv.Itoa(i)+".invalid", util.ErrorInvalidData("visit date"))
			return o
		}

		if v.StartDate.Before(today) {
			o.Failure("visit_date_"+strconv.Itoa(i)+".invalid", util.ErrorEqualGreater("visit date", "today date"))
			return o
		}

		// check glossary if task is exist
		v.GlossaryTask = &model.Glossary{
			Table:     "sales_assignment_item",
			Attribute: "task",
			ValueName: v.Task,
		}
		if err = v.GlossaryTask.Read("Table", "Attribute", "ValueName"); err != nil {
			o.Failure("task_"+strconv.Itoa(i)+".invalid", util.ErrorMustExistInDirectory("task"))
			return o
		}

		// get first index as initial start date and end date
		if i == 0 {
			c.StartDateBatch = v.StartDate
			c.EndDateBatch = v.EndDate
		}

		// compare to existing start date
		if v.StartDate.Before(c.StartDateBatch) {
			c.StartDateBatch = v.StartDate
		}

		// compare to existing end date
		if v.EndDate.After(c.EndDateBatch) {
			c.EndDateBatch = v.EndDate
		}

		// validation objective codes
		if v.ObjectiveCode != "" {
			codes := strings.Split(v.ObjectiveCode, ",")
			for _, code := range codes {
				v.ObjectiveCodeModel = &model.SalesAssignmentObjective{
					Code: code,
				}
				if err = v.ObjectiveCodeModel.Read("Code"); err != nil {
					o.Failure("objective_code_"+strconv.Itoa(i)+"_"+code+".invalid", util.ErrorInvalidData("objective code"))
				}

				if v.ObjectiveCodeModel.Status != 1 {
					o.Failure("objective_code_"+strconv.Itoa(i)+"_"+code+".invalid", util.ErrorActive("objective code"))
				}
			}
		}

	}
	return o
}

func (c *createRequest) Messages() map[string]string {
	messages := map[string]string{
		"assignments.required":               util.ErrorInputRequired("assignment"),
		"assignment.customer_type.required":  util.ErrorInputRequired("customer type"),
		"assignment.sales_group_id.required": util.ErrorInputRequired("sales group"),
		"assignment.branch_id.required":      util.ErrorInputRequired("branch"),
		"assignment.salesperson_id.required": util.ErrorInputRequired("salesperson"),
		"assignment.task.required":           util.ErrorInputRequired("task"),
		"assignment.visit_date.required":     util.ErrorInputRequired("visit date"),
		"sales_group_id.required":            util.ErrorInputRequired("sales group"),
	}

	return messages
}
