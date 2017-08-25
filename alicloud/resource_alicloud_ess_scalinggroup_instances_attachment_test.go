package alicloud

import (
	"fmt"
	"github.com/denverdino/aliyungo/common"
	"github.com/denverdino/aliyungo/ecs"
	"github.com/denverdino/aliyungo/ess"
	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/terraform"
	"strings"
	"testing"
)

func TestAccAlicloudEssScalingGroupInstancesAttachment_basic(t *testing.T) {
	var sg ess.ScalingGroupItemType
	var instance ecs.InstanceAttributesType

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},

		// module name
		IDRefreshName: "alicloud_ess_scaling_group_instances_attachment.attach",

		Providers:    testAccProviders,
		CheckDestroy: testAccCheckEssScalingGroupInstancesAttachmentDestroy,
		Steps: []resource.TestStep{
			resource.TestStep{
				Config: testAccEssScalingGroupInstancesAttachmentConfig,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckEssScalingGroupExists(
						"alicloud_ess_scaling_group.scaling", &sg),
					testAccCheckInstanceExists(
						"alicloud_instance.instance", &instance),
					testAccCheckEssScalingGroupInstancesAttachmentExists(
						"alicloud_ess_scaling_group_instances_attachment.attach", &sg, &instance),
				),
			},
		},
	})

}

func testAccCheckEssScalingGroupInstancesAttachmentExists(n string, group *ess.ScalingGroupItemType, instance *ecs.InstanceAttributesType) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("Not found: %s", n)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("No attachment ID is set")
		}

		client := testAccProvider.Meta().(*AliyunClient)
		essconn := client.essconn

		args := &ess.DescribeScalingInstancesArgs{
			RegionId:       client.Region,
			ScalingGroupId: group.ScalingGroupId,
			InstanceId:     common.FlattenArray([]string{instance.InstanceId}),
		}

		instances, _, err := essconn.DescribeScalingInstances(args)

		if err != nil {
			return err
		}

		if len(instances) == 0 {
			return fmt.Errorf("Attachment not found.")
		}

		for _, inst := range instances {
			if inst.InstanceId == instance.InstanceId {
				return nil
			}
		}

		return fmt.Errorf("Instances which attached with group not found.")
	}
}

func testAccCheckEssScalingGroupInstancesAttachmentDestroy(s *terraform.State) error {
	client := testAccProvider.Meta().(*AliyunClient)
	essconn := client.essconn

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "alicloud_ess_scaling_group_instances_attachment" {
			continue
		}

		parts := strings.Split(rs.Primary.ID, COLON_SEPARATED)
		args := &ess.DescribeScalingInstancesArgs{
			RegionId:       client.Region,
			ScalingGroupId: parts[0],
			InstanceId:     common.FlattenArray(strings.Split(parts[1], ",")),
		}

		instances, _, err := essconn.DescribeScalingInstances(args)

		if len(instances) > 0 {
			return fmt.Errorf("Error Attachment still exists.")
		}

		if err != nil {
			return err
		}
	}

	return nil
}

const testAccEssScalingGroupInstancesAttachmentConfig = `
variable "removal_policies" {
  type    = "list"
  default = ["OldestInstance", "NewestInstance"]
}

variable "ecs_instance_type" {
  default = "ecs.n4.large"
}

variable "security_group_name" {
  default = "tf-sg"
}

data "alicloud_images" "ecs_image" {
  most_recent = true
  name_regex =  "^centos_6\\w{1,5}[64].*"
}

resource "alicloud_ess_scaling_group" "scaling" {
  min_size = 1
  max_size = 5
  scaling_group_name = "tf-scaling"
  removal_policies = "${var.removal_policies}"
}

resource "alicloud_security_group" "sg" {
  name = "${var.security_group_name}"
  description = "tf-sg"
}

resource "alicloud_ess_scaling_configuration" "config" {
  scaling_group_id = "${alicloud_ess_scaling_group.scaling.id}"
  enable = true

  image_id = "${data.alicloud_images.ecs_image.images.0.id}"
  instance_type = "${var.ecs_instance_type}"
  security_group_id = "${alicloud_security_group.sg.id}"
}

resource "alicloud_instance" "instance" {
  image_id = "${data.alicloud_images.ecs_image.images.0.id}"

  system_disk_category = "cloud_ssd"
  system_disk_size = 80

  instance_type = "${var.ecs_instance_type}"
  internet_charge_type = "PayByBandwidth"
  security_groups = ["${alicloud_security_group.sg.id}"]
  instance_name = "youyouyou"

  tags {
    foo = "bar"
    work = "test"
  }
}

resource "alicloud_ess_scaling_group_instances_attachment" "attach" {
  scaling_group_id = "${alicloud_ess_scaling_group.scaling.id}"
  instance_ids = ["${alicloud_instance.instance.id}"]
}`
