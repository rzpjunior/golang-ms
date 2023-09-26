// Copyright 2020 PT. Eden Pangan Indonesia. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

package report

import (
	"fmt"
	"os"
	"strconv"
	"time"

	"github.com/tealeg/xlsx"

	"git.edenfarm.id/project-version2/api/log"
	"git.edenfarm.id/project-version2/api/util"
	"git.edenfarm.id/project-version2/datamodel/model"
)

// getPaymentGatewayXls : function to create excel file of payment gateway report
func getPaymentGatewayXls(date time.Time, data []*reportPaymentGateway, staff *model.Staff, note string) (filePath string, err error) {
	var file *xlsx.File
	var sheet *xlsx.Sheet
	var row *xlsx.Row

	dir := util.ExportDirectory
	filename := fmt.Sprintf("ReportPaymentGateway_%s_%s.xlsx", time.Now().Format("2006-01-02"), util.GenerateRandomDoc(5))

	fileDir := fmt.Sprintf("%s/%s", dir, filename)
	filePath = fmt.Sprintf("%s/%s", util.ResponseURLFromPostUpload+"report", filename)

	file = xlsx.NewFile()
	if sheet, err = file.AddSheet("Sheet1"); err == nil {
		row = sheet.AddRow()
		row.AddCell().Value = "Payment Gateway Report"
		row = sheet.AddRow()
		row.AddCell().Value = "Download Timestamp: " + date.Format("02/01/2006 15:04 WIB")

		row = sheet.AddRow()
		row.SetHeight(20)
		row.AddCell().Value = "No"
		row.AddCell().Value = "Main Outlet / Agent Code"
		row.AddCell().Value = "Main Outlet / Agent Name"
		row.AddCell().Value = "Outlet / Address Code"
		row.AddCell().Value = "Outlet / Address Name"
		row.AddCell().Value = "Sales Order Code"
		row.AddCell().Value = "Type"
		row.AddCell().Value = "Channel"
		row.AddCell().Value = "Account Number"
		row.AddCell().Value = "Total Amount"
		row.AddCell().Value = "Transaction Date"
		row.AddCell().Value = "Transaction Time"

		for i, v := range data {
			row = sheet.AddRow()
			row.AddCell().SetInt(i + 1)
			row.AddCell().Value = v.MerchantCode                    // Merchant Code
			row.AddCell().Value = v.MerchantName                    // Merchant Name
			row.AddCell().Value = v.BranchCode                      // Branch Code
			row.AddCell().Value = v.BranchName                      // Branch Name
			row.AddCell().Value = v.SalesOrderCode                  // SO Code
			row.AddCell().Value = v.Type                            // Type
			row.AddCell().Value = v.Channel                         // Channel
			row.AddCell().Value = v.AccountNumber                   // Account Number
			row.AddCell().SetFloatWithFormat(v.TotalAmount, "0.00") // Total Amount
			row.AddCell().Value = v.TransactionDate                 // Transaction Date
			row.AddCell().Value = v.TransactionTime                 // Transaction Time
		}

	}
	err = file.Save(fileDir)
	filePath, err = util.UploadToS3(filename, fileDir, "application/xlsx")

	// fungsi ini berguna untuk menghapus kembali file yang sudah diupload ke server
	os.Remove(fileDir)

	if err = log.AuditLogByUser(staff, 0, "report_payment_gateway", "Download", note); err != nil {
		return "", err
	}

	return
}

