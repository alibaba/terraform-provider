// common functions used by 'project' 'topic' modules of datahub etc
package alicloud

import (
	"time"

	"github.com/aliyun/aliyun-datahub-sdk-go/datahub/models"
	"github.com/aliyun/aliyun-datahub-sdk-go/datahub/types"
)

func convUint64ToDate(t uint64) string {
	return time.Unix(int64(t), 0).Format("2006-01-02 15:04:05")
}

func getNow() string {
	return time.Now().Format("2006-01-02 15:04:05")
}

func getRecordSchema(typeMap map[string]interface{}) (recordSchema *models.RecordSchema) {
	recordSchema = models.NewRecordSchema()

	for k, v := range typeMap {
		recordSchema.AddField(models.Field{Name: string(k), Type: types.FieldType(v.(string))})
	}

	return recordSchema
}
