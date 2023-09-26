// Copyright 2020 PT. Eden Pangan Indonesia. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

package sms

import (
	"fmt"
	"os"
	"time"

	"github.com/tealeg/xlsx"

	"git.edenfarm.id/project-version2/api/log"
	"git.edenfarm.id/project-version2/api/util"
	"git.edenfarm.id/project-version2/datamodel/model"
)

// GetPurchaseOrderXls : function to create excel file of purchase order report
func GetPurchaseOrderXls(date string, data []*reportPurchaseOrder, area *model.Area) (filePath string, err error) {
	var file *xlsx.File
	var sheet *xlsx.Sheet
	var row *xlsx.Row

	dir := util.ExportDirectory
	filename := fmt.Sprintf("ReportPurchaseOrder_%s_%s_%s.xlsx", time.Now().Format("2006-01-02"), util.ReplaceSpace(area.Name), util.GenerateRandomDoc(5))

	fileDir := fmt.Sprintf("%s/%s", dir, filename)
	file = xlsx.NewFile()
	if sheet, err = file.AddSheet("Sheet1"); err == nil {
		row = sheet.AddRow()
		row.SetHeight(20)
		row.AddCell().Value = "No"
		row.AddCell().Value = "Supplier Code"
		row.AddCell().Value = "Supplier Name"
		row.AddCell().Value = "Warehouse"
		row.AddCell().Value = "Purchase Order Code"
		row.AddCell().Value = "Order Date"
		row.AddCell().Value = "ETA Date"
		row.AddCell().Value = "ETA Time"
		row.AddCell().Value = "Order Status"
		row.AddCell().Value = "Order Payment Term"
		row.AddCell().Value = "Order Delivery Fee"
		row.AddCell().Value = "Order Tax Amount"
		row.AddCell().Value = "Grand Total"
		row.AddCell().Value = "Order Note"
		row.AddCell().Value = "Order Total"
		row.AddCell().Value = "Supplier Badge"
		row.AddCell().Value = "Goods Receipt"

		for i, v := range data {
			row = sheet.AddRow()
			row.AddCell().SetInt(i + 1)
			row.AddCell().Value = v.SupplierCode                    // Supplier Code
			row.AddCell().Value = v.SupplierName                    // Supplier Name
			row.AddCell().Value = v.WarehouseName                   // Warehouse Name
			row.AddCell().Value = v.OrderCode                       // Purchase Order Code
			row.AddCell().Value = v.OrderDate                       // Order Date
			row.AddCell().Value = v.EtaDate                         // Eta Date
			row.AddCell().Value = v.EtaTime                         // Eta Time
			row.AddCell().Value = v.OrderStatus                     // Order Status
			row.AddCell().Value = v.PaymentTerm                     // Order Payment Term
			row.AddCell().SetFloatWithFormat(v.DeliveryFee, "0.00") // Delivery Fee
			row.AddCell().SetFloatWithFormat(v.TaxAmount, "0.00")   // Tax Amount
			row.AddCell().SetFloatWithFormat(v.GrandTotal, "0.00")  // Grand Total
			row.AddCell().Value = v.OrderNote                       // Order Note
			row.AddCell().SetFloatWithFormat(v.TotalPrice, "0.00")  // Total Price
			row.AddCell().Value = v.SupplierBadge                   // Supplier Badge
			row.AddCell().Value = v.GoodsReceipt                    // Goods Receipt
		}

	}
	err = file.Save(fileDir)
	filePath, err = util.UploadToS3(filename, fileDir, "application/xlsx")

	// fungsi ini berguna untuk menghapus kembali file yang sudah diupload ke server
	os.Remove(fileDir)
	return
}

