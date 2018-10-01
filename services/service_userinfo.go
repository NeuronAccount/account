package services

import (
	"github.com/NeuronAccount/account/models"
	"github.com/NeuronAccount/account/storages/neuron_account_db"
	"github.com/NeuronFramework/rand"
	"github.com/NeuronFramework/rest"
	"github.com/NeuronFramework/sql/wrap"
)

func (s *AccountService) createUserInfo(ctx *rest.Context, tx *wrap.Tx) (userInfo *neuron_account_db.UserInfo, err error) {
	dbUserInfo := &neuron_account_db.UserInfo{}
	dbUserInfo.UserId = rand.NextHex(16)
	dbUserInfo.UserName = "用户" + rand.NextNumberFixedLength(8)
	dbUserInfo.UserIcon = ""
	dbUserInfo.PasswordHash = ""
	_, err = s.accountDB.UserInfo.Query().Insert(ctx, tx, dbUserInfo)
	if err != nil {
		return nil, err
	}

	return dbUserInfo, nil
}

func (s *AccountService) GetUserInfo(ctx *rest.Context, userId string) (userInfo *models.UserInfo, err error) {
	dbUserInfo, err := s.accountDB.UserInfo.Query().
		UserIdEqual(userId).Select(ctx, nil)
	if err != nil {
		return nil, err
	}
	if dbUserInfo == nil {
		return nil, rest.NotFound("用户不存在")
	}

	return fromUserInfo(dbUserInfo), nil
}

func (s *AccountService) SetUserName(ctx *rest.Context, userId string, userName string) (err error) {
	//检查名称是否已被使用
	dbOtherUserInfo, err := s.accountDB.UserInfo.Query().
		UserNameEqual(userName).Select(ctx, nil)
	if err != nil {
		return err
	}
	if dbOtherUserInfo != nil && dbOtherUserInfo.UserId != userId {
		return rest.NotFound("该名称已被使用")
	}

	//更新
	dbUserInfo, err := s.accountDB.UserInfo.Query().
		UserIdEqual(userId).Select(ctx, nil)
	if err != nil {
		return err
	}
	if dbUserInfo == nil {
		return rest.NotFound("用户不存在")
	}

	_, err = s.accountDB.UserInfo.Query().IdEqual(dbUserInfo.Id).
		SetUserName(userName).Update(ctx, nil)
	if err != nil {
		return err
	}

	return nil
}

func (s *AccountService) SetUserIcon(ctx *rest.Context, userId string, userIcon string) (err error) {
	dbUserInfo, err := s.accountDB.UserInfo.Query().UserIdEqual(userId).Select(ctx, nil)
	if err != nil {
		return err
	}
	if dbUserInfo == nil {
		return rest.NotFound("用户不存在")
	}

	_, err = s.accountDB.UserInfo.Query().IdEqual(dbUserInfo.Id).
		SetUserIcon(userIcon).Update(ctx, nil)
	if err != nil {
		return err
	}

	return nil
}
