// Copyright 2020 PT. Eden Pangan Indonesia. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

package report

import (
	"fmt"
	"os"
	"time"

	"git.edenfarm.id/cuxs/orm"
	"github.com/tealeg/xlsx"

	"git.edenfarm.id/project-version2/api/log"
	"git.edenfarm.id/project-version2/api/util"
	"git.edenfarm.id/project-version2/datamodel/model"
)

func getSalesOrderXls(dateDownload time.Time, date string, data []*reportSalesOrder, area *model.Area, staff *model.Staff, note string) (filePath string, err error) {
	var file *xlsx.File
	var sheet *xlsx.Sheet
	var row *xlsx.Row

	dir := util.ExportDirectory
	filename := fmt.Sprintf("ReportSalesOrder_%s_%s_%s.xlsx", time.Now().Format("2006-01-02"), util.ReplaceSpace(area.Name), util.GenerateRandomDoc(5))

	fileDir := fmt.Sprintf("%s/%s", dir, filename)

	file = xlsx.NewFile()
	if sheet, err = file.AddSheet("Sheet1"); err == nil {
		row = sheet.AddRow()
		row.AddCell().Value = "Sales Order Report"
		row = sheet.AddRow()
		row.AddCell().Value = "Download Timestamp: " + dateDownload.Format("02/01/2006 15:04 WIB")

		row = sheet.AddRow()
		row.SetHeight(20)
		row.AddCell().Value = "No"
		row.AddCell().Value = "Order Code"
		row.AddCell().Value = "Customer Code"
		row.AddCell().Value = "Customer Tag"
		row.AddCell().Value = "Customer Name"
		row.AddCell().Value = "Customer Phone Number"
		row.AddCell().Value = "Recipient Name"
		row.AddCell().Value = "Recipient Phone Number"
		row.AddCell().Value = "Shipping Address"
		row.AddCell().Value = "Warehouse"
		row.AddCell().Value = "Business Type"
		row.AddCell().Value = "Archetype"
		row.AddCell().Value = "Order Type"
		row.AddCell().Value = "Area"
		row.AddCell().Value = "City"
		row.AddCell().Value = "Salesperson"
		row.AddCell().Value = "Sales Group"
		row.AddCell().Value = "Order Date"
		row.AddCell().Value = "Order Delivery Date"
		row.AddCell().Value = "Order Status"
		row.AddCell().Value = "Order Note"
		row.AddCell().Value = "Total SKU Discount Amount"
		row.AddCell().Value = "Grand Total"
		row.AddCell().Value = "Order Channel"
		row.AddCell().Value = "Promo Code"
		row.AddCell().Value = "Delivery Fee"
		row.AddCell().Value = "Created At"
		row.AddCell().Value = "Created By"
		row.AddCell().Value = "Updated At"
		row.AddCell().Value = "Updated By"
		row.AddCell().Value = "Cancel Type"
		row.AddCell().Value = "Estimated Outgoing Time"

		for i, v := range data {

			row = sheet.AddRow()
			row.AddCell().SetInt(i + 1)
			row.AddCell().Value = v.SalesOrderCode                             // SO Code
			row.AddCell().Value = v.CustomerCode                               // Customer Code
			row.AddCell().Value = v.CustomerTag                                // Customer Tag
			row.AddCell().Value = v.CustomerName                               // Customer Name
			row.AddCell().Value = v.CustomerPhoneNumber                        // Customer Phone Number
			row.AddCell().Value = v.RecipientName                              // Recipient Name
			row.AddCell().Value = v.RecipientPhoneNumber                       // Recipient Phone Number
			row.AddCell().Value = v.ShippingAddress                            // Shipping Address
			row.AddCell().Value = v.WarehouseName                              // Warehouse Name
			row.AddCell().Value = v.BusinessType                               // Business Type
			row.AddCell().Value = v.ArchetypeName                              // Archetype Name
			row.AddCell().Value = v.OrderType                                  // Order Type
			row.AddCell().Value = v.AreaName                                   // Area Name
			row.AddCell().Value = v.City                                       // City
			row.AddCell().Value = v.Salesperson                                // Salesperson
			row.AddCell().Value = v.SalesGroup                                 // Sales Group
			row.AddCell().Value = v.OrderDate                                  // Order Date
			row.AddCell().Value = v.OrderDeliveryDate                          // Order Delivery Date
			row.AddCell().Value = v.OrderStatus                                // Order Status
			row.AddCell().Value = v.OrderNote                                  // Order Note
			row.AddCell().SetFloatWithFormat(v.TotalSKUDiscountAmount, "0.00") // Total SKU Discount Amount
			row.AddCell().SetFloatWithFormat(v.GrandTotal, "0.00")             // Grand Total
			row.AddCell().Value = v.OrderChannel                               // Order Channel
			row.AddCell().Value = v.PromoCode                                  // Promo Code
			row.AddCell().SetFloatWithFormat(v.DeliveryFee, "0.00")            // Delivery Fee
			row.AddCell().Value = v.CreatedAt                                  // Created At
			row.AddCell().Value = v.CreatedBy                                  // Created By
			row.AddCell().Value = v.LastUpdatedAt                              // Updated At
			row.AddCell().Value = v.LastUpdatedBy                              // Updated By
			row.AddCell().Value = v.CancelType                                 // Cancel Type
			row.AddCell().Value = v.ETDSt                                      // Estimated Outgoing Time
		}

	}
	err = file.Save(fileDir)
	filePath, err = util.UploadToS3(filename, fileDir, "application/xlsx")

	// fungsi ini berguna untuk menghapus kembali file yang sudah diupload ke server
	os.Remove(fileDir)

	if err = log.AuditLogByUser(staff, 0, "report_sales_order", "download", note); err != nil {
		return "", err
	}

	return
}

