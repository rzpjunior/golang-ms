package dto

import (
	"time"
)

type SalesPaymentResponse struct {
	ID               string    `json:"id"`
	Code             string    `json:"code"`
	Status           int8      `json:"status"`
	RecognitionDate  time.Time `json:"recognition_date"`
	Amount           float64   `json:"amount"`
	BankReceiveNum   string    `json:"bank_receive_num"`
	PaidOff          int8      `json:"paid_off"`
	ImageUrl         string    `json:"image_url"`
	Note             string    `json:"note"`
	CreatedAt        time.Time `json:"created_at"`
	CreatedBy        int64     `json:"created_by"`
	CancellationNote string    `json:"cancellation_note,omitempty"`
	ReceivedDate     time.Time `json:"received_date"`
	CustomerID       string    `json:"customer_id"`
	CustomerName     string    `json:"customer_name"`

	PaymentMethod string                `json:"payment_method"`
	SalesInvoice  *SalesInvoiceResponse `json:"sales_invoice"`
}

type SalesPaymentListRequest struct {
	Limit                     int32     `json:"limit"`
	Offset                    int32     `json:"offset"`
	Status                    string    `json:"status"`
	Search                    string    `json:"search"`
	OrderBy                   string    `json:"order_by"`
	SiteID                    string    `json:"-"`
	RecognitionDateFrom       time.Time `json:"-"`
	RecognitionDateTo         time.Time `json:"-"`
	RecognitionDateFromString string    `json:"recognition_date_from"`
	RecognitionDateToString   string    `json:"recognition_date_to"`
	CustomerID                string    `json:"-"`
}

type SalesPaymentDetailRequest struct {
	Id string `json:"id"`
}

type CreateSalesPaymentGPRequest struct {
	Interid        string  `json:"interid"`
	Sopnumbe       string  `json:"sopnumbe"`
	AmountReceived float64 `json:"amount_received"`
	Chekbkid       string  `json:"chekbkid"`
	Docdate        string  `json:"docdate"`
	RegionID       string  `json:"-"`
}

type CreateSalesPaymentRequest struct {
	Code                   string    `json:"-"`
	PaymentDateStr         string    `json:"payment_date" `
	PaymentMethodID        string    `json:"payment_method_id" valid:"required"`
	PaymentChannelID       string    `json:"payment_channel_id"`
	Amount                 float64   `json:"amount" valid:"required"`
	PaidOff                int8      `json:"paid_off"`
	Note                   string    `json:"note"`
	ImageUrl               string    `json:"image_url"`
	SalesInvoiceID         string    `json:"sales_invoice_id" valid:"required"`
	BankReceiveNum         string    `json:"bank_receive_num"`
	PaymentDate            time.Time `json:"-"`
	HaveCreditLimit        bool      `json:"-"`
	CreditLimitBefore      float64   `json:"-"`
	CreditLimitAfter       float64   `json:"-"`
	RemainingInvoiceAmount float64   `json:"-"`
	RegionID               string    `json:"-"`
	CheckbookID            string    `json:"checkbook_id" valid:"required"`
}

type SI_SalesPaymentGP struct {
	Sopnumbe  string `json:"sopnumbe"`
	Docdate   string `json:"docdate"`
	GnlRegion string `json:"gnl_region"`
	Locncode  string `json:"locncode"`
}

type SO_SalesPaymentGP struct {
	Orignumb string `json:"orignumb"`
	Ordrdate string `json:"ordrdate"`
}

type SalesPaymentGP struct {
	Docnumbr      string               `json:"docnumbr"`
	Docdate       string               `json:"docdate"`
	Custnmbr      string               `json:"custnmbr"`
	Custname      string               `json:"custname"`
	Curncyid      string               `json:"curncyid"`
	Cshrctyp      int                  `json:"cshrctyp"`
	PaymentMethod string               `json:"payment_method"`
	Dcstatus      int                  `json:"dcstatus"`
	Ortrxamt      float64              `json:"ortrxamt"`
	Creatddt      string               `json:"creatddt"`
	SalesInvoice  []*SI_SalesPaymentGP `json:"sales_invoice"`
	SalesOrder    []*SO_SalesPaymentGP `json:"sales_order"`
	Ordocamt      float64              `json:"ordocamt"`
}

