package handler

import (
	"github.com/NeuronFramework/log"
	"github.com/NeuronFramework/restful"
	"github.com/NeuronGroup/Account/api/private/gen/models"
	"github.com/NeuronGroup/Account/api/private/gen/restapi/operations"
	"github.com/NeuronGroup/Account/services"
	"github.com/go-openapi/runtime/middleware"
	"go.uber.org/zap"
)

type AccountHandlerOptions struct {
	AccountStorageConnectionString string
	OAuthConnectionString          string
}

type AccountHandler struct {
	logger  *zap.Logger
	options *AccountHandlerOptions
	service *services.AccountService
}

func NewAccountHandler(options *AccountHandlerOptions) (h *AccountHandler, err error) {
	h = &AccountHandler{}
	h.logger = log.TypedLogger(h)
	h.options = options
	h.service, err = services.NewAccountService(&services.AccountServiceOptions{
		AccountStorageConnectionString: options.AccountStorageConnectionString,
		OAuthConnectionString:          options.OAuthConnectionString,
	})
	if err != nil {
		return nil, err
	}

	return h, nil
}

func (h AccountHandler) SmsCode(p operations.SmsCodeParams) middleware.Responder {
	captchaId := ""
	if p.CaptchaID != nil {
		captchaId = *p.CaptchaID
	}

	captchaCode := ""
	if p.CaptchaCode != nil {
		captchaCode = *p.CaptchaCode
	}

	err := h.service.SmsCode(p.Scene, p.Phone, captchaId, captchaCode)
	if err != nil {
		return restful.Responder(err)
	}

	return operations.NewSmsCodeOK()
}

func (h AccountHandler) SmsSignup(p operations.SmsSignupParams) middleware.Responder {
	jwt, err := h.service.SmsSignup(p.Phone, p.SmsCode, p.Password, toOAuth2Param(p.Oauth2AuthorizeParams))
	if err != nil {
		return restful.Responder(err)
	}

	return operations.NewSmsSignupOK().WithPayload(&models.LoginResponse{Jwt: jwt})
}

func (h AccountHandler) SmsLogin(p operations.SmsLoginParams) middleware.Responder {
	jwt, err := h.service.SmsLogin(p.Phone, p.SmsCode, toOAuth2Param(p.Oauth2AuthorizeParams))
	if err != nil {
		return restful.Responder(err)
	}

	return operations.NewSmsLoginOK().WithPayload(&models.LoginResponse{Jwt: jwt})
}

func (h AccountHandler) Login(p operations.LoginParams) middleware.Responder {
	jwt, err := h.service.Login(p.Name, p.Password, toOAuth2Param(p.Oauth2AuthorizeParams))
	if err != nil {
		return restful.Responder(err)
	}

	return operations.NewLoginOK().WithPayload(&models.LoginResponse{Jwt: jwt})
}

func (h AccountHandler) Logout(p operations.LogoutParams) middleware.Responder {
	err := h.service.Logout(p.Jwt)
	if err != nil {
		return restful.Responder(err)
	}

	return operations.NewLogoutOK()
}
