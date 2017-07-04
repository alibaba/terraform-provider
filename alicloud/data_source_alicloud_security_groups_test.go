package alicloud

import (
	"github.com/hashicorp/terraform/helper/resource"
	"testing"
)

func TestAccAlicloudSecurityGroupsDataSource_basic(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckAlicloudSecurityGroupsDataSourceBasicConfig,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAlicloudDataSourceID("data.alicloud_security_groups.foo1"),
				),
			},
		},
	})
}

func TestAccAlicloudSecurityGroupsDataSource_name_regex(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckAlicloudSecurityGroupsDataSourceNameRegexConfig,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAlicloudDataSourceID("data.alicloud_security_groups.foo2"),
				),
			},
		},
	})
}

const testAccCheckAlicloudSecurityGroupsDataSourceBasicConfig = `
data "alicloud_security_groups" "foo1"{
}
`
const testAccCheckAlicloudSecurityGroupsDataSourceNameRegexConfig = `
data "alicloud_security_groups" "foo2"{
"security_group_name_regex"="^tf"
}
`
