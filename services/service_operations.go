package services

import (
	"github.com/NeuronAccount/account/models"
	"github.com/NeuronFramework/rest"
	"go.uber.org/zap"
	"strconv"
)

func (s *AccountService) addOperation(ctx *rest.Context, operation *models.AccountOperation) (err error) {
	operation.UserAgent = ctx.UserAgent
	dbOperation := toOperation(operation)
	_, err = s.accountDB.AccountOperation.Query().Insert(ctx, nil, dbOperation)
	if err != nil {
		s.logger.Error("addOperation", zap.Error(err))
		return err
	}

	return nil
}

func (s *AccountService) GetOperationList(ctx *rest.Context, userId string, query *models.OperationQuery) (items []*models.AccountOperation, nextPageToken string, err error) {
	q := s.accountDB.AccountOperation.Query()
	if query.OperationType != "" {
		q.OperationTypeEqual(query.OperationType)
	}
	pageToken := int64(0)
	pageSize := int64(40)
	if query.PageToken != "" {
		pageTokenI, err := strconv.Atoi(query.PageToken)
		if err != nil {
			return nil, "", rest.InvalidParam("PageToken无效")
		}
		pageToken = int64(pageTokenI)
	}
	if query.PageSize > 0 {
		pageSize = int64(query.PageSize)
	}

	dbOperationList, err := q.Limit(pageToken, pageSize).SelectList(ctx, nil)
	if err != nil {
		return nil, "", err
	}
	if dbOperationList == nil {
		return nil, "", nil
	}

	items = fromOperationList(dbOperationList)
	if items != nil {
		for _, v := range items {
			v.PhoneEncrypted = s.maskPhone(v.PhoneEncrypted)
		}
	}
	nextPageToken = strconv.FormatInt(pageToken+pageSize, 10)

	return items, nextPageToken, nil
}
