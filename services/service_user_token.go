package services

import (
	"github.com/NeuronAccount/account/models"
	"github.com/NeuronAccount/account/storages/neuron_account_db"
	"github.com/NeuronFramework/rand"
	"github.com/NeuronFramework/rest"
	"github.com/NeuronFramework/sql/wrap"
	"github.com/dgrijalva/jwt-go"
	"go.uber.org/zap"
	"time"
)

const CreateRefreshTokenMaxRetry = 10

// todo use redis or nothing
func (s *AccountService) createAccessToken(ctx *rest.Context, userId string) (accessToken string, err error) {
	//生成AccessToken
	expiresTime := time.Now().Add(time.Second * models.UserAccessTokenExpireSeconds)
	jwtToken := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.StandardClaims{
		Subject:   userId,
		ExpiresAt: expiresTime.Unix(),
		Id:        rand.NextNumberFixedLength(16), //防重，ExpiresAt精确到秒
	})
	accessToken, err = jwtToken.SignedString([]byte("0123456789"))
	if err != nil {
		return "", err
	}
	dbAccessToken := &neuron_account_db.AccessToken{}
	dbAccessToken.UserId = userId
	dbAccessToken.AccessToken = accessToken
	_, err = s.accountDB.AccessToken.Query().Insert(ctx, nil, dbAccessToken)
	if err != nil {
		return "", err
	}

	return accessToken, nil
}

func (s *AccountService) createRefreshToken(ctx *rest.Context, userId string) (refreshToken string, err error) {

	for i := 0; i < CreateRefreshTokenMaxRetry; i++ {
		//生成新token
		refreshToken := rand.NextHex(16)

		//直接更新，如果影响行数为0再插入
		result, err := s.accountDB.RefreshToken.Query().UserIdEqual(userId).
			SetRefreshToken(refreshToken).Update(ctx, nil)
		if err != nil {
			//token重复，重新生成
			if err == wrap.ErrDuplicated {
				s.logger.Warn("createRefreshToken Update ErrDuplicated",
					zap.String("refreshToken", refreshToken),
					zap.String("userId", userId))
				continue
			}

			return "", nil
		}

		//获取影响的行数
		affectRows, err := result.RowsAffected()
		if err != nil {
			return "", err
		}
		//未更新，插入新纪录
		if affectRows == 0 {
			dbRefreshToken := &neuron_account_db.RefreshToken{}
			dbRefreshToken.UserId = userId
			dbRefreshToken.RefreshToken = refreshToken
			_, err = s.accountDB.RefreshToken.Query().Insert(ctx, nil, dbRefreshToken)
			if err != nil {
				//UserId重复，到第一步直接更新
				//RefreshToken重复，到第一步直接更新，影响行数为0继续插入
				if err == wrap.ErrDuplicated {
					s.logger.Warn("createRefreshToken Insert ErrDuplicated",
						zap.String("refreshToken", refreshToken),
						zap.String("userId", userId))
					continue
				}

				return "", err
			}

		}

		return refreshToken, nil
	}

	return "", rest.Unknown("服务器正忙，请稍后再试")
}

func (s *AccountService) createUserToken(ctx *rest.Context, userId string) (userToken *models.UserToken, err error) {

	//创建AccessToken
	accessToken, err := s.createAccessToken(ctx, userId)
	if err != nil {
		return nil, err
	}

	//创建ARefreshToken
	refreshToken, err := s.createRefreshToken(ctx, userId)
	if err != nil {
		return nil, err
	}

	return &models.UserToken{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
	}, nil
}
