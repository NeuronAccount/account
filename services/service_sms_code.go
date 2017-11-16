package services

import (
	"context"
	"fmt"
	"github.com/NeuronGroup/Accounts/storages/account"
	"time"
)

func (s *AccountService) SmsCode(scene string, phone string, captchaId string, captchaCode string) (err error) {
	dbSmsCode := &account.SmsCode{}
	dbSmsCode.SceneType = scene
	dbSmsCode.PhoneNumber = phone
	dbSmsCode.SmsCode = "1234"
	dbSmsCode.CreateTime = time.Now()

	fmt.Println("sms code: " + dbSmsCode.SmsCode)

	_, err = s.db.SmsCode.Insert(context.Background(), nil, dbSmsCode)
	if err != nil {
		return err
	}

	return nil
}
