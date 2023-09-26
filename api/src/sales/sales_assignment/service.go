// Copyright 2020 PT. Eden Pangan Indonesia. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

package sales_assignment

import (
	"time"

	"fmt"
	"os"
	"strconv"

	"git.edenfarm.id/cuxs/common"
	"git.edenfarm.id/cuxs/env"
	"git.edenfarm.id/cuxs/orm"
	"git.edenfarm.id/project-version2/api/log"
	"git.edenfarm.id/project-version2/api/util"
	"git.edenfarm.id/project-version2/datamodel/model"
	"github.com/tealeg/xlsx"
)

// getBranchBySalesGroupXls : function to create excel file of task assignment
func getBranchBySalesGroupXls(data []*templateBranchBySalesGroup) (filePath string, err error) {
	var file *xlsx.File
	var sheet *xlsx.Sheet
	var row *xlsx.Row

	dir := env.GetString("EXPORT_DIRECTORY", "")

	filename := fmt.Sprintf("TaskAssignment%s_%s.xlsx", time.Now().Format("2006-01-02"), util.GenerateRandomDoc(5))
	fileDir := fmt.Sprintf("%s/%s", dir, filename)

	file = xlsx.NewFile()
	if sheet, err = file.AddSheet("Sheet1"); err == nil {
		row = sheet.AddRow()
		row.SetHeight(20)
		row.AddCell().Value = "No"
		row.AddCell().Value = "Sales_Group_ID"
		row.AddCell().Value = "Sales Group"
		row.AddCell().Value = "Customer_Type"
		row.AddCell().Value = "Outlet_ID"
		row.AddCell().Value = "Outlet Code"
		row.AddCell().Value = "Outlet Name"
		row.AddCell().Value = "Subdistrict"
		row.AddCell().Value = "District"
		row.AddCell().Value = "Staff_ID"
		row.AddCell().Value = "Staff Code"
		row.AddCell().Value = "Staff Name"
		row.AddCell().Value = "Task"
		row.AddCell().Value = "Visit_Date"
		row.AddCell().Value = "Objective_Code"

		for i, v := range data {
			row = sheet.AddRow()
			row.AddCell().SetInt(i + 1)
			sgID, _ := strconv.Atoi(common.Encrypt(int(v.SalesGroupID)))
			spID, _ := strconv.Atoi(common.Encrypt(int(v.SalespersonID)))
			bID, _ := strconv.Atoi(common.Encrypt(int(v.BranchID)))
			row.AddCell().SetInt(sgID)              // Sales Group ID
			row.AddCell().Value = v.SalesGroupName  // Sales Group
			row.AddCell().Value = v.CustomerType    // Customer_Type
			row.AddCell().SetInt(bID)               // Outlet ID
			row.AddCell().Value = v.BranchCode      // Outlet Code
			row.AddCell().Value = v.OutletName      // Outlet Name
			row.AddCell().Value = v.SubDistrictName // Subdistrict Name
			row.AddCell().Value = v.DistrictName    // District Name
			row.AddCell().SetInt(spID)              // Staff ID
			row.AddCell().Value = v.StaffCode       // Staff Code
			row.AddCell().Value = v.StaffName       // Staff Name
		}

	}
	err = file.Save(fileDir)
	filePath, err = util.UploadToS3(filename, fileDir, "application/xlsx")

	// fungsi ini berguna untuk menghapus kembali file yang sudah diupload ke server
	os.Remove(fileDir)
	return
}

// Save: function for save task assignment (upload)
func Save(r createRequest) (sa *model.SalesAssignment, err error) {
	r.Code, err = util.GenerateCode(r.Code, "sales_assignment", 6)
	o := orm.NewOrm()
	o.Begin()

	sa = &model.SalesAssignment{
		SalesGroup: r.SalesGroup,
		Code:       r.Code,
		StartDate:  r.StartDateBatch,
		EndDate:    r.EndDateBatch,
		Status:     1,
	}
	if _, err = o.Insert(sa); err != nil {
		o.Rollback()
		return nil, err
	}

	for _, v := range r.AssignmentList {
		if err == nil {
			sla := &model.SalesAssignmentItem{
				SalesAssignment:     sa,
				Branch:              v.Branch,
				SalesPerson:         v.Salesperson,
				Task:                v.GlossaryTask.ValueInt,
				CustomerType:        v.GlossaryCustomerType.ValueInt,
				StartDate:           v.StartDate,
				EndDate:             v.EndDate,
				Status:              1,
				CustomerAcquisition: v.CustomerAcquisition,
				ObjectiveCodes:      v.ObjectiveCode,
			}

			if _, err := o.Insert(sla); err != nil {
				o.Rollback()
				return nil, err
			}

		}
	}

	if err = log.AuditLogByUser(r.Session.Staff, sa.ID, "sales assignment", "create", "Upload Task Assignment"); err != nil {
		o.Rollback()
		return nil, err
	}

	o.Commit()
	return nil, err
}

// Cancel : function to change data sales assignment status into 3
func Cancel(r cancelRequest) (sa *model.SalesAssignment, err error) {
	o := orm.NewOrm()
	o.Begin()

	sa = &model.SalesAssignment{
		ID:     r.ID,
		Status: 3,
	}

	if _, err = o.Update(sa, "Status"); err != nil {
		o.Rollback()
		return nil, err
	}

	for _, v := range r.SalesAssignmentItem {
		sai := &model.SalesAssignmentItem{
			ID:     v.ID,
			Status: 3,
		}

		if _, err = o.Update(sai, "Status"); err != nil {
			o.Rollback()
			return nil, err
		}
	}

	if err = log.AuditLogByUser(r.Session.Staff, sa.ID, "sales assignment", "cancel batch", ""); err != nil {
		o.Rollback()
		return nil, err
	}

	o.Commit()
	return sa, err
}

// CancelItem : function to change data sales assignment item status into 3
func CancelItem(r cancelItemRequest) (sai *model.SalesAssignmentItem, err error) {
	o := orm.NewOrm()
	o.Begin()

	sai = &model.SalesAssignmentItem{
		ID:     r.ID,
		Status: 3,
	}

	if _, err = o.Update(sai, "Status"); err != nil {
		o.Rollback()
		return nil, err
	}

	o.Commit()
	return sai, err
}