// getMainOutletXls : function to create excel file of main outlet report
func getMainOutletXls(date time.Time, data []*reportMainOutlet, staff *model.Staff, note string) (filePath string, err error) {
	var file *xlsx.File
	var sheet *xlsx.Sheet
	var row *xlsx.Row

	dir := util.ExportDirectory
	filename := fmt.Sprintf("ReportMainOutlet_%s_%s.xlsx", time.Now().Format("2006-01-02"), util.GenerateRandomDoc(5))

	fileDir := fmt.Sprintf("%s/%s", dir, filename)
	filePath = fmt.Sprintf("%s/%s", util.ResponseURLFromPostUpload+"report", filename)

	file = xlsx.NewFile()
	if sheet, err = file.AddSheet("Sheet1"); err == nil {
		row = sheet.AddRow()
		row.AddCell().Value = "Main Outlet Report"
		row = sheet.AddRow()
		row.AddCell().Value = "Download Timestamp: " + date.Format("02/01/2006 15:04 WIB")

		row = sheet.AddRow()
		row.SetHeight(20)
		row.AddCell().Value = "No"
		row.AddCell().Value = "Code"
		row.AddCell().Value = "Name"
		row.AddCell().Value = "Billing Address"
		row.AddCell().Value = "Finance Area"
		row.AddCell().Value = "PIC Name"
		row.AddCell().Value = "Phone Number"
		row.AddCell().Value = "Email"
		row.AddCell().Value = "Current EdenPoint"
		row.AddCell().Value = "Default Payment Term"
		row.AddCell().Value = "Default Invoice Term"
		row.AddCell().Value = "Payment Group"
		row.AddCell().Value = "Business Type"
		row.AddCell().Value = "Business Type Credit Limit"
		row.AddCell().Value = "Credit Limit Amount"
		row.AddCell().Value = "Credit Limit Remaining"
		row.AddCell().Value = "Suspended"
		row.AddCell().Value = "Customer Group"
		row.AddCell().Value = "Customer Tag"
		row.AddCell().Value = "Status"
		row.AddCell().Value = "Created At"
		row.AddCell().Value = "Created By"
		row.AddCell().Value = "Last Updated At"
		row.AddCell().Value = "Last Updated By"

		for i, v := range data {
			row = sheet.AddRow()
			row.AddCell().SetInt(i + 1)
			row.AddCell().Value = v.Code               // Code
			row.AddCell().Value = v.Name               // Name
			row.AddCell().Value = v.BillingAddress     // Billing Address
			row.AddCell().Value = v.FinanceArea        // Finance Area
			row.AddCell().Value = v.PicName            // PIC Name
			row.AddCell().Value = v.PhoneNumber        // Phone Number
			row.AddCell().Value = v.Email              // Email
			row.AddCell().Value = v.CurrentEdenPoint   // Current EdenPoint
			row.AddCell().Value = v.DefaultPaymentTerm // Default Payment Term
			row.AddCell().Value = v.DefaultInvoiceTerm // Default Invoice Term
			row.AddCell().Value = v.PaymentGroup       // Payment Group
			row.AddCell().Value = v.BusinessType       // Business Type
			if v.BusinessTypeCreditLimit == "1" {
				row.AddCell().Value = "Badan Usaha" // Business Type Credit Limit
			} else {
				row.AddCell().Value = "Personal"
			}
			row.AddCell().Value = v.CreditLimitAmount          // Credit Limit Amount
			row.AddCell().Value = v.RemainingCreditLimitAmount // Credit Limit Remaining
			if v.Suspended == "1" {
				row.AddCell().Value = "Yes" // Suspended
			} else {
				row.AddCell().Value = "No"
			}
			row.AddCell().Value = v.CustomerGroup // Customer Group
			row.AddCell().Value = v.CustomerTag   // Customer Tag
			switch v.Status {                     // Status
			case "1":
				row.AddCell().Value = "Active"
			case "2":
				row.AddCell().Value = "Archive"
			case "3":
				row.AddCell().Value = "Deleted"
			default:
				row.AddCell().Value = ""
			}
			row.AddCell().Value = v.CreatedAt     // Created At
			row.AddCell().Value = v.CreatedBy     // Created By
			row.AddCell().Value = v.LastUpdatedAt // Updated At
			row.AddCell().Value = v.LastUpdatedBy // Updated By
		}

	}
	err = file.Save(fileDir)
	filePath, err = util.UploadToS3(filename, fileDir, "application/xlsx")

	// fungsi ini berguna untuk menghapus kembali file yang sudah diupload ke server
	os.Remove(fileDir)

	if err = log.AuditLogByUser(staff, 0, "main_outlet", "export", note); err != nil {
		return "", err
	}

	return
}

