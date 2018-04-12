package services

import (
	"github.com/NeuronAccount/account/models"
	"github.com/NeuronFramework/errors"
	"github.com/NeuronFramework/restful"
)

func (s *AccountService) GetUserInfo(ctx *restful.Context, userId string) (userInfo *models.UserInfo, err error) {
	dbUserInfo, err := s.userDB.User.GetQuery().UserId_Equal(userId).QueryOne(ctx, nil)
	if err != nil {
		return nil, err
	}
	if dbUserInfo == nil {
		return nil, errors.NotFound("用户信息不存在")
	}

	return fromUserInfo(dbUserInfo), nil
}
