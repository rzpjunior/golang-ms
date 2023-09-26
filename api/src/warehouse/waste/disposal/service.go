package disposal

import (
	"fmt"
	"os"
	"time"

	"git.edenfarm.id/cuxs/common"
	"git.edenfarm.id/cuxs/env"
	"git.edenfarm.id/cuxs/orm"
	"git.edenfarm.id/project-version2/api/log"
	"git.edenfarm.id/project-version2/api/util"
	"git.edenfarm.id/project-version2/datamodel/model"
	"github.com/tealeg/xlsx"
)

// Save : function to save data requested into database
func Save(r createRequest) (wd *model.WasteDisposal, e error) {
	//generate codes for document
	r.Code, _ = util.GenerateDocCode("WD", r.Warehouse.Code, "waste_disposal")
	o := orm.NewOrm()
	o.Begin()
	wd = &model.WasteDisposal{
		Warehouse:       r.Warehouse,
		Code:            r.Code,
		Status:          int8(1),
		RecognitionDate: r.RecognitionDateAt,
		Note:            r.Note,
	}
	if _, e = o.Insert(wd); e == nil {
		var arrWdi []*model.WasteDisposalItem
		for _, row := range r.WasteDisposalItems {
			//if no changes did on the product, proceed checking the next product
			if row.Quantity == 0 {
				continue
			}

			item := &model.WasteDisposalItem{
				Product:       row.Product,
				WasteDisposal: wd,
				DisposeQty:    row.Quantity,
				Note:          row.Note,
			}
			arrWdi = append(arrWdi, item)
		}
		if _, e = o.InsertMulti(100, &arrWdi); e == nil {
			e = log.AuditLogByUser(r.Session.Staff, wd.ID, "waste_disposal", "create", "")
		} else {
			o.Rollback()
		}
	} else {
		o.Rollback()
	}
	o.Commit()
	return wd, e
}

func ExportFormXls(date time.Time, r []*model.Stock, warehouse *model.Warehouse) (filePath string, err error) {
	var file *xlsx.File
	var sheet *xlsx.Sheet
	var row *xlsx.Row
	dir := env.GetString("EXPORT_DIRECTORY", "")

	filename := fmt.Sprintf("TemplateWasteDisposal_%s_%s_%s.xlsx", util.ReplaceSpace(warehouse.Name), date.Format("02012006"), util.GenerateRandomDoc(5))

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
		row.AddCell().Value = "Waste_Stock"
		row.AddCell().Value = "Quantity"
		row.AddCell().Value = "Note"

		for i, v := range r {
			v.Product.Read("ID")
			v.Product.Uom.Read("ID")

			row = sheet.AddRow()
			row.AddCell().SetInt(i + 1)
			row.AddCell().Value = common.Encrypt(v.Product.ID)     // ProductId
			row.AddCell().Value = v.Product.Code                   // Product Code
			row.AddCell().Value = v.Product.Name                   // Product Name
			row.AddCell().Value = v.Product.Uom.Name               // UOM
			row.AddCell().SetFloatWithFormat(v.WasteStock, "0.00") // Waste Stok
		}

		err = file.Save(fileDir)
	}
	err = file.Save(fileDir)
	filePath, err = util.UploadToS3(filename, fileDir, "application/xlsx")

	// fungsi ini berguna untuk menghapus kembali file yang sudah diupload ke server
	os.Remove(fileDir)
	return
}

func Cancel(r cancelRequest) (wd *model.WasteDisposal, e error) {
	o := orm.NewOrm()
	o.Begin()

	wd = &model.WasteDisposal{
		ID: r.ID,
	}

	wd.Status = 3

	if _, e = o.Update(wd, "Status"); e == nil {
		e = log.AuditLogByUser(r.Session.Staff, wd.ID, "waste_disposal", "cancel", r.CancellationNote)
	} else {
		o.Rollback()
		return nil, e
	}

	o.Commit()

	return r.WasteDisposal, e
}

func Confirm(r confirmRequest) (wd *model.WasteDisposal, e error) {
	o := orm.NewOrm()
	o.Begin()

	for i, v := range r.Stocks {
		prevWasteStock := v.WasteStock
		disposeQty := r.WasteDisposal.WasteDisposalItems[i].DisposeQty

		v.WasteStock = prevWasteStock - disposeQty

		if _, e = o.Update(v, "WasteStock"); e == nil {
			wl := &model.WasteLog{
				Warehouse:    r.WasteDisposal.Warehouse,
				Product:      v.Product,
				Ref:          r.WasteDisposal.ID,
				RefType:      3,
				Type:         2,
				InitialStock: prevWasteStock,
				Quantity:     disposeQty,
				FinalStock:   v.WasteStock,
				DocNote:      r.WasteDisposal.Note,
				ItemNote:     r.WasteDisposal.WasteDisposalItems[i].Note,
				Status:       1,
			}

			if _, e = o.Insert(wl); e != nil {
				o.Rollback()
				return nil, e
			}
		} else {
			o.Rollback()
			return nil, e
		}
	}

	r.WasteDisposal.Status = 2
	if _, e = o.Update(r.WasteDisposal, "Status"); e != nil {
		o.Rollback()
		return nil, e
	}

	if e = log.AuditLogByUser(r.Session.Staff, r.WasteDisposal.ID, "waste_disposal", "confirm", ""); e != nil {
		o.Rollback()
		return nil, e
	}

	o.Commit()
	return r.WasteDisposal, nil
}
