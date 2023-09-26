package migration

import (
	"strings"
	"time"

	"git.edenfarm.id/cuxs/orm"
	"git.edenfarm.id/project-version2/api/datastore/repository"
	"git.edenfarm.id/project-version2/datamodel/model"
)

type listPointMerchant struct {
	ID         int64   `orm:"column(id)" json:"id"`
	TotalPoint float64 `orm:"column(total_point)" json:"total_point"`
}

// migrationEdenPoint : function to get merchant and insert/update data to table merchant_point_expiration
func migrationEdenPoint(limit, offset, update string) (e error) {
	o := orm.NewOrm()
	orSelect := orm.NewOrm()
	orSelect.Using("read")
	var (
		getListPointMerchant    []*listPointMerchant
		configApp               *model.ConfigApp
		businesTypeEligibleEarn []string
		where                   string
	)

	// Get business type that eligible earn eden point
	if configApp, e = repository.GetConfigApp("attribute", "edenpoint_business_type"); e != nil {
		return e
	}

	businesTypeEligibleEarn = strings.Split(configApp.Value, ",")

	for range businesTypeEligibleEarn {
		where += " ?,"
	}

	where = strings.TrimSuffix(where, ",")

	q := "SELECT id, total_point FROM merchant WHERE business_type_id IN (" + where + ") AND total_point > 0 AND status IN (1, 2) LIMIT ? OFFSET ?"

	if _, e = orSelect.Raw(q, businesTypeEligibleEarn, limit, offset).QueryRows(&getListPointMerchant); e != nil {
		return e
	}
	currentPeriodDate, _ := time.Parse("2006-01-02", "2023-03-31")
	nextPeriodDate, _ := time.Parse("2006-01-02", "2023-06-30")
	for _, v := range getListPointMerchant {
		data := &model.MerchantPointExpiration{
			ID:                 v.ID,
			CurrentPeriodPoint: v.TotalPoint,
			CurrentPeriodDate:  currentPeriodDate,
			NextPeriodDate:     nextPeriodDate,
			LastUpdatedAt:      time.Now(),
		}

		if update == "1" {
			if e = data.Save(); e != nil {
				return e
			}
		} else {
			if _, e = o.Insert(data); e != nil {
				return e
			}
		}
	}

	return nil
}
