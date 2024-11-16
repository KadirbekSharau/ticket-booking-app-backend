package responses

type TokenResponse struct {
	Success               bool    `json:"success"`
	Token                 string `json:"token"`
	TokenType             string `json:"token_type"`
	ExpiresAt             int64  `json:"expires_at"`
	RefreshToken          string `json:"refresh_token"`
	RefreshTokenType      string `json:"refresh_token_type"`
	RefreshTokenExpiresAt int64  `json:"refresh_token_expires_at"`
}
