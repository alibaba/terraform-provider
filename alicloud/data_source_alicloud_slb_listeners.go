package alicloud

import (
	"github.com/hashicorp/terraform/helper/schema"
	"github.com/aliyun/alibaba-cloud-sdk-go/services/slb"
	"fmt"
	"log"
	"strconv"
)

func dataSourceAlicloudSlbListeners() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceAlicloudSlbListenersRead,

		Schema: map[string]*schema.Schema{
			"load_balancer_id": {
				Type:     schema.TypeString,
				Required: true,
			},

			// TODO

			// Computed values
			"slb_listeners": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"port": {
							Type:     schema.TypeInt,
							Computed: true,
						},
						"protocol": {
							Type:     schema.TypeString,
							Computed: true,
						},
						// TODO
					},
				},
			},
		},
	}
}

func dataSourceAlicloudSlbListenersRead(d *schema.ResourceData, meta interface{}) error {
	conn := meta.(*AliyunClient).slbconn

	args := slb.CreateDescribeLoadBalancerAttributeRequest()
	args.LoadBalancerId = d.Get("load_balancer_id").(string)

	resp, err := conn.DescribeLoadBalancerAttribute(args)
	if err != nil {
		return fmt.Errorf("DescribeLoadBalancerAttribute got an error: %#v", err)
	}
	if resp == nil {
		return fmt.Errorf("there is no SLB with the ID %s. Please change your search criteria and try again", args.LoadBalancerId)
	}

	if len(resp.ListenerPortsAndProtocol.ListenerPortAndProtocol) < 1 {
		return fmt.Errorf("Your query returned no results. Please change your search criteria and try again.")
	}

	log.Printf("[DEBUG] alicloud_slb_listeners - Slb listeners found: %#v", resp.ListenerPortsAndProtocol.ListenerPortAndProtocol)

	return slbListenersDescriptionAttributes(d, resp.ListenerPortsAndProtocol.ListenerPortAndProtocol, meta)
}

func slbListenersDescriptionAttributes(d *schema.ResourceData, listeners []slb.ListenerPortAndProtocol, meta interface{}) error {
	var ids []string
	var s []map[string]interface{}

	for _, listener := range listeners {
		mapping := map[string]interface{}{
			"port":     listener.ListenerPort,
			"protocol": listener.ListenerProtocol,
		}

		// TODO get more info

		log.Printf("[DEBUG] alicloud_slb_listeners - adding slb_listener mapping: %v", mapping)
		ids = append(ids, strconv.Itoa(listener.ListenerPort))
		s = append(s, mapping)
	}

	d.SetId(dataResourceIdHash(ids))
	if err := d.Set("slb_listeners", s); err != nil {
		return err
	}

	// create a json file in current directory and write data source to it.
	if output, ok := d.GetOk("output_file"); ok && output.(string) != "" {
		writeToFile(output.(string), s)
	}
	return nil
}
