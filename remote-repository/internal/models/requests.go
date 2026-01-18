package models

type RegisterRequest struct {
	Username     string `json:"username"`
	EmailAddress string `json:"email"`
	Password     string `json:"password"`
	FullName     string `json:"full_name,omitempty"`
}

type LoginRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

type TokenResponse struct {
	Status       bool   `json:"status,omitempty"`
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
}

type LogoutRequest struct {
	RefreshToken string `json:"refresh_token"`
}