// GetPurchaseOrderItemXls : function to create excel file of purchase order item report
func GetPurchaseOrderItemXls(date string, data []*reportPurchaseOrderItem, area *model.Area) (filePath string, err error) {
	var file *xlsx.File
	var sheet *xlsx.Sheet
	var row *xlsx.Row

	dir := util.ExportDirectory
	filename := fmt.Sprintf("ReportPurchaseOrderItem_%s_%s_%s.xlsx", time.Now().Format("2006-01-02"), util.ReplaceSpace(area.Name), util.GenerateRandomDoc(5))

	fileDir := fmt.Sprintf("%s/%s", dir, filename)

	file = xlsx.NewFile()
	if sheet, err = file.AddSheet("Sheet1"); err == nil {
		row = sheet.AddRow()
		row.SetHeight(20)
		row.AddCell().Value = "No"
		row.AddCell().Value = "Purchase Order Code"
		row.AddCell().Value = "Product Code"
		row.AddCell().Value = "Product Name"
		row.AddCell().Value = "UOM"
		row.AddCell().Value = "Taxable Item"
		row.AddCell().Value = "Order Item Note"
		row.AddCell().Value = "Ordered Qty"
		row.AddCell().Value = "Order Unit Price"
		row.AddCell().Value = "Include Tax"
		row.AddCell().Value = "Tax Percentage"
		row.AddCell().Value = "Unit Price Tax"
		row.AddCell().Value = "Tax Amount"
		row.AddCell().Value = "Subtotal"
		row.AddCell().Value = "Total Weight (Kg)"
		row.AddCell().Value = "Area"
		row.AddCell().Value = "Warehouse"
		row.AddCell().Value = "Supplier Code"
		row.AddCell().Value = "Supplier Name"
		row.AddCell().Value = "Order Date"
		row.AddCell().Value = "Eta Date"
		row.AddCell().Value = "Invoiced Qty"
		row.AddCell().Value = "Purchase Qty"

		for i, v := range data {

			row = sheet.AddRow()
			row.AddCell().SetInt(i + 1)
			row.AddCell().Value = v.OrderCode                              // Purchase Order Code
			row.AddCell().Value = v.ProductCode                            // Product Code
			row.AddCell().Value = v.ProductName                            // Product Name
			row.AddCell().Value = v.UOM                                    // UOM
			row.AddCell().Value = v.TaxableStr                             // Is Taxable Item?
			row.AddCell().Value = v.OrderItemNote                          // Order Item Note
			row.AddCell().SetFloatWithFormat(v.OrderedQty, "0.00")         // Ordered Qty
			row.AddCell().SetFloatWithFormat(v.OrderUnitPrice, "0.00")     // Order Unit Price
			row.AddCell().Value = v.IncludeTaxStr                          // Is Include Tax?
			row.AddCell().SetFloatWithFormat(v.OrderTaxPercentage, "0.00") // Order Unit Tax Percentage
			row.AddCell().SetFloatWithFormat(v.OrderUnitPriceTax, "0.00")  // Order Unit Price Tax
			row.AddCell().SetFloatWithFormat(v.OrderTaxAmount, "0.00")     // Order Tax Amount
			row.AddCell().SetFloatWithFormat(v.Subtotal, "0.00")           // Subtotal
			row.AddCell().SetFloatWithFormat(v.TotalWeight, "0.00")        // Total Weight
			row.AddCell().Value = v.AreaName                               // Area Name
			row.AddCell().Value = v.WarehouseName                          // Warehouse Name
			row.AddCell().Value = v.SupplierCode                           // Supplier Code
			row.AddCell().Value = v.SupplierName                           // Supplier Name
			row.AddCell().Value = v.OrderDate                              // Order Date
			row.AddCell().Value = v.EtaDate                                // Eta Date
			row.AddCell().SetFloatWithFormat(v.InvoicedQty, "0.00")        // Invoiced Qty
			row.AddCell().SetFloatWithFormat(v.PurchaseQty, "0.00")        // Purchased Qty
		}

	}
	err = file.Save(fileDir)
	filePath, err = util.UploadToS3(filename, fileDir, "application/xlsx")

	// fungsi ini berguna untuk menghapus kembali file yang sudah diupload ke server
	os.Remove(fileDir)
	return
}

