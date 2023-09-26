package model

type SearchSuggestion struct {
	ID   int64  `orm:"column(id);auto" json:"-"`
	Code string `orm:"column(code)" json:"code"`
	Name string `orm:"column(name)" json:"name"`
}