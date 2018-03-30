package services

import (
	"github.com/NeuronAccount/account/models"
	"github.com/NeuronAccount/account/storages/account_db"
	"github.com/NeuronFramework/errors"
	"strconv"
)

func fromOperation(p *account_db.AccountOperation) (r *models.Operation) {
	if p == nil {
		return nil
	}

	r = &models.Operation{}
	r.OperationId = strconv.Itoa(int(p.Id))
	r.OperationTime = p.CreateTime
	r.OperationType = models.OperationType(p.OperationType)
	r.Error = &errors.Error{Status: int(p.ErrorStatus), Code: p.ErrorCode, Message: p.ErrorMessage}
	r.UserAgent = p.UserAgent
	r.SmsScene = p.SmsScene
	r.Phone = p.PhoneNumber
	r.LoginName = p.LoginName
	r.AccountID = p.AccountId

	return r
}

func fromOperationList(p []*account_db.AccountOperation) (r []*models.Operation) {
	if p == nil {
		return nil
	}

	r = make([]*models.Operation, len(p))
	for i, v := range p {
		r[i] = fromOperation(v)
	}

	return r
}

func toOperation(p *models.Operation) (r *account_db.AccountOperation) {
	if p == nil {
		return nil
	}

	r = &account_db.AccountOperation{}
	r.OperationType = string(p.OperationType)
	if p.Error != nil {
		r.ErrorStatus = int32(p.Error.Status)
		r.ErrorCode = p.Error.Code
		r.ErrorMessage = p.Error.Message
	}
	r.UserAgent = p.UserAgent
	r.SmsScene = p.SmsScene
	r.PhoneNumber = p.Phone
	r.LoginName = p.LoginName
	r.AccountId = p.AccountID

	return r
}