// getOutletXls : function to create excel file of outlet report
func getOutletXls(date time.Time, data []*reportOutlet, staff *model.Staff) (filePath string, err error) {
	var file *xlsx.File
	var sheet *xlsx.Sheet
	var row *xlsx.Row

	dir := util.ExportDirectory
	filename := fmt.Sprintf("ReportOutlet_%s_%s.xlsx", time.Now().Format("2006-01-02"), util.GenerateRandomDoc(5))

	fileDir := fmt.Sprintf("%s/%s", dir, filename)
	filePath = fmt.Sprintf("%s/%s", util.ResponseURLFromPostUpload+"report", filename)

	file = xlsx.NewFile()
	if sheet, err = file.AddSheet("Sheet1"); err == nil {
		row = sheet.AddRow()
		row.AddCell().Value = "Outlet Report"
		row = sheet.AddRow()
		row.AddCell().Value = "Download Timestamp: " + date.Format("02/01/2006 15:04 WIB")

		row = sheet.AddRow()
		row.SetHeight(20)
		row.AddCell().Value = "No"
		row.AddCell().Value = "Main Outlet Code"
		row.AddCell().Value = "Outlet Code"
		row.AddCell().Value = "Outlet Name"
		row.AddCell().Value = "PIC Name"
		row.AddCell().Value = "Outlet PIC Phone Number"
		row.AddCell().Value = "Archetype"
		row.AddCell().Value = "Warehouse Default"
		row.AddCell().Value = "Outlet Area"
		row.AddCell().Value = "Province"
		row.AddCell().Value = "City"
		row.AddCell().Value = "District"
		row.AddCell().Value = "Sub District"
		row.AddCell().Value = "Shipping Address"
		row.AddCell().Value = "Price Set"
		row.AddCell().Value = "Salesperson"
		row.AddCell().Value = "Outlet Status"
		row.AddCell().Value = "Created At"
		row.AddCell().Value = "Created By"
		row.AddCell().Value = "Last Updated At"
		row.AddCell().Value = "Last Updated By"

		for i, v := range data {
			row = sheet.AddRow()
			row.AddCell().SetInt(i + 1)
			row.AddCell().Value = v.MerchantCode         // Merchant Code
			row.AddCell().Value = v.OutletCode           // Outlet Code
			row.AddCell().Value = v.OutletName           // Outlet Name
			row.AddCell().Value = v.PicName              // PIC Name
			row.AddCell().Value = v.OutletPicPhoneNumber // Outlet PIC Phone Number
			row.AddCell().Value = v.Archetype            // Archetype
			row.AddCell().Value = v.WarehouseDefault     // Warehouse Default
			row.AddCell().Value = v.OutletArea           // Outlet Area
			row.AddCell().Value = v.Province             // Province
			row.AddCell().Value = v.City                 // City
			row.AddCell().Value = v.District             // District
			row.AddCell().Value = v.SubDistrict          // Sub District
			row.AddCell().Value = v.ShippingAddress      // Shipping Address
			row.AddCell().Value = v.PriceSet             // Price Set
			row.AddCell().Value = v.SalesPerson          // Salesperson
			switch v.OutletStatus {                      // Outlet Status
			case "1":
				row.AddCell().Value = "Active"
			case "2":
				row.AddCell().Value = "Archive"
			case "3":
				row.AddCell().Value = "Deleted"
			default:
				row.AddCell().Value = ""
			}
			row.AddCell().Value = v.CreatedAt     // Created At
			row.AddCell().Value = v.CreatedBy     // Created By
			row.AddCell().Value = v.LastUpdatedAt // Last Updated At
			row.AddCell().Value = v.LastUpdatedBy // Last Updated By
		}

	}
	err = file.Save(fileDir)
	filePath, err = util.UploadToS3(filename, fileDir, "application/xlsx")

	// fungsi ini berguna untuk menghapus kembali file yang sudah diupload ke server
	os.Remove(fileDir)

	if err = log.AuditLogByUser(staff, 0, "branch", "Download", "Download Outlet"); err != nil {
		return "", err
	}
	return
}

