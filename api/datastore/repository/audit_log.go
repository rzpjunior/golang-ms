package repository

import (
	"encoding/json"
	"strconv"

	"git.edenfarm.id/cuxs/mongodb"
	"git.edenfarm.id/cuxs/orm"
	"git.edenfarm.id/project-version2/datamodel/model"
)

// GetAreas get all data user that matched with query request parameters.
// returning slices of User, total data without limit and error.
func GetLog(rq *orm.RequestQuery) (m []*model.AuditLog, total int64, err error) {
	md := mongodb.NewMongo()

	// make new orm query
	q, _ := rq.QueryReadOnly(new(model.AuditLog))

	// get total data
	if total, err = q.Count(); err != nil || total == 0 {
		return nil, total, err
	}
	// get data requested
	var mx []*model.AuditLog
	if _, err = q.All(&mx, rq.Fields...); err == nil {
		for _, v := range mx {
			var documentHistoryLog model.DocumentHistoryLog
			filter := map[string]interface{}{
				"audit_log_id": strconv.FormatInt(v.ID, 10),
			}
			ret, err := md.GetOneDataWithFilter("Document_History_Log", filter)
			if err != nil {
				md.DisconnectMongoClient()
			}
			json.Unmarshal(ret, &documentHistoryLog)

			v.ChangesLog = documentHistoryLog.ChangesLog
		}

		return mx, total, nil
	}

	md.DisconnectMongoClient()

	// return error some thing went wrong
	return nil, total, err
}
