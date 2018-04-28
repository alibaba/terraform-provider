package alicloud

import (
	"fmt"
	"testing"

	"github.com/aliyun/alibaba-cloud-sdk-go/services/drds"
	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/terraform"
)

func TestAccAlicloudDRDSInstance_Basic(t *testing.T) {
	var instance drds.DescribeDrdsInstanceResponse
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},

		IDRefreshName: "alicloud_drds_instance.basic",
		Providers:     testAccProviders,
		CheckDestroy:  testAccCheckDRDSInstanceDestroy,
		Steps: []resource.TestStep{
			resource.TestStep{
				Config: testAccDrdsInstance,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckDRDSInstanceExist(
						"alicloud_drds_instance.basic", &instance),
				),
			},
		},
	})
}

func testAccCheckDRDSInstanceExist(n string, instance *drds.DescribeDrdsInstanceResponse) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("not found: %s", n)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("no DRDS Instance ID is set")
		}

		client := testAccProvider.Meta().(*AliyunClient).drdsconn
		req := drds.CreateDescribeDrdsInstanceRequest()
		req.DrdsInstanceId = rs.Primary.ID

		response, err := client.DescribeDrdsInstance(req)

		if err == nil && response != nil && response.Data.DrdsInstanceId != "" {
			instance = response
			return nil
		}
		return fmt.Errorf("error finding DRDS instance %#v", rs.Primary.ID)
	}
}

func testAccCheckDRDSInstanceDestroy(s *terraform.State) error {
	for _, rs := range s.RootModule().Resources {
		if rs.Type != "alicloud_drds_instance" {
			continue
		}

		client := testAccProvider.Meta().(*AliyunClient)
		conn := client.drdsconn
		req := drds.CreateDescribeDrdsInstanceRequest()
		req.DrdsInstanceId = rs.Primary.ID
		response, err := conn.DescribeDrdsInstance(req)

		if err == nil && response != nil && response.Data.Status != "5" {
			return fmt.Errorf("error! DRDS instance still exists")
		}
	}

	return nil
}

const testAccDrdsInstance = `
provider "alicloud" {
	region = "cn-hangzhou"
}
resource "alicloud_drds_instance" "basic" {
  provider = "alicloud"
  description = "for rds"
  type = "PRIVATE"
  zone_id = "cn-hangzhou-f"
  specification = "drds.sn1.4c8g.8C16G"
  pay_type = "drdsPost"
  instance_series = "drds.sn1.4c8g"
  quantity = 1
}
`
