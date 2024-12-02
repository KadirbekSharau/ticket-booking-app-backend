package responses

type TokenResponse struct {
	Success               bool    `json:"success"`
	Token                 string `json:"token"`
	TokenType             string `json:"token_type"`
	ExpiresAt             int64  `json:"expires_at"`
}
