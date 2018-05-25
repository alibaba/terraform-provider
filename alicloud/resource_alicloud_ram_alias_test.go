package alicloud

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/terraform"
)

func TestAccAlicloudRamAlias_basic(t *testing.T) {
	var v string

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},

		// module name
		IDRefreshName: "alicloud_ram_alias.alias",

		Providers:    testAccProviders,
		CheckDestroy: testAccCheckRamAliasDestroy,
		Steps: []resource.TestStep{
			resource.TestStep{
				Config: testAccRamAliasConfig,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckRamAliasExists(
						"alicloud_ram_alias.alias", &v),
					resource.TestCheckResourceAttr(
						"alicloud_ram_alias.alias",
						"account_alias",
						"hallo"),
				),
			},
		},
	})

}

func testAccCheckRamAliasExists(n string, alias *string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("Not found: %s", n)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("No Alias ID is set")
		}

		client := testAccProvider.Meta().(*AliyunClient)
		conn := client.ramconn

		response, err := conn.GetAccountAlias()

		if err == nil {
			*alias = response.AccountAlias
			return nil
		}
		return fmt.Errorf("Error finding alias %s.", rs.Primary.ID)
	}
}

func testAccCheckRamAliasDestroy(s *terraform.State) error {

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "alicloud_ram_alias" {
			continue
		}

		// Try to find the alias
		client := testAccProvider.Meta().(*AliyunClient)
		conn := client.ramconn

		_, err := conn.GetAccountAlias()

		if err != nil {
			if RamEntityNotExist(err) {
				return nil
			}
			return err
		}
	}
	return nil
}

const testAccRamAliasConfig = `
resource "alicloud_ram_alias" "alias" {
  account_alias = "hallo"
}`
