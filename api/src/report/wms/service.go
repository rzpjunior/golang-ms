// Copyright 2020 PT. Eden Pangan Indonesia. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

package wms

import (
	"encoding/json"
	"fmt"
	"os"
	"time"

	"git.edenfarm.id/cuxs/mongodb"
	"git.edenfarm.id/project-version2/api/datastore/repository"
	"go.mongodb.org/mongo-driver/bson"

	"git.edenfarm.id/project-version2/api/util"
	"git.edenfarm.id/project-version2/datamodel/model"
	"github.com/tealeg/xlsx"
)

// GetStockLogXls : function to create excel file of stock log report
func GetStockLogXls(data []*reportStockLog, warehouse *model.Warehouse) (filePath string, err error) {
	var file *xlsx.File
	var sheet *xlsx.Sheet
	var row *xlsx.Row
	var filename string

	dir := util.ExportDirectory
	if warehouse.Name == "All Warehouse" {
		filename = fmt.Sprintf("ReportStockLog_%s_%s.xlsx", time.Now().Format("2006-01-02"), util.GenerateRandomDoc(5))
	} else {
		filename = fmt.Sprintf("ReportStockLog_%s_%s_%s.xlsx", time.Now().Format("2006-01-02"), util.ReplaceSpace(warehouse.Name), util.GenerateRandomDoc(5))
	}

	fileDir := fmt.Sprintf("%s/%s", dir, filename)

	file = xlsx.NewFile()
	if sheet, err = file.AddSheet("Sheet1"); err == nil {
		row = sheet.AddRow()
		row.SetHeight(20)
		row.AddCell().Value = "Timestamp"
		row.AddCell().Value = "Log_Type"
		row.AddCell().Value = "Ref_Type"
		row.AddCell().Value = "Reference_Code"
		row.AddCell().Value = "Product_Code"
		row.AddCell().Value = "Product_Name"
		row.AddCell().Value = "UOM"
		row.AddCell().Value = "Initial_Stock"
		row.AddCell().Value = "Quantity"
		row.AddCell().Value = "Final_Stock"
		row.AddCell().Value = "Warehouse"
		row.AddCell().Value = "Area"
		row.AddCell().Value = "Status"
		row.AddCell().Value = "Doc_Note"
		row.AddCell().Value = "Note"

		var RefType string
		var LogType string
		var Status string
		for _, v := range data {

			switch v.RefType {
			case "1":
				RefType = "delivery order"
			case "2":
				RefType = "delivery return"
			case "3":
				RefType = "goods receipt"
			case "4":
				RefType = "goods transfer"
			case "5":
				RefType = "stock opname"
			case "6":
				RefType = "waste entry"
			case "7":
				RefType = "transfer sku"
			}

			switch v.LogType {
			case "1":
				LogType = "IN"
			case "2":
				LogType = "OUT"
			}

			switch v.Status {
			case "1":
				Status = "Active"
			case "2":
				Status = "Archived"
			}

			row = sheet.AddRow()
			row.AddCell().Value = v.TimeStamp     // Timestamp
			row.AddCell().Value = LogType         // Log_Type
			row.AddCell().Value = RefType         // Ref_Type
			row.AddCell().Value = v.ReferenceCode // Reference_Code
			row.AddCell().Value = v.ProductCode
			row.AddCell().Value = v.ProductName                      // Product_Name
			row.AddCell().Value = v.Uom                              // Uom
			row.AddCell().SetFloatWithFormat(v.InitialStock, "0.00") // Initial_Stock
			row.AddCell().SetFloatWithFormat(v.Quantity, "0.00")     // Quantity
			row.AddCell().SetFloatWithFormat(v.FinalStock, "0.00")   // Final_Stock
			row.AddCell().Value = v.Warehouse                        // Warehouse
			row.AddCell().Value = v.Area                             // Area
			row.AddCell().Value = Status                             // Status
			row.AddCell().Value = v.DocNote                          // Doc_Note
			row.AddCell().Value = v.Note                             // Note
		}

	}
	err = file.Save(fileDir)
	filePath, err = util.UploadToS3(filename, fileDir, "application/xlsx")

	// fungsi ini berguna untuk menghapus kembali file yang sudah diupload ke server
	os.Remove(fileDir)
	return
}

// GetWasteLogXls : function to create excel file of waste log report
func GetWasteLogXls(data []*reportWasteLog, warehouse *model.Warehouse) (filePath string, err error) {
	var file *xlsx.File
	var sheet *xlsx.Sheet
	var row *xlsx.Row
	var filename string

	dir := util.ExportDirectory

	if warehouse.Name == "All Warehouse" {
		filename = fmt.Sprintf("ReportWasteLog_%s_%s.xlsx", time.Now().Format("2006-01-02"), util.GenerateRandomDoc(5))
	} else {
		filename = fmt.Sprintf("ReportWasteLog_%s_%s_%s.xlsx", time.Now().Format("2006-01-02"), util.ReplaceSpace(warehouse.Name), util.GenerateRandomDoc(5))
	}

	fileDir := fmt.Sprintf("%s/%s", dir, filename)

	file = xlsx.NewFile()
	if sheet, err = file.AddSheet("Sheet1"); err == nil {
		row = sheet.AddRow()
		row.SetHeight(20)
		row.AddCell().Value = "Timestamp"
		row.AddCell().Value = "Log_Type"
		row.AddCell().Value = "Ref_Type"
		row.AddCell().Value = "Reference_Code"
		row.AddCell().Value = "Waste_Reason"
		row.AddCell().Value = "Good_Receipt_Code"
		row.AddCell().Value = "Good_Transfer_Code"
		row.AddCell().Value = "Purchase_Order_Code"
		row.AddCell().Value = "Suplier_Name"
		row.AddCell().Value = "Suplier_Type"
		row.AddCell().Value = "Warehouse_Origin"
		row.AddCell().Value = "Product_Code"
		row.AddCell().Value = "Product_Name"
		row.AddCell().Value = "UOM"
		row.AddCell().Value = "Quantity"
		row.AddCell().Value = "Final_Stock"
		row.AddCell().Value = "Warehouse"
		row.AddCell().Value = "Area"
		row.AddCell().Value = "Doc_Note"
		row.AddCell().Value = "Note"

		for _, v := range data {
			row = sheet.AddRow()
			row.AddCell().Value = v.TimeStamp                      // Timestamp
			row.AddCell().Value = v.LogType                        // Log_Type
			row.AddCell().Value = v.RefType                        // Ref_Type
			row.AddCell().Value = v.ReferenceCode                  // Reference_Code
			row.AddCell().Value = v.WasteReason                    // Waste_Reason
			row.AddCell().Value = v.GoodReceiptCode                // Good_Receipt_Code
			row.AddCell().Value = v.GoodTransferCode               // Good_Transfer_Code
			row.AddCell().Value = v.PurchaseOrderCode              // Purchase_Order_Code
			row.AddCell().Value = v.SuplierName                    // Suplier_Name
			row.AddCell().Value = v.SuplierType                    // Suplier_Type
			row.AddCell().Value = v.WarehouseOrigin                // Warehouse_Origin
			row.AddCell().Value = v.ProductCode                    // ProductCode
			row.AddCell().Value = v.ProductName                    // Product_Name
			row.AddCell().Value = v.Uom                            // Uom
			row.AddCell().SetFloatWithFormat(v.Quantity, "0.00")   // Quantity
			row.AddCell().SetFloatWithFormat(v.FinalStock, "0.00") // Final_Stock
			row.AddCell().Value = v.Warehouse                      // Warehouse
			row.AddCell().Value = v.Area                           // Area
			row.AddCell().Value = v.DocNote                        // Doc_Note
			row.AddCell().Value = v.ItemNote                       // Item_Note
		}

	}
	err = file.Save(fileDir)
	filePath, err = util.UploadToS3(filename, fileDir, "application/xlsx")
	// fungsi ini berguna untuk menghapus kembali file yang sudah diupload ke server
	os.Remove(fileDir)
	return
}

