package models

const SmsCodeValidSeconds = 5 * 60 //验证码有效期5分钟
const SmsCodeLength = 4            //验证码长度为4

type SmsScene string

const (
	SmsSceneSmsLogin    = SmsScene("SMS_LOGIN")
	SmsSceneBindPhone   = SmsScene("BIND_PHONE")
	SmsSceneUnbindPhone = SmsScene("UNBIND_PHONE")
)

type SendSmsCodeParams struct {
	UserId      string
	Scene       SmsScene
	Phone       string
	CaptchaId   string
	CaptchaCode string
}
