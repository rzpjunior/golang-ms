// Copyright 2020 PT. Eden Pangan Indonesia. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

package util

import (
	"fmt"
	"math/rand"
	"strconv"
	"strings"

	"time"

	"git.edenfarm.id/cuxs/orm"
	"github.com/labstack/echo/v4"
)

// check table for generate code
func CheckTable(tableName string) (initial string, e error) {
	if tableName == "archetype" {
		tableName = "ARC"
	} else if tableName == "area" {
		tableName = "ARE"
	} else if tableName == "business_type" {
		tableName = "BTY"
	} else if tableName == "category" {
		tableName = "CTG"
	} else if tableName == "city" {
		tableName = "CTY"
	} else if tableName == "country" {
		tableName = "COU"
	} else if tableName == "courier" {
		tableName = "CRR"
	} else if tableName == "user_courier" {
		tableName = "UCR"
	} else if tableName == "tag_customer" {
		tableName = "CTA"
	} else if tableName == "delivery_order" {
		tableName = "SJ"
	} else if tableName == "district" {
		tableName = "DIS"
	} else if tableName == "division" {
		tableName = "DIV"
	} else if tableName == "term_invoice_sls" {
		tableName = "INT"
	} else if tableName == "user_merchant" {
		tableName = "UMC"
	} else if tableName == "notification" {
		tableName = "NOT"
	} else if tableName == "payment_method" {
		tableName = "PYM"
	} else if tableName == "permission" {
		tableName = "PMS"
	} else if tableName == "price_set" {
		tableName = "PCS"
	} else if tableName == "product" {
		tableName = "PRD"
	} else if tableName == "prospect_customer" {
		tableName = "PCT"
	} else if tableName == "prospect_supplier" {
		tableName = "PSP"
	} else if tableName == "province" {
		tableName = "PRV"
	} else if tableName == "term_payment_pur" {
		tableName = "PPT"
	} else if tableName == "role" {
		tableName = "ROL"
	} else if tableName == "term_payment_sls" {
		tableName = "SPT"
	} else if tableName == "staff" {
		tableName = "STF"
	} else if tableName == "sub_district" {
		tableName = "SDS"
	} else if tableName == "supplier_type" {
		tableName = "SUT"
	} else if tableName == "uom" {
		tableName = "UOM"
	} else if tableName == "user" {
		tableName = "USR"
	} else if tableName == "courier_vehicle" {
		tableName = "CVH"
	} else if tableName == "courier_vendor" {
		tableName = "CVN"
	} else if tableName == "voucher" {
		tableName = "VOU"
	} else if tableName == "warehouse" {
		tableName = "WRH"
	} else if tableName == "wrt" {
		tableName = "WRT"
	} else if tableName == "tag_product" {
		tableName = "PTA"
	} else if tableName == "merchant" {
		tableName = "M"
	} else if tableName == "branch" {
		tableName = "B"
	} else if tableName == "packing_order" {
		tableName = "PC"
	} else if tableName == "picking_order" {
		tableName = "PIO"
	} else if tableName == "koli" {
		tableName = "KOL"
	} else if tableName == "sales_group" {
		tableName = "SGR"
	} else if tableName == "sales_assignment" {
		tableName = "SLA"
	} else if tableName == "supplier_commodity" {
		tableName = "SUC"
	} else if tableName == "supplier_badge" {
		tableName = "SUB"
	} else if tableName == "supplier_organization" {
		tableName = "SOR"
	} else if tableName == "sku_discount" {
		tableName = "SKD"
	} else if tableName == "transfer_sku" {
		tableName = "TF"
	} else if tableName == "bin" {
		tableName = "BIN"
	} else if tableName == "purchase_plan" {
		tableName = "PP"
	} else if tableName == "field_purchase_order" {
		tableName = "FPO"
	} else if tableName == "purchase_deliver" {
		tableName = "PD"
	} else if tableName == "consolidated_purchase_deliver" {
		tableName = "CP"
	} else if tableName == "notification_campaign" {
		tableName = "PNT-" + time.Now().Format("0601")
	} else if tableName == "banner" {
		tableName = "BNR"
	} else if tableName == "consolidated_shipment" {
		tableName = "CS"
	} else if tableName == "eden_point_campaign" {
		tableName = "EPC"
	} else if tableName == "product_section" {
		tableName = "PSC"
	} else if tableName == "sales_assignment_objective" {
		tableName = "SOB"
	} else if tableName == "user_fridge" {
		tableName = "USF"
	} else {
		e = echo.ErrNotFound
	}

	return tableName, e
}

// for generate code general on has 3 prefix in first code like USR, UCR, etc..
func GenerateCode(format string, tableName string, codeLength ...int) (code string, e error) {
	var initialCode, template string
	var lenAbbr, lenInitialCode int
	orSelect := orm.NewOrm()
	orSelect.Using("read_only")

	codeLen := 4
	if len(codeLength) > 0 {
		codeLen = codeLength[0]
	}
	template = format + "#" + fmt.Sprintf("%0"+strconv.Itoa(codeLen)+"d", 1)
	lenAbbr = len(format)
	if e == nil {
		if template != "" {
			if e = orSelect.Raw("SELECT code FROM code_generator WHERE code_name = ? AND code LIKE ? ORDER BY id DESC LIMIT 1", tableName, format+"%").QueryRow(&initialCode); e == nil {
				lenInitialCode = len(initialCode)
				tempIncrement := initialCode[lenAbbr:lenInitialCode]
				increment, _ := strconv.Atoi(tempIncrement)
				increments := fmt.Sprintf("%0"+strconv.Itoa(codeLen)+"d", increment+1)
				code = fmt.Sprintf("%s%s", format, increments)
			} else {
				code = fmt.Sprintf("%s%s", format, fmt.Sprintf("%0"+strconv.Itoa(codeLen)+"d", 1))
			}
		}
	}

	// simpan code dokumen ke table code_generator untuk cek duplikat atau tidak
	_, e = orm.NewOrm().Raw("INSERT INTO `code_generator` (`code`, `code_name`) VALUES (?, ?);", code, tableName).Exec()
	if e != nil {
		var isDuplicate = strings.Contains(e.Error(), "Duplicate entry")
		if isDuplicate {
			code, e = GenerateCode(format, tableName)
		}

	}

	return code, e
}

