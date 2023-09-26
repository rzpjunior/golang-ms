// Copyright 2020 PT. Eden Pangan Indonesia. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

package fulfillment

import (
	"fmt"
	"os"
	"strconv"
	"strings"
	"time"

	"git.edenfarm.id/cuxs/env"
	"git.edenfarm.id/project-version2/datamodel/model"
	"git.edenfarm.id/project-version2/api/util"
	"github.com/tealeg/xlsx"
)

// PrintReportFulfillmentXls : function to create report fulfillment dashboard in xls
func PrintReportFulfillmentXls(data []*model.ReportFulfillment, warehouse *model.Warehouse, year string, lastUpdatedAt time.Time) (filePath string, err error) {
	var (
		file    *xlsx.File
		sheet   *xlsx.Sheet
		row     *xlsx.Row
		weekStr string
		weekInt int
	)

	dir := env.GetString("EXPORT_DIRECTORY", "")

	filename := fmt.Sprintf("ReportFulfillment_%s_%s_%s.xlsx", year, warehouse.Name, util.GenerateRandomDoc(5))
	fileDir := fmt.Sprintf("%s/%s", dir, filename)

	file = xlsx.NewFile()
	if sheet, err = file.AddSheet("Sheet1"); err == nil {
		row = sheet.AddRow()
		row.AddCell().Value = "Downloaded at : " + time.Now().Format("02/01/2006 15:04:05")

		row = sheet.AddRow()
		row.AddCell().Value = "Data Last Updated at : " + lastUpdatedAt.Format("02/01/2006 15:04:05")

		row = sheet.AddRow()
		row.AddCell().Value = "Year : " + year

		row = sheet.AddRow()
		row.AddCell().Value = "Warehouse : " + warehouse.Name

		sheet.AddRow()

		row = sheet.AddRow()
		row.AddCell().Value = "Week Number"
		row.AddCell().Value = "Date Range"
		row.AddCell().Value = "SO Fulfillment Rate"

		row.Sheet.SetColWidth(1, 1, 15)

		for _, v := range data {
			row = sheet.AddRow()
			if year == "2021" {
				weekStr = strings.TrimPrefix(v.WeekNumber[len(v.WeekNumber)-2:], "0")
				weekInt, _ = strconv.Atoi(weekStr)
				weekStr = strconv.Itoa(weekInt + 1)
			}

			row.AddCell().Value = weekStr                                                                   // Week Number
			row.AddCell().Value = v.StartDate.Format("02/01/2006") + " - " + v.EndDate.Format("02/01/2006") // Date Range
			row.AddCell().SetFloatWithFormat(v.FulfillmentRate/100, "0.00%")                                // SO Fulfillment Rate
		}

		boldStyle := xlsx.NewStyle()
		boldFont := xlsx.NewFont(10, "Liberation Sans")
		boldFont.Bold = true
		boldStyle.Font = *boldFont
		boldStyle.ApplyFont = true

		rightAlignmentStyle := xlsx.NewStyle()
		rightAlignmentStyle.Alignment.Horizontal = "right"
		rightAlignmentStyle.ApplyAlignment = true

		// looping to make BOLD font header
		for col := 0; col < 3; col++ {
			sheet.Cell(5, col).SetStyle(boldStyle)
		}

		// looping to make BOLD font for the first 4 rows
		for row := 0; row < 4; row++ {
			sheet.Cell(row, 0).SetStyle(boldStyle)
		}

		// looping to set right alignment of first column
		for row := 6; row < 60; row++ {
			sheet.Cell(row, 0).SetStyle(rightAlignmentStyle)
		}

		err = file.Save(fileDir)
		filePath, err = util.UploadToS3(filename, fileDir, "application/xlsx")

		// fungsi ini berguna untuk menghapus kembali file yang sudah diupload ke server
		os.Remove(fileDir)
	}

	return
}
