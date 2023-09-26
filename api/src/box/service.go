//// Copyright 2021 PT. Eden Pangan Indonesia. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

package box

import (
	"fmt"
	"os"
	"time"

	"git.edenfarm.id/cuxs/orm"
	"github.com/tealeg/xlsx"

	"git.edenfarm.id/cuxs/env"
	"git.edenfarm.id/project-version2/api/log"
	"git.edenfarm.id/project-version2/api/util"
	"git.edenfarm.id/project-version2/datamodel/model"
)

// // Save : function to save data requested into database
// func Save(r createRequest) (tempBoxFridge []*model.BoxFridge, e error) {
// 	// docCode, e := util.GenerateDocCode("BF", r.Code, "sales_order")
// 	// if e != nil {
// 	// 	return nil, e
// 	// }

// 	//still in testing
// 	docCode := "TES"
// 	o := orm.NewOrm()
// 	o.Begin()

// 	r.BranchFridge.LastSeenAt = time.Now()
// 	if _, e := o.Update(r.BranchFridge); e != nil {
// 		o.Rollback()
// 		return nil, e
// 	}

// 	if len(r.ExistBoxFridges) > 0 {
// 		for _, boxFridge := range r.ExistBoxFridges {
// 			boxFridge.Code = docCode
// 			_, e = o.Update(boxFridge)
// 			if e != nil {
// 				o.Rollback()
// 				return nil, e
// 			}
// 		}
// 	}

// 	if len(r.NewBoxFridges) > 0 {
// 		_, e = o.InsertMulti(100, r.NewBoxFridges)
// 		if e != nil {
// 			o.Rollback()
// 			return nil, e
// 		}
// 	}

// 	//	if e = log.AuditLogByUser(r.Session.Staff, 1, "box_fridge", "create", ""); e != nil {

// 	//for testing purpose
// 	Staff := &model.Staff{ID: 222}
// 	if e = Staff.Read("ID"); e != nil {
// 		fmt.Println(e)
// 	}
// 	if e = log.AuditLogByUser(Staff, 1, "box_fridge", "create", ""); e != nil {
// 		o.Rollback()
// 		return nil, e
// 	}

// 	o.Commit()

// 	return tempBoxFridge, e
// }

// Update : function to update data requested into database
func UpdateFinish(r updateFinishRequest) (tempBoxFridge *model.BoxFridge, e error) {

	//still in testing
	o := orm.NewOrm()
	o.Begin()

	for _, boxItem := range r.ListBoxItem {
		boxItem.Status = 3
		boxItem.FinishedAt = time.Now()
		boxItem.FinishedBy = r.Session.Staff.ID
		boxItem.LastUpdatedAt = time.Now()
		boxItem.LastUpdatedBy = r.Session.Staff.ID
		_, e = o.Update(boxItem)
		if e != nil {
			o.Rollback()
			return nil, e
		}
	}

	if e = log.AuditLogByUser(r.Session.Staff, 1, "box_fridge", "update_finish", ""); e != nil {
		o.Rollback()
		return nil, e
	}
	o.Commit()

	return r.BoxFridge, e
}

