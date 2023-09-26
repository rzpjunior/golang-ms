// Copyright 2020 PT. Eden Pangan Indonesia. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

package forecast_demand

import (
	"fmt"
	"os"
	"strconv"
	"strings"
	"time"

	"git.edenfarm.id/cuxs/env"
	"git.edenfarm.id/cuxs/orm"
	"git.edenfarm.id/project-version2/datamodel/model"
	"git.edenfarm.id/project-version2/api/util"
	"github.com/tealeg/xlsx"
)

// DownloadForecastDemandXls : download template Forecast Demand for update
func DownloadForecastDemandXls(date time.Time, r []orm.Params, arrDate []string) (filePath string, err error) {
	var file *xlsx.File
	var sheet *xlsx.Sheet
	var row *xlsx.Row

	dir := env.GetString("EXPORT_DIRECTORY", "")

	var warehouse string
	if len(r) != 0 {
		warehouse = r[0]["warehouse_name"].(string)
		warehouse = strings.ReplaceAll(warehouse, " ", "_")
	}

	filename := fmt.Sprintf("TemplateForecastDemand_%s_%s_%s.xlsx", date.Format("2006-01-02"), warehouse, util.GenerateRandomDoc(5))

	fileDir := fmt.Sprintf("%s/%s", dir, filename)

	file = xlsx.NewFile()
	if sheet, err = file.AddSheet("Sheet1"); err == nil {
		row = sheet.AddRow()
		row.SetHeight(20)
		row.AddCell().Value = "No"
		row.AddCell().Value = "Warehouse_Code"
		row.AddCell().Value = "Warehouse_Name"
		row.AddCell().Value = "Product_Category"
		row.AddCell().Value = "Product_Code"
		row.AddCell().Value = "Product_Name"
		row.AddCell().Value = "UOM"
		for _, v := range arrDate {
			row.AddCell().Value = v
		}

		for i, v := range r {
			row = sheet.AddRow()
			row.AddCell().SetInt(i + 1)                          // No
			row.AddCell().Value = v["warehouse_code"].(string)   // Warehouse Code
			row.AddCell().Value = v["warehouse_name"].(string)   // Warehouse Name
			row.AddCell().Value = v["product_category"].(string) // Category Name
			row.AddCell().Value = v["product_code"].(string)     // Product Code
			row.AddCell().Value = v["product_name"].(string)     // Product Name
			row.AddCell().Value = v["uom"].(string)              // UOM
			for _, val := range arrDate {
				floatVal, _ := strconv.ParseFloat(v[val].(string), 64)
				row.AddCell().SetFloatWithFormat(floatVal, "0.00")
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

	}

	return
}

// Update : function to update data of table forecast_demand
func Update(r updateRequest) error {
	o := orm.NewOrm()
	var e error
	var isCreated bool

	o.Begin()

	for _, v := range r.Data {
		if v.ForecastQty >= 0 {
			product := &model.Product{Code: v.ProductCode}
			product.Read("Code")

			warehouse := &model.Warehouse{Code: v.WarehouseCode}
			warehouse.Read("Code")

			if product.ID != 0 && warehouse.ID != 0 {
				forecastDemand := &model.ForecastDemand{
					Product:      product,
					Warehouse:    warehouse,
					ForecastDate: v.ForecastDate,
					ForecastQty:  v.ForecastQty,
				}

				if isCreated, forecastDemand.ID, e = o.ReadOrCreate(forecastDemand, "Product", "Warehouse", "ForecastDate"); e == nil {
					if !isCreated {
						forecastDemand.ForecastQty = v.ForecastQty
						if _, e = o.Update(forecastDemand, "ForecastQty"); e != nil {
							o.Rollback()
							return e
						}
					}
				} else {
					o.Rollback()
					return e
				}
			}
		}
	}

	o.Commit()
	return e
}
