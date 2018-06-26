package services

import (
	"github.com/NeuronAccount/account/models"
	"github.com/NeuronAccount/account/storages/neuron_account_db"
	"github.com/NeuronFramework/errors"
	"github.com/NeuronFramework/rest"
	"github.com/NeuronFramework/sql/wrap"
)

const CreateUserInfoMaxRetryCount = 10

func (s *AccountService) SmsLogin(ctx *rest.Context, phone string, smsCode string) (
	userToken *models.UserToken, err error) {

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
	dbPhoneAccount, err := s.accountDB.PhoneAccount.Query().PhoneEncryptedEqual(phoneEncrypted).Select(ctx, nil)
	if err != nil {
		return nil, err
	}
	if dbPhoneAccount == nil {
		err = s.accountDB.TransactionReadCommitted(ctx, false, func(tx *wrap.Tx) (err error) {
			dbUserInfo, err := s.retryCreateUserInfo(ctx, tx, CreateUserInfoMaxRetryCount)
			if err != nil {
				return err
			}

			dbPhoneAccount = &neuron_account_db.PhoneAccount{}
			dbPhoneAccount.PhoneEncrypted = phoneEncrypted
			dbPhoneAccount.UserId = dbUserInfo.UserId
			_, err = s.accountDB.PhoneAccount.Query().Insert(ctx, tx, dbPhoneAccount)
			if err != nil {
				//手机号已被注册的情况不需要单独判断
				//因为其UserID不太可能和新建的UserID相同，故直接返回错误
				return err
			}

			return nil
		})
		if err != nil {
			return nil, err
		}
	}

	//创建token
	userToken, err = s.createUserToken(ctx, dbPhoneAccount.UserId)
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

func (s *AccountService) PhonePasswordLogin(ctx *rest.Context, phone string, passwordHash1 string) (
	userToken *models.UserToken, err error) {

	//加密手机号
	phoneEncrypted, err := s.encryptPhone(phone)
	if err != nil {
		return nil, err
	}

	passwordHash2, err := s.calcPasswordHash(passwordHash1)
	if err != nil {
		return nil, err
	}

	//获取手机帐号
	dbPhoneAccount, err := s.accountDB.PhoneAccount.Query().PhoneEncryptedEqual(phoneEncrypted).Select(ctx, nil)
	if err != nil {
		return nil, err
	}
	if dbPhoneAccount == nil {
		return nil, errors.NotFound("手机号尚未注册")
	}

	//获取帐号信息
	dbUserInfo, err := s.accountDB.UserInfo.Query().UserIdEqual(dbPhoneAccount.UserId).Select(ctx, nil)
	if err != nil {
		return nil, err
	}
	if dbPhoneAccount == nil {
		return nil, errors.NotFound("内部错误，手机帐号不存在")
	}

	if dbUserInfo.PasswordHash != passwordHash2 {
		return nil, errors.BadRequest("AuthorizationFailed", "密码不正确")
	}

	//创建token
	userToken, err = s.createUserToken(ctx, dbPhoneAccount.UserId)
	if err != nil {
		return nil, err
	}

	//操作纪录
	s.addOperation(ctx, &models.AccountOperation{
		OperationType:  models.OperationPhonePasswordLogin,
		UserId:         dbPhoneAccount.UserId,
		PhoneEncrypted: phoneEncrypted,
	})

	return nil, nil
}
