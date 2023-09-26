package dto

type ItemCategoryResponse struct {
	ID       string `orm:"column(id);auto" json:"id"`
	Region   string `orm:"column(region)" json:"region"`
	Name     string `orm:"column(name)" json:"name"`
	ImageUrl string `orm:"column(image_url)" json:"image_url"`
	Status   string `orm:"column(status)" json:"status"`
}

type ItemCategoryMobileRequest struct {
	Platform string              `json:"platform" valid:"required"`
	Data     dataGetItemCategory `json:"data"`
}

type dataGetItemCategory struct {
	AddressID string `json:"address_id"`
}
