package services

import (
	"context"
	"github.com/NeuronFramework/errors"
	"github.com/NeuronGroup/account/storages/account"
	"strings"
)

func (s *AccountService) calcPasswordHash(password string) (passwordHash string) {
	return password
}

func (s *AccountService) Login(name string, password string) (jwt string, err error) {
	var dbAccount *account.Account
	if strings.Contains(name, "@") { //email
		dbAccount, err = s.db.Account.GetQuery().EmailAddress_Equal(name).QueryOne(context.Background(), nil)
	} else { //phone
		dbAccount, err = s.db.Account.GetQuery().PhoneNumber_Equal(name).QueryOne(context.Background(), nil)
	}

	if err != nil {
		return "", err
	}

	if dbAccount == nil {
		return "", errors.NotFound("帐号不存在")
	}

	passwordHash := s.calcPasswordHash(password)
	if dbAccount.PasswordHash != passwordHash {
		return "", errors.Unauthorized("帐号或密码错误")
	}

	jwt, err = generateJwt(dbAccount.AccountId)
	if err != nil {
		return "", err
	}

	return jwt, nil
}
