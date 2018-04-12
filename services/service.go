package services

import (
	"github.com/NeuronAccount/account/remotes/sms"
	"github.com/NeuronAccount/account/storages/neuron_account_db"
	"github.com/NeuronFramework/log"
	"github.com/dgrijalva/jwt-go"
	"go.uber.org/zap"
	"time"
)

type AccountServiceOptions struct {
}

type AccountService struct {
	logger     *zap.Logger
	options    *AccountServiceOptions
	userDB     *neuron_account_db.DB
	smsService *sms.Service
}

func NewAccountService(options *AccountServiceOptions) (s *AccountService, err error) {
	s = &AccountService{}
	s.logger = log.TypedLogger(s)
	s.options = options
	s.userDB, err = neuron_account_db.NewDB()
	if err != nil {
		return nil, err
	}
	s.smsService, err = sms.New()
	if err != nil {
		return nil, err
	}

	return s, nil
}

func (s *AccountService) generateJwt(accountId string) (tokenString string, err error) {
	expiresTime := time.Now().Add(time.Hour)
	userToken := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.StandardClaims{
		Subject:   accountId,
		ExpiresAt: expiresTime.Unix(),
	})
	tokenString, err = userToken.SignedString([]byte("0123456789"))
	if err != nil {
		return "", err
	}

	return tokenString, nil
}

func (s *AccountService) calcPasswordHash(password string) (passwordHash string) {
	return password
}
