// Copyright 2020 PT. Eden Pangan Indonesia. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

package stock

import (
	"time"

	"git.edenfarm.id/cuxs/orm"
	"git.edenfarm.id/project-version2/api/datastore/repository"
	"git.edenfarm.id/project-version2/datamodel/model"
)

func UpdateCommitedStock(r updateCommitedRequest, date time.Time) (e error) {
	o := orm.NewOrm()
	o.Begin()

	// region expected stock
	qExpectedStock := "UPDATE `stock` s " +
		"JOIN(SELECT warehouse_id, product_id, SUM(qty) 'expected_qty' " +
		"FROM ( SELECT gt.destination_id 'warehouse_id', gti.product_id,gti.deliver_qty 'qty' " +
		"FROM goods_transfer gt " +
		"JOIN goods_transfer_item gti ON gti.goods_transfer_id = gt.id " +
		"WHERE gt.`status` = 5 UNION ALL SELECT po.warehouse_id, poi.product_id,poi.order_qty 'qty' " +
		"FROM purchase_order po JOIN purchase_order_item poi ON poi.purchase_order_id = po.id " +
		"WHERE po.`status` = 5 " +
		"UNION ALL SELECT pp.warehouse_id, pp.product_id ,(plan_qty - IF (order_qty is null , 0, order_qty)) 'qty' " +
		"FROM (SELECT pp.warehouse_id, ppi.product_id, SUM(ppi.purchase_plan_qty) 'plan_qty' " +
		"FROM purchase_plan pp " +
		"JOIN purchase_plan_item ppi ON ppi.purchase_plan_id = pp.id " +
		"WHERE pp.`status` = 1 GROUP BY pp.warehouse_id, ppi.product_id) pp " +
		"LEFT JOIN (SELECT po.warehouse_id, poi.product_id, SUM(poi.order_qty) 'order_qty' " +
		"FROM purchase_order po JOIN purchase_plan pp ON pp.id = po.purchase_plan_id and pp.status = 1 " +
		"JOIN purchase_order_item poi ON poi.purchase_order_id = po.id " +
		"WHERE po.`status` = 1 " +
		"GROUP BY po.warehouse_id, poi.product_id) po ON pp.warehouse_id = po.warehouse_id " +
		"AND pp.product_id = po.product_id) es GROUP BY 1,2) es2 " +
		"ON s.warehouse_id = es2.warehouse_id AND s.product_id = es2.product_id SET s.`expected_qty` = es2.expected_qty;"
	if _, e = o.Raw(qExpectedStock).Exec(); e != nil {
		o.Rollback()
	}

	// region intransit stock
	qIntransitStock := "UPDATE stock s join(" +
		"SELECT product_id, warehouse_id, SUM(qty) 'intransit_qty' " +
		"FROM( SELECT poi.product_id , po.warehouse_id, poi.order_qty 'qty' " +
		"FROM purchase_order po JOIN purchase_order_item poi ON poi.purchase_order_id = po.id " +
		"LEFT JOIN goods_receipt gr ON gr.purchase_order_id = po.id " +
		"WHERE po.`status` = 1 AND (gr.id IS NULL OR gr.`status` NOT IN (1,2)) " +
		"UNION ALL SELECT gti.product_id , gt.destination_id 'warehouse_id', gti.deliver_qty 'qty' " +
		"FROM goods_transfer gt JOIN goods_transfer_item gti ON gti.goods_transfer_id = gt.id " +
		"LEFT JOIN goods_receipt gr ON gr.goods_transfer_id = gt.id " +
		"WHERE gt.`status` = 1 AND gt.stock_type != 2 AND (gr.id IS NULL OR gr.`status` NOT IN (1,2))) it GROUP BY 1,2) x " +
		"on s.warehouse_id = x.warehouse_id and s.product_id = x.product_id set s.intransit_qty = x.intransit_qty"
	if _, e = o.Raw(qIntransitStock).Exec(); e != nil {
		o.Rollback()
	}

	// region received stock
	qReceived := "UPDATE stock s join(" +
		"SELECT gr.warehouse_id, gri.product_id, SUM(gri.receive_qty) 'received_qty' " +
		"FROM goods_receipt gr " +
		"JOIN goods_receipt_item gri ON gri.goods_receipt_id = gr.id " +
		"WHERE `gr`.`status` = 1 GROUP BY gr.warehouse_id, gri.product_id) x" +
		" on s.warehouse_id = x.warehouse_id and s.product_id = x.product_id set s.received_qty = x.received_qty"
	if _, e = o.Raw(qReceived).Exec(); e != nil {
		o.Rollback()
	}

	// region intransit waste stock
	qIntransitWasteStock := "UPDATE stock s join(" +
		"SELECT gt.destination_id 'warehouse_id', gti.product_id, SUM(gti.deliver_qty) 'intransit_Waste_qty' " +
		"FROM goods_transfer gt " +
		"JOIN goods_transfer_item gti ON gti.goods_transfer_id = gt.id " +
		"WHERE gt.stock_type = 2 and gt.`status` = 1 " +
		"GROUP BY gt.destination_id, gti.product_id " +
		") x " +
		"on s.warehouse_id = x.warehouse_id and s.product_id = x.product_id set s.intransit_Waste_qty = x.intransit_waste_qty"
	if _, e = o.Raw(qIntransitWasteStock).Exec(); e != nil {
		o.Rollback()
	}
	// endregion

	if _, e = o.QueryTable(new(model.Stock)).Update(orm.Params{
		"CommitedInStock":  0,
		"CommitedOutStock": 0,
	}); e == nil {
		var nextDate time.Time
		dayOff := new(model.DayOff)
		for nextDate = date.AddDate(0, 0, 1); ; nextDate = nextDate.AddDate(0, 0, 1) {
			dayOff, _ = repository.GetDayOff("off_date", nextDate)
			if dayOff == nil && int(nextDate.Weekday()) != 0 {
				break
			}
		}

		// update commited out
		if _, e = o.Raw("update stock s "+
			"join "+
			"( "+
			"select so.warehouse_id, soi.product_id, so.delivery_date, "+
			"sum(case when so.status = 1 then soi.order_qty else 0 end) qty_act, "+
			"sum(case when so.status = 9 then soi.order_qty else 0 end) qty_ind, "+
			"sum(case when so.status = 12 then soi.order_qty else 0 end) qty_pnd "+
			"from sales_order so "+
			"join sales_order_item soi on so.id = soi.sales_order_id "+
			"join order_type_sls ots on so.order_type_sls_id = ots.id "+
			"where (so.status in (9,12) or (so.status = 1 and so.payment_group_sls_id != 1)) and so.delivery_date = ? and ots.value != \"draft\" "+
			"group by so.warehouse_id, soi.product_id "+
			") so on s.warehouse_id = so.warehouse_id and s.product_id = so.product_id "+
			"set s.commited_out_stock = coalesce(so.qty_act, 0) + coalesce(so.qty_ind, 0) + coalesce(so.qty_pnd, 0)", nextDate.Format("2006-01-02")).Exec(); e == nil {
			// update commited in
			if _, e = o.Raw("update stock s "+
				"join  "+
				"( "+
				"select po.warehouse_id, poi.product_id, sum(poi.order_qty) order_qty "+
				"from purchase_order po "+
				"join purchase_order_item poi on po.id = poi.purchase_order_id "+
				"where (po.status = 1 or po.status = 5) and po.has_finished_gr = 2 and po.eta_date = ? "+
				"group by po.warehouse_id, poi.product_id "+
				") po on s.warehouse_id = po.warehouse_id and s.product_id = po.product_id "+
				"set s.commited_in_stock = coalesce(po.order_qty, 0)", date.Format("2006-01-02")).Exec(); e == nil {
				configApp := &model.ConfigApp{Attribute: "lst_stc_upd_cmt"}
				if e = configApp.Read("attribute"); e == nil {
					configApp.Value = time.Now().Format("2006-01-02 15:04:05")
					if _, e = o.Update(configApp, "Value"); e == nil {
						o.Commit()
						return nil
					}
				}
			}
		}
	}

	o.Rollback()
	return e
}
