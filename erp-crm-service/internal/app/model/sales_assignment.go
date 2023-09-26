package model

import (
	"time"

	"git.edenfarm.id/edenlabs/edenlabs/orm"
)

type SalesAssignment struct {
	ID            int64     `orm:"column(id)" json:"-"`
	Code          string    `orm:"column(code)" json:"code"`
	TerritoryID   int64     `orm:"column(territory_id)" json:"territory_id"`
	TerritoryIDGP string    `orm:"column(territory_id_gp)" json:"territory_id_gp"`
	StartDate     time.Time `orm:"column(start_date)" json:"start_date"`
	EndDate       time.Time `orm:"column(end_date)" json:"end_date"`
	Status        int8      `orm:"column(status)" json:"status"`
}

func init() {
	orm.RegisterModel(new(SalesAssignment))
}

func (m *SalesAssignment) TableName() string {
	return "sales_assignment"
}