// getAgentXls : function to create excel file of agent report
func getAgentXls(date time.Time, data []*reportAgent, staff *model.Staff, note string) (filePath string, err error) {
	var file *xlsx.File
	var sheet *xlsx.Sheet
	var row *xlsx.Row

	dir := util.ExportDirectory
	filename := fmt.Sprintf("ReportAgent_%s_%s.xlsx", time.Now().Format("2006-01-02"), util.GenerateRandomDoc(5))

	fileDir := fmt.Sprintf("%s/%s", dir, filename)
	filePath = fmt.Sprintf("%s/%s", util.ResponseURLFromPostUpload+"report", filename)

	file = xlsx.NewFile()
	if sheet, err = file.AddSheet("Sheet1"); err == nil {
		row = sheet.AddRow()
		row.AddCell().Value = "Agent Report"
		row = sheet.AddRow()
		row.AddCell().Value = "Download Timestamp: " + date.Format("02/01/2006 15:04 WIB")

		row = sheet.AddRow()
		row.SetHeight(20)
		row.AddCell().Value = "No"
		row.AddCell().Value = "Agent Code"
		row.AddCell().Value = "Agent Name"
		row.AddCell().Value = "Main Outlet Name"
		row.AddCell().Value = "PIC Name"
		row.AddCell().Value = "PIC Phone Number"
		row.AddCell().Value = "Email"
		row.AddCell().Value = "Current EdenPoint"
		row.AddCell().Value = "Default Payment Term"
		row.AddCell().Value = "Default Invoice Term"
		row.AddCell().Value = "Payment Group"
		row.AddCell().Value = "Business Type"
		row.AddCell().Value = "Business Type Credit Limit"
		row.AddCell().Value = "Credit Limit Amount"
		row.AddCell().Value = "Credit Limit Remaining"
		row.AddCell().Value = "Suspended"
		row.AddCell().Value = "Archetype"
		row.AddCell().Value = "Customer Tag"
		row.AddCell().Value = "Warehouse Default"
		row.AddCell().Value = "Agent Area"
		row.AddCell().Value = "Address Area"
		row.AddCell().Value = "Province"
		row.AddCell().Value = "City"
		row.AddCell().Value = "District"
		row.AddCell().Value = "Sub District"
		row.AddCell().Value = "Shipping Address"
		row.AddCell().Value = "Price Set"
		row.AddCell().Value = "Salesperson"
		row.AddCell().Value = "Agent Status"
		row.AddCell().Value = "Address Status"
		row.AddCell().Value = "Created at"
		row.AddCell().Value = "Created by"
		row.AddCell().Value = "Last Updated at"
		row.AddCell().Value = "Last Updated by"

		for i, v := range data {
			row = sheet.AddRow()
			row.AddCell().SetInt(i + 1)
			row.AddCell().Value = v.AgentCode          // Agent Code
			row.AddCell().Value = v.AgentName          // Agent Name
			row.AddCell().Value = v.MainOutletName     // Main Outlet Name
			row.AddCell().Value = v.PicName            // PIC Name
			row.AddCell().Value = v.PicPhoneNumber     // PIC Phone Number
			row.AddCell().Value = v.AgentEmail         // Email
			row.AddCell().Value = v.CurrentEdenPoint   // Current EdenPoint
			row.AddCell().Value = v.DefaultPaymentTerm // Default Payment Term
			row.AddCell().Value = v.DefaultInvoiceTerm // Default Invoice Term
			row.AddCell().Value = v.PaymentGroup       // Payment Group
			row.AddCell().Value = v.BusinessType       // Business Type
			if v.BusinessTypeCreditLimit == "1" {
				row.AddCell().Value = "Badan Usaha" // Business Type Credit Limit
			} else {
				row.AddCell().Value = "Personal"
			}
			row.AddCell().Value = v.CreditLimitAmount          //	 Credit Limit Amount
			row.AddCell().Value = v.RemainingCreditLimitAmount // Credit Limit Remaining
			if v.Suspended == "1" {
				row.AddCell().Value = "Yes" // Suspended
			} else {
				row.AddCell().Value = "No"
			}
			row.AddCell().Value = v.Archetype        // Archetype
			row.AddCell().Value = v.CustomerTag      // Customer Tag
			row.AddCell().Value = v.WarehouseDefault // Warehouse Default
			row.AddCell().Value = v.AgentArea        // Agent Area
			row.AddCell().Value = v.AddressArea      // Agent Area
			row.AddCell().Value = v.Province         // Province
			row.AddCell().Value = v.City             // City
			row.AddCell().Value = v.District         // District
			row.AddCell().Value = v.SubDistrict      // Sub District
			row.AddCell().Value = v.ShippingAddress  // Shipping Address
			row.AddCell().Value = v.PriceSet         // Price Set
			row.AddCell().Value = v.SalesPerson      // Salesperson
			switch v.AgentStatus {                   // Agent Status
			case "1":
				row.AddCell().Value = "Active"
			case "2":
				row.AddCell().Value = "Archive"
			case "3":
				row.AddCell().Value = "Deleted"
			default:
				row.AddCell().Value = ""
			}
			switch v.AddressStatus { // Address Status
			case "1":
				row.AddCell().Value = "Active"
			case "2":
				row.AddCell().Value = "Archive"
			case "3":
				row.AddCell().Value = "Deleted"
			default:
				row.AddCell().Value = ""
			}
			row.AddCell().Value = v.CreatedAt     // Created At
			row.AddCell().Value = v.CreatedBy     // Created By
			row.AddCell().Value = v.LastUpdatedAt // Last Updated At
			row.AddCell().Value = v.LastUpdatedBy // Last Updated By
		}

	}
	err = file.Save(fileDir)
	filePath, err = util.UploadToS3(filename, fileDir, "application/xlsx")

	// fungsi ini berguna untuk menghapus kembali file yang sudah diupload ke server
	os.Remove(fileDir)

	if err = log.AuditLogByUser(staff, 0, "agent", "export", note); err != nil {
		return "", err
	}

	return
}

