package models

type OauthAccountInfo struct {
	ProviderId   string
	ProviderName string
	OpenId       string
	UserName     string
	UserIcon     string
}

type AccountInfo struct {
	UserId          string
	UserName        string
	UserIcon        string
	PhoneBinded     string
	OauthBindedList []*OauthAccountInfo
}