func getSalesOrderItemXls(dateDownload time.Time, date string, data []*reportSalesOrderItem, area *model.Area, staff *model.Staff, note string) (filePath string, err error) {
	var file *xlsx.File
	var sheet *xlsx.Sheet
	var row *xlsx.Row

	dir := util.ExportDirectory
	filename := fmt.Sprintf("ReportSalesOrderItem_%s_%s_%s.xlsx", time.Now().Format("2006-01-02"), util.ReplaceSpace(area.Name), util.GenerateRandomDoc(5))

	fileDir := fmt.Sprintf("%s/%s", dir, filename)
	file = xlsx.NewFile()
	if sheet, err = file.AddSheet("Sheet1"); err == nil {
		row = sheet.AddRow()
		row.AddCell().Value = "Sales Order Report"
		row = sheet.AddRow()
		row.AddCell().Value = "Download Timestamp: " + dateDownload.Format("02/01/2006 15:04 WIB")

		row = sheet.AddRow()
		row.SetHeight(20)
		row.AddCell().Value = "No"
		row.AddCell().Value = "Order Code"
		row.AddCell().Value = "Order Type"
		row.AddCell().Value = "Product Code"
		row.AddCell().Value = "Product Name"
		row.AddCell().Value = "Category"
		row.AddCell().Value = "UOM"
		row.AddCell().Value = "Order Item Note"
		row.AddCell().Value = "Ordered Qty"
		row.AddCell().Value = "Invoice Qty"
		row.AddCell().Value = "Order Unit Price"
		row.AddCell().Value = "Order Unit Shadow Price"
		row.AddCell().Value = "Taxable"
		row.AddCell().Value = "Tax Percentage"
		row.AddCell().Value = "Subtotal"
		row.AddCell().Value = "Total Weight"
		row.AddCell().Value = "Order Date"
		row.AddCell().Value = "Order Delivery Date"
		row.AddCell().Value = "Area"
		row.AddCell().Value = "Warehouse"
		row.AddCell().Value = "WRT"
		row.AddCell().Value = "Discount Qty"
		row.AddCell().Value = "Unit Price Discount"
		row.AddCell().Value = "SKU Discount Amount"
		row.AddCell().Value = "SKU Discount Name"

		for i, v := range data {
			row = sheet.AddRow()
			row.AddCell().SetInt(i + 1)
			row.AddCell().Value = v.SalesOrderCode                           // SO Code
			row.AddCell().Value = v.SalesOrderTypeName                       // SO Type
			row.AddCell().Value = v.ProductCode                              // Product Code
			row.AddCell().Value = v.ProductName                              // Product Name
			row.AddCell().Value = v.CategoryName                             // Category Name
			row.AddCell().Value = v.UomName                                  // UOM Name
			row.AddCell().Value = v.OrderItemNote                            // Order Item Note
			row.AddCell().SetFloatWithFormat(v.OrderedQty, "0.00")           // Ordered Qty
			row.AddCell().SetFloatWithFormat(v.InvoiceQty, "0.00")           // Invoice Qty
			row.AddCell().SetFloatWithFormat(v.OrderUnitPrice, "0.00")       // Order Unit Price
			row.AddCell().SetFloatWithFormat(v.OrderUnitShadowPrice, "0.00") // Order Unit Shadow Price
			row.AddCell().Value = v.TaxableStr                               // Is Taxable Item?
			row.AddCell().SetFloatWithFormat(v.OrderTaxPercentage, "0.00")   // Order Unit Tax Percentage
			row.AddCell().SetFloatWithFormat(v.Subtotal, "0.00")             // Subtotal
			row.AddCell().SetFloatWithFormat(v.TotalWeight, "0.00")          // Total Weight
			row.AddCell().Value = v.OrderDate                                // Order Date
			row.AddCell().Value = v.OrderDeliveryDate                        // Order Delivery Date
			row.AddCell().Value = v.AreaName                                 // Area Name
			row.AddCell().Value = v.WarehouseName                            // Warehouse Name
			row.AddCell().Value = v.WrtName                                  // Wrt Name
			row.AddCell().SetFloatWithFormat(v.DiscountQty, "0.00")          // Discount Qty
			row.AddCell().SetFloatWithFormat(v.UnitPriceDiscount, "0.00")    // Unit Price Discount
			row.AddCell().SetFloatWithFormat(v.SkuDiscAmount, "0.00")        // SKU Discount Amount
			row.AddCell().Value = v.SKUDiscountName                          // SKU Discount Name
		}

	}
	err = file.Save(fileDir)
	filePath, err = util.UploadToS3(filename, fileDir, "application/xlsx")

	// fungsi ini berguna untuk menghapus kembali file yang sudah diupload ke server
	os.Remove(fileDir)

	if err = log.AuditLogByUser(staff, 0, "report_sales_order_item", "download", note); err != nil {
		return "", err
	}

	return
}

