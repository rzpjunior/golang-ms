package prospect_customer

import (
	"git.edenfarm.id/cuxs/common"
	"git.edenfarm.id/cuxs/validation"
	"git.edenfarm.id/project-version2/datamodel/model"
	"git.edenfarm.id/project-version2/api/util"
)

type createRequest struct {
	CodeProspectCustomer string
	Name                 string `json:"name" valid:"required"`
	ArchetypeID          string `json:"archetype_id" valid:"required"`
	PhoneNumber          string `json:"phone_number" valid:"required"`
	AltPhoneNumber       string `json:"alt_phone_number"`
	StreetAddress        string `json:"street_address" valid:"required"`
	PicName              string `json:"pic_name" valid:"required"`
	Email                string `json:"email" valid:"required"`
	PicPhoneNumber       string `json:"pic_phone_number"`
	SubDistrictID        string `json:"sub_district_id" valid:"required"`
	TimeConsent          int8   `json:"time_consent" valid:"required"`
	ReferenceInfo        string `json:"reference_info" valid:"required"`
	ReferralCode         string `json:"referral_code"`

	SubDistrict *model.SubDistrict `json:"-"`
	ArcheType   *model.Archetype   `json:"-"`
}

func (c *createRequest) Validate() *validation.Output {
	o := &validation.Output{Valid: true}
	var err error

	if c.CodeProspectCustomer, err = util.CheckTable("prospect_customer"); err != nil {
		o.Failure("code.invalid", util.ErrorInvalidData("code"))
	}

	subDistrictID, e := common.Decrypt(c.SubDistrictID)
	if e != nil {
		o.Failure("subDistrictID.invalid", util.ErrorInvalidData("sub district"))
	}
	c.SubDistrict = &model.SubDistrict{ID: subDistrictID}
	if e := c.SubDistrict.Read(); e != nil {
		o.Failure("subDistrictID.invalid", util.ErrorInvalidData("sub district"))
	}

	archetypeID, e := common.Decrypt(c.ArchetypeID)
	if e != nil {
		o.Failure("ArchetypeID.invalid", util.ErrorInvalidData("archetype"))
	}
	c.ArcheType = &model.Archetype{ID: archetypeID}
	if e := c.ArcheType.Read(); e != nil {
		o.Failure("ArchetypeID.invalid", util.ErrorInvalidData("archetype"))
	}

	return o
}

func (c *createRequest) Messages() map[string]string {
	return map[string]string{
		"street_address.required":  util.ErrorInputRequiredIndo("alamat usaha"),
		"archetype_id.required":    util.ErrorSelectRequiredIndo("jenis usaha"),
		"email.required":           util.ErrorInputRequiredIndo("email"),
		"name.required":            util.ErrorInputRequiredIndo("nama usaha"),
		"phone_number.required":    util.ErrorInputRequiredIndo("nomor telepon aktif"),
		"pic_name.required":        util.ErrorInputRequiredIndo("nama pemilik usaha"),
		"reference_info.required":  util.ErrorSelectRequiredIndo("info referensi"),
		"sub_district_id.required": util.ErrorSelectRequiredIndo("kelurahan"),
		"time_consent.required":    util.ErrorSelectRequiredIndo("waktu terbaik"),
	}
}
