variable "name" {
  default = "slb_alicloud"
}
variable "vpc_id" {
  default = "vpc-2zefstl5ld0ysotubvzm6"
}
variable "vswitch_id" {
  default = "vsw-2zerfc6gu3o7s9cbtcc3u"
}

variable "internet_charge_type" {
  default = "paybytraffic"
}

variable "internet" {
  default = "false"
}


resource "alicloud_slb" "instance" {
  name = "${var.name}"
  vpc_id = "${var.vpc_id}"
  vswitch_id = "${var.vswitch_id}"
  internet_charge_type = "${var.internet_charge_type}"
  internet = "${var.internet}"
  listener = [
    {
      "instance_port" = "3375"
      "instance_protocol" = "tcp"
      "lb_port" = "3376"
      "lb_protocol" = "tcp"
      "bandwidth" = "5"
    }]
}

output "slb_id" {
  value = "${alicloud_slb.instance.id}"
}

output "slbname" {
  value = "${alicloud_slb.instance.name}"
}