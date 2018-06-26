package models

const SmsCodeValidSeconds = 5 * 60 //验证码有效期5分钟
const SmsCodeLength = 4            //验证码长度为4

const (
	SmsSceneSmsLogin      = "SMS_LOGIN"
	SmsSceneResetPassword = "RESET_PASSWORD"
	SmsSceneBindPhone     = "BIND_PHONE"
	SmsSceneUnbindPhone   = "UNBIND_PHONE"
)

type SendSmsCodeParams struct {
	UserId      string
	Scene       string
	Phone       string
	CaptchaId   string
	CaptchaCode string
}
