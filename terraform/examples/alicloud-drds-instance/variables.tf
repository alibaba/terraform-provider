variable "description" {
  default = "-DRDS instance for RDS"
}

variable "type" {
  default = "PRIVATE"
}

variable "zone_id" {
  default = "cn-hangzhou-b"
}

variable "specification" {
  default = "drds.sn1.4c8g.16C32G"
}

variable "pay_type" {
  default = "drdsPost"
}

variable "instance_series" {
  default = "drds.sn1.4c8g"
}
