package alicloud

import (
	"github.com/hashicorp/terraform/helper/schema"
)

func resourceAliyunSlbAttachment() *schema.Resource {
	return &schema.Resource{
		Create: resourceAliyunSlbAttachmentCreate,
		Read:   resourceAliyunSlbAttachmentRead,
		Delete: resourceAliyunSlbAttachmentDelete,

		Schema: map[string]*schema.Schema{

			"slb_id": &schema.Schema{
				Type:     schema.TypeSet,
				Elem:     &schema.Schema{Type: schema.TypeString},
				Required: true,
			},

			"instances": &schema.Schema{
				Type:     schema.TypeSet,
				Elem:     &schema.Schema{Type: schema.TypeString},
				Required: true,
				Set:      schema.HashString,
			},
		},
	}
}

func resourceAliyunSlbAttachmentCreate(d *schema.ResourceData, meta interface{}) error {

	slbId := d.Get("slb_id").(string)

	slbconn := meta.(*AliyunClient).slbconn

	loadBalancer, err := slbconn.DescribeLoadBalancerAttribute(slbId)
	if err != nil {
		if notFoundError(err) {
			d.SetId("")
			return nil
		}

		return err
	}

	d.SetId(loadBalancer.LoadBalancerId)

	o, n := d.GetChange("instances")
	os := o.(*schema.Set)
	ns := n.(*schema.Set)
	remove := expandBackendServers(os.Difference(ns).List())
	add := expandBackendServers(ns.Difference(os).List())

	if len(add) > 0 {
		_, err := slbconn.AddBackendServers(d.Id(), add)
		if err != nil {
			return err
		}
	}
	if len(remove) > 0 {
		removeBackendServers := make([]string, 0, len(remove))
		for _, e := range remove {
			removeBackendServers = append(removeBackendServers, e.ServerId)
		}
		_, err := slbconn.RemoveBackendServers(d.Id(), removeBackendServers)
		if err != nil {
			return err
		}
	}

	d.SetPartial("instances")

	return nil
}

func resourceAliyunSlbAttachmentRead(d *schema.ResourceData, meta interface{}) error {

	//todo

	return nil
}

func resourceAliyunSlbAttachmentDelete(d *schema.ResourceData, meta interface{}) error {
	//todo
	return nil
}
