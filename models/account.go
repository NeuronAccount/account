package models

type OAuth2AuthorizeParams struct {
	ResponseType string
	ClientID     string
	Scope        string
	RedirectURI  string
	State        string
}
