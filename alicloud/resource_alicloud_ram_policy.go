package alicloud

import (
	"fmt"
	"time"

	"github.com/denverdino/aliyungo/ram"
	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/helper/schema"
)

func resourceAlicloudRamPolicy() *schema.Resource {
	return &schema.Resource{
		Create: resourceAlicloudRamPolicyCreate,
		Read:   resourceAlicloudRamPolicyRead,
		Update: resourceAlicloudRamPolicyUpdate,
		Delete: resourceAlicloudRamPolicyDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},

		Schema: map[string]*schema.Schema{
			"name": &schema.Schema{
				Type:         schema.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: validateRamPolicyName,
			},
			"statement": &schema.Schema{
				Type:     schema.TypeSet,
				Required: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"effect": &schema.Schema{
							Type:     schema.TypeString,
							Required: true,
							ValidateFunc: func(v interface{}, k string) (ws []string, errors []error) {
								value := Effect(v.(string))
								if value != Allow && value != Deny {
									errors = append(errors, fmt.Errorf(
										"%q must be '%s' or '%s'.", k, Allow, Deny))
								}
								return
							},
						},
						"action": &schema.Schema{
							Type:     schema.TypeSet,
							Required: true,
							Elem: &schema.Schema{
								Type: schema.TypeString,
							},
						},
						"resource": &schema.Schema{
							Type:     schema.TypeSet,
							Required: true,
							Elem: &schema.Schema{
								Type: schema.TypeString,
							},
						},
					},
				},
			},
			"description": &schema.Schema{
				Type:         schema.TypeString,
				Optional:     true,
				ForceNew:     true,
				ValidateFunc: validateRamDesc,
			},
			"version": &schema.Schema{
				Type:         schema.TypeString,
				Optional:     true,
				Default:      "1",
				ValidateFunc: validatePolicyDocVersion,
			},
			"force": &schema.Schema{
				Type:     schema.TypeBool,
				Optional: true,
				Default:  false,
			},
			"type": &schema.Schema{
				Type:     schema.TypeString,
				Computed: true,
			},
			"document": &schema.Schema{
				Type:     schema.TypeString,
				Computed: true,
			},
			"attachment_count": &schema.Schema{
				Type:     schema.TypeInt,
				Computed: true,
			},
		},
	}
}

func resourceAlicloudRamPolicyCreate(d *schema.ResourceData, meta interface{}) error {
	conn := meta.(*AliyunClient).ramconn

	document, err := AssemblePolicyDocument(d.Get("statement").(*schema.Set).List(), d.Get("version").(string))
	if err != nil {
		return err
	}

	args := ram.PolicyRequest{
		PolicyName:     d.Get("name").(string),
		PolicyDocument: document,
	}

	if v, ok := d.GetOk("description"); ok && v.(string) != "" {
		args.Description = v.(string)
	}

	response, err := conn.CreatePolicy(args)
	if err != nil {
		return fmt.Errorf("CreatePolicy got an error: %#v", err)
	}

	d.SetId(response.Policy.PolicyName)
	return resourceAlicloudRamPolicyUpdate(d, meta)
}

func resourceAlicloudRamPolicyUpdate(d *schema.ResourceData, meta interface{}) error {
	conn := meta.(*AliyunClient).ramconn
	d.Partial(true)

	args := ram.PolicyRequest{
		PolicyName:   d.Id(),
		SetAsDefault: "true",
	}

	attributeUpdate := false
	if d.HasChange("statement") {
		d.SetPartial("statement")
		attributeUpdate = true
	}
	if d.HasChange("version") {
		d.SetPartial("version")
		attributeUpdate = true
	}

	if !d.IsNewResource() && attributeUpdate {
		document, err := AssemblePolicyDocument(d.Get("statement").(*schema.Set).List(), d.Get("version").(string))
		if err != nil {
			return err
		}
		args.PolicyDocument = document

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
		PolicyName: d.Id(),
		PolicyType: ram.Custom,
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

	policyDocument, err := ParsePolicyDocument(policyVersionResp.PolicyVersion.PolicyDocument)
	if err != nil {
		return err
	}

	d.Set("name", policy.PolicyName)
	d.Set("type", policy.PolicyType)
	d.Set("description", policy.Description)
	d.Set("attachment_count", policy.AttachmentCount)
	d.Set("version", policyDocument.Version)
	d.Set("statement", policyDocument.Statement)
	d.Set("document", policyVersionResp.PolicyVersion.PolicyDocument)

	return nil
}

func resourceAlicloudRamPolicyDelete(d *schema.ResourceData, meta interface{}) error {
	conn := meta.(*AliyunClient).ramconn

	args := ram.PolicyRequest{
		PolicyName: d.Id(),
	}

	if d.Get("force").(bool) {
		args.PolicyType = ram.Custom

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
