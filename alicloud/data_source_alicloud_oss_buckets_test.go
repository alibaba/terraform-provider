package alicloud

import (
	"github.com/hashicorp/terraform/helper/resource"
	"testing"
)

func TestAccAlicloudOssBucketsDataSource_basic(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckAlicloudOssBucketsDataSourceBasic,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAlicloudDataSourceID("data.alicloud_oss_buckets.balancers"),
					resource.TestCheckResourceAttr("data.alicloud_oss_buckets.balancers", "slbs.#", "1"),
					resource.TestCheckResourceAttrSet("data.alicloud_oss_buckets.balancers", "slbs.0.id"),
				),
			},
		},
	})
}

const testAccCheckAlicloudOssBucketsDataSourceBasic = `
variable "name" {
	default = "tf-testAccCheckAlicloudOssBucketsDataSourceBasic"
}

resource "alicloud_oss_bucket" "sample_bucket" {
	bucket = "${var.name}"
	acl = "public-read"
}

data "alicloud_oss_buckets" "buckets" {
    name_regex = "${alicloud_oss_bucket.sample_bucket.name}"
}
`
