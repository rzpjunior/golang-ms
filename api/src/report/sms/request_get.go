// Copyright 2020 PT. Eden Pangan Indonesia. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

package sms

import (
	"reflect"
	"strings"

	"git.edenfarm.id/cuxs/orm"
)

// getPurchaseOrder : query to get data of report purchase order
func getPurchaseOrder(cond map[string]interface{}) (po []*reportPurchaseOrder, e error) {
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

		if v != "" {
			values = append(values, v)
		}
	}
	where = strings.TrimSuffix(where, " and")

	q := "SELECT s.code AS supplier_code, s.name AS supplier_name, w.name AS warehouse_name, po.code AS order_code, DATE_FORMAT(po.recognition_date, '%d-%m-%Y') AS order_date, " +
		"DATE_FORMAT(po.eta_date, '%d-%m-%Y') AS eta_date, po.eta_time AS eta_time, g.value_name AS order_status, pur.name AS order_payment_term, po.delivery_fee AS delivery_fee, " +
		"po.tax_amount AS tax_amount, po.total_charge AS grand_total, po.note AS order_note, po.total_price AS total_price, sub.name AS supplier_badge, " +
		"if(po.has_finished_gr = 1, 'Received', 'Not Received') good_receipt " +
		"FROM purchase_order po " +
		"JOIN supplier s ON s.id = po.supplier_id " +
		"JOIN warehouse w ON w.id = po.warehouse_id " +
		"JOIN glossary g ON g.value_int = po.status AND g.attribute = 'doc_status' " +
		"JOIN term_payment_pur pur ON pur.id = po.term_payment_pur_id " +
		"JOIN supplier_badge sub ON po.supplier_badge_id = sub.id " +
		"WHERE po.status != 4 and " + where

	_, e = o.Raw(q, values).QueryRows(&po)

	return
}

// getPurchaseOrderItem : query to get data of report purchase order item
func getPurchaseOrderItem(cond map[string]interface{}) (poi []*reportPurchaseOrderItem, e error) {
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

	q := "SELECT po.code AS order_code, p.code AS product_code, p.name AS product_name, u.name AS uom, poi.note AS order_item_note, poi.order_qty AS ordered_qty, " +
		"poi.unit_price AS order_unit_price, poi.subtotal AS subtotal, poi.weight AS total_weight, a.name AS area_name, w.name AS warehouse_name, s.code AS supplier_code, " +
		"s.name AS supplier_name, DATE_FORMAT(po.recognition_date, '%d-%m-%Y') AS order_date, DATE_FORMAT(po.eta_date, '%d-%m-%Y') AS eta_date, pii.invoice_qty AS invoiced_qty, poi.purchase_qty, " +
		"IF(poi.include_tax = 1, 'Yes', 'No') include_tax_str, poi.tax_percentage order_tax_percentage, poi.tax_amount order_tax_amount, poi.unit_price_tax order_unit_price_tax, " +
		"IF(poi.taxable_item = 1, 'Yes', 'No') taxable_item_str " +
		"FROM purchase_order po " +
		"JOIN purchase_order_item poi ON poi.purchase_order_id = po.id " +
		"LEFT JOIN purchase_invoice_item pii ON pii.purchase_order_item_id = poi.id " +
		"JOIN product p ON p.id = poi.product_id " +
		"JOIN uom u ON u.id = p.uom_id " +
		"JOIN warehouse w ON w.id = po.warehouse_id " +
		"JOIN area a ON a.id = w.area_id " +
		"JOIN supplier s ON s.id = po.supplier_id " +
		"WHERE po.status != 4 and " + where

	_, e = o.Raw(q, values).QueryRows(&poi)

	return
}