// GetPurchaseInvoiceXls : function to create excel file of purchase invoice report
func GetPurchaseInvoiceXls(date string, data []*reportPurchaseInvoice, area *model.Area, staff *model.Staff, note string) (filePath string, err error) {
	var (
		file  *xlsx.File
		sheet *xlsx.Sheet
		row   *xlsx.Row
	)

	dir := util.ExportDirectory
	filename := fmt.Sprintf("ReportPurchaseInvoice_%s_%s_%s.xlsx", time.Now().Format("2006-01-02"), util.ReplaceSpace(area.Name), util.GenerateRandomDoc(5))

	fileDir := fmt.Sprintf("%s/%s", dir, filename)

	file = xlsx.NewFile()
	if sheet, err = file.AddSheet("Sheet1"); err == nil {
		row = sheet.AddRow()
		row.SetHeight(20)
		row.AddCell().Value = "No"
		row.AddCell().Value = "Area"
		row.AddCell().Value = "Warehouse"
		row.AddCell().Value = "Purchase Order Code"
		row.AddCell().Value = "Order Date"
		row.AddCell().Value = "Order ETA Date"
		row.AddCell().Value = "Supplier Code"
		row.AddCell().Value = "Supplier Name"
		row.AddCell().Value = "Total Order"
		row.AddCell().Value = "Invoice Code"
		row.AddCell().Value = "Invoice Date"
		row.AddCell().Value = "Invoice Due Date"
		row.AddCell().Value = "Invoice Status"
		row.AddCell().Value = "Invoice Note"
		row.AddCell().Value = "Delivery Fee"
		row.AddCell().Value = "Invoice Amount"
		row.AddCell().Value = "Total Tax Amount"
		row.AddCell().Value = "Total Invoice"
		row.AddCell().Value = "Adjustment Amount"
		row.AddCell().Value = "Adjustment Note"
		row.AddCell().Value = "Created At"
		row.AddCell().Value = "Created By"
		row.AddCell().Value = "Supplier Type"
		row.AddCell().Value = "Payment Term"
		row.AddCell().Value = "Total Payment"
		row.AddCell().Value = "ATA Date"

		for i, v := range data {
			ataDate := ""
			if v.AtaDate.Format("02-01-2006") != "01-01-0001" {
				ataDate = v.AtaDate.Format("02-01-2006")
			}

			row = sheet.AddRow()
			row.AddCell().SetInt(i + 1)
			row.AddCell().Value = v.AreaName                             // Area Name
			row.AddCell().Value = v.WarehouseName                        // Warehouse Name
			row.AddCell().Value = v.OrderCode                            // Order Code
			row.AddCell().Value = v.OrderDate                            // Order Date
			row.AddCell().Value = v.EtaDate                              // Eta Date
			row.AddCell().Value = v.SupplierCode                         // Supplier Code
			row.AddCell().Value = v.SupplierName                         // Supplier Name
			row.AddCell().SetFloatWithFormat(v.TotalOrder, "0.00")       // Total Order
			row.AddCell().Value = v.InvoiceCode                          // Invoice Code
			row.AddCell().Value = v.InvoiceDate                          // Invoice Date
			row.AddCell().Value = v.InvoiceDueDate                       // Invoice Due Date
			row.AddCell().Value = v.InvoiceStatus                        // Invoice Status
			row.AddCell().Value = v.InvoiceNote                          // Invoice Note
			row.AddCell().SetFloatWithFormat(v.DeliveryFee, "0.00")      // Delivery Fee
			row.AddCell().SetFloatWithFormat(v.InvoiceAmount, "0.00")    // Invoice Amount
			row.AddCell().SetFloatWithFormat(v.TaxAmount, "0.00")        // Tax Amount
			row.AddCell().SetFloatWithFormat(v.TotalInvoice, "0.00")     // Total Invoice
			row.AddCell().SetFloatWithFormat(v.AdjustmentAmount, "0.00") // Adjustment Amount
			row.AddCell().Value = v.AdjustmentNote                       // Adjustment Note
			row.AddCell().Value = v.CreatedAt                            // Created At
			row.AddCell().Value = v.CreatedBy                            // Created By
			row.AddCell().Value = v.SupplierType                         // Supplier Type
			row.AddCell().Value = v.PaymentTerm                          // Payment Term
			row.AddCell().SetFloatWithFormat(v.TotalPayment, "0.00")     // Total Payment
			row.AddCell().Value = ataDate                                // ATA Date
		}
	}
	err = file.Save(fileDir)
	filePath, err = util.UploadToS3(filename, fileDir, "application/xlsx")
	// fungsi ini berguna untuk menghapus kembali file yang sudah diupload ke server
	os.Remove(fileDir)

	if err = log.AuditLogByUser(staff, 0, "report_purchase_invoice", "Download", note); err != nil {
		return "", err
	}

	return
}

