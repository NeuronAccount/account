package services

import (
	"context"
	"github.com/NeuronAccount/account/models"
	"github.com/NeuronAccount/account/storages/account_db"
	"github.com/NeuronFramework/errors"
	"time"
)

func (s *AccountService) SmsCode(ctx context.Context, scene string, phone string, captchaId string, captchaCode string) (err error) {
	dbAccount, err := s.accountDB.Account.GetQuery().PhoneNumber_Equal(phone).QueryOne(ctx, nil)
	if err != nil {
		return err
	}

	if scene == models.SmsSceneLogin {
		if dbAccount == nil {
			return errors.NotFound("帐号不存在")
		}
	} else if scene == models.SmsSceneSignup {
		if dbAccount != nil {
			return errors.AlreadyExists("手机号已注册")
		}
	} else if scene == models.SmsSceneResetPassword {
		if dbAccount == nil {
			return errors.NotFound("帐号不存在")
		}
	} else {
		return errors.InvalidParam("验证码场景错误")
	}

	dbSmsCode := &account_db.SmsCode{}
	dbSmsCode.SceneType = scene
	dbSmsCode.PhoneNumber = phone
	dbSmsCode.SmsCode = "1234"
	dbSmsCode.CreateTime = time.Now()
	_, err = s.accountDB.SmsCode.Insert(ctx, nil, dbSmsCode)
	if err != nil {
		return err
	}

	return nil
}
