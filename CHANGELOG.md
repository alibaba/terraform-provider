## 1.0.7 (unreleased)

## 1.0.6 (May 2, 2017)

IMPROVEMENTS:

  * resource/alicloud_instance: add replaceing system disk function ([#119](https://github.com/alibaba/terraform-provider/pull/119))


## 1.0.5 (April 18, 2017)

IMPROVEMENTS:

  * resource/alicloud_instance: delete ecs instance with retry ([#113](https://github.com/alibaba/terraform-provider/pull/113))

BUG FIXES:

  * resource/resource_alicloud_security_group_rule: check ptr before use it.([#113](https://github.com/alibaba/terraform-provider/pull/113))
  * resource/alicloud_instance: fix ecs internet_max_bandwidth_out cannot set zero bug. cause if don't want allocate public ip, then must set internet_max_bandwidth_out is zero.([#105](https://github.com/alibaba/terraform-provider/pull/105))

FEATURES:

  * **New Resource:** `alicloud_ess_scalinggroup` ([#113](https://github.com/alibaba/terraform-provider/pull/113))
  * **New Resource:** `alicloud_ess_scalingconfiguration` ([#113](https://github.com/alibaba/terraform-provider/pull/113))
  * **New Resource:** `alicloud_ess_scalingrule` ([#113](https://github.com/alibaba/terraform-provider/pull/113))
  * **New Resource:** `alicloud_ess_schedule` ([#113](https://github.com/alibaba/terraform-provider/pull/113))
  * **New Resource:** `alicloud_snat_entry` ([#105](https://github.com/alibaba/terraform-provider/pull/105))
  * **New Resource:** `alicloud_forward_entry` ([#105](https://github.com/alibaba/terraform-provider/pull/105))
  * add snat entry and forward entry template sample in /terraform/examples/alicloud-vpc-snat.


## 1.0.4 (March 17, 2017)

BUG FIXES:
  
  * resource/alicloud_db_instance: fix rds update failed bug ([#102](https://github.com/alibaba/terraform-provider/pull/102))
  * resource/alicloud_instance: fix ecs instance system disk size not work bug ([#100](https://github.com/alibaba/terraform-provider/pull/100))
  
IMPROVEMENTS:

  * alicloud/config: add businessinfo to sdk client ([#96](https://github.com/alibaba/terraform-provider/pull/96))

## 1.0.3 (March 4, 2017)

FEATURES:

  * **New Resource:** `alicloud_db_instance` ([#85](https://github.com/alibaba/terraform-provider/pull/85))

IMPROVEMENTS:

  * resource/alicloud_slb: support slb listener persistence_timeout and health check ([#86](https://github.com/alibaba/terraform-provider/pull/86))

## 1.0.2 (February 24, 2017)

IMPROVEMENTS:

  * resource/alicloud_instance: change create ecs postpaid instance API form createInstance to runInstances, support BusinessInfo ([#80](https://github.com/alibaba/terraform-provider/pull/80))
  * resource/alicloud_instance: change ecs parameter zoneId from required to optional ([#74](https://github.com/alibaba/terraform-provider/pull/74))
  * resource/alicloud_instance: support userdata ([#71](https://github.com/alibaba/terraform-provider/pull/71))
  
BUG FIXES:
  
  * resource/alicloud_security_group_rule: fix security group egress rule delete failed ([#79](https://github.com/alibaba/terraform-provider/pull/79))
  * data resource/alicloud_images: data alicloud_images filter all images not only the first page ([#78](https://github.com/alibaba/terraform-provider/pull/78))

## 1.0.1 (January 17, 2017)

FEATURES:

  * **New Data Resource:** `alicloud_regions` ([#67](https://github.com/alibaba/terraform-provider/pull/67))
  * **New Data Resource:** `alicloud_images` ([#66](https://github.com/alibaba/terraform-provider/pull/66))
  * **New Data Resource:** `alicloud_instance_types` ([#64](https://github.com/alibaba/terraform-provider/pull/64))
  * **New Data Resource:** `alicloud_zones` ([#64](https://github.com/alibaba/terraform-provider/pull/64))
  * **New Resource:** `alicloud_route_entry` ([#58](https://github.com/alibaba/terraform-provider/pull/58))
  * **New Resource:** `alicloud_security_group_rule` ([#49](https://github.com/alibaba/terraform-provider/pull/49))
  * **New Resource:** `alicloud_slb_attachment` ([#31](https://github.com/alibaba/terraform-provider/pull/31))

IMPROVEMENTS:

  * resource/alicloud_instance: update instance tags ([#57](https://github.com/alibaba/terraform-provider/pull/57))
  * resource/alicloud_instance: create ecs instance with multi security groups ([#28](https://github.com/alibaba/terraform-provider/pull/28))
  * resource/alicloud_nat_gateway: support modify nat gateway spec ([#22](https://github.com/alibaba/terraform-provider/pull/22))
  * resource/alicloud_nat_gateway: support multi bandwidthPackage ([#22](https://github.com/alibaba/terraform-provider/pull/22))

BUG FIXES:

  * resource/alicloud_instance: bug fix io_optimized, remove default value, required is true ([#68](https://github.com/alibaba/terraform-provider/pull/68))
  * resource/alicloud_instance: bug fix cannot read internet_charge_type ([#55](https://github.com/alibaba/terraform-provider/pull/55))
  * resource/alicloud_instance: bug fix tags, io_optimized, private_ip output ([#47](https://github.com/alibaba/terraform-provider/pull/47))
  * resource/alicloud_slb: slb output backendsever ([#45](https://github.com/alibaba/terraform-provider/pull/45))
  * resource/alicloud_disk: fix some disk defects ([#42](https://github.com/alibaba/terraform-provider/pull/42))
  * resource/alicloud_slb_attachment: bug fix slb attachment ([#36](https://github.com/alibaba/terraform-provider/pull/36))
  * resource/alicloud_slb: fix slb internetchartype param for go sdk updated ([#32](https://github.com/alibaba/terraform-provider/pull/32))
  * resource/alicloud_slb: add udp listener, remove instance_protocol in listener ([#24](https://github.com/alibaba/terraform-provider/pull/24))
  * resource/alicloud_slb: fix slb bandwidth bug and modify listener default bandwidth ([#20](https://github.com/alibaba/terraform-provider/pull/20))
  

## 1.0.0(December 6, 2016)

  * **New Resource:** `alicloud_instance`
  * **New Resource:** `alicloud_security_group`
  * **New Resource:** `alicloud_slb`
  * **New Resource:** `alicloud_eip`
  * **New Resource:** `alicloud_vpc`
  * **New Resource:** `alicloud_vswitch`
  * **New Resource:** `alicloud_nat_gateway`
