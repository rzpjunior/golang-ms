// Copyright 2020 PT. Eden Pangan Indonesia. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

package datascraping

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

// SavePublicProduct1 : function to save scraped public data 1 into scrape db
func SavePublicProduct1(items []*model.PublicData1Items, area *model.PublicArea1) (e error) {
	o := orm.NewOrm()
	o.Using("scrape")
	o.Begin()

	var isCreated bool

	for _, v := range items {
		productImages := "["
		for _, v := range v.Product.ProductImages {
			productImages += `"` + v.ImageURL + `",`
		}
		productImages = strings.TrimSuffix(productImages, ",")
		productImages += "]"

		productName := strings.ReplaceAll(v.Product.Name, `"`, `\"`)

		publicProduct := &model.PublicProduct1{
			ID:            v.Product.ID,
			UOM:           v.Product.Unit.Description,
			Name:          productName,
			ProductImages: productImages,
			CreatedAt:     time.Now(),
		}

		if isCreated, publicProduct.ID, e = o.ReadOrCreate(publicProduct, "ID"); e == nil {
			if !isCreated && (publicProduct.UOM != v.Product.Unit.Description || publicProduct.Name != productName || publicProduct.ProductImages != productImages) {
				publicProduct.UOM = v.Product.Unit.Description
				publicProduct.Name = productName
				publicProduct.ProductImages = productImages
				publicProduct.LastUpdatedAt = time.Now()
				if _, e = o.Update(publicProduct, "UOM", "Name", "ProductImages", "LastUpdatedAt"); e != nil {
					_, e = o.Raw("insert into error_log (public_price_set, public_product, table_name, error_message) values (?, ?, ?, ?)", 1, v.Product.ID, "public_product", e.Error()).Exec()
				}
			}
		} else {
			_, e = o.Raw("insert into error_log (public_price_set, public_product, table_name, error_message) values (?, ?, ?, ?)", 1, v.Product.ID, "public_product", e.Error()).Exec()
		}

		publicPrice := &model.PublicPrice1{
			Product:            publicProduct,
			ScrapedDate:        time.Now().Format("2006-01-02"),
			Area:               area,
			Price:              v.ShowedPrice,
			Discount:           v.DiscountInRupiah,
			PriceAfterDiscount: v.ShowedPriceAfterDiscount,
			CreatedAt:          time.Now(),
		}

		if isCreated, publicPrice.ID, e = o.ReadOrCreate(publicPrice, "Product", "ScrapedDate", "Area"); e == nil {
			if !isCreated && (publicPrice.Price != float64(v.ShowedPrice) || publicPrice.Discount != float64(v.DiscountInRupiah)) {
				publicPrice.Price = v.ShowedPrice
				publicPrice.Discount = v.DiscountInRupiah
				publicPrice.PriceAfterDiscount = v.ShowedPriceAfterDiscount
				publicPrice.LastUpdatedAt = time.Now()
				if _, e = o.Update(publicPrice, "Price", "Discount", "PriceAfterDiscount", "LastUpdatedAt"); e != nil {
					_, e = o.Raw("insert into error_log (public_price_set, public_product, table_name, error_message) values (?, ?, ?, ?)", 1, v.Product.ID, "public_price", e.Error()).Exec()
				}
			}
		} else {
			_, e = o.Raw("insert into error_log (public_price_set, public_product, table_name, error_message) values (?, ?, ?, ?)", 1, v.Product.ID, "public_price", e.Error()).Exec()
		}
	}

	o.Commit()
	return
}

