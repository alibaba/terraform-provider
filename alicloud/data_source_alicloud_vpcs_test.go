package alicloud

import (
	"github.com/hashicorp/terraform/helper/resource"
	"testing"
)

func TestAccAlicloudVpcsDataSource_cidr_block(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckAlicloudVpcsDataSourceCidrBlockConfig,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAlicloudDataSourceID("data.alicloud_vpcs.vpc"),
					resource.TestCheckResourceAttr("data.alicloud_vpcs.vpc", "vpcs.#", "2"),
					resource.TestCheckResourceAttr("data.alicloud_vpcs.vpc", "vpcs.0.id", "vpc-2zef7qk11rgq3w053zxrr"),
					resource.TestCheckResourceAttr("data.alicloud_vpcs.vpc", "vpcs.0.region_id", "vpc-cn-beijing"),
					resource.TestCheckResourceAttr("data.alicloud_vpcs.vpc", "vpcs.0.status", "Available"),
					resource.TestCheckResourceAttr("data.alicloud_vpcs.vpc", "vpcs.0.vpc_name", "tf_test_foo"),
					resource.TestCheckResourceAttr("data.alicloud_vpcs.vpc", "is_default", "false"),
					resource.TestCheckResourceAttr("data.alicloud_vpcs.vpc", "vrouter_id", "vrt-2ze1hlrw039lr1ieyutvc"),
				),
			},
		},
	})
}
func TestAccAlicloudVpcsDataSource_name_regex(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckAlicloudVpcsDataSourceNameRegexConfig,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAlicloudDataSourceID("data.alicloud_vpcs.vpc"),
					resource.TestCheckResourceAttr("data.alicloud_vpcs.vpc", "vpcs.#", "4"),
				),
			},
		},
	})
}
func TestAccAlicloudVpcsDataSource_vswitch_id(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckAlicloudVpcsDataSourceVswitchIdConfig,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAlicloudDataSourceID("data.alicloud_vpcs.vpc"),
					resource.TestCheckResourceAttr("data.alicloud_vpcs.vpc", "vpcs.#", "1"),
				),
			},
		},
	})
}

const testAccCheckAlicloudVpcsDataSourceCidrBlockConfig = `
data "alicloud_vpcs" "vpc" {
cidr_block="172.16.0.0/12"
}
`
const testAccCheckAlicloudVpcsDataSourceNameRegexConfig = `
data "alicloud_vpcs" "vpc" {
name_regex="^tf"
}
`
const testAccCheckAlicloudVpcsDataSourceVswitchIdConfig = `
data "alicloud_vpcs" "vpc" {
vswitch_id="vsw-2zeay0imlgcn0et3ouizk"
}
`
