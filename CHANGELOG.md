## 1.4.0 (unreleased)

## 1.3.0 (December 6, 2017)
FEATURES:

  * **New Resource:** `alicloud_slb_listener` ([#290](https://github.com/alibaba/terraform-provider/pull/290))

IMPROVMENTS:
  * modify client and support endpoint: ([#290](https://github.com/alibaba/terraform-provider/pull/290))

## 1.2.11 (November 30, 2017)
IMPROVMENTS:
  * fix creating multiple route entries bug: ([#286](https://github.com/alibaba/terraform-provider/pull/286))
  * add creating multiple vpcs and vswitches test case: ([#286](https://github.com/alibaba/terraform-provider/pull/286))
  * modify ess scaling group maxsize/minsize/default_cooldown type to int pointer: ([#286](https://github.com/alibaba/terraform-provider/pull/286))

## 1.2.11 (November 30, 2017)
IMPROVMENTS:
  * fix creating multiple route entries bug: ([#286](https://github.com/alibaba/terraform-provider/pull/286))
  * add creating multiple vpcs and vswitches test case: ([#286](https://github.com/alibaba/terraform-provider/pull/286))
  * modify ess scaling group maxsize/minsize/default_cooldown type to int pointer: ([#286](https://github.com/alibaba/terraform-provider/pull/286))


## 1.2.11 (November 30, 2017)
IMPROVMENTS:
  * fix creating multiple route entries bug: ([#286](https://github.com/alibaba/terraform-provider/pull/286))
  * add creating multiple vpcs and vswitches test case: ([#286](https://github.com/alibaba/terraform-provider/pull/286))
  * modify ess scaling group maxsize/minsize/default_cooldown type to int pointer: ([#286](https://github.com/alibaba/terraform-provider/pull/286))

## 1.2.10 (November 16, 2017)
IMPROVMENTS:
  * fix slb listener max healthy check timeout and interval: ([#276](https://github.com/alibaba/terraform-provider/pull/276))


## 1.2.9 (November 16, 2017)
IMPROVMENTS:
  * fix retriving instance types bug: ([#270](https://github.com/alibaba/terraform-provider/pull/270))
  * add tags for ESS: ([#274](https://github.com/alibaba/terraform-provider/pull/274))


## 1.2.8 (November 12, 2017)
IMPROVMENTS:
  * add WaitForInstanceAsyn to ensure correct status before creating or deleting database([#267](https://github.com/alibaba/terraform-provider/pull/267))

## 1.2.7 (November 10, 2017)
IMPROVMENTS:
  * add keypair, ram and userdata for ess scaling configuration ([#265](https://github.com/alibaba/terraform-provider/pull/265))
  * add force_delete to delete scaling configuration when only one configuration is existing ([#265](https://github.com/alibaba/terraform-provider/pull/265))
  * add substitute to active another scaling configuration ([#265](https://github.com/alibaba/terraform-provider/pull/265))
  
## 1.2.6 (October 17, 2017)
IMPROVMENTS:
  * modify CPU to Core(s) and improve some test case ([#254](https://github.com/alibaba/terraform-provider/pull/254))

BUG FIXES:
  * fix datasource dns domain ttl type bug([#254](https://github.com/denverdino/aliyungo/pull/254))
  * fix security group rule destroy failed bug([#259](https://github.com/denverdino/aliyungo/pull/254))

## 1.2.5 (October 17, 2017)
IMPROVMENTS:
  * add nexthop type 'RouterInterface' ([#252](https://github.com/alibaba/terraform-provider/pull/252))

BUG FIXES:
  * fix route entry creating and deleting bug([#252](https://github.com/denverdino/aliyungo/pull/252))
  * fix security group rule creating defect([#252](https://github.com/denverdino/aliyungo/pull/252))

## 1.2.4 (September 29, 2017)
BUG FIXES:
  * fix OSS bucket bug([#241](https://github.com/denverdino/aliyungo/pull/241))


## 1.2.3 (September 21, 2017)
BUG FIXES:
  * fix SDK bug([#186](https://github.com/denverdino/aliyungo/pull/186))

## 1.2.2 (September 13, 2017)
IMPROVMENTS:
  * validate instance type ([#238](https://github.com/alibaba/terraform-provider/pull/238))
BUG FIXES:
  * set ess system disk category default([#238](https://github.com/alibaba/terraform-provider/pull/238))

## 1.2.1 (September 11, 2017)
BUG FIXES:
  * fix internet_charge_type diff bug ([#235](https://github.com/alibaba/terraform-provider/pull/235))
  * fix dns marshal bug ([#237](https://github.com/alibaba/terraform-provider/pull/237))

## 1.2.0 (September 9, 2017)
FEATURES:

  * **New Resource:** `alicloud_router_interface` ([#228](https://github.com/alibaba/terraform-provider/pull/228))

IMPROVEMENTS:
  * remove runinstance api ([#227](https://github.com/alibaba/terraform-provider/pull/227))

BUG FIXES:
  * fix security group egress rules diff bug ([#223](https://github.com/alibaba/terraform-provider/pull/223))

## 1.1.10 (August 31, 2017)
IMPROVEMENTS:
  * add role_name for instance ([#216](https://github.com/alibaba/terraform-provider/pull/216))
  * deprecate router_id from `alicloud_route_entry` ([#219](https://github.com/alibaba/terraform-provider/pull/219))
  * modify router_table_id to route_table_id in `alicloud_vpc` ([#219](https://github.com/alibaba/terraform-provider/pull/219))
  * add route_table_id output in `data.alicloud_vpcs` ([#219](https://github.com/alibaba/terraform-provider/pull/219))
  * modify 'output != nil' to 'output.(string) != ""' to fix bug for all datasource ([#219](https://github.com/alibaba/terraform-provider/pull/219))

## 1.1.9 (August 18, 2017)
FEATURES:

  * **New Resource:** `alicloud_ram_role_attachment` ([#204](https://github.com/alibaba/terraform-provider/pull/204))

IMPROVEMENTS:
  * improve RAM role and policy's document. ([#204](https://github.com/alibaba/terraform-provider/pull/204))
  * modify some word misspellings. ([#205](https://github.com/alibaba/terraform-provider/pull/205))

## 1.1.8 (August 9, 2017)
FEATURES:

  * **New Resource:** `alicloud_container_cluster` ([#197](https://github.com/alibaba/terraform-provider/pull/197))
  * **New Resource:** `alicloud_cdn_domain` ([#198](https://github.com/alibaba/terraform-provider/pull/198))
  * **New DataSource:** `alicloud_dns_domains` ([#187](https://github.com/alibaba/terraform-provider/pull/187))
  * **New DataSource:** `alicloud_dns_domain_groups` ([#187](https://github.com/alibaba/terraform-provider/pull/187))
  * **New DataSource:** `alicloud_dns_domin_records` ([#187](https://github.com/alibaba/terraform-provider/pull/187))
  * **New DataSource:** `alicloud_ram_account_alias` ([#187](https://github.com/alibaba/terraform-provider/pull/187))
  * **New DataSource:** `alicloud_ram_groups` ([#187](https://github.com/alibaba/terraform-provider/pull/187))
  * **New DataSource:** `alicloud_dns_policies` ([#187](https://github.com/alibaba/terraform-provider/pull/187))
  * **New DataSource:** `alicloud_dns_roles` ([#187](https://github.com/alibaba/terraform-provider/pull/187))
  * **New DataSource:** `alicloud_dns_users` ([#187](https://github.com/alibaba/terraform-provider/pull/187))
IMPROVEMENTS:
  * sort data source alicloud zones ([#197](https://github.com/alibaba/terraform-provider/pull/197))

## 1.1.7 (August 4, 2017)
BUG FIXES:
  * fix creating rds failed result from Error Code changed ([#191](https://github.com/alibaba/terraform-provider/pull/191))

IMPROVEMENTS:
  * add resetting rds password function ([#191](https://github.com/alibaba/terraform-provider/pull/191))

## 1.1.6 (July 27, 2017)
FEATURES:

  * **New Resource:** `alicloud_ram_user` ([#181](https://github.com/alibaba/terraform-provider/pull/181))
  * **New Resource:** `alicloud_ram_group` ([#181](https://github.com/alibaba/terraform-provider/pull/181))
  * **New Resource:** `alicloud_ram_access_key` ([#181](https://github.com/alibaba/terraform-provider/pull/181))
  * **New Resource:** `alicloud_ram_alias` ([#181](https://github.com/alibaba/terraform-provider/pull/181))
  * **New Resource:** `alicloud_ram_group_membership` ([#181](https://github.com/alibaba/terraform-provider/pull/181))
  * **New Resource:** `alicloud_ram_group_policy` ([#181](https://github.com/alibaba/terraform-provider/pull/181))
  * **New Resource:** `alicloud_ram_group_policy_attachment` ([#181](https://github.com/alibaba/terraform-provider/pull/181))
  * **New Resource:** `alicloud_ram_login_profile` ([#181](https://github.com/alibaba/terraform-provider/pull/181))
  * **New Resource:** `alicloud_ram_role` ([#181](https://github.com/alibaba/terraform-provider/pull/181))
  * **New Resource:** `alicloud_ram_role_policy_attachment` ([#181](https://github.com/alibaba/terraform-provider/pull/181))
  * **New Resource:** `alicloud_ram_user_policy_attachment` ([#181](https://github.com/alibaba/terraform-provider/pull/181))

## 1.1.5 (July 20, 2017)
FEATURES:

  * **New Resource:** `alicloud_key_pair` ([#177](https://github.com/alibaba/terraform-provider/pull/177))
  * **New Resource:** `alicloud_key_pair_attachment` ([#177](https://github.com/alibaba/terraform-provider/pull/177))
  * **New DataSource:** `alicloud_key_pairs` ([#177](https://github.com/alibaba/terraform-provider/pull/177))

IMPROVEMENTS:

  * add keypair for alicloud_instance ([#177](https://github.com/alibaba/terraform-provider/pull/177))
  * set some field as computed for rds ([#178](https://github.com/alibaba/terraform-provider/pull/178))

## 1.1.4 (July 17, 2017)
BUG FIXES:

  * fix diff error when authorizing security group rules ([#175](https://github.com/alibaba/terraform-provider/pull/175))

## 1.1.3 (July 14, 2017)
FEATURES:

  * **New Resource:** `alicloud_dns` ([#167](https://github.com/alibaba/terraform-provider/pull/167))
  * **New Resource:** `alicloud_dns_group` ([#167](https://github.com/alibaba/terraform-provider/pull/167))
  * **New Resource:** `alicloud_dns_record` ([#167](https://github.com/alibaba/terraform-provider/pull/167))

IMPROVEMENTS:

  * add terraform provider import function ([#172](https://github.com/alibaba/terraform-provider/pull/172))
  * filter instance types of ess and retain series III ([#168](https://github.com/alibaba/terraform-provider/pull/168))

## 1.1.2 (July 7, 2017)
IMPROVEMENTS:

  * add terraform user agent ([#151](https://github.com/alibaba/terraform-provider/pull/151))
  * filter instance types and retain series III ([#161](https://github.com/alibaba/terraform-provider/pull/161))
  * add data source's output function ([#161](https://github.com/alibaba/terraform-provider/pull/161))
  * set 'io_optimized' to deprecated ([#161](https://github.com/alibaba/terraform-provider/pull/161))
  * set 'device_name' to deprecated ([#155](https://github.com/alibaba/terraform-provider/pull/155))

BUG FIXES:

  * fix attaching multiple disks error ([#155](https://github.com/alibaba/terraform-provider/pull/155))
  * fix rds bug that security_ips will be changed at every time ([#157](https://github.com/alibaba/terraform-provider/pull/157))
  * fix eip bug that there is no ip address output after running 'terraform apply' at first time ([#160](https://github.com/alibaba/terraform-provider/pull/160))


## 1.1.1 (June 6, 2017)
FEATURES:

  * **New DataSource:** `alicloud_vpcs` ([#145](https://github.com/alibaba/terraform-provider/pull/145))

## 1.1.0 (June 2, 2017)
BUG FIXES:

  * add retry while creating vpc and vswitch result from vpc and vsw don't solve concurrent bug.([#135](https://github.com/alibaba/terraform-provider/pull/135))
  * change security_groups's attribution optional to required when creating ecs instance.([#135](https://github.com/alibaba/terraform-provider/pull/135))
  * fix security group rule bug about setting 'nic_type' to 'intranet' when security group in vpc or authorizing permission for source/dest security group.([#136](https://github.com/alibaba/terraform-provider/pull/136))

## 1.0.9 (May 19, 2017)
FEATURES:

  * **New Resource:** `alicloud_oss_bucket_object` ([#132](https://github.com/alibaba/terraform-provider/pull/132))

## 1.0.8 (May 15, 2017)
BUG FIXES:

  * Fix the bug that doesn't unify error code in master account and RAM.([#130](https://github.com/alibaba/terraform-provider/pull/130))
  * resource/resource_alicloud_slb: add WaitForListenerAysn before operation StartLoadBalancerListener([#130](https://github.com/alibaba/terraform-provider/pull/130))

## 1.0.7 (May 10, 2017)

BUG FIXES:

  * resource/resource_alicloud_instance: wait for instance creating successfully before allocate public ip.([#122](https://github.com/alibaba/terraform-provider/pull/122))
  * resource/provider: if ak or region is empty in the template, it will get ak or region from env.([#122](https://github.com/alibaba/terraform-provider/pull/122))
  * resource/resource_alicloud_nat_gateway: set bandwidth packages zone's attribution 'Compute' as true for solving diff bug.([#123](https://github.com/alibaba/terraform-provider/pull/123))

FEATURES:

  * **New Resource:** `alicloud_oss_bucket` ([#122](https://github.com/alibaba/terraform-provider/pull/122))

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
