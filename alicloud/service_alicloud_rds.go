package alicloud

import (
	"github.com/denverdino/aliyungo/common"
	"github.com/denverdino/aliyungo/rds"
)

// when getInstance is empty, then throw InstanceNotfound error
func (client *AliyunClient) DescribeDBInstanceById(id string) (instance *rds.DBInstanceAttribute, err error) {
	arrtArgs := rds.DescribeDBInstancesArgs{
		DBInstanceId: id,
	}
	resp, err := client.rdsconn.DescribeDBInstanceAttribute(&arrtArgs)
	if err != nil {
		return nil, err
	}

	attr := resp.Items.DBInstanceAttribute

	if len(attr) <= 0 {
		return nil, common.GetClientErrorFromString(InstanceNotfound)
	}

	return &attr[0], nil
}
