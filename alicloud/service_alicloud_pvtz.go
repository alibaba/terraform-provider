package alicloud

import (
	"github.com/alibaba/terraform-provider/alicloud/aliyunclient"
	"strconv"

	"github.com/aliyun/alibaba-cloud-sdk-go/services/pvtz"
)

func DescribePvtzZoneInfo(zoneId string, client *aliyunclient.AliyunClient) (zone pvtz.DescribeZoneInfoResponse, err error) {
	request := pvtz.CreateDescribeZoneInfoRequest()
	request.ZoneId = zoneId

	invoker := NewInvoker()
	err = invoker.Run(func() error {
		raw, err := client.RunSafelyWithPvtzClient(func(pvtzClient *pvtz.Client) (interface{}, error) {
			return pvtzClient.DescribeZoneInfo(request)
		})
		if err != nil {
			if IsExceptedErrors(err, []string{ZoneNotExists, ZoneVpcNotExists}) {
				return GetNotFoundErrorFromString(GetNotFoundMessage("PrivateZone", zoneId))
			}
			return err
		}
		resp := raw.(*pvtz.DescribeZoneInfoResponse)
		if resp == nil || resp.ZoneId != zoneId {
			return GetNotFoundErrorFromString(GetNotFoundMessage("PrivateZone", zoneId))
		}
		zone = *resp
		return nil
	})

	return

}

func DescribeZoneRecord(recordId int, zoneId string, client *aliyunclient.AliyunClient) (record pvtz.Record, err error) {
	request := pvtz.CreateDescribeZoneRecordsRequest()
	request.ZoneId = zoneId

	invoker := NewInvoker()
	err = invoker.Run(func() error {
		raw, err := client.RunSafelyWithPvtzClient(func(pvtzClient *pvtz.Client) (interface{}, error) {
			return pvtzClient.DescribeZoneRecords(request)
		})

		recordIdStr := strconv.Itoa(recordId)

		if err != nil {
			if IsExceptedErrors(err, []string{ZoneNotExists}) {
				return GetNotFoundErrorFromString(GetNotFoundMessage("PrivateZoneRecord", recordIdStr))
			}
			return err
		}
		resp := raw.(*pvtz.DescribeZoneRecordsResponse)
		if resp == nil {
			return GetNotFoundErrorFromString(GetNotFoundMessage("PrivateZoneRecord", recordIdStr))
		}

		var found bool
		for _, rec := range resp.Records.Record {
			if rec.RecordId == recordId {
				record = rec
				found = true
			}
		}

		if found == false {
			return GetNotFoundErrorFromString(GetNotFoundMessage("PrivateZoneRecord", recordIdStr))
		}

		return nil
	})

	return
}
