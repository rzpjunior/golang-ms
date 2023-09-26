package dto

import "time"

type IsPointUsedRequest struct {
	Platform string `json:"platform" valid:"required"`

	Session *SessionDataCustomer
}

type IsPointUsedResponse struct {
	IsPointUsed bool `json:"is_point_used"`
}

type GetPotentialEdenPointRequest struct {
	Platform string                 `json:"platform" valid:"required"`
	Data     *GetPotentialEdenPoint `json:"data" valid:"required"`

	Session *SessionDataCustomer
}

type GetPotentialEdenPoint struct {
	AddressID   string                       `json:"address_id" valid:"required"`
	RedeemCode  string                       `json:"redeem_code"`
	OrderTypeID string                       `json:"order_type_id"`
	Items       []*GetItemPotentialEdenPoint `json:"items" valid:"required"`
}

type GetItemPotentialEdenPoint struct {
	ID          string  `json:"id"`
	Price       float64 `json:"price"`
	TotalWeight float64 `json:"total_weight"`
}

type GetItemPotentialEdenPointResponse struct {
	Points float64 `json:"points"`
}

type RequestGetPointHistory struct {
	Platform string `json:"platform" valid:"required"`
	Offset   int16  `json:"offset" valid:"required"`
	Limit    int16  `json:"limit" valid:"required"`
	Session  *SessionDataCustomer
	//Data     []*PointHistoryList
}

type PointHistoryList struct {
	SalesOrderCode string ` json:"code"`
	CreatedDate    string ` json:"created_date"`
	PointValue     string ` json:"point_value"`
	StatusType     string ` json:"status_type"`
	Status         string ` json:"status"`
}

type GetCustomerPointExpirationRequest struct {
	Platform string `json:"platform" valid:"required"`

	Session *SessionDataCustomer
}

type GetCustomerPointExpirationResponse struct {
	ID                    int64     `json:"id"`
	CustomerID            int64     `json:"customer_id"`
	CurrentPeriodPoint    float64   `json:"current_period_point"`
	NextPeriodPoint       float64   `json:"next_period_point"`
	CurrentPeriodDate     time.Time `json:"current_period_date"`
	NextPeriodDate        time.Time `json:"next_period_date"`
	LastUpdatedAt         time.Time `json:"last_updated_at"`
	IsHavePointExpiration bool      `orm:"-" json:"is_have_point_expiration"`
}