type GetSalesPaymentGPListRequest struct {
	Limit         int32  `json:"limit"`
	Offset        int32  `json:"offset"`
	Docnumbr      string `json:"docnumbr"`
	DocdateFrom   string `json:"docdate_from"`
	DocdateTo     string `json:"docdate_to"`
	SiDocdateFrom string `json:"si_docdate_from"`
	SiDocdateTo   string `json:"si_docdate_to"`
	SoDocdateFrom string `json:"so_docdate_from"`
	SoDocdateTo   string `json:"so_docdate_to"`
	Custnmbr      string `json:"custnmbr"`
	Sopnumbe      string `json:"sopnumbe"`
	GnlRegion     string `json:"gnl_region"`
	Locncode      string `json:"locncode"`
}

// Purchase Order: struct to hold model data for database
type SalesPayment struct {
	ID               int64     `orm:"column(id);auto" json:"-"`
	Code             string    `orm:"column(code);size(50);null" json:"code"`
	Status           int8      `orm:"column(status)" json:"status"`
	RecognitionDate  time.Time `orm:"column(recognition_date)" json:"recognition_date"`
	Amount           float64   `orm:"column(amount)" json:"amount"`
	BankReceiveNum   string    `orm:"column(bank_receive_num)" json:"bank_receive_num"`
	PaidOff          int8      `orm:"column(paid_off)" json:"paid_off"`
	ImageUrl         string    `orm:"column(image_url);null" json:"image_url"`
	Note             string    `orm:"column(note)" json:"note"`
	CreatedAt        time.Time `orm:"column(created_at);type(timestamp);null" json:"created_at"`
	CreatedBy        int64     `orm:"column(created_by)" json:"created_by"`
	CancellationNote string    `orm:"-" json:"cancellation_note,omitempty"`
	ReceivedDate     time.Time `orm:"column(received_date);type(date);null" json:"received_date"`

	// TxnXendit      *TxnXendit      `orm:"column(txn_xendit_id);null;rel(fk)" json:"txn_xendit"`
	// SalesInvoice   *SalesInvoice   `orm:"column(sales_invoice_id);null;rel(fk)" json:"sales_invoice"`
	// PaymentMethod  *PaymentMethod  `orm:"column(payment_method_id);null;rel(fk)" json:"payment_method"`
	// PaymentChannel *PaymentChannel `orm:"column(payment_channel_id);null;rel(fk)" json:"payment_channel"`
}

type PerformancePaymentRequest struct {
	Limit                     int32     `json:"limit"`
	Offset                    int32     `json:"offset"`
	Status                    string    `json:"status"`
	Search                    string    `json:"search"`
	OrderBy                   string    `json:"order_by"`
	SiteID                    string    `json:"-"`
	RecognitionDateFrom       time.Time `json:"-"`
	RecognitionDateTo         time.Time `json:"-"`
	RecognitionDateFromString string    `json:"recognition_date_from"`
	RecognitionDateToString   string    `json:"recognition_date_to"`
	CustomerID                string    `json:"customer_id"`
}

type PaymentPerformance struct {
	CreditLimitAmount                   float64 `json:"credit_limit_amount"`
	CreditLimitRemainingAmount          float64 `json:"credit_limit_remaining_amount"`
	RemainingOutstanding                float64 `json:"remaining_outstanding"`
	CreditLimitUsageRemainingPercentage float64 `json:"credit_limit_usage_remaining_percentage"`
	OverdueDebtAmount                   float64 `json:"overdue_debt_amount"`
	OverdueDebtRemainingPercentage      float64 `json:"overdue_debt_remaining_percentage"`
	AveragePaymentAmount                float64 `json:"average_payment_amount"`
	AveragePaymentPercentage            float64 `json:"average_payment_percentage"`
	AveragePaymentPeriod                int     `json:"average_payment_period"`
}
