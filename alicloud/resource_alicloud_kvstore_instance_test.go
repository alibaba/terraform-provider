package alicloud

import (
	"fmt"
	"log"
	"testing"

	"github.com/aliyun/alibaba-cloud-sdk-go/services/r-kvstore"
	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/terraform"
)

func TestAccAlicloudKVStoreInstance_basic(t *testing.T) {
	var instance r_kvstore.DBInstanceAttribute

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},

		// module name
		IDRefreshName: "alicloud_kvstore_instance.foo",

		Providers:    testAccProviders,
		CheckDestroy: testAccCheckKVStoreInstanceDestroy,
		Steps: []resource.TestStep{
			resource.TestStep{
				Config: testAccKVStoreInstanceConfig,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckKVStoreInstanceExists(
						"alicloud_kvstore_instance.foo", &instance),
					resource.TestCheckResourceAttr(
						"alicloud_kvstore_instance.foo",
						"instance_class",
						"redis.master.small.default"),
					resource.TestCheckResourceAttr(
						"alicloud_kvstore_instance.foo",
						"engine_version",
						"2.8"),
					resource.TestCheckResourceAttr(
						"alicloud_kvstore_instance.foo",
						"instance_type",
						"Redis"),
					resource.TestCheckResourceAttr(
						"alicloud_kvstore_instance.foo",
						"instance_name",
						"testAccKVStoreInstanceConfig"),
				),
			},
		},
	})
}

func TestAccAlicloudKVStoreInstance_vpc(t *testing.T) {
	var instance r_kvstore.DBInstanceAttribute

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},

		// module name
		IDRefreshName: "alicloud_kvstore_instance.foo",

		Providers:    testAccProviders,
		CheckDestroy: testAccCheckKVStoreInstanceDestroy,
		Steps: []resource.TestStep{
			resource.TestStep{
				Config: testAccKVStoreInstance_vpc,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckKVStoreInstanceExists(
						"alicloud_kvstore_instance.foo", &instance),
					resource.TestCheckResourceAttr(
						"alicloud_kvstore_instance.foo",
						"instance_class",
						"redis.master.small.default"),
					resource.TestCheckResourceAttr(
						"alicloud_kvstore_instance.foo",
						"engine_version",
						"2.8"),
					resource.TestCheckResourceAttr(
						"alicloud_kvstore_instance.foo",
						"instance_type",
						"Redis"),
					resource.TestCheckResourceAttr(
						"alicloud_kvstore_instance.foo",
						"instance_name",
						"testAccKVStoreInstance_vpc"),
				),
			},
		},
	})

}

func TestAccAlicloudKVStoreInstance_upgradeClass(t *testing.T) {
	var instance r_kvstore.DBInstanceAttribute

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},

		// module name
		IDRefreshName: "alicloud_kvstore_instance.foo",

		Providers:    testAccProviders,
		CheckDestroy: testAccCheckKVStoreInstanceDestroy,
		Steps: []resource.TestStep{
			resource.TestStep{
				Config: testAccKVStoreInstance_class,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckKVStoreInstanceExists(
						"alicloud_kvstore_instance.foo", &instance),
					resource.TestCheckResourceAttr("alicloud_kvstore_instance.foo", "instance_class", "redis.master.small.default"),
				),
			},

			resource.TestStep{
				Config: testAccKVStoreInstance_classUpgrade,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckKVStoreInstanceExists(
						"alicloud_kvstore_instance.foo", &instance),
					resource.TestCheckResourceAttr("alicloud_kvstore_instance.foo", "instance_class", "redis.master.mid.default"),
				),
			},
		},
	})

}

func testAccCheckKVStoreInstanceExists(n string, d *r_kvstore.DBInstanceAttribute) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("Not found: %s", n)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("No DB Instance ID is set")
		}

		client := testAccProvider.Meta().(*AliyunClient)
		attr, err := client.DescribeRKVInstanceById(rs.Primary.ID)
		log.Printf("[DEBUG] check instance %s attribute %#v", rs.Primary.ID, attr)

		if err != nil {
			return err
		}

		if attr == nil {
			return fmt.Errorf("DB Instance not found")
		}

		*d = *attr
		return nil
	}
}

func testAccCheckKVStoreInstanceDestroy(s *terraform.State) error {
	client := testAccProvider.Meta().(*AliyunClient)

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "alicloud_kvstore_instance" {
			continue
		}

		ins, err := client.DescribeRKVInstanceById(rs.Primary.ID)

		if ins != nil {
			return fmt.Errorf("Error DB Instance still exist")
		}

		// Verify the error is what we want
		if err != nil {
			if NotFoundError(err) || IsExceptedError(err, InvalidKVStoreInstanceIdNotFound) {
				continue
			}
			return err
		}
	}

	return nil
}

const testAccKVStoreInstanceConfig = `
resource "alicloud_kvstore_instance" "foo" {
	instance_class = "redis.master.small.default"
	instance_name  = "testAccKVStoreInstanceConfig"
	password       = "Test12345"
	engine_version = "2.8"
  }
`

const testAccKVStoreInstance_vpc = `
data "alicloud_zones" "default" {
	available_resource_creation = "KVStore"
}
variable "name" {
	default = "testAccKVStoreInstance_vpc"
}
resource "alicloud_vpc" "foo" {
	name = "${var.name}"
	cidr_block = "172.16.0.0/12"
}

resource "alicloud_vswitch" "foo" {
 	vpc_id = "${alicloud_vpc.foo.id}"
 	cidr_block = "172.16.0.0/21"
 	availability_zone = "${data.alicloud_zones.default.zones.0.id}"
}

resource "alicloud_kvstore_instance" "foo" {
	instance_class = "redis.master.small.default"
	instance_name  = "${var.name}"
	password       = "Test12345"
	vswitch_id     = "${alicloud_vswitch.foo.id}"
}
`
const testAccKVStoreInstance_class = `
variable "name" {
	default = "testAccKVStoreInstance_class"
}
resource "alicloud_kvstore_instance" "foo" {
	instance_class = "redis.master.small.default"
	instance_name  = "${var.name}"
	password       = "Test12345"
}
`
const testAccKVStoreInstance_classUpgrade = `
variable "name" {
	default = "testAccKVStoreInstance_class"
}
resource "alicloud_kvstore_instance" "foo" {
	instance_class = "redis.master.mid.default"	
	instance_name  = "${var.name}"
	password       = "Test12345"
}
`
