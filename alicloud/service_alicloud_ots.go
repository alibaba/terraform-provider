package alicloud

import (
	"fmt"
	"github.com/alibaba/terraform-provider/alicloud/aliyunclient"
	"strings"

	"time"

	"github.com/aliyun/alibaba-cloud-sdk-go/services/ots"
	"github.com/aliyun/aliyun-tablestore-go-sdk/tablestore"
)

func buildTableClient(instanceName string, client *aliyunclient.AliyunClient) *tablestore.TableStoreClient {
	endpoint := LoadEndpoint(client.RegionId, OTSCode)
	if endpoint == "" {
		endpoint = fmt.Sprintf("%s.%s.ots.aliyuncs.com", instanceName, client.RegionId)
	}
	if !strings.HasPrefix(endpoint, string(Https)) && !strings.HasPrefix(endpoint, string(Http)) {
		endpoint = fmt.Sprintf("%s://%s", Https, endpoint)
	}
	return tablestore.NewClient(endpoint, instanceName, client.AccessKey, client.SecretKey)
}

func getPrimaryKeyType(primaryKeyType string) tablestore.PrimaryKeyType {
	var keyType tablestore.PrimaryKeyType
	t := PrimaryKeyTypeString(primaryKeyType)
	switch t {
	case IntegerType:
		keyType = tablestore.PrimaryKeyType_INTEGER
	case StringType:
		keyType = tablestore.PrimaryKeyType_STRING
	case BinaryType:
		keyType = tablestore.PrimaryKeyType_BINARY
	}
	return keyType
}

func DescribeOtsTable(instanceName, tableName string, client *aliyunclient.AliyunClient) (table *tablestore.DescribeTableResponse, err error) {
	describeTableReq := new(tablestore.DescribeTableRequest)
	describeTableReq.TableName = tableName

	table, err = buildTableClient(instanceName, client).DescribeTable(describeTableReq)
	if err != nil {
		if strings.HasPrefix(err.Error(), OTSObjectNotExist) {
			err = GetNotFoundErrorFromString(GetNotFoundMessage("OTS Table", tableName))
		}
		return
	}
	if table == nil || table.TableMeta == nil || table.TableMeta.TableName != tableName {
		err = GetNotFoundErrorFromString(GetNotFoundMessage("OTS Table", tableName))
	}
	return
}

func DeleteOtsTable(instanceName, tableName string, client *aliyunclient.AliyunClient) (bool, error) {

	deleteReq := new(tablestore.DeleteTableRequest)
	deleteReq.TableName = tableName
	if _, err := buildTableClient(instanceName, client).DeleteTable(deleteReq); err != nil {
		if NotFoundError(err) {
			return true, nil
		}
		return false, err
	}

	describ, err := DescribeOtsTable(instanceName, tableName, client)

	if err != nil {
		if NotFoundError(err) {
			return true, nil
		}
		return false, err
	}

	if describ.TableMeta != nil {
		return false, err
	}

	return true, err
}

// Convert tablestore.PrimaryKeyType to PrimaryKeyTypeString
func convertPrimaryKeyType(t tablestore.PrimaryKeyType) PrimaryKeyTypeString {
	var typeString PrimaryKeyTypeString
	switch t {
	case tablestore.PrimaryKeyType_INTEGER:
		typeString = IntegerType
	case tablestore.PrimaryKeyType_BINARY:
		typeString = BinaryType
	case tablestore.PrimaryKeyType_STRING:
		typeString = StringType
	}
	return typeString
}

func DescribeOtsInstance(name string, client *aliyunclient.AliyunClient) (inst ots.InstanceInfo, err error) {
	req := ots.CreateGetInstanceRequest()
	req.InstanceName = name
	req.Method = "GET"
	raw, err := client.RunSafelyWithOtsClient(func(otsClient *ots.Client) (interface{}, error) {
		return otsClient.GetInstance(req)
	})

	// OTS instance not found error code is "NotFound"
	if err != nil {
		return
	}
	resp := raw.(*ots.GetInstanceResponse)
	if resp == nil || resp.InstanceInfo.InstanceName != name {
		return inst, GetNotFoundErrorFromString(GetNotFoundMessage("OTS Instance", name))
	}
	return resp.InstanceInfo, nil
}

func DescribeOtsInstanceVpc(name string, client *aliyunclient.AliyunClient) (inst ots.VpcInfo, err error) {
	req := ots.CreateListVpcInfoByInstanceRequest()
	req.Method = "GET"
	req.InstanceName = name
	raw, err := client.RunSafelyWithOtsClient(func(otsClient *ots.Client) (interface{}, error) {
		return otsClient.ListVpcInfoByInstance(req)
	})
	if err != nil {
		return inst, err
	}
	resp := raw.(*ots.ListVpcInfoByInstanceResponse)
	if resp == nil || resp.TotalCount < 1 {
		return inst, GetNotFoundErrorFromString(GetNotFoundMessage("OTS Instance VPC", name))
	}
	return resp.VpcInfos.VpcInfo[0], nil
}

func WaitForOtsInstance(name string, status Status, timeout int, client *aliyunclient.AliyunClient) error {
	if timeout <= 0 {
		timeout = DefaultTimeout
	}

	for {
		inst, err := DescribeOtsInstance(name, client)
		if err != nil {
			return err
		}

		if inst.Status == convertOtsInstanceStatus(status) {
			break
		}
		timeout = timeout - DefaultIntervalShort
		if timeout <= 0 {
			return GetTimeErrorFromString(GetTimeoutMessage("OTS Instance", string(status)))
		}
		time.Sleep(DefaultIntervalShort * time.Second)
	}
	return nil
}
