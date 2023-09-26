package dto

import (
	"time"
)

type SalespersonResponse struct {
	ID            int64     `json:"id"`
	Code          string    `json:"code"`
	FirstName     string    `json:"firstname"`
	MiddleName    string    `json:"namemiddle"`
	LastName      string    `json:"lastname"`
	Status        int8      `json:"status"`
	StatusConvert string    `json:"status_convert"`
	CreatedAt     time.Time `json:"created_at"`
	UpdatedAt     time.Time `json:"updated_at"`
}

type SalespersonListRequest struct {
	Limit       int32  `json:"limit"`
	Offset      int32  `json:"offset"`
	OffsetQuery int32  `json:"-"`
	Status      int32  `json:"status"`
	Search      string `json:"search"`
	OrderBy     string `json:"order_by"`
}

type SalespersonDetailRequest struct {
	Id int32 `json:"id"`
}

type SalespersonGP struct {
	Slprsnid        string  `json:"slprsnid,omitempty"`
	Employid        string  `json:"employid,omitempty"`
	Vendorid        string  `json:"vendorid,omitempty"`
	Slprsnfn        string  `json:"slprsnfn,omitempty"`
	Sprsnsmn        string  `json:"sprsnsmn,omitempty"`
	Sprsnsln        string  `json:"sprsnsln,omitempty"`
	AddresS1        string  `json:"addresS1,omitempty"`
	AddresS2        string  `json:"addresS2,omitempty"`
	AddresS3        string  `json:"addresS3,omitempty"`
	City            string  `json:"city,omitempty"`
	State           string  `json:"state,omitempty"`
	Zip             string  `json:"zip,omitempty"`
	Country         string  `json:"country,omitempty"`
	PhonE1          string  `json:"phonE1,omitempty"`
	PhonE2          string  `json:"phonE2,omitempty"`
	PhonE3          string  `json:"phonE3,omitempty"`
	Fax             string  `json:"fax,omitempty"`
	Commcode        string  `json:"commcode,omitempty"`
	Comprcnt        int64   `json:"comprcnt,omitempty"`
	Comappto        int64   `json:"comappto,omitempty"`
	Costtodt        float64 `json:"costtodt,omitempty"`
	Cstlstyr        float64 `json:"cstlstyr,omitempty"`
	Ttlcomtd        float64 `json:"ttlcomtd,omitempty"`
	Ttlcomly        float64 `json:"ttlcomly,omitempty"`
	Comsltdt        float64 `json:"comsltdt,omitempty"`
	Comsllyr        float64 `json:"comsllyr,omitempty"`
	Ncomsltd        float64 `json:"ncomsltd,omitempty"`
	Ncomslyr        float64 `json:"ncomslyr,omitempty"`
	Noteindex       int64   `json:"noteindex,omitempty"`
	GnL_Supervisor1 string  `json:"gnL_Supervisor1,omitempty"`
}

type SalesPerson struct {
	ID                 int64  `orm:"column(id);auto" json:"-"`
	StaffID            string `orm:"-" json:"staff_id,omitempty"`
	Code               string `orm:"column(code);size(50);null" json:"code"`
	Name               string `orm:"column(name);size(100);null" json:"name"`
	DisplayName        string `orm:"column(display_name);size(100);null" json:"display_name"`
	EmployeeCode       string `orm:"column(employee_code);size(50);null" json:"employee_code"`
	RoleGroup          int8   `orm:"column(role_group);null" json:"role_group"`
	PhoneNumber        string `orm:"column(phone_number);size(15);null" json:"phone_number"`
	Status             int8   `orm:"column(status);null" json:"status"`
	SalesGroupID       int64  `orm:"column(sales_group_id);null;" json:"sales_group_id,omitempty"`
	SalesGroupName     string `orm:"-" json:"sales_group_name,omitempty"`
	WarehouseAccessStr string `orm:"column(warehouse_access)" json:"warehouse_access_str"`
	StatusConvert      string `orm:"-" json:"status_convert"`
	// WarehouseAccess    []*Warehouse `orm:"-" json:"warehouse_access"`

	// Picking List Module
	UsedStaff bool `orm:"-" json:"used_staff"`
	IsBusy    bool `orm:"-" json:"is_busy"`

	// Role      *Role      `orm:"column(role_id);null;rel(fk)" json:"role,omitempty"`
	// User      *User      `orm:"column(user_id);null;rel(fk)" json:"user,omitempty"`
	// Area      *Area      `orm:"column(area_id);null;rel(fk)" json:"area,omitempty"`
	// Parent    *Staff     `orm:"column(parent_id);null;rel(fk)" json:"parent,omitempty"`
	// Warehouse *Warehouse `orm:"column(warehouse_id);null;rel(fk)" json:"warehouse,omitempty"`

}
