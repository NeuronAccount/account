package services

import (
	"github.com/NeuronAccount/account/models"
	"github.com/NeuronAccount/account/storages/user_db"
	"github.com/NeuronFramework/errors"
	"github.com/NeuronFramework/rand"
	"github.com/NeuronFramework/restful"
	"github.com/dgrijalva/jwt-go"
	"time"
)

func (s *AccountService) RefreshToken(ctx *restful.Context, refreshToken string) (userToken *models.UserToken, err error) {
	//检查RefreshToken是否最新且有效
	dbRefreshToken, err := s.userDB.RefreshToken.GetQuery().
		OrderBy(user_db.REFRESH_TOKEN_FIELD_ID, false).
		Limit(0, 1).
		QueryOne(ctx, nil)
	if err != nil {
		return nil, err
	}
	if dbRefreshToken == nil {
		return nil, errors.NotFound("Token已失效，请重新登录")
	}
	if dbRefreshToken.IsLogout == 1 {
		return nil, errors.NotFound("Token已失效，请重新登录")
	}
	if dbRefreshToken.RefreshToken != refreshToken {
		dbRefreshTokenOld, err := s.userDB.RefreshToken.GetQuery().
			RefreshToken_Equal(refreshToken).QueryOne(ctx, nil)
		if err != nil {
			return nil, err
		}

		if dbRefreshTokenOld != nil {
			return nil, errors.NotFound("您已在其它地方登录，请重新登录")
		}

		return nil, errors.NotFound("Token已失效，请重新登录")
	}

	//生成AccessToken
	expiresTime := time.Now().Add(time.Second * models.UserAccessTokenExpireSeconds)
	jwtToken := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.StandardClaims{
		Subject:   dbRefreshToken.UserId,
		ExpiresAt: expiresTime.Unix(),
	})
	accessToken, err := jwtToken.SignedString([]byte("0123456789"))
	if err != nil {
		return nil, err
	}
	dbAccessToken := &user_db.AccessToken{}
	dbAccessToken.UserId = dbRefreshToken.UserId
	dbAccessToken.AccessToken = accessToken
	_, err = s.userDB.AccessToken.Insert(ctx, nil, dbAccessToken)
	if err != nil {
		return nil, err
	}

	//生成RefreshToken
	dbRefreshToken.RefreshToken = rand.NextHex(16)
	err = s.userDB.RefreshToken.Update(ctx, nil, dbRefreshToken)
	if err != nil {
		return nil, err
	}

	return &models.UserToken{
		AccessToken:  accessToken,
		RefreshToken: dbRefreshToken.RefreshToken,
	}, nil
}