// GetPurchasePaymentXls : function to create excel file of purchase payment report
func GetPurchasePaymentXls(date string, data []*reportPurchasePayment, area *model.Area, staff *model.Staff, note string) (filePath string, err error) {
	var file *xlsx.File
	var sheet *xlsx.Sheet
	var row *xlsx.Row

	dir := util.ExportDirectory
	filename := fmt.Sprintf("ReportPurchasePayment_%s_%s_%s.xlsx", time.Now().Format("2006-01-02"), util.ReplaceSpace(area.Name), util.GenerateRandomDoc(5))

	fileDir := fmt.Sprintf("%s/%s", dir, filename)

	file = xlsx.NewFile()
	if sheet, err = file.AddSheet("Sheet1"); err == nil {
		row = sheet.AddRow()
		row.SetHeight(20)
		row.AddCell().Value = "No"
		row.AddCell().Value = "Area"
		row.AddCell().Value = "Supplier Code"
		row.AddCell().Value = "Supplier Name"
		row.AddCell().Value = "Payment Code"
		row.AddCell().Value = "Payment Date"
		row.AddCell().Value = "Payment Amount"
		row.AddCell().Value = "Payment Method"
		row.AddCell().Value = "Payment Status"
		row.AddCell().Value = "Invoice Code"
		row.AddCell().Value = "Total Invoice"
		row.AddCell().Value = "Bank Payment Voucher Number"
		row.AddCell().Value = "Created At"
		row.AddCell().Value = "Created By"

		for i, v := range data {
			row = sheet.AddRow()
			row.AddCell().SetInt(i + 1)
			row.AddCell().Value = v.AreaName                          // Area Name
			row.AddCell().Value = v.SupplierCode                      // Supplier Code
			row.AddCell().Value = v.SupplierName                      // Supplier Name
			row.AddCell().Value = v.PaymentCode                       // Payment Code
			row.AddCell().Value = v.PaymentDate                       // Payment Date
			row.AddCell().SetFloatWithFormat(v.PaymentAmount, "0.00") // Payment Amount
			row.AddCell().Value = v.PaymentMethod                     // Payment Method
			row.AddCell().Value = v.PaymentStatus                     // Payment Status
			row.AddCell().Value = v.InvoiceCode                       // Invoice Code
			row.AddCell().SetFloatWithFormat(v.TotalInvoice, "0.00")  // Total Invoice
			row.AddCell().Value = v.PaymentNumber                     // Bank Payment Voucher Number
			row.AddCell().Value = v.CreatedAt                         // Created At
			row.AddCell().Value = v.CreatedBy                         // Created By
		}

	}
	err = file.Save(fileDir)
	filePath, err = util.UploadToS3(filename, fileDir, "application/xlsx")

	// fungsi ini berguna untuk menghapus kembali file yang sudah diupload ke server
	os.Remove(fileDir)

	if err = log.AuditLogByUser(staff, 0, "report_purchase_payment", "Download", note); err != nil {
		return "", err
	}

	return
}

