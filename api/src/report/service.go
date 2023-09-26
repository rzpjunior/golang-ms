// Copyright 2020 PT. Eden Pangan Indonesia. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

package report

import (
	"fmt"
	"os"
	"time"

	"git.edenfarm.id/project-version2/api/log"
	"git.edenfarm.id/project-version2/api/util"
	"git.edenfarm.id/project-version2/datamodel/model"
	"github.com/tealeg/xlsx"
)

func getDeliveryOrderItemXls(date time.Time, r []*reportDeliveryOrderItem, warehouse *model.Warehouse) (filePath string, err error) {
	var file *xlsx.File
	var sheet *xlsx.Sheet
	var fileName string
	var row *xlsx.Row
	dir := util.ExportDirectory

	if warehouse != nil {
		fileName = fmt.Sprintf("ReportDeliveryOrderItem_%s_%s_%s.xlsx", date.Format("2006-01-02"), util.ReplaceSpace(warehouse.Name), util.GenerateRandomDoc(5))
	}

	fileName = fmt.Sprintf("ReportDeliveryOrderItem_%s_%s.xlsx", date.Format("2006-01-02"), util.GenerateRandomDoc(5))

	fileDir := fmt.Sprintf("%s/%s", dir, fileName)

	file = xlsx.NewFile()
	if sheet, err = file.AddSheet("Sheet1"); err == nil {
		row = sheet.AddRow()
		row.SetHeight(20)
		row.AddCell().Value = "No"
		row.AddCell().Value = "Order Code"
		row.AddCell().Value = "Delivery Code"
		row.AddCell().Value = "Product Code"
		row.AddCell().Value = "Product Name"
		row.AddCell().Value = "UOM"
		row.AddCell().Value = "Delivery Item Note"
		row.AddCell().Value = "Delivered Qty"
		row.AddCell().Value = "Received Qty"
		row.AddCell().Value = "Delivery Weight (Kg)"
		row.AddCell().Value = "Area"
		row.AddCell().Value = "Warehouse"
		row.AddCell().Value = "Order Delivery Date"
		row.AddCell().Value = "Delivery Date"
		row.AddCell().Value = "WRT"
		row.AddCell().Value = "Delivery Status"
		row.AddCell().Value = "Delivery Note"

		for index, i := range r {
			row = sheet.AddRow()
			row.AddCell().SetInt(index + 1)
			row.AddCell().Value = i.OrderCode                          // OrderCode
			row.AddCell().Value = i.DeliveryCode                       // DeliveryCode
			row.AddCell().Value = i.ProductCode                        // ProductCode
			row.AddCell().Value = i.ProductName                        // ProductName
			row.AddCell().Value = i.Uom                                // Uom
			row.AddCell().Value = i.DeliveryItemNote                   // DeliveryItemNote
			row.AddCell().SetFloatWithFormat(i.DeliveredQty, "0.00")   // DeliveredQty
			row.AddCell().SetFloatWithFormat(i.ReceivedQty, "0.00")    // ReceivedQty
			row.AddCell().SetFloatWithFormat(i.DeliveryWeight, "0.00") // DeliveryWeight
			row.AddCell().Value = i.Area                               //Area
			row.AddCell().Value = i.Warehouse                          //Warehouse
			row.AddCell().Value = i.OrderDeliveryDate                  // OrderDeliveryDate
			row.AddCell().Value = i.DeliveryDate                       // DeliveryDate
			row.AddCell().Value = i.Wrt                                // Wrt
			row.AddCell().Value = i.DeliveryStatus                     // DeliveryStatus
			row.AddCell().Value = i.DeliveryNote                       // DeliveryNote
		}

	}
	err = file.Save(fileDir)
	filePath, err = util.UploadToS3(fileName, fileDir, "application/xlsx")

	// fungsi ini berguna untuk menghapus kembali file yang sudah diupload ke server
	os.Remove(fileDir)
	return
}

