package alicloud

import (
	"github.com/denverdino/aliyungo/common"
	"github.com/hashicorp/terraform/helper/mutexkv"
	"github.com/hashicorp/terraform/helper/schema"
	"github.com/hashicorp/terraform/terraform"
	"os"
)

// Provider returns a schema.Provider for alicloud
func Provider() terraform.ResourceProvider {
	return &schema.Provider{
		Schema: map[string]*schema.Schema{
			"access_key": &schema.Schema{
				Type:        schema.TypeString,
				Required:    true,
				DefaultFunc: schema.EnvDefaultFunc("ALICLOUD_ACCESS_KEY", nil),
				Description: descriptions["access_key"],
			},
			"secret_key": &schema.Schema{
				Type:        schema.TypeString,
				Required:    true,
				DefaultFunc: schema.EnvDefaultFunc("ALICLOUD_SECRET_KEY", nil),
				Description: descriptions["secret_key"],
			},
			"region": &schema.Schema{
				Type:        schema.TypeString,
				Required:    true,
				DefaultFunc: schema.EnvDefaultFunc("ALICLOUD_REGION", DEFAULT_REGION),
				Description: descriptions["region"],
			},
		},
		DataSourcesMap: map[string]*schema.Resource{

			"alicloud_images":             dataSourceAlicloudImages(),
			"alicloud_regions":            dataSourceAlicloudRegions(),
			"alicloud_zones":              dataSourceAlicloudZones(),
			"alicloud_instance_types":     dataSourceAlicloudInstanceTypes(),
			"alicloud_vpcs":               dataSourceAlicloudVpcs(),
			"alicloud_key_pairs":          dataSourceAlicloudKeyPairs(),
			"alicloud_dns_domains":        dataSourceAlicloudDnsDomains(),
			"alicloud_dns_domain_groups":  dataSourceAlicloudDnsDomainGroups(),
			"alicloud_dns_domain_records": dataSourceAlicloudDnsDomainRecords(),
			"alicloud_ram_account_alias":  dataSourceAlicloudRamAccountAlias(),
			"alicloud_ram_groups":         dataSourceAlicloudRamGroups(),
			"alicloud_ram_users":          dataSourceAlicloudRamUsers(),
			"alicloud_ram_roles":          dataSourceAlicloudRamRoles(),
			"alicloud_ram_policies":       dataSourceAlicloudRamPolicies(),
		},
		ResourcesMap: map[string]*schema.Resource{
			"alicloud_instance":                               resourceAliyunInstance(),
			"alicloud_ram_role_attachment":                    resourceAlicloudRamRoleAttachment(),
			"alicloud_disk":                                   resourceAliyunDisk(),
			"alicloud_disk_attachment":                        resourceAliyunDiskAttachment(),
			"alicloud_security_group":                         resourceAliyunSecurityGroup(),
			"alicloud_security_group_rule":                    resourceAliyunSecurityGroupRule(),
			"alicloud_db_instance":                            resourceAlicloudDBInstance(),
			"alicloud_ess_scaling_group":                      resourceAlicloudEssScalingGroup(),
			"alicloud_ess_scaling_group_instances_attachment": resourceAlicloudEssScalingGroupInstancesAttachment(),
			"alicloud_ess_scaling_configuration":              resourceAlicloudEssScalingConfiguration(),
			"alicloud_ess_scaling_rule":                       resourceAlicloudEssScalingRule(),
			"alicloud_ess_schedule":                           resourceAlicloudEssSchedule(),
			"alicloud_vpc":                                    resourceAliyunVpc(),
			"alicloud_nat_gateway":                            resourceAliyunNatGateway(),
			//both subnet and vswith exists,cause compatible old version, and compatible aws habit.
			"alicloud_subnet":                      resourceAliyunSubnet(),
			"alicloud_vswitch":                     resourceAliyunSubnet(),
			"alicloud_route_entry":                 resourceAliyunRouteEntry(),
			"alicloud_snat_entry":                  resourceAliyunSnatEntry(),
			"alicloud_forward_entry":               resourceAliyunForwardEntry(),
			"alicloud_eip":                         resourceAliyunEip(),
			"alicloud_eip_association":             resourceAliyunEipAssociation(),
			"alicloud_slb":                         resourceAliyunSlb(),
			"alicloud_slb_attachment":              resourceAliyunSlbAttachment(),
			"alicloud_oss_bucket":                  resourceAlicloudOssBucket(),
			"alicloud_oss_bucket_object":           resourceAlicloudOssBucketObject(),
			"alicloud_dns_record":                  resourceAlicloudDnsRecord(),
			"alicloud_dns":                         resourceAlicloudDns(),
			"alicloud_dns_group":                   resourceAlicloudDnsGroup(),
			"alicloud_key_pair":                    resourceAlicloudKeyPair(),
			"alicloud_key_pair_attachment":         resourceAlicloudKeyPairAttachment(),
			"alicloud_ram_user":                    resourceAlicloudRamUser(),
			"alicloud_ram_access_key":              resourceAlicloudRamAccessKey(),
			"alicloud_ram_login_profile":           resourceAlicloudRamLoginProfile(),
			"alicloud_ram_group":                   resourceAlicloudRamGroup(),
			"alicloud_ram_role":                    resourceAlicloudRamRole(),
			"alicloud_ram_policy":                  resourceAlicloudRamPolicy(),
			"alicloud_ram_alias":                   resourceAlicloudRamAlias(),
			"alicloud_ram_group_membership":        resourceAlicloudRamGroupMembership(),
			"alicloud_ram_user_policy_attachment":  resourceAlicloudRamUserPolicyAtatchment(),
			"alicloud_ram_role_policy_attachment":  resourceAlicloudRamRolePolicyAttachment(),
			"alicloud_ram_group_policy_attachment": resourceAlicloudRamGroupPolicyAtatchment(),
			"alicloud_container_cluster":           resourceAlicloudContainerCluster(),
			"alicloud_cdn_domain":                  resourceAlicloudCdnDomain(),
		},

		ConfigureFunc: providerConfigure,
	}
}

func providerConfigure(d *schema.ResourceData) (interface{}, error) {
	accesskey, ok := d.GetOk("access_key")
	if !ok {
		accesskey = os.Getenv("ALICLOUD_ACCESS_KEY")
	}
	secretkey, ok := d.GetOk("secret_key")
	if !ok {
		secretkey = os.Getenv("ALICLOUD_SECRET_KEY")
	}
	region, ok := d.GetOk("region")
	if !ok {
		region = os.Getenv("ALICLOUD_REGION")
		if region == "" {
			region = DEFAULT_REGION
		}
	}
	config := Config{
		AccessKey: accesskey.(string),
		SecretKey: secretkey.(string),
		Region:    common.Region(region.(string)),
	}

	client, err := config.Client()
	if err != nil {
		return nil, err
	}

	return client, nil
}

// This is a global MutexKV for use within this plugin.
var alicloudMutexKV = mutexkv.NewMutexKV()

var descriptions map[string]string

func init() {
	descriptions = map[string]string{
		"access_key": "Access key of alicloud",
		"secret_key": "Secret key of alicloud",
		"region":     "Region of alicloud",
	}
}
