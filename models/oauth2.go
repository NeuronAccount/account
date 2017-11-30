package models

type OAuth2Client struct {
	ClientId     string
	AccountId    string
	PasswordHash string
	RedirectUri  string
}

type OauthAuthorizationCode struct {
	Code          string
	ExpireSeconds int64
}

type OAuth2TokenRequest struct {
	GrantType         string
	AuthorizationCode string
	ClientID          string
	RedirectURI       string
	ResponseType      string
	Scope             string
	State             string
	RefreshToken      string
}

type OAuth2AccessToken struct {
	AccessToken  string
	TokenType    string
	ExpiresIn    int64
	RefreshToken string
	Scope        string
}
