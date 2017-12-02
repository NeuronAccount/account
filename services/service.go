package services

import (
	"github.com/NeuronFramework/log"
	"github.com/NeuronGroup/account/storages/account"
	"go.uber.org/zap"
)

type AccountServiceOptions struct {
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
	s.db, err = account.NewDB("root:123456@tcp(127.0.0.1:3307)/account?parseTime=true")
	if err != nil {
		return nil, err
	}

	return s, nil
}
