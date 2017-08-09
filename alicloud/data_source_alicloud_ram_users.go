package alicloud

import (
	"fmt"
	"log"
	"regexp"

	"github.com/denverdino/aliyungo/ram"
	"github.com/hashicorp/terraform/helper/schema"
)

func dataSourceAlicloudRamUsers() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceAlicloudRamUsersRead,

		Schema: map[string]*schema.Schema{
			"type": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
				ValidateFunc: func(v interface{}, k string) (ws []string, errors []error) {
					value := v.(string)
					if value != "group" && value != "policy" {
						errors = append(errors, fmt.Errorf("%q must be 'group' or 'policy'.", k))
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
			"user_name_regex": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
			},
			"output_file": {
				Type:     schema.TypeString,
				Optional: true,
			},

			// Computed values
			"users": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"user_id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"user_name": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"create_date": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"last_login_date": {
							Type:     schema.TypeString,
							Computed: true,
						},
					},
				},
			},
		},
	}
}

func dataSourceAlicloudRamUsersRead(d *schema.ResourceData, meta interface{}) error {
	conn := meta.(*AliyunClient).ramconn
	var allUsers []ram.User

	if v, ok := d.GetOk("type"); ok {
		// users for this group
		if v.(string) == "group" {
			if v, ok = d.GetOk("group_name"); !ok {
				return fmt.Errorf("If 'type' value is 'group', you must set 'group_name' at one time.")
			} else {
				resp, err := conn.ListUsersForGroup(ram.GroupQueryRequest{GroupName: v.(string)})
				if err != nil {
					return fmt.Errorf("ListUsersForGroup got an error: %#v", err)
				}
				allUsers = append(allUsers, resp.Users.User...)
			}
		}

		// users which has this policy
		if v.(string) == "policy" {
			policyName, nameOk := d.GetOk("policy_name")
			policyType, typeOk := d.GetOk("policy_type")
			if !nameOk || !typeOk {
				return fmt.Errorf("If 'type' value is 'policy', you must set 'policy_name' and 'policy_type' at one time.")
			} else {
				resp, err := conn.ListEntitiesForPolicy(ram.PolicyRequest{PolicyName: policyName.(string), PolicyType: policyType.(string)})
				if err != nil {
					return fmt.Errorf("ListEntitiesForPolicy got an error: %#v", err)
				}
				allUsers = append(allUsers, resp.Users.User...)
			}
		}
	} else {
		args := ram.ListUserRequest{}
		for {
			resp, err := conn.ListUsers(args)
			if err != nil {
				return fmt.Errorf("ListUsers got an error: %#v", err)
			}
			allUsers = append(allUsers, resp.Users.User...)
			if !resp.IsTruncated {
				break
			}
			args.Marker = resp.Marker
		}
	}

	var filteredUsers []ram.User
	if v, ok := d.GetOk("user_name_regex"); ok && v.(string) != "" {
		r := regexp.MustCompile(v.(string))
		for _, user := range allUsers {
			if r.MatchString(user.UserName) {
				filteredUsers = append(filteredUsers, user)
			}
		}
	} else {
		filteredUsers = allUsers
	}

	if len(filteredUsers) < 1 {
		return fmt.Errorf("Your query returned no results. Please change your search criteria and try again.")
	}

	log.Printf("[DEBUG] alicloud_ram_users - Users found: %#v", allUsers)

	return ramUsersDecriptionAttributes(d, filteredUsers, meta)
}

func ramUsersDecriptionAttributes(d *schema.ResourceData, users []ram.User, meta interface{}) error {
	var ids []string
	var s []map[string]interface{}
	for _, user := range users {
		mapping := map[string]interface{}{
			"user_id":         user.UserId,
			"user_name":       user.UserName,
			"create_date":     user.CreateDate,
			"last_login_date": user.LastLoginDate,
		}
		log.Printf("[DEBUG] alicloud_ram_users - adding user: %v", mapping)
		ids = append(ids, user.UserId)
		s = append(s, mapping)
	}

	d.SetId(dataResourceIdHash(ids))
	if err := d.Set("users", s); err != nil {
		return err
	}

	// create a json file in current directory and write data source to it.
	if output, ok := d.GetOk("output_file"); ok && output != nil {
		writeToFile(output.(string), s)
	}
	return nil
}