// getPurchaseInvoice : query to get data of report purchase invoice
func getPurchaseInvoice(cond map[string]interface{}) (pi []*reportPurchaseInvoice, e error) {
	var (
		where  string
		values []interface{}
	)

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

	q := "SELECT pi.id invoice_id, a.name AS area_name, w.name AS warehouse_name, po.code AS order_code, DATE_FORMAT(po.recognition_date, '%d-%m-%Y') AS order_date, " +
		"DATE_FORMAT(po.eta_date, '%d-%m-%Y') AS eta_date, s.code AS supplier_code, s.name AS supplier_name, po.total_charge AS total_order, pi.code AS invoice_code, " +
		"DATE_FORMAT(pi.recognition_date, '%d-%m-%Y') AS invoice_date, DATE_FORMAT(pi.due_date, '%d-%m-%Y') AS invoice_due_date, UPPER(g.value_name) AS invoice_status, " +
		"pi.note AS invoice_note, pi.delivery_fee AS delivery_fee, pi.total_price AS invoice_amount, pi.total_charge AS total_invoice, pi.tax_amount AS tax_amount, " +
		"case when pi.adjustment = 2 then pi.adj_amount * -1 else pi.adj_amount end AS adjustment_amount, pi.adj_note AS adjustment_note, " +
		"pi.created_at, st.name created_by, sty.name supplier_type, tpp.name payment_term, gr.ata_date ata_date, pp.total_payment " +
		"FROM purchase_order po " +
		"LEFT JOIN purchase_invoice pi ON po.id = pi.purchase_order_id " +
		"JOIN warehouse w ON w.id = po.warehouse_id " +
		"JOIN area a ON a.id = w.area_id " +
		"JOIN supplier s ON s.id = po.supplier_id " +
		"JOIN supplier_type sty ON s.supplier_type_id = sty.id " +
		"LEFT JOIN glossary g ON g.value_int = pi.status AND g.attribute = 'doc_status' " +
		"LEFT JOIN staff st ON pi.created_by = st.id " +
		"JOIN term_payment_pur tpp ON po.term_payment_pur_id = tpp.id " +
		"LEFT JOIN goods_receipt gr ON po.id = gr.purchase_order_id " +
		"LEFT JOIN " +
		"(" +
		"SELECT purchase_invoice_id, SUM(amount) total_payment FROM purchase_payment WHERE status = 2 GROUP BY purchase_invoice_id" +
		") pp ON pi.id = pp.purchase_invoice_id " +
		"WHERE (pi.id is null OR pi.status != 4) AND " + where + " " +
		"ORDER BY po.id DESC, pi.id DESC"
	_, e = o.Raw(q, values).QueryRows(&pi)

	return
}

// getPurchasePayment : query to get data of report purchase payment
func getPurchasePayment(cond map[string]interface{}) (pp []*reportPurchasePayment, e error) {
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

	q := "SELECT ar.name AS area_name, s.code AS supplier_code, s.name AS supplier_name, pp.code AS payment_code, DATE_FORMAT(pp.recognition_date, '%d-%m-%Y') AS payment_date, " +
		"pp.amount AS payment_amount, pm.name AS payment_method, UPPER(g.value_name) AS payment_status, pi.code AS invoice_code, pi.total_charge AS total_invoice, " +
		"pp.bank_payment_voucher_number AS payment_number, pp.created_at, st.name created_by " +
		"FROM purchase_payment pp " +
		"JOIN purchase_invoice pi ON pi.id = pp.purchase_invoice_id " +
		"JOIN purchase_order po ON po.id = pi.purchase_order_id " +
		"JOIN warehouse w ON w.id = po.warehouse_id " +
		"JOIN area ar ON ar.id = w.area_id " +
		"JOIN supplier s ON s.id = po.supplier_id " +
		"JOIN payment_method pm ON pm.id = pp.payment_method_id " +
		"JOIN glossary g ON g.value_int = pp.status AND g.attribute = 'doc_status' " +
		"LEFT JOIN staff st ON pp.created_by = st.id " +
		"WHERE pi.status != 4 and " + where

	_, e = o.Raw(q, values).QueryRows(&pp)

	return
}

