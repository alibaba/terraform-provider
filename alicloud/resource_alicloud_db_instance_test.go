package alicloud

import (
	"fmt"
	"github.com/denverdino/aliyungo/common"
	"github.com/denverdino/aliyungo/rds"
	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/terraform"
	"log"
	"testing"
)

func TestAccAlicloudDBInstance_basic(t *testing.T) {
	var instance rds.DBInstanceAttribute

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},

		// module name
		IDRefreshName: "alicloud_db_instance.foo",

		Providers:    testAccProviders,
		CheckDestroy: testAccCheckDBInstanceDestroy,
		Steps: []resource.TestStep{
			resource.TestStep{
				Config: testAccDBInstanceConfig,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckDBInstanceExists(
						"alicloud_db_instance.foo", &instance),
					resource.TestCheckResourceAttr(
						"alicloud_db_instance.foo",
						"security_ip_list",
						"127.0.0.1"),
					resource.TestCheckResourceAttr(
						"alicloud_db_instance.foo",
						"Port",
						"3306"),
					resource.TestCheckResourceAttr(
						"alicloud_db_instance.foo",
						"DBInstanceStorage",
						"10"),
					resource.TestCheckResourceAttr(
						"alicloud_db_instance.foo",
						"InstanceNetworkType",
						"Classic"),
					resource.TestCheckResourceAttr(
						"alicloud_db_instance.foo",
						"DBInstanceNetType",
						"Intranet"),
					resource.TestCheckResourceAttr(
						"alicloud_db_instance.foo",
						"EngineVersion",
						"5.6"),
					resource.TestCheckResourceAttr(
						"alicloud_db_instance.foo",
						"Engine",
						"MySQL"),
				),
			},
		},
	})

}

func testAccCheckDBInstanceExists(n string, d *rds.DBInstanceAttribute) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("Not found: %s", n)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("No DB Instance ID is set")
		}

		client := testAccProvider.Meta().(*AliyunClient)
		attr, err := client.DescribeDBInstanceById(rs.Primary.ID)
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

func testAccCheckDBInstanceDestroy(s *terraform.State) error {
	client := testAccProvider.Meta().(*AliyunClient)

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "alicloud_db_instance.foo" {
			continue
		}

		ins, err := client.DescribeDBInstanceById(rs.Primary.ID)

		if ins != nil {
			return fmt.Errorf("Error DB Instance still exist")
		}

		// Verify the error is what we want
		if err != nil {
			// Verify the error is what we want
			e, _ := err.(*common.Error)
			if e.ErrorResponse.Code == InstanceNotfound {
				continue
			}
			return err
		}
	}

	return nil
}

const testAccDBInstanceConfig = `
resource "alicloud_db_instance" "foo" {
	commodity_code = "bards"
	engine = "MySQL"
	engine_version = "5.6"
	db_instance_class = "rds.mysql.t1.small"
	db_instance_storage = "10"
	instance_charge_type = "Postpaid"
	db_instance_net_type = "Intranet"
}
`
