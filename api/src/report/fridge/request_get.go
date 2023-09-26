package fridge

import (
	"reflect"
	"strings"

	"git.edenfarm.id/cuxs/orm"
)

func getSoldProductFridge(cond map[string]interface{}) (sor []*reportSoldProductFridge, e error) {
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

	q := "select now() as created_at,bf.last_seen_at,m.id as merchant_id,m.name as merchant_name,b.id as branch_id, " +
		"b.name as branch_name,w.name as warehouse_name,p.name as product_name ,bi.total_weight , u.name as uom_name, " +
		"a.name as area_name " +
		"from box_fridge bf " +
		"join branch_fridge bf2 on bf.warehouse_id =bf2.warehouse_id " +
		"join branch b on bf2.branch_id=b.id " +
		"join merchant m on m.id = b.merchant_id " +
		"join warehouse w on w.id=bf.warehouse_id " +
		"join box_item bi on bi.box_id=bf.box_id " +
		"join product p on p.id = bi.product_id " +
		"join uom u on u.id=p.uom_id " +
		"join area a on b.area_id =a.id " +
		"where bf.status =2 and " + where + " " +
		" " +
		"order by bf.last_seen_at desc "

	_, e = o.Raw(q, values).QueryRows(&sor)

	return
}

func getAllProductFridge(cond map[string]interface{}) (sor []*reportAllProductFridge, e error) {
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

	if where != "" {
		where = " WHERE " + where
	}

	q := "select now() as created_at,bf.last_seen_at,m.id as merchant_id,m.name as merchant_name,b.id as branch_id, " +
		"b.name as branch_name,w.name as warehouse_name,p.name as product_name ,bi.total_weight , u.name as uom_name, " +
		"p.code as product_code,bi.unit_price,bi.total_price, " +
		"a.name as area_name, bi.status as box_item_status,bf.status as box_fridge_status,g.value_name as status,bf.image_url,bi.finished_at " +
		"from box_fridge bf " +
		"join (select * from branch_fridge bf2 group by bf2.warehouse_id) bf2 on bf.warehouse_id =bf2.warehouse_id " +
		"join branch b on bf2.branch_id=b.id " +
		"join merchant m on m.id = b.merchant_id " +
		"join warehouse w on w.id=bf.warehouse_id " +
		"join box_item bi on bi.id=bf.box_item_id " +
		"join product p on p.id = bi.product_id " +
		"join uom u on u.id=p.uom_id " +
		"join area a on b.area_id =a.id " +
		"join glossary g ON g.value_int = bf.status AND g.`attribute` = 'status' and g.`table` = 'box_fridge' " +
		" " + where + " " +
		" " +
		"order by bf.last_seen_at desc "

	_, e = o.Raw(q, values).QueryRows(&sor)

	return
}
