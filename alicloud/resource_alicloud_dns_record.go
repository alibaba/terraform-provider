package alicloud

import (
	"fmt"
	"time"

	"github.com/denverdino/aliyungo/common"
	"github.com/denverdino/aliyungo/dns"
	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/helper/schema"
)

func resourceAlicloudDnsRecord() *schema.Resource {
	return &schema.Resource{
		Create: resourceAlicloudDnsRecordCreate,
		Read:   resourceAlicloudDnsRecordRead,
		Update: resourceAlicloudDnsRecordUpdate,
		Delete: resourceAlicloudDnsRecordDelete,

		Schema: map[string]*schema.Schema{
			"name": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
			},
			"host_record": {
				Type:         schema.TypeString,
				Required:     true,
				ValidateFunc: validateRR,
			},
			"type": {
				Type:         schema.TypeString,
				Required:     true,
				ValidateFunc: validateDomainRecordType,
			},
			"value": {
				Type:     schema.TypeString,
				Required: true,
				//ValidateFunc: validateDnsRecordValue,
			},
			"ttl": {
				Type:     schema.TypeInt,
				Optional: true,
				Default:  600,
			},
			"priority": {
				Type:     schema.TypeInt,
				Optional: true,
			},
			"routing": {
				Type:     schema.TypeString,
				Optional: true,
				Default:  "default",
			},
			"status": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"locked": {
				Type:     schema.TypeBool,
				Computed: true,
			},
		},
	}
}

func resourceAlicloudDnsRecordCreate(d *schema.ResourceData, meta interface{}) error {
	conn := meta.(*AliyunClient).dnsconn

	args := &dns.AddDomainRecordArgs{
		DomainName: d.Get("name").(string),
		RR:         d.Get("host_record").(string),
		Type:       d.Get("type").(string),
		Value:      d.Get("value").(string),
	}

	if v, ok := d.GetOk("ttl"); ok {
		args.TTL = int32(v.(int))
	}

	if v, ok := d.GetOk("priority"); ok {
		priority := v.(int)
		if priority > 10 || priority < 1 {
			return fmt.Errorf("Priority value is 1-10")
		}
		args.Priority = priority

	} else if args.Type == dns.MXRecord {
		return fmt.Errorf("MXRecord needs priority param")
	}

	response, err := conn.AddDomainRecord(args)
	if err != nil {
		return fmt.Errorf("AddDomainRecord got a error: %#v", err)
	}
	recordID := response.RecordId
	d.SetId(recordID)
	return resourceAlicloudDnsRecordUpdate(d, meta)
}

func resourceAlicloudDnsRecordUpdate(d *schema.ResourceData, meta interface{}) error {
	conn := meta.(*AliyunClient).dnsconn

	d.Partial(true)
	attributeUpdate := false
	args := &dns.UpdateDomainRecordArgs{
		RecordId: d.Id(),
		RR:       d.Get("host_record").(string),
		Type:     d.Get("type").(string),
		Value:    d.Get("value").(string),
	}

	requiredParams := []string{"host_record", "type", "value"}
	for _, v := range requiredParams {
		if d.HasChange(v) {
			d.SetPartial(v)
			attributeUpdate = true
		}
	}

	if d.HasChange("priority") {
		d.SetPartial("priority")
		args.Priority = int32(d.Get("priority").(int))
		attributeUpdate = true
	}

	if d.HasChange("ttl") {
		d.SetPartial("ttl")
		args.TTL = int32(d.Get("ttl").(int))
		attributeUpdate = true
	}

	if d.HasChange("routing") {
		d.SetPartial("routing")
		args.Line = d.Get("routing").(string)
		attributeUpdate = true
	}

	if attributeUpdate && !d.IsNewResource() {
		if _, err := conn.UpdateDomainRecord(args); err != nil {
			return err
		}
	}

	d.Partial(false)

	return resourceAlicloudDnsRecordRead(d, meta)
}

func resourceAlicloudDnsRecordRead(d *schema.ResourceData, meta interface{}) error {
	conn := meta.(*AliyunClient).dnsconn

	args := &dns.DescribeDomainRecordInfoArgs{
		RecordId: d.Id(),
	}
	response, err := conn.DescribeDomainRecordInfo(args)
	if err != nil {
		if NotFoundError(err) {
			d.SetId("")
			return nil
		}
		return err
	}

	record := response.RecordType
	d.Set("host_record", record.RR)
	d.Set("type", record.Type)
	d.Set("value", record.Value)
	d.Set("ttl", record.TTL)
	d.Set("priority", record.Priority)
	d.Set("routing", record.Line)
	d.Set("status", record.Status)
	d.Set("locked", record.Locked)

	return nil
}

func resourceAlicloudDnsRecordDelete(d *schema.ResourceData, meta interface{}) error {
	conn := meta.(*AliyunClient).dnsconn
	args := &dns.DeleteDomainRecordArgs{
		RecordId: d.Id(),
	}
	return resource.Retry(5*time.Minute, func() *resource.RetryError {
		_, err := conn.DeleteDomainRecord(args)
		if err != nil {
			e, _ := err.(*common.Error)
			if e.ErrorResponse.Code == RecordForbiddenDNSChange {
				return resource.RetryableError(fmt.Errorf("Operation forbidden because DNS is changing - trying again after change complete."))
			}
			return resource.NonRetryableError(fmt.Errorf("Error deleting domain record %s: %s", d.Id(), err))
		}
		return nil
	})
}