package report

import (
	"reflect"
	"strings"
	"time"

	"git.edenfarm.id/cuxs/orm"
)

func getSalesOrder(cond map[string]interface{}) (sor []*reportSalesOrder, e error) {
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

	q := "SELECT so.code sales_order_code, m.code customer_code, GROUP_CONCAT(tc.name ORDER BY tc.name SEPARATOR ', ') customer_tag, " +
		"m.name customer_name, m.phone_number customer_phone_number, b.pic_name recipient_name, " +
		"b.phone_number recipient_phone_number, so.shipping_address shipping_address, w.name warehouse_name, bt.name business_type, at.name archetype_name, a.name area_name, ad.city_name city, " +
		"ots.name order_type_name, s.display_name salesperson, date_format(so.recognition_date, '%d-%m-%Y') order_date, date_format(so.delivery_date, '%d-%m-%Y') order_delivery_date, so.status, " +
		"g2.value_name order_status, so.note order_note, so.total_sku_disc_amount, so.total_charge grand_total, so.order_channel, g.value_name order_channel, " +
		"so.vou_redeem_code promo_code, so.delivery_fee delivery_fee, date_format(so.created_at, '%d-%m-%Y %H:%i:%s') created_at, s2.name created_by, date_format(so.last_updated_at, '%d-%m-%Y %H:%i:%s') last_updated_at, " +
		"s3.name last_updated_by, if(so.cancel_type is null or so.cancel_type = 0, '', if(so.cancel_type = 1, 'Unfulfilled', 'Regular')) cancel_type, sg.name sales_group, so.estimate_time_departure " +
		"FROM sales_order so " +
		"JOIN branch b ON b.id = so.branch_id " +
		"JOIN merchant m ON m.id = b.merchant_id " +
		"JOIN business_type bt on bt.id = m.business_type_id " +
		"JOIN warehouse w ON w.id = so.warehouse_id " +
		"LEFT JOIN archetype at ON at.id = so.archetype_id " +
		"JOIN area a ON a.id = b.area_id " +
		"JOIN adm_division ad ON ad.sub_district_id = so.sub_district_id " +
		"LEFT JOIN staff s ON s.id = so.salesperson_id " +
		"LEFT JOIN tag_customer tc ON FIND_IN_SET(tc.id, m.tag_customer) != 0 " +
		"LEFT JOIN sales_group sg ON sg.id = so.sales_group_id " +
		"JOIN staff s2 ON s2.id = so.created_by " +
		"LEFT JOIN staff s3 ON s3.id = so.last_updated_by " +
		"JOIN glossary g ON g.value_int = so.order_channel AND g.attribute = 'order_channel' " +
		"JOIN glossary g2 ON g2.value_int = so.status AND g2.attribute = 'doc_status' " +
		"JOIN order_type_sls ots ON so.order_type_sls_id = ots.id " +
		"where so.status != 4 and " + where + " " +
		"GROUP BY so.id " +
		"order by so.recognition_date desc "

	_, e = o.Raw(q, values).QueryRows(&sor)

	for _, v := range sor {
		etdTime := time.Unix(v.ETD, 0)
		v.ETDSt = etdTime.Format("2006-01-02 15:04:05")

		if v.ETDSt == "1970-01-01 07:00:00" {
			v.ETDSt = ""
		}
	}

	return
}

