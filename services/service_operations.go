package services

import (
	"github.com/NeuronAccount/account/models"
	"github.com/NeuronFramework/restful"
	"go.uber.org/zap"
)

func (s *AccountService) addOperation(ctx *restful.Context, operation *models.AccountOperation) (err error) {
	operation.UserAgent = ctx.UserAgent
	dbOperation := toOperation(operation)
	_, err = s.userDB.AccountOperation.Insert(ctx, nil, dbOperation)
	if err != nil {
		s.logger.Error("addOperation", zap.Error(err))
		return err
	}

	return nil
}
