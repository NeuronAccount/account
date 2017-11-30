package handler

import (
	"github.com/NeuronFramework/log"
	"github.com/NeuronFramework/restful"
	"github.com/NeuronGroup/Account/api/oauth2/gen/restapi/operations"
	"github.com/NeuronGroup/Account/models"
	"github.com/NeuronGroup/Account/services"
	"github.com/go-openapi/runtime/middleware"
	"go.uber.org/zap"
	"github.com/NeuronFramework/errors"
)

type OAuth2HandlerOptions struct {
	AccountStorageConnectionString string
	OAuthConnectionString          string
}

type OAuth2Handler struct {
	logger  *zap.Logger
	options *OAuth2HandlerOptions
	service *services.AccountService
}

func NewOAuth2Handler(options *OAuth2HandlerOptions) (h *OAuth2Handler, err error) {
	h = &OAuth2Handler{}
	h.logger = log.TypedLogger(h)
	h.options = options
	h.service, err = services.NewAccountService(&services.AccountServiceOptions{
		AccountStorageConnectionString: options.AccountStorageConnectionString,
		OAuthConnectionString:options.OAuthConnectionString,
	})
	if err != nil {
		return nil, err
	}

	return h, nil
}

func (h *OAuth2Handler) BasicAuth(clientId string, password string) (interface{}, error) {
	return h.service.OAuth2ClientLogin(clientId, password)
}

func (h *OAuth2Handler) Token(p operations.OAuth2TokenParams, oauthClient interface{}) middleware.Responder {
	if oauthClient == nil {
		return restful.Responder(errors.Unauthorized("client认证失败"))
	}

	result, err := h.service.OAuth2Token(toTokenRequest(&p), oauthClient.(*models.OAuth2Client))
	if err != nil {
		return restful.Responder(err)
	}

	return operations.NewOAuth2TokenOK().WithPayload(fromTokenResponse(result))
}
