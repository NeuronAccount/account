package services

import (
	"context"
	"github.com/NeuronAccount/account/models"
	"github.com/NeuronAccount/account/storages/account"
	"github.com/NeuronFramework/errors"
	"time"
)

func (s *AccountService) SmsLogin(phone string, smsCode string) (jwt string, err error) {
	//验证码
	dbSmsCode, err := s.db.SmsCode.GetQuery().
		SceneType_Equal(models.SCENE_TYPE_SMS_LOGIN).And().PhoneNumber_Equal(phone).
		OrderBy(account.SMS_CODE_FIELD_CREATE_TIME, false).QueryOne(context.Background(), nil)
	if err != nil {
		return "", err
	}

	if dbSmsCode == nil || dbSmsCode.SmsCode != smsCode {
		return "", errors.BadRequest("InvalidSmsCode", "验证码错误")
	}

	dbAccount, err := s.db.Account.GetQuery().
		PhoneNumber_Equal(phone).QueryOne(context.Background(), nil)
	if err != nil {
		return "", err
	}

	//是否新建帐号
	if dbAccount == nil {
		tx, err := s.db.BeginReadCommittedTx(context.Background(), false)
		if err != nil {
			return "", err
		}
		defer tx.Rollback()

		dbAccount, err := s.db.Account.GetQuery().ForUpdate().PhoneNumber_Equal(phone).QueryOne(context.Background(), tx)
		if err != nil {
			return "", err
		}

		if dbAccount == nil {
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
				return "", err
			}
		} //nothing

		err = tx.Commit()
		if err != nil {
			return "", err
		}
	}

	//生成Token
	jwt, err = generateJwt(dbAccount.AccountId)
	if err != nil {
		return "", err
	}

	return jwt, nil
}
