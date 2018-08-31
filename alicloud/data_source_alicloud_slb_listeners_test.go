package alicloud

import (
	"github.com/hashicorp/terraform/helper/resource"
	"testing"
)

func TestAccAlicloudSlbListenersDataSource_http(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckAlicloudSlbListenersDataSourceHttp,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAlicloudDataSourceID("data.alicloud_slb_listeners.slb_listeners"),
					resource.TestCheckResourceAttr("data.alicloud_slb_listeners.slb_listeners", "slb_listeners.#", "1"),
					resource.TestCheckResourceAttr("data.alicloud_slb_listeners.slb_listeners", "slb_listeners.0.frontend_port", "80"),
					resource.TestCheckResourceAttr("data.alicloud_slb_listeners.slb_listeners", "slb_listeners.0.protocol", "http"),

					testAccCheckAlicloudDataSourceID("data.alicloud_slb_listeners.slb_listeners_with_filters"),
					resource.TestCheckResourceAttr("data.alicloud_slb_listeners.slb_listeners_with_filters", "slb_listeners.#", "1"),
				),
			},
		},
	})
}

const testAccCheckAlicloudSlbListenersDataSourceHttp = `
variable "name" {
	default = "testAccCheckAlicloudSlbListenersDataSourceHttp"
}

data "alicloud_zones" "az" {
	"available_resource_creation"= "VSwitch"
}

resource "alicloud_vpc" "sample_vpc" {
  name = "${var.name}"
  cidr_block = "172.16.0.0/12"
}

resource "alicloud_vswitch" "sample_vswitch" {
  vpc_id = "${alicloud_vpc.sample_vpc.id}"
  cidr_block = "172.16.0.0/16"
  availability_zone = "${data.alicloud_zones.az.zones.0.id}"
}

resource "alicloud_slb" "sample_slb" {
  name = "${var.name}"
  vswitch_id = "${alicloud_vswitch.sample_vswitch.id}"
}

resource "alicloud_slb_listener" "sample_slb_listener" {
  load_balancer_id = "${alicloud_slb.sample_slb.id}"
  backend_port = 80
  frontend_port = 80
  protocol = "http"
  sticky_session = "on"
  sticky_session_type = "insert"
  cookie = "${var.name}"
  cookie_timeout = 86400
  health_check = "on"
  health_check_uri = "/cons"
  health_check_connect_port = 20
  healthy_threshold = 8
  unhealthy_threshold = 8
  health_check_timeout = 8
  health_check_interval = 5
  health_check_http_code = "http_2xx,http_3xx"
  bandwidth = 10
  x_forwarded_for = {
    retrive_slb_ip = true
    retrive_slb_id = true
  }
}

data "alicloud_slb_listeners" "slb_listeners" {
  load_balancer_id = "${alicloud_slb_listener.sample_slb_listener.load_balancer_id}"
}

data "alicloud_slb_listeners" "slb_listeners_with_filters" {
  load_balancer_id = "${alicloud_slb_listener.sample_slb_listener.load_balancer_id}"
  frontend_port = 80
  protocol = "http"
}
`
