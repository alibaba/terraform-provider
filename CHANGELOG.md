## 1.9.7 (Unreleased)

IMPROVEMENTS:

- Support new field 'instance_name' for _alicloud_ots_table_ ([#548](https://github.com/alibaba/terraform-provider/pull/548)))
- *New Resource*: _alicloud_ots_instance_attachment_ ([#547](https://github.com/alibaba/terraform-provider/pull/547)))
- *New Resource*: _alicloud_ots_instance_ ([#546](https://github.com/alibaba/terraform-provider/pull/546)))

## 1.9.6 (July 16, 2018)

IMPROVEMENTS:

- Remove the number limitation of vswitch_ids, slb_ids and db_instance_ids ([#539](https://github.com/alibaba/terraform-provider/pull/539)))
- Reduce test nat gateway cost ([#538](https://github.com/alibaba/terraform-provider/pull/538)))
- Improve cs application resource test case ([#537](https://github.com/alibaba/terraform-provider/pull/537)))
- Improve ecs resource test case ([#536](https://github.com/alibaba/terraform-provider/pull/536)))
- Improve eip resource test case ([#535](https://github.com/alibaba/terraform-provider/pull/535)))
- Improve rds resource test case ([#534](https://github.com/alibaba/terraform-provider/pull/534)))
- Improve ess resource test case ([#533](https://github.com/alibaba/terraform-provider/pull/533)))
- Improve vpc and vswitch resource test case ([#532](https://github.com/alibaba/terraform-provider/pull/532)))
- Improve slb resource test case ([#531](https://github.com/alibaba/terraform-provider/pull/531)))
- Improve security group resource test case ([#530](https://github.com/alibaba/terraform-provider/pull/530)))
- Improve ram resource test case ([#529](https://github.com/alibaba/terraform-provider/pull/529)))
- Improve container cluster resource test case ([#528](https://github.com/alibaba/terraform-provider/pull/528)))
- Improve cloud monitor resource test case ([#527](https://github.com/alibaba/terraform-provider/pull/527)))
- Improve route resource test case ([#526](https://github.com/alibaba/terraform-provider/pull/526)))
- Improve nat gateway resource test case ([#526](https://github.com/alibaba/terraform-provider/pull/526)))
- Improve log resource test case ([#525](https://github.com/alibaba/terraform-provider/pull/525)))
- Improve ots resource test case ([#524](https://github.com/alibaba/terraform-provider/pull/524)))
- Improve dns resource test case ([#523](https://github.com/alibaba/terraform-provider/pull/523)))
- Improve oss resource test case ([#522](https://github.com/alibaba/terraform-provider/pull/522)))
- Support changing ecs charge type from Prepaid to PostPaid ([#521](https://github.com/alibaba/terraform-provider/pull/521)))
- Add method to compare json template is equal ([#508](https://github.com/alibaba/terraform-provider/pull/508)))

BUG FIXES:

- Fix CS kubernetes error and CS app timeout ([#528](https://github.com/alibaba/terraform-provider/pull/528)))
- Fix Oss bucket diff error ([#522](https://github.com/alibaba/terraform-provider/pull/522)))

## 1.9.5 (June 20, 2018)

IMPROVEMENTS:

- Support user agent for log service ([#504](https://github.com/alibaba/terraform-provider/pull/504)))
- Support sts token for some resources ([#504](https://github.com/alibaba/terraform-provider/pull/504)))
- *New Resource*: _alicloud_log_machine_group_ ([#503](https://github.com/alibaba/terraform-provider/pull/503)))
- *New Resource*: _alicloud_log_store_index_ ([#503](https://github.com/alibaba/terraform-provider/pull/503)))
- *New Resource*: _alicloud_log_store_ ([#503](https://github.com/alibaba/terraform-provider/pull/503)))
- *New Resource*: _alicloud_log_project_ ([#503](https://github.com/alibaba/terraform-provider/pull/503)))
- Improve example about cs_kubernetes ([#501](https://github.com/alibaba/terraform-provider/pull/501)))

## 1.9.4 (June 8, 2018)

IMPROVEMENTS:

- cs_kubernetes supports output worker nodes and master nodes ([#489](https://github.com/alibaba/terraform-provider/pull/489)))
- Add vendor ([#488](https://github.com/alibaba/terraform-provider/pull/488)))
- cs_kubernetes supports to output kube config and certificate ([#488](https://github.com/alibaba/terraform-provider/pull/488)))
- Add a example to deploy mysql and wordpress on kubernetes ([#488](https://github.com/alibaba/terraform-provider/pull/488)))
- Add a example to create swarm and deploy wordpress on it ([#487](https://github.com/alibaba/terraform-provider/pull/487)))
- Change ECS, ESS sdk to official go sdk ([#485](https://github.com/alibaba/terraform-provider/pull/485)))
- Add website in this repo ([#479](https://github.com/alibaba/terraform-provider/pull/479)))


## 1.9.3 (May 25, 2018)

IMPROVEMENTS:

- *New Data Source*: _alicloud_db_instances_ ([#478](https://github.com/alibaba/terraform-provider/pull/478)))
- Improve alicloud.tf ([#477](https://github.com/alibaba/terraform-provider/pull/477)))
- Support to set auto renew for ECS instance ([#476](https://github.com/alibaba/terraform-provider/pull/476)))
- Add missing code for describing RDS zones ([#475](https://github.com/alibaba/terraform-provider/pull/475)))
- Add filter parameters and export parameters for instance types data source. ([#472](https://github.com/alibaba/terraform-provider/pull/472)))
- Add filter parameters for zones data source. ([#472](https://github.com/alibaba/terraform-provider/pull/472)))
- Remove kubernetes work_number limitation ([#471](https://github.com/alibaba/terraform-provider/pull/471)))
- Improve kubernetes examples ([#471](https://github.com/alibaba/terraform-provider/pull/471)))

BUG FIXES:

- Fix getting some instance types failed bug ([#471](https://github.com/alibaba/terraform-provider/pull/471)))
- Fix kubernetes out range index error ([#470](https://github.com/alibaba/terraform-provider/pull/470)))

## 1.9.2 (May 8, 2018)

IMPROVEMENTS:

- *New Resource*: _alicloud_ots_table_ (([#442](https://github.com/alibaba/terraform-provider/pull/442)))
- Prolong waiting time for creating kubernetes cluster to avoid timeout (([#465](https://github.com/alibaba/terraform-provider/pull/465)))
- Support load endpoint from environment variable or specified file (([#462](https://github.com/alibaba/terraform-provider/pull/462)))
- Update example (([#457](https://github.com/alibaba/terraform-provider/pull/457)))
- Remove terraform/example (([#458](https://github.com/alibaba/terraform-provider/pull/458)))

BUG FIXES:

- Fix modifying instance host name failed bug ((([#465](https://github.com/alibaba/terraform-provider/pull/465)))
- Fix SLB listener "OperationBusy" error (([#465](https://github.com/alibaba/terraform-provider/pull/465)))
- Fix deleting forward table not found error (([#457](https://github.com/alibaba/terraform-provider/pull/457)))
- Fix deleting slb listener error (([#439](https://github.com/alibaba/terraform-provider/pull/439)))
- Fix creating vswitch error (([#439](https://github.com/alibaba/terraform-provider/pull/439)))

## 1.9.1 (April 15, 2018)

In order to be consistent with official, the following update will also be added into version 1.9.1.

IMPROVEMENTS:

- *New Resource*: _alicloud_cms_alarm_ (([#438](https://github.com/alibaba/terraform-provider/pull/438)))
- Output application attribution service block (([#435](https://github.com/alibaba/terraform-provider/pull/435)))
- Add kubernetes example (([#436](https://github.com/alibaba/terraform-provider/pull/436)))
- Add connections output for kubernetes cluster (([#437](https://github.com/alibaba/terraform-provider/pull/437)))

## 1.9.1 (March 30, 2018)

IMPROVEMENTS:

- *New Resource*: _alicloud_cs_application_ (([#419](https://github.com/alibaba/terraform-provider/pull/419)))
- *New Datasource*: _alicloud_security_group_rules_ (([#421](https://github.com/alibaba/terraform-provider/pull/421)))
- Output swarm attribution 'vpc_id' (([#434](https://github.com/alibaba/terraform-provider/pull/434)))
- Output swarm and kubernetes's nodes information and other attribution (([#420](https://github.com/alibaba/terraform-provider/pull/420)))
- Set swarm ID before waiting its status (([#419](https://github.com/alibaba/terraform-provider/pull/419)))
- Add is_outdated for cs_swarm and cs_kubernetes (([#418](https://github.com/alibaba/terraform-provider/pull/418)))
- Add warning when creating postgresql and ppas database (([#417](https://github.com/alibaba/terraform-provider/pull/417)))
- Add eip unassociation retry times to avoid needless error (([#437](https://github.com/alibaba/terraform-provider/pull/437)))

BUG FIXES:

- Fix vpc not found when vpc has been deleted (([#416](https://github.com/alibaba/terraform-provider/pull/416)))

## 1.9.0 (March 20, 2018)

- New release aims to keep in sync with official release version. This release is same as version 1.8.2.

## 1.8.2 (March 16, 2018)

IMPROVEMENTS:

- *New Resource*: _alicloud_cs_kubernetes_ (([#414](https://github.com/alibaba/terraform-provider/pull/414)))
- *New Datasource*: _alicloud_security_groups_ (([#415](https://github.com/alibaba/terraform-provider/pull/415)))
- Modify _alicloud_container_cluster_ to _alicloud_cs_swarm_ (([#414](https://github.com/alibaba/terraform-provider/pull/414)))
- SLB listener supports server_group_id (([#415](https://github.com/alibaba/terraform-provider/pull/415)))

BUG FIXES:

- Fix vpc description validate (([#411](https://github.com/alibaba/terraform-provider/pull/411)))
- Add waiting time after ECS bind ECS to ensure network is ok (([#414](https://github.com/alibaba/terraform-provider/pull/414)))

## 1.8.1 (March 9, 2018)

IMPROVEMENTS:

- Improve some examples (([#409](https://github.com/alibaba/terraform-provider/pull/409)))
- DB instance supports multiple zone (([#408](https://github.com/alibaba/terraform-provider/pull/408)))
- Data source zones support to retrieve multiple zone (([#407](https://github.com/alibaba/terraform-provider/pull/407)))
- Disk support encrypt (([#400](https://github.com/alibaba/terraform-provider/pull/400)))
- VPC supports alibaba cloud official go sdk (([#406](https://github.com/alibaba/terraform-provider/pull/406)))

BUG FIXES:

- Fix not found db instance bug when allocating connection (([#410](https://github.com/alibaba/terraform-provider/pull/410)))


## 1.8.0 (March 1, 2018)

IMPROVEMENTS:

- RDS supports alibaba cloud official go sdk (([#397](https://github.com/alibaba/terraform-provider/pull/397)))
- Deprecated 'in_use' in eips datasource to fix conflict (([#397](https://github.com/alibaba/terraform-provider/pull/397)))
- Support new region 'ap-southeast-5' and 'ap-south-1' (([#398](https://github.com/alibaba/terraform-provider/pull/398)))

BUG FIXES:

- Fix reading router interface failed bug (([#399](https://github.com/alibaba/terraform-provider/pull/399)))

## 1.7.2 (February 9, 2018)

IMPROVEMENTS:

- *New DataSource*: _alicloud_eips_ (([#386](https://github.com/alibaba/terraform-provider/pull/386)))
- *New DataSource*: _alicloud_vswitches_ (([#385](https://github.com/alibaba/terraform-provider/pull/385)))
- Support inner network segregation in one security group (([#390](https://github.com/alibaba/terraform-provider/pull/390)))

BUG FIXES:

- Fix creating Classic instance failed result in role_name (([#389](https://github.com/alibaba/terraform-provider/pull/389)))
- Fix eip is not exist in nat gateway when creating snat (([#384](https://github.com/alibaba/terraform-provider/pull/384)))

## 1.7.1 (February 2, 2018)

IMPROVEMENTS:

- Support setting instance_name for ESS scaling configuration (([#377](https://github.com/alibaba/terraform-provider/pull/377)))
- Support multiple vswitches for ESS scaling group and output slbIds and dbIds (([#376](https://github.com/alibaba/terraform-provider/pull/376)))
- Modify EIP default to PayByTraffic for international account (([#373](https://github.com/alibaba/terraform-provider/pull/373)))
- Deprecate nat gateway fileds 'spec' and 'bandwidth_packages' ([#368](https://github.com/alibaba/terraform-provider/pull/368))
- Support to associate EIP with SLB and Nat Gateway ([#367](https://github.com/alibaba/terraform-provider/pull/367))

BUG FIXES:

- fix a bug that can't create multiple VPC, vswitch and nat gateway at one time ([#374](https://github.com/terraform-providers/terraform-provider-alicloud/pull/374))
- fix a bug that can't import instance 'role_name' ([#375](https://github.com/terraform-providers/terraform-provider-alicloud/pull/375))
- fix a bug that creating ESS scaling group and configuration results from 'Throttling' ([#377](https://github.com/terraform-providers/terraform-provider-alicloud/pull/377))

## 1.7.0 (January 25, 2018)

IMPROVEMENTS:

- *New Resource*: _alicloud_kms_key_ ([#356](https://github.com/alibaba/terraform-provider/pull/356))
- *New DataSource*: _alicloud_kms_keys_ ([#357](https://github.com/alibaba/terraform-provider/pull/357))
- *New DataSource*: _alicloud_instances_ ([#358](https://github.com/alibaba/terraform-provider/pull/358))
- Add a new field "specification" for _alicloud_slb_ ([#358](https://github.com/alibaba/terraform-provider/pull/358))
- Improve security group rule's port range for "-1/-1" ([#359](https://github.com/alibaba/terraform-provider/pull/359))

BUG FIXES:

- fix slb invalid status error when launching ESS scaling group ([#360](https://github.com/alibaba/terraform-provider/pull/360))

## 1.6.2 (January 18, 2018)

IMPROVEMENTS:

- Support to set instnace name for RDS ([#353](https://github.com/alibaba/terraform-provider/pull/353))
- Avoid container cluster cidr block conflicts with vswitch's ([#352](https://github.com/alibaba/terraform-provider/pull/352))

BUG FIXES:

- fix several bugs about db result from its status and id not found ([#354](https://github.com/alibaba/terraform-provider/pull/354))
- fix deleting slb_attachment resource failed bug ([#351](https://github.com/alibaba/terraform-provider/pull/351))

## 1.6.1 (January 18, 2018)

IMPROVEMENTS:

- Support to modify instance type and network spec ([#344](https://github.com/alibaba/terraform-provider/pull/344))
- Avoid needless error when creating security group rule ([#344](https://github.com/alibaba/terraform-provider/pull/344))

BUG FIXES:

- fix creating cluster container failed bug ([#344](https://github.com/alibaba/terraform-provider/pull/344))

## 1.6.0 (January 15, 2018)

IMPROVMENTS:

  * *New Resource*: _alicloud_ess_attachment_  ([#341](https://github.com/alibaba/terraform-provider/pull/341))
  * *New Resource*: _alicloud_slb_rule_ ([#340](https://github.com/alibaba/terraform-provider/pull/340))
  * *New Resource*: _alicloud_slb_server_group_ ([#339](https://github.com/alibaba/terraform-provider/pull/339))
  * Standardize the order of imports packages ([#335](https://github.com/alibaba/terraform-provider/pull/335))
  * Output tip message when international account create SLB failed ([#336](https://github.com/alibaba/terraform-provider/pull/336))
  * Support spot instance ([#338](https://github.com/alibaba/terraform-provider/pull/338))
  * Add "weight" for slb_attachment to improve the resource ([#341](https://github.com/alibaba/terraform-provider/pull/341))

BUG FIXES:

  * fix allocating RDS public connection conflict error ([#337](https://github.com/terraform-providers/terraform-provider-alicloud/pull/337))


## 1.5.3 (January 9, 2018)

IMPROVMENTS:
  * support to go version 1.8.1 for travis ([#334](https://github.com/alibaba/terraform-provider/pull/334))

BUG FIXES:
  * fix getting OSS endpoint failed error  ([#332](https://github.com/alibaba/terraform-provider/pull/332))
  * fix describing dns record not found when deleting record ([#333](https://github.com/alibaba/terraform-provider/pull/333))


## 1.5.2 (January 8, 2018)

BUG FIXES:
  * fix creating rds 'Prepaid' instance failed error  ([#329](https://github.com/alibaba/terraform-provider/pull/329))

## 1.5.1 (January 5, 2018)

BUG FIXES:
  * roll back instance zone to compute ([#327](https://github.com/alibaba/terraform-provider/pull/327))
  * modify security_token to Optional ([#328](https://github.com/alibaba/terraform-provider/pull/328))


BUG FIXES:

  * fix allocating RDS public connection conflict error ([#336](https://github.com/terraform-providers/terraform-provider-alicloud/pull/336))


## 1.5.3 (January 9, 2018)

IMPROVMENTS:
  * support to go version 1.8.1 for travis ([#334](https://github.com/alibaba/terraform-provider/pull/334))

BUG FIXES:
  * fix getting OSS endpoint failed error  ([#332](https://github.com/alibaba/terraform-provider/pull/332))
  * fix describing dns record not found when deleting record ([#333](https://github.com/alibaba/terraform-provider/pull/333))


## 1.5.2 (January 8, 2018)

BUG FIXES:
  * fix creating rds 'Prepaid' instance failed error  ([#329](https://github.com/alibaba/terraform-provider/pull/329))

## 1.5.1 (January 5, 2018)

BUG FIXES:
  * roll back instance zone to compute ([#327](https://github.com/alibaba/terraform-provider/pull/327))
  * modify security_token to Optional ([#328](https://github.com/alibaba/terraform-provider/pull/328))


## 1.5.0 (January 3, 2018)

FEATURES:

  * **New Resource:** `alicloud_db_account` ([#324](https://github.com/alibaba/terraform-provider/pull/324))
  * **New Resource:** `alicloud_db_account_privilege` ([#324](https://github.com/alibaba/terraform-provider/pull/324))
  * **New Resource:** `alicloud_db_backup_policy` ([#324](https://github.com/alibaba/terraform-provider/pull/324))
  * **New Resource:** `alicloud_db_connection` ([#324](https://github.com/alibaba/terraform-provider/pull/324))
  * **New Resource:** `alicloud_db_database` ([#324](https://github.com/alibaba/terraform-provider/pull/324))


IMPROVMENTS:
  * support to modify instance spec including instnaceType, bandwidth ([#322](https://github.com/alibaba/terraform-provider/pull/322))
  * support to modify instance privateIp and vswitch ([#322](https://github.com/alibaba/terraform-provider/pull/322))
  * support to modify instance charge type ([#322](https://github.com/alibaba/terraform-provider/pull/322))
  * output more useful error message ([#322](https://github.com/alibaba/terraform-provider/pull/322))
  * support to modify load balance internet attribute ([#323](https://github.com/alibaba/terraform-provider/pull/323))
  * modify `alicloud_db_instance` ([#324](https://github.com/alibaba/terraform-provider/pull/324))
  * give import route entry tips when it already exist ([#325](https://github.com/alibaba/terraform-provider/pull/325))


BUG FIXES:
  * fix SLB not found when describing SLB ([#322](https://github.com/alibaba/terraform-provider/pull/322))
  * fix attaching disk failed ([#323](https://github.com/alibaba/terraform-provider/pull/323))
  * fix dns record deleting failed ([#324](https://github.com/alibaba/terraform-provider/pull/324))
  * fix route entry cannot be deleted ([#325](https://github.com/alibaba/terraform-provider/pull/325))

## 1.3.3 (December 14, 2017)

IMPROVMENTS:
  * wait for SLB active before return back ([#310](https://github.com/alibaba/terraform-provider/pull/310))

## 1.3.3 (December 14, 2017)

IMPROVMENTS:
  * wait for SLB active before return back ([#310](https://github.com/alibaba/terraform-provider/pull/310))

IMPROVMENTS:

  * *New Resource*: _alicloud_slb_server_group_ ([#339](https://github.com/alibaba/terraform-provider/pull/339))
  * Standardize the order of imports packages ([#335](https://github.com/alibaba/terraform-provider/pull/335))
  * Output tip message when international account create SLB failed ([#336](https://github.com/alibaba/terraform-provider/pull/336))
  * Support spot instance ([#338](https://github.com/alibaba/terraform-provider/pull/338))

BUG FIXES:

  * fix allocating RDS public connection conflict error ([#337](https://github.com/terraform-providers/terraform-provider-alicloud/pull/337))


## 1.5.3 (January 9, 2018)

IMPROVMENTS:
  * support to go version 1.8.1 for travis ([#334](https://github.com/alibaba/terraform-provider/pull/334))

BUG FIXES:
  * fix getting OSS endpoint failed error  ([#332](https://github.com/alibaba/terraform-provider/pull/332))
  * fix describing dns record not found when deleting record ([#333](https://github.com/alibaba/terraform-provider/pull/333))


## 1.5.2 (January 8, 2018)

BUG FIXES:
  * fix creating rds 'Prepaid' instance failed error  ([#329](https://github.com/alibaba/terraform-provider/pull/329))

## 1.5.1 (January 5, 2018)

BUG FIXES:
  * roll back instance zone to compute ([#327](https://github.com/alibaba/terraform-provider/pull/327))
  * modify security_token to Optional ([#328](https://github.com/alibaba/terraform-provider/pull/328))


## 1.5.0 (January 3, 2018)

FEATURES:

  * **New Resource:** `alicloud_db_account` ([#324](https://github.com/alibaba/terraform-provider/pull/324))
  * **New Resource:** `alicloud_db_account_privilege` ([#324](https://github.com/alibaba/terraform-provider/pull/324))
  * **New Resource:** `alicloud_db_backup_policy` ([#324](https://github.com/alibaba/terraform-provider/pull/324))
  * **New Resource:** `alicloud_db_connection` ([#324](https://github.com/alibaba/terraform-provider/pull/324))
  * **New Resource:** `alicloud_db_database` ([#324](https://github.com/alibaba/terraform-provider/pull/324))


IMPROVMENTS:
  * support to modify instance spec including instnaceType, bandwidth ([#322](https://github.com/alibaba/terraform-provider/pull/322))
  * support to modify instance privateIp and vswitch ([#322](https://github.com/alibaba/terraform-provider/pull/322))
  * support to modify instance charge type ([#322](https://github.com/alibaba/terraform-provider/pull/322))
  * output more useful error message ([#322](https://github.com/alibaba/terraform-provider/pull/322))
  * support to modify load balance internet attribute ([#323](https://github.com/alibaba/terraform-provider/pull/323))
  * modify `alicloud_db_instance` ([#324](https://github.com/alibaba/terraform-provider/pull/324))
  * give import route entry tips when it already exist ([#325](https://github.com/alibaba/terraform-provider/pull/325))


BUG FIXES:
  * fix SLB not found when describing SLB ([#322](https://github.com/alibaba/terraform-provider/pull/322))
  * fix attaching disk failed ([#323](https://github.com/alibaba/terraform-provider/pull/323))
  * fix dns record deleting failed ([#324](https://github.com/alibaba/terraform-provider/pull/324))
  * fix route entry cannot be deleted ([#325](https://github.com/alibaba/terraform-provider/pull/325))

## 1.3.3 (December 14, 2017)

IMPROVMENTS:
  * wait for SLB active before return back ([#310](https://github.com/alibaba/terraform-provider/pull/310))

## 1.3.2 (December 13, 2017)

IMPROVMENTS:
  * deprecated ram_alias and add ram_account_alias ([#305](https://github.com/alibaba/terraform-provider/pull/305))
  * deprecated dns_domain_groups and add dns_groups ([#305](https://github.com/alibaba/terraform-provider/pull/305))
  * deprecated dns_domain_records and add dns_records ([#305](https://github.com/alibaba/terraform-provider/pull/305))
  * add slb listener importing test ([#305](https://github.com/alibaba/terraform-provider/pull/305))

BUG FIXES:
  * fix dns records bug ([#305](https://github.com/alibaba/terraform-provider/pull/305))
  * fix ESS bind SLB failed bug ([#308](https://github.com/alibaba/terraform-provider/pull/308))
  * fix security group not found bug ([#308](https://github.com/alibaba/terraform-provider/pull/308))


## 1.3.1 (December 7, 2017)

IMPROVMENTS:
  * fix slb attachment failed and heath_check_domain diff ignore: ([#296](https://github.com/alibaba/terraform-provider/pull/296))
  * match sdk changes: ([#300](https://github.com/alibaba/terraform-provider/pull/300))

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
