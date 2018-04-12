package services

import (
	"github.com/NeuronAccount/account/models"
	"github.com/NeuronAccount/account/storages/user_db"
)

func fromOperation(p *user_db.UserOperation) (r *models.Operation) {
	if p == nil {
		return nil
	}

	r = &models.Operation{}

	return r
}

func fromOperationList(p []*user_db.UserOperation) (r []*models.Operation) {
	if p == nil {
		return nil
	}

	r = make([]*models.Operation, len(p))
	for i, v := range p {
		r[i] = fromOperation(v)
	}

	return r
}

func toOperation(p *models.Operation) (r *user_db.UserOperation) {
	if p == nil {
		return nil
	}

	r = &user_db.UserOperation{}
	r.UserId = p.UserId
	r.OperationType = string(p.OperationType)
	r.UserAgent = p.UserAgent
	r.PhoneNumber = p.Phone

	return r
}
