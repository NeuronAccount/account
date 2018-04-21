package handler

import "github.com/NeuronAccount/account/models"
import (
	api "github.com/NeuronAccount/account/api/gen/models"
	"github.com/NeuronFramework/errors"
	"github.com/go-openapi/strfmt"
)

func fromUserToken(p *models.UserToken) (r *api.UserToken) {
	if p == nil {
		return nil
	}

	r = &api.UserToken{}
	r.AccessToken = &p.AccessToken
	r.RefreshToken = &p.RefreshToken

	return r
}

func fromUserInfo(p *models.UserInfo) (r *api.UserInfo) {
	if p == nil {
		return nil
	}

	r = &api.UserInfo{}
	r.UserID = &p.UserID
	r.UserName = &p.UserName
	r.UserIcon = &p.UserIcon

	return r
}

func fromOauthAccountInfo(p *models.OauthAccountInfo) (r *api.OauthAccountInfo) {
	if p == nil {
		return nil
	}

	r = &api.OauthAccountInfo{}
	r.ProviderID = &p.ProviderId
	r.ProviderName = &p.ProviderName
	r.OpenID = &p.OpenId
	r.UserName = &p.UserName
	r.UserIcon = &p.UserIcon

	return r
}

func fromOauthAccountInfoList(p []*models.OauthAccountInfo) (r []*api.OauthAccountInfo) {
	if p == nil {
		return nil
	}

	r = make([]*api.OauthAccountInfo, len(p))
	for i, v := range p {
		r[i] = fromOauthAccountInfo(v)
	}

	return r
}

func fromAccountInfo(p *models.AccountInfo) (r *api.AccountInfo) {
	if p == nil {
		return nil
	}

	r = &api.AccountInfo{}
	r.UserID = &p.UserId
	r.UserName = &p.UserName
	r.UserIcon = &p.UserIcon
	r.PhoneBinded = &p.PhoneBinded
	r.OauthBindedList = fromOauthAccountInfoList(p.OauthBindedList)

	return r
}

func fromOperationType(p models.OperationType) (r api.OperationType) {
	switch p {
	case models.OperationSendSmsCode:
		return api.OperationTypeSENDSMSCODE
	case models.OperationSmsLogin:
		return api.OperationTypeSMSLOGIN
	case models.OperationLogout:
		return api.OperationTypeLOGOUT
	default:
		panic(errors.Unknown("未知的操作类型" + string(p)))
	}
}

func fromSmsScene(p models.SmsScene) (r api.SmsScene) {
	switch p {
	case models.SmsSceneSmsLogin:
		return api.SmsSceneSMSLOGIN
	case models.SmsSceneBindPhone:
		return api.SmsSceneBINDPHONE
	case models.SmsSceneUnbindPhone:
		return api.SmsSceneUNBINDPHONE
	default:
		panic(errors.Unknown("未知的操作类型" + string(p)))
	}
}

func fromOperation(p *models.AccountOperation) (r *api.Operation) {
	if p == nil {
		return nil
	}

	r = &api.Operation{}
	r.OperationID = &p.OperationId
	r.OperationType = fromOperationType(p.OperationType)
	operationTime := strfmt.DateTime(p.OperationTime)
	r.OperationTime = &operationTime
	r.PhoneMasked = p.PhoneEncrypted
	r.SmsScene = fromSmsScene(p.SmsScene)
	r.UserID = &p.UserId
	r.UserAgent = p.UserAgent

	return r
}

func fromOperationList(p []*models.AccountOperation) (r []*api.Operation) {
	if p == nil {
		return nil
	}

	r = make([]*api.Operation, len(p))
	for i, v := range p {
		r[i] = fromOperation(v)
	}

	return r
}