// GetStockXls : function to create excel file of stock report
func GetStockXls(data []*reportStock, warehouse *model.Warehouse) (filePath string, err error) {
	var file *xlsx.File
	var sheet *xlsx.Sheet
	var row *xlsx.Row
	var filename string

	dir := util.ExportDirectory
	if warehouse.Name == "All Warehouse" {
		filename = fmt.Sprintf("ReportStocks_%s_%s.xlsx", time.Now().Format("2006-01-02"), util.GenerateRandomDoc(5))
	} else {
		filename = fmt.Sprintf("ReportStocks_%s_%s_%s.xlsx", time.Now().Format("2006-01-02"), util.ReplaceSpace(warehouse.Name), util.GenerateRandomDoc(5))
	}

	fileDir := fmt.Sprintf("%s/%s", dir, filename)

	file = xlsx.NewFile()
	if sheet, err = file.AddSheet("Sheet1"); err == nil {
		row = sheet.AddRow()
		row.SetHeight(20)
		row.AddCell().Value = "Product_Code"
		row.AddCell().Value = "Product_Name"
		row.AddCell().Value = "Warehouse_Name"
		row.AddCell().Value = "Available_Stock"
		row.AddCell().Value = "Waste_Stock"
		row.AddCell().Value = "Safety_Stock"
		row.AddCell().Value = "Commited_In_Stock"
		row.AddCell().Value = "Commited_Out_Stock"
		row.AddCell().Value = "Expected_Stock"
		row.AddCell().Value = "Intransit_Stock"
		row.AddCell().Value = "Received_Stock"
		row.AddCell().Value = "Intransit_Waste_Stock"
		row.AddCell().Value = "Salable"
		row.AddCell().Value = "Purchasable"
		row.AddCell().Value = "Status"

		for _, v := range data {
			row = sheet.AddRow()
			row.AddCell().Value = v.ProductCode                             // ProductCode
			row.AddCell().Value = v.ProductName                             // ProductName
			row.AddCell().Value = v.WarehouseName                           // WarehouseName
			row.AddCell().SetFloatWithFormat(v.AvailableStock, "0.00")      // AvailableStock
			row.AddCell().SetFloatWithFormat(v.WasteStock, "0.00")          // WasteStock
			row.AddCell().SetFloatWithFormat(v.SafetyStock, "0.00")         // SafetyStock
			row.AddCell().SetFloatWithFormat(v.CommitedInStock, "0.00")     // CommitedInStock
			row.AddCell().SetFloatWithFormat(v.CommitedOutStock, "0.00")    // CommitedOutStock
			row.AddCell().SetFloatWithFormat(v.ExpectedStock, "0.00")       // ExpectedStock
			row.AddCell().SetFloatWithFormat(v.IntransitStock, "0.00")      // IntransitStock
			row.AddCell().SetFloatWithFormat(v.ReceivedStock, "0.00")       // ReceivedStock
			row.AddCell().SetFloatWithFormat(v.IntransitWasteStock, "0.00") // IntransitWasteStock
			row.AddCell().Value = v.Salable                                 // Salable
			row.AddCell().Value = v.Purchasable                             // Purchasable
			row.AddCell().Value = v.Status                                  // Status
		}

	}
	err = file.Save(fileDir)
	filePath, err = util.UploadToS3(filename, fileDir, "application/xlsx")
	// fungsi ini berguna untuk menghapus kembali file yang sudah diupload ke server
	os.Remove(fileDir)
	return
}

