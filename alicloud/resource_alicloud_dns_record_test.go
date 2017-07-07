package alicloud

import (
	"fmt"
	"testing"

	"github.com/denverdino/aliyungo/dns"
	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/terraform"
	"log"
)

func TestAccAlicloudDnsRecord_basic(t *testing.T) {
	var v dns.RecordType

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},

		// module name
		IDRefreshName: "alicloud_dns_record.record",

		Providers:    testAccProviders,
		CheckDestroy: testAccCheckDnsRecordDestroy,
		Steps: []resource.TestStep{
			resource.TestStep{
				Config: testAccDnsRecordConfig,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckDnsRecordExists(
						"alicloud_dns_record.record", &v),
					resource.TestCheckResourceAttr(
						"alicloud_dns_record.record",
						"name",
						"heguimin.top"),
					resource.TestCheckResourceAttr(
						"alicloud_dns_record.record",
						"type",
						"CNAME"),
				),
			},
		},
	})

}

func testAccCheckDnsRecordExists(n string, record *dns.RecordType) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("Not found: %s", n)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("No Domain Record ID is set")
		}

		client := testAccProvider.Meta().(*AliyunClient)
		conn := client.dnsconn

		request := &dns.DescribeDomainRecordInfoArgs{
			RecordId: rs.Primary.ID,
		}

		response, err := conn.DescribeDomainRecordInfo(request)
		log.Printf("[WARN] Domain record id %#v", rs.Primary.ID)

		if err == nil {
			*record = response.RecordType
			return nil
		}
		return fmt.Errorf("Error finding domain record %#v", rs.Primary.ID)
	}
}

func testAccCheckDnsRecordDestroy(s *terraform.State) error {

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "alicloud_dns_record" {
			continue
		}

		// Try to find the domain record
		client := testAccProvider.Meta().(*AliyunClient)
		conn := client.dnsconn

		request := &dns.DescribeDomainRecordInfoArgs{
			RecordId: rs.Primary.ID,
		}

		response, err := conn.DescribeDomainRecordInfo(request)

		if response.RecordId != "" || err != nil {
			return fmt.Errorf("Error Domain record still exist.")
		}
	}

	return nil
}

const testAccDnsRecordConfig = `
resource "alicloud_dns_record" "record" {
  name = "heguimin.top"
  host_record = "alimailskajdh"
  type = "CNAME"
  value = "mail.mxhichin.com"
  count = 1
}
`
