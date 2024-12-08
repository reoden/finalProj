package pkgs

import (
	"errors"
	"os"

	"github.com/aliyun/alibaba-cloud-sdk-go/services/dysmsapi"
)

// 向手机发送验证码
func SendMsg(tel string, code string) (*dysmsapi.SendSmsResponse, error) {
	smsKeyID := os.Getenv("SMSKeyID")
	smsKeySecret := os.Getenv("SMSKeySecret")
	client, err := dysmsapi.NewClientWithAccessKey("cn-hangzhou", smsKeyID, smsKeySecret)
	if err != nil {
		return nil, err
	}

	smsTemplateCode := os.Getenv("SMSTemplateCode")
	request := dysmsapi.CreateSendSmsRequest()
	request.Scheme = "https"
	request.PhoneNumbers = tel             // 手机号变量值
	request.SignName = "餐鉴文化传播"            // 签名
	request.TemplateCode = smsTemplateCode // 模板编码
	request.TemplateParam = "{\"code\":\"" + code + "\"}"
	response, err := client.SendSms(request)
	if err != nil {
		return response, err
	}

	if response.Code == "isv.BUSINESS_LIMIT_CONTROL" {
		return response, errors.New("frequency_limit")
	}

	return response, nil
}
