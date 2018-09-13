package alicloud

import (
	"fmt"
	"log"
	"regexp"

	"github.com/alibaba/terraform-provider/alicloud/aliyunclient"

	"github.com/denverdino/aliyungo/dns"
	"github.com/hashicorp/terraform/helper/schema"
)

func dataSourceAlicloudDnsGroups() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceAlicloudDnsGroupsRead,

		Schema: map[string]*schema.Schema{
			"name_regex": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
			},
			"output_file": {
				Type:     schema.TypeString,
				Optional: true,
			},

			// Computed values
			"groups": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"group_id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"group_name": {
							Type:     schema.TypeString,
							Computed: true,
						},
					},
				},
			},
		},
	}
}

func dataSourceAlicloudDnsGroupsRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*aliyunclient.AliyunClient)

	args := &dns.DescribeDomainGroupsArgs{}

	var allGroups []dns.DomainGroupType
	pagination := getPagination(1, 50)
	for {
		args.Pagination = pagination
		raw, err := client.RunSafelyWithDnsClient(func(dnsClient *dns.Client) (interface{}, error) {
			return dnsClient.DescribeDomainGroups(args)
		})
		if err != nil {
			return err
		}
		groups := raw.([]dns.DomainGroupType)
		allGroups = append(allGroups, groups...)

		if len(groups) < pagination.PageSize {
			break
		}
		pagination.PageNumber += 1
	}

	var filteredGroups []dns.DomainGroupType
	if v, ok := d.GetOk("name_regex"); ok && v.(string) != "" {
		r := regexp.MustCompile(v.(string))

		for _, group := range allGroups {
			if r.MatchString(group.GroupName) {
				filteredGroups = append(filteredGroups, group)
			}
		}
	} else {
		filteredGroups = allGroups[:]
	}

	if len(filteredGroups) < 1 {
		return fmt.Errorf("Your query returned no results. Please change your search criteria and try again.")
	}

	log.Printf("[DEBUG] alicloud_dns_groups - Groups found: %#v", allGroups)

	return groupsDecriptionAttributes(d, filteredGroups, meta)
}

func groupsDecriptionAttributes(d *schema.ResourceData, groupTypes []dns.DomainGroupType, meta interface{}) error {
	var ids []string
	var s []map[string]interface{}
	for _, group := range groupTypes {
		mapping := map[string]interface{}{
			"group_id":   group.GroupId,
			"group_name": group.GroupName,
		}
		log.Printf("[DEBUG] alicloud_dns_groups - adding group: %v", mapping)
		ids = append(ids, group.GroupId)
		s = append(s, mapping)
	}

	d.SetId(dataResourceIdHash(ids))
	if err := d.Set("groups", s); err != nil {
		return err
	}

	// create a json file in current directory and write data source to it.
	if output, ok := d.GetOk("output_file"); ok && output.(string) != "" {
		writeToFile(output.(string), s)
	}
	return nil
}
