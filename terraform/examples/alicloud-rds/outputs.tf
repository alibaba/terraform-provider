
output "ports" {
  value = "${join(",", alicloud_db_instance.dc.*.port)}"
}

output "connections" {
  value = "${jsonencode(alicloud_db_instance.dc.connections)}"
}