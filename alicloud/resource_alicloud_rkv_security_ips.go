package alicloud

import (
	"fmt"
	"strings"
	"time"

	"github.com/aliyun/alibaba-cloud-sdk-go/services/r-kvstore"
	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/helper/schema"
)

func resourceAlicloudRKVSecurityIPs() *schema.Resource {
	return &schema.Resource{
		Create: resourceAlicloudRKVSecurityIPsCreate,
		Read:   resourceAlicloudRKVSecurityIPsRead,
		Update: resourceAlicloudRKVSecurityIPsUpdate,
		Delete: resourceAlicloudRKVSecurityIPsDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Schema: map[string]*schema.Schema{
			"instance_id": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"security_ips": &schema.Schema{
				Type:     schema.TypeSet,
				Elem:     &schema.Schema{Type: schema.TypeString},
				Required: true,
			},
			"security_group_name": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
			},
		},
	}
}

func resourceAlicloudRKVSecurityIPsCreate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*AliyunClient)
	conn := client.rkvconn

	request := r_kvstore.CreateModifySecurityIpsRequest()
	request.InstanceId = d.Get("instance_id").(string)
	request.SecurityIpGroupName = d.Get("security_group_name").(string)
	request.SecurityIps = LOCAL_HOST_IP

	if len(d.Get("security_ips").(*schema.Set).List()) > 0 {
		request.SecurityIps = strings.Join(expandStringList(d.Get("security_ips").(*schema.Set).List())[:], COMMA_SEPARATED)
	}

	err := resource.Retry(5*time.Minute, func() *resource.RetryError {
		if _, err := conn.ModifySecurityIps(request); err != nil {
			return resource.NonRetryableError(fmt.Errorf("Create security whitelist ips got an error: %#v", err))
		}
		return nil
	})

	if err != nil {
		return err
	}

	// A security ip whitelist does not have a native IP.
	d.SetId(fmt.Sprintf("%s%s%s", request.InstanceId, COLON_SEPARATED, request.SecurityIpGroupName))

	return resourceAlicloudDBDatabaseRead(d, meta)
}

func resourceAlicloudRKVSecurityIPsRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*AliyunClient)
	conn := client.rkvconn

	request := r_kvstore.CreateDescribeSecurityIpsRequest()
	request.InstanceId = d.Get("instance_id").(string)
	attribs, err := conn.DescribeSecurityIps(request)
	if err != nil {
		if NotFoundRKVInstance(err) {
			d.SetId("")
			return nil
		}
		return fmt.Errorf("Error Describe RKV Security IPs: %#v", err)
	}

	if attribs == nil || len(attribs.SecurityIpGroups.SecurityIpGroup) == 0 {
		d.SetId("")
		return nil
	}

	secIP := &attribs.SecurityIpGroups.SecurityIpGroup[0]

	instanceID := strings.Split(d.Id(), COLON_SEPARATED)[0]
	d.Set("instance_id", instanceID)
	d.Set("security_group_name", secIP.SecurityIpGroupName)
	d.Set("security_ips", secIP.SecurityIpList)

	return nil
}

func resourceAlicloudRKVSecurityIPsUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*AliyunClient)
	conn := client.rkvconn

	if d.HasChange("security_group_name") && d.HasChange("security_ips") {
		ipstr := strings.Join(expandStringList(d.Get("security_ips").(*schema.Set).List())[:], COMMA_SEPARATED)
		if ipstr == "" {
			ipstr = LOCAL_HOST_IP
		}

		request := r_kvstore.CreateModifySecurityIpsRequest()
		request.SecurityIps = ipstr
		request.SecurityIpGroupName = d.Get("security_group_name").(string)
		if _, err := conn.ModifySecurityIps(request); err != nil {
			return err
		}
	}

	return resourceAlicloudRKVSecurityIPsRead(d, meta)
}

func resourceAlicloudRKVSecurityIPsDelete(d *schema.ResourceData, meta interface{}) error {
	// There is no explicit delete, only update with modified security ips
	return nil
}
