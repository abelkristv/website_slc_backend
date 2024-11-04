package config

type TokenResponse struct {
	AccessToken  string `json:"access_token"`
	TokenType    string `json:"token_type"`
	ExpiresIn    int    `json:"expires_in"`
	RefreshToken string `json:"refresh_token"`
}

const BaseURL = "https://bluejack.binus.ac.id/lapi/api"

var AuthToken TokenResponse
