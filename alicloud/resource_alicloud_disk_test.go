package alicloud

import (
	"fmt"
	"testing"

	"github.com/denverdino/aliyungo/ecs"
	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/terraform"
	"log"
)

func TestAccAlicloudDisk_basic(t *testing.T) {
	var v ecs.DiskItemType
	// todo: create disk

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},

		// module name
		IDRefreshName: "alicloud_disk.foo",

		Providers:    testAccProviders,
		CheckDestroy: testAccCheckDiskDestroy,
		Steps: []resource.TestStep{
			resource.TestStep{
				Config: testAccDiskConfig,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckDiskExists(
						"alicloud_disk.foo", &v),
					resource.TestCheckResourceAttr(
						"alicloud_disk.foo",
						"category",
						"cloud"),
					resource.TestCheckResourceAttr(
						"alicloud_disk.foo",
						"size",
						"10"),
				),
			},
		},
	})

}

func TestAccAlicloudDisk_withTags(t *testing.T) {
	var v ecs.DiskItemType

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},

		//module name
		IDRefreshName: "alicloud_disk.foo",

		Providers:    testAccProviders,
		CheckDestroy: testAccCheckDiskDestroy,
		Steps: []resource.TestStep{
			resource.TestStep{
				Config: testAccDiskConfigWithTags,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckDiskExists("alicloud_disk.foo", &v),
				),
			},
		},
	})
}

func testAccCheckDiskExists(n string, disk *ecs.DiskItemType) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("Not found: %s", n)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("No Disk ID is set")
		}

		client := testAccProvider.Meta().(*AliyunClient)
		conn := client.ecsconn

		request := &ecs.DescribeDisksArgs{
			RegionId: client.Region,
			DiskIds:  []string{rs.Primary.ID},
		}

		response, _, err := conn.DescribeDisks(request)
		log.Printf("[WARN] disk ids %#v", rs.Primary.ID)

		if err == nil {
			if response != nil && len(response) > 0 {
				*disk = response[0]
				return nil
			}
		}
		return fmt.Errorf("Error finding ECS Disk %#v", rs.Primary.ID)
	}
}

func testAccCheckDiskDestroy(s *terraform.State) error {

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "alicloud_disk" {
			continue
		}

		// Try to find the Disk
		client := testAccProvider.Meta().(*AliyunClient)
		conn := client.ecsconn

		request := &ecs.DescribeDisksArgs{
			RegionId: client.Region,
			DiskIds:  []string{rs.Primary.ID},
		}

		response, _, err := conn.DescribeDisks(request)

		if response != nil && len(response) > 0 {
			return fmt.Errorf("Error ECS Disk still exist")
		}

		if err != nil {
			// Verify the error is what we want
			return err
		}
	}

	return nil
}

const testAccDiskConfig = `
resource "alicloud_disk" "foo" {
	# cn-beijing
	availability_zone = "cn-beijing-b"
        size = "10"
}
`
const testAccDiskConfigWithTags = `
resource "alicloud_disk" "foo" {
	# cn-beijing
	category = "cloud_efficiency"
	availability_zone = "cn-beijing-b"
        size = "30"
        tags {
        	Name = "TerraformTest"
        }
}
`
