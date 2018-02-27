package services

import (
	"context"
	"github.com/NeuronAccount/account/models"
	"github.com/NeuronAccount/account/storages/account_db"
	"github.com/NeuronFramework/errors"
	"github.com/NeuronFramework/log"
	"github.com/dgrijalva/jwt-go"
	"go.uber.org/zap"
	"time"
)

type AccountServiceOptions struct {
}

type AccountService struct {
	logger    *zap.Logger
	options   *AccountServiceOptions
	accountDB *account_db.DB
}

func NewAccountService(options *AccountServiceOptions) (s *AccountService, err error) {
	s = &AccountService{}
	s.logger = log.TypedLogger(s)
	s.options = options
	s.accountDB, err = account_db.NewDB()
	if err != nil {
		return nil, err
	}

	return s, nil
}

func (s *AccountService) validateSmsCode(
	ctx context.Context, scene string, phone string, smsCode string) (err error) {
	dbSmsCode, err := s.accountDB.SmsCode.GetQuery().
		SceneType_Equal(scene).And().PhoneNumber_Equal(phone).
		OrderBy(account_db.SMS_CODE_FIELD_CREATE_TIME, false).QueryOne(ctx, nil)
	if err != nil {
		return err
	}

	if dbSmsCode == nil || dbSmsCode.SmsCode != smsCode {
		return errors.BadRequest("InvalidSmsCode", "验证码错误")
	}

	if time.Now().Sub(dbSmsCode.CreateTime).Seconds() > models.SmsCodeValidSeconds {
		return errors.BadRequest("InvalidSmsCode", "验证码已过期")
	}

	return nil
}

func (s *AccountService) generateJwt(accountId string) (tokenString string, err error) {
	expiresTime := time.Now().Add(time.Hour)
	userToken := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.StandardClaims{
		Subject:   accountId,
		ExpiresAt: expiresTime.Unix(),
	})
	tokenString, err = userToken.SignedString([]byte("0123456789"))
	if err != nil {
		return "", err
	}

	return tokenString, nil
}

func (s *AccountService) calcPasswordHash(password string) (passwordHash string) {
	return password
}
