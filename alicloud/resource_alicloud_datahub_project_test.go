package alicloud

import (
	"fmt"
	"testing"

	// DEBUG only
	// "github.com/aliyun/aliyun-datahub-sdk-go/datahub/utils"

	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/terraform"
)

func TestAccAlicloudDatahubProject_Basic(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},

		// module name
		IDRefreshName: "alicloud_datahub_project.basic",
		Providers:     testAccProviders,
		CheckDestroy:  testAccCheckDatahubProjectDestroy,

		Steps: []resource.TestStep{
			resource.TestStep{
				Config: testAccDatahubProject,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckDatahubProjectExist(
						"alicloud_datahub_project.basic"),
					resource.TestCheckResourceAttr(
						"alicloud_datahub_project.basic",
						"project_name", "tf_test_datahub_project"),
				),
			},
		},
	})
}

func TestAccAlicloudDatahubProject_Update(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},

		// module name
		IDRefreshName: "alicloud_datahub_project.basic",
		Providers:     testAccProviders,
		CheckDestroy:  testAccCheckDatahubProjectDestroy,
		Steps: []resource.TestStep{
			resource.TestStep{
				Config: testAccDatahubProject,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckDatahubProjectExist(
						"alicloud_datahub_project.basic"),
					resource.TestCheckResourceAttr(
						"alicloud_datahub_project.basic",
						"comment", "project for basic."),
				),
			},

			resource.TestStep{
				Config: testAccDatahubProjectUpdate,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckDatahubProjectExist(
						"alicloud_datahub_project.basic"),
					resource.TestCheckResourceAttr(
						"alicloud_datahub_project.basic",
						"comment", "project for update."),
				),
			},
		},
	})
}

func testAccCheckDatahubProjectExist(n string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("not found Datahub project: %s", n)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("no Datahub project ID is set")
		}

		dh := testAccProvider.Meta().(*AliyunClient).dhconn
		_, err := dh.GetProject(rs.Primary.ID)

		// XXX DEBUG only
		// prj, err := dh.GetProject(rs.Primary.ID)
		// fmt.Printf("\nXXX:name:%s\n", rs.Primary.ID)
		// fmt.Printf("XXX:comment:%s\n", prj.Comment)
		// fmt.Printf("XXX:create_time:%s\n", utils.Uint64ToTimeString(prj.CreateTime))
		// fmt.Printf("XXX:last_modify_time:%s\n", utils.Uint64ToTimeString(prj.LastModifyTime))

		if err != nil {
			return err
		}
		return nil
	}
}

func testAccCheckDatahubProjectDestroy(s *terraform.State) error {
	for _, rs := range s.RootModule().Resources {
		if rs.Type != "alicloud_datahub_project" {
			continue
		}

		dh := testAccProvider.Meta().(*AliyunClient).dhconn
		_, err := dh.GetProject(rs.Primary.ID)

		if err != nil {
			continue
		}

		return fmt.Errorf("Datahub project %s still exists", rs.Primary.ID)
	}

	return nil
}

const testAccDatahubProject = `
provider "alicloud" {
    region = "cn-beijing"
}
variable "project_name" {
  default = "tf_test_datahub_project"
}
resource "alicloud_datahub_project" "basic" {
  project_name = "${var.project_name}"
  comment = "project for basic."
}
`

const testAccDatahubProjectUpdate = `
provider "alicloud" {
    region = "cn-beijing"
}
variable "project_name" {
  default = "tf_test_datahub_project"
}
resource "alicloud_datahub_project" "basic" {
  project_name = "${var.project_name}"
  comment = "project for update."
}
`
