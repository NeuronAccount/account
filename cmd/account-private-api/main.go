package main

import (
	"github.com/NeuronAccount/account/api-private/gen/restapi"
	"github.com/NeuronAccount/account/api-private/gen/restapi/operations"
	"github.com/NeuronAccount/account/cmd/account-private-api/handler"
	"github.com/NeuronFramework/restful"
	"github.com/go-openapi/loads"
	"net/http"
)

func main() {
	restful.Run(func() (http.Handler, error) {
		h, err := handler.NewAccountHandler()
		if err != nil {
			return nil, err
		}

		swaggerSpec, err := loads.Analyzed(restapi.SwaggerJSON, "")
		if err != nil {
			return nil, err
		}

		api := operations.NewAccountPrivateAPI(swaggerSpec)
		api.SmsCodeHandler = operations.SmsCodeHandlerFunc(h.SmsCode)
		api.SmsSignupHandler = operations.SmsSignupHandlerFunc(h.SmsSignup)
		api.SmsLoginHandler = operations.SmsLoginHandlerFunc(h.SmsLogin)
		api.LoginHandler = operations.LoginHandlerFunc(h.Login)
		api.LogoutHandler = operations.LogoutHandlerFunc(h.Logout)
		api.ResetPasswordHandler = operations.ResetPasswordHandlerFunc(h.ResetPassword)

		return api.Serve(nil), nil
	})
}
