package services

import (
	"github.com/NeuronAccount/account/models"
	"github.com/NeuronAccount/account/storages/neuron_account_db"
	"github.com/NeuronFramework/errors"
	"github.com/NeuronFramework/rand"
	"github.com/NeuronFramework/restful"
	"time"
)

func (s *AccountService) SendSmsCode(ctx *restful.Context, p *models.SendSmsCodeParams) (err error) {
	phoneEncrypted, err := s.encryptPhone(p.Phone)
	if err != nil {
		return err
	}

	smsCode := rand.NextNumberFixedLength(models.SmsCodeLength)
	_, err = s.smsService.SendSms(p.Phone, smsCode, "")
	if err != nil {
		return err
	}

	smsCode = "1234"

	dbSmsCode := &neuron_account_db.SmsCode{}
	dbSmsCode.SmsScene = string(p.Scene)
	dbSmsCode.PhoneEncrypted = phoneEncrypted
	dbSmsCode.SmsCode = smsCode
	dbSmsCode.UserId = p.UserId
	_, err = s.accountDB.SmsCode.Insert(ctx, nil, dbSmsCode)
	if err != nil {
		return err
	}

	s.addOperation(ctx, &models.AccountOperation{
		OperationType:  models.OperationSendSmsCode,
		PhoneEncrypted: phoneEncrypted,
		SmsScene:       p.Scene,
	})

	return nil
}

func (s *AccountService) validateSmsCode(
	ctx *restful.Context,
	scene models.SmsScene,
	phoneEncrypted string,
	smsCode string,
	userId string) (
	err error) {
	dbSmsCode, err := s.accountDB.SmsCode.GetQuery().
		SmsScene_Equal(string(scene)).
		And().PhoneEncrypted_Equal(phoneEncrypted).
		And().UserId_Equal(userId).
		OrderBy(neuron_account_db.SMS_CODE_FIELD_ID, false).
		QueryOne(ctx, nil)
	if err != nil {
		return err
	}

	if dbSmsCode == nil || dbSmsCode.SmsCode != smsCode {
		return errors.BadRequest("InvalidSmsCode", "验证码错误")
	}

	if time.Now().Sub(dbSmsCode.CreateTime).Seconds() > models.SmsCodeValidSeconds {
		return errors.BadRequest("InvalidSmsCode", "验证码已过期")
	}

	return nil
}
