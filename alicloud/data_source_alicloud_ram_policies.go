package alicloud

import (
	"fmt"
	"log"
	"regexp"

	"github.com/denverdino/aliyungo/ram"
	"github.com/hashicorp/terraform/helper/schema"
)

func dataSourceAlicloudRamPolicies() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceAlicloudRamPoliciesRead,

		Schema: map[string]*schema.Schema{
			"policy_type": {
				Type:         schema.TypeString,
				Optional:     true,
				ForceNew:     true,
				ValidateFunc: validatePolicyType,
			},
			"type": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
				ValidateFunc: func(v interface{}, k string) (ws []string, errors []error) {
					value := v.(string)
					if value != "group" && value != "user" && value != "role" {
						errors = append(errors, fmt.Errorf("%q must be 'group' or 'user' or 'role'.", k))
					}
					return
				},
			},
			"group_name": {
				Type:         schema.TypeString,
				Optional:     true,
				ForceNew:     true,
				ValidateFunc: validateRamGroupName,
			},
			"user_name": {
				Type:         schema.TypeString,
				Optional:     true,
				ForceNew:     true,
				ValidateFunc: validateRamName,
			},
			"role_name": {
				Type:         schema.TypeString,
				Optional:     true,
				ForceNew:     true,
				ValidateFunc: validateRamName,
			},
			"policy_name_regex": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
			},
			"output_file": {
				Type:     schema.TypeString,
				Optional: true,
			},

			// Computed values
			"policies": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"policy_name": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"policy_type": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"description": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"default_version": {
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
						"attachment_count": {
							Type:     schema.TypeInt,
							Computed: true,
						},
						"policy_document": {
							Type:     schema.TypeString,
							Computed: true,
						},
					},
				},
			},
		},
	}
}

func dataSourceAlicloudRamPoliciesRead(d *schema.ResourceData, meta interface{}) error {
	conn := meta.(*AliyunClient).ramconn
	var allPolicies []ram.Policy

	if v, ok := d.GetOk("type"); ok {
		// policies for this user
		if v.(string) == "user" {
			if v, ok = d.GetOk("user_name"); !ok {
				return fmt.Errorf("If 'type' value is 'user', you must set 'user_name' at one time.")
			} else {
				resp, err := conn.ListPoliciesForUser(ram.UserQueryRequest{UserName: v.(string)})
				if err != nil {
					return fmt.Errorf("ListPoliciesForUser got an error: %#v", err)
				}
				allPolicies = append(allPolicies, resp.Policies.Policy...)
			}
		}

		// policies for this group
		if v.(string) == "group" {
			if v, ok = d.GetOk("group_name"); !ok {
				return fmt.Errorf("If 'type' value is 'group', you must set 'group_name' at one time.")
			} else {
				resp, err := conn.ListPoliciesForGroup(ram.GroupQueryRequest{GroupName: v.(string)})
				if err != nil {
					return fmt.Errorf("ListPoliciesForGroup got an error: %#v", err)
				}
				allPolicies = append(allPolicies, resp.Policies.Policy...)
			}
		}

		// policies for this role
		if v.(string) == "role" {
			if v, ok = d.GetOk("role_name"); !ok {
				return fmt.Errorf("If 'type' value is 'role', you must set 'role_name' at one time.")
			} else {
				resp, err := conn.ListPoliciesForRole(ram.RoleQueryRequest{RoleName: v.(string)})
				if err != nil {
					return fmt.Errorf("ListPoliciesForRole got an error: %#v", err)
				}
				allPolicies = append(allPolicies, resp.Policies.Policy...)
			}
		}
	} else {
		args := ram.PolicyQueryRequest{}
		for {
			resp, err := conn.ListPolicies(args)
			if err != nil {
				return fmt.Errorf("ListPolicies got an error: %#v", err)
			}
			allPolicies = append(allPolicies, resp.Policies.Policy...)
			if !resp.IsTruncated {
				break
			}
			args.Marker = resp.Marker
		}
	}

	var filteredPolicies []ram.Policy
	for _, policy := range allPolicies {
		if v, ok := d.GetOk("policy_type"); ok && policy.PolicyType != v.(string) {
			continue
		}
		if v, ok := d.GetOk("policy_name_regex"); ok && v.(string) != "" {
			r := regexp.MustCompile(v.(string))
			if r.MatchString(policy.PolicyName) {
				continue
			}
		}
		filteredPolicies = append(filteredPolicies, policy)
	}

	if len(filteredPolicies) < 1 {
		return fmt.Errorf("Your query returned no results. Please change your search criteria and try again.")
	}

	log.Printf("[DEBUG] alicloud_ram_policies - Policies found: %#v", allPolicies)

	return ramPoliciesDecriptionAttributes(d, filteredPolicies, meta)
}

func ramPoliciesDecriptionAttributes(d *schema.ResourceData, policies []ram.Policy, meta interface{}) error {
	conn := meta.(*AliyunClient).ramconn
	var ids []string
	var s []map[string]interface{}
	for _, policy := range policies {
		resp, err := conn.GetPolicyVersionNew(ram.PolicyRequest{
			PolicyName: policy.PolicyName,
			PolicyType: ram.Type(policy.PolicyType),
			VersionId:  policy.DefaultVersion,
		})
		if err != nil {
			return err
		}

		mapping := map[string]interface{}{
			"policy_name":      policy.PolicyName,
			"policy_type":      policy.PolicyType,
			"description":      policy.Description,
			"default_version":  policy.DefaultVersion,
			"attachment_count": int(policy.AttachmentCount),
			"create_date":      policy.CreateDate,
			"update_date":      policy.UpdateDate,
			"policy_document":  resp.PolicyVersion.PolicyDocument,
		}

		log.Printf("[DEBUG] alicloud_ram_policies - adding policy: %v", mapping)
		ids = append(ids, policy.PolicyName)
		s = append(s, mapping)
	}

	d.SetId(dataResourceIdHash(ids))
	if err := d.Set("policies", s); err != nil {
		return err
	}

	// create a json file in current directory and write data source to it.
	if output, ok := d.GetOk("output_file"); ok && output != nil {
		writeToFile(output.(string), s)
	}
	return nil
}
