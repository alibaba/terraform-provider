package alicloud

import (
	"fmt"
	"log"
	"testing"

	"github.com/denverdino/aliyungo/common"
	"github.com/denverdino/aliyungo/ecs"
	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/helper/schema"
	"github.com/hashicorp/terraform/terraform"
)

func TestAccInstance_basic(t *testing.T) {
	var v ecs.InstanceAttributesType
	// todo: create and mount volume

	testCheck := func(*terraform.State) error {
		if v.ZoneId == "" {
			return fmt.Errorf("bad availability zone")
		}

		if len(v.SecurityGroupIds.SecurityGroupId) == 0 {
			return fmt.Errorf("no security group: %#v", v.SecurityGroupIds.SecurityGroupId)
		}

		return nil
	}

	resource.Test(t, resource.TestCase{
		PreCheck: func() { testAccPreCheck(t) },

		// module name
		IDRefreshName: "alicloud_instance.foo",

		Providers:    testAccProviders,
		CheckDestroy: testAccCheckInstanceDestroy,
		Steps: []resource.TestStep{
			resource.TestStep{
				Config: testAccInstanceConfig,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckInstanceExists(
						"alicloud_instance.foo", &v),
					testCheck,
					resource.TestCheckResourceAttr(
						"alicloud_instance.foo",
						"image_id",
						"m-2zeistfz0zw14pqx9c43"),
					resource.TestCheckResourceAttr(
						"alicloud_instance.foo",
						"instance_name",
						"test_foo"),
				),
			},

			// We repeat the exact same test so that we can be sure
			// that the user data hash stuff is working without generating
			// an incorrect diff.
			resource.TestStep{
				Config: testAccInstanceConfig,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckInstanceExists(
						"alicloud_instance.foo", &v),
					testCheck,
					resource.TestCheckResourceAttr(
						"alicloud_instance.foo",
						"image_id",
						"m-2zeistfz0zw14pqx9c43"),
					resource.TestCheckResourceAttr(
						"alicloud_instance.foo",
						"instance_name",
						"test_foo"),
				),
			},
		},
	})

}

func testAccCheckInstanceExists(n string, i *ecs.InstanceAttributesType) resource.TestCheckFunc {
	providers := []*schema.Provider{testAccProvider}
	return testAccCheckInstanceExistsWithProviders(n, i, &providers)
}

func testAccCheckInstanceExistsWithProviders(n string, i *ecs.InstanceAttributesType, providers *[]*schema.Provider) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		log.Printf("[DEBUG] check instance exists resource name is: %s", n)
		log.Printf("[DEBUG] check instance exists root module: %s", s.RootModule())
		log.Printf("[DEBUG] check instance exists all resources: %s", s.RootModule().Resources)
		rs, ok := s.RootModule().Resources[n]
		log.Printf("[DEBUG] check instance exists get resource: %s", rs)
		if !ok {
			return fmt.Errorf("Not found: %s", n)
		}
		log.Printf("[DEBUG] check instance exists: %s", rs.Primary.ID)
		if rs.Primary.ID == "" {
			return fmt.Errorf("No ID is set")
		}
		for _, provider := range *providers {
			// Ignore if Meta is empty, this can happen for validation providers
			if provider.Meta() == nil {
				continue
			}

			conn := provider.Meta().(*AliyunClient).ecsconn
			// todo: describeInstance or DescribeInstances?
			instance, err := conn.DescribeInstanceAttribute(rs.Primary.ID)

			if err == nil && instance != nil {
				*i = *instance
				return nil
			}

			// Verify the error is what we want
			e, _ := err.(*common.Error)
			if e.ErrorResponse.Code == "Instance.NotFound" {
				continue
			}
			if err != nil {
				return err
			}
		}

		return fmt.Errorf("Instance not found")
	}
}

func testAccCheckInstanceDestroy(s *terraform.State) error {
	return testAccCheckInstanceDestroyWithProvider(s, testAccProvider)
}

func testAccCheckInstanceDestroyWithProviders(providers *[]*schema.Provider) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		for _, provider := range *providers {
			if provider.Meta() == nil {
				continue
			}
			if err := testAccCheckInstanceDestroyWithProvider(s, provider); err != nil {
				return err
			}
		}
		return nil
	}
}

func testAccCheckInstanceDestroyWithProvider(s *terraform.State, provider *schema.Provider) error {
	conn := provider.Meta().(*AliyunClient).ecsconn

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "alicloud_instance" {
			continue
		}
		log.Printf("[DEBUG] check instance destroy: %v", rs.Primary.ID)

		// Try to find the resource
		// todo: describeInstance or DescribeInstances?
		instance, err := conn.DescribeInstanceAttribute(rs.Primary.ID)
		if err == nil {
			if instance.Status != "" && instance.Status != "Stopped" {
				return fmt.Errorf("Found unstopped instance: %s", instance.InstanceId)
			}
		}

		// Verify the error is what we want
		e, _ := err.(*common.Error)
		if e.ErrorResponse.Code == "InvalidInstanceId.NotFound" {
			continue
		}

		return err
	}

	return nil
}

const testAccInstanceConfig_1 = `
resource "alicloud_security_group" "tf_test_foo" {
	name = "tf_test_foo"
	description = "foo"
}

resource "alicloud_instance" "foo" {
	# cn-beijing
	availability_zone = "cn-beijing-b"
	image_id = "m-2zeistfz0zw14pqx9c43"

	instance_type = "ecs.s2.large"
	instance_network_type = "Classic"
	security_group_id = "${alicloud_security_group.tf_test_foo.id}"
	instance_name = "test_foo"
}
`

const testAccInstanceConfig = `
resource "alicloud_security_group" "tf_test_foo" {
	name = "tf_test_foo"
	description = "foo"
}

resource "alicloud_instance" "foo" {
	# cn-beijing
	availability_zone = "cn-beijing-b"
	image_id = "m-2zeistfz0zw14pqx9c43"

	instance_type = "ecs.s2.large"
	instance_network_type = "Classic"
	internet_charge_type = "PayByBandwidth"
	security_group_id = "${alicloud_security_group.tf_test_foo.id}"
	instance_name = "test_foo"
}
`
