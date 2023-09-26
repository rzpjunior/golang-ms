// Copyright 2020 PT. Eden Pangan Indonesia. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

package dashboard

import (
	"strings"

	"git.edenfarm.id/cuxs/orm"
	"git.edenfarm.id/project-version2/datamodel/model"
)

func GetOverviewByQuery(rq *orm.RequestQuery, date string) (m *model.DashboardOverview, total int64, err error) {

	// get data requested
	f := orm.NewOrm()
	f.Using("read_only")
	o := new(model.DashboardOverview)

	dateSeparator := strings.Split(date, "-")
	mo := dateSeparator[0]
	ye := dateSeparator[1]

	if err := f.Raw("SELECT count(id), sum(total_charge), sum(total_weight) from sales_order"+
		" where status not in (3,4)"+
		" and (month(recognition_date) = ? "+
		"and year(recognition_date) = ? "+
		")", mo, ye).QueryRow(&o.TotalTransaction, &o.SumTotalCharge, &o.TotalTonnage); err != nil {
		return nil, total, err
	}

	if _, err := f.Raw("select p.id, p.name, sum(soi.subtotal) as total from sales_order so"+
		" join sales_order_item soi on soi.sales_order_id = so.id"+
		" join product p on p.id = soi.product_id"+
		" where so.status not in (3,4)"+
		" and (month(so.recognition_date) = ?"+

		" and year(so.recognition_date) = ?"+

		") group by p.id"+
		" order by total desc"+
		" limit 5", mo, ye).QueryRows(&o.TopRevenue); err != nil {
		return nil, total, err
	}

	return o, total, nil

}

func GetGraphByQuery(rq *orm.RequestQuery, date string) (m []*model.DashboardGraph, total int64, err error) {

	// get data requested
	//o := new([]model.DashboardGraph)
	f := orm.NewOrm()
	f.Using("read_only")

	var o []*model.DashboardGraph
	var dateSO string

	dateSeparator := strings.Split(date, "-")
	mo := dateSeparator[0]
	ye := dateSeparator[1]
	dateSO = ye + "-" + mo + "-01"

	if _, err := f.Raw("SELECT Date, EXTRACT(DAY FROM Date) day, coalesce(sum(total_price), 0) total_price"+
		" from ("+
		" select last_day(?) - INTERVAL (a.a + (10 * b.a) + (100 * c.a)) DAY as date"+
		" from (select 0 as a union all select 1 union all select 2 union all select 3 union all select 4 union all select 5 union all select 6 union all select 7 union all select 8 union all select 9) as a"+
		" cross join (select 0 as a union all select 1 union all select 2 union all select 3 union all select 4 union all select 5 union all select 6 union all select 7 union all select 8 union all select 9) as b"+
		" cross join (select 0 as a union all select 1 union all select 2 union all select 3 union all select 4 union all select 5 union all select 6 union all select 7 union all select 8 union all select 9) as c"+
		" ) a left join (select * from sales_order where status not in(3,4)) so on a.date = so.recognition_date"+
		" where a.Date between ?"+
		" and last_day(?"+
		") group by a.Date", dateSO, dateSO, dateSO).QueryRows(&o); err != nil {
		return nil, total, err
	}

	return o, total, nil

}
