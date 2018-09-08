package alicloud

import (
	"github.com/alibaba/terraform-provider/alicloud/aliyunclient"
	"github.com/aliyun/aliyun-oss-go-sdk/oss"
)

type OssService struct {
	client *aliyunclient.AliyunClient
}

func (s *OssService) QueryOssBucketById(id string) (info *oss.BucketInfo, err error) {
	raw, err := s.client.RunSafelyWithOssClient(func(ossClient *oss.Client) (interface{}, error) {
		return ossClient.GetBucketInfo(id)
	})
	if err != nil {
		return nil, err
	}
	bucket := raw.(*oss.GetBucketInfoResult)
	return &bucket.BucketInfo, nil
}
