package dto

import "time"

type UserCustomerResponse struct {
	ID            int64     `json:"id"`
	Code          string    `json:"code,omitempty"`
	CustomerID    int64     `json:"customer_id,omitempty"`
	FirebaseToken string    `json:"firebase_token,omitempty"`
	Verification  int8      `json:"verification,omitempty"`
	TncAccVersion string    `json:"tnc_acc_version,omitempty"`
	TncAccAt      time.Time `json:"tnc_acc_at"`
	LastLoginAt   time.Time `json:"last_login_at"`
	Note          string    `json:"note,omitempty"`
	Status        int8      `json:"status,omitempty"`
	ForceLogout   int8      `json:"force_logout,omitempty"`
	LoginToken    string    `json:"login_token,omitempty"`
}

type GetDetailUserCustomerRequest struct {
	CustomerID int64 `json:"customer_id,omitempty"`
}