// GetGoodsReceiptItemXls : function to create excel file of goods receipt item report
func GetGoodsReceiptItemXls(data []*reportGoodsReceiptItem, warehouse *model.Warehouse) (filePath string, err error) {
	var file *xlsx.File
	var sheet *xlsx.Sheet
	var row *xlsx.Row
	var filename string

	dir := util.ExportDirectory
	if warehouse.Name == "All Warehouse" {
		filename = fmt.Sprintf("ReportGoodsReceiptItem_%s_%s.xlsx", time.Now().Format("2006-01-02"), util.GenerateRandomDoc(5))
	} else {
		filename = fmt.Sprintf("ReportGoodsReceiptItem_%s_%s_%s.xlsx", time.Now().Format("2006-01-02"), util.ReplaceSpace(warehouse.Name), util.GenerateRandomDoc(5))
	}

	fileDir := fmt.Sprintf("%s/%s", dir, filename)

	file = xlsx.NewFile()
	if sheet, err = file.AddSheet("Sheet1"); err == nil {
		row = sheet.AddRow()
		row.SetHeight(20)
		row.AddCell().Value = "Inbound_Code"
		row.AddCell().Value = "Supplier_Code"
		row.AddCell().Value = "Supplier_Name"
		row.AddCell().Value = "Inbound_Status"
		row.AddCell().Value = "GR_Code"
		row.AddCell().Value = "GR_Status"
		row.AddCell().Value = "SR_Code"
		row.AddCell().Value = "SR_Status"
		row.AddCell().Value = "DN_Code"
		row.AddCell().Value = "DN_Status"
		row.AddCell().Value = "Warehouse_Origin"
		row.AddCell().Value = "Warehouse_Destination"
		row.AddCell().Value = "Estimation_Arrival_Date"
		row.AddCell().Value = "Estimation_Arrival_Time"
		row.AddCell().Value = "Actual_Arrival_Date"
		row.AddCell().Value = "Actual_Arrival_Time"
		row.AddCell().Value = "GR_Note"
		row.AddCell().Value = "Product_Code"
		row.AddCell().Value = "Product_Name"
		row.AddCell().Value = "UOM"
		row.AddCell().Value = "GR_Item_Note"
		row.AddCell().Value = "Ordered_Qty"
		row.AddCell().Value = "Delivered_Qty"
		row.AddCell().Value = "Reject_Qty"
		row.AddCell().Value = "Received_Qty"
		row.AddCell().Value = "Return_Qty"
		row.AddCell().Value = "After_Sortir_Qty(Good)"
		row.AddCell().Value = "After_Sortir_Qty(Waste)"
		row.AddCell().Value = "After_Sortir_Qty(Down_Grade)"
		row.AddCell().Value = "Product_DownGrade"
		row.AddCell().Value = "TS_Code"
		row.AddCell().Value = "TS_Status"

		for _, v := range data {
			row = sheet.AddRow()
			row.AddCell().Value = v.InboundCode                                 // InboundCode
			row.AddCell().Value = v.SupplierCode                                // SupplierCode
			row.AddCell().Value = v.SupplierName                                // SupplierName
			row.AddCell().Value = v.InboundStatus                               // InboundStatus
			row.AddCell().Value = v.GRCode                                      // GRCode
			row.AddCell().Value = v.GRStatus                                    // GRStatus
			row.AddCell().Value = v.SupplierReturnCode                          // SR_Code
			row.AddCell().Value = v.SupplierReturnStatus                        // SR_Status
			row.AddCell().Value = v.DebitNoteCode                               // DN_Code
			row.AddCell().Value = v.DebitNoteStatus                             // DN_Status
			row.AddCell().Value = v.WarehouseOrigin                             // WarehouseOrigin
			row.AddCell().Value = v.WarehouseDestination                        // WarehouseDestination
			row.AddCell().Value = v.EstimationArrivalDate                       // EstimationArrivalDate
			row.AddCell().Value = v.EstimationArrivalTime                       // EstimationArrivalTime
			row.AddCell().Value = v.ActualArrivalDate                           // ActualArrivalDate
			row.AddCell().Value = v.ActualArrivalTime                           // ActualArrivalTime
			row.AddCell().Value = v.GRNote                                      // GRNote
			row.AddCell().Value = v.ProductCode                                 // ProductCode
			row.AddCell().Value = v.ProductName                                 // ProductName
			row.AddCell().Value = v.UOM                                         // UOM
			row.AddCell().Value = v.GRItemNote                                  // GRItemNote
			row.AddCell().SetFloatWithFormat(v.OrderedQty, "0.00")              // OrderedQty
			row.AddCell().SetFloatWithFormat(v.DeliveredQty, "0.00")            // DeliveredQty
			row.AddCell().SetFloatWithFormat(v.RejectQty, "0.00")               // RejectQty
			row.AddCell().SetFloatWithFormat(v.ReceivedQty, "0.00")             // ReceivedQty
			row.AddCell().SetFloatWithFormat(v.SRIReturnQuantity, "0.00")       // ReturnQty
			row.AddCell().SetFloatWithFormat(v.AfterSortirGoodQty, "0.00")      // GoodQty
			row.AddCell().SetFloatWithFormat(v.AfterSortirWasteQty, "0.00")     // WasteQty
			row.AddCell().SetFloatWithFormat(v.AfterSortirDownGradeQty, "0.00") // DownGrade
			row.AddCell().Value = v.Product_DownGrade                           // Product_DownGrade
			row.AddCell().Value = v.TS_Code                                     // TS_Code
			row.AddCell().Value = v.TS_Status                                   // TS_Status

		}

	}
	err = file.Save(fileDir)
	filePath, err = util.UploadToS3(filename, fileDir, "application/xlsx")

	// fungsi ini berguna untuk menghapus kembali file yang sudah diupload ke server
	os.Remove(fileDir)
	return
}

// GetDeliveryReturnItemXls : function to create excel file of delivery return item report
func GetDeliveryReturnItemXls(data []*reportDeliveryReturnItem, warehouse *model.Warehouse) (filePath string, err error) {
	var file *xlsx.File
	var sheet *xlsx.Sheet
	var row *xlsx.Row
	var filename string

	dir := util.ExportDirectory
	if warehouse.Name == "All Warehouse" {
		filename = fmt.Sprintf("ReportDeliveryReturnItem_%s_%s.xlsx", time.Now().Format("2006-01-02"), util.GenerateRandomDoc(5))
	} else {
		filename = fmt.Sprintf("ReportDeliveryReturnItem_%s_%s_%s.xlsx", time.Now().Format("2006-01-02"), util.ReplaceSpace(warehouse.Name), util.GenerateRandomDoc(5))
	}

	fileDir := fmt.Sprintf("%s/%s", dir, filename)

	file = xlsx.NewFile()
	if sheet, err = file.AddSheet("Sheet1"); err == nil {
		row = sheet.AddRow()
		row.SetHeight(20)
		row.AddCell().Value = "Return_Date"
		row.AddCell().Value = "Product_Code"
		row.AddCell().Value = "Product_Name"
		row.AddCell().Value = "Unit"
		row.AddCell().Value = "Good_Stock_Return_Qty"
		row.AddCell().Value = "Waste_Return_Qty"
		row.AddCell().Value = "Total_Return_Qty"
		row.AddCell().Value = "Product_Price"
		row.AddCell().Value = "Area"
		row.AddCell().Value = "Warehouse"
		row.AddCell().Value = "Order_Code"
		row.AddCell().Value = "Delivery_Code"
		row.AddCell().Value = "Delivery_Date"
		row.AddCell().Value = "Customer_Code"
		row.AddCell().Value = "Customer_Name"
		row.AddCell().Value = "Delivery_Return_Note"
		row.AddCell().Value = "Delivery_Return_Item_Note"

		for _, v := range data {
			row = sheet.AddRow()
			row.AddCell().Value = v.ReturnDate                             // ReturnDate
			row.AddCell().Value = v.ProductCode                            // ProductCode
			row.AddCell().Value = v.ProductName                            // ProductName
			row.AddCell().Value = v.Unit                                   // Unit
			row.AddCell().SetFloatWithFormat(v.GoodStockReturnQty, "0.00") // GoodStockReturnQty
			row.AddCell().SetFloatWithFormat(v.WasteReturnQty, "0.00")     // WasteReturnQty
			row.AddCell().SetFloatWithFormat(v.TotalReturnQty, "0.00")     // TotalReturnQty
			row.AddCell().SetFloatWithFormat(v.ProductPrice, "0.00")       // ProductPrice
			row.AddCell().Value = v.Area                                   // Area
			row.AddCell().Value = v.Warehouse                              // Warehouse
			row.AddCell().Value = v.OrderCode                              // OrderCode
			row.AddCell().Value = v.DeliveryCode                           // DeliveryCode
			row.AddCell().Value = v.DeliveryDate                           // DeliveryDate
			row.AddCell().Value = v.CustomerCode                           // CustomerCode
			row.AddCell().Value = v.CustomerName                           // CustomerName
			row.AddCell().Value = v.DeliveryReturnNote                     // DeliveryReturnNote
			row.AddCell().Value = v.DeliveryReturnItemNote                 // DeliveryReturnItemNote
		}

	}
	err = file.Save(fileDir)
	filePath, err = util.UploadToS3(filename, fileDir, "application/xlsx")
	// fungsi ini berguna untuk menghapus kembali file yang sudah diupload ke server
	os.Remove(fileDir)
	return
}

