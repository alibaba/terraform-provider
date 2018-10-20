package alicloud

import (
	"github.com/alibaba/terraform-provider/alicloud/connectivity"
	"github.com/aliyun/aliyun-oss-go-sdk/oss"
)

type OssService struct {
	client *connectivity.AliyunClient
}

func (s *OssService) QueryOssBucketById(id string) (info *oss.BucketInfo, err error) {
	raw, err := s.client.WithOssClient(func(ossClient *oss.Client) (interface{}, error) {
		return ossClient.GetBucketInfo(id)
	})
	if err != nil {
		return nil, err
	}
	bucket, _ := raw.(oss.GetBucketInfoResult)
	return &bucket.BucketInfo, nil
}
