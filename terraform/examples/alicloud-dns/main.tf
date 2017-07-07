resource "alicloud_dns_group" "group" {
  name = "${var.group_name}"
  count = "${var.count}"
}


resource "alicloud_dns" "dns" {
  name = "${var.domain_name}"
  group_id = "${alicloud_dns_group.group.id}"
  count = "${var.count}"
}


resource "alicloud_dns_record" "record" {
  name = "${alicloud_dns.dns.name}"
  host_record = "alimailskajdh"
  type = "CNAME"
  value = "mail.mxhichind.com"
  count = "${var.count}"
}
