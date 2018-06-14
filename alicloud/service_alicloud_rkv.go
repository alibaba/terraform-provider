package alicloud

import (
	"fmt"
	"time"

	"github.com/aliyun/alibaba-cloud-sdk-go/services/r-kvstore"
	"github.com/denverdino/aliyungo/common"
)

func (client *AliyunClient) DescribeRKVInstanceById(id string) (instance *r_kvstore.DBInstanceAttribute, err error) {
	request := r_kvstore.CreateDescribeInstanceAttributeRequest()
	request.InstanceId = id
	resp, err := client.rkvconn.DescribeInstanceAttribute(request)
	if err != nil {
		return nil, err
	}

	attr := resp.Instances.DBInstanceAttribute

	if len(attr) <= 0 {
		return nil, GetNotFoundErrorFromString(fmt.Sprintf("DB instance %s is not found.", id))
	}

	return &attr[0], nil
}

func (client *AliyunClient) WaitForRKVInstance(instanceId string, status Status, timeout int) error {
	if timeout <= 0 {
		timeout = DefaultTimeout
	}
	for {
		instance, err := client.DescribeRKVInstanceById(instanceId)
		if err != nil && !NotFoundError(err) && !IsExceptedError(err, InvalidRKVInstanceIdNotFound) {
			return err
		}

		if instance != nil && instance.InstanceStatus == string(status) {
			break
		}

		if timeout <= 0 {
			return common.GetClientErrorFromString("Timeout")
		}

		timeout = timeout - DefaultIntervalMedium
		time.Sleep(DefaultIntervalMedium * time.Second)
	}
	return nil
}

func NotFoundRKVInstance(err error) bool {
	if NotFoundError(err) || IsExceptedError(err, InvalidRKVInstanceIdNotFound) {
		return true
	}

	return false
}
