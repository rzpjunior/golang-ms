package entry

import (
	"fmt"
	"os"
	"time"

	"git.edenfarm.id/cuxs/common"
	"git.edenfarm.id/cuxs/env"
	"git.edenfarm.id/cuxs/orm"
	"git.edenfarm.id/project-version2/api/datastore/repository"
	"git.edenfarm.id/project-version2/api/log"
	"git.edenfarm.id/project-version2/api/util"
	"git.edenfarm.id/project-version2/datamodel/model"
	"github.com/tealeg/xlsx"
)

// Save : function to save data requested into database
func Save(r createRequest) (id string, e error) {
	//generate codes for document
	r.CodeWasteEntry, _ = util.GenerateDocCode("WE", r.Warehouse.Code, "waste_entry")
	o := orm.NewOrm()
	o.Begin()
	u := &model.WasteEntry{
		Warehouse:       r.Warehouse,
		Code:            r.CodeWasteEntry,
		Status:          1,
		RecognitionDate: r.RecognitionDateAt,
		Note:            r.Note,
	}
	if _, e = o.Insert(u); e == nil {
		e = log.AuditLogByUser(r.Session.Staff, u.ID, "waste_entry", "create", "")
		for _, row := range r.WasteEntryItems {
			//if no changes did on the product, proceed checking the next product
			if row.WasteStock == 0 {
				continue
			}
			// get value_int from glossary
			wasteReasonValue, e := repository.GetGlossaryMultipleValue("table", "all", "attribute", "waste_reason", "value_name", row.WasteReason)
			if e != nil {
				o.Rollback()
			}

			item := &model.WasteEntryItem{
				WasteEntry:  &model.WasteEntry{ID: u.ID},
				Product:     row.Product,
				WasteQty:    row.WasteStock,
				Note:        row.Note,
				WasteReason: wasteReasonValue.ValueInt,
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

func Confirm(r confirmRequest) (we *model.WasteEntry, e error) {
	o := orm.NewOrm()
	o.Begin()
	orSelect := orm.NewOrm()
	orSelect.Using("read_only")

	we = &model.WasteEntry{
		ID:     r.ID,
		Status: 2,
	}
	if _, e = o.Update(we, "Status"); e == nil {
		for _, row := range r.WasteEntry.WasteEntryItems {

			var stocks []*model.Stock
			var sls []*model.StockLog
			var sl *model.StockLog
			var wls []*model.WasteLog
			var wl *model.WasteLog

			orSelect.Raw("SELECT * FROM stock where warehouse_id = ? AND product_id = ?", r.WasteEntry.Warehouse.ID, row.Product.ID).QueryRows(&stocks)
			for _, stock := range stocks {
				sl = &model.StockLog{
					Warehouse:    r.WasteEntry.Warehouse,
					Ref:          r.WasteEntry.ID,
					Product:      row.Product,
					Quantity:     row.WasteQty,
					RefType:      6,
					Type:         2,
					InitialStock: stock.AvailableStock,
					FinalStock:   stock.AvailableStock - row.WasteQty,
					UnitCost:     0,
					DocNote:      r.WasteEntry.Note,
					Status:       1,
					ItemNote:     row.Note,
					CreatedAt:    time.Now(),
				}

				sls = append(sls, sl)

				wl = &model.WasteLog{
					Warehouse:    r.WasteEntry.Warehouse,
					Ref:          r.WasteEntry.ID,
					Product:      row.Product,
					RefType:      2,
					Type:         1,
					InitialStock: stock.WasteStock,
					Quantity:     row.WasteQty,
					FinalStock:   stock.WasteStock + row.WasteQty,
					Status:       1,
					DocNote:      r.WasteEntry.Note,
					ItemNote:     row.Note,
					WasteReason:  row.WasteReason,
				}

				wls = append(wls, wl)

				stock.WasteStock = stock.WasteStock + row.WasteQty
				stock.AvailableStock = stock.AvailableStock - row.WasteQty

				if _, e = o.Update(stock, "WasteStock", "AvailableStock"); e != nil {
					o.Rollback()
					return nil, e
				}
			}

			if _, e = o.InsertMulti(100, wls); e != nil {
				o.Rollback()
				return nil, e
			}

			if _, e = o.InsertMulti(100, sls); e != nil {
				o.Rollback()
				return nil, e
			}
		}

		e = log.AuditLogByUser(r.Session.Staff, r.ID, "waste_entry", "confirm", "")

	} else {
		o.Rollback()
		return nil, e
	}

	o.Commit()

	return
}

func ExportFormXls(date time.Time, r []*model.Stock, warehouse *model.Warehouse, wr []*model.Glossary) (filePath string, err error) {
	var file *xlsx.File
	var sheet *xlsx.Sheet
	var row *xlsx.Row
	dir := env.GetString("EXPORT_DIRECTORY", "")

	filename := fmt.Sprintf("TemplateWasteEntryItem%s_%s_%s.xlsx", util.ReplaceSpace(warehouse.Name), date.Format("02012006"), util.GenerateRandomDoc(5))

	fileDir := fmt.Sprintf("%s/%s", dir, filename)

	file = xlsx.NewFile()
	if sheet, err = file.AddSheet("Sheet1"); err == nil {
		row = sheet.AddRow()
		row.SetHeight(20)
		row.AddCell().Value = "No"
		row.AddCell().Value = "Product_ID"
		row.AddCell().Value = "Product_Code"
		row.AddCell().Value = "Product_Name"
		row.AddCell().Value = "UOM"
		row.AddCell().Value = "Stock"
		row.AddCell().Value = "Waste_Stock"
		row.AddCell().Value = "Quantity"
		row.AddCell().Value = "Waste_Reason"
		row.AddCell().Value = "Note"

		for i, v := range r {
			v.Product.Read("ID")
			v.Product.Uom.Read("ID")

			row = sheet.AddRow()
			row.AddCell().SetInt(i + 1)
			row.AddCell().Value = common.Encrypt(v.Product.ID)         // ProductId
			row.AddCell().Value = v.Product.Code                       // Product Code
			row.AddCell().Value = v.Product.Name                       // Product Name
			row.AddCell().Value = v.Product.Uom.Name                   // UOM
			row.AddCell().SetFloatWithFormat(v.AvailableStock, "0.00") // Stock
			row.AddCell().SetFloatWithFormat(v.WasteStock, "0.00")     // Waste Stock
		}
		err = file.Save(fileDir)
	}

	if sheet, err = file.AddSheet("Waste Reason Choices"); err == nil {
		row = sheet.AddRow()
		row.SetHeight(20)
		row.AddCell().Value = "No"
		row.AddCell().Value = "Waste Reason Choices"

		for i, v := range wr {
			row = sheet.AddRow()
			row.AddCell().SetInt(i + 1)
			row.AddCell().Value = v.ValueName //Waste Reason Choice
		}
		err = file.Save(fileDir)
	}

	err = file.Save(fileDir)
	filePath, err = util.UploadToS3(filename, fileDir, "application/xlsx")

	// fungsi ini berguna untuk menghapus kembali file yang sudah diupload ke server
	os.Remove(fileDir)
	return
}

func Cancel(r cancelRequest) (we *model.WasteEntry, e error) {
	o := orm.NewOrm()
	o.Begin()

	we = &model.WasteEntry{
		ID: r.ID,
	}

	we.Status = 3

	if _, e = o.Update(we, "Status"); e == nil {
		e = log.AuditLogByUser(r.Session.Staff, we.ID, "waste_entry", "cancel", r.Note)
	} else {
		o.Rollback()
		return nil, e
	}

	o.Commit()

	return r.WasteEntry, e
}
