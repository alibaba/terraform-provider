package alicloud

import (
	"testing"

	"github.com/hashicorp/terraform/helper/resource"
)

func TestAccAlicloudDnsRecordsDataSource_host_record_regex(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckAlicloudDnsRecordsDataSourceHostRecordRegexConfig,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAlicloudDataSourceID("data.alicloud_dns_records.record"),
					resource.TestCheckResourceAttr("data.alicloud_dns_records.record", "records.#", "1"),
					resource.TestCheckResourceAttr("data.alicloud_dns_records.record", "records.0.host_record", "smtp"),
					resource.TestCheckResourceAttr("data.alicloud_dns_records.record", "records.0.type", "CNAME"),
				),
			},
		},
	})
}

func TestAccAlicloudDnsRecordsDataSource_type(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckAlicloudDnsRecordsDataSourceTypeConfig,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAlicloudDataSourceID("data.alicloud_dns_records.record"),
					resource.TestCheckResourceAttr("data.alicloud_dns_records.record", "records.0.type", "CNAME"),
				),
			},
		},
	})
}

func TestAccAlicloudDnsRecordsDataSource_value_regex(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckAlicloudDnsRecordsDataSourceValueRegexConfig,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAlicloudDataSourceID("data.alicloud_dns_records.record"),
					resource.TestCheckResourceAttr("data.alicloud_dns_records.record", "records.0.value", "mail.mxhichina.com"),
				),
			},
		},
	})
}

func TestAccAlicloudDnsRecordsDataSource_line(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckAlicloudDnsRecordsDataSourceLineConfig,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAlicloudDataSourceID("data.alicloud_dns_records.record"),
					resource.TestCheckResourceAttr("data.alicloud_dns_records.record", "records.0.line", "default"),
				),
			},
		},
	})
}

func TestAccAlicloudDnsRecordsDataSource_status(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckAlicloudDnsRecordsDataSourceStatusConfig,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAlicloudDataSourceID("data.alicloud_dns_records.record"),
					resource.TestCheckResourceAttr("data.alicloud_dns_records.record", "records.0.status", "enable"),
				),
			},
		},
	})
}

func TestAccAlicloudDnsRecordsDataSource_is_locked(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckAlicloudDnsRecordsDataSourceIsLockedConfig,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAlicloudDataSourceID("data.alicloud_dns_records.record"),
					resource.TestCheckResourceAttr("data.alicloud_dns_records.record", "records.0.locked", "false"),
				),
			},
		},
	})
}

const testAccCheckAlicloudDnsRecordsDataSourceHostRecordRegexConfig = `
data "alicloud_dns_domains" "domains" {}

data "alicloud_dns_records" "record" {
  domain_name = "${data.alicloud_dns_domains.domains.domains.0.domain_name}"
  host_record_regex = "^smtp"
}`

const testAccCheckAlicloudDnsRecordsDataSourceTypeConfig = `
data "alicloud_dns_domains" "domains" {}

data "alicloud_dns_records" "record" {
  domain_name = "${data.alicloud_dns_domains.domains.domains.0.domain_name}"
  type = "CNAME"
}`

const testAccCheckAlicloudDnsRecordsDataSourceValueRegexConfig = `
data "alicloud_dns_domains" "domains" {}

data "alicloud_dns_records" "record" {
  domain_name = "${data.alicloud_dns_domains.domains.domains.0.domain_name}"
  value_regex = "^mail"
}`

const testAccCheckAlicloudDnsRecordsDataSourceStatusConfig = `
data "alicloud_dns_domains" "domains" {}

data "alicloud_dns_records" "record" {
  domain_name = "${data.alicloud_dns_domains.domains.domains.0.domain_name}"
  status = "enable"
}`

const testAccCheckAlicloudDnsRecordsDataSourceIsLockedConfig = `
data "alicloud_dns_domains" "domains" {}

data "alicloud_dns_records" "record" {
  domain_name = "${data.alicloud_dns_domains.domains.domains.0.domain_name}"
  is_locked = false
}`

const testAccCheckAlicloudDnsRecordsDataSourceLineConfig = `
data "alicloud_dns_domains" "domains" {}

data "alicloud_dns_records" "record" {
  domain_name = "${data.alicloud_dns_domains.domains.domains.0.domain_name}"
  line = "default"
}`
