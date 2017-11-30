package services

import (
	"context"
	"github.com/NeuronFramework/errors"
	"github.com/NeuronGroup/Account/models"
	"github.com/NeuronGroup/Account/storages/account"
	"math/rand"
	"strconv"
	"time"
)

func (s *AccountService) SmsSignup(phone string, smsCode string, password string, oAuth2Params *models.OAuth2AuthorizeParams) (jwt string, err error) {
	dbSmsCode, err := s.db.SmsCode.GetQuery().
		SceneType_Equal(models.SCENE_TYPE_SMS_SIGNUP).And().PhoneNumber_Equal(phone).
		OrderBy(account.SMS_CODE_FIELD_CREATE_TIME, false).QueryOne(context.Background(), nil)
	if err != nil {
		return "", err
	}

	if dbSmsCode == nil || dbSmsCode.SmsCode != smsCode {
		return "", errors.BadRequest("InvalidSmsCode", "验证码错误")
	}

	tx, err := s.db.BeginReadCommittedTx(context.Background(), false)
	if err != nil {
		return "", err
	}
	defer tx.Rollback()

	dbAccount, err := s.db.Account.GetQuery().ForUpdate().
		PhoneNumber_Equal(phone).QueryOne(context.Background(), tx)
	if err != nil {
		return "", err
	}

	if dbAccount != nil {
		return "", errors.AlreadyExists("帐号已存在")
	}

	dbAccountIdGen, err := s.db.AccountIdGen.GetQuery().ForUpdate().QueryOne(context.Background(), tx)
	if err != nil {
		return "", err
	}

	if dbAccountIdGen == nil {
		dbAccountIdGen = &account.AccountIdGen{}
		dbAccountIdGen.MaxId = models.ACCOUNT_ID_GEN_START + 1 + rand.Int63n(models.ACCOUNT_ID_GEN_STEP)
		_, err := s.db.AccountIdGen.Insert(context.Background(), tx, dbAccountIdGen)
		if err != nil {
			return "", err
		}
	} else {
		dbAccountIdGen.MaxId = dbAccountIdGen.MaxId + 1 + rand.Int63n(models.ACCOUNT_ID_GEN_STEP)
		err := s.db.AccountIdGen.Update(context.Background(), tx, dbAccountIdGen)
		if err != nil {
			return "", err
		}
	}

	dbAccount = &account.Account{}
	dbAccount.CreateTime = time.Now()
	dbAccount.UpdateTime = time.Now()
	dbAccount.PhoneNumber = phone
	dbAccount.EmailAddress = ""
	dbAccount.PasswordHash = s.calcPasswordHash(password)
	dbAccount.OauthProvider = ""
	dbAccount.OauthAccountId = ""
	dbAccount.AccountId = strconv.Itoa(int(dbAccountIdGen.MaxId))
	_, err = s.db.Account.Insert(context.Background(), tx, dbAccount)
	if err != nil {
		return "", err
	}

	err = tx.Commit()
	if err != nil {
		return "", err
	}

	//生成Token
	jwt, err = generateJwt(dbAccount.AccountId, oAuth2Params)
	if err != nil {
		return "", err
	}

	return jwt, nil
}
