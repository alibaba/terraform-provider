variable "name" {
  description = "Solution Name"
  default = "APIM_VPC"
}

variable "access_key" {
default=""
}

variable "secret_key" {
default=""
}

variable "region" {
default=""
}

variable "cidr" {
  description = "CIDR range to use for the VPC"
  default     = "192.168.0.0/16"
}


variable "az_count" {
  description = "Number of availability zones to use"
  default = 3
}

variable "nat_name"{
default = "Nat_terraform"
}
