package models

import "time"

type OperationType string

const (
	OperationSendSmsCode = OperationType("SEND_SMS_CODE")
	OperationSmsLogin    = OperationType("SMS_LOGIN")
	OperationLogout      = OperationType("LOGOUT")
	OperationBindPhone   = OperationType("BIND_PHONE")
	OperationUnbindPhone = OperationType("UNBIND_PHONE")
)

type AccountOperation struct {
	OperationId    string
	UserId         string
	OperationType  OperationType
	OperationTime  time.Time
	UserAgent      string
	PhoneEncrypted string
	SmsScene       SmsScene
	OtherUserId    string
}

type OperationQuery struct {
	OperationType string
	PageToken     string
	PageSize      int32
}
