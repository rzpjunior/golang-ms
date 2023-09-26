package report

import (
	"reflect"
	"strconv"
	"strings"

	"git.edenfarm.id/cuxs/orm"
)

func getPackingOrderReport(deliveryDate string, warehouseID int64) (poi []*reportPackingOrder, e error) {

	o := orm.NewOrm()
	o.Using("read_only")
	warehouseIDStr := strconv.Itoa(int(warehouseID))

	q := "select po.delivery_date, p2.name product_name, u.name uom, poi.total_order, poia.subtotal_pack, poia.subtotal_weight, s2.code helper_code , s2.name helper_name " +
		"from packing_order po inner join packing_order_item poi on po.id = poi.packing_order_id " +
		"inner join packing_order_item_assign poia on poia.packing_order_item_id = poi.id " +
		"inner join product p2 on poi.product_id = p2.id inner join uom u on p2.uom_id = u.id " +
		"inner join staff s2 on s2.id = poia.staff_id " +
		"where po.status = 2 and poia.staff_id is not NULL and po.delivery_date = ? and po.warehouse_id = ?"

	_, e = o.Raw(q, deliveryDate, warehouseIDStr).QueryRows(&poi)
	if e != nil {
		return nil, e
	}
	return poi, nil
}

func getDeliveryOrderItem(cond map[string]interface{}) (m []*reportDeliveryOrderItem, e error) {
	var where string
	var values []interface{}
	o := orm.NewOrm()
	o.Using("read_only")

	for k, v := range cond {
		if reflect.TypeOf(v).Kind().String() == "slice" {
			where = where + " " + k + "? and ? and"
		} else {
			where = where + " " + k + "? and"
		}

		values = append(values, v)
	}
	where = strings.TrimSuffix(where, " and")

	// get data requested
	q := "SELECT so.code order_code,do.code delivery_code,p.code product_code ,p.name product_name," +
		"uom.name uom,doi.note delivery_item_note ,doi.deliver_qty delivered_qty,doi.receive_qty received_qty ," +
		"doi.weight delivery_weight,a.name area,w.name warehouse ,so.delivery_date order_delivery_date," +
		"do.recognition_date delivery_date,wrt.name wrt ,g.value_name delivery_status,do.note delivery_note " +
		"FROM sales_order_item soi " +
		"JOIN sales_order so ON so.id = soi.sales_order_id " +
		"LEFT JOIN delivery_order_item doi ON doi.sales_order_item_id = soi.id " +
		"LEFT JOIN delivery_order do ON do.id = doi.delivery_order_id " +
		"JOIN wrt ON wrt.id = so.wrt_id " +
		"JOIN warehouse w ON w.id = so.warehouse_id " +
		"JOIN area a ON a.id = so.area_id " +
		"JOIN product p ON p.id = soi.product_id " +
		"JOIN uom ON uom.id = p.uom_id " +
		"LEFT JOIN glossary g ON g.value_int = do.status AND g.table='delivery_order' AND g.attribute= 'status' " +
		"WHERE so.status != 4 and " + where

	_, e = o.Raw(q, values).QueryRows(&m)
	if e != nil {
		return nil, e
	}
	// return error some thing went wrong
	return m, nil

}

