package alicloud

import (
	"fmt"
	"github.com/denverdino/aliyungo/ecs"
	"github.com/hashicorp/terraform/helper/schema"
	"log"
	"regexp"
)

func dataSourceAlicloudSecurityGroups() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceAlicloudSecurityGroupsRead,
		Schema: map[string]*schema.Schema{
			"vpc_id": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
			},
			"security_group_name_regex": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
			},
			// Computed Values
			"security_groups": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"security_group_id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"security_group_name": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"region_id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"description": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"permissions": {
							Type:     schema.TypeList,
							Computed: true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"ip_protocol": {
										Type:     schema.TypeString,
										Computed: true,
									},
									"port_range": {
										Type:     schema.TypeString,
										Computed: true,
									},
									"source_security_group_id": {
										Type:     schema.TypeString,
										Computed: true,
									},
									"source_group_owner_account": {
										Type:     schema.TypeString,
										Computed: true,
									},
									"source_cidr_ip": {
										Type:     schema.TypeString,
										Computed: true,
									},
									"dest_cidr_ip": {
										Type:     schema.TypeString,
										Computed: true,
									},
									"policy": {
										Type:     schema.TypeString,
										Computed: true,
									},
									"nic_type": {
										Type:     schema.TypeString,
										Computed: true,
									},
									"priority": {
										Type:     schema.TypeInt,
										Computed: true,
									},
									"direction": {
										Type:     schema.TypeString,
										Computed: true,
									},
									"description": {
										Type:     schema.TypeString,
										Computed: true,
									},
								},
							},
						},
						"vpc_id": {
							Type:     schema.TypeString,
							Computed: true,
						},
					},
				},
			},
		},
	}
}

func dataSourceAlicloudSecurityGroupsRead(d *schema.ResourceData, meta interface{}) error {
	conn := meta.(*AliyunClient).ecsconn

	securityGroupArgs := &ecs.DescribeSecurityGroupsArgs{
		RegionId: getRegion(d, meta),
	}

	if vpcId, ok := d.GetOk("vpc_id"); ok {
		securityGroupArgs.VpcId = vpcId.(string)
	}

	var securityGroups []ecs.SecurityGroupItemType
	for {
		items, pageResult, err := conn.DescribeSecurityGroups(securityGroupArgs)

		if err != nil {
			return err
		}

		securityGroups = append(securityGroups, items...)

		pagination := pageResult.NextPage()
		if pagination == nil {
			break
		}

		securityGroupArgs.Pagination = *pagination
	}

	var allGroupAttrs []ecs.DescribeSecurityGroupAttributeResponse
	for _, securityGroup := range securityGroups {
		securityGroupAttrArgs := &ecs.DescribeSecurityGroupAttributeArgs{
			SecurityGroupId: securityGroup.SecurityGroupId,
			RegionId:        getRegion(d, meta),
		}

		resp, err := conn.DescribeSecurityGroupAttribute(securityGroupAttrArgs)
		if err != nil {
			return err
		}

		allGroupAttrs = append(allGroupAttrs, *resp)
	}

	var attrs []ecs.DescribeSecurityGroupAttributeResponse
	if nameRegex, ok := d.GetOk("security_group_name_regex"); ok {
		if r, err := regexp.Compile(nameRegex.(string)); err == nil {
			for _, attr := range allGroupAttrs {
				if r.MatchString(attr.SecurityGroupName) {
					attrs = append(attrs, attr)
				}
			}
		}
	} else {
		attrs = allGroupAttrs
	}

	if len(attrs) < 1 {
		return fmt.Errorf("Your query returned no results. Please change your search criteria and try again.")
	}

	log.Printf("[DEBUG] alicloud_security_groups - security groups found: %#v", attrs)

	return securityGroupsDescriptionAttributes(d, attrs, meta)
}

func securityGroupsDescriptionAttributes(data *schema.ResourceData, attrs []ecs.DescribeSecurityGroupAttributeResponse, meta interface{}) error {
	var ids []string
	var s []map[string]interface{}
	for _, attr := range attrs {
		mapping := map[string]interface{}{
			"security_group_id":   attr.SecurityGroupId,
			"security_group_name": attr.SecurityGroupName,
			"region_id":           attr.RegionId,
			"description":         attr.Description,
			// Complex types get their own functions
			"permissions": securityGroupPermissionMappings(attr.Permissions.Permission),
			"vpc_id":      attr.VpcId,
		}
		log.Printf("[DEBUG] alicloud_security_group - adding security group mapping: %v", mapping)
		ids = append(ids, attr.SecurityGroupId)
		s = append(s, mapping)
	}

	data.SetId(dataResourceIdHash(ids))
	if err := data.Set("security_groups", s); err != nil {
		return err
	}

	return nil
}
func securityGroupPermissionMappings(permissionTypes []ecs.PermissionType) []map[string]interface{} {
	var s []map[string]interface{}

	for _, permission := range permissionTypes {
		mapping := map[string]interface{}{
			"ip_protocol":                permission.IpProtocol,
			"port_range":                 permission.PortRange,
			"source_security_group_id":   permission.SourceGroupId,
			"source_group_owner_account": permission.SourceGroupOwnerAccount,
			"source_cidr_ip":             permission.SourceCidrIp,
			"dest_cidr_ip":               permission.DestCidrIp,
			"policy":                     permission.Policy,
			"nic_type":                   permission.NicType,
			"priority":                   permission.Priority,
			"direction":                  permission.Direction,
			"description":                permission.Description,
		}

		log.Printf("[DEBUG] alicloud_security_group - adding permisson mapping: %v", mapping)
		s = append(s, mapping)
	}
	return s
}
