package main

import (
	"github.com/NeuronFramework/log"
	"github.com/NeuronFramework/restful"
	"github.com/NeuronGroup/Account/api/private/gen/restapi"
	"github.com/NeuronGroup/Account/api/private/gen/restapi/operations"
	"github.com/NeuronGroup/Account/cmd/accounts-private-api/handler"
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

	var bind_addr string
	var storageConnectionString = ""

	cmd := cobra.Command{}
	cmd.PersistentFlags().StringVar(&bind_addr, "bind-addr", ":8083", "api server bind addr")
	cmd.PersistentFlags().StringVar(&storageConnectionString, "mysql-connection-string",
		"root:123456@tcp(127.0.0.1:3307)/account?parseTime=true", "mysql connection string")
	cmd.RunE = func(cmd *cobra.Command, args []string) error {
		swaggerSpec, err := loads.Analyzed(restapi.SwaggerJSON, "")
		if err != nil {
			return err
		}
		api := operations.NewAccountsAPI(swaggerSpec)

		h, err := handler.NewAccountHandler(&handler.AccountHandlerOptions{
			AccountStorageConnectionString: storageConnectionString})
		if err != nil {
			return err
		}

		api.SmsCodeHandler = operations.SmsCodeHandlerFunc(h.SmsCode)
		api.SmsSignupHandler = operations.SmsSignupHandlerFunc(h.SmsSignup)
		api.SmsLoginHandler = operations.SmsLoginHandlerFunc(h.SmsLogin)
		api.LoginHandler = operations.LoginHandlerFunc(h.Login)
		api.LogoutHandler = operations.LogoutHandlerFunc(h.Logout)

		logger.Info("Start server", zap.String("addr", bind_addr))
		err = http.ListenAndServe(bind_addr,
			restful.Recovery(cors.AllowAll().Handler(api.Serve(nil))))
		if err != nil {
			return err
		}

		return nil
	}
	cmd.Execute()
}
