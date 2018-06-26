package services

import (
	"github.com/NeuronAccount/account/models"
	"github.com/NeuronAccount/account/storages/neuron_account_db"
	"github.com/NeuronFramework/errors"
	"github.com/NeuronFramework/rest"
	"go.uber.org/zap"
)

func (s *AccountService) ValidateBindPhone(ctx *rest.Context, userId string, phone string) (err error) {
	phoneEncrypted, err := s.encryptPhone(phone)
	if err != nil {
		return err
	}

	//手机号是否已被绑定
	dbPhoneAccount, err := s.accountDB.PhoneAccount.Query().
		PhoneEncryptedEqual(phoneEncrypted).Select(ctx, nil)
	if err != nil {
		return err
	}
	if dbPhoneAccount == nil {
		//尚未绑定，返回成功
		return nil
	}

	//手机号已与当前帐号绑定
	if dbPhoneAccount.UserId == userId {
		return errors.Unknown("当前用户已绑定该手机号")
	}

	//手机号已绑定到其它用户，获取绑定到的用户信息
	dbUserInfo, err := s.accountDB.UserInfo.Query().
		UserIdEqual(dbPhoneAccount.UserId).Select(ctx, nil)
	if err != nil {
		return err
	}
	if dbUserInfo == nil {
		s.logger.Error("ValidateBindPhone",
			zap.String("phone", phone),
			zap.String("userId", dbPhoneAccount.UserId))
		return errors.Unknown("该手机已绑定到其它帐号")
	}

	return errors.BadRequest(
		"AlreadyBinded",
		"该手机号已绑定到帐号"+s.maskString(dbUserInfo.UserName, 2, 2))
}

func (s *AccountService) BindPhone(ctx *rest.Context, userId string, phone string, smsCode string) (err error) {
	phoneEncrypted, err := s.encryptPhone(phone)
	if err != nil {
		return err
	}

	//校验手机验证码
	err = s.validateSmsCode(ctx, models.SmsSceneBindPhone, phoneEncrypted, smsCode, userId)
	if err != nil {
		return err
	}

	//插入或更新绑定纪录
	dbPhoneAccount := &neuron_account_db.PhoneAccount{}
	dbPhoneAccount.PhoneEncrypted = phoneEncrypted
	dbPhoneAccount.UserId = userId
	_, err = s.accountDB.PhoneAccount.Query().Insert(ctx, nil, dbPhoneAccount)
	if err != nil {
		return err
	}

	// 增加操作日志
	s.addOperation(ctx, &models.AccountOperation{
		OperationType:  models.OperationBindPhone,
		UserId:         userId,
		PhoneEncrypted: phoneEncrypted,
		OtherUserId:    "",
	})

	return nil
}

func (s *AccountService) UnbindPhone(ctx *rest.Context, userId string, phone string, smsCode string) (err error) {
	phoneEncrypted, err := s.encryptPhone(phone)
	if err != nil {
		return err
	}

	//校验手机验证码
	err = s.validateSmsCode(ctx, models.SmsSceneUnbindPhone, phoneEncrypted, smsCode, userId)
	if err != nil {
		return err
	}

	//删除匹配的数据
	_, err = s.accountDB.PhoneAccount.Query().
		UserIdEqual(userId).And().PhoneEncryptedEqual(phoneEncrypted).
		Delete(ctx, nil)
	if err != nil {
		return err
	}

	//增加操作日志
	s.addOperation(ctx, &models.AccountOperation{
		OperationType:  models.OperationUnbindPhone,
		UserId:         userId,
		PhoneEncrypted: phoneEncrypted,
	})

	return nil
}
