package alicloud

import (
	"github.com/denverdino/aliyungo/slb"
	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/terraform"
	"log"
	"testing"
)

func TestAccAlicloudSlbAttachment_basic(t *testing.T) {
	var slb slb.LoadBalancerType

	testCheckAttr := func() resource.TestCheckFunc {
		return func(*terraform.State) error {
			log.Printf("testCheckAttr slb BackendServers is: %s", slb.BackendServers)
			return nil
		}
	}

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},

		// module name
		IDRefreshName: "alicloud_slb_attachment.foo",
		Providers:     testAccProviders,
		CheckDestroy:  testAccCheckSlbDestroy,
		Steps: []resource.TestStep{
			//test internet_charge_type is paybybandwidth
			resource.TestStep{
				Config: testAccSlbAttachment,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckSlbExists("alicloud_slb_attachment.foo", &slb),
					testCheckAttr(),
					//resource.TestCheckResourceAttr(
					//	"alicloud_slb_attachment.foo", "internet_charge_type", "paybybandwidth"),
				),
			},
			resource.TestStep{
				Config: testAccSlbAttachment,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckSlbExists("alicloud_slb_attachment.update", &slb),
					testCheckAttr(),
					//resource.TestCheckResourceAttr(
					//	"alicloud_slb_attachment.foo", "internet_charge_type", "paybybandwidth"),
				),
			},
		},
	})
}

const testAccSlbAttachment = `
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
	internet_charge_type = "PayByBandwidth"
	internet_max_bandwidth_out = "5"
	system_disk_category = "cloud_efficiency"

	security_group_id = "${alicloud_security_group.foo.id}"
	instance_name = "test_foo"
}

resource "alicloud_slb" "foo" {
	name = "tf_test_slb_bind"
	internet_charge_type = "paybybandwidth"
	bandwidth = "5"
	internet = "true"
}

resource "alicloud_slb_attachment" "foo" {
	slb_id = "${alicloud_slb.foo.id}"
	//slb_id = "lb-2ze5bnhqq9q3ubbosun2b"
	instances = ["${alicloud_instance.foo.id}"]
	//instances = ["i-2ze2o5ndq3w3wwy9v52j"]
}

resource "alicloud_slb_attachment" "update" {
	slb_id = "${alicloud_slb.foo.id}"
	//slb_id = "lb-2ze5bnhqq9q3ubbosun2b"
	instances = ["${alicloud_instance.foo.id}"]
	//instances = ["i-25kojm48j"]
}
`
