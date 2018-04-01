package sms

import (
	"crypto"
	"crypto/hmac"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"github.com/NeuronFramework/errors"
	"github.com/NeuronFramework/log"
	"github.com/NeuronFramework/rand"
	"go.uber.org/zap"
	"net/http"
	"net/url"
	"sort"
	"strings"
	"time"
)

const aliyunSmsUrl = "http://dysmsapi.aliyuncs.com"
const accessKey = "GVI7f5CudVDQ9tIqvCMAB8HPmI39UD"
const queryAccessKeyId = "LTAIiOpcTYrEjWqX"
const queryTimestampFormat = "2006-01-02T15:04:05Z"
const queryFormat = "JSON"
const querySignatureMethod = "HMAC-SHA1"
const querySignatureVersion = "1.0"
const queryAction = "SendSms"
const queryVersion = "2017-05-25"
const queryRegionId = "cn-hangzhou"
const querySignName = "火星基地"
const queryTemplateCode = "SMS_127158405"
const queryTemplateParam = "{\"code\":\"%s\"}"

type KV struct {
	K string
	V string
}

type KVArray []*KV

func (arr KVArray) Len() int {
	return len(arr)
}

func (arr KVArray) Less(i, j int) bool {
	return arr[i].K < arr[j].K
}

func (arr KVArray) Swap(i, j int) {
	temp := arr[i]
	arr[i] = arr[j]
	arr[j] = temp
}

type AliyunSmsResponse struct {
	RequestId string
	Code      string
	Message   string
	BizId     string
}

type Service struct {
	logger *zap.Logger
}

func New() (s *Service, err error) {
	s = &Service{}
	s.logger = log.TypedLogger(s)

	return s, nil
}

func (s *Service) encodeUrl(str string) (r string) {
	r = url.QueryEscape(str)
	r = strings.Replace(r, "+", "%20", -1)
	r = strings.Replace(r, "*", "%2A", -1)
	r = strings.Replace(r, "%7E", "~", -1)

	return r
}

func (s *Service) signature(stringToSign string) (signature string, err error) {
	hmacSha1 := hmac.New(crypto.SHA1.New, []byte(accessKey+"&"))
	_, err = hmacSha1.Write([]byte(stringToSign))
	if err != nil {
		return "", err
	}
	sig := hmacSha1.Sum(nil)

	return base64.StdEncoding.EncodeToString(sig), nil
}

func (s *Service) buildUrl(phone string, smsCode string, smsCodeId string) (urlString string, err error) {
	query := KVArray{}
	query = append(query, &KV{K: "AccessKeyId", V: queryAccessKeyId})
	query = append(query, &KV{K: "Timestamp", V: time.Now().UTC().Format(queryTimestampFormat)})
	query = append(query, &KV{K: "Format", V: queryFormat})
	query = append(query, &KV{K: "SignatureMethod", V: querySignatureMethod})
	query = append(query, &KV{K: "SignatureVersion", V: querySignatureVersion})
	query = append(query, &KV{K: "SignatureNonce", V: rand.NextHex(16)})

	query = append(query, &KV{K: "Action", V: queryAction})
	query = append(query, &KV{K: "Version", V: queryVersion})
	query = append(query, &KV{K: "RegionId", V: queryRegionId})
	query = append(query, &KV{K: "PhoneNumbers", V: phone})
	query = append(query, &KV{K: "SignName", V: querySignName})
	query = append(query, &KV{K: "TemplateCode", V: queryTemplateCode})
	query = append(query, &KV{K: "TemplateParam", V: fmt.Sprintf(queryTemplateParam, smsCode)})
	query = append(query, &KV{K: "OutId", V: smsCodeId})

	sort.Sort(query)

	queryStrings := make([]string, 0)
	for _, v := range query {
		kvString := s.encodeUrl(v.K) + "=" + s.encodeUrl(v.V)
		queryStrings = append(queryStrings, kvString)
	}
	sortedQueryString := strings.Join(queryStrings, "&")

	stringToSign := "GET" + "&" + s.encodeUrl("/") + "&" + s.encodeUrl(sortedQueryString)
	signature, err := s.signature(stringToSign)
	if err != nil {
		return "", err
	}

	urlString = aliyunSmsUrl + "/?" + "Signature=" + s.encodeUrl(signature) + "&" + sortedQueryString

	return urlString, nil
}

func (s *Service) wrapError(code string, message string) (err error) {
	switch code {
	case "OK":
		return nil
	case "isv.MOBILE_NUMBER_ILLEGAL":
		return errors.InvalidParam("手机号格式错误")
	case "isv.BUSINESS_LIMIT_CONTROL":
		return errors.BadRequest("SendLimit", "每天最多发送5次")
	default:
		return errors.Unknown("短信发送后端服务失败，code＝" + code + ",message=" + message)
	}
}

func (s *Service) SendSms(phone string, smsCode string, smsCodeId string) (requestId string, err error) {
	urlString, err := s.buildUrl(phone, smsCode, smsCodeId)
	if err != nil {
		return "", err
	}

	s.logger.Info("SendSms", zap.String("url", urlString))

	req, err := http.NewRequest("GET", urlString, nil)
	if err != nil {
		return "", err
	}

	c := &http.Client{}
	resp, err := c.Do(req)
	if err != nil {
		return "", err
	}

	smsResponse := AliyunSmsResponse{}
	decoder := json.NewDecoder(resp.Body)
	decoder.Decode(&smsResponse)
	if err != nil {
		return "", err
	}

	s.logger.Info("SendSms", zap.Any("resp", smsResponse))

	return smsResponse.RequestId, s.wrapError(smsResponse.Code, smsResponse.Message)
}
