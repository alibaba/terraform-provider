package alicloud

import (
	"github.com/denverdino/aliyungo/common"
	"github.com/denverdino/aliyungo/rds"
	"log"
)

// when getInstance is empty, then throw InstanceNotfound error
func (client *AliyunClient) DescribeDBInstanceById(id string) (instance *rds.DBInstanceAttribute, err error) {
	arrtArgs := rds.DescribeDBInstancesArgs{
		DBInstanceId: id,
	}
	resp, err := client.rdsconn.DescribeDBInstanceAttribute(&arrtArgs)
	if err != nil {
		return nil, err
	}

	attr := resp.Items.DBInstanceAttribute

	if len(attr) <= 0 {
		return nil, common.GetClientErrorFromString(InstanceNotfound)
	}

	return &attr[0], nil
}

func (client *AliyunClient) CreateAccountByInfo(instanceId, username, pwd string) error {
	conn := client.rdsconn
	args := rds.CreateAccountArgs{
		DBInstanceId:    instanceId,
		AccountName:     username,
		AccountPassword: pwd,
	}

	if _, err := conn.CreateAccount(&args); err != nil {
		return err
	}

	if err := conn.WaitForAccount(instanceId, username, rds.Available, 200); err != nil {
		return err
	}
	return nil
}

func (client *AliyunClient) CreateDatabaseByInfo(instanceId, dbName, charset, desp string) error {
	conn := client.rdsconn
	args := rds.CreateDatabaseArgs{
		DBInstanceId:     instanceId,
		DBName:           dbName,
		CharacterSetName: charset,
		DBDescription:    desp,
	}
	_, err := conn.CreateDatabase(&args)
	return err
}

func (client *AliyunClient) GrantDBPrivilege2Account(instanceId, username, dbName string) error {
	conn := client.rdsconn
	pargs := rds.GrantAccountPrivilegeArgs{
		DBInstanceId:     instanceId,
		AccountName:      username,
		DBName:           dbName,
		AccountPrivilege: rds.ReadWrite,
	}
	if _, err := conn.GrantAccountPrivilege(&pargs); err != nil {
		return err
	}

	if err := conn.WaitForAccountPrivilege(instanceId, username, dbName, rds.ReadWrite, 200); err != nil {
		return err
	}
	return nil
}

func (client *AliyunClient) AllocateDBPublicConnection(instanceId, port string) error {
	conn := client.rdsconn
	args := rds.AllocateInstancePublicConnectionArgs{
		DBInstanceId:           instanceId,
		ConnectionStringPrefix: instanceId + "o",
		Port: port,
	}

	if _, err := conn.AllocateInstancePublicConnection(&args); err != nil {
		return err
	}

	if err := conn.WaitForPublicConnection(instanceId, 600); err != nil {
		return err
	}
	return nil
}

func (client *AliyunClient) ConfigDBBackup(instanceId, backupTime, backupPeriod string, retentionPeriod int) error {
	bargs := rds.BackupPolicy{
		PreferredBackupTime:   backupTime,
		PreferredBackupPeriod: backupPeriod,
		BackupRetentionPeriod: retentionPeriod,
	}
	args := rds.ModifyBackupPolicyArgs{
		DBInstanceId: instanceId,
		BackupPolicy: bargs,
	}

	if _, err := client.rdsconn.ModifyBackupPolicy(&args); err != nil {
		return err
	}
	return nil
}

func (client *AliyunClient) ModifySecurityIps(instanceId, ips string) error {
	sargs := rds.DBInstanceIPArray{
		SecurityIps: ips,
	}

	args := rds.ModifySecurityIpsArgs{
		DBInstanceId:      instanceId,
		DBInstanceIPArray: sargs,
	}

	if _, err := client.rdsconn.ModifySecurityIps(&args); err != nil {
		return err
	}
	return nil
}

// turn period to TimeType
func TransformPeriod2Time(period int, chargeType string) (ut int, tt common.TimeType) {
	log.Printf("get period %d chargeType %s", period, chargeType)
	if chargeType == string(rds.Postpaid) {
		return 1, common.Day
	}

	if period >= 1 && period <= 9 {
		return period, common.Month
	}

	if period == 12 {
		return 1, common.Year
	}

	if period == 24 {
		return 2, common.Year
	}
	return 0, common.Day

}

// turn TimeType to Period
func TransformTime2Period(ut int, tt common.TimeType) (period int) {
	if tt == common.Year {
		return 12 * ut
	}

	return ut

}

// Flattens an array of databases into a []map[string]interface{}
func flattenDatabaseMappings(list []rds.Database) []map[string]interface{} {
	result := make([]map[string]interface{}, 0, len(list))
	for _, i := range list {
		l := map[string]interface{}{
			"db_name":            i.DBName,
			"character_set_name": i.CharacterSetName,
			"db_description":     i.DBDescription,
		}
		result = append(result, l)
	}
	return result
}

func flattenDBBackup(list []rds.BackupPolicy) []map[string]interface{} {
	result := make([]map[string]interface{}, 0, len(list))
	for _, i := range list {
		l := map[string]interface{}{
			"preferred_backup_period": i.PreferredBackupPeriod,
			"preferred_backup_time":   i.PreferredBackupTime,
			"backup_retention_period": i.LogBackupRetentionPeriod,
		}
		result = append(result, l)
	}
	return result
}

func flattenDBSecurityIPs(list []rds.DBInstanceIPArray) []map[string]interface{} {
	result := make([]map[string]interface{}, 0, len(list))
	for _, i := range list {
		l := map[string]interface{}{
			"security_ips": i.SecurityIps,
		}
		result = append(result, l)
	}
	return result
}

// Flattens an array of databases connection into a []map[string]interface{}
func flattenDBConnections(list []rds.DBInstanceNetInfo) []map[string]interface{} {
	result := make([]map[string]interface{}, 0, len(list))
	for _, i := range list {
		l := map[string]interface{}{
			"connection_string": i.ConnectionString,
			"ip_type":           i.IPType,
			"ip_address":        i.IPAddress,
		}
		result = append(result, l)
	}
	return result
}