// exportTemplateXls: download template export product price
func exportTemplateXls(date time.Time) (filePath string, err error) {
	var file *xlsx.File
	var sheet *xlsx.Sheet
	var row *xlsx.Row

	dir := env.GetString("EXPORT_DIRECTORY", "")

	filename := fmt.Sprintf("ProductBox%s_%s.xlsx", date.Format("2006-01-02"), util.GenerateRandomDoc(5))
	fileDir := fmt.Sprintf("%s/%s", dir, filename)

	file = xlsx.NewFile()
	if sheet, err = file.AddSheet("Sheet1"); err == nil {
		row = sheet.AddRow()
		row.SetHeight(20)
		//row.AddCell().Value = "No"
		row.AddCell().Value = "Customer_Name"
		row.AddCell().Value = "Rfid_Code"
		//row.AddCell().Value = "Size_Box(1,2,3,4)"
		row.AddCell().Value = "Product_Code"
		row.AddCell().Value = "Total_Weight"
		// row.AddCell().Value = "Unit_Price"
		// row.AddCell().Value = "Total_Price"
		row.AddCell().Value = "Note"

		boldStyle := xlsx.NewStyle()
		boldFont := xlsx.NewFont(10, "Liberation Sans")
		boldFont.Bold = true
		boldStyle.Font = *boldFont
		boldStyle.ApplyFont = true

		// looping to get column range 0-7. making BOLD font header
		for col := 0; col < 28; col++ {
			sheet.Cell(0, col).SetStyle(boldStyle)
		}

		err = file.Save(fileDir)
		filePath, err = util.UploadToS3(filename, fileDir, "application/xlsx")
		// fungsi ini berguna untuk menghapus kembali file yang sudah diupload ke server
		os.Remove(fileDir)
	}
	return
}

// Save : function to save data requested into database
func SaveTemplate(r createRequestTemplate) (tempBoxFridge []*model.BoxFridge, e error) {
	//still in testing
	//docCode := "TES"
	var newBoxItem []*model.BoxItem
	o := orm.NewOrm()
	o.Begin()

	if len(r.ExistProductBox) > 0 {
		for _, boxFridge := range r.ExistProductBox {
			var productBox *model.BoxItem
			productBox = &model.BoxItem{Box: boxFridge.Box,
				Product:       boxFridge.Product,
				CreatedAt:     time.Now(),
				CreatedBy:     r.Session.Staff.ID,
				LastUpdatedAt: time.Now(),
				// TotalPrice:    boxFridge.UnitPrice * boxFridge.TotalWeight,
				TotalWeight: boxFridge.TotalWeight,
				// UnitPrice:     boxFridge.UnitPrice,
				Note: boxFridge.Note, Status: 1}
			newBoxItem = append(newBoxItem, productBox)
		}
	}

	if len(r.NewProductBox) > 0 {
		for _, boxFridge := range r.NewProductBox {
			var box *model.Box
			var productBox *model.BoxItem
			box = &model.Box{
				//Code:      docCode,
				CreatedAt: time.Now(),
				CreatedBy: r.Session.Staff.ID,
				Note:      boxFridge.Note,
				Rfid:      boxFridge.Rfid,
				//Size:          int8(boxFridge.Size),
				Status:        1,
				LastUpdatedAt: time.Now()}

			idBox, e := o.Insert(box)
			if e != nil {
				o.Rollback()
				return nil, e
			}
			box.ID = idBox

			productBox = &model.BoxItem{Box: box,
				Product:       boxFridge.Product,
				CreatedAt:     time.Now(),
				CreatedBy:     r.Session.Staff.ID,
				LastUpdatedAt: time.Now(),
				// TotalPrice:    boxFridge.UnitPrice * boxFridge.TotalWeight,
				TotalWeight: boxFridge.TotalWeight,
				// UnitPrice:     boxFridge.UnitPrice,
				Note: boxFridge.Note, Status: 1}
			newBoxItem = append(newBoxItem, productBox)
		}
	}

	_, e = o.InsertMulti(100, newBoxItem)
	if e != nil {
		o.Rollback()
		return nil, e
	}

	if e = log.AuditLogByUser(r.Session.Staff, 1, "box_fridge", "create", ""); e != nil {
		o.Rollback()
		return nil, e
	}
	// //for testing purpose
	// Staff := &model.Staff{ID: 222}
	// if e = Staff.Read("ID"); e != nil {
	// 	fmt.Println(e)
	// }
	// if e = log.AuditLogByUser(Staff, 1, "box_fridge", "create", ""); e != nil {
	// 	o.Rollback()
	// 	return nil, e
	// }

	o.Commit()

	return tempBoxFridge, e
}
