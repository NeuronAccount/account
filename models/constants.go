package models

const SmsSceneSignup = "SMS_SIGNUP"
const SmsSceneLogin = "SMS_LOGIN"
const SmsSceneResetPassword = "RESET_PASSWORD"

const SmsCodeValidSeconds = 5 * 60 //5m

const MAX_LOGIN_NAME_LENGTH = 24
const MAX_PHONE_LENGTH = 11
const MAX_PASSWORD_LENGTH = 20
const MAX_SMS_CODE_LENGTH = 6
