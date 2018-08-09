package alicloud

import (
	"testing"

	"github.com/hashicorp/terraform/helper/resource"
)

func TestAccAlicloudRKVInstancesDataSource(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckAlicloudRKVInstancesDataSourceConfig,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAlicloudDataSourceID("data.alicloud_kvstore_instances.rkvs"),
					resource.TestCheckResourceAttr("data.alicloud_kvstore_instances.rkvs", "instances.#", "1"),
					resource.TestCheckResourceAttr("data.alicloud_kvstore_instances.rkvs", "instances.0.instance_class", "redis.master.small.default"),
					resource.TestCheckResourceAttr("data.alicloud_kvstore_instances.rkvs", "instances.0.name", "tf-testAccCheckAlicloudRKVInstancesDataSourceConfig"),
					resource.TestCheckResourceAttr("data.alicloud_kvstore_instances.rkvs", "instances.0.instance_type", "Redis"),
					resource.TestCheckResourceAttr("data.alicloud_kvstore_instances.rkvs", "instances.0.charge_type", string(PostPaid)),
				),
			},
		},
	})
}

const testAccCheckAlicloudRKVInstancesDataSourceConfig = `
data "alicloud_kvstore_instances" "rkvs" {
  name_regex = "${alicloud_kvstore_instance.rkv.instance_name}"
}

resource "alicloud_kvstore_instance" "rkv" {
	instance_class = "redis.master.small.default"
	instance_name  = "tf-testAccCheckAlicloudRKVInstancesDataSourceConfig"
	password       = "Test12345"
	engine_version = "2.8"
  }
`