// getPurchaseInvoiceItem : query to get data of report purchase invoice item
func getPurchaseInvoiceItem(cond map[string]interface{}) (pii []*reportPurchaseInvoiceItem, e error) {
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

	q := "SELECT s.name AS supplier_name, w.name AS warehouse_name, a.name AS area, po.code AS order_code, g.value_name AS order_status, pi.code AS invoice_code, " +
		"g3.value_name AS invoice_status, gr.code AS gr_code, g2.value_name AS gr_status, po.eta_date AS eta_date, p.code AS product_code, p.name AS product_name, " +
		"u.name AS uom, poi.unit_price AS unit_price, poi.order_qty AS order_qty, gri.deliver_qty AS delivered_qty, gri.receive_qty AS received_qty, pii.invoice_qty AS invoice_qty, " +
		"gri.reject_qty AS reject_qty, pi.delivery_fee AS delivery_fee, pii.subtotal AS total_invoice, " +
		"IF(pii.include_tax = 1, 'Yes', 'No') include_tax_str, pii.tax_percentage tax_percentage, pii.tax_amount tax_amount, pii.unit_price_tax unit_price_tax, " +
		"IF(pii.taxable_item = 1, 'Yes', 'No') taxable_item_str " +
		"FROM purchase_order_item poi " +
		"JOIN purchase_order po ON po.id = poi.purchase_order_id AND po.status NOT IN (3 , 4, 5) " +
		"LEFT JOIN purchase_invoice pi ON pi.purchase_order_id = po.id AND pi.status NOT IN (3 , 4) " +
		"LEFT JOIN purchase_invoice_item pii ON pii.purchase_order_item_id = poi.id AND pi.id = pii.purchase_invoice_id " +
		"LEFT JOIN goods_receipt gr ON gr.purchase_order_id = po.id AND gr.status = 2 " +
		"LEFT JOIN goods_receipt_item gri ON gri.purchase_order_item_id = poi.id AND gr.id = gri.goods_receipt_id " +
		"JOIN supplier s ON s.id = po.supplier_id " +
		"JOIN warehouse w ON w.id = po.warehouse_id " +
		"JOIN area a ON a.id = w.area_id " +
		"JOIN product p ON p.id = poi.product_id " +
		"JOIN uom u ON u.id = p.uom_id " +
		"LEFT JOIN glossary g ON g.value_int = po.status AND g.attribute = 'doc_status' " +
		"LEFT JOIN glossary g2 ON g2.value_int = gr.status AND g2.attribute = 'doc_status' " +
		"LEFT JOIN glossary g3 ON g3.value_int = pi.status AND g3.attribute = 'doc_status'" +
		"WHERE " + where

	_, e = o.Raw(q, values).QueryRows(&pii)

	return
}

// getCogs : query to get data of report cogs
func getCogs(cond map[string]interface{}) (c []*reportCogs, e error) {
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

	q := "SELECT a.name area_name, w.name warehouse_name, c.eta_date eta_date, p.code product_code, p.name product_name, u.name uom, c.total_avg avg_price " +
		"FROM cogs c " +
		"JOIN product p ON c.product_id = p.id " +
		"JOIN uom u ON p.uom_id = u.id " +
		"JOIN warehouse w ON c.warehouse_id = w.id " +
		"JOIN area a ON w.area_id = a.id " +
		"WHERE " + where

	_, e = o.Raw(q, values).QueryRows(&c)

	return
}

// getPriceComparison : query to get data of report price comparison
func getPriceComparison(cond map[string]interface{}) (pc []*reportPriceComparison, e error) {
	var where string
	var values []interface{}

	o := orm.NewOrm()
	o.Using("scrape")

	for k, v := range cond {
		if reflect.TypeOf(v).Kind().String() == "slice" {
			where = where + " " + k + "? and ? and"
		} else {
			where = where + " " + k + "? and"
		}

		values = append(values, v)
	}
	where = strings.TrimSuffix(where, " and")

	q := "SELECT dpri.scraped_date survey_date, da.name area_name, dpro.code product_code, dpro.name product_name, dpro.uom uom, dpri.price selling_price, COALESCE(ppri1.price_after_discount, 0) public_price_1, COALESCE(ppri2.price_after_discount, 0) public_price_2 " +
		"FROM dashboard_product dpro " +
		"JOIN dashboard_price dpri ON dpro.id = dpri.product_id " +
		"JOIN dashboard_area da ON dpri.area_id = da.id " +
		"JOIN matched_product mp ON dpro.id = mp.dashboard_product_id " +
		"JOIN matched_area ma ON dpri.area_id = ma.dashboard_area_id " +
		"LEFT JOIN public_product_1 ppro1 ON mp.public_product_1_id = ppro1.id " +
		"LEFT JOIN public_price_1 ppri1 ON ppro1.id = ppri1.product_id AND dpri.scraped_date = ppri1.scraped_date AND ma.public_data_area_1_id = ppri1.area_id " +
		"LEFT JOIN public_product_2 ppro2 ON mp.public_product_2_id = ppro2.id " +
		"LEFT JOIN public_price_2 ppri2 ON ppro2.id = ppri2.product_id AND dpri.scraped_date = ppri2.scraped_date AND ma.public_data_area_2_id = ppri2.area_id " +
		"WHERE " + where

	_, e = o.Raw(q, values).QueryRows(&pc)

	return
}

