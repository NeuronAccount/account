package services

import (
	"database/sql"
	"github.com/NeuronAccount/account/models"
	"github.com/NeuronAccount/account/storages/account_db"
	"github.com/NeuronFramework/errors"
	"github.com/NeuronFramework/rand"
	"github.com/NeuronFramework/restful"
)

func (s *AccountService) SmsSignup(ctx *restful.Context, phone string, smsCode string, password string) (jwt string, err error) {
	//check account exists
	dbAccount, err := s.accountDB.Account.GetQuery().PhoneNumber_Equal(phone).QueryOne(ctx, nil)
	if err != nil {
		return "", err
	}
	if dbAccount != nil {
		return "", errors.AlreadyExists("帐号已存在，请直接登录")
	}

	err = s.validateSmsCode(ctx, models.SmsSceneSignup, phone, smsCode)
	if err != nil {
		return "", err
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

	s.addOperation(ctx, &models.Operation{
		OperationType: models.OperationSmsSignup,
		Phone:         phone,
		AccountID:     dbAccount.AccountId,
	})

	return jwt, nil
}
