package utils

import (
	"log"
	"psyWeb/configuration"

	"github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/common"
	"github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/common/errors"
	"github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/common/profile"
	sms "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/sms/v20210111" // 引入sms
)

func SendSMSByTencentCloud(phone_number string, code string) {
	credential := common.NewCredential(
		configuration.GetConfigInstance().SMS.SecretId,
		configuration.GetConfigInstance().SMS.SecretKey,
	)
	cpf := profile.NewClientProfile()
	client, _ := sms.NewClient(credential, "ap-guangzhou", cpf)
	r := sms.NewSendSmsRequest()
	/* 短信应用ID: 短信SdkAppId */
	r.SmsSdkAppId = common.StringPtr(configuration.GetConfigInstance().SMS.SdkAppId)
	/* 短信签名内容 */
	r.SignName = common.StringPtr(configuration.GetConfigInstance().SMS.Signature)
	/* 模板 ID: 必须填写已审核通过的模板 ID */
	r.TemplateId = common.StringPtr(configuration.GetConfigInstance().SMS.TemplateId)
	/* 模板参数: 模板参数的个数需要与 TemplateId 对应模板的变量个数保持一致，若无模板参数，则设置为空*/
	r.TemplateParamSet = common.StringPtrs([]string{code})
	/* 下发手机号码，采用 E.164 标准，+[国家或地区码][手机号]*/
	r.PhoneNumberSet = common.StringPtrs([]string{"+86" + phone_number})
	// 通过client对象调用想要访问的接口，需要传入请求对象
	_, err := client.SendSms(r)
	// 处理异常
	if _, ok := err.(*errors.TencentCloudSDKError); ok {
		log.Printf("An API error has returned: %s", err)
		return
	}
	// 非SDK异常，直接失败。实际代码中可以加入其他的处理。
	if err != nil {
		log.Printf("An Non-SDK error has returned: %s", err)
		return
	}
}
