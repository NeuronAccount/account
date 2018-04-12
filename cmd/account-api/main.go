package main

import (
	"github.com/NeuronAccount/account/api/gen/restapi"
	"github.com/NeuronAccount/account/api/gen/restapi/operations"
	"github.com/NeuronAccount/account/cmd/account-api/handler"
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

		api := operations.NewAccountAPI(swaggerSpec)
		api.SendLoginSmsCodeHandler = operations.SendLoginSmsCodeHandlerFunc(h.SendLoginSmsCode)
		api.SmsLoginHandler = operations.SmsLoginHandlerFunc(h.SmsLogin)
		api.LogoutHandler = operations.LogoutHandlerFunc(h.Logout)
		api.RefreshTokenHandler = operations.RefreshTokenHandlerFunc(h.RefreshToken)

		return api.Serve(nil), nil
	})
}
