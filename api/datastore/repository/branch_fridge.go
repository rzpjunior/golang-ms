package repository

import (
	"strconv"
	"time"

	"git.edenfarm.id/project-version2/datamodel/model"

	"git.edenfarm.id/cuxs/dbredis"
	"git.edenfarm.id/cuxs/orm"
)

// GetBranchFridge find a single data Branch Fridge using field and value condition.
func GetBranchFridge(field string, values ...interface{}) (*model.BranchFridge, error) {
	m := new(model.BranchFridge)
	o := orm.NewOrm()
	o.Using("read_only")

	if err := o.QueryTable(m).Filter(field, values...).RelatedSel().Limit(1).One(m); err != nil {
		return nil, err
	}
	return m, nil
}

// GetBranchFridges : function to get data from database based on parameters
func GetBranchFridges(rq *orm.RequestQuery) (m []*model.BranchFridgeListQuery, total int64, err error) {
	q, _ := rq.QueryReadOnly(new(model.BranchFridge))
	o := orm.NewOrm()
	o.Using("read_only")

	var mx []*model.BranchFridge
	var data []*model.BranchFridgeListQuery
	var to_offline_timer int64

	expPeriod := time.Duration(2) * time.Hour
	to_offline_key := "config_timer_to_offline"

	if dbredis.Redis.CheckExistKey(to_offline_key) {
		dbredis.Redis.GetCache(to_offline_key, &to_offline_timer)
	} else {
		if e := o.Raw("SELECT value from config_app where attribute = 'fridge_offline_timer'").QueryRow(&to_offline_timer); e != nil {

		}
		dbredis.Redis.SetCache(to_offline_key, to_offline_timer, expPeriod)
	}
	if _, err = q.GroupBy("warehouse_id").Exclude("status", 3).All(&mx, rq.Fields...); err == nil {
		for _, branchFridge := range mx {
			var tempData model.BranchFridgeListQuery
			branchFridge.Branch.Read()
			branchFridge.Branch.Merchant.Read()
			branchFridge.Warehouse.Read()

			tempData.BranchName = branchFridge.Branch.Name
			tempData.CustomerName = branchFridge.Branch.Merchant.Name
			tempData.WarehouseName = branchFridge.Warehouse.Name

			var tempOnlineFridge []*model.BranchFridge
			var tempOfflineFridge []*model.BranchFridge

			if _, e := o.Raw("select id, branch_id, warehouse_id, mac_address, code, note, 1 as status, created_at, created_by, last_updated_at, last_updated_by, last_seen_at "+
				" from "+
				"branch_fridge bf "+
				"where bf.status=1 and bf.warehouse_id =? and bf.last_seen_at >= NOW()-INTERVAL ? SECOND", branchFridge.Warehouse.ID, to_offline_timer).QueryRows(&tempOnlineFridge); e != nil {
			}

			if _, e := o.Raw("select id, branch_id, warehouse_id, mac_address, code, note, 0 as status, created_at, created_by, last_updated_at, last_updated_by, last_seen_at "+
				" from "+
				"branch_fridge bf "+
				"where bf.status=1 and bf.warehouse_id =? and bf.last_seen_at <= NOW()-INTERVAL ? SECOND", branchFridge.Warehouse.ID, to_offline_timer).QueryRows(&tempOfflineFridge); e != nil {
			}

			tempData.AllFridge = append(tempData.AllFridge, tempOnlineFridge...)
			tempData.AllFridge = append(tempData.AllFridge, tempOfflineFridge...)
			onlineFridge := int64(len(tempOnlineFridge))
			offlineFridge := int64(len(tempOfflineFridge))
			allFridge := onlineFridge + offlineFridge
			tempData.Status = strconv.Itoa(int(onlineFridge)) + "/" + strconv.Itoa(int(allFridge))
			data = append(data, &tempData)
		}
		total = int64(len(data))
		return data, total, nil
	}
	return nil, total, err
}

// GetFilterBranchFridges : function to get data from database based on parameters with filtered permission
func GetFilterBranchFridges(rq *orm.RequestQuery) (m []*model.BranchFridge, total int64, err error) {
	q, _ := rq.QueryReadOnly(new(model.BranchFridge))

	if total, err = q.Exclude("status", 3).Count(); err != nil || total == 0 {
		return nil, total, err
	}

	var mx []*model.BranchFridge
	if _, err = q.Exclude("status", 3).All(&mx, rq.Fields...); err == nil {
		return mx, total, nil
	}

	return nil, total, err
}

// ValidBranchFridge : function to check if id is valid in database
func ValidBranchFridge(id int64) (Branch *model.BranchFridge, e error) {
	Branch = &model.BranchFridge{ID: id}
	e = Branch.Read("ID")

	return
}
