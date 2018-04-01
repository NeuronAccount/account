package services

import (
	"github.com/NeuronAccount/account/models"
	"github.com/NeuronAccount/account/storages/account_db"
	"github.com/NeuronFramework/errors"
	"github.com/NeuronFramework/rand"
	"github.com/NeuronFramework/restful"
	"time"
)

func (s *AccountService) SmsCode(ctx *restful.Context, scene string, phone string, captchaId string, captchaCode string) (err error) {
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

	smsCode := rand.NextNumberFixedLength(6)
	_, err = s.smsService.SendSms(phone, smsCode, "")
	if err != nil {
		return err
	}

	dbSmsCode := &account_db.SmsCode{}
	dbSmsCode.SceneType = scene
	dbSmsCode.PhoneNumber = phone
	dbSmsCode.SmsCode = smsCode
	dbSmsCode.CreateTime = time.Now()
	_, err = s.accountDB.SmsCode.Insert(ctx, nil, dbSmsCode)
	if err != nil {
		return err
	}

	op := &models.Operation{
		OperationType: models.OperationSmsCode,
		SmsScene:      scene,
		Phone:         phone,
	}
	if dbAccount != nil {
		op.AccountID = dbAccount.AccountId
	}
	s.addOperation(ctx, op)

	return nil
}
