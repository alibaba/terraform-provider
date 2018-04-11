package alicloud

import (
	"fmt"
	"log"
	"os"

	"github.com/hashicorp/terraform/helper/schema"
	"github.com/aliyun/aliyun-tablestore-go-sdk/tablestore"
)

func resourceAlicloudOtsTable() *schema.Resource {
	return &schema.Resource{
		Create: resourceAliyunOtsTableCreate,
		Read: resourceAliyunOtsTableRead,
		Update: resourceAliyunOtsTableUpdate,
		Delete: resourceAliyunOtsTableDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},

		Schema: map[string]*schema.Schema{
			"table_name": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"primary_key_1_name": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
				//ValidateFunc: validateDiskName,

			},
			"primary_key_1_type": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
				//ValidateFunc: validateDiskName,
				//Default: tablestore.PrimaryKeyType_INTEGER,
			},
			"primary_key_2_name": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
			},
			"primary_key_2_type": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
			},
			"primary_key_3_name": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
			},
			"primary_key_3_type": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
			},
			"primary_key_4_name": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
			},
			"primary_key_4_type": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
			},
			"time_to_live": &schema.Schema{
				Type:     schema.TypeInt,
				Optional: true,
				//ValidateFunc: validateDiskName,
			},
			"max_version": &schema.Schema{
				Type:     schema.TypeInt,
				Optional: true,
			},
			"read_cap": &schema.Schema{  //read capacity
				Type:     schema.TypeInt,
				Optional: true,
			},
			"write_cap": &schema.Schema{  //write capacity
				Type:     schema.TypeInt,
				Optional: true,
			},

		},
	}
}

func getOtsClient() *tablestore.TableStoreClient{
	/*
	Get OTS Client
	 */
	endpoint := os.Getenv("OTS_ENDPOINT")
	instanceName := os.Getenv("OTS_INSTANCE_NAME")
	accessKeyId := os.Getenv("ALICLOUD_ACCESS_KEY")
	accessKeySecret := os.Getenv("ALICLOUD_SECRET_KEY")
	client := tablestore.NewClient(endpoint, instanceName, accessKeyId, accessKeySecret)
	return client
}

func getPrimaryKeyType(primaryKeyType string) tablestore.PrimaryKeyType{
	var keyType tablestore.PrimaryKeyType
	switch primaryKeyType{
	case "INTEGER":
		keyType = tablestore.PrimaryKeyType_INTEGER
	case "STRING":
		keyType = tablestore.PrimaryKeyType_STRING
	case "BINARY":
		keyType = tablestore.PrimaryKeyType_BINARY
	}
	return keyType
}

func resourceAliyunOtsTableCreate(d *schema.ResourceData, meta interface{}) error {
	log.Printf("Creating OTS Table...")
	client := getOtsClient()

	tableMeta := new(tablestore.TableMeta)

	table_name := d.Get("table_name").(string)
	tableMeta.TableName = table_name

	primaryKeyType := d.Get("primary_key_1_type").(string)

	keyType := getPrimaryKeyType(primaryKeyType)

	tableMeta.AddPrimaryKeyColumn(d.Get("primary_key_1_name").(string), keyType)

	primaryKey2Name := d.Get("primary_key_2_name")
	primaryKey2Type := getPrimaryKeyType(d.Get("primary_key_2_type").(string))
	//Whether need to check the value of primaryKey2Type?
	if primaryKey2Name != nil {
		tableMeta.AddPrimaryKeyColumn(primaryKey2Name.(string), primaryKey2Type)
	}

	primaryKey3Name := d.Get("primary_key_3_name")
	primaryKey3Type := getPrimaryKeyType(d.Get("primary_key_3_type").(string))
	if primaryKey3Name != nil {
		tableMeta.AddPrimaryKeyColumn(primaryKey3Name.(string), primaryKey3Type)
	}

	primaryKey4Name := d.Get("primary_key_4_name")
	primaryKey4Type := getPrimaryKeyType(d.Get("primary_key_4_type").(string))
	if primaryKey4Name != nil {
		tableMeta.AddPrimaryKeyColumn(primaryKey4Name.(string), primaryKey4Type)
	}

	tableOption := new(tablestore.TableOption)
	tableOption.TimeToAlive = d.Get("time_to_live").(int)
	tableOption.MaxVersion = d.Get("max_version").(int)
	reservedThroughput := new(tablestore.ReservedThroughput)
	//reservedThroughput.Readcap = d.Get("read_cap").(int)
	//reservedThroughput.Writecap = d.Get("write_cap").(int)

	createTableRequest := new(tablestore.CreateTableRequest)
	createTableRequest.TableMeta = tableMeta
	createTableRequest.TableOption = tableOption
	createTableRequest.ReservedThroughput = reservedThroughput
	_, err := client.CreateTable(createTableRequest)
	if err != nil {
		return fmt.Errorf("Failed to create table with error: %s", err)
	} else {
		log.Printf("Create table finished")
	}

	// Need to set id before calling read method or terraform.state won't be generated.
	d.SetId(table_name)
	//return resourceAliyunOtsTableUpdate(d, meta)
	return resourceAliyunOtsTableRead(d, meta)
}

