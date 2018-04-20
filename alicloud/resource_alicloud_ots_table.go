package alicloud

import (
	"fmt"
	"github.com/aliyun/aliyun-tablestore-go-sdk/tablestore"
	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/helper/schema"
	"time"
)

func resourceAlicloudOtsTable() *schema.Resource {
	return &schema.Resource{
		Create: resourceAliyunOtsTableCreate,
		Read:   resourceAliyunOtsTableRead,
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
				Required: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"name": &schema.Schema{
							Type:     schema.TypeString,
							Required: true,
						},
						"type": &schema.Schema{
							Type:         schema.TypeString,
							Required:     true,
							ValidateFunc: validatePrimaryTypeKey,
						},
					},
				},
				MaxItems: 4,
			},
			"time_to_live": &schema.Schema{
				Type:     schema.TypeInt,
				Required: true,
			},
			"max_version": &schema.Schema{
				Type:     schema.TypeInt,
				Required: true,
			},
		},
	}
}

func resourceAliyunOtsTableCreate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*AliyunClient).otsconn

	tableMeta := new(tablestore.TableMeta)
	table_name := d.Get("table_name").(string)
	tableMeta.TableName = table_name

	for _, primaryKey := range d.Get("primary_key").([]interface{}) {
		pk := primaryKey.(map[string]interface{})
		pkValue := getPrimaryKeyType(pk["type"].(string))
		tableMeta.AddPrimaryKeyColumn(pk["name"].(string), pkValue)
	}
	tableOption := new(tablestore.TableOption)
	tableOption.TimeToAlive = d.Get("time_to_live").(int)
	tableOption.MaxVersion = d.Get("max_version").(int)

	reservedThroughput := new(tablestore.ReservedThroughput)

	createTableRequest := new(tablestore.CreateTableRequest)
	createTableRequest.TableMeta = tableMeta
	createTableRequest.TableOption = tableOption
	createTableRequest.ReservedThroughput = reservedThroughput

	_, err := client.CreateTable(createTableRequest)
	if err != nil {
		return fmt.Errorf("Failed to create table with error: %s", err)
	}

	// Need to set id before calling read method or terraform.state won't be generated.
	d.SetId(table_name)
	return resourceAliyunOtsTableUpdate(d, meta)
}

func resourceAliyunOtsTableRead(d *schema.ResourceData, meta interface{}) error {
	tableName := d.Get("table_name").(string)
	describ, err := describeOtsTable(tableName, meta)

	if err != nil {
		return fmt.Errorf("failed to describe table with error: %s", err)
	}

	d.Set("table_name", describ.TableMeta.TableName)

	var pks []map[string]interface{}
	keys := describ.TableMeta.SchemaEntry
	for k, v := range keys {
		item := make(map[string]interface{})
		item["name"] = k
		item["type"] = v
		pks = append(pks, item)
	}
	d.Set("primary_key", pks)

	d.Set("time_to_live", describ.TableOption.TimeToAlive)
	d.Set("max_version", describ.TableOption.MaxVersion)

	return nil
}

func resourceAliyunOtsTableUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*AliyunClient).otsconn
	update := false

	updateTableReq := new(tablestore.UpdateTableRequest)
	table_name := d.Get("table_name").(string)
	updateTableReq.TableName = table_name

	// As the issue of ots sdk, time_to_live and max_version need to be updated together at present.
	// For the issue, please refer to https://github.com/aliyun/aliyun-tablestore-go-sdk/issues/18
	tableOption := new(tablestore.TableOption)
	if d.HasChange("time_to_live") && !d.IsNewResource() {
		update = true
		tableOption.TimeToAlive = d.Get("time_to_live").(int)
	}

	if d.HasChange("max_version") && !d.IsNewResource() {
		update = true
		tableOption.MaxVersion = d.Get("max_version").(int)
	}

	if update {
		updateTableReq.TableOption = tableOption
		_, err := client.UpdateTable(updateTableReq)

		if err != nil {
			return fmt.Errorf("failed to update table with error: %s", err)
		}
	}
	return resourceAliyunOtsTableRead(d, meta)
}

func resourceAliyunOtsTableDelete(d *schema.ResourceData, meta interface{}) error {
	tableName := d.Get("table_name").(string)
	return resource.Retry(2*time.Minute, func() *resource.RetryError {
		successFlag, err := deleteOtsTable(tableName, meta)
		if !successFlag {
			return resource.RetryableError(fmt.Errorf("Delete instance timeout and got an error: %#v.", err))
		}
		return nil
	})
}