// GetProductsXls : function to create excel file of products report
func GetProductsXls(data []*reportProducts) (filePath string, err error) {
	var file *xlsx.File
	var sheet *xlsx.Sheet
	var row *xlsx.Row

	dir := util.ExportDirectory
	filename := fmt.Sprintf("ReportProducts_%s_%s.xlsx", time.Now().Format("2006-01-02"), util.GenerateRandomDoc(5))

	fileDir := fmt.Sprintf("%s/%s", dir, filename)

	file = xlsx.NewFile()
	if sheet, err = file.AddSheet("Sheet1"); err == nil {
		row = sheet.AddRow()
		row.SetHeight(20)
		row.AddCell().Value = "Product_Code"
		row.AddCell().Value = "Product_Name"
		row.AddCell().Value = "Category"
		row.AddCell().Value = "Category_Code"
		row.AddCell().Value = "Parent"
		row.AddCell().Value = "Parent_Code"
		row.AddCell().Value = "Grand_Parent"
		row.AddCell().Value = "Grand_Parent_Code"
		row.AddCell().Value = "UOM"
		row.AddCell().Value = "Total_Weight"
		row.AddCell().Value = "Min_Qty_Order"
		row.AddCell().Value = "Product_Note"
		row.AddCell().Value = "Product_Description"
		row.AddCell().Value = "Product_Tag"
		row.AddCell().Value = "Product_Status"
		row.AddCell().Value = "Warehouse_Salability"
		row.AddCell().Value = "Warehouse_Purchasability"
		row.AddCell().Value = "Warehouse_Storability"
		row.AddCell().Value = "Spare Percentage"

		for _, v := range data {
			row = sheet.AddRow()
			row.AddCell().Value = v.ProductCode                         // ProductCode
			row.AddCell().Value = v.ProductName                         // ProductName
			row.AddCell().Value = v.Category                            // Category
			row.AddCell().Value = v.CategoryCode                        // Category
			row.AddCell().Value = v.Parent                              // Parent
			row.AddCell().Value = v.ParentCode                          // Parent
			row.AddCell().Value = v.GrandParent                         // GrandParent
			row.AddCell().Value = v.GrandParentCode                     // GrandParent
			row.AddCell().Value = v.UOM                                 // UOM
			row.AddCell().SetFloatWithFormat(v.TotalWeight, "0.00")     // TotalWeight
			row.AddCell().SetFloatWithFormat(v.MinimalOrderQty, "0.00") // MinimalOrderQty
			row.AddCell().Value = v.ProductNote                         // ProductNote
			row.AddCell().Value = v.ProductDescription                  // ProductDescription
			row.AddCell().Value = v.ProductTag                          // ProductTag
			row.AddCell().Value = v.ProductStatus                       // ProductStatus
			row.AddCell().Value = v.WarehouseSalability                 // WarehouseSalability
			row.AddCell().Value = v.WarehousePurchasability             // WarehousePurchasability
			row.AddCell().Value = v.WarehouseStorability                // WarehouseStorability
			row.AddCell().SetFloatWithFormat(v.SparePercentage, "0.00") // Spare Percentage
		}

	}
	err = file.Save(fileDir)
	filePath, err = util.UploadToS3(filename, fileDir, "application/xlsx")

	// fungsi ini berguna untuk menghapus kembali file yang sudah diupload ke server
	os.Remove(fileDir)
	return
}

func getDeliveryOrderXls(date time.Time, r []*reportDeliveryOrder, warehouse *model.Warehouse) (filePath string, err error) {
	var file *xlsx.File
	var sheet *xlsx.Sheet
	var row *xlsx.Row
	var fileName string
	dir := util.ExportDirectory

	if warehouse != nil {
		fileName = fmt.Sprintf("ReportDeliveryOrder_%s_%s_%s.xlsx", date.Format("2006-01-02"), util.ReplaceSpace(warehouse.Name), util.GenerateRandomDoc(5))
	}

	fileName = fmt.Sprintf("ReportDeliveryOrder_%s_%s.xlsx", date.Format("2006-01-02"), util.GenerateRandomDoc(5))

	fileDir := fmt.Sprintf("%s/%s", dir, fileName)

	file = xlsx.NewFile()
	if sheet, err = file.AddSheet("Sheet1"); err == nil {
		row = sheet.AddRow()
		row.SetHeight(20)
		row.AddCell().Value = "No"
		row.AddCell().Value = "Warehouse"
		row.AddCell().Value = "Business_Type"
		row.AddCell().Value = "Order_Code"
		row.AddCell().Value = "Order_Type"
		row.AddCell().Value = "Merchant_Name"
		row.AddCell().Value = "Order_Status"
		row.AddCell().Value = "Delivery_Code"
		row.AddCell().Value = "Delivery_Status"
		row.AddCell().Value = "Shipping_Address"
		row.AddCell().Value = "Province"
		row.AddCell().Value = "City"
		row.AddCell().Value = "District"
		row.AddCell().Value = "Sub District"
		row.AddCell().Value = "Postal_Code"
		row.AddCell().Value = "WRT"
		row.AddCell().Value = "Order_Weight"
		row.AddCell().Value = "Delivery_Date"
		row.AddCell().Value = "Payment_Term"
		row.AddCell().Value = "Tag_Customer"
		row.AddCell().Value = "Area"

		for i, v := range r {

			row = sheet.AddRow()
			row.AddCell().SetInt(i + 1)
			row.AddCell().Value = v.Warehouse    // Warehouse
			row.AddCell().Value = v.BusinessType // Business Type
			row.AddCell().Value = v.OrderCode    // Order Code
			row.AddCell().Value = v.OrderType    // Order Type
			row.AddCell().Value = v.MerchantName // Merchant Name
			row.AddCell().Value = v.OrderStatus  // Order Status

			if v.DeliveryCode != "" {
				row.AddCell().Value = v.DeliveryCode // Delivery Code
			} else {
				row.AddCell().Value = "-" // Delivery Code
			}

			if v.DeliveryStatus != "" {
				row.AddCell().Value = v.DeliveryStatus // Delivery Status
			} else {
				row.AddCell().Value = "-" // Delivery Status
			}
			row.AddCell().Value = v.ShippingAddress                 // Shipping Address
			row.AddCell().Value = v.Province                        // Province
			row.AddCell().Value = v.City                            // City
			row.AddCell().Value = v.District                        // District
			row.AddCell().Value = v.SubDistrict                     // Sub District
			row.AddCell().Value = v.PostalCode                      // Postal Code
			row.AddCell().Value = v.Wrt                             // WRT
			row.AddCell().SetFloatWithFormat(v.OrderWeight, "0.00") // Order Weight
			row.AddCell().Value = v.DeliveryDate                    // Delivery Date
			row.AddCell().Value = v.PaymentTerm                     // Payment Term
			row.AddCell().Value = v.TagCustomer                     // Tag Customer
			row.AddCell().Value = v.AreaName                        // Area
		}

	}
	err = file.Save(fileDir)
	filePath, err = util.UploadToS3(fileName, fileDir, "application/xlsx")
	// fungsi ini berguna untuk menghapus kembali file yang sudah diupload ke server
	os.Remove(fileDir)
	return
}

