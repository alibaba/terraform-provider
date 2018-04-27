package alicloud

import (
	"github.com/aliyun/aliyun-tablestore-go-sdk/tablestore"
)

func getPrimaryKeyType(primaryKeyType string) tablestore.PrimaryKeyType {
	var keyType tablestore.PrimaryKeyType
	t := PrimaryKeyType(primaryKeyType)
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

func describeOtsTable(tableName string, meta interface{}) (*tablestore.DescribeTableResponse, error) {
	client := meta.(*AliyunClient).otsconn

	describeTableReq := new(tablestore.DescribeTableRequest)
	describeTableReq.TableName = tableName

	return client.DescribeTable(describeTableReq)
}

func deleteOtsTable(tableName string, meta interface{}) (bool, error) {
	client := meta.(*AliyunClient).otsconn

	deleteReq := new(tablestore.DeleteTableRequest)
	deleteReq.TableName = tableName
	_, err := client.DeleteTable(deleteReq)

	describ, _ := describeOtsTable(tableName, meta)

	if describ.TableMeta != nil {
		return false, err
	}

	return true, err
}
