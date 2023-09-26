package prospect_customer

import "git.edenfarm.id/cuxs/validation"

type requestGetListDeclineType struct {
	DeclineType []*declineType
}

type declineType struct {
	ValueInt    int8   `json:"-"`
	ValueIntEnc string `json:"value_int"`
	ValueName   string `json:"value_name"`
}

func (c *requestGetListDeclineType) Validate() *validation.Output {
	o := &validation.Output{Valid: true}

	return o
}

func (c *requestGetListDeclineType) Messages() map[string]string {
	return map[string]string{}
}
