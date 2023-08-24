package utils

import (
	"point-manage/config"
	"crypto/hmac"
	"crypto/sha1"
	"encoding/base64"
	"encoding/json"
	"github.com/wonderivan/logger"
	"hash"
	"io"
	"time"
)

type ConfigStruct struct {
	Expiration string     `json:"expiration"`
	Conditions [][]string `json:"conditions"`
}

type CallbackParam struct {
	CallbackURL      string `json:"callbackURL"`
	CallbackBody     string `json:"callbackBody"`
	CallbackBodyType string `json:"callbackBodyType"`
}
type PolicyToken struct {
	AccessKeyID string `json:"accessid"`
	Host        string `json:"host"`
	Expire      int64  `json:"expire"`
	Signature   string `json:"signature"`
	Policy      string `json:"policy"`
	Directory   string `json:"dir"`
	Callback    string `json:"callback"`
}

func GetPolicyToken(bucket string) string {
	ossBucket := config.Instance.OSSConfig[bucket]
	expireAt := time.Now().Add(time.Duration(ossBucket.ExpireTime) * time.Second).Format("2006-01-02T15:04:05Z")

	//create post policy json
	var config ConfigStruct
	config.Expiration = expireAt
	var condition []string
	condition = append(condition, "starts-with")
	condition = append(condition, "$key")
	condition = append(condition, ossBucket.UploadDir)
	config.Conditions = append(config.Conditions, condition)

	//calculate signature
	result, err := json.Marshal(config)
	bytes := base64.StdEncoding.EncodeToString(result)
	h := hmac.New(func() hash.Hash { return sha1.New() }, []byte(ossBucket.AccessKeySecret))
	_, _ = io.WriteString(h, bytes)
	signedStr := base64.StdEncoding.EncodeToString(h.Sum(nil))

	var callbackParam CallbackParam
	callbackParam.CallbackURL = ossBucket.CallbackUrl
	callbackParam.CallbackBody = "filename=${object}&size=${size}&mimeType=${mimeType}&height=${imageInfo.height}&width=${imageInfo.width}"
	callbackParam.CallbackBodyType = "application/x-www-form-urlencoded"
	callbackStr, err := json.Marshal(callbackParam)
	if err != nil {
		logger.Info("callback json err: %+v", err)
		return ""
	}
	callbackBase64 := base64.StdEncoding.EncodeToString(callbackStr)

	var policyToken PolicyToken

	policyToken.AccessKeyID = ossBucket.AccessKeyId
	policyToken.Host = ossBucket.Host
	expire, err := time.ParseInLocation("2006-01-02T15:04:05Z", expireAt, time.Local)
	if err != nil {
		logger.Info("Parse expire time failed with error: %+v", err)
		return ""
	}
	policyToken.Expire = expire.Unix()
	policyToken.Signature = signedStr
	policyToken.Directory = ossBucket.UploadDir
	policyToken.Policy = bytes
	policyToken.Callback = callbackBase64
	/////////

	response, err := json.Marshal(policyToken)
	if err != nil {
		logger.Info("json.Marshal(policyToken) failed with error: %+v", err)
		return ""
	}
	return string(response)
}
