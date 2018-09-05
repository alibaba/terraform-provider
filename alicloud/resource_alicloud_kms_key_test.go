package alicloud

import (
	"fmt"
	"testing"

	"github.com/denverdino/aliyungo/kms"
	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/terraform"
	"strings"
)


func init() {
	resource.AddTestSweepers("alicloud_kms_key", &resource.Sweeper{
		Name: "alicloud_kms_key",
		F:    testSweepKmsKeys,
	})
}

func testSweepKeyPairs(region string) error {
	client, err := sharedClientForRegion(region)
	if err != nil {
		return fmt.Errorf("error getting Alicloud client: %s", err)
	}
	conn := client.(*AliyunClient)

	prefixes := []string{
		"tf-testAcc",
		"tf_testAcc",
		"tf_test_",
		"tf-test-",
		"testAcc",
		"terraform-test-",
	}

	var pairs []ecs.KeyPair
	req := ecs.CreateDescribeKeyPairsRequest()
	req.RegionId = conn.RegionId
	req.PageSize = requests.NewInteger(PageSizeLarge)
	req.PageNumber = requests.NewInteger(1)
	for {
		resp, err := conn.ecsconn.DescribeKeyPairs(req)
		if err != nil {
			return fmt.Errorf("Error retrieving Key Pairs: %s", err)
		}
		if resp == nil || len(resp.KeyPairs.KeyPair) < 1 {
			break
		}
		pairs = append(pairs, resp.KeyPairs.KeyPair...)

		if len(resp.KeyPairs.KeyPair) < PageSizeLarge {
			break
		}

		if page, err := getNextpageNumber(req.PageNumber); err != nil {
			return err
		} else {
			req.PageNumber = page
		}
	}

	for _, v := range pairs {
		name := v.KeyPairName
		skip := true
		for _, prefix := range prefixes {
			if strings.HasPrefix(strings.ToLower(name), strings.ToLower(prefix)) {
				skip = false
				break
			}
		}
		if skip {
			log.Printf("[INFO] Skipping Key Pair: %s", name)
			continue
		}
		log.Printf("[INFO] Deleting Key Pair: %s", name)
		req := ecs.CreateDeleteKeyPairsRequest()
		req.KeyPairNames = convertListToJsonString(append(make([]interface{}, 0, 1), name))
		if _, err := conn.ecsconn.DeleteKeyPairs(req); err != nil {
			log.Printf("[ERROR] Failed to delete Key Pair (%s): %s", name, err)
		}
	}
	return nil
}

func TestAccAlicloudKmsKey_basic(t *testing.T) {
	var keyBefore, keyAfter kms.KeyMetadata

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckAlicloudKmsKeyDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAlicloudKmsKeyBasic,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAlicloudKmsKeyExists("alicloud_kms_key.key", &keyBefore),
					resource.TestCheckResourceAttr("alicloud_kms_key.key", "is_enabled", "true"),
				),
			},
			{
				Config: testAlicloudKmsKeyUpdate,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAlicloudKmsKeyExists("alicloud_kms_key.key", &keyAfter),
					resource.TestCheckResourceAttr("alicloud_kms_key.key", "is_enabled", "false"),
				),
			},
		},
	})
}

func testAccCheckAlicloudKmsKeyExists(name string, key *kms.KeyMetadata) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[name]
		if !ok {
			return fmt.Errorf("Not found: %s", name)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("No KMS Key ID is set")
		}

		conn := testAccProvider.Meta().(*AliyunClient).kmsconn

		o, err := conn.DescribeKey(rs.Primary.ID)
		if err != nil {
			return err
		}

		meta := o.KeyMetadata
		key = &meta

		return nil
	}
}

func testAccCheckAlicloudKmsKeyDestroy(s *terraform.State) error {
	conn := testAccProvider.Meta().(*AliyunClient).kmsconn

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "alicloud_kms_key" {
			continue
		}

		out, err := conn.DescribeKey(rs.Primary.ID)

		if err != nil && !IsExceptedError(err, ForbiddenKeyNotFound) {
			return err
		}

		if KeyState(out.KeyMetadata.KeyState) == PendingDeletion {
			return nil
		}

		return fmt.Errorf("KMS key still exists:\n%#v", out.KeyMetadata)
	}

	return nil
}

const testAlicloudKmsKeyBasic = `
resource "alicloud_kms_key" "key" {
    description = "Terraform acc test"
    deletion_window_in_days = 7
}`

const testAlicloudKmsKeyUpdate = `
resource "alicloud_kms_key" "key" {
    description = "Terraform acc test"
    deletion_window_in_days = 7
    is_enabled = false
}`
