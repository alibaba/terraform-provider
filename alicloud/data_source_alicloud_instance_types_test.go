package alicloud

import (
	"fmt"
	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/terraform"
	"testing"
)

func TestAccAlicloudInstanceTypesDataSource_basic(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckAlicloudInstanceTypesDataSourceBasicConfig,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAlicloudInstanceTypesDataSourceID("data.alicloud_instance_types.4c8g"),

					resource.TestCheckResourceAttr("data.alicloud_instance_types.4c8g", "instance_types.#", "4"),

					resource.TestCheckResourceAttr("data.alicloud_instance_types.4c8g", "instance_types.0.cpu_core_count", "4"),
					resource.TestCheckResourceAttr("data.alicloud_instance_types.4c8g", "instance_types.0.memory_size", "8"),
				),
			},

			resource.TestStep{
				Config: testAccCheckAlicloudInstanceTypesDataSourceBasicConfigUpdate,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAlicloudInstanceTypesDataSourceID("data.alicloud_instance_types.4c8g"),

					resource.TestCheckResourceAttr("data.alicloud_instance_types.4c8g", "instance_types.#", "1"),

					resource.TestCheckResourceAttr("data.alicloud_instance_types.4c8g", "instance_types.0.cpu_core_count", "4"),
					resource.TestCheckResourceAttr("data.alicloud_instance_types.4c8g", "instance_types.0.memory_size", "8"),
				),
			},
		},
	})
}

func testAccCheckAlicloudInstanceTypesDataSourceID(n string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("Can't find instance type data source: %s", n)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("instance type data source ID not set")
		}
		return nil
	}
}

const testAccCheckAlicloudInstanceTypesDataSourceBasicConfig = `
data "alicloud_instance_types" "4c8g" {
	cpu_core_count = 4
	memory_size = 8
}
`

const testAccCheckAlicloudInstanceTypesDataSourceBasicConfigUpdate = `
data "alicloud_instance_types" "4c8g" {
	instance_type_family= "ecs.s3"
	cpu_core_count = 4
	memory_size = 8
}
`
