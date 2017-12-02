package handler

import (
	"github.com/NeuronAccount/account/api/private/gen/restapi/operations"
	"github.com/NeuronAccount/account/services"
	"github.com/NeuronFramework/log"
	"github.com/NeuronFramework/restful"
	"github.com/go-openapi/runtime/middleware"
	"go.uber.org/zap"
)

type AccountHandlerOptions struct {
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
	h.service, err = services.NewAccountService(&services.AccountServiceOptions{})
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
	jwt, err := h.service.SmsSignup(p.Phone, p.SmsCode, p.Password)
	if err != nil {
		return restful.Responder(err)
	}

	return operations.NewSmsSignupOK().WithPayload(jwt)
}

func (h AccountHandler) SmsLogin(p operations.SmsLoginParams) middleware.Responder {
	jwt, err := h.service.SmsLogin(p.Phone, p.SmsCode)
	if err != nil {
		return restful.Responder(err)
	}

	return operations.NewSmsLoginOK().WithPayload(jwt)
}

func (h AccountHandler) Login(p operations.LoginParams) middleware.Responder {
	jwt, err := h.service.Login(p.Name, p.Password)
	if err != nil {
		return restful.Responder(err)
	}

	return operations.NewLoginOK().WithPayload(jwt)
}

func (h AccountHandler) Logout(p operations.LogoutParams) middleware.Responder {
	err := h.service.Logout(p.Jwt)
	if err != nil {
		return restful.Responder(err)
	}

	return operations.NewLogoutOK()
}
