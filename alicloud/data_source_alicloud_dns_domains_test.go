package alicloud

import (
	"github.com/hashicorp/terraform/helper/resource"
	"testing"
)

func TestAccAlicloudDnsDomainsDataSource_ali_domain(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckAlicloudDomainsDataSourceAliDomainConfig,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAlicloudDataSourceID("data.alicloud_dns_domains.domain"),
					resource.TestCheckResourceAttr("data.alicloud_dns_domains.domain", "domains.#", "1"),
					resource.TestCheckResourceAttr("data.alicloud_dns_domains.domain", "domains.0.domain_id", "6f1a920c-c4a0-4231-98ea-7c4e9a89218a"),
					resource.TestCheckResourceAttr("data.alicloud_dns_domains.domain", "domains.0.domain_name", "heguimin.top"),
					resource.TestCheckResourceAttr("data.alicloud_dns_domains.domain", "domains.0.version_code", "mianfei"),
					resource.TestCheckResourceAttr("data.alicloud_dns_domains.domain", "domains.0.group_name", "newfish"),
					resource.TestCheckResourceAttr("data.alicloud_dns_domains.domain", "domains.0.group_id", "85ab8713-4a30-4de4-9d20-155ff830f651"),
					resource.TestCheckResourceAttr("data.alicloud_dns_domains.domain", "domains.0.puny_code", "heguimin.top"),
				),
			},
		},
	})
}

func TestAccAlicloudDnsDomainsDataSource_version_code(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckAlicloudDomainsDataSourceVersionCodeConfig,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAlicloudDataSourceID("data.alicloud_dns_domains.domain"),
					resource.TestCheckResourceAttr("data.alicloud_dns_domains.domain", "domains.#", "2"),
				),
			},
		},
	})
}

func TestAccAlicloudDnsDomainsDataSource_name_regex(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckAlicloudDomainsDataSourceNameRegexConfig,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAlicloudDataSourceID("data.alicloud_dns_domains.domain"),
					resource.TestCheckResourceAttr("data.alicloud_dns_domains.domain", "domains.#", "1"),
				),
			},
		},
	})
}

func TestAccAlicloudDnsDomainsDataSource_group_name_regex(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckAlicloudDomainsDataSourceGroupNameRegexConfig,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAlicloudDataSourceID("data.alicloud_dns_domains.domain"),
					resource.TestCheckResourceAttr("data.alicloud_dns_domains.domain", "domains.#", "1"),
				),
			},
		},
	})
}

const testAccCheckAlicloudDomainsDataSourceAliDomainConfig = `
data "alicloud_dns_domains" "domain" {
  ali_domain = true
}`

const testAccCheckAlicloudDomainsDataSourceVersionCodeConfig = `
data "alicloud_dns_domains" "domain" {
  version_code = "mianfei"
}`

const testAccCheckAlicloudDomainsDataSourceNameRegexConfig = `
data "alicloud_dns_domains" "domain" {
  domain_name_regex = "^hegui"
}`

const testAccCheckAlicloudDomainsDataSourceGroupNameRegexConfig = `
data "alicloud_dns_domains" "domain" {
  group_name_regex = ".*"
}`
