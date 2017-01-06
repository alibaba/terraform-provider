output "route_table_id" {
  value = "${alicloud_security_group_rule.ssh.route_table_id}"
}

output "nexthop_type" {
  value = "${alicloud_security_group_rule.ssh.nexthop_type}"
}

output "nexthop_id" {
  value = "${alicloud_security_group_rule.ssh.nexthop_id}"
}