func GenerateDocCode(firstCode string, middleCode string, docType string) (code string, e error) {
	var codeExist, template string
	orSelect := orm.NewOrm()
	orSelect.Using("read_only")
	template = firstCode + "#-#CODE#-#0220#0001"
	if e == nil {
		if template != "" {
			if e = orSelect.Raw(`SELECT code FROM code_generator WHERE code_name = ?
				AND code LIKE ? ORDER BY id DESC LIMIT 1`, docType, fmt.Sprint(firstCode, "-", middleCode+"-", "%")).QueryRow(&codeExist); e == nil {
				tempInitial := strings.Split(codeExist, "-"+middleCode+"-")
				if len(tempInitial) != 2 {
					year, month, _ := time.Now().Date()
					years := strconv.Itoa(year)
					tempDate := fmt.Sprintf("%02d%s", month, years[2:4])
					code = fmt.Sprintf("%s%s%s%s%s%s", firstCode, "-", middleCode, "-", tempDate, "0001")
				} else {
					year, month, _ := time.Now().Date()
					tempDate := string([]rune(tempInitial[1]))[0:4]
					tempIncrement := string([]rune(tempInitial[1]))[4:len(tempInitial[1])]
					tempMonth, _ := strconv.Atoi(tempDate[0:2])
					if int(month) == tempMonth {
						increment, _ := strconv.Atoi(tempIncrement)
						increments := fmt.Sprintf("%0"+strconv.Itoa(4)+"d", increment+1)
						code = fmt.Sprintf("%s%s%s%s%s%s", firstCode, "-", middleCode, "-", tempDate, increments)
					} else {
						years := strconv.Itoa(year)
						tempDate = fmt.Sprintf("%02d%s", month, years[2:4])
						code = fmt.Sprintf("%s%s%s%s%s%s", firstCode, "-", middleCode, "-", tempDate, "0001")
					}
				}
			} else {
				templates := strings.Split(template, "#")
				year, month, _ := time.Now().Date()
				years := strconv.Itoa(year)
				tempDate := fmt.Sprintf("%02d%s", month, years[2:4])
				code = fmt.Sprintf("%s%s%s%s%s%s", templates[0], "-", middleCode, "-", tempDate, "0001")
			}
		} else {
			year, month, _ := time.Now().Date()
			years := strconv.Itoa(year)
			tempDate := fmt.Sprintf("%02d%s", month, years[2:4])
			code = fmt.Sprintf("%s%s%s%s%s%s", firstCode, "-", middleCode, "-", tempDate, "0001")
		}
	}

	// simpan code dokumen ke table code_generator untuk cek duplikat atau tidak
	_, e = orm.NewOrm().Raw("INSERT INTO `code_generator` (`code`, `code_name`) VALUES (?, ?);", code, docType).Exec()
	if e != nil {
		var isDuplicate = strings.Contains(e.Error(), "Duplicate entry")
		if isDuplicate {
			code, e = GenerateDocCode(firstCode, middleCode, docType)
		}
	}

	return code, e
}

// GenerateCustomerCode : function to generate new code for customer
func GenerateCustomerCode(tableCode, tableName string) (code string, e error) {
	randCode := GenerateRandomString("ABCDEFGHIJKLMNOPQRSTUVWXYZ", 3)
	randNumber := GenerateRandomString("0123456789", 4)

	code = tableCode + randCode + randNumber

	// simpan code dokumen ke table code_generator untuk cek duplikat atau tidak
	if _, e = orm.NewOrm().Raw("insert into `code_generator` (`code`, `code_name`) values (?, ?)", code, tableName).Exec(); e != nil {
		var isDuplicate = strings.Contains(e.Error(), "Duplicate entry")
		if isDuplicate {
			code, e = GenerateCustomerCode(tableCode, tableName)
		}
	}

	return
}

func GenerateRandomString(charset string, length int) string {
	byteString := make([]byte, length)
	for i := range byteString {
		byteString[i] = charset[rand.Intn(len(charset))]
	}

	return string(byteString)
}

// GenerateCodeReferral : function to generate referral code in table merchant
func GenerateCodeReferral() (code string, e error) {
	code = GenerateRandomString("ABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789", 8)

	if _, e = orm.NewOrm().Raw("insert into code_generator_referral (code) values (?)", code).Exec(); e != nil {
		var isDuplicate = strings.Contains(e.Error(), "Duplicate entry")
		if isDuplicate {
			code, e = GenerateCodeReferral()
		}
	}

	return
}

// GenerateRangeDates : function to generate dates between two dates
func GenerateRangeDates(start, end time.Time) func() time.Time {
	y, m, d := start.Date()
	start = time.Date(y, m, d, 0, 0, 0, 0, time.UTC)
	y, m, d = end.Date()
	end = time.Date(y, m, d, 0, 0, 0, 0, time.UTC)

	return func() time.Time {
		if start.After(end) {
			return time.Time{}
		}
		date := start
		start = start.AddDate(0, 0, 1)
		return date
	}
}
