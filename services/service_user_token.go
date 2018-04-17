package services

import (
	"github.com/NeuronAccount/account/models"
	"github.com/NeuronAccount/account/storages/neuron_account_db"
	"github.com/NeuronFramework/errors"
	"github.com/NeuronFramework/rand"
	"github.com/NeuronFramework/rest"
	"github.com/dgrijalva/jwt-go"
	"time"
)

func (s *AccountService) createUserToken(
	ctx *rest.Context,
	userId string,
	dbRefreshToken *neuron_account_db.RefreshToken) (
	userToken *models.UserToken, err error) {
	//防刷
	dbRefreshTokenCount, err := s.accountDB.RefreshToken.GetQuery().
		UserId_Equal(userId).
		And().UpdateTime_Greater(time.Now().UTC().Add(-time.Minute)).
		And().IsLogout_Equal(0).QueryCount(ctx, nil)
	if err != nil {
		return nil, err
	}
	if dbRefreshTokenCount > 5 {
		return nil, errors.BadRequest("FreqLimit", "登录过于频繁，请稍后再试")
	}

	//生成AccessToken
	expiresTime := time.Now().Add(time.Second * models.UserAccessTokenExpireSeconds)
	jwtToken := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.StandardClaims{
		Subject:   userId,
		ExpiresAt: expiresTime.Unix(),
		Id:        rand.NextNumberFixedLength(16), //防重，ExpiresAt精确到秒
	})
	accessToken, err := jwtToken.SignedString([]byte("0123456789"))
	if err != nil {
		return nil, err
	}
	dbAccessToken := &neuron_account_db.AccessToken{}
	dbAccessToken.UserId = userId
	dbAccessToken.AccessToken = accessToken
	_, err = s.accountDB.AccessToken.Insert(ctx, nil, dbAccessToken)
	if err != nil {
		return nil, err
	}

	//生成RefreshToken
	refreshToken := rand.NextHex(16)
	if dbRefreshToken == nil {
		dbRefreshToken = &neuron_account_db.RefreshToken{}
		dbRefreshToken.RefreshToken = refreshToken
		dbRefreshToken.UserId = userId
		dbRefreshToken.IsLogout = 0
		dbRefreshToken.LogoutTime = time.Now()
		_, err = s.accountDB.RefreshToken.Insert(ctx, nil, dbRefreshToken)
		if err != nil {
			return nil, err
		}
	} else {
		err = s.accountDB.RefreshToken.GetUpdate().
			RefreshToken(refreshToken).
			IsLogout(0).
			Update(ctx, nil, dbRefreshToken.Id)
		if err != nil {
			return nil, err
		}
	}

	return &models.UserToken{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
	}, nil
}
