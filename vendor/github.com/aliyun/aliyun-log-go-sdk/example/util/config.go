package util

import sls "github.com/aliyun/aliyun-log-go-sdk"

// Project define Project for test
var Client = &sls.Client{
	Endpoint:        "cn-hangzhou.log.aliyuncs.com",
	AccessKeyID:     "",
	AccessKeySecret: "",
}

var ProjectName = "test-project"
