package services

import (
	"github.com/NeuronAccount/account/remotes/sms"
	"github.com/NeuronAccount/account/storages/neuron_account_db"
	"github.com/NeuronFramework/log"
	"github.com/NeuronFramework/rand"
	"github.com/dgrijalva/jwt-go"
	"go.uber.org/zap"
	"strings"
	"time"
)

type AccountServiceOptions struct {
}

type AccountService struct {
	logger     *zap.Logger
	options    *AccountServiceOptions
	accountDB  *neuron_account_db.DB
	smsService *sms.Service
}

func NewAccountService(options *AccountServiceOptions) (s *AccountService, err error) {
	s = &AccountService{}
	s.logger = log.TypedLogger(s)
	s.options = options
	s.accountDB, err = neuron_account_db.NewDB()
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
		Id:        rand.NextHex(16),
		Subject:   accountId,
		ExpiresAt: expiresTime.Unix(),
	})
	tokenString, err = userToken.SignedString([]byte("0123456789"))
	if err != nil {
		return "", err
	}

	return tokenString, nil
}

func (s *AccountService) encryptPhone(phone string) (phoneEncrypted string, err error) {
	return phone, nil
}

func (s *AccountService) maskPhone(phone string) (phoneMasked string) {
	l := len(phone)
	if l <= 6 {
		return phone
	}

	return phone[0:2] + strings.Repeat("*", l-6) + phone[l-3:l-1]
}

func (s *AccountService) maskString(p string, nPrefix int, nSuffix int) (r string) {
	l := len(p)
	if nPrefix < 0 {
		nPrefix = 0
	}
	if nPrefix > l {
		nPrefix = l
	}
	if nSuffix < 0 {
		nSuffix = 0
	}
	if nSuffix > l {
		nSuffix = l
	}

	if l <= nPrefix+nSuffix {
		return p
	}

	return p[0:nPrefix] + strings.Repeat("*", l-nPrefix-nSuffix) + p[l-nSuffix-1:l-1]
}

func (s *AccountService) calcPasswordHash(passwordHash1 string) (passwordHash2 string, err error) {
	return passwordHash1, nil
}
