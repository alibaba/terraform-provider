package alicloud

import (
	"github.com/denverdino/aliyungo/common"
	"github.com/hashicorp/terraform/helper/schema"
)

type InstanceNetWork string

const (
	ClassicNet = InstanceNetWork("Classic")
	VpcNet = InstanceNetWork("Vpc")
)

const defaultTimeout = 120

func getRegion(d *schema.ResourceData, meta interface{}) common.Region {
	return meta.(*AliyunClient).Region
}

func notFoundError(err error) bool {
	if e, ok := err.(*common.Error); ok && (e.StatusCode == 404 || e.ErrorResponse.Message == "Not found") {
		return true
	}

	return false
}