func DownloadPackingReportXls(date time.Time, r []*reportPackingOrder, wh *model.Warehouse) (filePath string, err error) {
	var file *xlsx.File
	var sheet *xlsx.Sheet
	var row *xlsx.Row
	dir := util.ExportDirectory
	filename := fmt.Sprintf("ReportPackingOrder_%s_%s_%s.xlsx", util.ReplaceUnderscore(wh.Name), date.Format("2006-01-02"), util.GenerateRandomDoc(5))

	fileDir := fmt.Sprintf("%s/%s", dir, filename)
	file = xlsx.NewFile()
	if sheet, err = file.AddSheet("Sheet1"); err == nil {
		row = sheet.AddRow()
		row.SetHeight(20)
		row.AddCell().Value = "No"
		row.AddCell().Value = "Packing_Date"
		row.AddCell().Value = "Product_Name"
		row.AddCell().Value = "UOM"
		row.AddCell().Value = "Total_Order"
		row.AddCell().Value = "Total_Packing"
		row.AddCell().Value = "Total_Weight"
		row.AddCell().Value = "Helper_Code"
		row.AddCell().Value = "Helper_Name"

		for i, v := range r {
			row = sheet.AddRow()
			row.AddCell().SetInt(i + 1)
			row.AddCell().Value = v.DeliveryDate                                // Delivery Date
			row.AddCell().Value = v.ProductName                                 // Product Name
			row.AddCell().Value = v.Uom                                         // Uom
			row.AddCell().SetFloatWithFormat(float64(v.TotalOrder), "0.00")     // Total Order
			row.AddCell().SetFloatWithFormat(float64(v.SubtotalPack), "0.00")   // Total Pack
			row.AddCell().SetFloatWithFormat(float64(v.SubtotalWeight), "0.00") // Total Weigth
			row.AddCell().Value = v.HelperCode                                  // Helper Code
			row.AddCell().Value = v.HelperName                                  // Helper Name

		}
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
	}
	filePath, err = util.UploadToS3(filename, fileDir, "application/xlsx")
	// fungsi ini berguna untuk menghapus kembali file yang sudah diupload ke server
	os.Remove(fileDir)
	return
}

