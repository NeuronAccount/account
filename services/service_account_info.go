package services

import (
	"github.com/NeuronAccount/account/models"
	"github.com/NeuronAccount/account/storages/neuron_account_db"
	"github.com/NeuronFramework/errors"
	"github.com/NeuronFramework/rest"
	"go.uber.org/zap"
)

func (s *AccountService) GetAccountInfo(ctx *rest.Context, userId string) (accountInfo *models.AccountInfo, err error) {
	accountInfo = &models.AccountInfo{UserId: userId}

	//基本信息
	dbUserInfo, err := s.accountDB.UserInfo.GetQuery().UserId_Equal(userId).QueryOne(ctx, nil)
	if err != nil {
		return nil, err
	}
	if dbUserInfo == nil {
		return nil, errors.NotFound("用户不存在")
	}
	accountInfo.UserName = dbUserInfo.UserName
	accountInfo.UserIcon = dbUserInfo.UserIcon

	//手机绑定
	dbPhoneAccount, err := s.accountDB.PhoneAccount.GetQuery().UserId_Equal(userId).QueryOne(ctx, nil)
	if err != nil {
		return nil, err
	}
	if dbPhoneAccount != nil {
		accountInfo.PhoneBinded = s.maskPhone(dbPhoneAccount.PhoneEncrypted)
	}

	//第三方帐号绑定
	dbOauthAccountList, err := s.accountDB.OauthAccount.GetQuery().UserId_Equal(userId).QueryList(ctx, nil)
	if err != nil {
		return nil, err
	}
	accountInfo.OauthBindedList = fromOauthAccountList(dbOauthAccountList)

	return accountInfo, nil
}

func (s *AccountService) ValidateBindPhone(ctx *rest.Context, userId string, phone string) (err error) {
	phoneEncrypted, err := s.encryptPhone(phone)
	if err != nil {
		return err
	}

	//手机号是否已被绑定
	dbPhoneAccount, err := s.accountDB.PhoneAccount.GetQuery().
		PhoneEncrypted_Equal(phoneEncrypted).
		QueryOne(ctx, nil)
	if err != nil {
		return err
	}
	if dbPhoneAccount == nil {
		return nil
	}

	//手机号已与当前帐号绑定
	if dbPhoneAccount.UserId == userId {
		return errors.Unknown("当前用户已绑定该手机号")
	}

	//手机号已绑定到其它用户，获取绑定到的用户信息
	dbUserInfo, err := s.accountDB.UserInfo.GetQuery().UserId_Equal(dbPhoneAccount.UserId).QueryOne(ctx, nil)
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
		"该用户已绑定到帐号"+s.maskString(dbUserInfo.UserName, 2, 2))
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

	//检测当前帐号是否已绑定手机号
	dbPhoneAccount, err := s.accountDB.PhoneAccount.GetQuery().
		UserId_Equal(userId).
		QueryOne(ctx, nil)
	if err != nil {
		return err
	}
	if dbPhoneAccount != nil {
		return errors.BadRequest("AlreadyBinded", "该用户已绑定手机")
	}

	//如果手机号绑定到其它帐号，先解除绑定
	otherUserId := ""
	dbPhoneAccountOther, err := s.accountDB.PhoneAccount.GetQuery().
		PhoneEncrypted_Equal(phoneEncrypted).
		QueryOne(ctx, nil)
	if err != nil {
		return err
	}
	if dbPhoneAccountOther != nil {
		otherUserId = dbPhoneAccount.UserId
		err = s.accountDB.PhoneAccount.Delete(ctx, nil, dbPhoneAccount.Id)
		if err != nil {
			return err
		}
	}

	//绑定当前帐号
	dbPhoneAccount = &neuron_account_db.PhoneAccount{}
	_, err = s.accountDB.PhoneAccount.Insert(ctx, nil, dbPhoneAccount)
	if err != nil {
		return err
	}

	s.addOperation(ctx, &models.AccountOperation{
		OperationType:  models.OperationBindPhone,
		UserId:         userId,
		PhoneEncrypted: phoneEncrypted,
		OtherUserId:    otherUserId,
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

	dbPhoneAccount, err := s.accountDB.PhoneAccount.GetQuery().UserId_Equal(userId).QueryOne(ctx, nil)
	if err != nil {
		return err
	}
	if dbPhoneAccount == nil {
		return errors.BadRequest("NotBinded", "该用户尚未绑定手机号")
	}
	err = s.accountDB.PhoneAccount.Delete(ctx, nil, dbPhoneAccount.Id)
	if err != nil {
		return err
	}

	s.addOperation(ctx, &models.AccountOperation{
		OperationType:  models.OperationUnbindPhone,
		UserId:         userId,
		PhoneEncrypted: phoneEncrypted,
	})

	return nil
}

func (s *AccountService) BindOauthAccount(ctx *rest.Context, userId string) (err error) {
	return nil
}

func (s *AccountService) UnbindOauthAccount(ctx *rest.Context, userId string) (err error) {
	return nil
}
