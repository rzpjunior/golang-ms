package model

type Division struct {
	ID            int64  `orm:"column(id);auto" json:"-"`
	Code          string `orm:"column(code);size(50);null" json:"code"`
	Name          string `orm:"column(name);size(100);null" json:"name"`
	Note          string `orm:"column(note);size(250);null" json:"note"`
	Status        int8   `orm:"column(status);null" json:"status"`
	StatusConvert string `orm:"-" json:"status_convert"`
}

// TODO: init to GP
// func init() {
// 	orm.RegisterModel(new(Division))
// }

// func (m *Division) TableName() string {
// 	return "division"
// }
