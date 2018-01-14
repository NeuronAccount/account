package services

import (
	"context"
	"github.com/NeuronAccount/account/storages/account_db"
	"github.com/NeuronFramework/errors"
	"strings"
)

func (s *AccountService) calcPasswordHash(password string) (passwordHash string) {
	return password
}

func (s *AccountService) Login(ctx context.Context, name string, password string) (jwt string, err error) {
	var dbAccount *account_db.Account
	if strings.Contains(name, "@") { //email
		dbAccount, err = s.accountDB.Account.GetQuery().EmailAddress_Equal(name).QueryOne(ctx, nil)
	} else { //phone
		dbAccount, err = s.accountDB.Account.GetQuery().PhoneNumber_Equal(name).QueryOne(ctx, nil)
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