// GetPurchaseInvoiceItemXls : function to create excel file of purchase invoice item report
func GetPurchaseInvoiceItemXls(date string, data []*reportPurchaseInvoiceItem, area *model.Area, staff *model.Staff, note string) (filePath string, err error) {
	var file *xlsx.File
	var sheet *xlsx.Sheet
	var row *xlsx.Row

	dir := util.ExportDirectory
	filename := fmt.Sprintf("ReportPurchaseInvoiceItem_%s_%s_%s.xlsx", time.Now().Format("2006-01-02"), util.ReplaceSpace(area.Name), util.GenerateRandomDoc(5))

	fileDir := fmt.Sprintf("%s/%s", dir, filename)

	file = xlsx.NewFile()
	if sheet, err = file.AddSheet("Sheet1"); err == nil {
		row = sheet.AddRow()
		row.SetHeight(20)
		row.AddCell().Value = "No"
		row.AddCell().Value = "Supplier Name"
		row.AddCell().Value = "Warehouse Name"
		row.AddCell().Value = "Area"
		row.AddCell().Value = "Order Code"
		row.AddCell().Value = "Order Status"
		row.AddCell().Value = "Invoice Code"
		row.AddCell().Value = "Invoice Status"
		row.AddCell().Value = "GR Code"
		row.AddCell().Value = "GR Status"
		row.AddCell().Value = "Eta Date"
		row.AddCell().Value = "Product Code"
		row.AddCell().Value = "Product Name"
		row.AddCell().Value = "UOM"
		row.AddCell().Value = "Taxable Item"
		row.AddCell().Value = "Include Tax"
		row.AddCell().Value = "Unit Price"
		row.AddCell().Value = "Tax Percentage"
		row.AddCell().Value = "Unit Price Tax"
		row.AddCell().Value = "Tax Amount"
		row.AddCell().Value = "Order Qty"
		row.AddCell().Value = "Delivery Qty"
		row.AddCell().Value = "Received Qty"
		row.AddCell().Value = "Invoice Qty"
		row.AddCell().Value = "Reject Qty"
		row.AddCell().Value = "Delivery Fee"
		row.AddCell().Value = "Total Invoice"

		for i, v := range data {
			row = sheet.AddRow()
			row.AddCell().SetInt(i + 1)
			row.AddCell().Value = v.SupplierName                      // Supplier Name
			row.AddCell().Value = v.WarehouseName                     // Warehouse Name
			row.AddCell().Value = v.Area                              // Area
			row.AddCell().Value = v.OrderCode                         // Order Code
			row.AddCell().Value = v.OrderStatus                       // Order Status
			row.AddCell().Value = v.InvoiceCode                       // Invoice Code
			row.AddCell().Value = v.InvoiceStatus                     // Invoice Status
			row.AddCell().Value = v.GRCode                            // GR Code
			row.AddCell().Value = v.GRStatus                          // GR Status
			row.AddCell().Value = v.EtaDate                           // Eta Date
			row.AddCell().Value = v.ProductCode                       // Product Code
			row.AddCell().Value = v.ProductName                       // Product Name
			row.AddCell().Value = v.UOM                               // UOM
			row.AddCell().Value = v.TaxableStr                        // Is Taxable Item?
			row.AddCell().Value = v.IncludeTaxStr                     // Is Include Tax?
			row.AddCell().SetFloatWithFormat(v.UnitPrice, "0.00")     // Unit Price
			row.AddCell().SetFloatWithFormat(v.TaxPercentage, "0.00") // Unit Tax Percentage
			row.AddCell().SetFloatWithFormat(v.UnitPriceTax, "0.00")  // Unit Price Tax
			row.AddCell().SetFloatWithFormat(v.TaxAmount, "0.00")     // Tax Amount
			row.AddCell().SetFloatWithFormat(v.OrderQty, "0.00")      // Order Qty
			row.AddCell().SetFloatWithFormat(v.DeliveredQty, "0.00")  // Delivered Qty
			row.AddCell().SetFloatWithFormat(v.ReceivedQty, "0.00")   // Received Qty
			row.AddCell().SetFloatWithFormat(v.InvoiceQty, "0.00")    // Invoice Qty
			row.AddCell().SetFloatWithFormat(v.RejectQty, "0.00")     // Reject Qty
			row.AddCell().SetFloatWithFormat(v.DeliveryFee, "0.00")   // Delivery Fee
			row.AddCell().SetFloatWithFormat(v.TotalInvoice, "0.00")  // Total Invoice
		}

	}
	err = file.Save(fileDir)
	filePath, err = util.UploadToS3(filename, fileDir, "application/xlsx")
	// fungsi ini berguna untuk menghapus kembali file yang sudah diupload ke server
	os.Remove(fileDir)

	if err = log.AuditLogByUser(staff, 0, "report_purchase_invoice_item", "Download", note); err != nil {
		return "", err
	}

	return
}

// GetCogsXls : function to create excel file of cogs report
func GetCogsXls(date string, data []*reportCogs, area *model.Area, staff *model.Staff, note string) (filePath string, err error) {
	var file *xlsx.File
	var sheet *xlsx.Sheet
	var row *xlsx.Row

	dir := util.ExportDirectory
	filename := fmt.Sprintf("ReportCogs_%s_%s_%s.xlsx", time.Now().Format("2006-01-02"), util.ReplaceSpace(area.Name), util.GenerateRandomDoc(5))

	fileDir := fmt.Sprintf("%s/%s", dir, filename)

	file = xlsx.NewFile()
	if sheet, err = file.AddSheet("Sheet1"); err == nil {
		row = sheet.AddRow()
		row.SetHeight(20)
		row.AddCell().Value = "No"
		row.AddCell().Value = "Area"
		row.AddCell().Value = "Warehouse"
		row.AddCell().Value = "PO ETA Date"
		row.AddCell().Value = "Product Code"
		row.AddCell().Value = "Product Name"
		row.AddCell().Value = "UOM"
		row.AddCell().Value = "Average Purchase Price"

		for i, v := range data {
			row = sheet.AddRow()
			row.AddCell().SetInt(i + 1)
			row.AddCell().Value = v.AreaName                     // Area Name
			row.AddCell().Value = v.WarehouseName                // Warehouse Name
			row.AddCell().Value = v.EtaDate                      // PO ETA Date
			row.AddCell().Value = v.ProductCode                  // Product Code
			row.AddCell().Value = v.ProductName                  // Product Name
			row.AddCell().Value = v.UOM                          // UOM
			row.AddCell().SetFloatWithFormat(v.AvgPrice, "0.00") // Average Purchase Price
		}

	}
	err = file.Save(fileDir)
	filePath, err = util.UploadToS3(filename, fileDir, "application/xlsx")
	// fungsi ini berguna untuk menghapus kembali file yang sudah diupload ke server
	os.Remove(fileDir)

	if err = log.AuditLogByUser(staff, 0, "report_cogs", "Download", note); err != nil {
		return "", err
	}

	return
}

