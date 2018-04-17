package services

import (
	"github.com/NeuronAccount/account/models"
	"github.com/NeuronAccount/account/storages/neuron_account_db"
	"strconv"
)

func fromOperation(p *neuron_account_db.AccountOperation) (r *models.AccountOperation) {
	if p == nil {
		return nil
	}

	r = &models.AccountOperation{}
	r.OperationId = strconv.Itoa(int(p.Id))
	r.UserId = p.UserId
	r.OperationType = models.OperationType(p.OperationType)
	r.OperationTime = p.CreateTime
	r.UserAgent = p.UserAgent
	r.PhoneEncrypted = p.PhoneEncrypted
	r.SmsScene = models.SmsScene(p.SmsScene)
	r.OtherUserId = p.OtherUserId

	return r
}

func fromOperationList(p []*neuron_account_db.AccountOperation) (r []*models.AccountOperation) {
	if p == nil {
		return nil
	}

	r = make([]*models.AccountOperation, len(p))
	for i, v := range p {
		r[i] = fromOperation(v)
	}

	return r
}

func toOperation(p *models.AccountOperation) (r *neuron_account_db.AccountOperation) {
	if p == nil {
		return nil
	}

	r = &neuron_account_db.AccountOperation{}
	r.UserId = p.UserId
	r.OperationType = string(p.OperationType)
	r.UserAgent = p.UserAgent
	r.PhoneEncrypted = p.PhoneEncrypted
	r.SmsScene = string(p.SmsScene)
	r.OtherUserId = p.OtherUserId

	return r
}

func fromUserInfo(p *neuron_account_db.UserInfo) (r *models.UserInfo) {
	if p == nil {
		return nil
	}

	r = &models.UserInfo{}
	r.UserID = p.UserId
	r.UserName = p.UserName
	r.UserIcon = p.UserIcon

	return r
}

func fromOauthAccount(p *neuron_account_db.OauthAccount) (r *models.OauthAccountInfo) {
	if p == nil {
		return nil
	}

	r = &models.OauthAccountInfo{}
	r.ProviderId = p.ProviderId
	r.ProviderName = p.ProviderName
	r.OpenId = p.OpenId
	r.UserName = p.UserName
	r.UserIcon = p.UserIcon

	return r
}

func fromOauthAccountList(p []*neuron_account_db.OauthAccount) (r []*models.OauthAccountInfo) {
	if p == nil {
		return nil
	}

	r = make([]*models.OauthAccountInfo, len(p))
	for i, v := range p {
		r[i] = fromOauthAccount(v)
	}

	return r
}
