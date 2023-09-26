package stock_opname

import (
	"fmt"
	"os"
	"time"

	"git.edenfarm.id/cuxs/common"
	"git.edenfarm.id/cuxs/cuxs/event"
	"git.edenfarm.id/cuxs/env"
	"git.edenfarm.id/cuxs/orm"
	"git.edenfarm.id/project-version2/api/log"
	"git.edenfarm.id/project-version2/api/util"
	"git.edenfarm.id/project-version2/datamodel/model"
	"github.com/tealeg/xlsx"
)

// Save : function to save data requested into database
func Save(r createRequest) (id string, e error) {
	//generate codes for document
	r.CodeStockOpname, _ = util.GenerateDocCode("ST", r.Warehouse.Code, "stock_opname")
	o := orm.NewOrm()
	o.Begin()
	u := &model.StockOpname{
		Warehouse:       r.Warehouse,
		Category:        r.Category,
		Code:            r.CodeStockOpname,
		StockType:       r.StockType.ValueInt,
		Status:          int8(1),
		RecognitionDate: r.RecognitionDateAt,
		Note:            r.Note,
	}
	if _, e = o.Insert(u); e == nil {
		e = log.AuditLogByUser(r.Session.Staff, u.ID, "stock_opname", "create", "")
		for _, row := range r.StockOpnameItems {
			item := &model.StockOpnameItem{
				StockOpname:  &model.StockOpname{ID: u.ID},
				Product:      row.Product,
				InitialStock: common.Rounder(row.InitialStock, 0.5, 2),
				AdjustQty:    common.Rounder(row.FinalStock-row.InitialStock, 0.5, 2),
				FinalStock:   common.Rounder(row.FinalStock, 0.5, 2),
				OpnameReason: row.OpnameReasonInt,
				Note:         row.Note,
			}

			if _, e = o.Insert(item); e != nil {
				o.Rollback()
			}
		}
	} else {
		o.Rollback()
	}
	if e == nil {
		o.Commit()
	}
	return common.Encrypt(u.ID), e
}

func ExportGoodStockFormXls(date time.Time, r []*model.Stock, warehouse *model.Warehouse, wr []*model.Glossary) (filePath string, err error) {
	var file *xlsx.File
	var sheet *xlsx.Sheet
	var row *xlsx.Row
	dir := env.GetString("EXPORT_DIRECTORY", "")

	filename := fmt.Sprintf("TemplateStockOpnameProducts_%s_%s_%s.xlsx", util.ReplaceSpace(warehouse.Name), date.Format("02012006"), util.GenerateRandomDoc(5))

	fileDir := fmt.Sprintf("%s/%s", dir, filename)

	file = xlsx.NewFile()
	if sheet, err = file.AddSheet("Sheet1"); err == nil {
		row = sheet.AddRow()
		row.SetHeight(20)
		row.AddCell().Value = "No"
		row.AddCell().Value = "Product_ID"
		row.AddCell().Value = "Product_Code"
		row.AddCell().Value = "Product_Name"
		row.AddCell().Value = "Category"
		row.AddCell().Value = "UOM"
		row.AddCell().Value = "Stock"
		row.AddCell().Value = "GoodStock"
		row.AddCell().Value = "Opname_Reason"
		row.AddCell().Value = "Note"

		for i, v := range r {
			v.Product.Read("ID")
			v.Product.Category.Read("ID")
			v.Product.Uom.Read("ID")

			row = sheet.AddRow()
			row.AddCell().SetInt(i + 1)
			row.AddCell().Value = common.Encrypt(v.Product.ID)         // ProductId
			row.AddCell().Value = v.Product.Code                       // Product Code
			row.AddCell().Value = v.Product.Name                       // Product Name
			row.AddCell().Value = v.Product.Category.Name              // Category
			row.AddCell().Value = v.Product.Uom.Name                   // UOM
			row.AddCell().SetFloatWithFormat(v.AvailableStock, "0.00") // AvailableStock
		}

		err = file.Save(fileDir)
	}

	if sheet, err = file.AddSheet("Opname Reason Choices"); err == nil {
		row = sheet.AddRow()
		row.SetHeight(20)
		row.AddCell().Value = "No"
		row.AddCell().Value = "Opname Reason Choices"

		for i, v := range wr {
			row = sheet.AddRow()
			row.AddCell().SetInt(i + 1)
			row.AddCell().Value = v.ValueName //Opname Reason Choice
		}
		err = file.Save(fileDir)
	}

	err = file.Save(fileDir)
	filePath, err = util.UploadToS3(filename, fileDir, "application/xlsx")
	// fungsi ini berguna untuk menghapus kembali file yang sudah diupload ke server
	os.Remove(fileDir)
	return
}

