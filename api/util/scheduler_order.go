package util

import (
	"strconv"
	"strings"
	"time"

	"git.edenfarm.id/cuxs/orm"
	"github.com/robfig/cron/v3"
)

var Cron = cron.New()
var SecCron = cron.New()
var GetAreaPolicy []areaPolicy
var entryIdNotifyWillRemove, entryIdCancelledWillRemove []cron.EntryID

type areaPolicy struct {
	OrderTimeLimit string `orm:"column(order_time_limit);null" json:"order_time_limit"`
	Area           string `orm:"column(area_id)" json:"area_id"`
}

func CronJobPushNotification() {
	Cron.AddFunc("CRON_TZ=Asia/Jakarta 59 23 * * *", CronGetAreaPolicy) // setiap jam 23:59 waktu jakarta
	Cron.Start()
}

func CronGetAreaPolicy() {
	o := orm.NewOrm()
	o.Using("read_only")

	if len(GetAreaPolicy) > 0 {
		// mengosongkan value times
		GetAreaPolicy = GetAreaPolicy[:0]
	}

	o.Raw("SELECT order_time_limit,area_id  FROM area_policy").QueryRows(&GetAreaPolicy)
	// setiap jam 12 malam akan remove cron job untuk semua area
	for _, eid := range entryIdNotifyWillRemove {
		SecCron.Remove(eid)
	}
	for _, eidCancelled := range entryIdCancelledWillRemove {
		SecCron.Remove(eidCancelled)
	}
	entryIdNotifyWillRemove = entryIdNotifyWillRemove[:len(entryIdNotifyWillRemove)-len(entryIdNotifyWillRemove)]
	entryIdCancelledWillRemove = entryIdCancelledWillRemove[:len(entryIdCancelledWillRemove)-len(entryIdCancelledWillRemove)]
	createCronJobEveryArea()
}

func createCronJobEveryArea() {
	// ini untuk nambah cronjob sesuai dengan banyaknya area policy
	for _, ap := range GetAreaPolicy {
		areaPolicyObject := ap // create a new "areaPolicyObject" variable on each iteration (goroutines on loop iterator variables)
		// this function for push notif before
		s := strings.Split(areaPolicyObject.OrderTimeLimit, ":") // split from string to array
		toIntHour, _ := strconv.Atoi(s[0])
		entryIDNotif, _ := SecCron.AddFunc("CRON_TZ=Asia/Jakarta "+s[1]+" "+strconv.Itoa(toIntHour-1)+" * * *", func() { // s[1] = menit, 15 = jam
			pushNotification(areaPolicyObject, "reminderPayment")
		})

		// this function for create cronjob Cancel order
		entryIDCancelled, _ := SecCron.AddFunc("CRON_TZ=Asia/Jakarta "+s[1]+" "+s[0]+" * * *", func() {
			pushNotification(areaPolicyObject, "cancelled")
		})
		entryIdNotifyWillRemove = append(entryIdNotifyWillRemove, entryIDNotif)
		entryIdCancelledWillRemove = append(entryIdCancelledWillRemove, entryIDCancelled)
	}

	SecCron.Start()
}

type so struct {
	ID                   int64   `orm:"column(id);" json:"id"`
	Code                 string  `orm:"column(code);size(45);null" json:"code"`
	FirebaseToken        string  `orm:"column(firebase_token);null" json:"firebase_token"`
	MerchantID           int64   `orm:"column(merchant_id);null" json:"merchant_id"`
	TotalCharge          float64 `orm:"column(total_charge);null" json:"total_charge"`
	CreditLimitAmount    float64 `orm:"column(credit_limit_amount);null" json:"credit_limit_amount"`
	CreditLimitRemaining float64 `orm:"column(credit_limit_remaining);null" json:"credit_limit_remaining"`
}

