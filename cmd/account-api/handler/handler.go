package handler

import (
	api "github.com/NeuronAccount/account/api/gen/models"
	"github.com/NeuronAccount/account/api/gen/restapi/operations"
	"github.com/NeuronAccount/account/models"
	"github.com/NeuronAccount/account/services"
	"github.com/NeuronFramework/errors"
	"github.com/NeuronFramework/log"
	"github.com/NeuronFramework/rest"
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
	if token == "" {
		return "", nil
	}

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
	err := h.service.SendSmsCode(rest.NewContext(p.HTTPRequest), &models.SendSmsCodeParams{
		UserId:      "",
		Scene:       models.SmsSceneSmsLogin,
		Phone:       p.Phone,
		CaptchaId:   p.CaptchaID,
		CaptchaCode: p.CaptchaCode,
	})
	if err != nil {
		return errors.Wrap(err)
	}

	return operations.NewSendLoginSmsCodeOK()
}

func (h *AccountHandler) SendSmsCode(p operations.SendSmsCodeParams, userId interface{}) middleware.Responder {
	err := h.service.SendSmsCode(rest.NewContext(p.HTTPRequest), &models.SendSmsCodeParams{
		UserId:      userId.(string),
		Scene:       models.SmsScene(p.Scene),
		Phone:       p.Phone,
		CaptchaId:   p.CaptchaID,
		CaptchaCode: p.CaptchaCode,
	})
	if err != nil {
		return errors.Wrap(err)
	}

	return operations.NewSendSmsCodeOK()
}

func (h *AccountHandler) SmsLogin(p operations.SmsLoginParams) middleware.Responder {
	userToken, err := h.service.SmsLogin(rest.NewContext(p.HTTPRequest), p.Phone, p.SmsCode)
	if err != nil {
		return errors.Wrap(err)
	}

	return operations.NewSmsLoginOK().WithPayload(fromUserToken(userToken))
}

func (h *AccountHandler) Logout(p operations.LogoutParams) middleware.Responder {
	err := h.service.Logout(rest.NewContext(p.HTTPRequest), p.AccessToken, p.RefreshToken)
	if err != nil {
		return errors.Wrap(err)
	}

	return operations.NewLogoutOK()
}

func (h *AccountHandler) RefreshToken(p operations.RefreshTokenParams) middleware.Responder {
	userToken, err := h.service.RefreshToken(rest.NewContext(p.HTTPRequest), p.RefreshToken)
	if err != nil {
		return errors.Wrap(err)
	}

	return operations.NewRefreshTokenOK().WithPayload(fromUserToken(userToken))
}

func (h *AccountHandler) OauthState(p operations.OauthStateParams) middleware.Responder {
	state, err := h.service.OauthState(rest.NewContext(p.HTTPRequest))
	if err != nil {
		return errors.Wrap(err)
	}

	return operations.NewOauthStateOK().WithPayload(state)
}

func (h *AccountHandler) OauthJump(p operations.OauthJumpParams) middleware.Responder {
	userToken, err := h.service.OauthJump(
		rest.NewContext(p.HTTPRequest),
		&models.OauthJumpParams{
			RedirectUri:       p.RedirectURI,
			AuthorizationCode: p.AuthorizationCode,
			State:             p.State,
		})

	if err != nil {
		return errors.Wrap(err)
	}

	return operations.NewOauthJumpOK().WithPayload(fromUserToken(userToken))
}

func (h *AccountHandler) GetUserInfo(p operations.GetUserInfoParams, userId interface{}) middleware.Responder {
	userInfo, err := h.service.GetUserInfo(rest.NewContext(p.HTTPRequest), userId.(string))
	if err != nil {
		return errors.Wrap(err)
	}

	return operations.NewGetUserInfoOK().WithPayload(fromUserInfo(userInfo))
}

func (h *AccountHandler) SetUserName(p operations.SetUserNameParams, userId interface{}) middleware.Responder {
	err := h.service.SetUserName(rest.NewContext(p.HTTPRequest), userId.(string), p.UserName)
	if err != nil {
		return errors.Wrap(err)
	}

	return operations.NewSetUserNameOK()
}

func (h *AccountHandler) SetUserIcon(p operations.SetUserIconParams, userId interface{}) middleware.Responder {
	err := h.service.SetUserIcon(rest.NewContext(p.HTTPRequest), userId.(string), p.UserIcon)
	if err != nil {
		return errors.Wrap(err)
	}

	return operations.NewSetUserIconOK()
}

func (h *AccountHandler) GetAccountInfo(p operations.GetAccountInfoParams, userId interface{}) middleware.Responder {
	accountInfo, err := h.service.GetAccountInfo(rest.NewContext(p.HTTPRequest), userId.(string))
	if err != nil {
		return errors.Wrap(err)
	}

	return operations.NewGetAccountInfoOK().WithPayload(fromAccountInfo(accountInfo))
}

func (h *AccountHandler) BindPhone(p operations.BindPhoneParams, userId interface{}) middleware.Responder {
	err := h.service.BindPhone(rest.NewContext(p.HTTPRequest), userId.(string), p.Phone, p.SmsCode)
	if err != nil {
		return errors.Wrap(err)
	}

	return operations.NewBindPhoneOK()
}

func (h *AccountHandler) UnbindPhone(p operations.UnbindPhoneParams, userId interface{}) middleware.Responder {
	err := h.service.BindPhone(rest.NewContext(p.HTTPRequest), userId.(string), p.Phone, p.SmsCode)
	if err != nil {
		return errors.Wrap(err)
	}

	return operations.NewUnbindPhoneOK()
}

func (h *AccountHandler) BindOauthAccount(p operations.BindOauthAccountParams, userId interface{}) middleware.Responder {
	err := h.service.BindOauthAccount(rest.NewContext(p.HTTPRequest), userId.(string))
	if err != nil {
		return errors.Wrap(err)
	}

	return operations.NewBindOauthAccountOK()
}

func (h *AccountHandler) UnbindOauthAccount(p operations.UnbindOauthAccountParams, userId interface{}) middleware.Responder {
	err := h.service.UnbindOauthAccount(rest.NewContext(p.HTTPRequest), userId.(string))
	if err != nil {
		return errors.Wrap(err)
	}

	return operations.NewUnbindOauthAccountOK()
}

func (h *AccountHandler) GetOperationList(p operations.GetOperationListParams, userId interface{}) middleware.Responder {
	query := &models.OperationQuery{}
	if p.OperationType != nil {
		query.OperationType = *p.OperationType
	}
	if p.PageToken != nil {
		query.PageToken = *p.PageToken
	}
	if p.PageSize != nil {
		query.PageSize = *p.PageSize
	}

	items, nextPageToken, err := h.service.GetOperationList(rest.NewContext(p.HTTPRequest), userId.(string), query)
	if err != nil {
		return errors.Wrap(err)
	}

	response := &api.OperationListResponse{
		Items:         fromOperationList(items),
		NextPageToken: &nextPageToken,
	}

	return operations.NewGetOperationListOK().WithPayload(response)
}