func getSalesOrderItem(cond map[string]interface{}) (soir []*reportSalesOrderItem, e error) {
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

	q := "select so.code sales_order_code, p.code product_code, p.name product_name, c.name category_name, " +
		"uom.name uom_name, soi.note order_item_note, soi.order_qty ordered_qty, " +
		"coalesce(sii.invoice_qty, sii.invoice_qty, 0) invoice_qty, soi.unit_price order_unit_price, " +
		"soi.discount_qty, soi.unit_price_discount, soi.sku_disc_amount, " +
		"case when sd.id is null then '-' else sd.name end sku_discount_name, " +
		"case when sdit.id is null then '0' else max(sdit.disc_amount) end disc_amount, " +
		"soi.shadow_price order_unit_shadow_price, soi.subtotal subtotal, soi.weight total_weight, " +
		"so.recognition_date order_date, so.delivery_date order_delivery_date, a.name area_name, " +
		"w.name warehouse_name, wrt.name wrt_name, g.value_name order_status, soi.tax_percentage order_tax_percentage, " +
		"tso.name sales_order_type_name, IF(soi.taxable_item = 1, 'Yes', 'No') taxable_item_str from sales_order_item soi " +
		"join sales_order so on so.id = soi.sales_order_id " +
		"join order_type_sls tso on tso.id = so.order_type_sls_id " +
		"left join sales_invoice si on si.sales_order_id = so.id and si.status not in (3, 4) " +
		"left join sales_invoice_item sii on sii.sales_order_item_id = soi.id and sii.sales_invoice_id = si.id " +
		"left join sku_discount_item sdi on soi.sku_discount_item_id = sdi.id " +
		"left join sku_discount sd on sdi.sku_discount_id = sd.id " +
		"left join sku_discount_item_tier sdit on sdi.id = sdit.sku_discount_item_id and soi.order_qty >= sdit.minimum_qty " +
		"join wrt on wrt.id = so.wrt_id join warehouse w on w.id = so.warehouse_id " +
		"join area a on a.id = so.area_id " +
		"join product p on p.id = soi.product_id " +
		"join uom on uom.id = p.uom_id " +
		"join category c on c.id = p.category_id " +
		"join glossary g on g.value_int = so.status and g.attribute = 'doc_status' " +
		"where so.status != 4 and " + where + " " +
		"group by soi.id " +
		"order by so.id;"

	_, e = o.Raw(q, values).QueryRows(&soir)

	return
}

func getSalesInvoice(cond map[string]interface{}) (soi []*reportSalesInvoice, e error) {
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

	q := "SELECT so.code order_code, m.code merchant_code, m.name merchant_name, b.code branch_code, b.name branch_name, " +
		"date_format(so.delivery_date, '%d-%m-%Y') order_delivery_date, si.code invoice_code, date_format(si.recognition_date, '%d-%m-%Y') invoice_date, date_format(si.due_date, '%d-%m-%Y') invoice_due_date, g.value_name invoice_status, " +
		"si.adj_note adjustment_note, sum(sp.amount) total_confirmed_payment, " +
		"si.total_price total_invoice, si.delivery_fee delivery_fee, si.vou_disc_amount voucher_amount,( " +
		"CASE WHEN (si.adjustment = 2 ) THEN -( si.adj_amount ) ELSE si.adj_amount END ) adjustment_amount, " +
		"si.total_charge total_charge, a.name area, w.name warehouse, g2.value_name customer_group, bt.name business_type, ar.name archetype, tps.name payment_term, " +
		"tis.name invoice_term, s.name created_by, DATE_FORMAT (si.created_at,'%d-%m-%Y %H:%i:%s') created_at, s2.name updated_by, DATE_FORMAT (si.last_updated_at,'%d-%m-%Y %H:%i:%s') updated_at, si.point_redeem_amount point_redeem_amount, " +
		"si.total_sku_disc_amount " +
		"FROM sales_order so " +
		"LEFT JOIN sales_invoice si ON si.sales_order_id = so.id " +
		"LEFT JOIN staff s ON s.id = si.created_by " +
		"LEFT JOIN staff s2 ON s2.id = si.last_updated_by " +
		"JOIN branch b ON b.id = so.branch_id " +
		"JOIN merchant m ON m.id = b.merchant_id " +
		"JOIN glossary g2 ON g2.value_int = m.customer_group AND g2.attribute = 'customer_group' AND g2.table = 'merchant' " +
		"JOIN business_type bt ON bt.id = m.business_type_id " +
		"JOIN sub_district sd ON sd.id = b.sub_district_id " +
		"JOIN area a ON a.id = b.area_id " +
		"JOIN archetype ar ON ar.id = so.archetype_id " +
		"LEFT JOIN glossary g ON g.value_int = si.status AND g.attribute = 'doc_status' " +
		"LEFT JOIN sales_payment sp ON si.id = sp.sales_invoice_id AND sp.status NOT IN ( 3, 4 ) " +
		"LEFT JOIN term_payment_sls tps ON tps.id = so.term_payment_sls_id " +
		"LEFT JOIN term_invoice_sls tis ON tis.id = si.term_invoice_sls_id " +
		"JOIN warehouse w ON w.id = so.warehouse_id   " +
		"WHERE so.status != 4 and " + where + " " +
		"GROUP BY so.id "

	_, e = o.Raw(q, values).QueryRows(&soi)

	return
}

