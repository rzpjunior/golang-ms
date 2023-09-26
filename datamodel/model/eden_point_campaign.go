package model

import (
	"encoding/json"
	"time"

	"git.edenfarm.id/cuxs/common"
	"git.edenfarm.id/cuxs/orm"
)

func init() {
	orm.RegisterModel(new(EdenPointCampaign))
}

// EdenPointCampaign : struct to hold model data for database
type EdenPointCampaign struct {
	ID                 int64     `orm:"column(id);auto" json:"-"`
	Code               string    `orm:"column(code)" json:"code"`
	Name               string    `orm:"column(name)" json:"name"`
	Area               string    `orm:"column(area)" json:"area"`
	AreaName           string    `orm:"-" json:"area_name"`
	AreaNameArr        []string  `orm:"-" json:"area_name_arr"`
	Archetype          string    `orm:"column(archetype)" json:"archetype"`
	ArchetypeName      string    `orm:"-" json:"archetype_name"`
	ArchetypeNameArr   []string  `orm:"-" json:"archetype_name_arr"`
	TagCustomer        string    `orm:"column(tag_customer)" json:"tag_customer"`
	TagCustomerName    string    `orm:"-" json:"tag_customer_name"`
	TagCustomerNameArr []string  `orm:"-" json:"tag_customer_name_arr"`
	CampaignFilterType int8      `orm:"column(campaign_filter_type)" json:"campaign_filter_type"`
	Multiple           int8      `orm:"column(multiple)" json:"multiple"`
	ImageUrl           string    `orm:"column(image_url)" json:"image_url"`
	StartDate          time.Time `orm:"column(start_date)" json:"start_date"`
	EndDate            time.Time `orm:"column(end_date)" json:"end_date"`
	Note               string    `orm:"column(note)" json:"note"`
	Status             int8      `orm:"column(status)" json:"status"`
	CreatedAt          time.Time `orm:"column(created_at)" json:"created_at"`
	CreatedBy          int64     `orm:"column(created_by)" json:"created_by"`
}

// MarshalJSON : function to customize data struct into json. encrypt all key ids
func (m *EdenPointCampaign) MarshalJSON() ([]byte, error) {
	type Alias EdenPointCampaign

	return json.Marshal(&struct {
		ID        string `json:"id"`
		StartDate string `json:"start_date"`
		EndDate   string `json:"end_date"`
		CreatedAt string `json:"created_at"`
		CreatedBy string `json:"created_by"`
		*Alias
	}{
		ID:        common.Encrypt(m.ID),
		StartDate: m.StartDate.Format("2006-01-02 15:04:05"),
		EndDate:   m.EndDate.Format("2006-01-02 15:04:05"),
		CreatedAt: m.CreatedAt.Format("2006-01-02 15:04:05"),
		CreatedBy: common.Encrypt(m.CreatedBy),
		Alias:     (*Alias)(m),
	})
}

// Save : function to save data into database. will update when it has valid id, create otherwise
func (m *EdenPointCampaign) Save(fields ...string) (err error) {
	o := orm.NewOrm()
	if m.ID > 0 {
		_, err = o.Update(m, fields...)
	} else {
		m.ID, err = o.Insert(m)
	}
	return
}

// Read : function to get data from database
func (m *EdenPointCampaign) Read(fields ...string) error {
	o := orm.NewOrm()
	o.Using("read_only")
	return o.Read(m, fields...)
}
