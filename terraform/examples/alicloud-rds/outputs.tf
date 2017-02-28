
output "port" {
  value = "${alicloud_db_instance.dc.port}"
}

output "connections" {
  value = "${alicloud_db_instance.dc.connections}"
}