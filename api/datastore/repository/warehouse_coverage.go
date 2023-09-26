package repository

import (
	"git.edenfarm.id/project-version2/datamodel/model"

	"git.edenfarm.id/cuxs/orm"
)

// GetWarehouseCoverage find a single data warehouse coverage using field and value condition.
func GetWarehouseCoverage(field string, values ...interface{}) (*model.WarehouseCoverage, error) {
	m := new(model.WarehouseCoverage)
	o := orm.NewOrm()
	o.Using("read_only")

	if err := o.QueryTable(m).Filter(field, values...).RelatedSel().Limit(1).One(m); err != nil {
		return nil, err
	}
	return m, nil
}

// GetWarehouseCoverages : function to get data from database based on parameters
func GetWarehouseCoverages(rq *orm.RequestQuery) (m []*model.WarehouseCoverage, total int64, err error) {
	q, _ := rq.QueryReadOnly(new(model.WarehouseCoverage))

	if total, err = q.Count(); err != nil || total == 0 {
		return nil, total, err
	}

	var mx []*model.WarehouseCoverage
	if _, err = q.All(&mx, rq.Fields...); err == nil {
		for _, v := range mx {
			warehouseType, err := GetGlossaryMultipleValue("table", "warehouse", "attribute", "warehouse_type", "value_int", v.Warehouse.WarehouseType)
			if err == nil {
				v.Warehouse.WarehouseTypeName = warehouseType.ValueName
			}
		}

		return mx, total, nil
	}

	return nil, total, err
}

// GetFilterWarehouseCoverages : function to get data from database based on parameters with filtered permission
func GetFilterWarehouseCoverages(rq *orm.RequestQuery) (m []*model.WarehouseCoverage, total int64, err error) {
	q, _ := rq.QueryReadOnly(new(model.WarehouseCoverage))

	if total, err = q.GroupBy("subdistrict__district__id").Count(); err != nil || total == 0 {
		return nil, total, err
	}

	var mx []*model.WarehouseCoverage
	if _, err = q.GroupBy("subdistrict__district__id").All(&mx, rq.Fields...); err == nil {
		return mx, total, nil
	}

	return nil, total, err
}

func ValidWarehouseCoverage(id int64) (wc *model.WarehouseCoverage, e error) {
	wc = &model.WarehouseCoverage{ID: id}
	e = wc.Read("ID")

	return
}

// CheckWarehouseCoverageData : function to check data based on filter and exclude parameters
func CheckWarehouseCoverageData(filter, exclude map[string]interface{}) (wc []*model.WarehouseCoverage, total int64, err error) {
	rq := orm.RequestQuery{}
	o, _ := rq.QueryReadOnly(new(model.WarehouseCoverage))

	for k, v := range filter {
		o = o.Filter(k, v)
	}

	for k, v := range exclude {
		o = o.Exclude(k, v)
	}

	if total, err := o.All(&wc); err == nil {
		return wc, total, nil
	}

	return nil, 0, err
}