// getVoucherLogXls : function to create excel file of voucher log report
func getVoucherLogXls(date time.Time, data []*reportVoucherLog) (filePath string, err error) {
	var file *xlsx.File
	var sheet *xlsx.Sheet
	var row *xlsx.Row

	dir := util.ExportDirectory
	filename := fmt.Sprintf("ReportVoucherLog_%s_%s.xlsx", time.Now().Format("2006-01-02"), util.GenerateRandomDoc(5))

	fileDir := fmt.Sprintf("%s/%s", dir, filename)
	filePath = fmt.Sprintf("%s/%s", util.ResponseURLFromPostUpload+"report", filename)

	file = xlsx.NewFile()
	if sheet, err = file.AddSheet("Sheet1"); err == nil {
		row = sheet.AddRow()
		row.AddCell().Value = "Voucher Log Report"
		row = sheet.AddRow()
		row.AddCell().Value = "Download Timestamp: " + date.Format("02/01/2006 15:04 WIB")

		row = sheet.AddRow()
		row.SetHeight(20)
		row.AddCell().Value = "No"
		row.AddCell().Value = "Voucher Code"
		row.AddCell().Value = "Redeem Code"
		row.AddCell().Value = "Redeem Date"
		row.AddCell().Value = "Voucher Type"
		row.AddCell().Value = "Order Code"
		row.AddCell().Value = "Order Date"
		row.AddCell().Value = "Order Status"
		row.AddCell().Value = "Discount Amount"
		row.AddCell().Value = "Total Order"
		row.AddCell().Value = "Main Outlet / Agent Code"
		row.AddCell().Value = "Main Outlet / Agent Name"
		row.AddCell().Value = "Outlet / Address Code"
		row.AddCell().Value = "Outlet / Address Name"
		row.AddCell().Value = "Area"
		row.AddCell().Value = "Business Type"
		row.AddCell().Value = "Archetype"

		for i, v := range data {
			row = sheet.AddRow()
			row.AddCell().SetInt(i + 1)
			row.AddCell().Value = v.VoucherCode                    // Voucher Code
			row.AddCell().Value = v.RedeemCode                     // Redeem Code
			row.AddCell().Value = v.RedeemDate                     // Redeem Date
			row.AddCell().Value = v.VoucherType                    // Voucher Type
			row.AddCell().Value = v.OrderCode                      // Order Code
			row.AddCell().Value = v.OrderDate                      // Order Date
			row.AddCell().Value = v.OrderStatus                    // Order Status
			row.AddCell().SetFloatWithFormat(v.DiscAmount, "0.00") // Discount Amount
			row.AddCell().SetFloatWithFormat(v.TotalOrder, "0.00") // Total Order
			row.AddCell().Value = v.MerchantCode                   // Main Outlet / Agent Code
			row.AddCell().Value = v.MerchantName                   // Main Outlet / Agent Name
			row.AddCell().Value = v.OutletCode                     // Outlet / Address Code
			row.AddCell().Value = v.OutletName                     // Outlet / Address Name
			row.AddCell().Value = v.Area                           // Area
			row.AddCell().Value = v.BusinessType                   // Business Type
			row.AddCell().Value = v.Archetype                      // Archetype
		}

	}
	err = file.Save(fileDir)
	filePath, err = util.UploadToS3(filename, fileDir, "application/xlsx")

	// fungsi ini berguna untuk menghapus kembali file yang sudah diupload ke server
	os.Remove(fileDir)
	return
}

