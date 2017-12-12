package services

import (
	"context"
	"github.com/NeuronAccount/account/models"
	"github.com/NeuronAccount/account/storages/account_db"
	"github.com/NeuronFramework/errors"
	"github.com/NeuronFramework/rand"
	"github.com/NeuronFramework/sql/wrap"
)

func (s *AccountService) SmsSignup(ctx context.Context, phone string, smsCode string, password string) (jwt string, err error) {
	dbSmsCode, err := s.accountDB.SmsCode.GetQuery().
		SceneType_Equal(models.SCENE_TYPE_SMS_SIGNUP).And().PhoneNumber_Equal(phone).
		OrderBy(account_db.SMS_CODE_FIELD_CREATE_TIME, false).QueryOne(ctx, nil)
	if err != nil {
		return "", err
	}

	if dbSmsCode == nil || dbSmsCode.SmsCode != smsCode {
		return "", errors.BadRequest("InvalidSmsCode", "验证码错误")
	}

	accountId := rand.NextHex(16)

	err = s.accountDB.TransactionReadCommitted(ctx, func(tx *wrap.Tx) (err error) {
		dbAccount, err := s.accountDB.Account.GetQuery().ForUpdate().
			PhoneNumber_Equal(phone).
			QueryOne(ctx, tx)
		if err != nil {
			return err
		}

		if dbAccount != nil {
			return errors.AlreadyExists("帐号已存在")
		}

		dbAccount = &account_db.Account{}
		dbAccount.AccountId = accountId
		dbAccount.PhoneNumber = phone
		dbAccount.EmailAddress = ""
		dbAccount.PasswordHash = s.calcPasswordHash(password)
		dbAccount.OauthProvider = ""
		dbAccount.OauthAccountId = ""

		_, err = s.accountDB.Account.Insert(ctx, tx, dbAccount)
		if err != nil {
			return err
		}

		return nil
	})
	if err != nil {
		return "", err
	}

	//生成Token
	jwt, err = generateJwt(accountId)
	if err != nil {
		return "", err
	}

	return jwt, nil
}