// getPricingInboundItem : query to get data of report pricing inbound item
func getPricingInboundItem(cond map[string]interface{}, cond2 map[string]interface{}) (pricingInboundItem []*reportPricingInboundItem, e error) {
	var where string
	var where2 string
	var values []interface{}
	var values2 []interface{}
	var reportPO []*reportPricingInboundItem

	o := orm.NewOrm()
	o.Using("read_only")

	for k, v := range cond {
		if reflect.TypeOf(v).Kind().String() == "slice" {
			where = where + " " + k + "? and ? and"
		} else {
			where = where + " " + k + "? and"
		}

		values = append(values, v)
	}

	for k, v := range cond2 {
		if reflect.TypeOf(v).Kind().String() == "slice" {
			where2 = where2 + " " + k + "? and ? and"
		} else {
			where2 = where2 + " " + k + "? and"
		}

		values2 = append(values2, v)
	}

	where = strings.TrimSuffix(where, " and")

	where2 = strings.TrimSuffix(where2, " and")

	q := "SELECT gr.id as id, gt.code as inbound_code, '-' as supplier_name,'-' as supplier_code, wo.name as warehouse_origin, " +
		"wd.name as warehouse_destination, a.name as area, gt.eta_date as eta_date, gr.ata_date as ata_date, p.code as product_code," +
		"p.name as product_name, u.name as uom, gti.request_qty as request_qty, gri.deliver_qty as delivered_qty, '-' as invoice_qty, gri.receive_qty as receive_qty, " +
		"gti.unit_cost as unit_price, g.value_name as inbound_status, wd.id warehouse_id, wd.area_id, gt.recognition_date order_date " +
		"FROM goods_receipt_item gri " +
		"LEFT JOIN goods_receipt gr ON gr.id = gri.goods_receipt_id " +
		"LEFT JOIN goods_transfer gt ON gt.id = gr.goods_transfer_id " +
		"LEFT JOIN goods_transfer_item gti ON gti.goods_transfer_id = gt.id AND gti.product_id = gri.product_id " +
		"LEFT JOIN warehouse wo ON wo.id = gt.origin_id " +
		"LEFT JOIN warehouse wd ON wd.id = gt.destination_id " +
		"LEFT JOIN area a ON a.id = wd.area_id " +
		"LEFT JOIN product p ON p.id = gri.product_id " +
		"LEFT JOIN uom u on u.id = p.uom_id " +
		"LEFT JOIN glossary g ON g.value_int = gr.status AND g.`attribute` = 'status' and g.`table` = 'goods_receipt' " +
		"WHERE gr.goods_transfer_id IS NOT NULL AND " + where + " GROUP BY gri.id"

	q2 := "SELECT po.code as inbound_code, s.name as supplier_name,s.code as supplier_code, '-' as warehouse_origin, " +
		"wd.name as warehouse_destination, a.name as area, po.eta_date as eta_date, gr.ata_date as ata_date, p.code as product_code, " +
		"p.name as product_name, u.name as uom, poi.order_qty as request_qty, gri.deliver_qty as delivered_qty, pii.invoice_qty as invoice_qty, gri.receive_qty as receive_qty, " +
		"IF(pii.id IS NOT NULL, IF(pii.unit_price_tax > 0, pii.unit_price_tax ,pii.unit_price), " +
		"IF(poi.unit_price_tax > 0, poi.unit_price_tax, poi.unit_price)) as unit_price, g.value_name as inbound_status, wd.id warehouse_id, wd.area_id area_id, po.recognition_date order_date, " +
		"IF(pii.id IS NOT NULL, pii.tax_percentage, poi.tax_percentage) tax_percentage, " +
		"p.taxable taxability  " +
		"FROM goods_receipt_item gri " +
		"LEFT JOIN goods_receipt gr ON gr.id = gri.goods_receipt_id " +
		"LEFT JOIN purchase_order po ON po.id = gr.purchase_order_id " +
		"LEFT JOIN supplier s ON s.id = po.supplier_id " +
		"LEFT JOIN warehouse wd ON wd.id = gr.warehouse_id " +
		"LEFT JOIN area a ON a.id = wd.area_id " +
		"LEFT JOIN product p ON p.id = gri.product_id " +
		"LEFT JOIN purchase_order_item poi ON poi.id = gri.purchase_order_item_id " +
		"LEFT JOIN purchase_invoice_item pii ON pii.purchase_order_item_id  = gri.purchase_order_item_id " +
		"LEFT JOIN uom u on u.id = p.uom_id " +
		"LEFT JOIN glossary g ON g.value_int = gr.status AND g.`attribute` = 'status' and g.`table` = 'goods_receipt' " +
		"WHERE gr.purchase_order_id IS NOT NULL AND " + where2 + " GROUP BY gri.id"

	if _, e = o.Raw(q, values).QueryRows(&pricingInboundItem); e != nil {
		return nil, e
	}

	if _, e = o.Raw(q2, values2).QueryRows(&reportPO); e != nil {
		return nil, e
	}

	pricingInboundItem = append(pricingInboundItem, reportPO...)

	return
}

func getPriceChangeHistoryReport(cond map[string]interface{}, countPriceSet int) (priceChangeHistory []*reportPriceChangeHistory, e error) {
	var where string
	var values []interface{}
	o := orm.NewOrm()
	o.Using("read_only")

	for k, v := range cond {
		if k == "p.price_set_id IN " {
			var counter int
			var condPriceSet string
			for counter < countPriceSet {
				condPriceSet = condPriceSet + "?, "
				counter++
			}
			condPriceSet = strings.TrimSuffix(condPriceSet, ", ")
			where = where + " " + k + "(" + condPriceSet + ")" + " and"
		} else {
			where = where + " " + k + "? and ? and"
		}

		values = append(values, v)
	}

	where = strings.TrimSuffix(where, " and")

	// get data requested
	q := "SELECT pl.created_at created_at, ps.name price_set,p2.name product_name, pl.unit_price unit_price," +
		"s.id staff_id, s.name created_by " +
		"FROM price_log pl " +
		"LEFT JOIN price p ON p.id = pl.price_id " +
		"LEFT JOIN product p2 ON p2.id = p.product_id " +
		"LEFT JOIN price_set ps ON ps.id = p.price_set_id " +
		"LEFT JOIN staff s ON s.id = pl.created_by " +
		"WHERE " + where + " GROUP BY pl.id"

	_, e = o.Raw(q, values).QueryRows(&priceChangeHistory)
	// return error some thing went wrong
	if e != nil {
		return nil, e
	}
	return priceChangeHistory, nil
}
