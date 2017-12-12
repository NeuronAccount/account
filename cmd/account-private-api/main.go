package main

import (
	"github.com/NeuronAccount/account/api-private/gen/restapi"
	"github.com/NeuronAccount/account/api-private/gen/restapi/operations"
	"github.com/NeuronAccount/account/cmd/account-private-api/handler"
	"github.com/NeuronFramework/log"
	"github.com/NeuronFramework/restful"
	"github.com/go-openapi/loads"
	"github.com/go-openapi/runtime/middleware"
	"github.com/rs/cors"
	"github.com/spf13/cobra"
	"go.uber.org/zap"
	"net/http"
)

func main() {
	log.Init(true)

	middleware.Debug = false

	logger := zap.L().Named("main")

	var bindAddr string

	cmd := cobra.Command{}
	cmd.PersistentFlags().StringVar(&bindAddr, "bind-addr", ":8083", "api server bind addr")
	cmd.RunE = func(cmd *cobra.Command, args []string) error {
		swaggerSpec, err := loads.Analyzed(restapi.SwaggerJSON, "")
		if err != nil {
			return err
		}
		api := operations.NewAccountPrivateAPI(swaggerSpec)

		h, err := handler.NewAccountHandler()
		if err != nil {
			return err
		}

		api.SmsCodeHandler = operations.SmsCodeHandlerFunc(h.SmsCode)
		api.SmsSignupHandler = operations.SmsSignupHandlerFunc(h.SmsSignup)
		api.SmsLoginHandler = operations.SmsLoginHandlerFunc(h.SmsLogin)
		api.LoginHandler = operations.LoginHandlerFunc(h.Login)
		api.LogoutHandler = operations.LogoutHandlerFunc(h.Logout)

		logger.Info("Start server", zap.String("addr", bindAddr))
		err = http.ListenAndServe(bindAddr,
			restful.Recovery(cors.AllowAll().Handler(api.Serve(nil))))
		if err != nil {
			return err
		}

		return nil
	}
	cmd.Execute()
}
