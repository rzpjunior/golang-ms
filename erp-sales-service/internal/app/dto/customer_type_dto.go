package dto

type CustomerTypeResponse struct {
	ID   string `json:"id"`
	Name string `json:"name,omitempty"`
}

type GetCustomerTypeRequest struct {
	Limit  int
	Offset int
}
