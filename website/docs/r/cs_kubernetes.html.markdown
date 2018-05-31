---
layout: "alicloud"
page_title: "Alicloud: alicloud_cs_kubernetes"
sidebar_current: "docs-alicloud-resource-cs-kubernetes"
description: |-
  Provides a Alicloud resource to manage container kubernetes cluster.
---

# alicloud\_cs\_kubernetes

This resource will help you to manager a Kubernetes Cluster. The cluster is same as container service created by web console.

-> **NOTE:** Kubernetes cluster only supports VPC network and it can access internet while creating kubernetes cluster.
A Nat Gateway and configuring a SNAT for it can ensure one VPC network access internet. If there is no nat gateway in the
VPC, you can set `new_nat_gateway` to "true" to create one automatically.

-> **NOTE:** If there is no specified `vswitch_id`, the resource will create a new VPC and VSwitch while creating kubernetes cluster.

-> **NOTE:** Each kubernetes cluster contains 3 master nodes and those number cannot be changed at now.

-> **NOTE:** Creating kubernetes cluster need to install several packages and it will cost more than one hour. Please be patient.

## Example Usage

Basic Usage

```
data "alicloud_zones" "default" {
  "available_resource_creation"= "VSwitch"
}

resource "alicloud_cs_kubernetes" "main" {
  name_prefix = "my-first-k8s"
  availability_zone = "${data.alicloud_zones.default.zones.0.id}"
  new_nat_gateway = true
  master_instance_type = "ecs.n4.small"
  worker_instance_type = "ecs.n4.small"
  worker_number = 3
  password = "Test12345"
  pod_cidr = "192.168.1.0/24"
  service_cidr = "192.168.2.0/24"
  enable_ssh = true
  install_cloud_monitor = true
}
```
## Argument Reference

The following arguments are supported:

* `name` - The kubernetes cluster's name. It is the only in one Alicloud account.
* `name_prefix` - The kubernetes cluster name's prefix. It is conflict with `name`. If it is specified, terraform will using it to build the only cluster name. Default to "Terraform-Creation".
* `availability_zone` - (Force new resource) The Zone where new kubernetes cluster will be located. If it is not be specified, the value will be vswitch's zone.
* `vswitch_id` - (Force new resource) The vswitch where new kubernetes cluster will be located. If it is not specified, a new VPC and VSwicth will be built. It must be in the zone which `availability_zone` specified.
* `new_nat_gateway` - (Force new resource) Whether to create a new nat gateway while creating kubernetes cluster. Default to true.
* `master_instance_type` - (Required, Force new resource) The instance type of master node.
* `worker_instance_type` - (Required, Force new resource) The instance type of worker node.
* `worker_number` - The worker node number of the kubernetes cluster. Default to 3. It is limited up to 50 and if you want to enlarge it, please apply white list or contact with us.
* `password` - (Required, Force new resource) The password of ssh login cluster node.
* `pod_cidr` - (Required, Force new resource) The CIDR block for the pod network. It will be allocated automatically when `vswitch_id` is not specified.
It cannot be duplicated with the VPC CIDR and CIDR used by Kubernetes cluster in VPC, cannot be modified after creation.
Maximum number of hosts allowed in the cluster: 256. Refer to [Plan Kubernetes CIDR blocks under VPC](https://www.alibabacloud.com/help/doc-detail/64530.htm).
* `service_cidr` - (Required, Force new resource) The CIDR block for the service network.  It will be allocated automatically when `vswitch_id` is not specified.
It cannot be duplicated with the VPC CIDR and CIDR used by Kubernetes cluster in VPC, cannot be modified after creation.
* `enable_ssh` - (Force new resource) Whether to allow to SSH login kubernetes. Default to false.
* `master_disk_category` - (Force new resource) The system disk category of master node. Its valid value are `cloud_ssd` and `cloud_efficiency`. Default to `cloud_efficiency`.
* `master_disk_size` - (Force new resource) The system disk size of master node. Its valid value range [20~32768] in GB. Default to 20.
* `worker_disk_category` - (Force new resource) The system disk category of worker node. Its valid value are `cloud_ssd` and `cloud_efficiency`. Default to `cloud_efficiency`.
* `worker_disk_size` - (Force new resource) The system disk size of worker node. Its valid value range [20~32768] in GB. Default to 20.
* `install_cloud_monitor` - (Force new resource) Whether to install cloud monitor for the kubernetes' node.
* `is_outdated` - (Optional) Whether to use outdated instance type. Default to false.

## Attributes Reference

The following attributes are exported:

* `id` - The ID of the container cluster.
* `name` - The name of the container cluster.
* `availability_zone` - The ID of availability zone.
* `worker_number` The ECS instance node number in the current container cluster.
* `vswitch_id` - The ID of VSwitch where the current cluster is located.
* `vpc_id` - The ID of VPC where the current cluster is located.
* `slb_id` - (Deprecated from version 1.9.2).
* `slb_internet` - The ID of public load balancer where the current cluster master node is located.
* `slb_intranet` - The ID of private load balancer where the current cluster master node is located.
* `security_group_id` - The ID of security group where the current cluster worker node is located.
* `image_id` - The ID of node image.
* `nat_gateway_id` - The ID of nat gateway used to launch kubernetes cluster.
* `master_instance_type` - The instance type of master node.
* `worker_instance_type` - The instance type of worker node.
* `master_disk_category` - The system disk category of master node.
* `master_disk_size` - The system disk size of master node.
* `worker_disk_category` - The system disk category of worker node.
* `worker_disk_size` - The system disk size of worker node.
* `nodes` - List of cluster nodes. It contains several attributes to `Block Nodes`.
* `connections` - Map of kubernetes cluster connection information. It contains several attributes to `Block Connections`.

### Block Nodes

* `id` - ID of the node.
* `name` - Node name.
* `private_ip` - The private IP address of node.
* `role` - Node role. "Master" or "Worker"

### Block Connections

* `api_server_internet` - API Server Internet endpoint.
* `api_server_intranet` - API Server Intranet endpoint.
* `master_public_ip` - Master node SSH IP address.
* `service_domain` - Service Access Domain.

## Import

Kubernetes cluster can be imported using the id, e.g.

```
$ terraform import alicloud_cs_kubernetes.main ce4273f9156874b46bb
```