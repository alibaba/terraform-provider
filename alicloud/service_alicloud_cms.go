package alicloud

import (
	"time"

	"strconv"

	"github.com/alibaba/terraform-provider/alicloud/connectivity"
	"github.com/aliyun/alibaba-cloud-sdk-go/sdk/requests"
	"github.com/aliyun/alibaba-cloud-sdk-go/services/cms"
)

type CmsService struct {
	client *connectivity.AliyunClient
}

func (s *CmsService) BuildCmsCommonRequest(region string) *requests.CommonRequest {
	request := requests.NewCommonRequest()
	return request
}

func (s *CmsService) BuildCmsAlarmRequest(id string) *requests.CommonRequest {

	request := s.BuildCmsCommonRequest(s.client.RegionId)
	request.QueryParams["Id"] = id

	return request
}

func (s *CmsService) DescribeAlarm(id string) (alarm cms.AlarmInListAlarm, err error) {

	request := cms.CreateListAlarmRequest()

	request.Id = id
	raw, err := s.client.RunSafelyWithCmsClient(func(cmsClient *cms.Client) (interface{}, error) {
		return cmsClient.ListAlarm(request)
	})
	if err != nil {
		return alarm, err
	}
	response, _ := raw.(*cms.ListAlarmResponse)

	if len(response.AlarmList.Alarm) < 1 {
		return alarm, GetNotFoundErrorFromString(GetNotFoundMessage("Alarm Rule", id))
	}

	return response.AlarmList.Alarm[0], nil
}

func (s *CmsService) WaitForCmsAlarm(id string, enabled bool, timeout int) error {
	if timeout <= 0 {
		timeout = DefaultTimeout
	}

	for {
		alarm, err := s.DescribeAlarm(id)
		if err != nil {
			return err
		}

		if alarm.Enable == enabled {
			break
		}
		timeout = timeout - DefaultIntervalShort
		if timeout <= 0 {
			return GetTimeErrorFromString(GetTimeoutMessage("Alarm", strconv.FormatBool(enabled)))
		}
		time.Sleep(DefaultIntervalShort * time.Second)
	}
	return nil
}
