package repository

import (
	"git.edenfarm.id/project-version2/datamodel/model"

	"git.edenfarm.id/cuxs/orm"
)

// GetPriceSet find a single data price set using field and value condition.
func GetPriceSet(field string, values ...interface{}) (*model.PriceSet, error) {
	m := new(model.PriceSet)
	o := orm.NewOrm()
	o.Using("read_only")

	if err := o.QueryTable(m).Filter(field, values...).RelatedSel().Limit(1).One(m); err != nil {
		return nil, err
	}
	return m, nil
}

// GetPriceSets : function to get data from database based on parameters
func GetPriceSets(rq *orm.RequestQuery) (m []*model.PriceSet, total int64, err error) {
	q, _ := rq.QueryReadOnly(new(model.PriceSet))

	if total, err = q.Exclude("status", 3).Count(); err != nil || total == 0 {
		return nil, total, err
	}

	var mx []*model.PriceSet
	if _, err = q.Exclude("status", 3).All(&mx, rq.Fields...); err == nil {
		return mx, total, nil
	}

	return nil, total, err
}

// GetFilterPriceSets : function to get data from database based on parameters with filtered permission
func GetFilterPriceSets(rq *orm.RequestQuery) (m []*model.PriceSet, total int64, err error) {
	q, _ := rq.QueryReadOnly(new(model.PriceSet))

	if total, err = q.Exclude("status", 3).Count(); err != nil || total == 0 {
		return nil, total, err
	}

	var mx []*model.PriceSet
	if _, err = q.Exclude("status", 3).All(&mx, rq.Fields...); err == nil {
		return mx, total, nil
	}

	return nil, total, err
}

// ValidPriceSet : function to check if id is valid in database
func ValidPriceSet(id int64) (priceSet *model.PriceSet, e error) {
	priceSet = &model.PriceSet{ID: id}
	e = priceSet.Read("ID")

	return
}

// CheckPriceSetData : function to check data based on filter and exclude parameters
func CheckPriceSetData(filter, exclude map[string]interface{}) (ps []*model.PriceSet, total int64, err error) {
	rq := orm.RequestQuery{}
	o, _ := rq.QueryReadOnly(new(model.PriceSet))

	for k, v := range filter {
		o = o.Filter(k, v)
	}

	for k, v := range exclude {
		o = o.Exclude(k, v)
	}
	if total, err := o.All(&ps); err == nil {
		return ps, total, nil
	}

	return nil, 0, err
}

// GetAllPriceSets : function to get all price set data
func GetAllPriceSets() (m []*model.PriceSet, err error) {
	w := new(*model.PriceSet)
	o := orm.NewOrm()
	o.Using("read_only")

	if _, err := o.QueryTable(w).Exclude("status", 3).All(&m); err != nil {
		return nil, err
	}

	return m, nil
}

// GetListPriceSetAgent : function to get list price set of agent
func GetListPriceSetAgent(agentID int64) (m []*model.MerchantPriceSet, total int64, err error) {
	o := orm.NewOrm()
	o.Using("read_only")

	o1 := o.QueryTable(new(model.MerchantPriceSet))

	if total, err = o1.Filter("merchant_id", agentID).All(&m); err == nil {
		return m, total, err
	}

	return nil, 0, err
}
