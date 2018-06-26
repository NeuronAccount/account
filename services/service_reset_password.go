package services

import (
	"fmt"
	"github.com/NeuronAccount/account/models"
	"github.com/NeuronFramework/errors"
	"github.com/NeuronFramework/rest"
)

func (s *AccountService) ResetPassword(ctx *rest.Context, phone string, smsCode string, newPasswordHash1 string) (err error) {

	//加密手机号
	phoneEncrypted, err := s.encryptPhone(phone)
	if err != nil {
		return err
	}

	//hash密码
	passwordHash2, err := s.calcPasswordHash(newPasswordHash1)
	if err != nil {
		return err
	}

	//校验验证码
	err = s.validateSmsCode(ctx, models.SmsSceneResetPassword, phoneEncrypted, smsCode, "")
	if err != nil {
		return err
	}

	//获取手机号的userId
	dbPhoneAccount, err := s.accountDB.PhoneAccount.Query().PhoneEncryptedEqual(phoneEncrypted).Select(ctx, nil)
	if err != nil {
		return err
	}
	if dbPhoneAccount == nil {
		return errors.NotFound("该手机号尚未注册")
	}

	//更新密码
	result, err := s.accountDB.UserInfo.Query().UserIdEqual(dbPhoneAccount.UserId).
		SetPasswordHash(passwordHash2).Update(ctx, nil)
	if err != nil {
		return err
	}
	affectedRows, err := result.RowsAffected()
	if err != nil {
		return err
	}
	if affectedRows != 1 {
		return errors.Unknown(fmt.Sprintf("更新失败，影响行数%d", affectedRows))
	}

	//操作纪录
	s.addOperation(ctx, &models.AccountOperation{
		OperationType:  models.OperationResetPassword,
		UserId:         dbPhoneAccount.UserId,
		PhoneEncrypted: phoneEncrypted,
	})

	return nil
}
