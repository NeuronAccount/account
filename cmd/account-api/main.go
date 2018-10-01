package main

import (
	"github.com/NeuronAccount/account/api/gen/restapi"
	"github.com/NeuronAccount/account/api/gen/restapi/operations"
	"github.com/NeuronAccount/account/cmd/account-api/handler"
	"github.com/NeuronFramework/rest"
	"github.com/go-openapi/loads"
	"net/http"
)

func main() {
	rest.Run(func() (http.Handler, error) {
		swaggerSpec, err := loads.Analyzed(restapi.SwaggerJSON, "")
		if err != nil {
			return nil, err
		}

		h, err := handler.NewAccountHandler()
		if err != nil {
			return nil, err
		}

		api := operations.NewAccountAPI(swaggerSpec)
		api.ServeError = rest.ServeError
		api.BearerAuth = h.BearerAuth
		api.SendSmsCodeHandler = operations.SendSmsCodeHandlerFunc(h.SendSmsCode)
		api.SendLoginSmsCodeHandler = operations.SendLoginSmsCodeHandlerFunc(h.SendLoginSmsCode)
		api.SmsLoginHandler = operations.SmsLoginHandlerFunc(h.SmsLogin)
		api.PhonePasswordLoginHandler = operations.PhonePasswordLoginHandlerFunc(h.PhonePasswordLogin)
		api.LogoutHandler = operations.LogoutHandlerFunc(h.Logout)
		api.RefreshTokenHandler = operations.RefreshTokenHandlerFunc(h.RefreshToken)
		api.OauthStateHandler = operations.OauthStateHandlerFunc(h.OauthState)
		api.OauthJumpHandler = operations.OauthJumpHandlerFunc(h.OauthJump)
		api.ResetPasswordHandler = operations.ResetPasswordHandlerFunc(h.ResetPassword)
		api.GetUserInfoHandler = operations.GetUserInfoHandlerFunc(h.GetUserInfo)
		api.SetUserNameHandler = operations.SetUserNameHandlerFunc(h.SetUserName)
		api.SetUserIconHandler = operations.SetUserIconHandlerFunc(h.SetUserIcon)
		api.GetAccountInfoHandler = operations.GetAccountInfoHandlerFunc(h.GetAccountInfo)
		api.BindPhoneHandler = operations.BindPhoneHandlerFunc(h.BindPhone)
		api.UnbindPhoneHandler = operations.UnbindPhoneHandlerFunc(h.UnbindPhone)
		api.BindOauthAccountHandler = operations.BindOauthAccountHandlerFunc(h.BindOauthAccount)
		api.UnbindOauthAccountHandler = operations.UnbindOauthAccountHandlerFunc(h.UnbindOauthAccount)
		api.GetOperationListHandler = operations.GetOperationListHandlerFunc(h.GetOperationList)

		return api.Serve(nil), nil
	})
}