// getInbound : query to get data of report inbound detail
func getInbound(cond map[string]interface{}) (in []*reportInboundDetail, e error) {
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

	q := "select w.name warehouse_name, s.name supplier_name, po.code po_code, po.eta_date, po.eta_time, po.committed_at, gr.ata_date, gr.ata_time, " +
		"if(sb.name = 'Central Market', sb.name, 'Supplier') source " +
		"from purchase_order po " +
		"join supplier s on po.supplier_id = s.id " +
		"left join goods_receipt gr on po.id = gr.purchase_order_id and gr.status = 2 " +
		"join warehouse w on po.warehouse_id = w.id " +
		"join supplier_badge sb on s.supplier_badge_id = sb.id " +
		"where po.status != 4 and " + where + " " +
		"order by if(sb.name = 'Central Market', 1, 0) asc, po.eta_date desc, po.id desc "
	_, e = o.Raw(q, values).QueryRows(&in)

	return
}

// getFieldPurchaser : query to get data of report field purchaser
func getFieldPurchaser(cond map[string]interface{}) (fpoi []*reportFieldPurchaser, e error) {
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

	q := "SELECT DATE_FORMAT(pp.created_at, '%d-%m-%Y') AS purchase_plan_date, pp.code AS purchase_plan_code, DATE_FORMAT(po.created_at, '%d-%m-%Y') AS purchase_order_date, po.code AS purchase_order_code, " +
		"so.name AS supplier_organization_name, sup.name AS supplier_name, gret.value_name AS returnable, grej.value_name AS rejectable, p.code AS product_code, p.name AS product_name, u.name AS uom, ppi.purchase_plan_qty AS purchase_plan_qty, " +
		"ppi.unit_price AS price_reference, poi.purchase_qty AS purchase_qty, poi.unit_price AS unit_price, poi.subtotal AS total_price, pt.name AS payment_term, " +
		"st.name AS field_purchaser_name, CONCAT('https://maps.google.com/?q=',po.latitude,',',po.longitude) AS order_location, cs.code AS consolidated_shipment_code, " +
		"w.name AS warehouse_name, cs.driver_name AS driver_name, cs.vehicle_number AS vehicle_number, cs.driver_phone_number AS driver_phone_number, " +
		"DATE_FORMAT(po.eta_date, '%d-%m-%Y') AS eta_date, po.eta_time AS eta_time, DATE_FORMAT(gr.ata_date, '%d-%m-%Y') AS ata_date, gr.ata_time AS ata_time, " +
		"DATE_FORMAT(gr.created_at, '%d-%m-%Y') AS inbound_date, gri.receive_qty AS receive_qty, g.value_name AS status " +
		"FROM purchase_order po " +
		"JOIN purchase_plan pp ON pp.id = po.purchase_plan_id " +
		"JOIN supplier_organization so ON so.id = pp.supplier_organization_id " +
		"JOIN supplier sup ON sup.id = po.supplier_id " +
		"LEFT JOIN glossary gret ON gret.table = 'supplier' AND gret.attribute = 'returnable' AND gret.value_int = sup.returnable " +
		"LEFT JOIN glossary grej ON grej.table = 'supplier' AND grej.attribute = 'rejectable' AND grej.value_int = sup.rejectable " +
		"JOIN warehouse w ON w.id = po.warehouse_id " +
		"JOIN purchase_order_item poi ON po.id = poi.purchase_order_id " +
		"JOIN product p ON p.id = poi.product_id " +
		"JOIN uom u ON u.id = p.uom_id " +
		"JOIN purchase_plan_item ppi ON ppi.id = poi.purchase_plan_item_id " +
		"JOIN term_payment_pur pt ON pt.id = po.term_payment_pur_id " +
		"JOIN staff st ON st.id = po.created_by " +
		"LEFT JOIN consolidated_shipment cs ON cs.id = po.consolidated_shipment_id " +
		"LEFT JOIN (goods_receipt_item gri JOIN goods_receipt gr ON gr.id = gri.goods_receipt_id AND gr.status IN (1, 2)) ON gri.purchase_order_item_id = poi.id " +
		"JOIN glossary g ON g.table = 'purchase_order' AND g.attribute = 'status' AND g.value_int = po.status " +
		"WHERE po.status IN (1, 2) and " + where

	if _, e = o.Raw(q, values).QueryRows(&fpoi); e != nil {
		return nil, e
	}

	return
}
