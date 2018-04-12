package services

import (
	"github.com/NeuronAccount/account/models"
	"github.com/NeuronAccount/account/storages/neuron_account_db"
	"github.com/NeuronFramework/errors"
	"github.com/NeuronFramework/rand"
	"github.com/NeuronFramework/restful"
	"time"
)

func (s *AccountService) SmsLogin(ctx *restful.Context, phone string, smsCode string) (userToken *models.UserToken, err error) {
	//校验验证码
	dbLoginSmsCode, err := s.userDB.LoginSmsCode.GetQuery().
		PhoneNumber_Equal(phone).
		OrderBy(neuron_account_db.LOGIN_SMS_CODE_FIELD_ID, false).
		QueryOne(ctx, nil)
	if err != nil {
		return nil, err
	}

	if dbLoginSmsCode == nil || dbLoginSmsCode.SmsCode != smsCode {
		return nil, errors.BadRequest("InvalidSmsCode", "验证码错误")
	}

	if time.Now().Sub(dbLoginSmsCode.CreateTime).Seconds() > models.SmsCodeValidSeconds {
		return nil, errors.BadRequest("InvalidSmsCode", "验证码已过期")
	}

	//获取帐号，如果不存在，新建一个
	dbPhoneAccount, err := s.userDB.PhoneAccount.GetQuery().PhoneNumber_Equal(phone).QueryOne(ctx, nil)
	if err != nil {
		return nil, err
	}
	if dbPhoneAccount == nil {
		dbUser := &neuron_account_db.User{}
		dbUser.UserId = rand.NextHex(16)
		dbUser.UserName = "用户" + rand.NextNumberFixedLength(8)
		dbUser.UserIcon = ""
		_, err = s.userDB.User.Insert(ctx, nil, dbUser)
		if err != nil {
			return nil, err
		}

		dbPhoneAccount = &neuron_account_db.PhoneAccount{}
		dbPhoneAccount.PhoneNumber = phone
		dbPhoneAccount.UserId = dbUser.UserId
		_, err = s.userDB.PhoneAccount.Insert(ctx, nil, dbPhoneAccount)
		if err != nil {
			return nil, err
		}
	}

	//生成AccessToken
	accessToken, err := s.generateJwt(dbPhoneAccount.UserId)
	if err != nil {
		return nil, err
	}
	dbAccessToken := &neuron_account_db.AccessToken{}
	dbAccessToken.UserId = dbPhoneAccount.UserId
	dbAccessToken.AccessToken = accessToken
	_, err = s.userDB.AccessToken.Insert(ctx, nil, dbAccessToken)
	if err != nil {
		return nil, err
	}

	//生成RefreshToken
	dbRefreshToken := &neuron_account_db.RefreshToken{}
	dbRefreshToken.UserId = dbPhoneAccount.UserId
	dbRefreshToken.RefreshToken = rand.NextHex(16)
	dbRefreshToken.IsLogout = 0
	dbRefreshToken.LogoutTime = time.Now()
	_, err = s.userDB.RefreshToken.Insert(ctx, nil, dbRefreshToken)
	if err != nil {
		return nil, err
	}

	//操作纪录
	s.addOperation(ctx, &models.AccountOperation{
		OperationType: models.OperationSmsLogin,
		UserId:        dbPhoneAccount.UserId,
		Phone:         phone,
	})

	return &models.UserToken{
		AccessToken:  accessToken,
		RefreshToken: dbRefreshToken.RefreshToken,
	}, nil
}