// SavePublicProduct2 : function to save scraped public data 2 into scrape db
func SavePublicProduct2(items []*model.PublicData2Items, area *model.PublicArea2) (e error) {
	o := orm.NewOrm()
	o.Using("scrape")
	o.Begin()

	var isCreated bool
	var uom string

	for _, v := range items {
		if v.PackDesc == v.PackNote {
			uom = v.PackDesc
		} else {
			uom = v.PackDesc + " (" + v.PackNote + ")"
		}

		publicProduct := &model.PublicProduct2{
			ProductKey:    v.Key,
			UOM:           uom,
			Name:          v.Name,
			ProductImages: v.Image.Lg,
			CreatedAt:     time.Now(),
		}

		if isCreated, publicProduct.ID, e = o.ReadOrCreate(publicProduct, "ProductKey"); e == nil {
			if !isCreated && (publicProduct.UOM != uom || publicProduct.Name != v.Name || publicProduct.ProductImages != v.Image.Lg) {
				publicProduct.UOM = uom
				publicProduct.Name = v.Name
				publicProduct.ProductImages = v.Image.Lg
				publicProduct.LastUpdatedAt = time.Now()
				if _, e = o.Update(publicProduct, "UOM", "Name", "ProductImages", "LastUpdatedAt"); e != nil {
					_, e = o.Raw("insert into error_log (public_price_set, public_product, table_name, error_message) values (?, ?, ?, ?)", 2, v.Key, "public_product", e.Error()).Exec()
				}
			}
		} else {
			_, e = o.Raw("insert into error_log (public_price_set, public_product, table_name, error_message) values (?, ?, ?, ?)", 2, v.Key, "public_product", e.Error()).Exec()
		}

		actualPrice, _ := strconv.ParseFloat(v.ActualPrice, 64)
		price, _ := strconv.ParseFloat(v.Price, 64)
		publicPrice := &model.PublicPrice2{
			Product:            publicProduct,
			ScrapedDate:        time.Now().Format("2006-01-02"),
			Area:               area,
			Price:              actualPrice,
			Discount:           actualPrice - price,
			PriceAfterDiscount: price,
			CreatedAt:          time.Now(),
		}

		if isCreated, publicPrice.ID, e = o.ReadOrCreate(publicPrice, "Product", "ScrapedDate", "Area"); e == nil {
			if !isCreated && (publicPrice.Price != actualPrice || publicPrice.Discount != (actualPrice-price)) {
				publicPrice.Price = actualPrice
				publicPrice.Discount = actualPrice - price
				publicPrice.PriceAfterDiscount = price
				publicPrice.LastUpdatedAt = time.Now()
				if _, e = o.Update(publicPrice, "Price", "Discount", "PriceAfterDiscount", "LastUpdatedAt"); e != nil {
					_, e = o.Raw("insert into error_log (public_price_set, public_product, table_name, error_message) values (?, ?, ?, ?)", 2, v.Key, "public_price", e.Error()).Exec()
				}
			}
		} else {
			_, e = o.Raw("insert into error_log (public_price_set, public_product, table_name, error_message) values (?, ?, ?, ?)", 2, v.Key, "public_price", e.Error()).Exec()
		}
	}

	o.Commit()
	return
}

// SaveDashboardProduct : function to save scraped dashboard data into scrape db
func SaveDashboardProduct(matchedArea []*model.MatchedArea) (e error) {
	var resData []*model.DataDashboard
	var isCreated bool

	for _, v := range matchedArea {
		o := orm.NewOrm()
		o.Using("read_only")
		query := "select pro.code product_code, pro.name product_name, u.name uom, pri.name area, pro.id product_id, coalesce(pri.unit_price, 0) price " +
			"from product pro " +
			"join uom u on pro.uom_id = u.id " +
			"left join " +
			"( " +
			"	select pri.*, ps.name " +
			"	from price pri " +
			"	join " +
			"	( " +
			"		select id, name " +
			"		from price_set " +
			"		where name = ? " +
			"	) ps on ps.id = pri.price_set_id  " +
			") pri on pro.id = pri.product_id " +
			"where pro.status = 1 "
		if _, e = o.Raw(query, "High "+v.DashboardArea.Name).QueryRows(&resData); e == nil {
			ormScrape := orm.NewOrm()
			ormScrape.Using("scrape")
			ormScrape.Begin()

			areaName := v.DashboardArea.Name
			for _, v := range resData {
				dashboardProduct := &model.DashboardProduct{
					ID:        v.ProductID,
					Code:      v.ProductCode,
					Name:      v.ProductName,
					UOM:       v.UOM,
					CreatedAt: time.Now(),
				}

				if isCreated, dashboardProduct.ID, e = ormScrape.ReadOrCreate(dashboardProduct, "ID"); e == nil {
					if !isCreated && (dashboardProduct.Code != v.ProductCode || dashboardProduct.Name != v.ProductName || dashboardProduct.UOM != v.UOM) {
						dashboardProduct.Code = v.ProductCode
						dashboardProduct.Name = v.ProductName
						dashboardProduct.UOM = v.UOM
						dashboardProduct.LastUpdatedAt = time.Now()

						if _, e = ormScrape.Update(dashboardProduct, "Code", "Name", "UOM", "LastUpdatedAt"); e != nil {
							ormScrape.Rollback()
							return e
						}
					}
				}

				if v.Price > 0 {
					area := &model.DashboardArea{Name: areaName}
					area.Read("Name")

					dashboardPrice := &model.DashboardPrice{
						Area:        area,
						Product:     dashboardProduct,
						ScrapedDate: time.Now().Format("2006-01-02"),
						Price:       v.Price,
						CreatedAt:   time.Now(),
					}

					if isCreated, dashboardPrice.ID, e = ormScrape.ReadOrCreate(dashboardPrice, "Area", "Product", "ScrapedDate"); e == nil {
						if !isCreated && dashboardPrice.Price != v.Price {
							dashboardPrice.Price = v.Price
							dashboardPrice.LastUpdatedAt = time.Now()

							if _, e = ormScrape.Update(dashboardPrice, "Price", "LastUpdatedAt"); e != nil {
								ormScrape.Rollback()
								return e
							}
						}
					}
				}
			}

			ormScrape.Commit()
		}
	}

	return nil
}

