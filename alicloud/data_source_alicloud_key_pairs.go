package alicloud

import (
	"fmt"
	"github.com/denverdino/aliyungo/common"
	"github.com/denverdino/aliyungo/ecs"
	"github.com/hashicorp/terraform/helper/schema"
	"log"
	"regexp"
)

func dataSourceAlicloudKeyPairs() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceAlicloudKeyPairsRead,

		Schema: map[string]*schema.Schema{
			"name_regex": &schema.Schema{
				Type:         schema.TypeString,
				Optional:     true,
				ForceNew:     true,
				ValidateFunc: validateNameRegex,
			},

			"finger_print": &schema.Schema{
				Type:     schema.TypeBool,
				Optional: true,
				Computed: true,
			},
			"output_file": {
				Type:     schema.TypeString,
				Optional: true,
			},

			//Computed value
			"keypairs": &schema.Schema{
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"key_name": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"finger_print": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"instances": {
							Type:     schema.TypeList,
							Computed: true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"instance_id": {
										Type:     schema.TypeString,
										Computed: true,
									},
									"instance_name": {
										Type:     schema.TypeString,
										Computed: true,
									},
									"instance_type": {
										Type:     schema.TypeString,
										Computed: true,
									},
									"vswitch_id": {
										Type:     schema.TypeString,
										Computed: true,
									},
								},
							},
						},
					},
				},
			},
		},
	}
}

func dataSourceAlicloudKeyPairsRead(d *schema.ResourceData, meta interface{}) error {
	conn := meta.(*AliyunClient).ecsconn

	var regex *regexp.Regexp
	if name, ok := d.GetOk("name_regex"); ok {
		regex = regexp.MustCompile(name.(string))
	}

	args := &ecs.DescribeKeyPairsArgs{
		RegionId: getRegion(d, meta),
	}
	if fingerPrint, ok := d.GetOk("finger_print"); ok {
		args.KeyPairFingerPrint = fingerPrint.(string)
	}
	var keyPairs []ecs.KeyPairItemType
	pagination := common.Pagination{
		PageSize: 50,
	}
	pagenumber := 1
	for true {
		pagination.PageNumber = pagenumber
		args.Pagination = pagination
		results, _, err := conn.DescribeKeyPairs(args)
		if err != nil {
			return fmt.Errorf("Error DescribekeyPairs: %#v", err)
		}
		for _, key := range results {
			if regex == nil || (regex != nil && regex.MatchString(key.KeyPairName)) {
				keyPairs = append(keyPairs, key)
			}
		}
		if len(results) < pagination.PageSize {
			break
		}
		pagenumber += 1
	}

	if len(keyPairs) < 1 {
		return fmt.Errorf("Your query key pairs returned no results. Please change your search criteria and try again.")
	}

	keyPairsAttach := make(map[string][]map[string]interface{})
	pagenumber = 1
	for true {
		pagination.PageNumber = pagenumber
		args.Pagination = pagination
		instances, _, err := conn.DescribeInstances(&ecs.DescribeInstancesArgs{
			RegionId: getRegion(d, meta),
		})
		if err != nil {
			return fmt.Errorf("Error DescribeInstances: %#v", err)
		}
		for _, inst := range instances {
			if inst.KeyPairName != "" {
				mapping := map[string]interface{}{
					"instance_id":   inst.InstanceId,
					"instance_name": inst.InstanceName,
					"instance_type": inst.InstanceType,
					"vswitch_id":    inst.VpcAttributes.VSwitchId,
				}
				if val, ok := keyPairsAttach[inst.KeyPairName]; ok {
					val = append(val, mapping)
					keyPairsAttach[inst.KeyPairName] = val
				} else {
					keyPairsAttach[inst.KeyPairName] = append(make([]map[string]interface{}, 0, 1), mapping)
				}
			}
		}
		if len(instances) < pagination.PageSize {
			break
		}
		pagenumber += 1
	}

	return keyPairsDescriptionAttributes(d, keyPairs, keyPairsAttach)
}

func keyPairsDescriptionAttributes(d *schema.ResourceData, keyPairs []ecs.KeyPairItemType, keyPairsAttach map[string][]map[string]interface{}) error {
	var names []string
	var s []map[string]interface{}
	for _, key := range keyPairs {
		mapping := map[string]interface{}{
			"id":           key.KeyPairName,
			"key_name":     key.KeyPairName,
			"finger_print": key.KeyPairFingerPrint,
			"instances":    keyPairsAttach[key.KeyPairName],
		}

		log.Printf("[DEBUG] alicloud_key_pairs - adding keypair mapping: %v", mapping)
		names = append(names, string(key.KeyPairName))
		s = append(s, mapping)
	}

	d.SetId(dataResourceIdHash(names))
	if err := d.Set("keypairs", s); err != nil {
		return err
	}

	// create a json file in current directory and write data source to it.
	if output, ok := d.GetOk("output_file"); ok && output != nil {
		writeToFile(output.(string), s)
	}
	return nil
}