// GetItemRecapXls: function for gathering all soi related
func GetItemRecapXls(date time.Time, r []*model.SalesOrderItem, area *model.Area) (filePath string, err error) {
	var file *xlsx.File
	var sheet *xlsx.Sheet
	var row *xlsx.Row
	dir := util.ExportDirectory

	filename := fmt.Sprintf("ReportOrderItemRecap_%s_%s_%s.xlsx", date.Format("2006-01-02"), util.ReplaceSpace(area.Name), util.GenerateRandomDoc(5))

	fileDir := fmt.Sprintf("%s/%s", dir, filename)
	file = xlsx.NewFile()
	if sheet, err = file.AddSheet("Sheet1"); err == nil {
		row = sheet.AddRow()
		row.SetHeight(20)
		row.AddCell().Value = "No"
		row.AddCell().Value = "Order_Delivery_Date"
		row.AddCell().Value = "Area"
		row.AddCell().Value = "Warehouse"
		row.AddCell().Value = "Product_Code"
		row.AddCell().Value = "Product_Name"
		row.AddCell().Value = "UOM"
		row.AddCell().Value = "Total_Quantity"
		row.AddCell().Value = "Total_Weight"

		for i, v := range r {
			v.Product.Read("ID")
			v.Product.Uom.Read("ID")

			row = sheet.AddRow()
			row.AddCell().SetInt(i + 1)
			row.AddCell().Value = v.SalesOrder.DeliveryDate.String()      // Order Delivery Date
			row.AddCell().Value = v.SalesOrder.Area.Name                  // Area
			row.AddCell().Value = v.SalesOrder.Warehouse.Name             // Warehouse
			row.AddCell().Value = v.Product.Code                          // Product Code
			row.AddCell().Value = v.Product.Name                          // Product Name
			row.AddCell().Value = v.Product.Uom.Name                      // UOM
			row.AddCell().SetFloatWithFormat(float64(v.OrderQty), "0.00") // Total Quantity
			row.AddCell().SetFloatWithFormat(float64(v.Weight), "0.00")   // Total Weight (Kg)

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

// GetItemRecap get all data sales_order_item that matched with query request parameters.
// returning slices of SalesOrderItem, total data without limit and error.
func GetItemRecap(rq *orm.RequestQuery) (m []*model.SalesOrderItem, total int64, err error) {
	// make new orm query
	q, _ := rq.QueryReadOnly(new(model.SalesOrderItem))

	// get total data
	if total, err = q.Count(); err != nil || total == 0 {
		return nil, total, err
	}

	// get data requested
	var mx []*model.SalesOrderItem
	if _, err = q.RelatedSel().All(&mx, rq.Fields...); err == nil {
		poi := removeDuplicates(mx)
		for _, i := range poi {
			i.Product.Category.Read("ID")
			i.Product.Uom.Read("ID")
		}
		return poi, total, nil
	}

	// return error some thing went wrong
	return nil, total, err
}

func removeDuplicates(poi []*model.SalesOrderItem) []*model.SalesOrderItem {
	// Use map to record duplicates as we find them.
	duplicate := map[int64]bool{}
	var result []*model.SalesOrderItem
	var storageTemp []*model.SalesOrderItem
	for _, v := range poi {
		if v.SalesOrder.Status != 3 && v.SalesOrder.Status != 4 {
			v.Product.Read("ID")
			if duplicate[v.Product.ID] == true {
				// data duplicate will append to storageTemp
				storageTemp = append(storageTemp, v)

			} else {
				// Record this element as an encountered element.
				duplicate[v.Product.ID] = true
				result = append(result, v)
				// Append to result slice.
			}
		}
	}

	for _, i := range result {
		if i.SalesOrder.Status != 3 && i.SalesOrder.Status != 4 {
			for _, a := range storageTemp {
				if i.Product.ID == a.Product.ID {
					i.OrderQty += a.OrderQty
				}
			}
		}
	}

	// Return the new slice.
	return result
}

func getSalesInvoiceXls(date time.Time, data []*reportSalesInvoice, area *model.Area, staff *model.Staff, note string) (filePath string, err error) {
	var file *xlsx.File
	var sheet *xlsx.Sheet
	var row *xlsx.Row

	dir := util.ExportDirectory
	filename := fmt.Sprintf("ReportSalesInvoice_%s_%s_%s.xlsx", time.Now().Format("2006-01-02"), util.ReplaceSpace(area.Name), util.GenerateRandomDoc(5))

	fileDir := fmt.Sprintf("%s/%s", dir, filename)
	filePath = fmt.Sprintf("%s/%s", util.ResponseURLFromPostUpload+"report", filename)

	file = xlsx.NewFile()
	if sheet, err = file.AddSheet("Sheet1"); err == nil {
		row = sheet.AddRow()
		row.AddCell().Value = "Sales Invoice Report"
		row = sheet.AddRow()
		row.AddCell().Value = "Download Timestamp: " + date.Format("02/01/2006 15:04 WIB")

		row = sheet.AddRow()
		row.SetHeight(20)
		row.AddCell().Value = "No"
		row.AddCell().Value = "Order Code"
		row.AddCell().Value = "Main Outlet / Agent Code"
		row.AddCell().Value = "Main Outlet / Agent Name"
		row.AddCell().Value = "Outlet / Address Code"
		row.AddCell().Value = "Outlet / Address Name"
		row.AddCell().Value = "Order Delivery Date"
		row.AddCell().Value = "Invoice Code"
		row.AddCell().Value = "Invoice Date"
		row.AddCell().Value = "Invoice Due Date"
		row.AddCell().Value = "Invoice Status"
		row.AddCell().Value = "Adjustment Note"
		row.AddCell().Value = "Total Confirmed Payment"
		row.AddCell().Value = "Total Invoice"
		row.AddCell().Value = "Delivery Fee"
		row.AddCell().Value = "Point Redeem Amount"
		row.AddCell().Value = "Voucher Amount"
		row.AddCell().Value = "Adjustment Amount"
		row.AddCell().Value = "Total SKU Discount Amount"
		row.AddCell().Value = "Total Charge"
		row.AddCell().Value = "Area"
		row.AddCell().Value = "Warehouse"
		row.AddCell().Value = "Customer Group"
		row.AddCell().Value = "Business Type"
		row.AddCell().Value = "Archetype"
		row.AddCell().Value = "Payment Term"
		row.AddCell().Value = "Invoice Term"
		row.AddCell().Value = "Created At"
		row.AddCell().Value = "Created By"
		row.AddCell().Value = "Updated At"
		row.AddCell().Value = "Updated By"

		for i, v := range data {
			row = sheet.AddRow()
			row.AddCell().SetInt(i + 1)
			row.AddCell().Value = v.SalesOrderCode                         // SO Code
			row.AddCell().Value = v.MerchantCode                           // Merchant Code
			row.AddCell().Value = v.MerchantName                           // Merchant Name
			row.AddCell().Value = v.BranchCode                             // Branch Code
			row.AddCell().Value = v.BranchName                             // Branch Name
			row.AddCell().Value = v.OrderDeliveryDate                      // Order Delivery Date
			row.AddCell().Value = v.InvoiceCode                            // Invoice Code
			row.AddCell().Value = v.InvoiceDate                            // Invoice Date
			row.AddCell().Value = v.InvoiceDueDate                         // Invoice Due Date
			row.AddCell().Value = v.InvoiceStatus                          // Invoice Status
			row.AddCell().Value = v.AdjNote                                // Adjustment Note
			row.AddCell().SetFloatWithFormat(v.TotalConfPay, "0.00")       // Total Confirmed Payment
			row.AddCell().SetFloatWithFormat(v.TotalInvoice, "0.00")       // Total Invoice
			row.AddCell().SetFloatWithFormat(v.DeliveryFee, "0.00")        // Delivery Fee
			row.AddCell().SetFloatWithFormat(v.PointRedeemAmount, "0.00")  // PointRedeemAmount
			row.AddCell().SetFloatWithFormat(v.VouAmount, "0.00")          // Voucher Amount
			row.AddCell().SetFloatWithFormat(v.AdjAmount, "0.00")          // Adjustment Amount
			row.AddCell().SetFloatWithFormat(v.TotalSkuDiscAmount, "0.00") // Total Sku Discount
			row.AddCell().SetFloatWithFormat(v.TotalCharge, "0.00")        // Total Charge
			row.AddCell().Value = v.Area                                   // Area
			row.AddCell().Value = v.Warehouse                              // Warehouse
			row.AddCell().Value = v.CustomerGroup                          // Customer Group
			row.AddCell().Value = v.BusinessType                           // Business Type
			row.AddCell().Value = v.Archetype                              // Archetype Name
			row.AddCell().Value = v.PaymentTerm                            // Payment Term
			row.AddCell().Value = v.InvoiceTerm                            // Invoice Term
			row.AddCell().Value = v.CreatedAt                              // Created At
			row.AddCell().Value = v.CreatedBy                              // Created By
			row.AddCell().Value = v.LastUpdatedAt                          // Updated At
			row.AddCell().Value = v.LastUpdatedBy                          // Updated By
		}
	}

	err = file.Save(fileDir)
	filePath, err = util.UploadToS3(filename, fileDir, "application/xlsx")
	// fungsi ini berguna untuk menghapus kembali file yang sudah diupload ke server
	os.Remove(fileDir)

	if err = log.AuditLogByUser(staff, 0, "report_sales_invoice", "download", note); err != nil {
		return "", err
	}

	return
}

func getSalesPaymentXls(date time.Time, data []*reportSalesPayment, area *model.Area, staff *model.Staff, note string) (filePath string, err error) {
	var file *xlsx.File
	var sheet *xlsx.Sheet
	var row *xlsx.Row

	dir := util.ExportDirectory
	filename := fmt.Sprintf("ReportSalesPayment_%s_%s_%s.xlsx", time.Now().Format("2006-01-02"), util.ReplaceSpace(area.Name), util.GenerateRandomDoc(5))

	fileDir := fmt.Sprintf("%s/%s", dir, filename)
	filePath = fmt.Sprintf("%s/%s", util.ResponseURLFromPostUpload+"report", filename)

	file = xlsx.NewFile()
	if sheet, err = file.AddSheet("Sheet1"); err == nil {
		row = sheet.AddRow()
		row.AddCell().Value = "Sales Payment Report"
		row = sheet.AddRow()
		row.AddCell().Value = "Download Timestamp: " + date.Format("02/01/2006 15:04 WIB")

		row = sheet.AddRow()
		row.SetHeight(20)
		row.AddCell().Value = "No"
		row.AddCell().Value = "Payment Code"
		row.AddCell().Value = "Bank Receive Voucher Number"
		row.AddCell().Value = "Payment Date"
		row.AddCell().Value = "Received Date"
		row.AddCell().Value = "Area"
		row.AddCell().Value = "Warehouse"
		row.AddCell().Value = "Payment Status"
		row.AddCell().Value = "Payment Method"
		row.AddCell().Value = "Payment Amount"
		row.AddCell().Value = "Order Code"
		row.AddCell().Value = "Invoice Code"
		row.AddCell().Value = "Customer Code"
		row.AddCell().Value = "Customer Name"
		row.AddCell().Value = "Created By"
		row.AddCell().Value = "Created At"

		for i, v := range data {
			row = sheet.AddRow()
			row.AddCell().SetInt(i + 1)
			row.AddCell().Value = v.PaymentCode                       // Payment Code
			row.AddCell().Value = v.BankReceiveVouNum                 // Bank Receive Voucher Number
			row.AddCell().Value = v.PaymentDate                       // Payment Date
			row.AddCell().Value = v.ReceivedDate                      // Received Date
			row.AddCell().Value = v.Area                              // Area
			row.AddCell().Value = v.WarehouseName                     // Warehouse Name
			row.AddCell().Value = v.PaymentStatus                     // Payment Status
			row.AddCell().Value = v.PaymentMethod                     // Payment Method
			row.AddCell().SetFloatWithFormat(v.PaymentAmount, "0.00") // Payment Amount
			row.AddCell().Value = v.OrderCode                         // Order Code
			row.AddCell().Value = v.InvoiceCode                       // Invoice Code
			row.AddCell().Value = v.OutletCode                        // Customer Code
			row.AddCell().Value = v.OutletName                        // Customer Name
			row.AddCell().Value = v.CreatedBy                         // Created By
			row.AddCell().Value = v.CreatedAt                         // Created At
		}
	}

	err = file.Save(fileDir)
	filePath, err = util.UploadToS3(filename, fileDir, "application/xlsx")
	// fungsi ini berguna untuk menghapus kembali file yang sudah diupload ke server
	os.Remove(fileDir)

	if err = log.AuditLogByUser(staff, 0, "report_sales_payment", "download", note); err != nil {
		return "", err
	}

	return
}

// getProspectiveCustomerXls : function to create excel file of prospective customer report
func getProspectiveCustomerXls(date time.Time, data []*reportProspectiveCustomer, area *model.Area, staff *model.Staff) (filePath string, err error) {
	var file *xlsx.File
	var sheet *xlsx.Sheet
	var row *xlsx.Row

	dir := util.ExportDirectory
	filename := fmt.Sprintf("ReportProspectiveCustomer_%s_%s_%s.xlsx", time.Now().Format("2006-01-02"), util.ReplaceSpace(area.Name), util.GenerateRandomDoc(5))

	fileDir := fmt.Sprintf("%s/%s", dir, filename)
	filePath = fmt.Sprintf("%s/%s", util.ResponseURLFromPostUpload+"report", filename)

	file = xlsx.NewFile()
	if sheet, err = file.AddSheet("Sheet1"); err == nil {
		row = sheet.AddRow()
		row.AddCell().Value = "Prospective Customer Report"
		row = sheet.AddRow()
		row.AddCell().Value = "Download Timestamp: " + date.Format("02/01/2006 15:04 WIB")

		row = sheet.AddRow()
		row.SetHeight(20)
		row.AddCell().Value = "No"
		row.AddCell().Value = "Prospective Customer Code"
		row.AddCell().Value = "Prospective Customer Name"
		row.AddCell().Value = "Business Type"
		row.AddCell().Value = "Archetype"
		row.AddCell().Value = "PIC Name"
		row.AddCell().Value = "Phone Number"
		row.AddCell().Value = "Pic Finance Name"
		row.AddCell().Value = "Pic Finance Contact"
		row.AddCell().Value = "Pic Business Name"
		row.AddCell().Value = "Pic Business Contact"
		row.AddCell().Value = "Term Of Payment"
		row.AddCell().Value = "Term Of Invoice"
		row.AddCell().Value = "Billing Address"
		row.AddCell().Value = "Note"
		row.AddCell().Value = "Area"
		row.AddCell().Value = "Address"
		row.AddCell().Value = "Province"
		row.AddCell().Value = "City"
		row.AddCell().Value = "District"
		row.AddCell().Value = "Sub District"
		row.AddCell().Value = "Postal Code"
		row.AddCell().Value = "Best Time to Call"
		row.AddCell().Value = "Referral Code"
		row.AddCell().Value = "Existing Customer"
		row.AddCell().Value = "Request Upgrade"
		row.AddCell().Value = "Created At"
		row.AddCell().Value = "Processed At"
		row.AddCell().Value = "Processed By"
		row.AddCell().Value = "Status"
		row.AddCell().Value = "Salesperson"
		row.AddCell().Value = "Decline Type"
		row.AddCell().Value = "Decline Note"

		for i, v := range data {
			row = sheet.AddRow()
			row.AddCell().SetInt(i + 1)
			row.AddCell().Value = v.CustomerCode       // Prospective Customer Code
			row.AddCell().Value = v.CustomerName       // Prospective Customer Name
			row.AddCell().Value = v.BusinessType       // Business Type
			row.AddCell().Value = v.Archetype          // Archetype
			row.AddCell().Value = v.PicName            // PIC Name
			row.AddCell().Value = v.PhoneNumber        // Phone Number
			row.AddCell().Value = v.PicFinanceName     // Pic Finance Name
			row.AddCell().Value = v.PicFinanceContact  // Pic Finance Contact
			row.AddCell().Value = v.PicBusinessName    // Pic Business Name
			row.AddCell().Value = v.PicBusinessContact // Pic Business Contact
			row.AddCell().Value = v.PaymentTerm        // Payment Term
			row.AddCell().Value = v.InvoiceTerm        // Invoice Term
			row.AddCell().Value = v.BillingAddress     // Billing Address
			row.AddCell().Value = v.Note               // Note
			row.AddCell().Value = v.Area               // Area
			row.AddCell().Value = v.StreetAddress      // Address
			row.AddCell().Value = v.Province           // Province
			row.AddCell().Value = v.City               // City
			row.AddCell().Value = v.District           // District
			row.AddCell().Value = v.SubDistrict        // Sub District
			row.AddCell().Value = v.PostalCode         // Postal Code
			row.AddCell().Value = v.BestTimeToCall     // Best Time to Call
			row.AddCell().Value = v.ReferralCode       // Referral Code
			row.AddCell().Value = v.ExistingCustomer   // Existing Customer
			row.AddCell().Value = v.ReqUpgrade         // Request Upgrade
			row.AddCell().Value = v.CreatedAt          // Created At
			row.AddCell().Value = v.ProcessedAt        // Processed At
			row.AddCell().Value = v.ProcessedBy        // Processed By
			row.AddCell().Value = v.Status             // Status
			row.AddCell().Value = v.Salesperson        // Salesperson
			row.AddCell().Value = v.DeclineType        // Decline Type
			row.AddCell().Value = v.DeclineNote        // Decline Note
		}
	}

	err = file.Save(fileDir)
	filePath, err = util.UploadToS3(filename, fileDir, "application/xlsx")
	// fungsi ini berguna untuk menghapus kembali file yang sudah diupload ke server
	os.Remove(fileDir)

	if err = log.AuditLogByUser(staff, 0, "prospect_customer", "Download", "Download Prospective Customer"); err != nil {
		return "", err
	}

	return
}

// getSkuDiscountXls : function to create excel of sku discount report
func getSkuDiscountXls(date time.Time, data []*reportSkuDiscount, staff *model.Staff, note string) (filePath string, err error) {
	var file *xlsx.File
	var sheet *xlsx.Sheet
	var row *xlsx.Row

	dir := util.ExportDirectory
	filename := fmt.Sprintf("ReportSkuDiscount_%s_%s.xlsx", time.Now().Format("2006-01-02"), util.GenerateRandomDoc(5))

	fileDir := fmt.Sprintf("%s/%s", dir, filename)

	file = xlsx.NewFile()
	if sheet, err = file.AddSheet("Sheet1"); err == nil {
		row = sheet.AddRow()
		row.AddCell().Value = "SKU Discount Report"
		row = sheet.AddRow()
		row.AddCell().Value = "Download Timestamp: " + date.Format("02/01/2006 15:04 WIB")

		row = sheet.AddRow()
		row.SetHeight(20)
		row.AddCell().Value = "No"
		row.AddCell().Value = "Promotion Name"
		row.AddCell().Value = "Price Set"
		row.AddCell().Value = "Start Period"
		row.AddCell().Value = "End Period"
		row.AddCell().Value = "Division"
		row.AddCell().Value = "Order Channel"
		row.AddCell().Value = "Note"

		for i, v := range data {
			row = sheet.AddRow()
			row.AddCell().SetInt(i + 1)
			row.AddCell().Value = v.SkuDiscName                      // Sku Disc Name
			row.AddCell().Value = v.PriceSet                         // PriceSet
			row.AddCell().Value = v.StartPeriod.Format("2006-01-02") // Start Period
			row.AddCell().Value = v.EndPeriod.Format("2006-01-02")   // End Period
			row.AddCell().Value = v.Division                         // Division
			row.AddCell().Value = v.OrderChannel                     // Order Channel
			row.AddCell().Value = v.Note                             // Note
		}
	}

	err = file.Save(fileDir)
	if err != nil {
		fmt.Println(err)
	}
	filePath, err = util.UploadToS3(filename, fileDir, "application/xlsx")

	// fungsi ini berguna untuk menghapus kembali file yang sudah diupload ke server
	os.Remove(fileDir)

	if err = log.AuditLogByUser(staff, 0, "report_sku_discount", "download", note); err != nil {
		return "", err
	}

	return
}

// getSkuDiscounItemtXls : function to create excel of sku discount item report
func getSkuDiscounItemXls(date time.Time, data []*reportSkuDiscountItem, staff *model.Staff, note string) (filePath string, err error) {
	var file *xlsx.File
	var sheet *xlsx.Sheet
	var row *xlsx.Row

	dir := util.ExportDirectory
	filename := fmt.Sprintf("ReportSkuDiscountItem_%s_%s.xlsx", time.Now().Format("2006-01-02"), util.GenerateRandomDoc(5))

	fileDir := fmt.Sprintf("%s/%s", dir, filename)

	file = xlsx.NewFile()
	if sheet, err = file.AddSheet("Sheet1"); err == nil {
		row = sheet.AddRow()
		row.AddCell().Value = "SKU Discount Item Report"
		row = sheet.AddRow()
		row.AddCell().Value = "Download Timestamp: " + date.Format("02/01/2006 15:04 WIB")

		row = sheet.AddRow()
		row.SetHeight(20)
		row.AddCell().Value = "No"
		row.AddCell().Value = "Promotion Name"
		row.AddCell().Value = "Product Code"
		row.AddCell().Value = "Product Name"
		row.AddCell().Value = "UOM"
		row.AddCell().Value = "Tier"
		row.AddCell().Value = "Minimum Qty"
		row.AddCell().Value = "Amount"
		row.AddCell().Value = "Overall Quota"
		row.AddCell().Value = "Quota Per User"
		row.AddCell().Value = "Daily Quota Per User"
		row.AddCell().Value = "Budget"
		row.AddCell().Value = "Remaining Budget"

		for i, v := range data {
			row = sheet.AddRow()
			row.AddCell().SetInt(i + 1)
			row.AddCell().Value = v.SkuDiscName                      // Sku Disc Name
			row.AddCell().Value = v.ProductCode                      // Product Code
			row.AddCell().Value = v.ProductName                      // Product Name
			row.AddCell().Value = v.Uom                              // UOM
			row.AddCell().SetInt(int(v.TierLevel))                   // Tier Level
			row.AddCell().SetFloatWithFormat(v.MinimumQty, "0.00")   // Minimum Qty
			row.AddCell().SetFloatWithFormat(v.Amount, "0.00")       // Percentage
			row.AddCell().SetFloatWithFormat(v.OverallQuota, "0.00") // Overall Quota
			if v.OverallQuotaPerUser == 0 {
				row.AddCell().Value = "-"
			} else {
				row.AddCell().SetFloatWithFormat(v.OverallQuotaPerUser, "0.00") // Overall Quota Per User
			}
			if v.DailyQuotaPerUser == 0 {
				row.AddCell().Value = "-"
			} else {
				row.AddCell().SetFloatWithFormat(v.DailyQuotaPerUser, "0.00") // Daily Quota Per User
			}
			if v.Budget == 0 {
				row.AddCell().Value = "-"
			} else {
				row.AddCell().SetFloatWithFormat(v.Budget, "0.00") // Budget
			}
			if v.RemBudget == 0 {
				row.AddCell().Value = "-"
			} else {
				row.AddCell().SetFloatWithFormat(v.RemBudget, "0.00") // Remaining Budget
			}
		}
	}

	err = file.Save(fileDir)
	filePath, err = util.UploadToS3(filename, fileDir, "application/xlsx")

	// fungsi ini berguna untuk menghapus kembali file yang sudah diupload ke server
	os.Remove(fileDir)

	if err = log.AuditLogByUser(staff, 0, "report_sku_discount_item", "download", note); err != nil {
		return "", err
	}

	return
}

func getSalesOrderFeedbackXls(date time.Time, data []*reportSalesOrderFeedback, staff *model.Staff) (filePath string, err error) {
	var file *xlsx.File
	var sheet *xlsx.Sheet
	var row *xlsx.Row

	dir := util.ExportDirectory
	filename := fmt.Sprintf("ReportSalesOrderFeedback_%s_%s.xlsx", time.Now().Format("2006-01-02"), util.GenerateRandomDoc(5))

	fileDir := fmt.Sprintf("%s/%s", dir, filename)

	file = xlsx.NewFile()
	if sheet, err = file.AddSheet("Sheet1"); err == nil {
		row = sheet.AddRow()
		row.AddCell().Value = "Sales Order Feedback Report"
		row = sheet.AddRow()
		row.AddCell().Value = "Download Timestamp: " + date.Format("02/01/2006 15:04 WIB")

		row = sheet.AddRow()
		row.SetHeight(20)
		row.AddCell().Value = "No"
		row.AddCell().Value = "Merchant_Code"
		row.AddCell().Value = "Merchant_Name"
		row.AddCell().Value = "Merchant_Phone_Number"
		row.AddCell().Value = "Business_Type"
		row.AddCell().Value = "Archetype"
		row.AddCell().Value = "Branch_Shipping_Address"
		row.AddCell().Value = "City"
		row.AddCell().Value = "District"
		row.AddCell().Value = "Area"
		row.AddCell().Value = "Sales_Order_Code"
		row.AddCell().Value = "Delivery_Date"
		row.AddCell().Value = "Feedback_Created_At"
		row.AddCell().Value = "Rating_Score"
		row.AddCell().Value = "Tags"
		row.AddCell().Value = "Feedback_Description"
		row.AddCell().Value = "To_Be_Contacted"

		for i, v := range data {

			row = sheet.AddRow()
			row.AddCell().SetInt(i + 1)
			row.AddCell().Value = v.MerchantCode          // Merchant Code
			row.AddCell().Value = v.MerchantName          // Merchant Name
			row.AddCell().Value = v.MerchantPhoneNumber   // Merchant Phone Number
			row.AddCell().Value = v.BusinessType          // Business Type
			row.AddCell().Value = v.Archetype             // Archetype
			row.AddCell().Value = v.BranchShippingAddress // Customer Shipping Address
			row.AddCell().Value = v.City                  // City
			row.AddCell().Value = v.District              // District
			row.AddCell().Value = v.Area                  // Area Name
			row.AddCell().Value = v.SalesOrderCode        // SO Code
			row.AddCell().Value = v.DeliveryDate          // Delivery Date
			row.AddCell().Value = v.FeedBackCreatedAt     // Feedback Created At
			row.AddCell().Value = v.RatingScore           // Rating Score
			row.AddCell().Value = v.Tags                  // Merchant Tags
			row.AddCell().Value = v.FeedbackDescription   // Feedback Desciprtion
			row.AddCell().Value = v.ToBeContacted         // To Be Contacted
		}

	}
	err = file.Save(fileDir)
	filePath, err = util.UploadToS3(filename, fileDir, "application/xlsx")

	// fungsi ini berguna untuk menghapus kembali file yang sudah diupload ke server
	os.Remove(fileDir)

	if err = log.AuditLogByUser(staff, 0, "report_sales_order_feedback", "download", ""); err != nil {
		return "", err
	}

	return
}

// getEdenPointXls : function to create excel of EdenPoint Report
func getEdenPointXls(date time.Time, data []*reportEdenPointLog, staff *model.Staff, note string) (filePath string, err error) {
	var file *xlsx.File
	var sheet *xlsx.Sheet
	var row *xlsx.Row

	dir := util.ExportDirectory
	filename := fmt.Sprintf("ReportEdenPointLog_%s_%s.xlsx", time.Now().Format("2006-01-02"), util.GenerateRandomDoc(5))

	fileDir := fmt.Sprintf("%s/%s", dir, filename)

	file = xlsx.NewFile()
	if sheet, err = file.AddSheet("Sheet1"); err == nil {
		row = sheet.AddRow()
		row.AddCell().Value = "EdenPoint Log Report"
		row = sheet.AddRow()
		row.AddCell().Value = "Download Timestamp: " + date.Format("02/01/2006 15:04 WIB")

		row = sheet.AddRow()
		row.SetHeight(20)
		row.AddCell().Value = "No"
		row.AddCell().Value = "EdenPoint Date"
		row.AddCell().Value = "Merchant/Agent"
		row.AddCell().Value = "Previous EdenPoint"
		row.AddCell().Value = "EdenPoint Movement"
		row.AddCell().Value = "Status"
		row.AddCell().Value = "Current EdenPoint"
		row.AddCell().Value = "Transaction Type"
		row.AddCell().Value = "Advocate Merchant"
		row.AddCell().Value = "Referee Merchant"
		row.AddCell().Value = "Campaign Id"
		row.AddCell().Value = "Campaign Name"
		row.AddCell().Value = "Campaign Multiplier"
		row.AddCell().Value = "Log Note"
		row.AddCell().Value = "Order Code"
		row.AddCell().Value = "Order Date"
		row.AddCell().Value = "Created Order Date"
		row.AddCell().Value = "Finished Order Date"
		row.AddCell().Value = "Total Sales Order"
		row.AddCell().Value = "Order Status"

		for i, v := range data {
			row = sheet.AddRow()
			row.AddCell().SetInt(i + 1)
			row.AddCell().Value = v.EdenPointDate                         // EdenPoint Date
			row.AddCell().Value = v.MerchantName                          // Merchant Name
			row.AddCell().SetFloatWithFormat(v.PreviousEdenPoint, "0.00") // Previous EdenPoint
			row.AddCell().SetFloatWithFormat(v.EdenPoint, "0.00")         // EdenPoint Movement
			row.AddCell().Value = v.Status                                // Status
			row.AddCell().SetFloatWithFormat(v.CurrentEdenPoint, "0.00")  // Current EdenPoint
			row.AddCell().Value = v.TransactionType                       // Transaction Type
			row.AddCell().Value = v.AdvocateMerchant                      // Advocate Merchant
			row.AddCell().Value = v.RefereeMerchant                       // Referee Merchant
			row.AddCell().Value = v.CampaignId                            // Campaign Id
			row.AddCell().Value = v.CampaignName                          // Campaign Name
			row.AddCell().SetInt(int(v.CampaignMultiplier))               // Campaign Multiplier
			row.AddCell().Value = v.LogNote                               // Log Note
			row.AddCell().Value = v.OrderCode                             // Order Code
			row.AddCell().Value = v.OrderDate                             // Order Date
			row.AddCell().Value = v.CreatedOrderDate                      // Created Order Date
			row.AddCell().Value = v.FinishedOrderDate                     // Finished Order Date
			row.AddCell().SetFloatWithFormat(v.TotalSalesOrder, "0.00")   // Total Sales Order
			row.AddCell().Value = v.OrderStatus                           // Order Status
		}
	}

	err = file.Save(fileDir)
	filePath, err = util.UploadToS3(filename, fileDir, "application/xlsx")

	// fungsi ini berguna untuk menghapus kembali file yang sudah diupload ke server
	os.Remove(fileDir)

	if err = log.AuditLogByUser(staff, 0, "report_eden_point_log", "download", note); err != nil {
		return "", err
	}

	return
}
