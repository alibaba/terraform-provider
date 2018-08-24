package alicloud

import (
	"testing"

	"github.com/hashicorp/terraform/helper/resource"
)

func TestAccAlicloudEIPsDataSource(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckAlicloudEipsDataSourceConfig,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAlicloudDataSourceID("data.alicloud_eips.foo"),
					resource.TestCheckResourceAttr("data.alicloud_eips.foo", "eips.#", "2"),
					resource.TestCheckResourceAttr("data.alicloud_eips.foo", "eips.0.status", "Available"),
					resource.TestCheckResourceAttr("data.alicloud_eips.foo", "eips.1.bandwidth", "5"),
				),
			},
		},
	})
}

func TestAccAlicloudEIPsDataSourceWithStatus(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckAlicloudEipsDataSourceWithStatusConfig,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAlicloudDataSourceID("data.alicloud_eips.availableEips"),
					resource.TestCheckResourceAttr("data.alicloud_eips.availableEips", "eips.#", "1"),
					resource.TestCheckResourceAttr("data.alicloud_eips.availableEips", "eips.0.status", "Available"),
					resource.TestCheckResourceAttr("data.alicloud_eips.availableEips", "eips.0.bandwidth", "10"),

					testAccCheckAlicloudDataSourceID("data.alicloud_eips.inUseEips"),
					resource.TestCheckResourceAttr("data.alicloud_eips.inUseEips", "eips.#", "1"),
					resource.TestCheckResourceAttr("data.alicloud_eips.inUseEips", "eips.0.status", "InUse"),
					resource.TestCheckResourceAttr("data.alicloud_eips.inUseEips", "eips.0.bandwidth", "5"),
				),
			},
		},
	})
}

const testAccCheckAlicloudEipsDataSourceConfig = `
resource "alicloud_eip" "eip" {
  count = 2
  bandwidth = 5
}

data "alicloud_eips" "foo" {
  ids = ["${alicloud_eip.eip.*.id}"]
  ip_addresses = ["${alicloud_eip.eip.*.ip_address}"]
}
`

const testAccCheckAlicloudEipsDataSourceWithStatusConfig = `
resource "alicloud_eip" "availableEip" {
  bandwidth = 10
}
resource "alicloud_eip" "inUseEip" {
  bandwidth = 5
}

data "alicloud_zones" "zones" {
  "available_resource_creation" = "VSwitch"
}
resource "alicloud_vpc" "vpc" {
  name = "testAccCheckAlicloudEipsDataSourceWithStatusConfig"
  cidr_block = "10.1.0.0/21"
}
resource "alicloud_vswitch" "vswitch" {
  vpc_id = "${alicloud_vpc.vpc.id}"
  cidr_block = "10.1.1.0/24"
  availability_zone = "${data.alicloud_zones.zones.zones.0.id}"
  name = "testAccCheckAlicloudEipsDataSourceWithStatusConfig"
}
resource "alicloud_slb" "slb" {
  name = "testAccCheckAlicloudEipsDataSourceWithStatusConfig"
  specification = "slb.s2.small"
  vswitch_id = "${alicloud_vswitch.vswitch.id}"
}
resource "alicloud_eip_association" "eipAssociation" {
  allocation_id = "${alicloud_eip.inUseEip.id}"
  instance_id = "${alicloud_slb.slb.id}"
}

data "alicloud_eips" "availableEips" {
  ids = ["${alicloud_eip.availableEip.*.id}"]
  status = "Available"
}
data "alicloud_eips" "inUseEips" {
  ids = ["${alicloud_eip_association.eipAssociation.allocation_id}"]
  status = "InUse"
}
`
