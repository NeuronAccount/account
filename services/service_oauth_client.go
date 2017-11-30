package services

import (
	"context"
	"fmt"
	"github.com/NeuronGroup/Account/models"
	"github.com/NeuronGroup/Account/storages/oauth"
)

func (s *AccountService) OAuth2ClientLogin(clientId string, password string) (c *models.OAuth2Client, err error) {
	dbClient, err := s.oauthDB.OauthClient.GetQuery().ClientId_Equal(clientId).QueryOne(context.Background(), nil)
	if err != nil {
		return nil, err
	}

	if dbClient == nil {
		return nil, fmt.Errorf("clientId not exists")
	}

	if dbClient.PasswordHash != password {
		return nil, fmt.Errorf("password failed")
	}

	return oauth.FromOauthClient(dbClient), nil
}
