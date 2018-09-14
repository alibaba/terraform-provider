package alicloud

import (
	"github.com/hashicorp/terraform/helper/resource"
	"testing"
)

func TestAccAlicloudDisksDataSource_basic(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckAlicloudDisksDataSourceBasic,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAlicloudDataSourceID("data.alicloud_disks.disks"),
					resource.TestCheckResourceAttr("data.alicloud_disks.disks", "disks.#", "1"),
					resource.TestCheckResourceAttrSet("data.alicloud_disks.disks", "disks.0.id"),
					resource.TestCheckResourceAttr("data.alicloud_disks.disks", "disks.0.name", "tf-testAccCheckAlicloudDisksDataSourceBasic"),
					resource.TestCheckResourceAttr("data.alicloud_disks.disks", "disks.0.description", "tf-testAccCheckAlicloudDisksDataSourceBasic_description"),
					resource.TestCheckResourceAttrSet("data.alicloud_disks.disks", "disks.0.availability_zone"),
					resource.TestCheckResourceAttr("data.alicloud_disks.disks", "disks.0.status", "Available"),
					resource.TestCheckResourceAttr("data.alicloud_disks.disks", "disks.0.type", "data"),
					resource.TestCheckResourceAttr("data.alicloud_disks.disks", "disks.0.category", "cloud_efficiency"),
					resource.TestCheckResourceAttr("data.alicloud_disks.disks", "disks.0.encrypted", "off"),
					resource.TestCheckResourceAttr("data.alicloud_disks.disks", "disks.0.size", "20"),
					resource.TestCheckResourceAttr("data.alicloud_disks.disks", "disks.0.image_id", ""),
					resource.TestCheckResourceAttr("data.alicloud_disks.disks", "disks.0.snapshot_id", ""),
					resource.TestCheckResourceAttr("data.alicloud_disks.disks", "disks.0.instance_id", ""),
					resource.TestCheckResourceAttrSet("data.alicloud_disks.disks", "disks.0.creation_time"),
					resource.TestCheckResourceAttr("data.alicloud_disks.disks", "disks.0.attached_time", ""),
					resource.TestCheckResourceAttr("data.alicloud_disks.disks", "disks.0.detached_time", ""),
					resource.TestCheckResourceAttrSet("data.alicloud_disks.disks", "disks.0.expiration_time"),
					resource.TestCheckResourceAttr("data.alicloud_disks.disks", "disks.0.tags.%", "2"),
					resource.TestCheckResourceAttr("data.alicloud_disks.disks", "disks.0.tags.Name", "TerraformTest"),
				),
			},
		},
	})
}

func TestAccAlicloudDisksDataSource_filterByAllFields(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckAlicloudDisksDataSourceFilterByAllFields,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAlicloudDataSourceID("data.alicloud_disks.disks"),
					resource.TestCheckResourceAttr("data.alicloud_disks.disks", "disks.#", "1"),
					resource.TestCheckResourceAttr("data.alicloud_disks.disks", "disks.0.name", "tf-testAccCheckAlicloudDisksDataSourceFilterByAllFields"),
				),
			},
		},
	})
}

const testAccCheckAlicloudDisksDataSourceBasic = `
variable "name" {
	default = "tf-testAccCheckAlicloudDisksDataSourceBasic"
}

data "alicloud_zones" "az" {
	"available_resource_creation"= "VSwitch"
}

resource "alicloud_disk" "sample_disk" {
	availability_zone = "${data.alicloud_zones.az.zones.0.id}"
	category = "cloud_efficiency"
	name = "${var.name}"
    description = "${var.name}_description"
	size = "20"
	tags {
	    Name = "TerraformTest"
	    Name1 = "TerraformTest"
	}
}

data "alicloud_disks" "disks" {
    name_regex = "${alicloud_disk.sample_disk.name}"
}
`

const testAccCheckAlicloudDisksDataSourceFilterByAllFields = `
variable "name" {
	default = "tf-testAccCheckAlicloudDisksDataSourceFilterByAllFields"
}

data "alicloud_zones" "az" {
	"available_resource_creation"= "VSwitch"
}

resource "alicloud_disk" "sample_disk" {
	availability_zone = "${data.alicloud_zones.az.zones.0.id}"
	category = "cloud_efficiency"
	name = "${var.name}"
    description = "${var.name}_description"
	size = "20"
	tags {
	    Name = "TerraformTest"
	    Name1 = "TerraformTest"
	}
}

data "alicloud_disks" "disks" {
    ids = ["${alicloud_disk.sample_disk.id}"]
    name_regex = "${alicloud_disk.sample_disk.name}"
    type = "data"
    category = "cloud_efficiency"
    encrypted = "off"
    tags = {
        Name = "TerraformTest"
    }
}
`