func pushNotification(areaPolicy areaPolicy, status string) {
	var bulan string
	var salesOrder []so
	orSelect := orm.NewOrm()
	orSelect.Using("read_only")
	o := orm.NewOrm()
	mn := &MessageNotification{}
	modelNotif := &ModelNotification{}
	currentTime := time.Now()
	if status == "cancelled" {
		orSelect.Raw("SELECT so.id id, so.code code, um.firebase_token firebase_token, m.id merchant_id, "+
			"so.total_charge total_charge, m.credit_limit_amount credit_limit_amount, m.credit_limit_remaining credit_limit_remaining "+
			"FROM sales_order so "+
			"INNER JOIN branch b ON b.id = so.branch_id "+
			"INNER JOIN merchant m ON m.id = b.merchant_id "+
			"INNER JOIN user_merchant um ON um.id = m.user_merchant_id "+
			"WHERE so.area_id = ? AND so.status= 1 AND DATE(so.delivery_date) = ? AND so.payment_group_sls_id = 1 AND so.has_ext_invoice = 2;", areaPolicy.Area, currentTime.Add(time.Hour*24).Format("2006-01-02")).QueryRows(&salesOrder)
	} else if status == "reminderPayment" {
		orSelect.Raw("SELECT so.id id, so.code code, um.firebase_token firebase_token, m.id merchant_id "+
			"FROM sales_order so "+
			"INNER JOIN branch b ON b.id = so.branch_id "+
			"INNER JOIN merchant m ON m.id = b.merchant_id "+
			"INNER JOIN user_merchant um ON um.id = m.user_merchant_id "+
			"WHERE so.area_id = ? AND so.status= 1 AND DATE(so.delivery_date) = ? AND so.payment_group_sls_id = 1;", areaPolicy.Area, currentTime.Add(time.Hour*24).Format("2006-01-02")).QueryRows(&salesOrder)
	}

	if len(salesOrder) > 0 {
		for _, s := range salesOrder {
			if status == "cancelled" {
				orSelect.Raw("SELECT message, title FROM notification WHERE code= 'NOT0005'").QueryRow(&mn)
				mn.Message = ReplaceNotificationSalesOrder(mn.Message, "#sales_order_code#", s.Code)
				o.Raw("UPDATE sales_order SET status = 3 WHERE id = ?;", s.ID).Exec()

				if s.CreditLimitAmount > 0 {
					creditLimitBefore := s.CreditLimitRemaining
					creditLimitAfter := creditLimitBefore + s.TotalCharge
					o.Raw("UPDATE merchant SET credit_limit_remaining = ? WHERE id = ?", creditLimitAfter, s.MerchantID).Exec()

					o.Raw("INSERT INTO credit_limit_log (merchant_id, type, ref_id, credit_limit_before, credit_limit_after, note) "+
						"VALUES (?,'sales_order',?,?,?,'auto cancel sales order')", s.MerchantID, s.ID, creditLimitBefore, creditLimitAfter).Exec()
				}

				modelNotif.SendTo = s.FirebaseToken
				modelNotif.Title = mn.Title
				modelNotif.Message = mn.Message
				modelNotif.Type = "1"
				modelNotif.RefID = s.ID
				modelNotif.MerchantID = s.MerchantID
				modelNotif.ServerKey = ServerKeyFireBase
				PostModelNotification(modelNotif)

			} else if status == "reminderPayment" {
				orSelect.Raw("SELECT message, title FROM notification WHERE code= 'NOT0007'").QueryRow(&mn)
				year, month, day := time.Now().Date()
				if month.String() == "January" {
					bulan = "Jan"
				} else if month.String() == "February" {
					bulan = "Feb"
				} else if month.String() == "March" {
					bulan = "Mar"
				} else if month.String() == "April" {
					bulan = "Apr"
				} else if month.String() == "May" {
					bulan = "Mei"
				} else if month.String() == "June" {
					bulan = "Jun"
				} else if month.String() == "July" {
					bulan = "Jul"
				} else if month.String() == "August" {
					bulan = "Ags"
				} else if month.String() == "September" {
					bulan = "Sep"
				} else if month.String() == "October" {
					bulan = "Okt"
				} else if month.String() == "November" {
					bulan = "Nov"
				} else if month.String() == "December" {
					bulan = "Des"
				}
				mn.Message = ReplaceNotificationSalesOrder(mn.Message, "#current_date#", strconv.Itoa(day)+" "+bulan+" "+strconv.Itoa(year))
				mn.Message = ReplaceNotificationSalesOrder(mn.Message, "#time_limit#", areaPolicy.OrderTimeLimit)
				mn.Message = ReplaceNotificationSalesOrder(mn.Message, "#sales_order_code#", s.Code)

				modelNotif.SendTo = s.FirebaseToken
				modelNotif.Title = mn.Title
				modelNotif.Message = mn.Message
				modelNotif.Type = "1"
				modelNotif.RefID = s.ID
				modelNotif.MerchantID = s.MerchantID
				modelNotif.ServerKey = ServerKeyFireBase
				PostModelNotification(modelNotif)

			}
		}
	}
	return
}
