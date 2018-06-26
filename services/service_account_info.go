package services

import (
	"github.com/NeuronAccount/account/models"
	"github.com/NeuronFramework/errors"
	"github.com/NeuronFramework/rest"
)

func (s *AccountService) GetAccountInfo(ctx *rest.Context, userId string) (accountInfo *models.AccountInfo, err error) {
	accountInfo = &models.AccountInfo{UserId: userId}

	//基本信息
	dbUserInfo, err := s.accountDB.UserInfo.Query().
		UserIdEqual(userId).Select(ctx, nil)
	if err != nil {
		return nil, err
	}
	if dbUserInfo == nil {
		return nil, errors.NotFound("用户不存在")
	}
	accountInfo.UserName = dbUserInfo.UserName
	accountInfo.UserIcon = dbUserInfo.UserIcon

	//手机绑定
	dbPhoneAccount, err := s.accountDB.PhoneAccount.Query().
		UserIdEqual(userId).Select(ctx, nil)
	if err != nil {
		return nil, err
	}
	if dbPhoneAccount != nil {
		accountInfo.PhoneBinded = s.maskPhone(dbPhoneAccount.PhoneEncrypted)
	}

	//第三方帐号绑定
	dbOauthAccountList, err := s.accountDB.OauthAccount.Query().
		UserIdEqual(userId).SelectList(ctx, nil)
	if err != nil {
		return nil, err
	}
	accountInfo.OauthBindedList = fromOauthAccountList(dbOauthAccountList)

	return accountInfo, nil
}