// DownloadPublicProductXls : download list public product for update
func DownloadPublicProductXls(date time.Time, r []*model.PublicProductForXls, priceSet int64) (filePath string, err error) {
	var file *xlsx.File
	var sheet *xlsx.Sheet
	var row *xlsx.Row

	dir := env.GetString("EXPORT_DIRECTORY", "")

	filename := fmt.Sprintf("PublicProduct%d_%s_%s.xlsx", priceSet, date.Format("2006-01-02"), util.GenerateRandomDoc(5))

	fileDir := fmt.Sprintf("%s/%s", dir, filename)

	file = xlsx.NewFile()
	if sheet, err = file.AddSheet("Sheet1"); err == nil {
		row = sheet.AddRow()
		row.SetHeight(20)
		row.AddCell().Value = "No"
		row.AddCell().Value = "Product_Name"
		row.AddCell().Value = "UOM"
		row.AddCell().Value = "Product_Images"
		row.AddCell().Value = "Product_Edenfarm"

		for i, v := range r {
			row = sheet.AddRow()
			row.AddCell().SetInt(i + 1)                                                                                                  // No
			row.AddCell().Value = v.ProductName                                                                                          // Product Name
			row.AddCell().Value = v.UOM                                                                                                  // UOM
			row.AddCell().Value = strings.ReplaceAll(strings.ReplaceAll(strings.ReplaceAll(v.ProductImages, "[", ""), "]", ""), `"`, "") // Product Images
			row.AddCell().Value = v.DashboardProductName                                                                                 // Product Dashboard
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

// DownloadProductMatchingXls : download list of dashboard product to be matched with public product
func DownloadProductMatchingXls(date time.Time, r []*model.ProductMatchingTemplate) (filePath string, err error) {
	var file *xlsx.File
	var sheet *xlsx.Sheet
	var row *xlsx.Row

	dir := env.GetString("EXPORT_DIRECTORY", "")

	filename := fmt.Sprintf("TemplateMatchingProduct_%s_%s.xlsx", date.Format("2006-01-02"), util.GenerateRandomDoc(5))

	fileDir := fmt.Sprintf("%s/%s", dir, filename)

	file = xlsx.NewFile()
	if sheet, err = file.AddSheet("Sheet1"); err == nil {
		row = sheet.AddRow()
		row.SetHeight(20)
		row.AddCell().Value = "No"
		row.AddCell().Value = "Eden_Product_Code"
		row.AddCell().Value = "Eden_Product_Name"
		row.AddCell().Value = "Public_Product_1"
		row.AddCell().Value = "Public_Product_2"

		for i, v := range r {
			row = sheet.AddRow()
			row.AddCell().SetInt(i + 1)                  // No
			row.AddCell().Value = v.DashboardProductCode // Eden Product Code
			row.AddCell().Value = v.DashboardProductName // Eden Product Name
			row.AddCell().Value = v.PublicProduct1       // Public Product 1
			row.AddCell().Value = v.PublicProduct2       // Public Product 2
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

// SaveUpdateMatching : save matched product
func SaveUpdateMatching(r matchingRequest) (successCount int64, e error) {
	o := orm.NewOrm()
	o.Using("scrape")
	o.Begin()

	var isCreated bool
	var colName []string

	for _, v := range r.Data {
		dashboardProduct := &model.DashboardProduct{Code: v.EdenProductCode}
		dashboardProduct.Read("Code")

		if v.PublicProduct1 == "" && v.PublicProduct2 == "" {
			o.Raw("delete from matched_product where dashboard_product_id = ?", dashboardProduct.ID).Exec()
		} else {
			publicProduct1 := &model.PublicProduct1{Name: v.PublicProduct1}
			publicProduct1.Read("Name")

			publicProduct2 := &model.PublicProduct2{Name: v.PublicProduct2}
			publicProduct2.Read("Name")

			matchedProduct := &model.MatchedProduct{
				DashboardProduct: dashboardProduct,
				PublicProduct1:   publicProduct1,
				PublicProduct2:   publicProduct2,
			}

			if publicProduct1.ID != 0 || publicProduct2.ID != 0 {
				if isCreated, matchedProduct.ID, e = o.ReadOrCreate(matchedProduct, "DashboardProduct"); e == nil {
					if !isCreated {
						matchedProduct.PublicProduct1 = publicProduct1
						colName = append(colName, "PublicProduct1")

						matchedProduct.PublicProduct2 = publicProduct2
						colName = append(colName, "PublicProduct2")

						if _, e = o.Update(matchedProduct, colName...); e == nil {
							successCount++
						}
					} else {
						successCount++
					}
				}
			}
		}
	}

	o.Commit()
	return
}
