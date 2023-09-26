// Copyright 2020 PT. Eden Pangan Indonesia. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

package repository

import (
	"strings"
	"time"

	"git.edenfarm.id/project-version2/datamodel/model"
	"git.edenfarm.id/project-version2/api/util"

	"git.edenfarm.id/cuxs/orm"
)

// GetForecastDemand : find a single data using field and value condition.
func GetForecastDemand(field string, values ...interface{}) (*model.ForecastDemand, error) {
	m := new(model.ForecastDemand)
	o := orm.NewOrm()
	o.Using("read_only")

	if err := o.QueryTable(m).Filter(field, values...).RelatedSel().Limit(1).One(m); err != nil {
		return nil, err
	}

	return m, nil
}

// GetForecastDemands : function to get data from database based on parameters
func GetForecastDemands(rq *orm.RequestQuery) (fd []*model.ForecastDemand, total int64, err error) {
	q, _ := rq.QueryReadOnly(new(model.ForecastDemand))

	if total, err = q.All(&fd, rq.Fields...); err == nil {
		return fd, total, nil
	}

	return nil, 0, err
}

// ValidForecastDemand : function to check if id is valid in database
func ValidForecastDemand(id int64) (fd *model.ForecastDemand, e error) {
	fd = &model.ForecastDemand{ID: id}
	e = fd.Read("ID")

	return
}

// GetForecastDemandsForExport : function to get data from database based on parameters to be exported to xls
func GetForecastDemandsForExport(rq *orm.RequestQuery) (fdx []orm.Params, total int64, arrDate []string, err error) {
	o := orm.NewOrm()
	o.Using("read_only")

	var tempTable, where, addSelect string
	var valuesArr []string

	// double loop to set up string of query of where clause based on parameters
	for _, conditions := range rq.Conditions {
		for i, value := range conditions {

			if strings.Contains(i, ".between") {
				tab := "tmp_tab"
				col := strings.Split(i, ".")

				if len(col) >= 3 {
					tab = col[len(col)-3]
				}

				val := strings.Split(value, ".")

				// set up query to create temporary table that consists of dates ranging between two dates
				start, _ := time.Parse("2006-01-02", val[0])
				end, _ := time.Parse("2006-01-02", val[1])
				tempTable = "select '" + val[0] + "' forecast_date union "
				addSelect += "sum(case when forecast_demand.forecast_date = '" + val[0] + "' then forecast_demand.forecast_qty else 0 end) '" + val[0] + "', "
				arrDate = append(arrDate, val[0])
				for rd := util.GenerateRangeDates(start.AddDate(0, 0, 1), end); ; {
					date := rd()
					if date.IsZero() {
						break
					}
					tempTable += "select '" + date.Format("2006-01-02") + "' union "
					addSelect += "sum(case when forecast_demand.forecast_date = '" + date.Format("2006-01-02") + "' then forecast_demand.forecast_qty else 0 end) '" + date.Format("2006-01-02") + "', "
					arrDate = append(arrDate, date.Format("2006-01-02"))
				}
				tempTable = strings.TrimSuffix(tempTable, " union ")

				where += tab + "." + col[len(col)-2] + " " + col[len(col)-1] + " '" + val[0] + "' and '" + val[1] + "' and "
			} else {
				tab := "forecast_demand"
				col := strings.Split(i, ".")
				if len(col) >= 2 {
					tab = col[len(col)-2]
				}

				where += tab + "." + col[len(col)-1] + " = ? and "
				valuesArr = append(valuesArr, value)
			}
		}
	}
	where = strings.TrimSuffix(where, " and ")

	selectCol := "select warehouse.code warehouse_code, warehouse.name warehouse_name, category.name product_category, product.code product_code, product.name product_name, uom.name uom, " + addSelect
	selectCol = strings.TrimSuffix(selectCol, ", ") + " "

	query := selectCol +
		"from stock " +
		"cross join (" + tempTable + ") tmp_tab " +
		"left join forecast_demand on stock.product_id = forecast_demand.product_id and stock.warehouse_id = forecast_demand.warehouse_id and tmp_tab.forecast_date = forecast_demand.forecast_date " +
		"join product on stock.product_id = product.id " +
		"join warehouse on stock.warehouse_id = warehouse.id " +
		"join category on product.category_id = category.id " +
		"join uom on product.uom_id = uom.id " +
		"where stock.status = 1 and " + where + " " +
		"group by stock.product_id, stock.warehouse_id " +
		"order by tmp_tab.forecast_date, product.id"

	if total, err = o.Raw(query, valuesArr).Values(&fdx); err == nil {
		return fdx, total, arrDate, nil
	}

	return nil, 0, nil, err
}
