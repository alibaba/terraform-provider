package alicloud

import (
	"fmt"
	"github.com/denverdino/aliyungo/ram"
	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/helper/schema"
	"time"
)

func resourceAlicloudRamPolicy() *schema.Resource {
	return &schema.Resource{
		Create: resourceAlicloudRamPolicyCreate,
		Read:   resourceAlicloudRamPolicyRead,
		Update: resourceAlicloudRamPolicyUpdate,
		Delete: resourceAlicloudRamPolicyDelete,

		Schema: map[string]*schema.Schema{
			"policy_name": &schema.Schema{
				Type:         schema.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: validateRamPolicyName,
			},
			"policy_document": &schema.Schema{
				Type:         schema.TypeString,
				Required:     true,
				ValidateFunc: validateRamPolicyDoc,
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
			"default_version": &schema.Schema{
				Type:     schema.TypeString,
				Computed: true,
			},
			"policy_type": &schema.Schema{
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
			"attachment_count": &schema.Schema{
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}

func resourceAlicloudRamPolicyCreate(d *schema.ResourceData, meta interface{}) error {
	conn := meta.(*AliyunClient).ramconn

	args := ram.PolicyRequest{
		PolicyName:     d.Get("policy_name").(string),
		PolicyDocument: d.Get("policy_document").(string),
	}

	if v, ok := d.GetOk("description"); ok && v.(string) != "" {
		args.Description = v.(string)
	}

	response, err := conn.CreatePolicy(args)
	if err != nil {
		return fmt.Errorf("CreatePolicy got an error: %#v", err)
	}

	d.SetId(response.Policy.PolicyName)
	d.Set("policy_type", response.Policy.PolicyType)
	return resourceAlicloudRamPolicyUpdate(d, meta)
}

func resourceAlicloudRamPolicyUpdate(d *schema.ResourceData, meta interface{}) error {
	conn := meta.(*AliyunClient).ramconn
	d.Partial(true)

	if d.HasChange("policy_document") && !d.IsNewResource() {
		d.SetPartial("policy_document")
		args := ram.PolicyRequest{
			PolicyName:     d.Get("policy_name").(string),
			PolicyDocument: d.Get("policy_document").(string),
			SetAsDefault:   "true",
		}
		if _, err := conn.CreatePolicyVersion(args); err != nil {
			return fmt.Errorf("Error updating policy %s: %#v", d.Id(), err)
		}
	}

	d.Partial(false)

	return resourceAlicloudRamPolicyRead(d, meta)
}

func resourceAlicloudRamPolicyRead(d *schema.ResourceData, meta interface{}) error {
	conn := meta.(*AliyunClient).ramconn

	args := ram.PolicyRequest{
		PolicyType: d.Get("policy_type").(string),
		PolicyName: d.Get("policy_name").(string),
	}

	policyResp, err := conn.GetPolicy(args)
	if err != nil {
		if RamEntityNotExist(err) {
			d.SetId("")
		}
		return fmt.Errorf("GetPolicy got an error: %#v", err)
	}
	policy := policyResp.Policy

	args.VersionId = policy.DefaultVersion
	policyVersionResp, err := conn.GetPolicyVersionNew(args)
	if err != nil {
		return fmt.Errorf("GetPolicyVersion got an error: %#v", err)
	}

	d.Set("policy_name", policy.PolicyName)
	d.Set("policy_type", policy.PolicyType)
	d.Set("description", policy.Description)
	d.Set("create_date", policy.CreateDate)
	d.Set("update_date", policy.UpdateDate)
	d.Set("default_version", policy.DefaultVersion)
	d.Set("attachment_count", policy.AttachmentCount)
	d.Set("policy_document", policyVersionResp.PolicyVersion.PolicyDocument)

	return nil
}

func resourceAlicloudRamPolicyDelete(d *schema.ResourceData, meta interface{}) error {
	conn := meta.(*AliyunClient).ramconn

	args := ram.PolicyRequest{
		PolicyName: d.Get("policy_name").(string),
	}

	if d.Get("force").(bool) {
		args.PolicyType = d.Get("policy_type").(string)

		// list and detach entities for this policy
		response, err := conn.ListEntitiesForPolicy(args)
		if err != nil {
			return fmt.Errorf("Error listing entities for policy %s when trying to delete: %#v", d.Id(), err)
		}

		if len(response.Users.User) > 0 {
			for _, v := range response.Users.User {
				_, err := conn.DetachPolicyFromUser(ram.AttachPolicyRequest{
					PolicyRequest: args,
					UserName:      v.UserName,
				})
				if err != nil && !RamEntityNotExist(err) {
					return fmt.Errorf("Error detaching policy %s from user %s:%#v", d.Id(), v.UserId, err)
				}
			}
		}

		if len(response.Groups.Group) > 0 {
			for _, v := range response.Groups.Group {
				_, err := conn.DetachPolicyFromGroup(ram.AttachPolicyToGroupRequest{
					PolicyRequest: args,
					GroupName:     v.GroupName,
				})
				if err != nil && !RamEntityNotExist(err) {
					return fmt.Errorf("Error detaching policy %s from group %s:%#v", d.Id(), v.GroupName, err)
				}
			}
		}

		if len(response.Roles.Role) > 0 {
			for _, v := range response.Roles.Role {
				_, err := conn.DetachPolicyFromRole(ram.AttachPolicyToRoleRequest{
					PolicyRequest: args,
					RoleName:      v.RoleName,
				})
				if err != nil && !RamEntityNotExist(err) {
					return fmt.Errorf("Error detaching policy %s from role %s:%#v", d.Id(), v.RoleId, err)
				}
			}
		}

		// list and delete policy version which are not default
		pvResp, err := conn.ListPolicyVersionsNew(args)
		if err != nil {
			return fmt.Errorf("Error listing policy versions for policy %s:%#v", d.Id(), err)
		}
		if len(pvResp.PolicyVersions.PolicyVersion) > 1 {
			for _, v := range pvResp.PolicyVersions.PolicyVersion {
				if !v.IsDefaultVersion {
					args.VersionId = v.VersionId
					if _, err = conn.DeletePolicyVersion(args); err != nil && !RamEntityNotExist(err) {
						return fmt.Errorf("Error delete policy version %s for policy %s:%#v", v.VersionId, d.Id(), err)
					}
				}
			}
		}
	}

	return resource.Retry(5*time.Minute, func() *resource.RetryError {
		if _, err := conn.DeletePolicy(args); err != nil {
			if IsExceptedError(err, DeleteConflictPolicyUser) || IsExceptedError(err, DeleteConflictPolicyGroup) || IsExceptedError(err, DeleteConflictRolePolicy) {
				return resource.RetryableError(fmt.Errorf("The policy can not been attached to any user or group or role while deleting the policy. - you can set force with true to force delete the policy."))
			}
			if IsExceptedError(err, DeleteConflictPolicyVersion) {
				return resource.RetryableError(fmt.Errorf("The policy can not has any version except the defaul version. - you can set force with true to force delete the policy."))
			}
			return resource.NonRetryableError(fmt.Errorf("Error deleting policy %s: %#v", d.Id(), err))
		}
		return nil
	})
}
