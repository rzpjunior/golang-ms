package widget

import (
	"reflect"
	"strings"

	"git.edenfarm.id/cuxs/orm"
)

// getGrandTotalSalesOrderWithWRT : query to get grand total data of so with wrt
func getGrandTotalSalesOrderWithWRT(cond map[string]interface{}) (sowrt *getGrandTotalDashboardSalesOrderWithWRT, total int64, e error) {
	var where string
	var values []interface{}

	o := orm.NewOrm()
	o.Using("read_only")
	for k, v := range cond {
		if reflect.TypeOf(v).Kind().String() == "slice" {
			where = where + " " + k + " >= ? and " + k + " <= ? and"
		} else {
			where = where + " " + k + "? and"
		}

		values = append(values, v)
	}
	where = strings.TrimSuffix(where, " and")

	q := "select count(*) 'grand_total_so', sum(case when so.status in(1,9,12) then 1 else 0 end) 'grand_total_active'," +
		" sum(case when so.status in (3) then 1 else 0 end) 'grand_total_cancelled'," +
		" sum(case when so.status in (7,8,10,11,13) then 1 else 0 end) 'grand_total_on_delivery'," +
		" sum(case when so.status in (2) then 1 else 0 end) 'grand_total_finished'" +
		" from wrt w join sales_order so on so.wrt_id = w.id join warehouse wh on wh.id = so.warehouse_id" +
		" join area a on a.id = so.area_id" +
		" where " + where

	e = o.Raw(q, values).QueryRow(&sowrt)

	getSalesOrderWithWRT(cond, sowrt)

	return
}

// getSalesOrderWithWRT : query to get data of so with wrt
func getSalesOrderWithWRT(cond map[string]interface{}, sowrt *getGrandTotalDashboardSalesOrderWithWRT) {
	var where string
	var values []interface{}

	o := orm.NewOrm()
	o.Using("read_only")

	for k, v := range cond {
		if reflect.TypeOf(v).Kind().String() == "slice" {
			where = where + " " + k + " >= ? and " + k + " <= ? and"
		} else {
			where = where + " " + k + "? and"
		}

		values = append(values, v)
	}
	where = strings.TrimSuffix(where, " and")

	q := "select w.name 'wrt', count(*) 'total_so', sum(case when so.status in(1,9,12) then 1 else 0 end) 'active'," +
		" sum(case when so.status in (3) then 1 else 0 end) 'cancelled'," +
		" sum(case when so.status in (7,8,10,11,13) then 1 else 0 end) 'on_delivery'," +
		" sum(case when so.status in (2) then 1 else 0 end) 'finished'" +
		" from wrt w join sales_order so on so.wrt_id = w.id join warehouse wh on wh.id = so.warehouse_id" +
		" join area a on a.id = so.area_id" +
		" where " + where + " group by w.name"

	o.Raw(q, values).QueryRows(&sowrt.DashboardSOWRT)
	sowrt.TotalRow = len(sowrt.DashboardSOWRT)

}

// getTotalPickingOrderWithWRT : query to get data of so with wrt
func getTotalPickingOrderWithWRT(cond map[string]interface{}) (pickWrt *getDashboardTotalPickingOrderWithWRT, total int64, e error) {
	var where string
	var values []interface{}

	o := orm.NewOrm()
	o.Using("read_only")

	for k, v := range cond {
		if reflect.TypeOf(v).Kind().String() == "slice" {
			where = where + " " + k + " >= ? and " + k + " <= ? and"
		} else {
			where = where + " " + k + "? and"
		}

		values = append(values, v)
	}
	where = "and " + strings.TrimSuffix(where, " and")

	q := "select count(if(poa.status = 1,1,null)) 'new'," +
		" count(if(poa.status = 2, 1, null)) 'finished'," +
		" count(if(poa.status = 3, 1, null)) 'on_progress'," +
		" count(if(poa.status = 4, 1, null)) 'need_approval'," +
		" count(if(poa.status = 5, 1, null)) 'picked'," +
		" count(if(poa.status = 6, 1, null)) 'checking'" +
		" from picking_order po" +
		" join picking_order_assign poa on po.id = poa.picking_order_id" +
		" join sales_order so on so.id = poa.sales_order_id where so.status in (1,2,7,9,10,12,13)" + where

	e = o.Raw(q, values).QueryRow(&pickWrt)
	pickWrt.TotalSO = pickWrt.NewPickingStatus + pickWrt.FinishedPickingStatus +
		pickWrt.OnProgressPickingStatus + pickWrt.NeedApprovalPickingStatus +
		pickWrt.PickedPickingStatus + pickWrt.CheckingPickingStatus

	return

}

// getIdlePicking : query to get idle picker
func getIdlePicking(cond map[string]interface{}) (idlePickers []*getIdlePickingObj, total int64, e error) {
	var where string
	var values []interface{}

	o := orm.NewOrm()
	o.Using("read_only")
	for k, v := range cond {
		if reflect.TypeOf(v).Kind().String() == "slice" {
			where = where + " " + k + " >= ? and " + k + " <= ? and"
		} else {
			where = where + " " + k + "? and"
		}

		values = append(values, v)
	}
	where = strings.TrimSuffix(where, " and")
	if where != "" {
		where = "AND " + where
	}

	q := "SELECT CONCAT(s.code,' - ',s.name) 'staff' , timestampdiff(minute,MAX(poa.checkout_timestamp),now()) 'duration_minutes' ," +
		"CONCAT(FLOOR(timestampdiff(minute,MAX(poa.checkout_timestamp),now())/60),' Hour '," +
		"MOD(timestampdiff(minute,MAX(poa.checkout_timestamp),now()),60),' Minutes ') 'duration_idle' , w.name 'warehouse' " +
		"FROM picking_order_assign poa " +
		"JOIN staff s ON s.id = poa.staff_id " +
		"JOIN warehouse w ON w.id = s.warehouse_id " +
		"WHERE poa.checkout_timestamp is not null " + where + " " +
		" and poa.status NOT IN(3,1) " +
		"AND date(poa.checkout_timestamp) = CURRENT_DATE() GROUP BY poa.staff_id " +
		"ORDER BY duration_minutes DESC "

	var amount int64
	if amount, e = o.Raw(q, values).QueryRows(&idlePickers); e != nil {
		return nil, 0, e
	}

	return idlePickers, amount, e
}
