package services

import (
	"fmt"
	"github.com/NeuronFramework/log"
	"github.com/NeuronGroup/Account/storages/account"
	"github.com/NeuronGroup/Account/storages/oauth"
	"go.uber.org/zap"
)

type AccountServiceOptions struct {
	AccountStorageConnectionString string
	OAuthConnectionString          string
}

type AccountService struct {
	logger  *zap.Logger
	options *AccountServiceOptions
	db      *account.DB
	oauthDB *oauth.DB
}

func NewAccountService(options *AccountServiceOptions) (s *AccountService, err error) {
	s = &AccountService{}
	s.logger = log.TypedLogger(s)
	s.options = options
	s.db, err = account.NewDB(options.AccountStorageConnectionString)
	if err != nil {
		return nil, err
	}
	fmt.Println("conn", options.OAuthConnectionString)

	s.oauthDB, err = oauth.NewDB(options.OAuthConnectionString)
	if err != nil {
		return nil, err
	}

	return s, nil
}