func getSalesPayment(cond map[string]interface{}) (sp []*reportSalesPayment, e error) {
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

	q := "SELECT sp.code payment_code, sp.bank_receive_num, date_format(sp.recognition_date, '%d-%m-%Y') payment_date, " +
		"IF(sp.received_date,date_format(sp.received_date, '%d-%m-%Y'), NULL) received_date, " +
		"a.name area, g.value_name payment_status, pm.name payment_method, " +
		"sp.amount payment_amount, so.code order_code, si.code invoice_code, b.code outlet_code, " +
		"b.name outlet_name, s.name created_by, DATE_FORMAT (sp.created_at,'%d-%m-%Y %H:%i:%s') created_at, w.name warehouse_name " +
		"FROM sales_payment sp " +
		"JOIN payment_method pm ON sp.payment_method_id = pm.id " +
		"JOIN staff s ON s.id = sp.created_by " +
		"JOIN sales_invoice si ON si.id = sp.sales_invoice_id " +
		"JOIN sales_order so ON so.id = si.sales_order_id " +
		"JOIN branch b ON b.id = so.branch_id " +
		"JOIN area a ON a.id = b.area_id " +
		"JOIN warehouse w ON w.id = so.warehouse_id " +
		"LEFT JOIN glossary g ON g.value_int = sp.status " +
		"AND g.attribute = 'doc_status' " +
		"WHERE sp.status != 4 and " + where + " " +
		"GROUP BY sp.id "

	_, e = o.Raw(q, values).QueryRows(&sp)

	return
}

// getProspectiveCustomer : query to get data of report prospective customer
func getProspectiveCustomer(cond map[string]interface{}) (sp []*reportProspectiveCustomer, e error) {
	var where string
	var values []interface{}

	o := orm.NewOrm()
	o.Using("read_only")

	values = append(values, "%d-%m-%Y %H:%i:%s")
	values = append(values, "%d-%m-%Y %H:%i:%s")
	for k, v := range cond {
		if reflect.TypeOf(v).Kind().String() == "slice" {
			where = where + " " + k + "? and ? and"
		} else {
			where = where + " " + k + "? and"
		}
		values = append(values, v)
	}
	where = strings.TrimSuffix(where, " and")

	q := "SELECT " +
		"pc.code prospect_customer_code, " +
		"pc.name prospect_customer_name, " +
		"bt.name business_type, " +
		"a.name archetype, " +
		"pc.pic_name pic_name, " +
		"pc.phone_number phone_number, " +
		"pc.pic_finance_name,  " +
		"pc.pic_finance_contact,  " +
		"pc.pic_business_name,  " +
		"pc.pic_business_contact, " +
		"tp.name term_of_payment, " +
		"ti.name term_of_invoice, " +
		"pc.billing_address, " +
		"pc.note notes, " +
		"ad.area_name area, " +
		"pc.street_address, " +
		"ad.province_name province, " +
		"ad.city_name city, " +
		"ad.district_name district, " +
		"ad.sub_district_name sub_district, " +
		"ad.postal_code postal_code, " +
		"g.value_name best_time_to_call, " +
		"pc.referrer_code referral_code, " +
		"CONCAT(m.code,' ', m.name) existing_customer, " +
		"IF(pc.merchant_id IS NULL, 'No', 'Yes') request_upgrade, " +
		"DATE_FORMAT (pc.created_at,?) created_at, " +
		"DATE_FORMAT (pc.processed_at,?) processed_at, " +
		"s.name processed_by, " +
		"sp.name salesperson, " +
		"g3.value_name status, " +
		"g4.value_name decline_type, " +
		"pc.decline_note " +
		"FROM prospect_customer pc " +
		"JOIN archetype a ON a.id = pc.archetype_id " +
		"JOIN business_type bt ON bt.id = a.business_type_id " +
		"JOIN adm_division ad ON ad.sub_district_id = pc.sub_district_id " +
		"LEFT JOIN glossary g ON g.table = 'prospect_customer' AND g.attribute = 'time_consent' AND g.value_int = pc.time_consent " +
		"LEFT JOIN merchant m ON m.id = pc.merchant_id " +
		"LEFT JOIN staff s ON s.id = pc.processed_by " +
		"LEFT JOIN staff sp ON sp.id = pc.salesperson_id " +
		"JOIN glossary g3 ON g3.table = 'prospect_customer' AND g3.attribute = 'reg_status' AND g3.value_int = pc.reg_status " +
		"LEFT JOIN glossary g4 ON g4.table = 'prospect_customer' AND g4.attribute = 'decline_type' AND g4.value_int = pc.decline_type " +
		"LEFT JOIN term_invoice_sls ti ON ti.id=pc.term_invoice_sls_id " +
		"LEFT JOIN term_payment_sls tp ON tp.id=pc.term_payment_sls_id " +
		"WHERE " + where + " " +
		"ORDER BY pc.id;"

	_, e = o.Raw(q, values).QueryRows(&sp)

	return
}

