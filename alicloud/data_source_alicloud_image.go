package alicloud

import (
	"bytes"
	"fmt"
	"log"
	"regexp"
	"sort"

	"github.com/denverdino/aliyungo/ecs"
	"github.com/hashicorp/terraform/helper/hashcode"
	"github.com/hashicorp/terraform/helper/schema"
	"time"
)

func dataSourceAlicloudImage() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceAlicloudImageRead,

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
			"image_id": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"image_location": {
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
				Type:     schema.TypeString,
				Computed: true,
			},
			// Complex computed values
			"disk_device_mappings": {
				Type:     schema.TypeSet,
				Computed: true,
				Set:      imageDiskDeviceMappingHash,
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
	}
}

// dataSourceAlicloudImageDescriptionRead performs the Alicloud Image lookup.
func dataSourceAlicloudImageRead(d *schema.ResourceData, meta interface{}) error {
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

	var image ecs.ImageType
	if len(filteredImages) < 1 {
		return fmt.Errorf("Your query returned no results. Please change your search criteria and try again.")
	}

	if len(filteredImages) > 1 {
		recent := d.Get("most_recent").(bool)
		log.Printf("[DEBUG] alicloud_image - multiple results found and `most_recent` is set to: %t", recent)
		if recent {
			image = mostRecentImage(filteredImages)
		} else {
			return fmt.Errorf("Your query returned more than one result. Please try a more " +
				"specific search criteria, or set `most_recent` attribute to true.")
		}
	} else {
		// Query returned single result.
		image = filteredImages[0]
	}
	d.Set("image_location", getRegion(d, meta))

	log.Printf("[DEBUG] alicloud_image - Single Image found: %s", image.ImageId)
	return imageDescriptionAttributes(d, image)
}

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

// populate the numerous fields that the image description returns.
func imageDescriptionAttributes(d *schema.ResourceData, image ecs.ImageType) error {
	// Simple attributes first
	d.SetId(image.ImageId)
	d.Set("architecture", image.Architecture)
	d.Set("creation_time", image.CreationTime)
	d.Set("description", image.Description)
	d.Set("image_id", image.ImageId)
	d.Set("image_owner_alias", image.ImageOwnerAlias)
	d.Set("os_name", image.OSName)
	d.Set("name", image.ImageName)
	//if image.Platform != nil {
	//	d.Set("platform", image.Platform)
	//}
	//d.Set("is_self_shared", image.Public)
	d.Set("status", image.Status)
	d.Set("state", image.Status)
	d.Set("size", image.Size)
	d.Set("product_code", image.ProductCode)
	//d.Set("tags", tagsToMap(image.tags))
	// Complex types get their own functions
	if err := d.Set("disk_device_mappings", imageDiskDeviceMappings(image.DiskDeviceMappings.DiskDeviceMapping)); err != nil {
		return err
	}

	return nil
}

// Returns a set of disk device mappings.
func imageDiskDeviceMappings(m []ecs.DiskDeviceMapping) *schema.Set {
	s := &schema.Set{
		F: imageDiskDeviceMappingHash,
	}
	for _, v := range m {
		mapping := map[string]interface{}{
			"device":      v.Device,
			"size":        v.Size,
			"snapshot_id": v.SnapshotId,
		}

		log.Printf("[DEBUG] alicloud_image - adding disk device mapping: %v", mapping)
		s.Add(mapping)
	}

	return s
}

// Generates a hash for the set hash function used by the disk_device_mappings
// attribute.
func imageDiskDeviceMappingHash(v interface{}) int {
	var buf bytes.Buffer
	// All keys added in alphabetical order.
	m := v.(map[string]interface{})
	buf.WriteString(fmt.Sprintf("%s-", m["device"].(string)))
	buf.WriteString(fmt.Sprintf("%s-", m["size"].(string)))
	buf.WriteString(fmt.Sprintf("%s-", m["snapshot_id"].(string)))

	return hashcode.String(buf.String())
}
