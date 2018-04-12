package models

import (
	"time"
)

type OperationType string

const (
	OperationSmsLogin = OperationType("SMS_LOGIN")
	OperationLogout   = OperationType("LOGOUT")
)

type Operation struct {
	OperationId   string
	OperationTime time.Time
	OperationType OperationType
	UserId        string
	UserAgent     string
	Phone         string
}
