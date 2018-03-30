package models

import (
	"github.com/NeuronFramework/errors"
	"time"
)

type OperationType string

const (
	OperationSmsCode       = OperationType("SMS_CODE")
	OperationSmsSignup     = OperationType("SMS_SIGNUP")
	OperationSmsLogin      = OperationType("SMS_LOGIN")
	OperationLogin         = OperationType("LOGIN")
	OperationLogout        = OperationType("LOGOUT")
	OperationResetPassword = OperationType("RESET_PASSWORD")
)

type Operation struct {
	OperationId   string
	OperationTime time.Time
	OperationType OperationType
	Error         *errors.Error
	UserAgent     string
	SmsScene      string
	Phone         string
	LoginName     string
	AccountID     string
}
