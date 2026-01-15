package models

type RegisterRequest struct {
	Username string `json:"username"`
	EmailAddress string `json:"email"`
	Password string `json:"password"`
	FullName string `json:"full_name,omitempty"`
}

type RegisterTokenResponse struct {
	Status bool `json:"status"`
	AccessToken string `json:"access_token"` 
	RefreshToken string `json:"refresh_token"`
}

type LoginRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

type TokenResponse struct {
	AccessToken string `json:"access_token"` 
	RefreshToken string `json:"refresh_token"`
}
