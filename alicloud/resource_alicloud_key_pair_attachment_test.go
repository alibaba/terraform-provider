package alicloud

import (
	"fmt"
	"testing"

	"github.com/denverdino/aliyungo/ecs"
	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/terraform"
)

func TestAccAlicloudKeyPairAttachment_basic(t *testing.T) {
	var keypair ecs.KeyPairItemType
	var instance ecs.InstanceAttributesType

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},

		// module name
		IDRefreshName: "alicloud_key_pair_attachment.attach",
		Providers:     testAccProviders,
		CheckDestroy:  testAccCheckKeyPairAttachmentDestroy,
		Steps: []resource.TestStep{
			resource.TestStep{
				Config: testAccKeyPairAttachmentConfig,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckKeyPairExists(
						"alicloud_key_pair.key", &keypair),
					testAccCheckInstanceExists(
						"alicloud_instance.instance.0", &instance),
					testAccCheckKeyPairAttachmentExists(
						"alicloud_key_pair_attachment.attach", &instance, &keypair),
				),
			},
		},
	})

}

func testAccCheckKeyPairAttachmentExists(n string, instance *ecs.InstanceAttributesType, keypair *ecs.KeyPairItemType) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("Not found: %s", n)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("No Key Pair Attachment ID is set")
		}

		client := testAccProvider.Meta().(*AliyunClient)

		response, err := client.QueryInstancesById(instance.InstanceId)
		if err != nil {
			return fmt.Errorf("Error QueryInstancesById: %#v", err)
		}

		if response != nil && response.KeyPairName == keypair.KeyPairName {
			keypair.KeyPairName = response.KeyPairName
			instance = response
			return nil

		}
		return fmt.Errorf("Error KeyPairAttachment is not exist.")
	}
}

func testAccCheckKeyPairAttachmentDestroy(s *terraform.State) error {

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "alicloud_key_pair_attachment" {
			continue
		}
		// Try to find the Disk
		client := testAccProvider.Meta().(*AliyunClient)

		instance_ids := rs.Primary.Attributes["instance_ids"]

		for _, inst := range instance_ids {
			response, err := client.QueryInstancesById(string(inst))
			if err != nil {
				return err
			}

			if response != nil && response.KeyPairName != "" {
				return fmt.Errorf("Error Key Pair Attachment still exist")
			}

		}
	}

	return nil
}

const testAccKeyPairAttachmentConfig = `
variable "count_format" {
  default = "%02d"
}
variable "availability_zones" {
  default = "cn-beijing-d"
}

resource "alicloud_vpc" "main" {
  name = "vpc-for-keypair"
  cidr_block = "10.1.0.0/21"
}

resource "alicloud_vswitch" "main" {
  vpc_id = "${alicloud_vpc.main.id}"
  cidr_block = "10.1.1.0/24"
  availability_zone = "${var.availability_zones}"
  depends_on = [
    "alicloud_vpc.main"]
}
resource "alicloud_security_group" "group" {
  name = "test-for-keypair"
  description = "New security group"
  vpc_id = "${alicloud_vpc.main.id}"
}

resource "alicloud_instance" "instance" {
  instance_name = "test-keypair-${format(var.count_format, count.index+1)}"
  image_id = "ubuntu_140405_64_40G_cloudinit_20161115.vhd"
  instance_type = "ecs.n4.small"
  count = 2
  availability_zone = "${var.availability_zones}"
  security_groups = ["${alicloud_security_group.group.id}"]
  vswitch_id = "${alicloud_vswitch.main.id}"

  internet_charge_type = "PayByTraffic"
  internet_max_bandwidth_out = 5

  allocate_public_ip = "true"

  password = "Test12345"

  instance_charge_type = "PostPaid"
  system_disk_category = "cloud_ssd"
}

resource "alicloud_key_pair" "key" {
  key_name = "terraform-test-key-pair-attachment"
}

resource "alicloud_key_pair_attachment" "attach" {
  key_name = "${alicloud_key_pair.key.id}"
  instance_ids = ["${alicloud_instance.instance.*.id}"]
}
`
