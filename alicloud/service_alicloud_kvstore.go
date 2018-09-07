package alicloud

import (
	"github.com/alibaba/terraform-provider/alicloud/aliyunclient"
	"time"

	"github.com/aliyun/alibaba-cloud-sdk-go/services/r-kvstore"
	"github.com/denverdino/aliyungo/common"
)

func DescribeRKVInstanceById(id string, client *aliyunclient.AliyunClient) (instance *r_kvstore.DBInstanceAttribute, err error) {
	request := r_kvstore.CreateDescribeInstanceAttributeRequest()
	request.InstanceId = id
	raw, err := client.RunSafelyWithRkvClient(func(rkvClient *r_kvstore.Client) (interface{}, error) {
		return rkvClient.DescribeInstanceAttribute(request)
	})
	if err != nil {
		if IsExceptedError(err, InvalidKVStoreInstanceIdNotFound) {
			return nil, GetNotFoundErrorFromString(GetNotFoundMessage("KVStore instance", id))
		}
		return nil, err
	}
	resp := raw.(*r_kvstore.DescribeInstanceAttributeResponse)
	if resp == nil || len(resp.Instances.DBInstanceAttribute) <= 0 {
		return nil, GetNotFoundErrorFromString(GetNotFoundMessage("KVStore instance", id))
	}

	return &resp.Instances.DBInstanceAttribute[0], nil
}

func DescribeRKVInstancebackupPolicy(id string, client *aliyunclient.AliyunClient) (policy *r_kvstore.DescribeBackupPolicyResponse, err error) {
	request := r_kvstore.CreateDescribeBackupPolicyRequest()
	request.InstanceId = id
	raw, err := client.RunSafelyWithRkvClient(func(rkvClient *r_kvstore.Client) (interface{}, error) {
		return rkvClient.DescribeBackupPolicy(request)
	})
	if err != nil {
		if IsExceptedError(err, InvalidKVStoreInstanceIdNotFound) {
			return nil, GetNotFoundErrorFromString(GetNotFoundMessage("KVStore Instance Policy", id))
		}
		return nil, err
	}
	policy = raw.(*r_kvstore.DescribeBackupPolicyResponse)

	if policy == nil {
		err = GetNotFoundErrorFromString(GetNotFoundMessage("KVStore Instance Policy", id))
	}

	return
}

func WaitForRKVInstance(instanceId string, status Status, timeout int, client *aliyunclient.AliyunClient) error {
	if timeout <= 0 {
		timeout = DefaultTimeout
	}
	for {
		instance, err := DescribeRKVInstanceById(instanceId, client)
		if err != nil && !NotFoundError(err) {
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