// GetItemRecapXls : function to create excel file of item recap report
func GetItemRecapXls(data []*reportItemRecap, area *model.Area) (filePath string, err error) {
	var file *xlsx.File
	var sheet *xlsx.Sheet
	var row *xlsx.Row
	var filename string

	dir := util.ExportDirectory
	filename = fmt.Sprintf("ReportItemRecap_%s_%s_%s.xlsx", time.Now().Format("2006-01-02"), util.ReplaceSpace(area.Name), util.GenerateRandomDoc(5))

	fileDir := fmt.Sprintf("%s/%s", dir, filename)

	file = xlsx.NewFile()
	if sheet, err = file.AddSheet("Sheet1"); err == nil {
		row = sheet.AddRow()
		row.SetHeight(20)
		row.AddCell().Value = "Order_Delivery_Date"
		row.AddCell().Value = "Area"
		row.AddCell().Value = "Warehouse"
		row.AddCell().Value = "Product_Code"
		row.AddCell().Value = "Product_Name"
		row.AddCell().Value = "UOM"
		row.AddCell().Value = "Total_Quantity"
		row.AddCell().Value = "Total_Quantity_Zero_Waste"
		row.AddCell().Value = "Total_Weight"
		row.AddCell().Value = "Category"
		row.AddCell().Value = "Category_Code"
		row.AddCell().Value = "Parent"
		row.AddCell().Value = "Parent_Code"
		row.AddCell().Value = "Grand_Parent"
		row.AddCell().Value = "Grand_Parent_Code"

		for _, v := range data {
			row = sheet.AddRow()
			row.AddCell().Value = v.DeliveryDate                          // Order_Delivery_Date
			row.AddCell().Value = v.Area                                  // Area
			row.AddCell().Value = v.Warehouse                             // Warehouse
			row.AddCell().Value = v.ProductCode                           // Product_Code
			row.AddCell().Value = v.ProductName                           // Product_Name
			row.AddCell().Value = v.Uom                                   // UOM
			row.AddCell().SetFloatWithFormat(v.TotalQty, "0.00")          // Total_Quantity
			row.AddCell().SetFloatWithFormat(v.TotalQtyZeroWaste, "0.00") // Total_Quantity_Zero_Waste
			row.AddCell().SetFloatWithFormat(v.TotalWeight, "0.00")       // Total_Weight
			row.AddCell().Value = v.Category                              // Category
			row.AddCell().Value = v.CategoryCode                          // CategoryCode
			row.AddCell().Value = v.Parent                                // Parent
			row.AddCell().Value = v.ParentCode                            // ParentCode
			row.AddCell().Value = v.GrandParent                           // GrandParent
			row.AddCell().Value = v.GrandParentCode                       // GrandParentCode
		}

	}
	err = file.Save(fileDir)
	filePath, err = util.UploadToS3(filename, fileDir, "application/xlsx")
	// fungsi ini berguna untuk menghapus kembali file yang sudah diupload ke server
	os.Remove(fileDir)
	return
}

// GetMovementStockXls : function to create excel file of movement stock report
func GetMovementStockXls(data []*reportMovementStock, warehouse *model.Warehouse) (filePath string, err error) {
	var file *xlsx.File
	var sheet *xlsx.Sheet
	var row *xlsx.Row

	dir := util.ExportDirectory
	filename := fmt.Sprintf("ReportMovementStock_%s_%s_%s.xlsx", time.Now().Format("2006-01-02"), util.ReplaceSpace(warehouse.Name), util.GenerateRandomDoc(5))

	fileDir := fmt.Sprintf("%s/%s", dir, filename)

	file = xlsx.NewFile()
	if sheet, err = file.AddSheet("Sheet1"); err == nil {
		row = sheet.AddRow()
		row.SetHeight(20)
		row.AddCell().Value = "Product_Code"
		row.AddCell().Value = "Product_Name"
		row.AddCell().Value = "Category"
		row.AddCell().Value = "UOM"
		row.AddCell().Value = "Stock"
		row.AddCell().Value = "Plan_Inbound"
		row.AddCell().Value = "Actual_Inbound"
		row.AddCell().Value = "Waste"
		row.AddCell().Value = "Plan_Delivery"
		row.AddCell().Value = "Actual_Delivery"
		row.AddCell().Value = "Stock_Transfer(In)"
		row.AddCell().Value = "Stock_Transfer(Out)"
		row.AddCell().Value = "Goods_Return"
		row.AddCell().Value = "Final_Stock"
		row.AddCell().Value = "Actual_Stock"
		row.AddCell().Value = "Stock_Differential"

		for _, v := range data {
			row = sheet.AddRow()
			row.AddCell().Value = v.ProductCode                           // Product_Code
			row.AddCell().Value = v.ProductName                           // Product_Name
			row.AddCell().Value = v.Category                              // Category
			row.AddCell().Value = v.Uom                                   // UOM
			row.AddCell().SetFloatWithFormat(v.Stock, "0.00")             // Stock
			row.AddCell().SetFloatWithFormat(v.PlanInbound, "0.00")       // Plan_Inbound
			row.AddCell().SetFloatWithFormat(v.ActualInbound, "0.00")     // Actual_Inbound
			row.AddCell().SetFloatWithFormat(v.Waste, "0.00")             // Waste
			row.AddCell().SetFloatWithFormat(v.PlanDelivery, "0.00")      // Plan_Delivery
			row.AddCell().SetFloatWithFormat(v.ActualDelivery, "0.00")    // Actual_Delivery
			row.AddCell().SetFloatWithFormat(v.StockTransferIn, "0.00")   // Stock_Transfer(In)
			row.AddCell().SetFloatWithFormat(v.StockTransferOut, "0.00")  // Stock_Transfer(Out)
			row.AddCell().SetFloatWithFormat(v.GoodsReturn, "0.00")       // Goods_Return
			row.AddCell().SetFloatWithFormat(v.FinalStock, "0.00")        // Final_Stock
			row.AddCell().SetFloatWithFormat(v.ActualStock, "0.00")       // Actual_Stock
			row.AddCell().SetFloatWithFormat(v.StockDifferential, "0.00") // Stock_Differential
		}

	}
	err = file.Save(fileDir)
	filePath, err = util.UploadToS3(filename, fileDir, "application/xlsx")
	// fungsi ini berguna untuk menghapus kembali file yang sudah diupload ke server
	os.Remove(fileDir)
	return
}

