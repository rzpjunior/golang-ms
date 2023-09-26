// Copyright 2020 PT. Eden Pangan Indonesia. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

package product_price

import (
	"fmt"
	"math"
	"os"
	"strconv"
	"strings"
	"time"

	"git.edenfarm.id/cuxs/common"

	"git.edenfarm.id/cuxs/env"
	"git.edenfarm.id/cuxs/orm"
	"git.edenfarm.id/project-version2/datamodel/model"
	"git.edenfarm.id/project-version2/api/log"
	"git.edenfarm.id/project-version2/api/util"
	"github.com/tealeg/xlsx"
)

// getProductPrice: get Product price based on given request query
func getProductPrice(rq *orm.RequestQuery, tagProduct string) (m []*model.Price, total int64, err error) {
	q, _ := rq.QueryReadOnly(new(model.Price))

	cond := q.GetCond()

	if tagProduct != "" {
		cond1 := orm.NewCondition()
		cond1 = cond1.And("product__tag_product__icontains", ","+tagProduct+",").Or("product__tag_product__istartswith", tagProduct+",").Or("product__tag_product__iendswith", ","+tagProduct).Or("product__tag_product", tagProduct)

		cond = cond.AndCond(cond1)
	}

	q = q.SetCond(cond)

	// get total data
	if total, err = q.Count(); err != nil || total == 0 {
		return nil, total, err
	}

	// get data requested
	var mx []*model.Price
	if _, err = q.RelatedSel().All(&mx, rq.Fields...); err == nil {
		return mx, total, nil
	}

	// return error some thing went wrong
	return nil, total, err

}

// getProductPriceUpdateXls: download template update product price
func getProductPriceUpdateXls(date time.Time, r []*model.Price) (filePath string, err error) {
	var file *xlsx.File
	var sheet *xlsx.Sheet
	var row *xlsx.Row

	dir := env.GetString("EXPORT_DIRECTORY", "")

	var psName string
	if len(r) != 0 {
		psName = r[0].PriceSet.Name
		psName = strings.ReplaceAll(psName, " ", "_")
	}
	filename := fmt.Sprintf("Template_%s_%s_%s.xlsx", psName, date.Format("2006-01-02"), util.GenerateRandomDoc(5))
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
		row.AddCell().Value = "Category"
		row.AddCell().Value = "Salable"
		row.AddCell().Value = "Unit_Price"

		for index, i := range r {
			row = sheet.AddRow()
			row.AddCell().SetInt(index + 1) // No
			pID, _ := strconv.Atoi(common.Encrypt(int(i.Product.ID)))
			row.AddCell().SetInt(pID)                             // Product ID
			row.AddCell().Value = i.Product.Code                  // Product Code
			row.AddCell().Value = i.Product.Name                  // Product Name
			row.AddCell().Value = i.Product.Uom.Name              // UOM
			row.AddCell().Value = i.Product.Category.Name         // Category
			row.AddCell().SetInt(int(i.Product.Salability))       // Salable
			row.AddCell().SetFloatWithFormat(i.UnitPrice, "0.00") // Unit Price

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

	}
	return
}

// getShadowPriceUpdateXls: download template shadow product price
func getShadowPriceUpdateXls(date time.Time, r []*model.Price) (filePath string, err error) {
	var file *xlsx.File
	var sheet *xlsx.Sheet
	var row *xlsx.Row

	dir := env.GetString("EXPORT_DIRECTORY", "")

	var psName string
	if len(r) != 0 {
		psName = r[0].PriceSet.Name
		psName = strings.ReplaceAll(psName, " ", "_")
	}

	filename := fmt.Sprintf("Template_%s_%s_%s.xlsx", psName, date.Format("2006-01-02"), util.GenerateRandomDoc(5))
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
		row.AddCell().Value = "Category"
		row.AddCell().Value = "Salable"
		row.AddCell().Value = "Unit_Price"
		row.AddCell().Value = "Shadow_Price"

		for index, i := range r {
			row = sheet.AddRow()
			row.AddCell().SetInt(index + 1)                         // No
			row.AddCell().Value = common.Encrypt(int(i.Product.ID)) // Product ID
			row.AddCell().Value = i.Product.Code                    // Product Code
			row.AddCell().Value = i.Product.Name                    // Product Name
			row.AddCell().Value = i.Product.Uom.Name                // UOM
			row.AddCell().Value = i.Product.Category.Name           // Category
			row.AddCell().SetInt(int(i.Product.Salability))         // Salable
			row.AddCell().SetFloatWithFormat(i.UnitPrice, "0.00")   // Unit Price
			row.AddCell().SetFloatWithFormat(i.ShadowPrice, "0.00") // Shadow Price

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
	}
	return
}