func ExportWasteStockFormXls(date time.Time, r []*model.Stock, warehouse *model.Warehouse, wr []*model.Glossary) (filePath string, err error) {
	var file *xlsx.File
	var sheet *xlsx.Sheet
	var row *xlsx.Row
	dir := env.GetString("EXPORT_DIRECTORY", "")

	filename := fmt.Sprintf("TemplateStockOpnameProducts_%s_%s_%s.xlsx", util.ReplaceSpace(warehouse.Name), date.Format("02012006"), util.GenerateRandomDoc(5))

	fileDir := fmt.Sprintf("%s/%s", dir, filename)

	file = xlsx.NewFile()
	if sheet, err = file.AddSheet("Sheet1"); err == nil {
		row = sheet.AddRow()
		row.SetHeight(20)
		row.AddCell().Value = "No"
		row.AddCell().Value = "Product_ID"
		row.AddCell().Value = "Product_Code"
		row.AddCell().Value = "Product_Name"
		row.AddCell().Value = "Category"
		row.AddCell().Value = "UOM"
		row.AddCell().Value = "Waste_Stock"
		row.AddCell().Value = "Actual_Waste_Stock"
		row.AddCell().Value = "Opname_Reason"
		row.AddCell().Value = "Note"

		for i, v := range r {
			v.Product.Read("ID")
			v.Product.Category.Read("ID")
			v.Product.Uom.Read("ID")

			row = sheet.AddRow()
			row.AddCell().SetInt(i + 1)
			row.AddCell().Value = common.Encrypt(v.Product.ID)     // ProductId
			row.AddCell().Value = v.Product.Code                   // Product Code
			row.AddCell().Value = v.Product.Name                   // Product Name
			row.AddCell().Value = v.Product.Category.Name          // Category
			row.AddCell().Value = v.Product.Uom.Name               // UOM
			row.AddCell().SetFloatWithFormat(v.WasteStock, "0.00") // WasteStock
		}

		err = file.Save(fileDir)
	}

	if sheet, err = file.AddSheet("Opname Reason Choices"); err == nil {
		row = sheet.AddRow()
		row.SetHeight(20)
		row.AddCell().Value = "No"
		row.AddCell().Value = "Opname Reason Choices"

		for i, v := range wr {
			row = sheet.AddRow()
			row.AddCell().SetInt(i + 1)
			row.AddCell().Value = v.ValueName //Opname Reason Choice
		}
		err = file.Save(fileDir)
	}

	err = file.Save(fileDir)
	filePath, err = util.UploadToS3(filename, fileDir, "application/xlsx")
	// fungsi ini berguna untuk menghapus kembali file yang sudah diupload ke server
	os.Remove(fileDir)
	return
}

func Cancel(r cancelRequest) (st *model.StockOpname, e error) {
	o := orm.NewOrm()
	o.Begin()

	st = &model.StockOpname{
		ID: r.ID,
	}

	st.Status = 3

	if _, e = o.Update(st, "Status"); e == nil {
		e = log.AuditLogByUser(r.Session.Staff, st.ID, "stock_opname", "cancel", r.Note)
	} else {
		o.Rollback()
		return nil, e
	}

	o.Commit()

	return r.StockOpname, e
}

func Confirm(r confirmRequest) (st *model.StockOpname, e error) {
	o := orm.NewOrm()
	o.Begin()
	orSelect := orm.NewOrm()
	orSelect.Using("read_only")

	st = &model.StockOpname{
		ID:     r.ID,
		Status: 2,
	}
	if _, e = o.Update(st, "Status"); e == nil {

		orSelect.LoadRelated(st, "StockOpnameItems", 2)

		for _, row := range st.StockOpnameItems {

			if row.AdjustQty != 0 {
				go event.Call("stockopname::commited", row)
			}

		}

		e = log.AuditLogByUser(r.Session.Staff, st.ID, "stock_opname", "confirm", "")

	} else {
		o.Rollback()
		return nil, e
	}

	o.Commit()

	return
}

func Download(do *model.StockOpname) (filePath string, e error) {

	var file *xlsx.File
	var sheet *xlsx.Sheet
	var row *xlsx.Row
	var err error

	dir := env.GetString("EXPORT_DIRECTORY", "")

	t := time.Now()
	date := t.Format("02012006")
	filename := fmt.Sprintf("ExportStockOpname_%s_%s_%s.xlsx", util.ReplaceSpace(do.Warehouse.Name), date, util.GenerateRandomDoc(5))

	fileDir := fmt.Sprintf("%s/%s", dir, filename)

	file = xlsx.NewFile()
	if sheet, err = file.AddSheet("Sheet1"); err == nil {
		row = sheet.AddRow()
		row.SetHeight(20)
		row.AddCell().Value = "No"
		row.AddCell().Value = "Product_ID"
		row.AddCell().Value = "Product_Code"
		row.AddCell().Value = "Product_Name"
		row.AddCell().Value = "Category"
		row.AddCell().Value = "UOM"
		row.AddCell().Value = "Initial Stock"
		row.AddCell().Value = "Final Stock"
		row.AddCell().Value = "Adjustment Qty"
		row.AddCell().Value = "Note"

		for i, v := range do.StockOpnameItems {
			v.Product.Read("ID")
			v.Product.Category.Read("ID")
			v.Product.Uom.Read("ID")

			row = sheet.AddRow()
			row.AddCell().SetInt(i + 1)
			row.AddCell().Value = common.Encrypt(v.Product.ID)       // ProductId
			row.AddCell().Value = v.Product.Code                     // Product Code
			row.AddCell().Value = v.Product.Name                     // Product Name
			row.AddCell().Value = v.Product.Category.Name            // Category
			row.AddCell().Value = v.Product.Uom.Name                 // UOM
			row.AddCell().SetFloatWithFormat(v.InitialStock, "0.00") // Initial Stock
			row.AddCell().SetFloatWithFormat(v.FinalStock, "0.00")   // Final Stock
			row.AddCell().SetFloatWithFormat(v.AdjustQty, "0.00")    // Adjustment Qty
			row.AddCell().Value = v.Note                             // Note
		}

		err = file.Save(fileDir)
	}
	err = file.Save(fileDir)
	filePath, err = util.UploadToS3(filename, fileDir, "application/xlsx")
	//fungsi ini berguna untuk menghapus kembali file yang sudah diupload ke server
	os.Remove(fileDir)
	return

}