// getSubmissionXls : function to create excel file of submission report
func getSubmissionXls(date time.Time, data []*reportSubmission, staff *model.Staff) (filePath string, err error) {
	var file *xlsx.File
	var sheet *xlsx.Sheet
	var row *xlsx.Row

	dir := util.ExportDirectory
	filename := fmt.Sprintf("ReportSalesTask_%s_%s.xlsx", time.Now().Format("2006-01-02"), util.GenerateRandomDoc(5))

	fileDir := fmt.Sprintf("%s/%s", dir, filename)
	filePath = fmt.Sprintf("%s/%s", util.ResponseURLFromPostUpload+"report", filename)

	file = xlsx.NewFile()
	if sheet, err = file.AddSheet("Sheet1"); err == nil {
		row = sheet.AddRow()
		row.AddCell().Value = "Submission Report"
		row = sheet.AddRow()
		row.AddCell().Value = "Download Timestamp: " + date.Format("02/01/2006 15:04 WIB")

		row = sheet.AddRow()
		row.SetHeight(20)
		row.AddCell().Value = "No"
		row.AddCell().Value = "Sales Group"
		row.AddCell().Value = "Type"
		row.AddCell().Value = "Out of Route"
		row.AddCell().Value = "Customer Type"
		row.AddCell().Value = "Assignment Date"
		row.AddCell().Value = "Submission Date"
		row.AddCell().Value = "Finish Date"
		row.AddCell().Value = "Salesperson"
		row.AddCell().Value = "Customer Code"
		row.AddCell().Value = "Customer Name"
		row.AddCell().Value = "Phone Number"
		row.AddCell().Value = "Customer Location"
		row.AddCell().Value = "Actual FS Location"
		row.AddCell().Value = "Actual Distance"
		row.AddCell().Value = "Address"
		row.AddCell().Value = "Result"
		row.AddCell().Value = "Food Apps"
		row.AddCell().Value = "Status"
		row.AddCell().Value = "Effective Call"
		row.AddCell().Value = "Revenue Effective Call"
		row.AddCell().Value = "Objective Codes"
		row.AddCell().Value = "Objective Names"

		for i, v := range data {
			row = sheet.AddRow()
			row.AddCell().SetInt(i + 1)
			row.AddCell().Value = v.Salesgroup                                                                             //Sales Group
			row.AddCell().Value = v.Type                                                                                   //Type
			row.AddCell().Value = v.OutOfRouteStr                                                                          //Out of Route
			row.AddCell().Value = v.CustomerTypeStr                                                                        //Customer Type
			row.AddCell().Value = v.StartDate + " to " + v.EndDate                                                         //Assignment Date to End Date
			row.AddCell().Value = v.SubmissionDate                                                                         //Submission Date
			row.AddCell().Value = v.FinishDate                                                                             //Finish Date
			row.AddCell().Value = v.Salesperson                                                                            //Salesperson
			row.AddCell().Value = v.OutletCode                                                                             //Customer Code
			row.AddCell().Value = v.OutletName                                                                             //Customer Name
			row.AddCell().Value = v.PhoneNumber                                                                            //Phone Number
			row.AddCell().Value = fmt.Sprintf("%g", v.Latitude) + "," + fmt.Sprintf("%g", v.Longitude)                     //Customer Location
			row.AddCell().Value = fmt.Sprintf("%g", v.ActualTaskLatitude) + "," + fmt.Sprintf("%g", v.ActualTaskLongitude) //Actual FS Location
			row.AddCell().Value = fmt.Sprintf("%g", v.ActualDistance)                                                      //Actual Distance
			row.AddCell().Value = v.ShippingAddress                                                                        //Address
			row.AddCell().Value = v.Result                                                                                 //Result
			row.AddCell().Value = v.FoodApp                                                                                //Food Apps
			row.AddCell().Value = v.StatusStr                                                                              //Status
			row.AddCell().Value = strconv.FormatBool(v.EffectiveCall)                                                      //Effective Call
			row.AddCell().SetFloatWithFormat(v.RevenueEffectiveCall, "0.00")                                               //Revenue Effective Call
			row.AddCell().Value = v.ObjectiveCodes                                                                         //Objective Codes
			row.AddCell().Value = v.ObjectiveCodesStr                                                                      //Objective Names
		}
	}
	err = file.Save(fileDir)
	filePath, err = util.UploadToS3(filename, fileDir, "application/xlsx")

	// fungsi ini berguna untuk menghapus kembali file yang sudah diupload ke server
	os.Remove(fileDir)

	if err = log.AuditLogByUser(staff, 0, "sales_assignment_item", "Download", "Download Submission"); err != nil {
		return "", err
	}
	return
}
