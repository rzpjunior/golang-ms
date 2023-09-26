package reportx

import (
	"fmt"
	"os"
	"reflect"
	"time"

	"github.com/tealeg/xlsx"
)

type Excelx struct {
	Sheets []Sheet
}

type Sheet struct {
	WithNumbering bool
	Name          string
	Headers       []string
	Bodys         []interface{}
}

func InterfaceSlice(slice interface{}) []interface{} {
	v := reflect.ValueOf(slice)

	ret := make([]interface{}, v.NumField())
	for i := 0; i < v.NumField(); i++ {
		ret[i] = v.Field(i).Interface()
	}

	return ret
}

func GenerateXlsx(fileName string, excelx Excelx) (fileDir string, err error) {
	var file *xlsx.File
	var sheet *xlsx.Sheet
	var row *xlsx.Row

	dir, err := os.Getwd()
	if err != nil {
		return
	}

	fileDir = fmt.Sprintf("%s/%s", dir, fileName)

	file = xlsx.NewFile()
	for _, excelxSheet := range excelx.Sheets {
		if sheet, err = file.AddSheet(excelxSheet.Name); err == nil {

			// set header
			row = sheet.AddRow()
			row.SetHeight(20)

			if excelxSheet.WithNumbering {
				row.AddCell().Value = "No"
			}

			for _, excelxHeader := range excelxSheet.Headers {
				row.AddCell().Value = excelxHeader
			}

			// set body
			for i, excelxRow := range excelxSheet.Bodys {
				row = sheet.AddRow()

				if excelxSheet.WithNumbering {
					no := i + 1
					row.AddCell().SetInt(no)
				}

				var cells []interface{}

				fmt.Println("reflect.TypeOf(excelxRow).Kind(): ", reflect.TypeOf(excelxRow).Kind())
				if reflect.TypeOf(excelxRow).Kind() == reflect.Map || reflect.TypeOf(excelxRow).Kind() == reflect.Slice {
					v := reflect.ValueOf(excelxRow)
					if v.Kind() == reflect.Map {
						for _, key := range v.MapKeys() {
							strct := v.MapIndex(key)
							cells = append(cells, strct.Interface())
						}
					} else if v.Kind() == reflect.Slice {
						for i := 0; i < v.Len(); i++ {
							cells = append(cells, v.Index(i).Interface())
						}
					}
				} else {
					cells = InterfaceSlice(excelxRow)
				}

				for _, excelxCell := range cells {
					cellType := reflect.TypeOf(excelxCell)
					if cellType.Kind() == reflect.Int {
						row.AddCell().SetInt(excelxCell.(int))
					} else if cellType.Kind() == reflect.Int8 {
						row.AddCell().SetInt(int(excelxCell.(int8)))
					} else if cellType.Kind() == reflect.Int16 {
						row.AddCell().SetInt(int(excelxCell.(int16)))
					} else if cellType.Kind() == reflect.Int32 {
						row.AddCell().SetInt(int(excelxCell.(int32)))
					} else if cellType.Kind() == reflect.Int64 {
						row.AddCell().SetInt64(excelxCell.(int64))
					} else if cellType.Kind() == reflect.Float32 || cellType.Kind() == reflect.Float64 {
						row.AddCell().SetFloat(excelxCell.(float64))
					} else if cellType.Kind() == reflect.Bool {
						row.AddCell().SetBool(excelxCell.(bool))
					} else if _, ok := excelxCell.(time.Time); ok {
						row.AddCell().SetDateTime(excelxCell.(time.Time))
					} else {
						row.AddCell().Value = excelxCell.(string)
					}
				}
			}

		}

	}

	err = file.Save(fileDir)
	return
}
