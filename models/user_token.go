package models

const UserAccessTokenExpireSeconds = 5 * 60 //AccessToken有效期5分钟

type UserToken struct {
	AccessToken  string
	RefreshToken string
}