// GetPickingXls : function to create excel file of picking report
func GetPickingXls(data []*reportPicking, warehouse *model.Warehouse) (filePath string, err error) {
	var file *xlsx.File
	var sheet *xlsx.Sheet
	var row *xlsx.Row

	dir := util.ExportDirectory
	filename := fmt.Sprintf("ReportPickingOrder_%s_%s_%s.xlsx", time.Now().Format("2006-01-02"), util.ReplaceSpace(warehouse.Name), util.GenerateRandomDoc(5))

	fileDir := fmt.Sprintf("%s/%s", dir, filename)

	file = xlsx.NewFile()
	if sheet, err = file.AddSheet("Sheet1"); err == nil {
		row = sheet.AddRow()
		row.SetHeight(20)
		row.AddCell().Value = "Delivery_Date"
		row.AddCell().Value = "Time_Assign"
		row.AddCell().Value = "Picking_List_Code"
		row.AddCell().Value = "Sales_Order_Code"
		row.AddCell().Value = "Order_Type"
		row.AddCell().Value = "Payment_Term"
		row.AddCell().Value = "Merchant"
		row.AddCell().Value = "Business_Type"
		row.AddCell().Value = "Total_Item"
		row.AddCell().Value = "Sales_Order_Weight"
		row.AddCell().Value = "Total_Weight"
		row.AddCell().Value = "Total_Koli"
		row.AddCell().Value = "Shipping_Adrres"
		row.AddCell().Value = "WRT"
		row.AddCell().Value = "Warehouse"
		row.AddCell().Value = "Picker"
		row.AddCell().Value = "Time_Start_Picked"
		row.AddCell().Value = "Time_Finished_Picked"
		row.AddCell().Value = "Checker"
		row.AddCell().Value = "Time_Checkin"
		row.AddCell().Value = "Time_Checkout"
		row.AddCell().Value = "Vendor"
		row.AddCell().Value = "Planning"
		row.AddCell().Value = "Courier"
		row.AddCell().Value = "Dispatch_Time"
		row.AddCell().Value = "Status_Picking_Assigned"
		row.AddCell().Value = "Status_Sales_Order"

		for _, v := range data {
			row = sheet.AddRow()
			row.AddCell().Value = v.DeliveryDate                         // DeliveryDate
			row.AddCell().Value = v.TimestampAssign                      // TimestampAssign
			row.AddCell().Value = v.PickingListCode                      // PickingListCode
			row.AddCell().Value = v.SalesOrderCode                       // SalesOrderCode
			row.AddCell().Value = v.SalesOrderType                       // SalesOrderType
			row.AddCell().Value = v.PaymentTerm                          // SalesOrderPaymentTerm
			row.AddCell().Value = v.Merchant                             // Merchant
			row.AddCell().Value = v.BusinessType                         // MerchantBusinessType
			row.AddCell().SetFloatWithFormat(v.Item, "0.00")             // Item
			row.AddCell().SetFloatWithFormat(v.SalesOrderWeight, "0.00") // SalesOrderWeight
			row.AddCell().SetFloatWithFormat(v.TotalWeight, "0.00")      // TotalWeight
			row.AddCell().SetFloatWithFormat(v.TotalKoli, "0.00")        // TotalKoli
			row.AddCell().Value = v.ShippingAddress                      // ShippingAddress
			row.AddCell().Value = v.Wrt                                  // Wrt
			row.AddCell().Value = v.Warehouse                            // Warehouse
			row.AddCell().Value = v.Picker                               // Picker
			row.AddCell().Value = v.TimeStartPicked                      // TimeStartPicked
			row.AddCell().Value = v.TimeFinishPicked                     // TimeFinishPicked
			row.AddCell().Value = v.Checker                              // Checker
			row.AddCell().Value = v.TimeCheckin                          // TimeCheckin
			row.AddCell().Value = v.TimeCheckout                         // TimeCheckout
			row.AddCell().Value = v.Vendor                               // Vendor
			row.AddCell().Value = v.Planning                             // Planning
			row.AddCell().Value = v.Courier                              // Courier
			row.AddCell().Value = v.DispatchTime                         // DispatchTime
			row.AddCell().Value = v.StatusPickingAssigned                // StatusPickingAssigned
			row.AddCell().Value = v.StatusSalesOrder                     // StatusSalesOrder

		}

	}
	err = file.Save(fileDir)
	filePath, err = util.UploadToS3(filename, fileDir, "application/xlsx")
	// fungsi ini berguna untuk menghapus kembali file yang sudah diupload ke server
	os.Remove(fileDir)
	return
}

// GetPickingOrderItemXls : function to create excel file of picking order item report
func GetPickingOrderItemXls(data []*reportPickingOrderItem, warehouse *model.Warehouse) (filePath string, err error) {
	var file *xlsx.File
	var sheet *xlsx.Sheet
	var row *xlsx.Row

	dir := util.ExportDirectory
	fileName := fmt.Sprintf("ReportPickingOrderItem_%s_%s_%s.xlsx", time.Now().Format("2006-01-02"), util.ReplaceSpace(warehouse.Name), util.GenerateRandomDoc(5))

	fileDir := fmt.Sprintf("%s/%s", dir, fileName)

	file = xlsx.NewFile()
	if sheet, err = file.AddSheet("Sheet1"); err != nil {
		return
	}

	row = sheet.AddRow()
	row.SetHeight(20)
	row.AddCell().Value = "Delivery_Date"
	row.AddCell().Value = "Picking_List_Code"
	row.AddCell().Value = "Sales_Order_Code"
	row.AddCell().Value = "Merchant"
	row.AddCell().Value = "Product_Code"
	row.AddCell().Value = "Product_Name"
	row.AddCell().Value = "Uom"
	row.AddCell().Value = "Order_Quantity"
	row.AddCell().Value = "Quantity_Picker"
	row.AddCell().Value = "Quantity_Checker"
	row.AddCell().Value = "WRT"
	row.AddCell().Value = "Warehouse"
	row.AddCell().Value = "Unfulfilled_Note"

	for _, v := range data {
		row = sheet.AddRow()
		row.AddCell().Value = v.DeliveryDate                   // DeliveryDate
		row.AddCell().Value = v.PickingListCode                // PickingListCode
		row.AddCell().Value = v.SalesOrderCode                 // SalesOrderCode
		row.AddCell().Value = v.Merchant                       // Merchant
		row.AddCell().Value = v.ProductCode                    // ProductCode
		row.AddCell().Value = v.ProductName                    // ProductName
		row.AddCell().Value = v.Uom                            // UOM
		row.AddCell().SetFloatWithFormat(v.OrderQty, "0.00")   // OrderQty
		row.AddCell().SetFloatWithFormat(v.QtyPicker, "0.00")  // QtyPicker
		row.AddCell().SetFloatWithFormat(v.QtyChecker, "0.00") // QtyChecker
		row.AddCell().Value = v.Wrt                            // Wrt
		row.AddCell().Value = v.Warehouse                      // Category
		row.AddCell().Value = v.UnfullfilledNote               // UnfullfilledNote

	}

	err = file.Save(fileDir)
	filePath, err = util.UploadToS3(fileName, fileDir, "application/xlsx")
	// fungsi ini berguna untuk menghapus kembali file yang sudah diupload ke server
	os.Remove(fileDir)
	return
}

