package services

import (
	"github.com/NeuronFramework/log"
	"github.com/NeuronGroup/Account/storages/account"
	"go.uber.org/zap"
)

type AccountServiceOptions struct {
	AccountStorageConnectionString string
}

type AccountService struct {
	logger  *zap.Logger
	options *AccountServiceOptions
	db      *account.DB
}

func NewAccountService(options *AccountServiceOptions) (s *AccountService, err error) {
	s = &AccountService{}
	s.logger = log.TypedLogger(s)
	s.options = options
	s.db, err = account.NewDB(options.AccountStorageConnectionString)
	if err != nil {
		return nil, err
	}

	return s, nil
}
