package config

// Enum for the OAuth 2.0 flows
const (
	AuthCode    = 1
)

// AuthCodeConfig defines the variables required in the OAuth 2.0 Authorization Code flow
type AuthCodeConfig struct {
	ClientID     string `json:"clientID"`
	ClientSecret string `json:"clientSecret"`
}

// OA2Config defines the configurations for all the flows in OAuth 2.0
type OA2Config struct {
	BaseURL         string            `json:"baseURL"`
	AuthCodeCnfg    AuthCodeConfig    `json:"authCode"`
}