package alicloud

import (
	"github.com/alibaba/terraform-provider/alicloud/aliyunclient"
	"github.com/aliyun/aliyun-oss-go-sdk/oss"
)

func QueryOssBucketById(id string, client *aliyunclient.AliyunClient) (info *oss.BucketInfo, err error) {
	raw, err := client.RunSafelyWithOssClient(func(ossClient *oss.Client) (interface{}, error) {
		return ossClient.GetBucketInfo(id)
	})
	if err != nil {
		return nil, err
	}
	bucket := raw.(*oss.GetBucketInfoResult)
	return &bucket.BucketInfo, nil
}
