package services

import (
	"context"
	"database/sql"
	"github.com/NeuronAccount/account/models"
	"github.com/NeuronAccount/account/storages/account_db"
	"github.com/NeuronFramework/errors"
	"github.com/NeuronFramework/rand"
)

func (s *AccountService) SmsSignup(ctx context.Context, phone string, smsCode string, password string) (jwt string, err error) {
	err = s.validateSmsCode(ctx, models.SmsSceneSignup, phone, smsCode)
	if err != nil {
		return "", err
	}

	//check account exists
	dbAccount, err := s.accountDB.Account.GetQuery().PhoneNumber_Equal(phone).QueryOne(ctx, nil)
	if err != nil {
		return "", err
	}
	if dbAccount != nil {
		return "", errors.AlreadyExists("帐号已存在")
	}

	dbAccount = &account_db.Account{}
	dbAccount.AccountId = rand.NextHex(16)
	dbAccount.PhoneNumber = sql.NullString{Valid: true, String: phone}
	dbAccount.PasswordHash = s.calcPasswordHash(password)
	_, err = s.accountDB.Account.Insert(ctx, nil, dbAccount)
	if err != nil {
		return "", err
	}

	//gen gwt
	jwt, err = s.generateJwt(dbAccount.AccountId)
	if err != nil {
		return "", err
	}

	return jwt, nil
}
