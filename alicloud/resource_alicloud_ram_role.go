package alicloud

import (
	"fmt"
	"strings"
	"time"

	"github.com/denverdino/aliyungo/ram"
	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/helper/schema"
)

func resourceAlicloudRamRole() *schema.Resource {
	return &schema.Resource{
		Create: resourceAlicloudRamRoleCreate,
		Read:   resourceAlicloudRamRoleRead,
		Update: resourceAlicloudRamRoleUpdate,
		Delete: resourceAlicloudRamRoleDelete,

		Schema: map[string]*schema.Schema{
			"role_name": &schema.Schema{
				Type:         schema.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: validateRamName,
			},
			"assume_role_policy": &schema.Schema{
				Type:         schema.TypeString,
				Required:     true,
				ValidateFunc: validateJsonString,
			},
			"description": &schema.Schema{
				Type:         schema.TypeString,
				Optional:     true,
				ForceNew:     true,
				ValidateFunc: validateRamDesc,
			},
			"force": &schema.Schema{
				Type:     schema.TypeBool,
				Optional: true,
				Default:  false,
			},
			"arn": &schema.Schema{
				Type:     schema.TypeString,
				Computed: true,
			},
			"create_date": &schema.Schema{
				Type:     schema.TypeString,
				Computed: true,
			},
			"update_date": &schema.Schema{
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}

func resourceAlicloudRamRoleCreate(d *schema.ResourceData, meta interface{}) error {
	conn := meta.(*AliyunClient).ramconn
	args := ram.RoleRequest{
		RoleName:                 d.Get("role_name").(string),
		AssumeRolePolicyDocument: d.Get("assume_role_policy").(string),
	}
	if v, ok := d.GetOk("description"); ok && v.(string) != "" {
		args.Description = v.(string)
	}

	response, err := conn.CreateRole(args)
	if err != nil {
		return fmt.Errorf("CreateRole got an error: %#v", err)
	}

	d.SetId(response.Role.RoleId)
	return resourceAlicloudRamRoleUpdate(d, meta)
}

func resourceAlicloudRamRoleUpdate(d *schema.ResourceData, meta interface{}) error {
	conn := meta.(*AliyunClient).ramconn

	d.Partial(true)

	args := ram.UpdateRoleRequest{
		RoleName: d.Get("role_name").(string),
	}

	if d.HasChange("assume_role_policy") && !d.IsNewResource() {
		d.SetPartial("assume_role_policy")
		args.NewAssumeRolePolicyDocument = d.Get("assume_role_policy").(string)
		if _, err := conn.UpdateRole(args); err != nil {
			return fmt.Errorf("UpdateRole got an error: %v", err)
		}
	}

	d.Partial(false)
	return resourceAlicloudRamRoleRead(d, meta)
}

func resourceAlicloudRamRoleRead(d *schema.ResourceData, meta interface{}) error {
	conn := meta.(*AliyunClient).ramconn

	args := ram.RoleQueryRequest{
		RoleName: d.Get("role_name").(string),
	}

	response, err := conn.GetRole(args)
	if err != nil {
		if RamEntityNotExist(err) {
			d.SetId("")
		}
		return fmt.Errorf("GetRole got an error: %v", err)
	}

	role := response.Role
	d.Set("arn", role.Arn)
	d.Set("role_name", role.RoleName)
	d.Set("create_date", role.CreateDate)
	d.Set("update_date", role.UpdateDate)
	d.Set("description", role.Description)
	rolePolicy := strings.Replace(strings.Replace(role.AssumeRolePolicyDocument, "\n", "", -1), " ", "", -1)
	d.Set("assume_role_policy", rolePolicy)
	return nil
}

func resourceAlicloudRamRoleDelete(d *schema.ResourceData, meta interface{}) error {
	conn := meta.(*AliyunClient).ramconn

	args := ram.RoleQueryRequest{
		RoleName: d.Get("role_name").(string),
	}

	if d.Get("force").(bool) {
		resp, err := conn.ListPoliciesForRole(args)
		if err != nil {
			return fmt.Errorf("Error listing Policies for Role (%s) when trying to delete: %#v", d.Id(), err)
		}

		// Loop and remove the Policies from the Role
		if len(resp.Policies.Policy) > 0 {
			for _, v := range resp.Policies.Policy {
				_, err = conn.DetachPolicyFromRole(ram.AttachPolicyToRoleRequest{
					PolicyRequest: ram.PolicyRequest{
						PolicyName: v.PolicyName,
						PolicyType: v.PolicyType,
					},
					RoleName: d.Get("role_name").(string),
				})
				if err != nil && !RamEntityNotExist(err) {
					return fmt.Errorf("Error detach Policy from Role %s: %#v", d.Id(), err)
				}
			}
		}
	}
	return resource.Retry(5*time.Minute, func() *resource.RetryError {
		if _, err := conn.DeleteRole(args); err != nil {
			if IsExceptedError(err, DeleteConflictRolePolicy) {
				return resource.RetryableError(fmt.Errorf("The role can not has any attached policy while deleting the role. - you can set force with true to force delete the role."))
			}
			return resource.NonRetryableError(fmt.Errorf("Error deleting role %s: %#v, you can set force with true to force delete the role.", d.Id(), err))
		}
		return nil
	})
}