// GetGoodsTransferItemXls : function to create excel file of goods transfer item report
func GetGoodsTransferItemXls(data []*reportGoodsTransferItem, warehouse *model.Warehouse) (filePath string, err error) {
	var file *xlsx.File
	var sheet *xlsx.Sheet
	var row *xlsx.Row

	dir := util.ExportDirectory
	fileName := fmt.Sprintf("ReportGoodsTransferItem_%s_%s_%s.xlsx", time.Now().Format("2006-01-02"), util.ReplaceSpace(warehouse.Name), util.GenerateRandomDoc(5))

	fileDir := fmt.Sprintf("%s/%s", dir, fileName)

	file = xlsx.NewFile()
	if sheet, err = file.AddSheet("Sheet1"); err != nil {
		return
	}

	row = sheet.AddRow()
	row.SetHeight(20)
	row.AddCell().Value = "Timestamp"
	row.AddCell().Value = "Goods_Transfer_Code"
	row.AddCell().Value = "Product_Code"
	row.AddCell().Value = "Product_Name"
	row.AddCell().Value = "Uom"
	row.AddCell().Value = "Warehouse_Origin"
	row.AddCell().Value = "Warehouse_Destination"
	row.AddCell().Value = "Request_Qty"
	row.AddCell().Value = "Transfer_Qty"
	row.AddCell().Value = "Received_Qty"
	row.AddCell().Value = "Status"
	row.AddCell().Value = "Doc_Note"
	row.AddCell().Value = "Note"

	for _, v := range data {
		row = sheet.AddRow()
		row.AddCell().Value = v.Timestamp
		row.AddCell().Value = v.GoodsTransferCode
		row.AddCell().Value = v.ProductCode
		row.AddCell().Value = v.ProductName
		row.AddCell().Value = v.Uom
		row.AddCell().Value = v.WarehouseOrigin
		row.AddCell().Value = v.WarehouseDestination
		row.AddCell().SetFloatWithFormat(v.RequestQty, "0.00")
		row.AddCell().SetFloatWithFormat(v.TransferQty, "0.00")
		row.AddCell().SetFloatWithFormat(v.ReceivedQty, "0.00")
		row.AddCell().Value = v.Status
		row.AddCell().Value = v.DocNote
		row.AddCell().Value = v.Note

	}

	err = file.Save(fileDir)
	filePath, err = util.UploadToS3(fileName, fileDir, "application/xlsx")
	// fungsi ini berguna untuk menghapus kembali file yang sudah diupload ke server
	os.Remove(fileDir)
	return
}

// GetPickingRoutingXls : function to create excel file of Picking Routing report
func GetPickingRoutingXls(data []*reportPickingRouting, warehouse *model.Warehouse) (filePath string, err error) {
	var file *xlsx.File
	var sheet *xlsx.Sheet
	var row *xlsx.Row

	dir := util.ExportDirectory
	fileName := fmt.Sprintf("ReportPickingRoutingItem_%s_%s_%s.xlsx", time.Now().Format("2006-01-02"), util.ReplaceSpace(warehouse.Name), util.GenerateRandomDoc(5))

	fileDir := fmt.Sprintf("%s/%s", dir, fileName)

	file = xlsx.NewFile()
	if sheet, err = file.AddSheet("Sheet1"); err != nil {
		return
	}

	row = sheet.AddRow()
	row.SetHeight(20)
	row.AddCell().Value = "SO Code"
	row.AddCell().Value = "PL code"
	row.AddCell().Value = "Product Name"
	row.AddCell().Value = "Product Quantity"
	row.AddCell().Value = "Product UOM"
	row.AddCell().Value = "Bin Name"
	row.AddCell().Value = "Picker Name"
	row.AddCell().Value = "Lead Picker Name"
	row.AddCell().Value = "Route Sequence"
	row.AddCell().Value = "Step Type"
	row.AddCell().Value = "Expected Walking Duration"
	row.AddCell().Value = "Actual Walking Duration"
	row.AddCell().Value = "Expected Picking Duration"
	row.AddCell().Value = "Actual Picking Duration"
	row.AddCell().Value = "Status"

	for _, v := range data {
		row = sheet.AddRow()
		row.AddCell().Value = v.SalesOrderCode
		row.AddCell().Value = v.PickingListCode
		row.AddCell().Value = v.Product
		row.AddCell().SetFloatWithFormat(v.OrderQty, "0.00")
		row.AddCell().Value = v.UOM
		row.AddCell().Value = v.RackName
		row.AddCell().Value = v.Picker
		row.AddCell().Value = v.LeadPicker
		row.AddCell().SetInt(int(v.Sequence))
		row.AddCell().Value = v.StepTypeStr
		row.AddCell().SetInt(int(v.ExpectedWalkingDuration))
		row.AddCell().SetInt(int(v.ActualWalkingDuration))
		row.AddCell().SetInt(int(v.ExpectedServiceDuration))
		row.AddCell().SetInt(int(v.ActualPickingDuration))
		row.AddCell().Value = v.StatusStr

	}

	err = file.Save(fileDir)
	filePath, err = util.UploadToS3(fileName, fileDir, "application/xlsx")
	// fungsi ini berguna untuk menghapus kembali file yang sudah diupload ke server
	os.Remove(fileDir)
	return
}

