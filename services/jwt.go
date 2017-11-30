package services

import (
	"fmt"
	"github.com/NeuronGroup/Account/models"
)

func generateJwt(accountId string, oauth2AuthorizeParams *models.OAuth2AuthorizeParams) (jwt string, err error) {
	var scope string
	if oauth2AuthorizeParams == nil {
		scope = models.SCOPE_ROOT
	} else {
		scope = oauth2AuthorizeParams.Scope
	}

	fmt.Println(scope)

	return "jwt123456789", nil
}
