package alicloud

import (
	"fmt"
	"github.com/denverdino/aliyungo/ecs"
	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/helper/schema"
	"io/ioutil"
	"strings"
)

func resourceAlicloudKeyPair() *schema.Resource {
	return &schema.Resource{
		Create: resourceAlicloudKeyPairCreate,
		Read:   resourceAlicloudKeyPairRead,
		Update: nil,
		Delete: resourceAlicloudKeyPairDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},

		Schema: map[string]*schema.Schema{
			"key_name": &schema.Schema{
				Type:          schema.TypeString,
				Optional:      true,
				Computed:      true,
				ForceNew:      true,
				ValidateFunc:  validateKeyPairName,
				ConflictsWith: []string{"key_name_prefix"},
			},
			"key_name_prefix": &schema.Schema{
				Type:         schema.TypeString,
				Optional:     true,
				ForceNew:     true,
				ValidateFunc: validateKeyPairPrefix,
			},
			"public_key": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
				StateFunc: func(v interface{}) string {
					switch v.(type) {
					case string:
						return strings.TrimSpace(v.(string))
					default:
						return ""
					}
				},
			},
			"finger_print": &schema.Schema{
				Type:     schema.TypeString,
				Computed: true,
			},
			"key_file": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
				Default:  "alicloud_keypair.pem",
			},
		},
	}
}

func resourceAlicloudKeyPairCreate(d *schema.ResourceData, meta interface{}) error {
	conn := meta.(*AliyunClient).ecsconn

	var keyName string
	if v, ok := d.GetOk("key_name"); ok {
		keyName = v.(string)
	} else if v, ok := d.GetOk("key_name_prefix"); ok {
		keyName = resource.PrefixedUniqueId(v.(string))
	} else {
		keyName = resource.UniqueId()
	}

	if publicKey, ok := d.GetOk("public_key"); ok {
		keypair, err := conn.ImportKeyPair(&ecs.ImportKeyPairArgs{
			RegionId:      getRegion(d, meta),
			KeyPairName:   keyName,
			PublicKeyBody: publicKey.(string),
		})
		if err != nil {
			return fmt.Errorf("Error Import KeyPair: %s", err)
		}

		d.SetId(keypair.KeyPairName)
	} else {
		keypair, err := conn.CreateKeyPair(&ecs.CreateKeyPairArgs{
			RegionId:    getRegion(d, meta),
			KeyPairName: keyName,
		})
		if err != nil {
			return fmt.Errorf("Error Create KeyPair: %s", err)
		}

		d.SetId(keypair.KeyPairName)
		ioutil.WriteFile(d.Get("key_file").(string), []byte(keypair.PrivateKeyBody), 444)
	}

	return resourceAlicloudKeyPairRead(d, meta)
}

func resourceAlicloudKeyPairRead(d *schema.ResourceData, meta interface{}) error {
	conn := meta.(*AliyunClient).ecsconn

	keypairs, _, err := conn.DescribeKeyPairs(&ecs.DescribeKeyPairsArgs{
		RegionId:    getRegion(d, meta),
		KeyPairName: d.Id(),
	})
	if err != nil {
		if IsExceptedError(err, KeyPairNotFound) {
			d.SetId("")
			return nil
		}
		return fmt.Errorf("Error Retrieving KeyPair: %s", err)
	}

	if len(keypairs) > 0 {
		d.Set("key_name", keypairs[0].KeyPairName)
		d.Set("fingerprint", keypairs[0].KeyPairFingerPrint)
		return nil
	}

	return fmt.Errorf("Unable to find key pair within: %#v", keypairs)
}

func resourceAlicloudKeyPairDelete(d *schema.ResourceData, meta interface{}) error {
	conn := meta.(*AliyunClient).ecsconn

	err := conn.DeleteKeyPairs(&ecs.DeleteKeyPairsArgs{
		RegionId:     getRegion(d, meta),
		KeyPairNames: convertListToJsonString(append(make([]interface{}, 0, 1), d.Id())),
	})
	return err
}
