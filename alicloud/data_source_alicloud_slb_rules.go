package alicloud

import (
	"github.com/hashicorp/terraform/helper/schema"
	"github.com/aliyun/alibaba-cloud-sdk-go/services/slb"
	"github.com/aliyun/alibaba-cloud-sdk-go/sdk/requests"
	"fmt"
	"log"
)

func dataSourceAlicloudSlbRules() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceAlicloudSlbRulesRead,

		Schema: map[string]*schema.Schema{
			"load_balancer_id": {
				Type:     schema.TypeString,
				Required: true,
			},
			"frontend_port": {
				Type:     schema.TypeInt,
				Required: true,
			},
			// TODO add more filters

			// Computed values
			"slb_rules": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"name": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"domain": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"url": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"server_group_id": {
							Type:     schema.TypeString,
							Computed: true,
						},

						// TODO add more attributes
					},
				},
			},
		},
	}
}

func dataSourceAlicloudSlbRulesRead(d *schema.ResourceData, meta interface{}) error {
	conn := meta.(*AliyunClient).slbconn

	args := slb.CreateDescribeRulesRequest()
	args.LoadBalancerId = d.Get("load_balancer_id").(string)
	args.ListenerPort = requests.NewInteger(d.Get("frontend_port").(int))

	resp, err := conn.DescribeRules(args)
	if err != nil {
		return fmt.Errorf("DescribeRules got an error: %#v", err)
	}
	if resp == nil {
		return fmt.Errorf("there is no SLB with the ID %s. Please change your search criteria and try again", args.LoadBalancerId)
	}

	if len(resp.Rules.Rule) < 1 {
		return fmt.Errorf("Your query returned no results. Please change your search criteria and try again.")
	}

	log.Printf("[DEBUG] alicloud_slb_rules - Slb rules found: %#v", resp.Rules.Rule)

	return slbRulesDescriptionAttributes(d, resp.Rules.Rule, meta)
}

func slbRulesDescriptionAttributes(d *schema.ResourceData, rules []slb.Rule, meta interface{}) error {
	var ids []string
	var s []map[string]interface{}

	for _, rule := range rules {
		mapping := map[string]interface{}{
			"id":              rule.RuleId,
			"name":            rule.RuleName,
			"domain":          rule.Domain,
			"url":             rule.Url,
			"server_group_id": rule.VServerGroupId,
		}

		log.Printf("[DEBUG] alicloud_slb_rules - adding slb_ruler mapping: %v", mapping)
		ids = append(ids, rule.RuleId)
		s = append(s, mapping)
	}

	d.SetId(dataResourceIdHash(ids))
	if err := d.Set("slb_rules", s); err != nil {
		return err
	}

	// create a json file in current directory and write data source to it.
	if output, ok := d.GetOk("output_file"); ok && output.(string) != "" {
		writeToFile(output.(string), s)
	}
	return nil
}