func GetTransferSkuItemXls(data []*reportTransferSkuItem, warehouse *model.Warehouse) (filePath string, err error) {
	var file *xlsx.File
	var sheet *xlsx.Sheet
	var row *xlsx.Row
	var filename string

	dir := util.ExportDirectory
	if warehouse.Name == "All Warehouse" {
		filename = fmt.Sprintf("ReportTransferSkuItem_%s_%s.xlsx", time.Now().Format("2006-01-02"), util.GenerateRandomDoc(5))
	} else {
		filename = fmt.Sprintf("ReportTransferSkuItem_%s_%s_%s.xlsx", time.Now().Format("2006-01-02"), util.ReplaceSpace(warehouse.Name), util.GenerateRandomDoc(5))
	}

	fileDir := fmt.Sprintf("%s/%s", dir, filename)

	file = xlsx.NewFile()
	if sheet, err = file.AddSheet("Sheet1"); err == nil {
		row = sheet.AddRow()
		row.SetHeight(20)
		row.AddCell().Value = "TS_Code"
		row.AddCell().Value = "TS_Status"
		row.AddCell().Value = "Recognition_Date"
		row.AddCell().Value = "Inbound_Code"
		row.AddCell().Value = "Supplier_Code"
		row.AddCell().Value = "Supplier_Name"
		row.AddCell().Value = "GR_Code"
		row.AddCell().Value = "GR_Status"
		row.AddCell().Value = "Warehouse_Origin"
		row.AddCell().Value = "Warehouse_Destination"
		row.AddCell().Value = "Product_Code"
		row.AddCell().Value = "Product_Name"
		row.AddCell().Value = "UOM"
		row.AddCell().Value = "Goods Stock"
		row.AddCell().Value = "Received_Qty"
		row.AddCell().Value = "After_Sortir_Qty(Good)"
		row.AddCell().Value = "After_Sortir_Qty(Waste)"
		row.AddCell().Value = "Discrepancy"
		row.AddCell().Value = "After_Sortir_Qty(Down_Grade)"
		row.AddCell().Value = "Product_Code_Downgrade"
		row.AddCell().Value = "Product_DownGrade"
		row.AddCell().Value = "UOM"

		for _, v := range data {
			row = sheet.AddRow()
			row.AddCell().Value = v.TransferSkuCode
			row.AddCell().Value = v.TransferSkuStatus
			row.AddCell().Value = v.RecognitionDate
			row.AddCell().Value = v.InboundCode
			row.AddCell().Value = v.SupplierCode
			row.AddCell().Value = v.SupplierName
			row.AddCell().Value = v.GRCode
			row.AddCell().Value = v.GRStatus
			row.AddCell().Value = v.WarehouseOrigin
			row.AddCell().Value = v.WarehouseDestination
			row.AddCell().Value = v.ProductCode
			row.AddCell().Value = v.ProductName
			row.AddCell().Value = v.Uom
			row.AddCell().SetFloatWithFormat(v.GoodsStock, "0.00")
			row.AddCell().SetFloatWithFormat(v.ReceivedQty, "0.00")
			row.AddCell().SetFloatWithFormat(v.AfterSortirQty, "0.00")
			row.AddCell().SetFloatWithFormat(v.AfterSortirQtyForWaste, "0.00")
			row.AddCell().SetFloatWithFormat(v.Discrepancy, "0.00")
			row.AddCell().SetFloatWithFormat(v.AfterSortirQtyForDownGrade, "0.00")
			row.AddCell().Value = v.ProductCodeDowngrade
			row.AddCell().Value = v.ProductNameDowngrade
			row.AddCell().Value = v.UOMDowngrade
		}

	}
	err = file.Save(fileDir)
	filePath, err = util.UploadToS3(filename, fileDir, "application/xlsx")

	// fungsi ini berguna untuk menghapus kembali file yang sudah diupload ke server
	os.Remove(fileDir)
	return
}

func GetPackingRecommendationXls(data []*model.BarcodeModel, warehouse *model.Warehouse) (filePath string, err error) {
	var file *xlsx.File
	var sheet *xlsx.Sheet
	var row *xlsx.Row
	var filename string
	m := mongodb.NewMongo()

	dir := util.ExportDirectory
	if warehouse.Name == "All Warehouse" {
		filename = fmt.Sprintf("ReportPacking_%s_%s.xlsx", time.Now().Format("2006-01-02"), util.GenerateRandomDoc(5))
	} else {
		filename = fmt.Sprintf("ReportPacking_%s_%s_%s.xlsx", time.Now().Format("2006-01-02"), util.ReplaceSpace(warehouse.Name), util.GenerateRandomDoc(5))
	}

	fileDir := fmt.Sprintf("%s/%s", dir, filename)

	file = xlsx.NewFile()
	if sheet, err = file.AddSheet("Sheet1"); err == nil {
		row = sheet.AddRow()
		row.SetHeight(20)
		row.AddCell().Value = "Packing_Order_Code"
		row.AddCell().Value = "Delivery_Date"
		row.AddCell().Value = "Area"
		row.AddCell().Value = "Warehouse"
		row.AddCell().Value = "Product_Code"
		row.AddCell().Value = "Product_Name"
		row.AddCell().Value = "Pack_Expected_Weight"
		row.AddCell().Value = "Pack_Actual_Weight"
		row.AddCell().Value = "Packing_Order_Status"
		row.AddCell().Value = "Pack_Code"
		row.AddCell().Value = "Print_Count"
		row.AddCell().Value = "Created_At"
		row.AddCell().Value = "Created_By_Code"
		row.AddCell().Value = "Created_By_Name"
		row.AddCell().Value = "Deleted_By_Code"
		row.AddCell().Value = "Deleted_By_Name"
		row.AddCell().Value = "Pack_Status"

		for _, v := range data {

			// region read data
			v.PackingOrder, _ = repository.ValidPackingOrder(v.PackingOrderID)
			v.Product, _ = repository.ValidProduct(v.ProductID)
			v.CreatedObj, _ = repository.ValidStaff(v.CreatedBy)

			var deletedCode, deletedName string
			if v.DeletedBy != 0 {
				v.DeletedObj, _ = repository.ValidStaff(v.DeletedBy)
				deletedCode = v.DeletedObj.Code
				deletedName = v.DeletedObj.Name
			}
			// endregion

			var statusStr string
			if v.Status == 1 {
				statusStr = "Active"
			} else {
				statusStr = "Deleted"
			}

			// region read Packing Item
			filter := bson.D{
				{"packing_order_id", v.PackingOrder.ID},
				{"product_id", v.Product.ID},
				{"pack_type", v.PackType},
			}

			// endregion
			var rp *model.ResponseData
			var res2 []byte
			if res2, err = m.GetOneDataWithFilter("Packing_Item", filter); err != nil {
				return "", err
			}

			// region convert byte data to json data
			if err = json.Unmarshal(res2, &rp); err != nil {
				return "", err
			}

			v.PackingOrder.Warehouse.Read("ID")
			v.PackingOrder.Area.Read("ID")
			row = sheet.AddRow()
			row.AddCell().Value = v.PackingOrder.Code
			row.AddCell().Value = v.PackingOrder.DeliveryDate.Format("2006-01-02")
			row.AddCell().Value = v.PackingOrder.Area.Name
			row.AddCell().Value = v.PackingOrder.Warehouse.Name
			row.AddCell().Value = v.Product.Code
			row.AddCell().Value = v.Product.Name
			row.AddCell().SetFloatWithFormat(rp.PackType, "0.00")
			row.AddCell().SetFloatWithFormat(v.WeightScale, "0.00")
			row.AddCell().Value = util.ConvertStatusDoc(v.PackingOrder.Status)
			row.AddCell().Value = v.Code
			row.AddCell().SetInt(v.DeltaPrint)
			row.AddCell().Value = v.CreatedAt
			row.AddCell().Value = v.CreatedObj.Code
			row.AddCell().Value = v.CreatedObj.Name
			row.AddCell().Value = deletedCode
			row.AddCell().Value = deletedName
			row.AddCell().Value = statusStr
		}

	}
	err = file.Save(fileDir)
	filePath, err = util.UploadToS3(filename, fileDir, "application/xlsx")

	// fungsi ini berguna untuk menghapus kembali file yang sudah diupload ke server
	m.DisconnectMongoClient()
	os.Remove(fileDir)
	return
}
