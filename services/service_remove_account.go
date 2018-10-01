package services

import (
	"github.com/NeuronAccount/account/models"
	"github.com/NeuronFramework/rest"
	"github.com/NeuronFramework/sql/wrap"
)

func (s *AccountService) RemoveAccount(ctx *rest.Context, userId string) (err error) {
	if userId == "" {
		return rest.InvalidParam("userId不能为空")
	}

	return s.accountDB.TransactionReadCommitted(ctx, false, func(tx *wrap.Tx) (err error) {
		dbPhoneAccount, err := s.accountDB.PhoneAccount.Query().UserIdEqual(userId).Select(ctx, nil)
		if err != nil {
			return nil
		}
		if dbPhoneAccount != nil && dbPhoneAccount.PhoneEncrypted != "" {
			_, err = s.accountDB.SmsCode.Query().PhoneEncryptedEqual(dbPhoneAccount.PhoneEncrypted).Delete(ctx, tx)
			if err != nil {
				return err
			}

			_, err = s.accountDB.AccountOperation.Query().PhoneEncryptedEqual(dbPhoneAccount.PhoneEncrypted).Delete(ctx, tx)
			if err != nil {
				return err
			}
		}

		_, err = s.accountDB.AccessToken.Query().UserIdEqual(userId).Delete(ctx, tx)
		if err != nil {
			return err
		}

		_, err = s.accountDB.RefreshToken.Query().UserIdEqual(userId).Delete(ctx, tx)
		if err != nil {
			return err
		}

		_, err = s.accountDB.AccountOperation.Query().UserIdEqual(userId).Delete(ctx, tx)
		if err != nil {
			return err
		}

		_, err = s.accountDB.UserInfo.Query().UserIdEqual(userId).Delete(ctx, tx)
		if err != nil {
			return err
		}

		_, err = s.accountDB.PhoneAccount.Query().UserIdEqual(userId).Delete(ctx, tx)
		if err != nil {
			return err
		}

		_, err = s.accountDB.SmsCode.Query().UserIdEqual(userId).Delete(ctx, tx)
		if err != nil {
			return err
		}

		_, err = s.accountDB.AccountOperation.Query().UserIdEqual(userId).Delete(ctx, tx)
		if err != nil {
			return err
		}

		//操作纪录
		s.addOperation(ctx, &models.AccountOperation{
			OperationType: models.OperationRemoveAccount,
			UserId:        dbPhoneAccount.UserId,
		})

		return nil
	})
}