// GetPricingInboundItem : function to create excel file of pricing inbound item report
func GetPricingInboundItem(date string, data []*reportPricingInboundItem, area *model.Area, staff *model.Staff, note string) (filePath string, err error) {
	var file *xlsx.File
	var sheet *xlsx.Sheet
	var row *xlsx.Row

	dir := util.ExportDirectory
	filename := fmt.Sprintf("ReportPricingInboundItem%s_%s_%s.xlsx", time.Now().Format("2006-01-02"), util.ReplaceSpace(area.Name), util.GenerateRandomDoc(5))

	fileDir := fmt.Sprintf("%s/%s", dir, filename)

	file = xlsx.NewFile()
	sheet, err = file.AddSheet("Sheet1")
	if err != nil {
		return "", err
	}

	row = sheet.AddRow()
	row.SetHeight(20)
	row.AddCell().Value = "No"
	row.AddCell().Value = "Inbound Code"
	row.AddCell().Value = "Supplier Code"
	row.AddCell().Value = "Supplier Name"
	row.AddCell().Value = "Warehouse Origin"
	row.AddCell().Value = "Warehouse Destination"
	row.AddCell().Value = "Area"
	row.AddCell().Value = "Order_Date"
	row.AddCell().Value = "Eta Date"
	row.AddCell().Value = "Ata Date"
	row.AddCell().Value = "Product Code"
	row.AddCell().Value = "Product Name"
	row.AddCell().Value = "UOM"
	row.AddCell().Value = "Request Qty"
	row.AddCell().Value = "Delivered Qty"
	row.AddCell().Value = "Receive Qty"
	row.AddCell().Value = "Invoice Qty"
	row.AddCell().Value = "Taxability"
	row.AddCell().Value = "Tax Percentage"
	row.AddCell().Value = "Unit Price"
	row.AddCell().Value = "Inbound Status"

	for i, v := range data {
		row = sheet.AddRow()
		row.AddCell().SetInt(i + 1)
		row.AddCell().Value = v.InboundCode                      // Inbound Code
		row.AddCell().Value = v.SupplierCode                     //	Supplier Code
		row.AddCell().Value = v.SupplierName                     // Supplier Name
		row.AddCell().Value = v.WarehouseOrigin                  // warehouse Origin
		row.AddCell().Value = v.WarehouseDestination             // warehouse Destination
		row.AddCell().Value = v.Area                             // Area Code
		row.AddCell().Value = v.OrderDate                        // Order Date
		row.AddCell().Value = v.EtaDate                          // Eta Date
		row.AddCell().Value = v.AtaDate                          // Ata Date
		row.AddCell().Value = v.ProductCode                      // Product Code
		row.AddCell().Value = v.ProductName                      // Product Name
		row.AddCell().Value = v.Uom                              // Uom
		row.AddCell().SetFloatWithFormat(v.RequestQty, "0.00")   // Request Quantity
		row.AddCell().SetFloatWithFormat(v.DeliveredQty, "0.00") // Delivered Quantity
		row.AddCell().SetFloatWithFormat(v.ReceiveQty, "0.00")   // Receive Quantity
		row.AddCell().SetFloatWithFormat(v.InvoiceQty, "0.00")   // Invoice Quantity
		if v.Taxability == 1 {
			row.AddCell().Value = "yes" // Taxability
		} else {
			row.AddCell().Value = "no"
		}
		row.AddCell().SetFloatWithFormat(v.TaxPercentage, "0.00") // Tax Percentage
		row.AddCell().SetFloatWithFormat(v.UnitPrice, "0.00")     // Unit Price
		row.AddCell().Value = v.InboundStatus                     // Inbound Status
	}

	err = file.Save(fileDir)
	if err != nil {
		return "", err
	}

	filePath, err = util.UploadToS3(filename, fileDir, "application/xlsx")
	if err != nil {
		return "", err
	}

	// fungsi ini berguna untuk menghapus kembali file yang sudah diupload ke server
	os.Remove(fileDir)

	if err = log.AuditLogByUser(staff, 0, "report_pricing_inbound_item", "Download", note); err != nil {
		return "", err
	}

	return
}

func DownloadPriceChangeHistoryXls(date time.Time, r []*reportPriceChangeHistory, staff *model.Staff, note string) (filePath string, err error) {
	var file *xlsx.File
	var sheet *xlsx.Sheet
	var row *xlsx.Row
	dir := util.ExportDirectory
	filename := fmt.Sprintf("ReportPriceChangeHistory_%s_%s.xlsx", date.Format("2006-01-02"), util.GenerateRandomDoc(5))

	fileDir := fmt.Sprintf("%s/%s", dir, filename)
	file = xlsx.NewFile()
	if sheet, err = file.AddSheet("Sheet1"); err == nil {
		row = sheet.AddRow()
		row.SetHeight(20)
		row.AddCell().Value = "No"
		row.AddCell().Value = "Created At"
		row.AddCell().Value = "Price Set"
		row.AddCell().Value = "Product Name"
		row.AddCell().Value = "Unit Price"
		row.AddCell().Value = "Created By"

		for i, v := range r {
			row = sheet.AddRow()
			row.AddCell().SetInt(i + 1)
			row.AddCell().Value = v.CreatedAt                              // Created At
			row.AddCell().Value = v.PriceSet                               // Price Set
			row.AddCell().Value = v.ProductName                            // Product Name
			row.AddCell().SetFloatWithFormat(float64(v.UnitPrice), "0.00") // Unit Price
			row.AddCell().Value = v.CreatedBy                              // Created By
		}
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
	}
	filePath, err = util.UploadToS3(filename, fileDir, "application/xlsx")
	// fungsi ini berguna untuk menghapus kembali file yang sudah diupload ke server
	os.Remove(fileDir)

	if err = log.AuditLogByUser(staff, 0, "report_price_change_history", "Download", note); err != nil {
		return "", err
	}

	return
}
