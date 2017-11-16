package services

import (
	"context"
	"github.com/NeuronFramework/errors"
	"github.com/NeuronGroup/Account/models"
	"github.com/NeuronGroup/Account/storages/account"
	"time"
)

func (s *AccountService) SmsSignup(phone string, smsCode string) (err error) {
	dbSmsCode, err := s.db.SmsCode.GetQuery().
		SceneType_Equal(models.SCENE_TYPE_SMS_SIGNUP).And().PhoneNumber_Equal(phone).
		OrderBy(account.SMS_CODE_FIELD_CREATE_TIME, false).QueryOne(context.Background(), nil)
	if err != nil {
		return err
	}

	if dbSmsCode == nil || dbSmsCode.SmsCode != smsCode {
		return errors.BadRequest("InvalidSmsCode", "验证码错误")
	}

	tx, err := s.db.BeginReadCommittedTx(context.Background(), false)
	if err != nil {
		return err
	}
	defer tx.Rollback()

	dbAccount, err := s.db.Account.GetQuery().ForUpdate().
		PhoneNumber_Equal(phone).QueryOne(context.Background(), tx)
	if err != nil {
		return err
	}

	if dbAccount != nil {
		return errors.AlreadyExists("帐号已存在")
	}

	dbAccount = &account.Account{}
	dbAccount.CreateTime = time.Now()
	dbAccount.UpdateTime = time.Now()
	dbAccount.PhoneNumber = phone
	dbAccount.EmailAddress = ""
	dbAccount.PasswordHash = ""
	dbAccount.OauthProvider = ""
	dbAccount.OauthAccountId = ""
	dbAccount.AccountId = "1234567890"
	_, err = s.db.Account.Insert(context.Background(), tx, dbAccount)
	if err != nil {
		return nil
	}

	err = tx.Commit()
	if err != nil {
		return err
	}

	return nil
}
