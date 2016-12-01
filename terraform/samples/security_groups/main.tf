variable "short_name" {
}
variable "vpc_id" {
}

resource "alicloud_security_group" "default" {
  name = "${var.short_name}-default"
  description = "Default security group for VPC"
  vpc_id = "${var.vpc_id}"
}

resource "alicloud_security_group" "control" {
  name = "${var.short_name}-control"
  description = "Allow inboud traffic for control nodes"
  vpc_id = "${var.vpc_id}"
}

resource "alicloud_security_group" "edge" {
  name = "${var.short_name}-edge"
  description = "Allow inboud traffic for edge routing"
  vpc_id = "${var.vpc_id}"
}

resource "alicloud_security_group" "worker" {
  name = "${var.short_name}-worker"
  description = "Allow inboud traffic for worker nodes"
  vpc_id = "${var.vpc_id}"
}

output "default_security_group" {
  value = "${alicloud_security_group.default.id}"
}

output "edge_security_group" {
  value = "${alicloud_security_group.edge.id}"
}

output "control_security_group" {
  value = "${alicloud_security_group.control.id}"
}

output "worker_security_group" {
  value = "${alicloud_security_group.worker.id}"
}
