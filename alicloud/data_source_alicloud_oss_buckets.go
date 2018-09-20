package alicloud

import (
	"github.com/hashicorp/terraform/helper/schema"
	"github.com/aliyun/aliyun-oss-go-sdk/oss"
	"regexp"
	"fmt"
	"log"
)

func dataSourceAlicloudOssBuckets() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceAlicloudOssBucketsRead,

		Schema: map[string]*schema.Schema{
			"name_regex": {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: validateNameRegex,
				ForceNew:     true,
			},
			"acl": {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: validateOssBucketAcl,
			},
			"output_file": {
				Type:     schema.TypeString,
				Optional: true,
			},

			// Computed values
			"buckets": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"name": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"acl": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"extranet_endpoint": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"intranet_endpoint": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"location": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"owner": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"storage_class": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"logging_enabled": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"creation_date": {
							Type:     schema.TypeString,
							Computed: true,
						},

						"cors_rules": {
							Type:     schema.TypeList,
							Computed: true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"allowed_headers": {
										Type:     schema.TypeString,
										Computed: true,
									},
									"allowed_methods": {
										Type:     schema.TypeString,
										Computed: true,
									},
									"allowed_origins": {
										Type:     schema.TypeString,
										Computed: true,
									},
									"expose_headers": {
										Type:     schema.TypeString,
										Computed: true,
									},
									"max_age_seconds": {
										Type:     schema.TypeInt,
										Computed: true,
									},
								},
							},
						},

						"website": {
							Type:     schema.TypeMap,
							Computed: true,
						},

						"logging": {
							Type:     schema.TypeMap,
							Computed: true,
						},

						"referer_config": {
							Type:     schema.TypeMap,
							Computed: true,
						},

						"lifecycle_rule": {
							Type:     schema.TypeList,
							Computed: true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"id": {
										Type:     schema.TypeString,
										Computed: true,
									},
									"prefix": {
										Type:     schema.TypeString,
										Computed: true,
									},
									"enabled": {
										Type:     schema.TypeBool,
										Computed: true,
									},
									"expiration": {
										Type:     schema.TypeMap,
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

func dataSourceAlicloudOssBucketsRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*AliyunClient)

	var initialOptions []oss.Option
	if v, ok := d.GetOk("acl"); ok && v.(string) != "" {
		acl := oss.ACLType(v.(string))
		initialOptions = append(initialOptions, oss.ACL(acl))
	}

	var allBuckets []oss.BucketProperties
	nextMarker := ""
	for {
		var options []oss.Option
		options = append(options, initialOptions...)
		if nextMarker != "" {
			options = append(options, oss.Marker(nextMarker))
		}

		resp, err := client.ossconn.ListBuckets(options...)
		if err != nil {
			return err
		}

		if resp.Buckets == nil || len(resp.Buckets) < 1 {
			break
		}

		allBuckets = append(allBuckets, resp.Buckets...)

		nextMarker = resp.NextMarker
		if nextMarker == "" {
			break
		}
	}

	var filteredBucketsTemp []oss.BucketProperties
	nameRegex, ok := d.GetOk("name_regex")
	if ok && nameRegex.(string) != "" {
		var r *regexp.Regexp
		if nameRegex != "" {
			r = regexp.MustCompile(nameRegex.(string))
		}
		for _, bucket := range allBuckets {
			if r != nil && !r.MatchString(bucket.Name) {
				continue
			}
			filteredBucketsTemp = append(filteredBucketsTemp, bucket)
		}
	} else {
		filteredBucketsTemp = allBuckets
	}

	if len(filteredBucketsTemp) < 1 {
		return fmt.Errorf("Your query returned no results. Please change your search criteria and try again.")
	}

	log.Printf("[DEBUG] alicloud_oss_buckets - Bucket found: %#v", filteredBucketsTemp)

	return bucketsDescriptionAttributes(d, filteredBucketsTemp, meta)
}

