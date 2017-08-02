package alicloud

import (
	"github.com/hashicorp/terraform/helper/resource"
	"testing"
)

func TestAccAlicloudDnsDomainRecordsDataSource_host_record_regex(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckAlicloudDomainRecordsDataSourceHostRecordRegexConfig,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAlicloudDataSourceID("data.alicloud_dns_domain_records.record"),
					resource.TestCheckResourceAttr("data.alicloud_dns_domain_records.record", "records.#", "1"),
					resource.TestCheckResourceAttr("data.alicloud_dns_domain_records.record", "records.0.record_id", "3438492787133440"),
					resource.TestCheckResourceAttr("data.alicloud_dns_domain_records.record", "records.0.domain_name", "heguimin.top"),
					resource.TestCheckResourceAttr("data.alicloud_dns_domain_records.record", "records.0.host_record", "smtp"),
					resource.TestCheckResourceAttr("data.alicloud_dns_domain_records.record", "records.0.status", "ENABLE"),
					resource.TestCheckResourceAttr("data.alicloud_dns_domain_records.record", "records.0.ttl", "600"),
					resource.TestCheckResourceAttr("data.alicloud_dns_domain_records.record", "records.0.type", "CNAME"),
					resource.TestCheckResourceAttr("data.alicloud_dns_domain_records.record", "records.0.value", "smtp.mxhichina.com"),
					resource.TestCheckResourceAttr("data.alicloud_dns_domain_records.record", "records.0.line", "default"),
				),
			},
		},
	})
}

func TestAccAlicloudDnsDomainRecordsDataSource_type(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckAlicloudDomainRecordsDataSourceTypeConfig,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAlicloudDataSourceID("data.alicloud_dns_domain_records.record"),
					resource.TestCheckResourceAttr("data.alicloud_dns_domain_records.record", "records.#", "7"),
				),
			},
		},
	})
}

func TestAccAlicloudDnsDomainRecordsDataSource_value_regex(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckAlicloudDomainRecordsDataSourceValueRegexConfig,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAlicloudDataSourceID("data.alicloud_dns_domain_records.record"),
					resource.TestCheckResourceAttr("data.alicloud_dns_domain_records.record", "records.#", "3"),
				),
			},
		},
	})
}

func TestAccAlicloudDnsDomainRecordsDataSource_line(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckAlicloudDomainRecordsDataSourceLineConfig,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAlicloudDataSourceID("data.alicloud_dns_domain_records.record"),
					resource.TestCheckResourceAttr("data.alicloud_dns_domain_records.record", "records.#", "17"),
				),
			},
		},
	})
}

func TestAccAlicloudDnsDomainRecordsDataSource_status(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckAlicloudDomainRecordsDataSourceStatusConfig,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAlicloudDataSourceID("data.alicloud_dns_domain_records.record"),
					resource.TestCheckResourceAttr("data.alicloud_dns_domain_records.record", "records.#", "17"),
				),
			},
		},
	})
}

func TestAccAlicloudDnsDomainRecordsDataSource_is_locked(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckAlicloudDomainRecordsDataSourceIsLockedConfig,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAlicloudDataSourceID("data.alicloud_dns_domain_records.record"),
					resource.TestCheckResourceAttr("data.alicloud_dns_domain_records.record", "records.#", "17"),
				),
			},
		},
	})
}

const testAccCheckAlicloudDomainRecordsDataSourceHostRecordRegexConfig = `
data "alicloud_dns_domain_records" "record" {
  domain_name = "heguimin.top"
  host_record_regex = ".*smtp.*"
}`

const testAccCheckAlicloudDomainRecordsDataSourceTypeConfig = `
data "alicloud_dns_domain_records" "record" {
  domain_name = "heguimin.top"
  type = "CNAME"
}`

const testAccCheckAlicloudDomainRecordsDataSourceValueRegexConfig = `
data "alicloud_dns_domain_records" "record" {
  domain_name = "heguimin.top"
  value_regex = "^mail.mxhichina"
}`

const testAccCheckAlicloudDomainRecordsDataSourceStatusConfig = `
data "alicloud_dns_domain_records" "record" {
  domain_name = "heguimin.top"
  status = "enable"
}`

const testAccCheckAlicloudDomainRecordsDataSourceIsLockedConfig = `
data "alicloud_dns_domain_records" "record" {
  domain_name = "heguimin.top"
  is_locked = false
}`

const testAccCheckAlicloudDomainRecordsDataSourceLineConfig = `
data "alicloud_dns_domain_records" "record" {
  domain_name = "heguimin.top"
  line = "default"
}`
