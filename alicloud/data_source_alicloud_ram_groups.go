package alicloud

import (
	"fmt"
	"github.com/denverdino/aliyungo/ram"
	"github.com/hashicorp/terraform/helper/schema"
	"log"
	"regexp"
)

func dataSourceAlicloudRamGroups() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceAlicloudRamGroupsRead,

		Schema: map[string]*schema.Schema{
			"type": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
				ValidateFunc: func(v interface{}, k string) (ws []string, errors []error) {
					value := v.(string)
					if value != "user" && value != "policy" {
						errors = append(errors, fmt.Errorf("%q must be 'user' or 'policy'.", k))
					}
					return
				},
			},
			"user_name": {
				Type:         schema.TypeString,
				Optional:     true,
				ForceNew:     true,
				ValidateFunc: validateRamName,
			},
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
			"group_name_regex": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
			},
			"output_file": {
				Type:     schema.TypeString,
				Optional: true,
			},

			// Computed values
			"groups": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"group_name": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"comments": {
							Type:     schema.TypeString,
							Computed: true,
						},
					},
				},
			},
		},
	}
}

func dataSourceAlicloudRamGroupsRead(d *schema.ResourceData, meta interface{}) error {
	conn := meta.(*AliyunClient).ramconn
	var allGroups []ram.Group

	if v, ok := d.GetOk("type"); ok {
		// groups for this user
		if v.(string) == "user" {
			if v, ok = d.GetOk("user_name"); !ok {
				return fmt.Errorf("If 'type' value is 'user', you must set 'user_name' at one time.")
			} else {
				resp, err := conn.ListGroupsForUser(ram.UserQueryRequest{UserName: v.(string)})
				if err != nil {
					return fmt.Errorf("ListGroupsForUser got an error: %#v", err)
				}
				allGroups = append(allGroups, resp.Groups.Group...)
			}
		}

		// groups which has this policy
		if v.(string) == "policy" {
			policyName, nameOk := d.GetOk("policy_name")
			policyType, typeOk := d.GetOk("policy_type")
			if !nameOk || !typeOk {
				return fmt.Errorf("If 'type' value is 'policy', you must set 'policy_name' and 'policy_type' at one time.")
			} else {
				resp, err := conn.ListEntitiesForPolicy(ram.PolicyRequest{PolicyName: policyName.(string), PolicyType: ram.Type(policyType.(string))})
				if err != nil {
					return fmt.Errorf("ListEntitiesForPolicy got an error: %#v", err)
				}
				allGroups = append(allGroups, resp.Groups.Group...)
			}
		}
	} else {
		args := ram.GroupListRequest{}
		for {
			resp, err := conn.ListGroup(args)
			if err != nil {
				return fmt.Errorf("ListGroup got an error: %#v", err)
			}
			allGroups = append(allGroups, resp.Groups.Group...)
			if !resp.IsTruncated {
				break
			}
			args.Marker = resp.Marker
		}
	}

	var filteredGroups []ram.Group
	if v, ok := d.GetOk("group_name_regex"); ok && v.(string) != "" {
		r := regexp.MustCompile(v.(string))

		for _, group := range allGroups {
			if r.MatchString(group.GroupName) {
				filteredGroups = append(filteredGroups, group)
			}
		}
	} else {
		filteredGroups = allGroups[:]
	}

	if len(filteredGroups) < 1 {
		return fmt.Errorf("Your query returned no results. Please change your search criteria and try again.")
	}

	log.Printf("[DEBUG] alicloud_ram_groups - Groups found: %#v", allGroups)

	return ramGroupsDecriptionAttributes(d, filteredGroups, meta)
}

func ramGroupsDecriptionAttributes(d *schema.ResourceData, groups []ram.Group, meta interface{}) error {
	var ids []string
	var s []map[string]interface{}
	for _, group := range groups {
		mapping := map[string]interface{}{
			"group_name": group.GroupName,
			"comments":   group.Comments,
		}
		log.Printf("[DEBUG] alicloud_ram_groups - adding group: %v", mapping)
		ids = append(ids, group.GroupName)
		s = append(s, mapping)
	}

	d.SetId(dataResourceIdHash(ids))
	if err := d.Set("groups", s); err != nil {
		return err
	}

	// create a json file in current directory and write data source to it.
	if output, ok := d.GetOk("output_file"); ok && output != nil {
		writeToFile(output.(string), s)
	}
	return nil
}
