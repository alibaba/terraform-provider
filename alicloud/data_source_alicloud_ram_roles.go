package alicloud

import (
	"fmt"
	"log"
	"regexp"

	"github.com/denverdino/aliyungo/ram"
	"github.com/hashicorp/terraform/helper/schema"
)

func dataSourceAlicloudRamRoles() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceAlicloudRamRolesRead,

		Schema: map[string]*schema.Schema{
			"policy_name": {
				Type:         schema.TypeString,
				Optional:     true,
				ForceNew:     true,
				ValidateFunc: validateRamPolicyName,
			},
			"policy_type": {
				Type:         schema.TypeString,
				Optional:     true,
				ForceNew:     true,
				ValidateFunc: validatePolicyType,
			},
			"role_name_regex": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
			},
			"output_file": {
				Type:     schema.TypeString,
				Optional: true,
			},

			// Computed values
			"roles": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"role_id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"role_name": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"arn": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"description": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"assume_role_policy_document": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"create_date": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"update_date": {
							Type:     schema.TypeString,
							Computed: true,
						},
					},
				},
			},
		},
	}
}

func dataSourceAlicloudRamRolesRead(d *schema.ResourceData, meta interface{}) error {
	conn := meta.(*AliyunClient).ramconn
	var allRoles []ram.Role

	policyName, nameOk := d.GetOk("policy_name")
	policyType, typeOk := d.GetOk("policy_type")
	if nameOk && typeOk {
		resp, err := conn.ListEntitiesForPolicy(ram.PolicyRequest{PolicyName: policyName.(string), PolicyType: ram.Type(policyType.(string))})
		if err != nil {
			return fmt.Errorf("ListEntitiesForPolicy got an error: %#v", err)
		}
		allRoles = append(allRoles, resp.Roles.Role...)

	} else if !nameOk && !typeOk {
		resp, err := conn.ListRoles()
		if err != nil {
			return fmt.Errorf("ListRoles got an error: %#v", err)
		}
		allRoles = append(allRoles, resp.Roles.Role...)
	} else {
		return fmt.Errorf("you must set 'policy_name' and 'policy_type' at the same time.")
	}

	var filteredRoles []ram.Role
	if v, ok := d.GetOk("role_name_regex"); ok && v.(string) != "" {
		r := regexp.MustCompile(v.(string))

		for _, role := range allRoles {
			if r.MatchString(role.RoleName) {
				filteredRoles = append(filteredRoles, role)
			}
		}
	} else {
		filteredRoles = allRoles
	}

	if len(filteredRoles) < 1 {
		return fmt.Errorf("Your query returned no results. Please change your search criteria and try again.")
	}

	log.Printf("[DEBUG] alicloud_ram_roles - Roles found: %#v", allRoles)

	return ramRolesDecriptionAttributes(d, filteredRoles, meta)
}

func ramRolesDecriptionAttributes(d *schema.ResourceData, roles []ram.Role, meta interface{}) error {
	var ids []string
	var s []map[string]interface{}
	for _, role := range roles {
		mapping := map[string]interface{}{
			"role_id":                     role.RoleId,
			"role_name":                   role.RoleName,
			"arn":                         role.Arn,
			"description":                 role.Description,
			"create_date":                 role.CreateDate,
			"update_date":                 role.UpdateDate,
			"assume_role_policy_document": role.AssumeRolePolicyDocument,
		}
		log.Printf("[DEBUG] alicloud_ram_roles - adding role: %v", mapping)
		ids = append(ids, role.RoleId)
		s = append(s, mapping)
	}

	d.SetId(dataResourceIdHash(ids))
	if err := d.Set("roles", s); err != nil {
		return err
	}

	// create a json file in current directory and write data source to it.
	if output, ok := d.GetOk("output_file"); ok && output != nil {
		writeToFile(output.(string), s)
	}
	return nil
}
