package services

import (
	"github.com/NeuronAccount/account/storages/account_db"
	"github.com/NeuronFramework/log"
	"go.uber.org/zap"
)

type AccountServiceOptions struct {
}

type AccountService struct {
	logger    *zap.Logger
	options   *AccountServiceOptions
	accountDB *account_db.DB
}

func NewAccountService(options *AccountServiceOptions) (s *AccountService, err error) {
	s = &AccountService{}
	s.logger = log.TypedLogger(s)
	s.options = options
	s.accountDB, err = account_db.NewDB("root:123456@tcp(127.0.0.1:3307)/account?parseTime=true")
	if err != nil {
		return nil, err
	}

	return s, nil
}
