provider "alicloud" {
	region = "cn-hangzhou"
}
//data "alicloud_zones" "foo" {
//	available_instance_type= "ecs.c2.xlarge"
//	available_resource_creation= "IoOptimized"
//	available_disk_category= "cloud"
//}
data "alicloud_images" "slave" {
	most_recent = true
	owners = "system"
	name_regex = "^ubuntu_14.*_64"
}