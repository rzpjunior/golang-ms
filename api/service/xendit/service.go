package xendit

import (
	"git.edenfarm.id/cuxs/orm"
	"git.edenfarm.id/project-version2/datamodel/model"
	"git.edenfarm.id/project-version2/api/util"
	"github.com/xendit/xendit-go"
	"github.com/xendit/xendit-go/virtualaccount"
)

func PermataXenditFixedVA(m *model.Merchant) (sd *xendit.VirtualAccount, e *xendit.Error) {
	xendit.Opt.SecretKey = util.XenditKey
	data := &virtualaccount.CreateFixedVAParams{
		ExternalID: "PERMATA_FVA-" + m.Code,
		BankCode:   "PERMATA",
		Name:       "Eden Farm",
	}
	if sd, e = virtualaccount.CreateFixedVA(data); e == nil {
		var pc *model.PaymentChannel
		var merchantAccNum *model.MerchantAccNum
		orSelect := orm.NewOrm()
		orSelect.Using("read_only")

		orSelect.Raw("SELECT * FROM payment_channel WHERE value = ?", sd.BankCode+"_FVA").QueryRow(&pc)
		orSelect.Raw("SELECT * FROM merchant_acc_num WHERE merchant_id = ? AND payment_channel_id = ?", m.ID, 7).QueryRow(&merchantAccNum)
		if pc != nil {
			if merchantAccNum != nil && merchantAccNum.AccountNumber == "" {
				merchantAccNum.AccountNumber = sd.AccountNumber
				merchantAccNum.Save("AccountNumber")
			} else {
				man := &model.MerchantAccNum{
					Merchant:       m,
					PaymentChannel: pc,
					AccountNumber:  sd.AccountNumber,
					AccountName:    sd.Name,
				}
				man.Save()
			}
		}

	}
	return sd, e
}

func BCAXenditFixedVA(m *model.Merchant) (sd *xendit.VirtualAccount, e *xendit.Error) {
	xendit.Opt.SecretKey = util.XenditKey
	data := &virtualaccount.CreateFixedVAParams{
		ExternalID: "BCA_FVA-" + m.Code,
		BankCode:   "BCA",
		Name:       "Eden Farm",
	}
	if sd, e = virtualaccount.CreateFixedVA(data); e == nil {
		var pc *model.PaymentChannel
		var merchantAccNum *model.MerchantAccNum
		orSelect := orm.NewOrm()
		orSelect.Using("read_only")

		orSelect.Raw("SELECT * FROM payment_channel WHERE value = ?", sd.BankCode+"_FVA").QueryRow(&pc)
		orSelect.Raw("SELECT * FROM merchant_acc_num WHERE merchant_id = ? AND payment_channel_id = ?", m.ID, 6).QueryRow(&merchantAccNum)
		if pc != nil {
			if merchantAccNum != nil && merchantAccNum.AccountNumber == "" {
				merchantAccNum.AccountNumber = sd.AccountNumber
				merchantAccNum.Save("AccountNumber")
			} else {
				man := &model.MerchantAccNum{
					Merchant:       m,
					PaymentChannel: pc,
					AccountNumber:  sd.AccountNumber,
					AccountName:    sd.Name,
				}
				man.Save()
			}
		}

	}

	return sd, e
}
