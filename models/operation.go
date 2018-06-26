package models

import "time"

const (
	OperationSendSmsCode        = "SEND_SMS_CODE"
	OperationSmsLogin           = "SMS_LOGIN"
	OperationPhonePasswordLogin = "PHONE_PASSWORD_LOGIN"
	OperationLogout             = "LOGOUT"
	OperationBindPhone          = "BIND_PHONE"
	OperationUnbindPhone        = "UNBIND_PHONE"
	OperationResetPassword      = "RESET_PASSWORD"
	OperationRemoveAccount      = "REMOVE_ACCOUNT"
)

type AccountOperation struct {
	OperationId    string
	UserId         string
	OperationType  string
	OperationTime  time.Time
	UserAgent      string
	PhoneEncrypted string
	SmsScene       string
	OtherUserId    string
}

type OperationQuery struct {
	OperationType string
	PageToken     string
	PageSize      int32
}
