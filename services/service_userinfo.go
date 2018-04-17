package services

import (
	"github.com/NeuronAccount/account/models"
	"github.com/NeuronFramework/errors"
	"github.com/NeuronFramework/rest"
)

func (s *AccountService) GetUserInfo(ctx *rest.Context, userId string) (userInfo *models.UserInfo, err error) {
	dbUserInfo, err := s.accountDB.UserInfo.GetQuery().UserId_Equal(userId).QueryOne(ctx, nil)
	if err != nil {
		return nil, err
	}
	if dbUserInfo == nil {
		return nil, errors.NotFound("用户不存在")
	}

	return fromUserInfo(dbUserInfo), nil
}

func (s *AccountService) SetUserName(ctx *rest.Context, userId string, userName string) (err error) {
	//检查名称是否已被使用
	dbOtherUserInfo, err := s.accountDB.UserInfo.GetQuery().UserName_Equal(userName).QueryOne(ctx, nil)
	if err != nil {
		return err
	}
	if dbOtherUserInfo != nil && dbOtherUserInfo.UserId != userId {
		return errors.NotFound("该名称已被使用")
	}

	//更新
	dbUserInfo, err := s.accountDB.UserInfo.GetQuery().UserId_Equal(userId).QueryOne(ctx, nil)
	if err != nil {
		return err
	}
	if dbUserInfo == nil {
		return errors.NotFound("用户不存在")
	}

	err = s.accountDB.UserInfo.GetUpdate().UserName(userName).Update(ctx, nil, dbUserInfo.Id)
	if err != nil {
		return err
	}

	return nil
}

func (s *AccountService) SetUserIcon(ctx *rest.Context, userId string, userIcon string) (err error) {
	dbUserInfo, err := s.accountDB.UserInfo.GetQuery().UserId_Equal(userId).QueryOne(ctx, nil)
	if err != nil {
		return err
	}
	if dbUserInfo == nil {
		return errors.NotFound("用户不存在")
	}

	err = s.accountDB.UserInfo.GetUpdate().UserIcon(userIcon).Update(ctx, nil, dbUserInfo.Id)
	if err != nil {
		return err
	}

	return nil
}
