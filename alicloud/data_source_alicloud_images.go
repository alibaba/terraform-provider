package alicloud

import (
	"fmt"
	"log"
	"regexp"
	"sort"

	"github.com/denverdino/aliyungo/ecs"
	"github.com/hashicorp/terraform/helper/schema"
	"strings"
	"time"
)

func dataSourceAlicloudImages() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceAlicloudImagesRead,

		Schema: map[string]*schema.Schema{
			//"filter": dataSourceFiltersSchema(),
			"name_regex": {
				Type:         schema.TypeString,
				Optional:     true,
				ForceNew:     true,
				ValidateFunc: validateNameRegex,
			},
			"most_recent": {
				Type:     schema.TypeBool,
				Optional: true,
				Default:  false,
				ForceNew: true,
			},
			"owners": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
			},
			// Computed values.
			"images": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"image_id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"architecture": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"creation_time": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"description": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"image_owner_alias": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"os_type": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"os_name": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"name": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"platform": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"is_self_shared": {
							Type:     schema.TypeBool,
							Computed: true,
						},
						"status": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"state": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"size": {
							Type:     schema.TypeInt,
							Computed: true,
						},
						// Complex computed values
						"disk_device_mappings": {
							Type:     schema.TypeList,
							Computed: true,
							//Set:      imageDiskDeviceMappingHash,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"device": {
										Type:     schema.TypeString,
										Computed: true,
									},
									"size": {
										Type:     schema.TypeString,
										Computed: true,
									},
									"snapshot_id": {
										Type:     schema.TypeString,
										Computed: true,
									},
								},
							},
						},
						"product_code": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"tags": tagsSchema(),
					},
				},
			},
		},
	}
}

// dataSourceAlicloudImagesDescriptionRead performs the Alicloud Image lookup.
func dataSourceAlicloudImagesRead(d *schema.ResourceData, meta interface{}) error {
	conn := meta.(*AliyunClient).ecsconn

	//filters, filtersOk := d.GetOk("filter")
	nameRegex, nameRegexOk := d.GetOk("name_regex")
	owners, ownersOk := d.GetOk("owners")

	//if executableUsersOk == false && filtersOk == false && nameRegexOk == false && ownersOk == false {
	//	return fmt.Errorf("One of executable_users, filters, name_regex, or owners must be assigned")
	//}
	if nameRegexOk == false && ownersOk == false {
		return fmt.Errorf("One of name_regex, or owners must be assigned")
	}

	params := &ecs.DescribeImagesArgs{
		RegionId: getRegion(d, meta),
	}

	//if filtersOk {
	//	params.Filters = buildAwsDataSourceFilters(filters.(*schema.Set))
	//}
	if ownersOk {
		params.ImageOwnerAlias = ecs.ImageOwnerAlias(owners.(string))
	}

	resp, _, err := conn.DescribeImages(params)
	if err != nil {
		return err
	}

	var filteredImages []ecs.ImageType
	if nameRegexOk {
		r := regexp.MustCompile(nameRegex.(string))
		for _, image := range resp {
			// Check for a very rare case where the response would include no
			// image name. No name means nothing to attempt a match against,
			// therefore we are skipping such image.
			if image.ImageName == "" {
				log.Printf("[WARN] Unable to find Image name to match against "+
					"for image ID %q, nothing to do.",
					image.ImageId)
				continue
			}
			if r.MatchString(image.ImageName) {
				filteredImages = append(filteredImages, image)
			}
		}
	} else {
		filteredImages = resp[:]
	}

	var images []ecs.ImageType
	if len(filteredImages) < 1 {
		return fmt.Errorf("Your query returned no results. Please change your search criteria and try again.")
	}

	recent := d.Get("most_recent").(bool)
	log.Printf("[DEBUG] alicloud_image - multiple results found and `most_recent` is set to: %t", recent)
	if len(filteredImages) > 1 && recent {
		// Query returned single result.
		images = append(images, mostRecentImage(filteredImages))
	} else {
		images = filteredImages
	}

	log.Printf("[DEBUG] alicloud_image - Images found: %#v", images)
	return imagesDescriptionAttributes(d, images)
}

// populate the numerous fields that the image description returns.
func imagesDescriptionAttributes(d *schema.ResourceData, images []ecs.ImageType) error {
	var id []string
	var s []map[string]interface{}
	for _, image := range images {
		mapping := map[string]interface{}{
			"id":                image.ImageId,
			"architecture":      image.Architecture,
			"creation_time":     image.CreationTime.String(),
			"description":       image.Description,
			"image_id":          image.ImageId,
			"image_owner_alias": image.ImageOwnerAlias,
			"os_name":           image.OSName,
			//"os_type":           image.OSType,
			"name": image.ImageName,
			//"platform": image.Platform,
			//"is_self_shared", image.IsSelfShared,
			//"status":       image.Status,
			//"state":        image.Status,
			"size":         image.Size,
			"product_code": image.ProductCode,

			//d.Set("tags", tagsToMap(image.tags)),
			// Complex types get their own functions
			"disk_device_mappings": imageDiskDeviceMappings(image.DiskDeviceMappings.DiskDeviceMapping),
		}

		log.Printf("[DEBUG] alicloud_image - adding image mapping: %v", mapping)
		id = append(id, image.ImageId)
		s = append(s, mapping)
	}

	d.SetId(strings.Join(id, ";"))
	if err := d.Set("images", s); err != nil {
		return err
	}
	return nil
}

//Find most recent image
type imageSort []ecs.ImageType

func (a imageSort) Len() int      { return len(a) }
func (a imageSort) Swap(i, j int) { a[i], a[j] = a[j], a[i] }
func (a imageSort) Less(i, j int) bool {
	itime, _ := time.Parse(time.RFC3339, a[i].CreationTime.String())
	jtime, _ := time.Parse(time.RFC3339, a[j].CreationTime.String())
	return itime.Unix() < jtime.Unix()
}

// Returns the most recent Image out of a slice of images.
func mostRecentImage(images []ecs.ImageType) ecs.ImageType {
	sortedImages := images
	sort.Sort(imageSort(sortedImages))
	return sortedImages[len(sortedImages)-1]
}

// Returns a set of disk device mappings.
func imageDiskDeviceMappings(m []ecs.DiskDeviceMapping) []map[string]interface{} {
	var s []map[string]interface{}

	for _, v := range m {
		mapping := map[string]interface{}{
			"device":      v.Device,
			"size":        v.Size,
			"snapshot_id": v.SnapshotId,
		}

		log.Printf("[DEBUG] alicloud_image - adding disk device mapping: %v", mapping)
		s = append(s, mapping)
	}

	return s
}
