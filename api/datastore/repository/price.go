package repository

import (
	"strings"

	"git.edenfarm.id/project-version2/datamodel/model"

	"git.edenfarm.id/cuxs/orm"
)

// GetPrice find a single data price set using field and value condition.
func GetPrice(field string, values ...interface{}) (*model.Price, error) {
	m := new(model.Price)
	o := orm.NewOrm()
	o.Using("read_only")

	if err := o.QueryTable(m).Filter(field, values...).RelatedSel().Limit(1).One(m); err != nil {
		return nil, err
	}
	return m, nil
}

// GetPrices : function to get data from database based on parameters
func GetPrices(rq *orm.RequestQuery, tagProduct string) (m []*model.Price, total int64, err error) {
	q, _ := rq.QueryReadOnly(new(model.Price))
	o := orm.NewOrm()
	o.Using("read_only")

	cond := q.GetCond()

	if tagProduct != "" {
		cond1 := orm.NewCondition()
		cond1 = cond1.And("product__tag_product__icontains", ","+tagProduct+",").Or("product__tag_product__istartswith", tagProduct+",").Or("product__tag_product__iendswith", ","+tagProduct).Or("product__tag_product", tagProduct)

		cond = cond.AndCond(cond1)
	}

	q = q.SetCond(cond)

	if total, err = q.Filter("product__status", 1).Filter("priceset__status", 1).Filter("product__category__status", 1).
		Count(); err != nil || total == 0 {
		return nil, total, err
	}

	var mx []*model.Price
	if _, err = q.RelatedSel(2).OrderBy("id").Filter("product__status", 1).Filter("priceset__status", 1).Filter("product__category__status", 1).
		All(&mx, rq.Fields...); err == nil {

		for _, v := range mx {
			if v.Product.TagProduct != "" {
				qMark := ""
				tagProductArr := strings.Split(v.Product.TagProduct, ",")
				for range tagProductArr {
					qMark += "?,"
				}
				qMark = strings.TrimSuffix(qMark, ",")

				if err := o.Raw("select group_concat(name) from tag_product where value in ("+qMark+") order by id asc", tagProductArr).QueryRow(&v.Product.TagProductStr); err != nil {
					return nil, total, err
				}
			}
		}
		return mx, total, nil
	}

	return nil, total, err
}
