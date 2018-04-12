package models

import "time"

type OperationType string

const (
	OperationSmsLogin = OperationType("SMS_LOGIN")
	OperationLogout   = OperationType("LOGOUT")
)

type AccountOperation struct {
	OperationId   string
	OperationTime time.Time
	OperationType OperationType
	UserId        string
	UserAgent     string
	Phone         string
}

type UserToken struct {
	AccessToken  string
	RefreshToken string
}

type UserInfo struct {
	UserID string
	Name   string
	Icon   string
}
