package services

import (
	"github.com/NeuronAccount/account/models"
	"github.com/NeuronAccount/account/storages/user_db"
	"github.com/NeuronFramework/rand"
	"github.com/NeuronFramework/restful"
)

func (s *AccountService) SendLoginSmsCode(
	ctx *restful.Context,
	phone string,
	captchaId string,
	captchaCode string) (err error) {
	loginSmsCode := rand.NextNumberFixedLength(models.SmsCodeLength)
	_, err = s.smsService.SendSms(phone, loginSmsCode, "")
	if err != nil {
		return err
	}

	dbLoginSmsCode := &user_db.LoginSmsCode{}
	dbLoginSmsCode.PhoneNumber = phone
	dbLoginSmsCode.SmsCode = loginSmsCode
	_, err = s.userDB.LoginSmsCode.Insert(ctx, nil, dbLoginSmsCode)
	if err != nil {
		return err
	}

	return nil
}
