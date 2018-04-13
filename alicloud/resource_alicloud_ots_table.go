package alicloud

import (
	"fmt"
	"log"
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
			"primary_key": {
				Type:     schema.TypeList,
				Optional: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"name": &schema.Schema{
							Type:     schema.TypeString,
							Required: true,
						},
						"type": &schema.Schema{
							Type:     schema.TypeString,
							Required: true,
						},
					},
				},
				MinItems: 1,
				MaxItems: 4,
			},
			"table_option": {
				Type:     schema.TypeSet,
				Required: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"time_to_live": {
							Type:     schema.TypeInt,
							Optional: true,
						},
						"max_version": {
							Type:     schema.TypeInt,
							Optional: true,
						},
					},
				},
			},
		},
	}
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
	client := meta.(*AliyunClient).otsconn

	tableMeta := new(tablestore.TableMeta)
	table_name := d.Get("table_name").(string)
	tableMeta.TableName = table_name

	for _, primaryKey :=range d.Get("primary_key").([]interface{}){
		pk := primaryKey.(map[string]interface{})
		pkValue := getPrimaryKeyType(pk["type"].(string))
		tableMeta.AddPrimaryKeyColumn(pk["name"].(string), pkValue)

	}
	tableOption := new(tablestore.TableOption)
	to := d.Get("table_option").(*schema.Set).List()
	w := to[0].(map[string]interface{})

	tableOption.TimeToAlive = w["time_to_live"].(int)
	tableOption.MaxVersion = w["max_version"].(int)

	reservedThroughput := new(tablestore.ReservedThroughput)

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
	return resourceAliyunOtsTableUpdate(d, meta)
}

func resourceAliyunOtsTableRead(d *schema.ResourceData, meta interface{}) error {
	log.Println("Describing OTS Table...")
	client := meta.(*AliyunClient).otsconn

	describeTableReq := new(tablestore.DescribeTableRequest)
	describeTableReq.TableName = d.Get("table_name").(string)

	describ, err := client.DescribeTable(describeTableReq)

	if err != nil {
		return fmt.Errorf("failed to update table with error: %s", err)
	} else {
		d.Set("table_name", describ.TableMeta.TableName)
		t := make(map[string]interface{})
		t["time_to_live"] = describ.TableOption.TimeToAlive
		t["max_version"] = describ.TableOption.MaxVersion
		var tableOptions []map[string]interface{}
		tableOptions = append(tableOptions, t)

		d.Set("table_option", tableOptions)

		log.Println("DescribeTable finished. Table meta: %d, %d", describ.TableOption.MaxVersion, describ.TableOption.TimeToAlive)
	}

	return nil
}

func resourceAliyunOtsTableUpdate(d *schema.ResourceData, meta interface{}) error {
	log.Println("Updating OTS Table...")
	client := meta.(*AliyunClient).otsconn

	update := false

	updateTableReq := new(tablestore.UpdateTableRequest)
	table_name := d.Get("table_name").(string)
	updateTableReq.TableName = table_name
	//updateTableReq.TableOption = new(tablestore.TableOption)

	if d.HasChange("table_option") && !d.IsNewResource(){
		tableOption := new(tablestore.TableOption)
		update = true
		to := d.Get("table_option").(*schema.Set).List()
		w := to[0].(map[string]interface{})

		tableOption.TimeToAlive = w["time_to_live"].(int)
		tableOption.MaxVersion = w["max_version"].(int)
		updateTableReq.TableOption = tableOption
	}

	if update {
		_, err := client.UpdateTable(updateTableReq)

		if err != nil {
			return fmt.Errorf("failed to update table with error: %s", err)
		} else {
			log.Println("update finished")
			return resourceAliyunOtsTableRead(d, meta)
		}
	}
	return nil
}

func resourceAliyunOtsTableDelete(d *schema.ResourceData, meta interface{}) error {
	log.Println("Deleting OTS table...")
	client := meta.(*AliyunClient).otsconn

	deleteReq := new(tablestore.DeleteTableRequest)
	deleteReq.TableName = d.Get("table_name").(string)
	_, err := client.DeleteTable(deleteReq)
	if err != nil {
		return fmt.Errorf("Failed to delete table with error: %s", err)
	} else {
		log.Println("Delete table finished")
	}

	return nil
}