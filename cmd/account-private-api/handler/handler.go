package handler

import (
	"github.com/NeuronAccount/account/api-private/gen/restapi/operations"
	"github.com/NeuronAccount/account/services"
	"github.com/NeuronFramework/errors"
	"github.com/NeuronFramework/log"
	"github.com/NeuronFramework/restful"
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

func (h AccountHandler) SmsCode(p operations.SmsCodeParams) middleware.Responder {
	captchaId := ""
	if p.CaptchaID != nil {
		captchaId = *p.CaptchaID
	}

	captchaCode := ""
	if p.CaptchaCode != nil {
		captchaCode = *p.CaptchaCode
	}

	err := h.service.SmsCode(restful.NewContext(p.HTTPRequest), p.Scene, p.Phone, captchaId, captchaCode)
	if err != nil {
		return errors.Wrap(err)
	}

	return operations.NewSmsCodeOK()
}

func (h AccountHandler) SmsSignup(p operations.SmsSignupParams) middleware.Responder {
	jwt, err := h.service.SmsSignup(restful.NewContext(p.HTTPRequest), p.Phone, p.SmsCode, p.Password)
	if err != nil {
		return errors.Wrap(err)
	}

	return operations.NewSmsSignupOK().WithPayload(jwt)
}

func (h AccountHandler) SmsLogin(p operations.SmsLoginParams) middleware.Responder {
	jwt, err := h.service.SmsLogin(restful.NewContext(p.HTTPRequest), p.Phone, p.SmsCode)
	if err != nil {
		return errors.Wrap(err)
	}

	return operations.NewSmsLoginOK().WithPayload(jwt)
}

func (h AccountHandler) Login(p operations.LoginParams) middleware.Responder {
	jwt, err := h.service.Login(restful.NewContext(p.HTTPRequest), p.Name, p.Password)
	if err != nil {
		return errors.Wrap(err)
	}

	return operations.NewLoginOK().WithPayload(jwt)
}

func (h AccountHandler) Logout(p operations.LogoutParams) middleware.Responder {
	err := h.service.Logout(restful.NewContext(p.HTTPRequest), p.Jwt)
	if err != nil {
		return errors.Wrap(err)
	}

	return operations.NewLogoutOK()
}

func (h AccountHandler) ResetPassword(p operations.ResetPasswordParams) middleware.Responder {
	err := h.service.ResetPassword(restful.NewContext(p.HTTPRequest), p.Phone, p.SmsCode, p.NewPassword)
	if err != nil {
		return errors.Wrap(err)
	}

	return operations.NewResetPasswordOK()
}
