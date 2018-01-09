package alicloud

import (
	"fmt"
	"log"
	"testing"

	"github.com/denverdino/aliyungo/ram"
	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/terraform"
)

func TestAccAlicloudRamGroup_basic(t *testing.T) {
	var v ram.Group

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},

		// module name
		IDRefreshName: "alicloud_ram_group.group",

		Providers:    testAccProviders,
		CheckDestroy: testAccCheckRamGroupDestroy,
		Steps: []resource.TestStep{
			resource.TestStep{
				Config: testAccRamGroupConfig,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckRamGroupExists(
						"alicloud_ram_group.group", &v),
					resource.TestCheckResourceAttr(
						"alicloud_ram_group.group",
						"name",
						"groupname"),
					resource.TestCheckResourceAttr(
						"alicloud_ram_group.group",
						"comments",
						"group comments"),
				),
			},
		},
	})

}

func testAccCheckRamGroupExists(n string, group *ram.Group) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("Not found: %s", n)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("No Group ID is set")
		}

		client := testAccProvider.Meta().(*AliyunClient)
		conn := client.ramconn

		request := ram.GroupQueryRequest{
			GroupName: rs.Primary.Attributes["name"],
		}

		response, err := conn.GetGroup(request)
		log.Printf("[WARN] Group id %#v", rs.Primary.ID)

		if err == nil {
			*group = response.Group
			return nil
		}
		return fmt.Errorf("Error finding group %#v", rs.Primary.ID)
	}
}

func testAccCheckRamGroupDestroy(s *terraform.State) error {

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "alicloud_ram_group" {
			continue
		}

		// Try to find the group
		client := testAccProvider.Meta().(*AliyunClient)
		conn := client.ramconn

		request := ram.GroupQueryRequest{
			GroupName: rs.Primary.Attributes["name"],
		}

		_, err := conn.GetGroup(request)

		if err != nil {
			if RamEntityNotExist(err) {
				return nil
			}
			return err
		}
	}
	return nil
}

const testAccRamGroupConfig = `
resource "alicloud_ram_group" "group" {
  name = "groupname"
  comments = "group comments"
  force=true
}`
