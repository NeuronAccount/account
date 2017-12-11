package services

import (
	"context"
	"github.com/NeuronAccount/account/models"
	"github.com/NeuronAccount/account/storages/account_db"
	"github.com/NeuronFramework/errors"
)

func (s *AccountService) SmsLogin(ctx context.Context, phone string, smsCode string) (jwt string, err error) {
	//验证码
	dbSmsCode, err := s.accountDB.SmsCode.GetQuery().
		SceneType_Equal(models.SCENE_TYPE_SMS_LOGIN).And().PhoneNumber_Equal(phone).
		OrderBy(account_db.SMS_CODE_FIELD_CREATE_TIME, false).QueryOne(ctx, nil)
	if err != nil {
		return "", err
	}

	if dbSmsCode == nil || dbSmsCode.SmsCode != smsCode {
		return "", errors.BadRequest("InvalidSmsCode", "验证码错误")
	}

	dbAccount, err := s.accountDB.Account.GetQuery().
		PhoneNumber_Equal(phone).QueryOne(ctx, nil)
	if err != nil {
		return "", err
	}

	if dbAccount == nil {
		return "", errors.NotFound("帐号不存在")
	}

	//生成Token
	jwt, err = generateJwt(dbAccount.AccountId)
	if err != nil {
		return "", err
	}

	return jwt, nil
}