// getSkuDiscount : query to get data of report sku discount
func getSkuDiscount(cond map[string]interface{}) (sd []*reportSkuDiscount, e error) {
	var where string
	var values []interface{}

	o := orm.NewOrm()
	o.Using("read_only")

	for k, v := range cond {
		if reflect.TypeOf(v).Kind().String() == "slice" {
			where += " " + k + "? and ? and"
		} else if strings.Contains(k, "FIND_IN_SET") {
			where += " " + k + " and"
		} else {
			where += " " + k + "? and"
		}

		values = append(values, v)
	}
	where = strings.TrimSuffix(where, " and")

	q := "SELECT sd.name sku_disc_name, " +
		"GROUP_CONCAT(DISTINCT ps.name ORDER BY ps.id SEPARATOR ', ') price_set_name, " +
		"GROUP_CONCAT(DISTINCT g.note ORDER BY g.value_int SEPARATOR ', ') order_channel_name, " +
		"sd.start_timestamp, sd.end_timestamp, d.name division_name, sd.note " +
		"FROM sku_discount sd  " +
		"LEFT JOIN price_set ps ON FIND_IN_SET(ps.id, sd.price_set)" +
		"JOIN division d ON sd.division_id = d.id  " +
		"JOIN glossary g ON g.`attribute` = 'order_channel'  " +
		"WHERE " + where + " " +
		"AND FIND_IN_SET(g.value_int, sd.order_channel)  " +
		"GROUP BY sd.id  " +
		"ORDER BY sd.start_timestamp DESC"

	_, e = o.Raw(q, values).QueryRows(&sd)

	return
}

// getSkuDiscountItem : query to get data of report sku discount item
func getSkuDiscountItem(cond map[string]interface{}) (sdi []*reportSkuDiscountItem, e error) {
	var where string
	var values []interface{}

	o := orm.NewOrm()
	o.Using("read_only")

	for k, v := range cond {
		if strings.Contains(k, "FIND_IN_SET") {
			where += " " + k + " and"
		} else {
			where += " " + k + "? and"
		}

		values = append(values, v)
	}
	where = strings.TrimSuffix(where, " and")

	if where != "" {
		where = " WHERE " + where
	}

	q := "SELECT sd.name sku_disc_name, p.code product_code, p.name product_name, u.name uom_name, sdit.tier_level, sdit.minimum_qty, sdit.disc_amount, " +
		"sdi.overall_quota, sdi.overall_quota_per_user, sdi.daily_quota_per_user, sdi.budget, sdi.rem_budget " +
		"FROM sku_discount sd " +
		"JOIN sku_discount_item sdi ON sd.id = sdi.sku_discount_id " +
		"JOIN sku_discount_item_tier sdit ON sdi.id = sdit.sku_discount_item_id " +
		"JOIN product p ON sdi.product_id = p.id " +
		"JOIN uom u ON p.uom_id = u.id " +
		" " + where + " " +
		"ORDER BY sd.start_timestamp DESC; "

	_, e = o.Raw(q, values).QueryRows(&sdi)

	return
}

