package alicloud

import (
	"fmt"
	"github.com/denverdino/aliyungo/ecs"
	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/helper/schema"
	"strings"
	"time"
)

func resourceAliyunSecurityGroupRule() *schema.Resource {
	return &schema.Resource{
		Create: resourceAliyunSecurityGroupRuleCreate,
		Read:   resourceAliyunSecurityGroupRuleRead,
		Delete: resourceAliyunSecurityGroupRuleDelete,

		Schema: map[string]*schema.Schema{
			"type": &schema.Schema{
				Type:         schema.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: validateSecurityRuleType,
				Description:  "Type of rule, ingress (inbound) or egress (outbound).",
			},

			"ip_protocol": &schema.Schema{
				Type:         schema.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: validateSecurityRuleIpProtocol,
			},

			"nic_type": &schema.Schema{
				Type:         schema.TypeString,
				Optional:     true,
				ForceNew:     true,
				Computed:     true,
				ValidateFunc: validateSecurityRuleNicType,
			},

			"policy": &schema.Schema{
				Type:         schema.TypeString,
				Optional:     true,
				ForceNew:     true,
				ValidateFunc: validateSecurityRulePolicy,
			},

			"port_range": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},

			"priority": &schema.Schema{
				Type:         schema.TypeInt,
				Optional:     true,
				ForceNew:     true,
				ValidateFunc: validateSecurityPriority,
			},

			"security_group_id": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},

			"cidr_ip": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
			},

			"source_security_group_id": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
			},

			"source_group_owner_account": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
			},
		},
	}
}

func resourceAliyunSecurityGroupRuleCreate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*AliyunClient)
	conn := client.ecsconn

	direction := d.Get("type").(string)
	sgId := d.Get("security_group_id").(string)
	ptl := d.Get("ip_protocol").(string)
	port := d.Get("port_range").(string)
	nicType := d.Get("nic_type").(string)

	var autherr error
	switch GroupRuleDirection(direction) {
	case GroupRuleIngress:
		args, err := buildAliyunSecurityIngressArgs(d, meta)
		if err != nil {
			return err
		}
		autherr = conn.AuthorizeSecurityGroup(args)
	case GroupRuleEgress:
		args, err := buildAliyunSecurityEgressArgs(d, meta)
		if err != nil {
			return err
		}
		autherr = conn.AuthorizeSecurityGroupEgress(args)
	default:
		return fmt.Errorf("Security Group Rule must be type 'ingress' or type 'egress'")
	}

	if autherr != nil {
		return fmt.Errorf(
			"Error authorizing security group rule type %s: %s",
			direction, autherr)
	}

	d.SetId(sgId + ":" + direction + ":" + ptl + ":" + port + ":" + nicType)

	return resourceAliyunSecurityGroupRuleRead(d, meta)
}

func resourceAliyunSecurityGroupRuleRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*AliyunClient)
	parts := strings.Split(d.Id(), ":")
	sgId := parts[0]
	direction := parts[1]
	ip_protocol := parts[2]
	port_range := parts[3]
	nic_type := parts[4]
	rule, err := client.DescribeSecurityGroupRule(sgId, direction, nic_type, ip_protocol, port_range)

	if err != nil {
		if NotFoundError(err) {
			d.SetId("")
			return nil
		}
		return fmt.Errorf("Error SecurityGroup rule: %#v", err)
	}

	d.Set("type", rule.Direction)
	d.Set("ip_protocol", strings.ToLower(string(rule.IpProtocol)))
	d.Set("nic_type", rule.NicType)
	d.Set("policy", strings.ToLower(string(rule.Policy)))
	d.Set("port_range", rule.PortRange)
	d.Set("priority", rule.Priority)
	d.Set("security_group_id", sgId)
	//support source and desc by type
	if GroupRuleDirection(direction) == GroupRuleIngress {
		d.Set("cidr_ip", rule.SourceCidrIp)
		d.Set("source_security_group_id", rule.SourceGroupId)
		d.Set("source_group_owner_account", rule.SourceGroupOwnerAccount)
	} else {
		d.Set("cidr_ip", rule.DestCidrIp)
		d.Set("source_security_group_id", rule.DestGroupId)
		d.Set("source_group_owner_account", rule.DestGroupOwnerAccount)
	}

	return nil
}

func deleteSecurityGroupRule(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*AliyunClient)
	ruleType := d.Get("type").(string)

	if GroupRuleDirection(ruleType) == GroupRuleIngress {
		args, err := buildAliyunSecurityIngressArgs(d, meta)
		if err != nil {
			return err
		}
		revokeArgs := &ecs.RevokeSecurityGroupArgs{
			AuthorizeSecurityGroupArgs: *args,
		}
		return client.RevokeSecurityGroup(revokeArgs)
	}

	args, err := buildAliyunSecurityEgressArgs(d, meta)

	if err != nil {
		return err
	}
	revokeArgs := &ecs.RevokeSecurityGroupEgressArgs{
		AuthorizeSecurityGroupEgressArgs: *args,
	}
	return client.RevokeSecurityGroupEgress(revokeArgs)
}

func resourceAliyunSecurityGroupRuleDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*AliyunClient)
	parts := strings.Split(d.Id(), ":")
	sgId, direction, ip_protocol, port_range, nic_type := parts[0], parts[1], parts[2], parts[3], parts[4]

	return resource.Retry(5*time.Minute, func() *resource.RetryError {
		err := deleteSecurityGroupRule(d, meta)

		if err != nil {
			resource.RetryableError(fmt.Errorf("Security group rule in use - trying again while it is deleted."))
		}

		_, err = client.DescribeSecurityGroupRule(sgId, direction, nic_type, ip_protocol, port_range)
		if err != nil {
			if NotFoundError(err) {
				return nil
			}
			return resource.NonRetryableError(err)
		}

		return resource.RetryableError(fmt.Errorf("Security group rule in use - trying again while it is deleted."))
	})

}

func checkCidrAndSourceGroupId(cidrIp, sourceGroupId string) error {
	if cidrIp == "" && sourceGroupId == "" {
		return fmt.Errorf("Either cidr_ip or source_security_group_id is required.")
	}

	if cidrIp != "" && sourceGroupId != "" {
		return fmt.Errorf("You should set only one value of cidr_ip or source_security_group_id.")
	}
	return nil
}
func buildAliyunSecurityIngressArgs(d *schema.ResourceData, meta interface{}) (*ecs.AuthorizeSecurityGroupArgs, error) {
	conn := meta.(*AliyunClient).ecsconn

	args := &ecs.AuthorizeSecurityGroupArgs{
		RegionId: getRegion(d, meta),
	}

	if v := d.Get("ip_protocol").(string); v != "" {
		args.IpProtocol = ecs.IpProtocol(v)
	}

	if v := d.Get("port_range").(string); v != "" {
		args.PortRange = v
	}

	if v := d.Get("policy").(string); v != "" {
		args.Policy = ecs.PermissionPolicy(v)
	}

	if v := d.Get("priority").(int); v != 0 {
		args.Priority = v
	}

	if v := d.Get("nic_type").(string); v != "" {
		args.NicType = ecs.NicType(v)
	}

	cidrIp := d.Get("cidr_ip").(string)
	sourceGroupId := d.Get("source_security_group_id").(string)
	if err := checkCidrAndSourceGroupId(cidrIp, sourceGroupId); err != nil {
		return nil, err
	}
	if cidrIp != "" {
		args.SourceCidrIp = cidrIp
	}

	if sourceGroupId != "" {
		args.SourceGroupId = sourceGroupId
	}

	if v := d.Get("source_group_owner_account").(string); v != "" {
		args.SourceGroupOwnerAccount = v
	}

	sgId := d.Get("security_group_id").(string)

	sgArgs := &ecs.DescribeSecurityGroupAttributeArgs{
		SecurityGroupId: sgId,
		RegionId:        getRegion(d, meta),
	}

	_, err := conn.DescribeSecurityGroupAttribute(sgArgs)
	if err != nil {
		return nil, fmt.Errorf("Error get security group %s error: %#v", sgId, err)
	}

	args.SecurityGroupId = sgId

	return args, nil
}

func buildAliyunSecurityEgressArgs(d *schema.ResourceData, meta interface{}) (*ecs.AuthorizeSecurityGroupEgressArgs, error) {
	conn := meta.(*AliyunClient).ecsconn

	args := &ecs.AuthorizeSecurityGroupEgressArgs{
		RegionId: getRegion(d, meta),
	}

	if v := d.Get("ip_protocol").(string); v != "" {
		args.IpProtocol = ecs.IpProtocol(v)
	}

	if v := d.Get("port_range").(string); v != "" {
		args.PortRange = v
	}

	if v := d.Get("policy").(string); v != "" {
		args.Policy = ecs.PermissionPolicy(v)
	}

	if v := d.Get("priority").(int); v != 0 {
		args.Priority = v
	}

	if v := d.Get("nic_type").(string); v != "" {
		args.NicType = ecs.NicType(v)
	}

	cidrIp := d.Get("cidr_ip").(string)
	sourceGroupId := d.Get("source_security_group_id").(string)
	if err := checkCidrAndSourceGroupId(cidrIp, sourceGroupId); err != nil {
		return nil, err
	}
	if cidrIp != "" {
		args.DestCidrIp = cidrIp
	}

	if sourceGroupId != "" {
		args.DestGroupId = sourceGroupId
	}

	if v := d.Get("source_group_owner_account").(string); v != "" {
		args.DestGroupOwnerAccount = v
	}

	sgId := d.Get("security_group_id").(string)

	sgArgs := &ecs.DescribeSecurityGroupAttributeArgs{
		SecurityGroupId: sgId,
		RegionId:        getRegion(d, meta),
	}

	_, err := conn.DescribeSecurityGroupAttribute(sgArgs)
	if err != nil {
		return nil, fmt.Errorf("Error get security group %s error: %#v", sgId, err)
	}

	args.SecurityGroupId = sgId

	return args, nil
}
