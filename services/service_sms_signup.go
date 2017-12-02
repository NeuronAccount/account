package services

import (
	"context"
	"github.com/NeuronFramework/errors"
	"github.com/NeuronFramework/rand"
	"github.com/NeuronGroup/account/models"
	"github.com/NeuronGroup/account/storages/account"
)

func (s *AccountService) SmsSignup(phone string, smsCode string, password string) (jwt string, err error) {
	dbSmsCode, err := s.db.SmsCode.GetQuery().
		SceneType_Equal(models.SCENE_TYPE_SMS_SIGNUP).And().PhoneNumber_Equal(phone).
		OrderBy(account.SMS_CODE_FIELD_CREATE_TIME, false).QueryOne(context.Background(), nil)
	if err != nil {
		return "", err
	}

	if dbSmsCode == nil || dbSmsCode.SmsCode != smsCode {
		return "", errors.BadRequest("InvalidSmsCode", "验证码错误")
	}

	tx, err := s.db.BeginReadCommittedTx(context.Background(), false)
	if err != nil {
		return "", err
	}
	defer tx.Rollback()

	dbAccount, err := s.db.Account.GetQuery().ForUpdate().
		PhoneNumber_Equal(phone).QueryOne(context.Background(), tx)
	if err != nil {
		return "", err
	}

	if dbAccount != nil {
		return "", errors.AlreadyExists("帐号已存在")
	}

	dbAccount = &account.Account{}
	dbAccount.AccountId = rand.NextBase64(16)
	dbAccount.PhoneNumber = phone
	dbAccount.EmailAddress = ""
	dbAccount.PasswordHash = s.calcPasswordHash(password)
	dbAccount.OauthProvider = ""
	dbAccount.OauthAccountId = ""

	_, err = s.db.Account.Insert(context.Background(), tx, dbAccount)
	if err != nil {
		return "", err
	}

	err = tx.Commit()
	if err != nil {
		return "", err
	}

	//生成Token
	jwt, err = generateJwt(dbAccount.AccountId)
	if err != nil {
		return "", err
	}

	return jwt, nil
}