// GetPriceComparisonXls : function to create excel file of price comparison report
func GetPriceComparisonXls(date string, data []*reportPriceComparison, area *model.Area) (filePath string, err error) {
	var file *xlsx.File
	var sheet *xlsx.Sheet
	var row *xlsx.Row

	dir := util.ExportDirectory
	filename := fmt.Sprintf("ReportPriceComparison_%s_%s_%s.xlsx", time.Now().Format("2006-01-02"), util.ReplaceSpace(area.Name), util.GenerateRandomDoc(5))

	fileDir := fmt.Sprintf("%s/%s", dir, filename)

	file = xlsx.NewFile()
	if sheet, err = file.AddSheet("Sheet1"); err == nil {
		row = sheet.AddRow()
		row.SetHeight(20)
		row.AddCell().Value = "No"
		row.AddCell().Value = "Survey Date"
		row.AddCell().Value = "Area"
		row.AddCell().Value = "Product Code"
		row.AddCell().Value = "Product Name"
		row.AddCell().Value = "UOM"
		row.AddCell().Value = "Selling Price"
		row.AddCell().Value = "Public Price 1"
		row.AddCell().Value = "Public Price 2"

		for i, v := range data {
			row = sheet.AddRow()
			row.AddCell().SetInt(i + 1)
			row.AddCell().Value = v.SurveyDate                       // Survey Date
			row.AddCell().Value = v.AreaName                         // Area Name
			row.AddCell().Value = v.ProductCode                      // Product Code
			row.AddCell().Value = v.ProductName                      // Product Name
			row.AddCell().Value = v.UOM                              // UOM
			row.AddCell().SetFloatWithFormat(v.SellingPrice, "0.00") // Selling Price
			row.AddCell().SetFloatWithFormat(v.PublicPrice1, "0.00") // Public Price 1
			row.AddCell().SetFloatWithFormat(v.PublicPrice2, "0.00") // Public Price 2
		}
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
	filePath, err = util.UploadToS3(filename, fileDir, "application/xlsx")
	// fungsi ini berguna untuk menghapus kembali file yang sudah diupload ke server
	os.Remove(fileDir)
	return
}

// GetInboundXls : function to create excel file of inbound report
func GetInboundXls(date string, data []*reportInboundDetail, warehouse *model.Warehouse) (filePath string, err error) {
	var (
		file                   *xlsx.File
		sheet                  *xlsx.Sheet
		row                    *xlsx.Row
		dataSummary            []*reportInboundSummary
		etaDateSummary, source string
		j                      int64
	)

	dir := util.ExportDirectory
	warehouseName := "All Warehouse"
	if warehouse.ID != 0 {
		warehouseName = warehouse.Name
	}
	filename := fmt.Sprintf("ReportInbound_%s_%s_%s.xlsx", time.Now().Format("2006-01-02"), util.ReplaceSpace(warehouseName), util.GenerateRandomDoc(5))

	fileDir := fmt.Sprintf("%s/%s", dir, filename)
	file = xlsx.NewFile()

	boldStyle := xlsx.NewStyle()
	boldFont := xlsx.NewFont(10, "Liberation Sans")
	boldFont.Bold = true
	boldStyle.Font = *boldFont
	boldStyle.ApplyFont = true

	if sheet, err = file.AddSheet("Details"); err == nil {
		row = sheet.AddRow()
		row.SetHeight(20)
		row.Sheet.SetColWidth(0, 0, 5)
		row.AddCell().Value = "Warehouse : " + warehouseName

		sheet.AddRow()

		row = sheet.AddRow()
		row.SetHeight(20)
		row.AddCell().Value = "No"
		row.AddCell().Value = "Supplier Name"
		row.AddCell().Value = "PO Code"
		row.AddCell().Value = "PO ETA Date"
		row.AddCell().Value = "PO Committed at"
		row.AddCell().Value = "GR ATA Date"
		row.AddCell().Value = "On Time ETA"
		row.AddCell().Value = "On Time Committed PO"

		for i, v := range data {
			onTimeInbound := 0
			onTimeCommitPO := 0
			etaDate, _ := time.Parse("2006-01-02 15:04:05 (MST)", v.EtaDate.Format("2006-01-02")+" "+v.EtaTime+":00 (WIB)")
			ataDate, _ := time.Parse("2006-01-02 15:04:05 (MST)", v.AtaDate.Format("2006-01-02")+" "+v.AtaTime+":00 (WIB)")
			ataDateStr := ataDate.Format("02-01-2006 15:04:05")
			committedAtStr := v.CommittedAt.Format("02-01-2006 15:04:05")
			if !ataDate.IsZero() {
				if etaDate.After(ataDate) {
					onTimeInbound = 1
				}

				if !v.CommittedAt.IsZero() && (v.CommittedAt.Before(ataDate) || v.CommittedAt.Equal(ataDate)) {
					onTimeCommitPO = 1
				}
			} else {
				ataDateStr = ""
			}

			if v.CommittedAt.IsZero() {
				committedAtStr = ""
			}

			row = sheet.AddRow()
			row.AddCell().SetInt(i + 1)
			row.AddCell().Value = v.SupplierName                        // Supplier Name
			row.AddCell().Value = v.OrderCode                           // Purchase Order Code
			row.AddCell().Value = etaDate.Format("02-01-2006 15:04:05") // Eta Date
			row.AddCell().Value = committedAtStr                        // Committed At
			row.AddCell().Value = ataDateStr                            // Ata Date
			row.AddCell().SetInt(onTimeInbound)                         // On time Inbound
			row.AddCell().SetInt(onTimeCommitPO)                        // On time Commit PO

			if source != v.Source || etaDateSummary != v.EtaDate.Format("2006-01-02") {
				dataSummary = append(dataSummary, &reportInboundSummary{Source: v.Source, EtaDate: v.EtaDate, TotalData: 1, TotalFulfillInbound: int64(onTimeInbound), TotalFulfillCommit: int64(onTimeCommitPO)})
				j++
				etaDateSummary = v.EtaDate.Format("2006-01-02")
				source = v.Source
			} else {
				if onTimeInbound == 1 {
					dataSummary[j-1].TotalFulfillInbound++
				}
				if onTimeCommitPO == 1 {
					dataSummary[j-1].TotalFulfillCommit++
				}
				dataSummary[j-1].TotalData++
			}
		}
	}
	sheet.Cell(0, 0).SetStyle(boldStyle)
	for col := 0; col < 8; col++ {
		sheet.Cell(2, col).SetStyle(boldStyle)
	}

	if sheet, err = file.AddSheet("Summary"); err == nil {
		row = sheet.AddRow()
		row.SetHeight(20)
		row.Sheet.SetColWidth(0, 0, 5)
		row.AddCell().Value = "Warehouse : " + warehouse.Name

		sheet.AddRow()

		row = sheet.AddRow()
		row.SetHeight(20)
		row.AddCell().Value = "No"
		row.AddCell().Value = "Source"
		row.AddCell().Value = "ETA Date"
		row.AddCell().Value = "On Time Eta"
		row.AddCell().Value = "On Time Commit"

		for i, v := range dataSummary {
			row = sheet.AddRow()
			row.AddCell().SetInt(i + 1)
			row.AddCell().Value = v.Source                                                                 // Source
			row.AddCell().Value = v.EtaDate.Format("02-01-2006")                                           // Eta Date
			row.AddCell().SetFloatWithFormat(float64(v.TotalFulfillInbound)/float64(v.TotalData), "0.00%") // On Time Eta
			row.AddCell().SetFloatWithFormat(float64(v.TotalFulfillCommit)/float64(v.TotalData), "0.00%")  // On Time Commit
		}
	}
	sheet.Cell(0, 0).SetStyle(boldStyle)
	for col := 0; col < 5; col++ {
		sheet.Cell(2, col).SetStyle(boldStyle)
	}

	err = file.Save(fileDir)
	filePath, err = util.UploadToS3(filename, fileDir, "application/xlsx")

	// fungsi ini berguna untuk menghapus kembali file yang sudah diupload ke server
	os.Remove(fileDir)
	return
}

// GetFieldPurchaserXls : function to create excel file of field purchaser report
func GetFieldPurchaserXls(date string, data []*reportFieldPurchaser, area *model.Area) (filePath string, err error) {
	var file *xlsx.File
	var sheet *xlsx.Sheet
	var row *xlsx.Row

	dir := util.ExportDirectory
	filename := fmt.Sprintf("ReportFieldPurchaser_%s_%s_%s.xlsx", time.Now().Format("2006-01-02"), util.ReplaceSpace(area.Name), util.GenerateRandomDoc(5))

	fileDir := fmt.Sprintf("%s/%s", dir, filename)

	file = xlsx.NewFile()
	sheet, err = file.AddSheet("Sheet1")
	if err != nil {
		return "", err
	}

	row = sheet.AddRow()
	row.SetHeight(20)
	row.AddCell().Value = "No"
	row.AddCell().Value = "PP Date"
	row.AddCell().Value = "PP Code"
	row.AddCell().Value = "PO Date"
	row.AddCell().Value = "PO Code"
	row.AddCell().Value = "Sup. Organization"
	row.AddCell().Value = "Supplier"
	row.AddCell().Value = "Returnable"
	row.AddCell().Value = "Rejectable"
	row.AddCell().Value = "SKU Code"
	row.AddCell().Value = "SKU Name"
	row.AddCell().Value = "UOM"
	row.AddCell().Value = "Purchase Plan Qty"
	row.AddCell().Value = "Price Reference"
	row.AddCell().Value = "Purchase Quantity"
	row.AddCell().Value = "Unit Price"
	row.AddCell().Value = "Total Price"
	row.AddCell().Value = "Payment Term"
	row.AddCell().Value = "FP Name"
	row.AddCell().Value = "Order Location"
	row.AddCell().Value = "CS Code"
	row.AddCell().Value = "Warehouse"
	row.AddCell().Value = "Driver Name"
	row.AddCell().Value = "Vehicle Number"
	row.AddCell().Value = "Driver Phone Number"
	row.AddCell().Value = "ETA Date"
	row.AddCell().Value = "ETA Time"
	row.AddCell().Value = "ATA Date"
	row.AddCell().Value = "ATA Time"
	row.AddCell().Value = "Inbound Date"
	row.AddCell().Value = "Received Qty"
	row.AddCell().Value = "Status"

	for i, v := range data {
		row = sheet.AddRow()
		row.AddCell().SetInt(i + 1)
		row.AddCell().Value = v.PurchasePlanDate                                                           // Purchase Plan Created Date
		row.AddCell().Value = v.PurchasePlanCode                                                           // Purchase Plan Code
		row.AddCell().Value = v.PurchaseOrderDate                                                          // Purchase Order Date
		row.AddCell().Value = v.PurchaseOrderCode                                                          // Purchase Order Code
		row.AddCell().Value = v.SupplierOrganizationName                                                   // Supplier Organization Name
		row.AddCell().Value = v.SupplierName                                                               // Supplier Name
		row.AddCell().Value = v.Returnable                                                                 // Returnable
		row.AddCell().Value = v.Rejectable                                                                 // Rejectable
		row.AddCell().Value = v.ProductCode                                                                // Product Code
		row.AddCell().Value = v.ProductName                                                                // Product Name
		row.AddCell().Value = v.UOM                                                                        // UOM Name
		row.AddCell().SetFloatWithFormat(v.PurchasePlanQty, "0.00")                                        // Purchase Plan Qty
		row.AddCell().SetFloatWithFormat(v.PriceReference, "0.00")                                         // Price Reference
		row.AddCell().SetFloatWithFormat(v.PurchaseQty, "0.00")                                            // Purchase Quantity
		row.AddCell().SetFloatWithFormat(v.UnitPrice, "0.00")                                              // Unit Price
		row.AddCell().SetFloatWithFormat(v.TotalPrice, "0.00")                                             // Total Price
		row.AddCell().Value = v.PaymentTerm                                                                // Payment Term
		row.AddCell().Value = v.FieldPurchaserName                                                         // Field Purchaser Name
		row.AddCell().SetFormula(fmt.Sprintf(`HYPERLINK("%s","%s")`, v.OrderLocation, "Link Google Maps")) // Field Purchase Order Lat & Long
		row.AddCell().Value = v.ConsolidatedShipmentCode                                                   // Consolidated Shipment Code
		row.AddCell().Value = v.WarehouseName                                                              // Warehouse Name
		row.AddCell().Value = v.DriverName                                                                 // Driver Name
		row.AddCell().Value = v.VehicleNumber                                                              // Vehicle Number
		row.AddCell().Value = v.DriverPhoneNumber                                                          // Driver Phone Number
		row.AddCell().Value = v.EtaDate                                                                    // Purchase Order ETA Date
		row.AddCell().Value = v.EtaTime                                                                    // Purchase Order ETA Time
		row.AddCell().Value = v.AtaDate                                                                    // GR ATA Date
		row.AddCell().Value = v.AtaTime                                                                    // GR ATA Time
		row.AddCell().Value = v.InboundDate                                                                // Inbound Date
		row.AddCell().SetFloatWithFormat(v.ReceiveQty, "0.00")                                             // Receive Qty
		row.AddCell().Value = v.Status                                                                     // Status

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
	return
}
