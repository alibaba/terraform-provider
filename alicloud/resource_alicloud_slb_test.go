package alicloud

import (
	"testing"
	"github.com/hashicorp/terraform/helper/resource"
	"fmt"
	"github.com/hashicorp/terraform/terraform"
	"github.com/denverdino/aliyungo/slb"
	"log"
)

func TestAccAlicloudSlb_basic(t *testing.T) {
	var slb slb.LoadBalancerType

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},

		// module name
		IDRefreshName: "alicloud_slb.classic",
		Providers:     testAccProviders,
		CheckDestroy:  testAccCheckSlbDestroy,
		Steps: []resource.TestStep{
			//test internet_charge_type is paybybandwidth
			resource.TestStep{
				Config: testAccSlbConfig,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckSlbExists("alicloud_slb.classic_paybybandwidth", &slb),
					resource.TestCheckResourceAttr(
						"alicloud_slb.classic_paybybandwidth", "name", "tf_test_slb_classic"),
				),
			},

			//test internet_charge_type is paybytraffic
			resource.TestStep{
				Config: testAccSlbConfig,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckSlbExists("alicloud_slb.classic_paybytraffic", &slb),
					resource.TestCheckResourceAttr(
						"alicloud_slb.classic_paybytraffic", "name", "tf_test_slb_classic"),
				),
			},
		},
	})
}

func TestAccAlicloudSlb_listener(t *testing.T) {
	var slb slb.LoadBalancerType

	testListener := func() resource.TestCheckFunc {
		return func(*terraform.State) error {
			listenerPorts := slb.ListenerPorts.ListenerPort[0]
			log.Printf("[WARN] get listenerPorts %#v", listenerPorts)
			if listenerPorts != 3376 {
				return fmt.Errorf("bad loadbalancer listener: %s", listenerPorts)
			}

			return nil
		}
	}

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},

		// module name
		IDRefreshName: "alicloud_slb.listener",
		Providers:     testAccProviders,
		CheckDestroy:  testAccCheckSlbDestroy,
		Steps: []resource.TestStep{
			resource.TestStep{
				Config: testAccSlbListener,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckSlbExists("alicloud_slb.listener", &slb),
					resource.TestCheckResourceAttr(
						"alicloud_slb.listener", "name", "tf_test_slb"),
					testListener(),
				),
			},
		},
	})
}

func TestAccAlicloudSlb_vpc(t *testing.T) {
	var slb slb.LoadBalancerType

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},

		// module name
		IDRefreshName: "alicloud_slb.vpc",
		Providers:     testAccProviders,
		CheckDestroy:  testAccCheckSlbDestroy,
		Steps: []resource.TestStep{
			resource.TestStep{
				Config: testAccSlb4VpcConfig,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckSlbExists("alicloud_slb.vpc", &slb),
					resource.TestCheckResourceAttr(
						"alicloud_slb.vpc", "name", "tf_test_slb_vpc"),
				),
			},
		},
	})
}

func TestAccAlicloudSlb_bindECS(t *testing.T) {
	var slb slb.LoadBalancerType

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},

		// module name
		IDRefreshName: "alicloud_slb.bindecs",
		Providers:     testAccProviders,
		CheckDestroy:  testAccCheckSlbDestroy,
		Steps: []resource.TestStep{
			resource.TestStep{
				Config: testAccSlbBindECS,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckSlbExists("alicloud_slb.backendservice", &slb),
					resource.TestCheckResourceAttr(
						"alicloud_slb.backendservice", "name", "tf_test_slb_bindecs"),
				),
			},
		},
	})
}

func testAccCheckSlbExists(n string, slb *slb.LoadBalancerType) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("Not found: %s", n)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("No SLB ID is set")
		}

		client := testAccProvider.Meta().(*AliyunClient)
		instance, err := client.DescribeLoadBalancerAttribute(rs.Primary.ID)

		if err != nil {
			return err
		}
		if instance == nil {
			return fmt.Errorf("SLB not found")
		}

		*slb = *instance
		return nil
	}
}

func testAccCheckSlbDestroy(s *terraform.State) error {
	client := testAccProvider.Meta().(*AliyunClient)

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "alicloud_slb" {
			continue
		}

		// Try to find the Slb
		instance, err := client.DescribeLoadBalancerAttribute(rs.Primary.ID)

		if instance != nil {
			return fmt.Errorf("SLB still exist")
		}

		if err != nil {
			// Verify the error is what we want
			return err
		}

	}

	return nil
}

const testAccSlbConfig = `
resource "alicloud_slb" "classic_paybybandwidth" {
  name = "tf_test_slb_classic"
  internet_charge_type = "paybybandwidth"
  bandwidth = "5"
  internet = "true"
}

resource "alicloud_slb" "classic_paybytraffic" {
  name = "tf_test_slb_classic"
  internet_charge_type = "paybytraffic"
  bandwidth = "5"
  internet = "true"
}
`

const testAccSlbListener = `
resource "alicloud_slb" "listener" {
  name = "tf_test_slb"
  internet_charge_type = "paybybandwidth"
  bandwidth = "5"
  internet = "true"
  listener = [
    {
      "instance_port" = "2375"
      "instance_protocol" = "tcp"
      "lb_port" = "3376"
      "lb_protocol" = "tcp"
      "bandwidth" = "5"
    }]
}
`

const testAccSlb4VpcConfig = `
resource "alicloud_vpc" "foo" {
  name = "tf_test_foo"
  cidr_block = "172.16.0.0/12"
}

resource "alicloud_vswitch" "foo" {
  vpc_id = "${alicloud_vpc.foo.id}"
  cidr_block = "172.16.0.0/21"
  availability_zone = "cn-beijing-b"
}

resource "alicloud_slb" "vpc" {
  name = "tf_test_slb_vpc"
  internet_charge_type = "paybybandwidth"
  bandwidth = "5"
  internet = "true"
  vswitch_id = "${alicloud_vswitch.foo.id}"
}
`

const testAccSlbBindECS = `
resource "alicloud_security_group" "foo" {
	name = "tf_test_foo"
	description = "foo"
}

resource "alicloud_instance" "foo" {
	# cn-beijing
	availability_zone = "cn-beijing-b"
	image_id = "ubuntu1404_64_40G_cloudinit_20160727.raw"

	# series II
	instance_type = "ecs.n1.medium"
	io_optimized = "optimized"
	internet_charge_type = "PayByBandwidth"
	system_disk_category = "cloud_efficiency"

	security_group_id = "${alicloud_security_group.foo.id}"
	instance_name = "test_foo"

}

resource "alicloud_slb" "backendservice" {
  name = "tf_test_slb_bindecs"
  internet_charge_type = "paybybandwidth"
  bandwidth = "5"
  internet = "true"
  instances = ["${alicloud_instance.foo.id}"]
}
`

