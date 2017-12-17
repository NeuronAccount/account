package services

import (
	"context"
	"github.com/NeuronAccount/account/storages/account_db"
	"time"
)

func (s *AccountService) SmsCode(ctx context.Context, scene string, phone string, captchaId string, captchaCode string) (err error) {
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