func bucketsDescriptionAttributes(d *schema.ResourceData, buckets []oss.BucketProperties, meta interface{}) error {
	client := meta.(*AliyunClient)

	var ids []string
	var s []map[string]interface{}
	for _, bucket := range buckets {
		mapping := map[string]interface{}{
			"name":          bucket.Name,
			"location":      bucket.Location,
			"storage_class": bucket.StorageClass,
			"creation_date": bucket.CreationDate.Format("2016-01-01"),
		}

		// Add additional information
		resp, err := client.ossconn.GetBucketInfo(bucket.Name)
		if err != nil {
			mapping["acl"] = resp.BucketInfo.ACL
			mapping["extranet_endpoint"] = resp.BucketInfo.ExtranetEndpoint
			mapping["intranet_endpoint"] = resp.BucketInfo.IntranetEndpoint
			mapping["owner"] = resp.BucketInfo.Owner.ID
		} else {
			log.Printf("[WARN] Unable to get additional information for the bucket %s: %v", bucket.Name, err)
		}

		// Add CORS rule information
		cors, err := client.ossconn.GetBucketCORS(bucket.Name)
		if err != nil && !IsExceptedErrors(err, []string{NoSuchCORSConfiguration}) {
			log.Printf("[WARN] Unable to get CORS information for the bucket %s: %v", bucket.Name, err)
		} else if err == nil && cors.CORSRules != nil {
			rules := make([]map[string]interface{}, 0, len(cors.CORSRules))
			for _, rule := range cors.CORSRules {
				ruleMapping := make(map[string]interface{})
				ruleMapping["allowed_headers"] = rule.AllowedHeader
				ruleMapping["allowed_methods"] = rule.AllowedMethod
				ruleMapping["allowed_origins"] = rule.AllowedOrigin
				ruleMapping["expose_headers"] = rule.ExposeHeader
				ruleMapping["max_age_seconds"] = rule.MaxAgeSeconds
				rules = append(rules, ruleMapping)
			}
			mapping["cors_rules"] = rules
		}

		// Add website configuration
		ws, err := client.ossconn.GetBucketWebsite(bucket.Name)
		if err != nil && !IsExceptedErrors(err, []string{NoSuchWebsiteConfiguration}) {
			log.Printf("[WARN] Unable to get website information for the bucket %s: %v", bucket.Name, err)
		} else if err == nil && &ws != nil {
			websiteMapping := make(map[string]interface{})
			if v := &ws.IndexDocument; v != nil {
				websiteMapping["index_document"] = v.Suffix
			}
			if v := &ws.ErrorDocument; v != nil {
				websiteMapping["error_document"] = v.Key
			}
			mapping["website"] = websiteMapping
		}

		// Add logging information
		logEnabled := false
		logging, err := client.ossconn.GetBucketLogging(bucket.Name)
		if err != nil {
			log.Printf("[WARN] Unable to get logging information for the bucket %s: %v", bucket.Name, err)
		} else if logging.LoggingEnabled.TargetBucket != "" || logging.LoggingEnabled.TargetPrefix != "" {
			logEnabled = true
			mapping["logging"] = map[string]interface{}{
				"target_bucket": logging.LoggingEnabled.TargetBucket,
				"target_prefix": logging.LoggingEnabled.TargetPrefix,
			}
		}
		mapping["logging_enabled"] = logEnabled

		// Add referer information
		referer, err := client.ossconn.GetBucketReferer(bucket.Name)
		if err != nil {
			log.Printf("[WARN] Unable to get referer information for the bucket %s: %v", bucket.Name, err)
		} else {
			mapping["referer_config"] = map[string]interface{}{
				"allow_empty": referer.AllowEmptyReferer,
				"referers":    referer.RefererList,
			}
		}

		// Add lifecycle information
		lifecycle, err := client.ossconn.GetBucketLifecycle(bucket.Name)
		if err != nil {
			log.Printf("[WARN] Unable to get lifecycle information for the bucket %s: %v", bucket.Name, err)
		} else if len(lifecycle.Rules) > 0 {
			ruleMappings := make([]map[string]interface{}, 0, len(lifecycle.Rules))

			for _, lifecycleRule := range lifecycle.Rules {
				ruleMapping := make(map[string]interface{})
				ruleMapping["id"] = lifecycleRule.ID
				ruleMapping["prefix"] = lifecycleRule.Prefix
				if LifecycleRuleStatus(lifecycleRule.Status) == ExpirationStatusEnabled {
					ruleMapping["enabled"] = true
				} else {
					ruleMapping["enabled"] = false
				}

				// Expiration
				expirationMapping := make(map[string]interface{})
				if !lifecycleRule.Expiration.Date.IsZero() {
					expirationMapping["date"] = (lifecycleRule.Expiration.Date).Format("2006-01-02")
				}
				if &lifecycleRule.Expiration.Days != nil {
					expirationMapping["days"] = int(lifecycleRule.Expiration.Days)
				}
				ruleMapping["expiration"] = expirationMapping
				ruleMappings = append(ruleMappings, ruleMapping)
			}
			mapping["lifecycle_rule"] = ruleMappings
		}

		log.Printf("[DEBUG] alicloud_oss_buckets - adding bucket mapping: %v", mapping)
		ids = append(ids, bucket.Name)
		s = append(s, mapping)
	}

	d.SetId(dataResourceIdHash(ids))
	if err := d.Set("instances", s); err != nil {
		return err
	}

	// create a json file in current directory and write data source to it.
	if output, ok := d.GetOk("output_file"); ok && output.(string) != "" {
		writeToFile(output.(string), s)
	}
	return nil
}