func resourceAliyunOtsTableRead(d *schema.ResourceData, meta interface{}) error {
	log.Println("Describing OTS Table...")
	client := getOtsClient()

	describeTableReq := new(tablestore.DescribeTableRequest)
	describeTableReq.TableName = d.Get("table_name").(string)

	describ, err := client.DescribeTable(describeTableReq)

	if err != nil {
		fmt.Println("failed to update table with error: %s", err)
	} else {
		d.Set("table_name", describ.TableMeta.TableName)
		d.Set("time_to_live", describ.TableOption.TimeToAlive)
		d.Set("max_version", describ.TableOption.MaxVersion)
		d.Set("read_cap", describ.ReservedThroughput.Readcap)
		d.Set("write_cap", describ.ReservedThroughput.Writecap)
		log.Println("DescribeTable finished. Table meta: %d, %d", describ.TableOption.MaxVersion, describ.TableOption.TimeToAlive)
	}

	return nil
}

func resourceAliyunOtsTableUpdate(d *schema.ResourceData, meta interface{}) error {
	log.Println("Updating OTS Table...")
	client := getOtsClient()

	update := false

	updateTableReq := new(tablestore.UpdateTableRequest)
	table_name := d.Get("table_name").(string)
	updateTableReq.TableName = table_name
	updateTableReq.TableOption = new(tablestore.TableOption)

	if d.HasChange("time_to_live") {
		update = true
		time_to_live := d.Get("time_to_live").(int)
		log.Printf("time_to_live changed to: %d", time_to_live)
		updateTableReq.TableOption.TimeToAlive = time_to_live
	}

	if d.HasChange("max_version"){
		update = true
		max_version := d.Get("max_version").(int)
		log.Printf("time_to_live changed to: %d", max_version)
		updateTableReq.TableOption.MaxVersion = max_version
	}

	if update {
		_, err := client.UpdateTable(updateTableReq)

		if err != nil {
			log.Println("failed to update table with error: %s", err)
		} else {
			log.Println("update finished")
			d.SetId(table_name)
			return resourceAliyunOtsTableRead(d, meta)
		}
	}
	return nil
}

func resourceAliyunOtsTableDelete(d *schema.ResourceData, meta interface{}) error {
	log.Println("Deleting OTS table...")
	client := getOtsClient()

	deleteReq := new(tablestore.DeleteTableRequest)
	deleteReq.TableName = d.Get("table_name").(string)
	_, err := client.DeleteTable(deleteReq)
	if err != nil {
		log.Println("Failed to delete table with error: %s", err)
	} else {
		log.Println("Delete table finished")
	}

	return nil
}