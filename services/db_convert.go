package services

import (
	"github.com/NeuronAccount/account/models"
	"github.com/NeuronAccount/account/storages/user_db"
)

func fromOperation(p *user_db.AccountOperation) (r *models.AccountOperation) {
	if p == nil {
		return nil
	}

	r = &models.AccountOperation{}

	return r
}

func fromOperationList(p []*user_db.AccountOperation) (r []*models.AccountOperation) {
	if p == nil {
		return nil
	}

	r = make([]*models.AccountOperation, len(p))
	for i, v := range p {
		r[i] = fromOperation(v)
	}

	return r
}

func toOperation(p *models.AccountOperation) (r *user_db.AccountOperation) {
	if p == nil {
		return nil
	}

	r = &user_db.AccountOperation{}
	r.UserId = p.UserId
	r.OperationType = string(p.OperationType)
	r.UserAgent = p.UserAgent
	r.PhoneNumber = p.Phone

	return r
}

func fromUserInfo(p *user_db.User) (r *models.UserInfo) {
	if p == nil {
		return nil
	}

	r = &models.UserInfo{}
	r.UserID = p.UserId
	r.Name = p.UserName
	r.Icon = p.UserIcon

	return r
}