// exportTemplateXls: download template export product price
func exportTemplateXls(date time.Time, r []*model.Price) (filePath string, err error) {
	var file *xlsx.File
	var sheet *xlsx.Sheet
	var row *xlsx.Row

	dir := env.GetString("EXPORT_DIRECTORY", "")

	var psName string
	if len(r) != 0 {
		psName = r[0].PriceSet.Name
		psName = strings.ReplaceAll(psName, " ", "_")
	}
	filename := fmt.Sprintf("ProductPrice%s_%s_%s.xlsx", date.Format("2006-01-02"), psName, util.GenerateRandomDoc(5))
	fileDir := fmt.Sprintf("%s/%s", dir, filename)

	file = xlsx.NewFile()
	if sheet, err = file.AddSheet("Sheet1"); err == nil {
		row = sheet.AddRow()
		row.SetHeight(20)
		row.AddCell().Value = "No"
		row.AddCell().Value = "Product_Code"
		row.AddCell().Value = "Product_Name"
		row.AddCell().Value = "UOM"
		row.AddCell().Value = "Category"
		row.AddCell().Value = "Price_Set"
		row.AddCell().Value = "Salable"
		row.AddCell().Value = "Unit_Price"
		row.AddCell().Value = "Shadow_Price"
		row.AddCell().Value = "Shadow_Percentage"

		for index, i := range r {
			row = sheet.AddRow()
			row.AddCell().SetInt(index + 1)                            // No
			row.AddCell().Value = i.Product.Code                       // Product Code
			row.AddCell().Value = i.Product.Name                       // Product Name
			row.AddCell().Value = i.Product.Uom.Name                   // UOM
			row.AddCell().Value = i.Product.Category.Name              // Category
			row.AddCell().Value = i.PriceSet.Name                      // Price Set
			row.AddCell().SetInt(int(i.Product.Salability))            // Salable
			row.AddCell().SetFloatWithFormat(i.UnitPrice, "0.00")      // Unit Price
			row.AddCell().SetFloatWithFormat(i.ShadowPrice, "0.00")    // Shadow Price
			row.AddCell().Value = strconv.Itoa(i.ShadowPricePct) + "%" // Shadow Percentage

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
	}
	return
}

// Update: function for update product price (upload)
func Update(r updateRequest) (we []model.Price, e error) {
	o := orm.NewOrm()
	var err error
	var lp []model.Price
	o.Begin()

	for _, v := range r.UpdateProductPrice {

		p := &model.Price{Product: v.Product, PriceSet: v.PriceSet}
		err = p.Read("Product", "PriceSet")

		if err == nil {
			p = &model.Price{
				ID:             p.ID,
				UnitPrice:      math.Round(v.UnitPrice),
				ShadowPrice:    p.ShadowPrice,
				ShadowPricePct: p.ShadowPricePct,
			}

			if _, err := o.Update(p, "unit_price"); err != nil {
				o.Rollback()
				return nil, err
			}

			err := log.AuditLogByUser(r.Session.Staff, p.ID, "product price", "update", "Update Product Price")
			if err != nil {
				o.Rollback()
				return
			}

			pl := &model.PriceLog{
				PriceID:   p.ID,
				UnitPrice: p.UnitPrice,
				CreatedAt: time.Now(),
				CreatedBy: r.Session.Staff,
			}
			if _, err = o.Insert(pl); err != nil {
				o.Rollback()
				return nil, err
			}
			r.Price = p
			lp = append(lp, *p)
		}
	}

	o.Commit()
	return lp, err
}

// UpdateShadow: function for update shadow product price (upload)
func UpdateShadow(r shadowRequest) (we []model.Price, e error) {
	o := orm.NewOrm()
	var err error
	var lp []model.Price
	var shadowPct float64
	o.Begin()

	for _, v := range r.UpdateShadowPrice {

		p := &model.Price{Product: v.Product, PriceSet: v.PriceSet}
		err = p.Read("Product", "PriceSet")

		if v.ShadowPrice != 0 {
			shadowPct = (v.ShadowPrice - p.UnitPrice) / v.ShadowPrice
		} else {
			shadowPct = 0
		}

		if err == nil {
			p = &model.Price{
				ID:             p.ID,
				UnitPrice:      p.UnitPrice,
				ShadowPrice:    math.Round(v.ShadowPrice),
				ShadowPricePct: int(math.Round(shadowPct * 100)),
			}

			if _, err := o.Update(p, "shadow_price", "shadow_price_pct"); err != nil {
				o.Rollback()
				return nil, err
			}

			err := log.AuditLogByUser(r.Session.Staff, p.ID, "product price", "shadow", "Update Shadow Price")
			if err != nil {
				o.Rollback()
				return
			}
			r.Price = p
			lp = append(lp, *p)
		}
	}

	o.Commit()
	return lp, err
}
