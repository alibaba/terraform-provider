package alicloud

import (
	"fmt"
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
				Required: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"name": &schema.Schema{
							Type:     schema.TypeString,
							Required: true,
						},
						"type": &schema.Schema{
							Type:     schema.TypeString,
							Required: true,
							ValidateFunc: validatePrimaryTypeKey,
						},
					},
				},
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
							Default: -1,
						},
						"max_version": {
							Type:     schema.TypeInt,
							Optional: true,
							Default: 1,
						},
					},
				},
				MaxItems: 1,
			},
		},
	}
}

func resourceAliyunOtsTableCreate(d *schema.ResourceData, meta interface{}) error {
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
	} else {
		d.Set("table_name", describ.TableMeta.TableName)

		t := make(map[string]interface{})
		t["time_to_live"] = describ.TableOption.TimeToAlive
		t["max_version"] = describ.TableOption.MaxVersion
		var tableOptions []map[string]interface{}
		tableOptions = append(tableOptions, t)
		d.Set("table_option", tableOptions)

		var pks []map[string]interface{}
		keys := describ.TableMeta.SchemaEntry
		for k, v := range keys{
			item := make(map[string]interface{})
			item["name"] = k
			item["type"] = v
			pks = append(pks, item)
		}
		d.Set("primary_key", pks)
	}

	return nil
}

func resourceAliyunOtsTableUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*AliyunClient).otsconn
	d.Partial(true)
	update := false

	updateTableReq := new(tablestore.UpdateTableRequest)
	table_name := d.Get("table_name").(string)
	updateTableReq.TableName = table_name

	if d.HasChange("table_option") && !d.IsNewResource(){
		tableOption := new(tablestore.TableOption)
		update = true
		to := d.Get("table_option").(*schema.Set).List()
		w := to[0].(map[string]interface{})

		tableOption.TimeToAlive = w["time_to_live"].(int)
		tableOption.MaxVersion = w["max_version"].(int)
		updateTableReq.TableOption = tableOption
		d.SetPartial("table_option")
	}

	if update {
		_, err := client.UpdateTable(updateTableReq)

		if err != nil {
			return fmt.Errorf("failed to update table with error: %s", err)
		}
	}
	d.Partial(false)
	return resourceAliyunOtsTableRead(d, meta)
}

func resourceAliyunOtsTableDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*AliyunClient).otsconn

	deleteReq := new(tablestore.DeleteTableRequest)
	tableName := d.Get("table_name").(string)
	deleteReq.TableName = tableName
	_, err := client.DeleteTable(deleteReq)

	describ, _ := describeOtsTable(tableName, meta)

	if err != nil || describ.TableMeta != nil{
		return fmt.Errorf("Failed to delete table with error: %s", err)
	}

	return nil
}