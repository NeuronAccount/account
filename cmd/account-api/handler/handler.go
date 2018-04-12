package handler

import (
	"github.com/NeuronAccount/account/api/gen/restapi/operations"
	"github.com/NeuronAccount/account/services"
	"github.com/NeuronFramework/errors"
	"github.com/NeuronFramework/log"
	"github.com/NeuronFramework/restful"
	"github.com/dgrijalva/jwt-go"
	"github.com/go-openapi/runtime/middleware"
	"go.uber.org/zap"
)

type AccountHandler struct {
	logger  *zap.Logger
	service *services.AccountService
}

func NewAccountHandler() (h *AccountHandler, err error) {
	h = &AccountHandler{}
	h.logger = log.TypedLogger(h)
	h.service, err = services.NewAccountService(&services.AccountServiceOptions{})
	if err != nil {
		return nil, err
	}

	return h, nil
}

func (h *AccountHandler) BearerAuth(token string) (userId interface{}, err error) {
	claims := jwt.StandardClaims{}
	_, err = jwt.ParseWithClaims(token, &claims, func(t *jwt.Token) (interface{}, error) {
		return []byte("0123456789"), nil
	})
	if err != nil {
		return nil, err
	}

	if claims.Subject == "" {
		return nil, errors.Unknown("验证失败： claims.Subject nil")
	}

	return claims.Subject, nil
}

func (h *AccountHandler) SendLoginSmsCode(p operations.SendLoginSmsCodeParams) middleware.Responder {
	err := h.service.SendLoginSmsCode(restful.NewContext(p.HTTPRequest), p.Phone, p.CaptchaID, p.CaptchaCode)
	if err != nil {
		return errors.Wrap(err)
	}

	return operations.NewSendLoginSmsCodeOK()
}

func (h *AccountHandler) SmsLogin(p operations.SmsLoginParams) middleware.Responder {
	userToken, err := h.service.SmsLogin(restful.NewContext(p.HTTPRequest), p.Phone, p.SmsCode)
	if err != nil {
		return errors.Wrap(err)
	}

	return operations.NewSmsLoginOK().WithPayload(fromToken(userToken))
}

func (h *AccountHandler) Logout(p operations.LogoutParams) middleware.Responder {
	err := h.service.Logout(restful.NewContext(p.HTTPRequest), p.AccessToken, p.RefreshToken)
	if err != nil {
		return errors.Wrap(err)
	}

	return operations.NewLogoutOK()
}

func (h *AccountHandler) RefreshToken(p operations.RefreshTokenParams) middleware.Responder {
	userToken, err := h.service.RefreshToken(restful.NewContext(p.HTTPRequest), p.RefreshToken)
	if err != nil {
		return errors.Wrap(err)
	}

	return operations.NewRefreshTokenOK().WithPayload(fromToken(userToken))
}

func (h *AccountHandler) GetUserInfo(p operations.GetUserInfoParams, userId interface{}) middleware.Responder {
	userInfo, err := h.service.GetUserInfo(restful.NewContext(p.HTTPRequest), userId.(string))
	if err != nil {
		return errors.Wrap(err)
	}

	return operations.NewGetUserInfoOK().WithPayload(fromUserInfo(userInfo))
}
