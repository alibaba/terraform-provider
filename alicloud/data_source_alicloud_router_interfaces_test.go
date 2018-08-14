package alicloud

import (
	"testing"

	"github.com/hashicorp/terraform/helper/resource"
)

func TestAccAlicloudRouterInterfacesDataSource_basic(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckAlicloudRouterInterfacesDataSourceConfig,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAlicloudDataSourceID("data.alicloud_router_interfaces.router_interfaces"),
					resource.TestCheckResourceAttr("data.alicloud_router_interfaces.router_interfaces", "interfaces.#", "1"),
					resource.TestCheckResourceAttrSet("data.alicloud_router_interfaces.router_interfaces", "interfaces.0.id"),
					resource.TestCheckResourceAttr("data.alicloud_router_interfaces.router_interfaces", "interfaces.0.status", "Idle"),
					resource.TestCheckResourceAttr("data.alicloud_router_interfaces.router_interfaces", "interfaces.0.name", "testAccCheckAlicloudRouterInterfacesDataSourceConfig"),
					resource.TestCheckResourceAttr("data.alicloud_router_interfaces.router_interfaces", "interfaces.0.description", "testAccCheckAlicloudRouterInterfacesDataSourceConfig_descr"),
					resource.TestCheckResourceAttr("data.alicloud_router_interfaces.router_interfaces", "interfaces.0.role", "InitiatingSide"),
					resource.TestCheckResourceAttr("data.alicloud_router_interfaces.router_interfaces", "interfaces.0.specification", "Large.2"),
					resource.TestCheckResourceAttrSet("data.alicloud_router_interfaces.router_interfaces", "interfaces.0.router_id"),
					resource.TestCheckResourceAttr("data.alicloud_router_interfaces.router_interfaces", "interfaces.0.router_type", "VRouter"),
					resource.TestCheckResourceAttrSet("data.alicloud_router_interfaces.router_interfaces", "interfaces.0.vpc_id"),
					resource.TestCheckResourceAttr("data.alicloud_router_interfaces.router_interfaces", "interfaces.0.access_point_id", ""),
					resource.TestCheckResourceAttrSet("data.alicloud_router_interfaces.router_interfaces", "interfaces.0.creation_time"),
					resource.TestCheckResourceAttrSet("data.alicloud_router_interfaces.router_interfaces", "interfaces.0.opposite_region_id"),
					resource.TestCheckResourceAttr("data.alicloud_router_interfaces.router_interfaces", "interfaces.0.opposite_interface_id", ""),
					resource.TestCheckResourceAttr("data.alicloud_router_interfaces.router_interfaces", "interfaces.0.opposite_router_id", ""),
					resource.TestCheckResourceAttr("data.alicloud_router_interfaces.router_interfaces", "interfaces.0.opposite_router_type", "VRouter"),
					resource.TestCheckResourceAttr("data.alicloud_router_interfaces.router_interfaces", "interfaces.0.opposite_interface_owner_id", ""),
					resource.TestCheckResourceAttr("data.alicloud_router_interfaces.router_interfaces", "interfaces.0.health_check_source_ip", ""),
					resource.TestCheckResourceAttr("data.alicloud_router_interfaces.router_interfaces", "interfaces.0.health_check_target_ip", ""),
				),
			},
		},
	})
}

const testAccCheckAlicloudRouterInterfacesDataSourceConfig = `
variable "name" {
	default = "testAccCheckAlicloudRouterInterfacesDataSourceConfig"
}

resource "alicloud_vpc" "foo" {
  name = "${var.name}"
  cidr_block = "172.16.0.0/12"
}

data "alicloud_regions" "current_regions" {
  current = true
}

resource "alicloud_router_interface" "interface" {
  opposite_region = "${data.alicloud_regions.current_regions.regions.0.id}"
  router_type = "VRouter"
  router_id = "${alicloud_vpc.foo.router_id}"
  role = "InitiatingSide"
  specification = "Large.2"
  name = "${var.name}"
  description = "${var.name}_descr"
}

data "alicloud_router_interfaces" "router_interfaces" {
  router_id = "${alicloud_vpc.foo.router_id}"
  specification = "${alicloud_router_interface.interface.specification}"
}
`
