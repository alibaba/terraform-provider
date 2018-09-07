package alicloud

import (
	"time"

	"strconv"

	"github.com/aliyun/alibaba-cloud-sdk-go/sdk/requests"
	"github.com/aliyun/alibaba-cloud-sdk-go/services/cms"
	"github.com/alibaba/terraform-provider/alicloud/aliyunclient"
)

func BuildCmsCommonRequest(region string) *requests.CommonRequest {

	request := requests.NewCommonRequest()

	return request
}

func BuildCmsAlarmRequest(id string, client *aliyunclient.AliyunClient) *requests.CommonRequest {

	request := BuildCmsCommonRequest(client.RegionId)
	request.QueryParams["Id"] = id

	return request
}

func DescribeAlarm(id string, client *aliyunclient.AliyunClient) (alarm cms.AlarmInListAlarm, err error) {

	request := cms.CreateListAlarmRequest()

	request.Id = id
	rawResponse, err := client.RunSafelyWithCmsClient(func(cmsClient *cms.Client) (interface{}, error) {
		return cmsClient.ListAlarm(request)
	})
	if err != nil {
		return alarm, err
	}
	response := rawResponse.(*cms.ListAlarmResponse)

	if len(response.AlarmList.Alarm) < 1 {
		return alarm, GetNotFoundErrorFromString(GetNotFoundMessage("Alarm Rule", id))
	}

	return response.AlarmList.Alarm[0], nil
}

func WaitForCmsAlarm(id string, enabled bool, timeout int, client *aliyunclient.AliyunClient) error {
	if timeout <= 0 {
		timeout = DefaultTimeout
	}

	for {
		alarm, err := DescribeAlarm(id, client)
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
