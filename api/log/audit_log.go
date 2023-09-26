package log

import (
	"time"

	"git.edenfarm.id/project-version2/datamodel/model"
	"github.com/labstack/echo/v4"
)

// this func for make audit log who create by user or from dashboard,
// user = get from session, reffID = get from id like sales_order id (when you create sales order)
// types = get from name of module (sales_order, purchase_order, etc..), function =  get from endpoint (create, update, delete , etc..)
// note = get from deletion_note, cancellation_note, if function create, update fill note blank or this note is optional.
func AuditLogByUser(staff *model.Staff, reffID int64, types, function string, note ...string) (e error) {

	al := &model.AuditLog{
		Staff:     staff,
		RefID:     reffID,
		Type:      types,
		Function:  function,
		Timestamp: time.Now(),
	}
	if note != nil {
		al.Note = note[0]
	}
	if e = al.Save(); e != nil {
		e = echo.ErrBadRequest
	}

	return e
}

func AuditLogByMerchant(userMerchant *model.UserMerchant, reffID int64, types, function string, note ...string) (e error) {

	al := &model.AuditLog{
		UserMerchant: userMerchant,
		RefID:        reffID,
		Type:         types,
		Function:     function,
		Timestamp:    time.Now(),
		Note:         note[0],
	}
	if note != nil {
		al.Note = note[0]
	}
	if e = al.Save(); e != nil {
		e = echo.ErrBadRequest
	}

	return e
}

func AuditLogByMerchantAndUser(userMerchant *model.UserMerchant, staff *model.Staff, reffID int64, types, function string, note ...string) (e error) {

	al := &model.AuditLog{
		UserMerchant: userMerchant,
		Staff:        staff,
		RefID:        reffID,
		Type:         types,
		Function:     function,
		Timestamp:    time.Now(),
		Note:         note[0],
	}
	if note != nil {
		al.Note = note[0]
	}
	if e = al.Save(); e != nil {
		e = echo.ErrBadRequest
	}

	return e
}

func AuditLogByUserReturnAuditLogRow(staff *model.Staff, reffID int64, types, function string, note ...string) (e error, al *model.AuditLog) {

	al = &model.AuditLog{
		Staff:     staff,
		RefID:     reffID,
		Type:      types,
		Function:  function,
		Timestamp: time.Now(),
	}
	if note != nil {
		al.Note = note[0]
	}
	if e = al.Save(); e != nil {
		e = echo.ErrBadRequest
	}

	return e, al
}
