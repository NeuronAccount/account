package models

type OauthJumpParams struct {
	RedirectUri       string
	AuthorizationCode string
	State             string
}
