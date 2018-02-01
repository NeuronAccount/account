package services

import (
	"context"
	"fmt"
	"github.com/NeuronAccount/account/models"
	"github.com/NeuronFramework/errors"
)

func (s *AccountService) ResetPassword(ctx context.Context, phone string, smsCode string, newPassword string) (err error) {
	err = s.validateSmsCode(ctx, models.SmsSceneResetPassword, phone, smsCode)
	if err != nil {
		return err
	}

	fmt.Println("aaa")

	dbAccount, err := s.accountDB.Account.GetQuery().PhoneNumber_Equal(phone).QueryOne(ctx, nil)
	if err != nil {
		return nil
	}

	if dbAccount == nil {
		return errors.NotFound("帐号不存在")
	}

	dbAccount.PasswordHash = s.calcPasswordHash(newPassword)
	err = s.accountDB.Account.Update(ctx, nil, dbAccount)
	if err != nil {
		return err
	}

	return nil
}
