package prospect_customer

import (
	"time"

	"git.edenfarm.id/cuxs/orm"
	"git.edenfarm.id/project-version2/api/log"
	"git.edenfarm.id/project-version2/api/util"
	"git.edenfarm.id/project-version2/datamodel/model"
)

func Decline(r declineRequest) (u *model.ProspectCustomer, e error) {
	u = &model.ProspectCustomer{
		ID:             r.ID,
		RegStatus:      3,
		ProcessedAt:    time.Now(),
		ProcessedBy:    r.Session.Staff.ID,
		IDCardNumber:   "",
		IDCardImage:    "",
		SelfieImage:    "",
		TaxpayerNumber: "",
		TaxpayerImage:  "",
		DeclineTypeID:  r.DeclineType,
		DeclineNote:    r.DeclineNote,
	}

	if e = u.Save("RegStatus", "ProcessedAt", "ProcessedBy", "IDCardNumber", "IDCardImage", "SelfieImage", "TaxpayerNumber", "TaxpayerImage", "DeclineTypeID", "DeclineNote"); e == nil {
		if r.ProspectiveCustomer.Merchant != nil {
			m := &model.Merchant{
				ID:            r.ProspectiveCustomer.Merchant.ID,
				UpgradeStatus: 3,
			}

			if e = m.Save("UpgradeStatus"); e == nil {
				e = log.AuditLogByUser(r.Session.Staff, u.ID, "prospect_customer", "decline", "")
			}
		} else {
			e = log.AuditLogByUser(r.Session.Staff, u.ID, "prospect_customer", "decline", r.DeclineNote)
		}

		// notification FS Apps when Decline Prospective Customer
		if r.ProspectiveCustomer.SalespersonID != 0 {
			f := orm.NewOrm()
			f.Using("read_only")

			f.Raw("SELECT * FROM staff where id = ?", r.ProspectiveCustomer.SalespersonID).QueryRow(&r.Staff)
			r.Staff.User.Read("ID")

			messageNotif := &util.MessageNotification{}
			f.Raw("SELECT message, title FROM notification WHERE code= 'NOT0018'").QueryRow(&messageNotif)
			f.Raw("SELECT g.value_name FROM glossary g where g.table = 'prospect_customer' AND g.attribute = 'decline_type' AND g.value_int = ?", r.DeclineType).QueryRow(&r.DeclineTypeString)
			messageNotif.Message = util.ReplaceNotificationSalesOrder(messageNotif.Message, "#Name#", r.ProspectiveCustomer.Name)
			messageNotif.Message = util.ReplaceNotificationSalesOrder(messageNotif.Message, "#type#", r.DeclineTypeString)

			mn := &util.ModelNotification{
				SendTo:    r.Staff.User.SalesAppNotifToken,
				Title:     messageNotif.Title,
				Message:   messageNotif.Message,
				Type:      "2",
				RefID:     r.ID,
				ServerKey: util.FieldSalesServerKeyFireBase,
				StaffID:   r.Staff.ID,
			}
			util.PostModelNotificationFieldSales(mn)
		}

	}

	return u, e
}

func Save(r createRequest) (u *model.ProspectCustomer, e error) {
	r.CodeProspectCustomer, e = util.GenerateCode(r.CodeProspectCustomer, "prospect_customer", 6)

	o := orm.NewOrm()
	o.Begin()

	if e := r.ArcheType.BusinessType.Read("ID"); e != nil {
		o.Rollback()
		return nil, e
	}

	u = &model.ProspectCustomer{
		Code:             r.CodeProspectCustomer,
		Name:             r.Name,
		Archetype:        r.ArcheType,
		Email:            r.Email,
		PhoneNumber:      r.PhoneNumber,
		AltPhoneNumber:   r.AltPhoneNumber,
		StreetAddress:    r.StreetAddress,
		PicName:          r.PicName,
		TimeConsent:      r.TimeConsent,
		BusinessTypeName: r.ArcheType.BusinessType.Name,
		RegStatus:        1,
		RegChannel:       8,
		SubDistrict:      r.SubDistrict,
		ReferrerCode:     r.ReferralCode,
		ReferenceInfo:    r.ReferenceInfo,
		CreatedAt:        time.Now(),
	}
	if _, e = o.Insert(u); e == nil {
		e = log.AuditLogByUser(nil, u.ID, "prospect_customer", "create", "")
	} else {
		o.Rollback()
		return nil, e
	}

	o.Commit()

	return
}
