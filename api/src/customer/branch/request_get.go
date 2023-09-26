package branch

import (
	"reflect"
	"strings"

	"git.edenfarm.id/cuxs/orm"
)

// getBranchFilterBySalesperson : query to get data of template assignment task
func getBranchFilterBySalesperson(cond map[string]interface{}) (m []*templateBranchBySalesperson, e error) {
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

	q := "SELECT a.id area_id, a.name area_name, sg.id sales_group_id, sg.name sales_group_name, b.id branch_id, b.code branch_code, b.name branch_name, s.id salesperson_id," +
		" s.code staff_code, s.name staff_name" +
		" FROM branch b" +
		" JOIN staff s ON s.id = b.salesperson_id" +
		" LEFT JOIN sales_group sg ON sg.id = s.sales_group_id" +
		" LEFT JOIN area a ON a.id = s.area_id" +
		" WHERE b.status = 1 AND s.status = 1 AND " + where + ";"

	if _, e = o.Raw(q, values).QueryRows(&m); e != nil {
		return nil, e
	}

	return m, nil
}
