// common functions used by 'project' 'topic' modules of datahub etc
package alicloud

import (
	"time"
)

func convUint64ToDate(t uint64) string {
	return time.Unix(int64(t), 0).Format("2006-01-02 15:04:05")
}

func getNow() string {
	return time.Now().Format("2006-01-02 15:04:05")
}
