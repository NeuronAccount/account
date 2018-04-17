package services

import (
	"github.com/NeuronAccount/account/models"
	"github.com/NeuronAccount/account/storages/neuron_account_db"
	"github.com/NeuronFramework/rand"
	"github.com/NeuronFramework/rest"
)

func (s *AccountService) SmsLogin(ctx *rest.Context, phone string, smsCode string) (userToken *models.UserToken, err error) {
	//加密手机号
	phoneEncrypted, err := s.encryptPhone(phone)
	if err != nil {
		return nil, err
	}

	//校验验证码
	err = s.validateSmsCode(ctx, models.SmsSceneSmsLogin, phoneEncrypted, smsCode, "")
	if err != nil {
		return nil, err
	}

	//获取帐号，如果不存在，新建一个
	dbPhoneAccount, err := s.accountDB.PhoneAccount.GetQuery().
		PhoneEncrypted_Equal(phoneEncrypted).
		QueryOne(ctx, nil)
	if err != nil {
		return nil, err
	}
	if dbPhoneAccount == nil {
		dbUser := &neuron_account_db.UserInfo{}
		dbUser.UserId = rand.NextHex(16)
		dbUser.UserName = "用户" + rand.NextNumberFixedLength(8)
		dbUser.UserIcon = ""
		_, err = s.accountDB.UserInfo.Insert(ctx, nil, dbUser)
		if err != nil {
			return nil, err
		}

		dbPhoneAccount = &neuron_account_db.PhoneAccount{}
		dbPhoneAccount.PhoneEncrypted = phoneEncrypted
		dbPhoneAccount.UserId = dbUser.UserId
		_, err = s.accountDB.PhoneAccount.Insert(ctx, nil, dbPhoneAccount)
		if err != nil {
			return nil, err
		}

		//todo 已存在，回退
	}

	userToken, err = s.createUserToken(ctx, dbPhoneAccount.UserId, nil)
	if err != nil {
		return nil, err
	}

	//操作纪录
	s.addOperation(ctx, &models.AccountOperation{
		OperationType:  models.OperationSmsLogin,
		UserId:         dbPhoneAccount.UserId,
		PhoneEncrypted: phoneEncrypted,
	})

	return userToken, nil
}