func getSalesOrderFeedback() (sorf []*reportSalesOrderFeedback, e error) {
	o := orm.NewOrm()
	o.Using("read_only")

	q := "SELECT " +
		"m.code AS Merchant_Code, m.name AS Merchant_Name, m.phone_number AS Merchant_Phone_Number, bt.name AS Business_Type, a.name AS Archetype, b.shipping_address AS Branch_Shipping_Address" +
		", ad.city_name AS City, ad.district_name AS District, a2.name AS Area, sof.sales_order_code AS Sales_Order_Code, sof.delivery_date AS Delivery_Date" +
		", sof.created_at AS Feedback_Created_At, sof.rating_score AS Rating_Score, sof.tags AS Tags, sof.description AS Feedback_Description, g.value_name AS To_Be_Contacted " +
		"FROM sales_order_feedback sof JOIN sales_order so ON so.id = sof.sales_order_id " +
		"JOIN merchant m ON m.id = sof.merchant_id " +
		"JOIN business_type bt ON bt.id  = m.business_type_id " +
		"JOIN branch b ON b.id = so.branch_id " +
		"JOIN archetype a ON a.id = b.archetype_id " +
		"JOIN area a2 ON a2.id = b.area_id " +
		"JOIN adm_division ad ON ad.sub_district_id = b.sub_district_id " +
		"JOIN glossary g ON g.value_int = sof.to_be_contacted AND g.`attribute` = 'to_be_contacted';"

	_, e = o.Raw(q).QueryRows(&sorf)

	return
}

func getEdenPoint(cond map[string]interface{}) (sd []*reportEdenPointLog, e error) {
	var where string
	var values []interface{}

	o := orm.NewOrm()
	o.Using("read_only")

	for k, v := range cond {
		if reflect.TypeOf(v).Kind().String() == "slice" {
			where += " " + k + "? and ? and"
		} else if strings.Contains(k, "FIND_IN_SET") {
			where += " " + k + " and"
		} else {
			where += " " + k + "? and"
		}
		values = append(values, v)
	}
	where = strings.TrimSuffix(where, " and")

	if where != "" {
		where = " WHERE " + where
	}

	q := "SELECT mpl.created_date edenpoint_date, m.name merchant_name, " +
		"case when mpl.status = 1 then mpl.recent_point - mpl.point_value when mpl.status IN (2,4,6) then mpl.recent_point + mpl.point_value else 0 end previous_edenpoint, " +
		"mpl.point_value edenpoint, case when mpl.status = 1 then 'Earned' when mpl.status IN (2,4) then 'Redeemed' when mpl.status = 6 then 'Expired' else '-' end status, " +
		"mpl.recent_point current_edenpoint, g2.value_name transaction_type, m4.name advocate_merchant, m2.name referee_merchant, mpl.campaign_id, " +
		"mpl.campaign_name, mpl.campaign_multiplier, mpl.note log_note, so.code order_code, so.recognition_date order_date, so.created_at, so.finished_at, " +
		"so.total_charge total_sales_order, g.value_name order_status " +
		"FROM merchant_point_log mpl " +
		"JOIN merchant m on mpl.merchant_id = m.id " +
		"LEFT JOIN sales_order so on so.id = mpl.sales_order_id " +
		"LEFT JOIN glossary g on g.value_int = so.status AND g.`attribute` = 'doc_status' " +
		"JOIN glossary g2 on g2.value_int = mpl.transaction_type AND g2.`attribute` = 'transaction_type' " +
		"LEFT JOIN merchant m2 on m2.id = mpl.referee_id " +
		"LEFT JOIN merchant m3 on m3.id = mpl.referrer_id " +
		"LEFT JOIN merchant m4 on m3.referral_code = m4.referral_code " +
		" " + where + " " +
		"GROUP BY mpl.id " +
		"ORDER BY mpl.id desc, mpl.created_date desc;"

	_, e = o.Raw(q, values).QueryRows(&sd)

	return
}
