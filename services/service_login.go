package services

import (
	"github.com/NeuronAccount/account/models"
	"github.com/NeuronAccount/account/storages/account_db"
	"github.com/NeuronFramework/errors"
	"github.com/NeuronFramework/restful"
	"strings"
)

func (s *AccountService) Login(ctx *restful.Context, name string, password string) (jwt string, err error) {
	var dbAccount *account_db.Account
	if strings.Contains(name, "@") { //email
		dbAccount, err = s.accountDB.Account.GetQuery().EmailAddress_Equal(name).QueryOne(ctx, nil)
	} else { //phone
		dbAccount, err = s.accountDB.Account.GetQuery().PhoneNumber_Equal(name).QueryOne(ctx, nil)
	}
	if err != nil {
		return "", err
	}

	if dbAccount == nil {
		return "", errors.NotFound("帐号不存在")
	}

	passwordHash := s.calcPasswordHash(password)
	if dbAccount.PasswordHash != passwordHash {
		op := &models.Operation{
			OperationType: models.OperationLogin,
			LoginName:     name,
			Error:         errors.Unauthorized("帐号或密码错误"),
		}
		s.addOperation(ctx, op)
		return "", op.Error
	}

	jwt, err = s.generateJwt(dbAccount.AccountId)
	if err != nil {
		return "", err
	}

	s.addOperation(ctx, &models.Operation{
		OperationType: models.OperationLogin,
		LoginName:     name,
		AccountID:     dbAccount.AccountId,
	})

	return jwt, nil
}
