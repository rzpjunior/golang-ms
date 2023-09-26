// Copyright 2020 PT. Eden Pangan Indonesia. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

package fridge

import (
	"fmt"
	"os"
	"time"

	"github.com/tealeg/xlsx"

	"git.edenfarm.id/project-version2/api/util"
	"git.edenfarm.id/project-version2/datamodel/model"
)

func getSoldProductXls(date string, data []*reportSoldProductFridge, warehouse *model.Warehouse) (filePath string, err error) {
	var file *xlsx.File
	var sheet *xlsx.Sheet
	var row *xlsx.Row

	dir := util.ExportDirectory
	filename := fmt.Sprintf("ReportSoldProductFridge_%s_%s_%s.xlsx", time.Now().Format("2006-01-02"), warehouse.Name, util.GenerateRandomDoc(5))

	fileDir := fmt.Sprintf("%s/%s", dir, filename)

	file = xlsx.NewFile()
	if sheet, err = file.AddSheet("Sheet1"); err == nil {
		row = sheet.AddRow()
		row.SetHeight(20)
		row.AddCell().Value = "No"
		row.AddCell().Value = "Created Date"
		row.AddCell().Value = "Sold Date"
		row.AddCell().Value = "Merchant Name"
		row.AddCell().Value = "Branch Name"
		row.AddCell().Value = "Warehouse Name"
		row.AddCell().Value = "Product Name"
		row.AddCell().Value = "Qty"
		row.AddCell().Value = "UOM"

		for i, v := range data {

			row = sheet.AddRow()
			row.AddCell().SetInt(i + 1)
			row.AddCell().Value = v.CreatedAt     // SO Code
			row.AddCell().Value = v.SoldDate      // Customer Code
			row.AddCell().Value = v.MerchantName  // Customer Tag
			row.AddCell().Value = v.BranchName    // Customer Name
			row.AddCell().Value = v.WarehouseName // Customer Phone Number
			row.AddCell().Value = v.ProductName   // Recipient Name
			row.AddCell().Value = v.TotalWeight   // Recipient Phone Number
			row.AddCell().Value = v.UOMName       // Shipping Address
		}

	}
	err = file.Save(fileDir)
	filePath, err = util.UploadToS3(filename, fileDir, "application/xlsx")

	// fungsi ini berguna untuk menghapus kembali file yang sudah diupload ke server
	os.Remove(fileDir)
	return
}

func getAllProductXls(date string, data []*reportAllProductFridge, warehouse *model.Warehouse) (filePath string, err error) {
	var file *xlsx.File
	var sheet *xlsx.Sheet
	var row *xlsx.Row

	dir := util.ExportDirectory
	filename := fmt.Sprintf("ReportAllProductFridge_%s_%s_%s.xlsx", time.Now().Format("2006-01-02"), warehouse.Name, util.GenerateRandomDoc(5))

	fileDir := fmt.Sprintf("%s/%s", dir, filename)

	file = xlsx.NewFile()
	if sheet, err = file.AddSheet("Sheet1"); err == nil {
		row = sheet.AddRow()
		row.SetHeight(20)
		row.AddCell().Value = "No"
		row.AddCell().Value = "Created Date"
		row.AddCell().Value = "Processed Date"
		row.AddCell().Value = "Finished Date"
		row.AddCell().Value = "Merchant Name"
		row.AddCell().Value = "Branch Name"
		row.AddCell().Value = "Warehouse Name"
		row.AddCell().Value = "Product Code"
		row.AddCell().Value = "Product Name"
		row.AddCell().Value = "Unit Price"
		row.AddCell().Value = "Qty"
		row.AddCell().Value = "UOM"
		row.AddCell().Value = "Total Price"
		//row.AddCell().Value = "Box Item Status"
		//row.AddCell().Value = "Box Fridge Status"
		row.AddCell().Value = "Status"
		row.AddCell().Value = "Waste Image URL"

		for i, v := range data {
			row = sheet.AddRow()
			if v.BoxFridgeStatus == 1 {
				v.Status = ""
				v.ProcessedDate = ""
			}
			row.AddCell().SetInt(i + 1)
			row.AddCell().Value = v.CreatedAt     // SO Code
			row.AddCell().Value = v.ProcessedDate // Customer Code
			row.AddCell().Value = v.FinishedAt    // Customer Code
			row.AddCell().Value = v.MerchantName  // Customer Tag
			row.AddCell().Value = v.BranchName    // Customer Name
			row.AddCell().Value = v.WarehouseName // Customer Phone Number
			row.AddCell().Value = v.ProductCode   // Recipient Name
			row.AddCell().Value = v.ProductName   // Recipient Name
			row.AddCell().Value = v.UnitPrice     // Recipient Name
			row.AddCell().Value = v.TotalWeight   // Recipient Phone Number
			row.AddCell().Value = v.UOMName
			row.AddCell().Value = v.TotalPrice // Shipping Address
			//	row.AddCell().Value = v.BoxItemStatus   // Shipping Address
			//	row.AddCell().Value = v.BoxFridgeStatus // Shipping Address

			if v.BoxItemStatus == 3 {
				row.AddCell().Value = "finished " + v.Status // Shipping Address
			} else {
				row.AddCell().Value = v.Status // Shipping Address
			}
			row.AddCell().Value = v.ImageURL // Shipping Address
		}

	}
	err = file.Save(fileDir)
	filePath, err = util.UploadToS3(filename, fileDir, "application/xlsx")

	// fungsi ini berguna untuk menghapus kembali file yang sudah diupload ke server
	os.Remove(fileDir)
	return
}
