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

const (
	// Redis HSET which holds the issued tokens
	authCodeTokensSet = "OA2B_AC_Tokens"

	// Redis HSET which holds the issued grants until a token request is made.
	authCodeGrantSet = "OA2B_AC_Grants"

	// AuthCodeFlowID is prepended to a refresh token issued by the Authorization Code flow
	AuthCodeFlowID = "AUTHCODE"

	//auth userID
	AuthCodeFlowUserId = "UserIdAuthCodeFlow"
)