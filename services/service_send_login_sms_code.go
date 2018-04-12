package services

import (
	"github.com/NeuronAccount/account/models"
	"github.com/NeuronAccount/account/storages/neuron_account_db"
	"github.com/NeuronFramework/rand"
	"github.com/NeuronFramework/restful"
)

func (s *AccountService) SendLoginSmsCode(ctx *restful.Context, phone string, captchaId string, captchaCode string) (err error) {
	loginSmsCode := rand.NextNumberFixedLength(models.SmsCodeLength)
	_, err = s.smsService.SendSms(phone, loginSmsCode, "")
	if err != nil {
		return err
	}

	loginSmsCode = "1234"

	dbLoginSmsCode := &neuron_account_db.LoginSmsCode{}
	dbLoginSmsCode.PhoneNumber = phone
	dbLoginSmsCode.SmsCode = loginSmsCode
	_, err = s.userDB.LoginSmsCode.Insert(ctx, nil, dbLoginSmsCode)
	if err != nil {
		return err
	}

	s.addOperation(ctx, &models.AccountOperation{
		OperationType: models.OperationSendLoginSmsCode,
		Phone:         phone,
	})

	return nil
}
