package dto

type LoginRequest struct {
	UserName string `json:"userName"`
	Password string `json:"password"`
}

// type LoginResponse struct {
// 	Code    int    `json:"code"`
// 	Message string `json:"message"`
// 	Data    struct {
// 		User struct {
// 			ID         int    `json:"id"`
// 			Email      string `json:"email"`
// 			Activation bool   `json:"activation"`
// 		} `json:"user"`
// 		Token struct {
// 			AccessToken string `json:"access_token"`
// 			ExpiresIn   string `json:"access_token_expires_in"`
// 		} `json:"token"`
// 	} `json:"data"`
// }

type LoginResponse struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
	Data    struct {
		User struct {
			ID       int    `json:"id"`
			Username string `json:"username"`
		} `json:"user"`
		Token        string `json:"token"`
		RefreshToken string `json:"refreshToken"`
		Expiration   string `json:"expiration"`
	} `json:"data"`
}

type CommonGPResponse struct {
	Code     int    `json:"code"`
	Message  string `json:"message"`
	Sopnumbe string `json:"sopnumbe"`
}

type CommonPurchaseOrderGPResponse struct {
	Code     int    `json:"code"`
	Message  string `json:"message"`
	Ponumber string `json:"ponumber"`
}
