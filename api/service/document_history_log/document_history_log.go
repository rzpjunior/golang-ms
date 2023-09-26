package document_history_log

import (
	"encoding/json"
	"errors"
	"fmt"
	"strconv"
	"strings"

	"git.edenfarm.id/cuxs/mongodb"
	"git.edenfarm.id/project-version2/api/util"
	"git.edenfarm.id/project-version2/datamodel/model"
	"golang.org/x/text/cases"
	"golang.org/x/text/language"
)

type UserDocumentHistoryLog struct {
	OldStaff   *model.Staff
	NewStaff   *model.Staff
	OldUser    *model.User
	NewUser    *model.User
	Type       string
	RefID      int64
	AuditLogID int64
}

func makeUserDocumentHistoryLog(data *UserDocumentHistoryLog) (e error) {

	if e := data.OldStaff.Role.Read("ID"); e != nil {
		return errors.New(util.ErrorInvalidData("role"))
	}

	if data.OldStaff.Parent != nil {
		if e := data.OldStaff.Parent.Read("ID"); e != nil {
			return errors.New(util.ErrorInvalidData("parent"))
		}
	}

	if e := data.OldStaff.Area.Read("ID"); e != nil {
		return errors.New(util.ErrorInvalidData("area"))
	}

	if e := data.OldStaff.Warehouse.Read("ID"); e != nil {
		return errors.New(util.ErrorInvalidData("warehouse"))
	}

	oldStaffJson, _ := json.Marshal(data.OldStaff)
	oldStaffRaw := make(map[string]string)
	json.Unmarshal(oldStaffJson, &oldStaffRaw)

	newStaffJson, _ := json.Marshal(data.NewStaff)
	newStaffRaw := make(map[string]string)
	json.Unmarshal(newStaffJson, &newStaffRaw)

	oldUserJson, _ := json.Marshal(data.OldUser)
	oldUserRaw := make(map[string]string)
	json.Unmarshal(oldUserJson, &oldUserRaw)

	newUserJson, _ := json.Marshal(data.NewUser)
	newUserRaw := make(map[string]string)
	json.Unmarshal(newUserJson, &newUserRaw)

	var previousData, afterData []model.Data

	if data.NewStaff.Role != nil && data.OldStaff.Role.ID != data.NewStaff.Role.ID {
		previousData = append(previousData, model.Data{
			FieldName: "Role",
			Value:     data.OldStaff.Role.Name,
		})

		afterData = append(afterData, model.Data{
			FieldName: "Role",
			Value:     data.NewStaff.Role.Name,
		})
	}

	if data.OldStaff.Parent != nil && data.NewStaff.Parent != nil {
		if data.OldStaff.Parent.ID != data.NewStaff.Parent.ID {
			previousData = append(previousData, model.Data{
				FieldName: "Supervisor",
				Value:     data.OldStaff.Parent.Name,
			})

			afterData = append(afterData, model.Data{
				FieldName: "Supervisor",
				Value:     data.NewStaff.Parent.Name,
			})
		}
	}

	if data.OldStaff.Parent != nil && data.NewStaff.Parent == nil {
		previousData = append(previousData, model.Data{
			FieldName: "Supervisor",
			Value:     data.OldStaff.Parent.Name,
		})

		afterData = append(afterData, model.Data{
			FieldName: "Supervisor",
			Value:     "-",
		})
	}

	if data.OldStaff.Parent == nil && data.NewStaff.Parent != nil {
		previousData = append(previousData, model.Data{
			FieldName: "Supervisor",
			Value:     "-",
		})

		afterData = append(afterData, model.Data{
			FieldName: "Supervisor",
			Value:     data.NewStaff.Parent.Name,
		})
	}

	if data.NewStaff.Area != nil && data.OldStaff.Area.ID != data.NewStaff.Area.ID {
		previousData = append(previousData, model.Data{
			FieldName: "Area",
			Value:     data.OldStaff.Area.Name,
		})

		afterData = append(afterData, model.Data{
			FieldName: "Area",
			Value:     data.NewStaff.Area.Name,
		})
	}

	if data.NewStaff.Warehouse != nil && data.OldStaff.Warehouse.ID != data.NewStaff.Warehouse.ID {
		previousData = append(previousData, model.Data{
			FieldName: "Warehouse",
			Value:     data.OldStaff.Warehouse.Name,
		})

		afterData = append(afterData, model.Data{
			FieldName: "Warehouse",
			Value:     data.NewStaff.Warehouse.Name,
		})
	}

	userKeysToCompare := []string{"note"}
	staffKeysToCompare := []string{"name", "display_name", "phone_number", "role", "area", "warehouse"}

	for _, key := range userKeysToCompare {
		isNotEqual := oldUserRaw[key] != newUserRaw[key]

		if oldUserRaw[key] == "" {
			oldUserRaw[key] = "-"
		}

		if newUserRaw[key] == "" {
			newUserRaw[key] = "-"
		}

		if isNotEqual {
			previousData = append(previousData, model.Data{
				FieldName: cases.Title(language.Und).String(strings.ReplaceAll(key, "_", " ")),
				Value:     oldUserRaw[key],
			})

			afterData = append(afterData, model.Data{
				FieldName: cases.Title(language.Und).String(strings.ReplaceAll(key, "_", " ")),
				Value:     newUserRaw[key],
			})
		}
	}

	for _, key := range staffKeysToCompare {
		isNotEqual := oldStaffRaw[key] != newStaffRaw[key]
		isNewValueNotEmpty := newStaffRaw[key] != ""

		if newStaffRaw[key] == "" {
			newStaffRaw[key] = "-"
		}

		if newStaffRaw[key] == "" {
			newStaffRaw[key] = "-"
		}

		if isNotEqual && isNewValueNotEmpty {
			previousData = append(previousData, model.Data{
				FieldName: cases.Title(language.Und).String(strings.ReplaceAll(key, "_", " ")),
				Value:     oldStaffRaw[key],
			})

			afterData = append(afterData, model.Data{
				FieldName: cases.Title(language.Und).String(strings.ReplaceAll(key, "_", " ")),
				Value:     newStaffRaw[key],
			})
		}
	}

	changeLog := model.ChangesLog{
		PreviousData: previousData,
		AfterData:    afterData,
	}

	isChanged := len(previousData) > 0 && len(afterData) > 0

	if isChanged {
		md := mongodb.NewMongo()
		_, err := md.InsertOneData("Document_History_Log", &model.DocumentHistoryLog{
			AuditLogID: strconv.FormatInt(data.AuditLogID, 10),
			RefID:      strconv.FormatInt(data.RefID, 10),
			Type:       data.Type,
			ChangesLog: changeLog,
		})

		if err != nil {
			fmt.Println(err)
			md.DisconnectMongoClient()
			return e
		}

		md.DisconnectMongoClient()
	}

	return
}
