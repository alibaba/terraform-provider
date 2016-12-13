variable "slb_name" {
  default = "slb_worker"
}

variable "instances" {
  type = "list"
  default = [
    "i-2zecejialx1rx513qcyv",
    "i-2zedgb871dbnpc5x3w9n"]
}

variable "internet_charge_type" {
  default = "paybytraffic"
}

variable "internet" {
  default = true
}