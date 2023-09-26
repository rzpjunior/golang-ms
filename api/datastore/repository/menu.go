package repository

import (
	"git.edenfarm.id/project-version2/datamodel/model"

	"git.edenfarm.id/cuxs/orm"
)

// GetMenu find a single data price set using field and value condition.
func GetMenu(field string, values ...interface{}) (*model.Menu, error) {
	m := new(model.Menu)
	o := orm.NewOrm()
	o.Using("read_only")

	if err := o.QueryTable(m).Filter(field, values...).RelatedSel().Limit(1).One(m); err != nil {
		return nil, err
	}

	return m, nil
}

// GetMenus : function to get menu data from database based on parameters
func GetMenus(rq *orm.RequestQuery) (m []*model.Menu, total int64, err error) {
	q, _ := rq.QueryReadOnly(new(model.Menu))

	var mx []*model.Menu
	if total, err = q.RelatedSel(1).Filter("status", 1).Filter("parent_id__isnull", true).OrderBy("order").All(&mx, rq.Fields...); err == nil {
		for _, v := range mx {
			_, err = q.RelatedSel(1).Filter("status", 1).Filter("parent_id", v.ID).OrderBy("order").All(&v.Child)
		}
		return mx, total, nil
	}

	return nil, 0, err
}

// GetFilterMenus : function to get data from database based on parameters with filtered permission
func GetFilterMenus(rq *orm.RequestQuery) (m []*model.Menu, total int64, err error) {
	q, _ := rq.QueryReadOnly(new(model.Menu))

	var mx []*model.Menu
	if total, err = q.Filter("status", 1).All(&mx, rq.Fields...); err == nil {
		return mx, total, nil
	}

	return nil, 0, err
}

// ValidMenu : function to check if id is valid in database
func ValidMenu(id int64) (Menu *model.Menu, e error) {
	Menu = &model.Menu{ID: id}
	e = Menu.Read("ID")

	return
}

// GetMenusByUserID : function to get menu list based on user id
func GetMenusByUserID(userIds ...interface{}) (m []*model.Menu, err error) {
	o := orm.NewOrm()
	o.Using("read_only")

	where := "where m.status = 1 and p.status = 1 and m.parent_id is null "
	q := "select m.* " +
		"from menu m " +
		"join permission p on m.permission_id = p.id "

	if len(userIds) > 0 {
		join := "join user_permission up on p.id = up.permission_id "
		where = where + "and up.user_id = ? "
		q = q + join
	}

	q = q + where + "order by m.order"
	if _, err = o.Raw(q, userIds).QueryRows(&m); err == nil {
		for _, v := range m {
			v.Privilege.Read("ID")
			_, err = o.QueryTable(new(model.Menu)).RelatedSel(1).Filter("status", 1).Filter("parent_id", v.ID).OrderBy("order").All(&v.Child)
		}
		return m, nil
	}

	return nil, err
}
