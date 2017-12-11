package services

import (
	"context"
	"fmt"
	"github.com/NeuronAccount/account/storages/account_db"
	"time"
)

func (s *AccountService) SmsCode(ctx context.Context, scene string, phone string, captchaId string, captchaCode string) (err error) {
	dbSmsCode := &account_db.SmsCode{}
	dbSmsCode.SceneType = scene
	dbSmsCode.PhoneNumber = phone
	dbSmsCode.SmsCode = "1234"
	dbSmsCode.CreateTime = time.Now()

	fmt.Println("sms code: " + dbSmsCode.SmsCode)

	_, err = s.accountDB.SmsCode.Insert(ctx, nil, dbSmsCode)
	if err != nil {
		return err
	}

	return nil
}
