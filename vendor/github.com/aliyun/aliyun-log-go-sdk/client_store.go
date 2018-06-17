package sls

func convertLogstore(c *Client, project, logstore string) *LogStore {
	c.accessKeyLock.RLock()
	proj := &LogProject{
		Name:            project,
		Endpoint:        c.Endpoint,
		AccessKeyID:     c.AccessKeyID,
		AccessKeySecret: c.AccessKeySecret,
		SecurityToken:   c.SecurityToken,
	}
	c.accessKeyLock.RUnlock()
	return &LogStore{
		project: proj,
		Name:    logstore,
	}
}

// ListShards returns shard id list of this logstore.
func (c *Client) ListShards(project, logstore string) (shardIDs []int, err error) {
	ls := convertLogstore(c, project, logstore)
	return ls.ListShards()
}

// PutLogs put logs into logstore.
// The callers should transform user logs into LogGroup.
func (c *Client) PutLogs(project, logstore string, lg *LogGroup) (err error) {
	ls := convertLogstore(c, project, logstore)
	return ls.PutLogs(lg)
}

// PutLogsWithCompressType put logs into logstore with specific compress type.
// The callers should transform user logs into LogGroup.
func (c *Client) PutLogsWithCompressType(project, logstore string, lg *LogGroup, compressType int) (err error) {
	ls := convertLogstore(c, project, logstore)
	if err := ls.SetPutLogCompressType(compressType); err != nil {
		return err
	}
	return ls.PutLogs(lg)
}

// PutRawLogWithCompressType put raw log data to log service, no marshal
func (c *Client) PutRawLogWithCompressType(project, logstore string, rawLogData []byte, compressType int) (err error) {
	ls := convertLogstore(c, project, logstore)
	if err := ls.SetPutLogCompressType(compressType); err != nil {
		return err
	}
	return ls.PutRawLog(rawLogData)
}

// GetCursor gets log cursor of one shard specified by shardId.
// The from can be in three form: a) unix timestamp in seccond, b) "begin", c) "end".
// For more detail please read: https://help.aliyun.com/document_detail/29024.html
func (c *Client) GetCursor(project, logstore string, shardID int, from string) (cursor string, err error) {
	ls := convertLogstore(c, project, logstore)
	return ls.GetCursor(shardID, from)
}

// GetLogsBytes gets logs binary data from shard specified by shardId according cursor and endCursor.
// The logGroupMaxCount is the max number of logGroup could be returned.
// The nextCursor is the next curosr can be used to read logs at next time.
func (c *Client) GetLogsBytes(project, logstore string, shardID int, cursor, endCursor string,
	logGroupMaxCount int) (out []byte, nextCursor string, err error) {
	ls := convertLogstore(c, project, logstore)
	return ls.GetLogsBytes(shardID, cursor, endCursor, logGroupMaxCount)
}

// PullLogs gets logs from shard specified by shardId according cursor and endCursor.
// The logGroupMaxCount is the max number of logGroup could be returned.
// The nextCursor is the next cursor can be used to read logs at next time.
// @note if you want to pull logs continuous, set endCursor = ""
func (c *Client) PullLogs(project, logstore string, shardID int, cursor, endCursor string,
	logGroupMaxCount int) (gl *LogGroupList, nextCursor string, err error) {
	ls := convertLogstore(c, project, logstore)
	return ls.PullLogs(shardID, cursor, endCursor, logGroupMaxCount)
}

// GetHistograms query logs with [from, to) time range
func (c *Client) GetHistograms(project, logstore string, topic string, from int64, to int64, queryExp string) (*GetHistogramsResponse, error) {
	ls := convertLogstore(c, project, logstore)
	return ls.GetHistograms(topic, from, to, queryExp)
}

// GetLogs query logs with [from, to) time range
func (c *Client) GetLogs(project, logstore string, topic string, from int64, to int64, queryExp string,
	maxLineNum int64, offset int64, reverse bool) (*GetLogsResponse, error) {
	ls := convertLogstore(c, project, logstore)
	return ls.GetLogs(topic, from, to, queryExp, maxLineNum, offset, reverse)
}

// CreateIndex ...
func (c *Client) CreateIndex(project, logstore string, index Index) error {
	ls := convertLogstore(c, project, logstore)
	return ls.CreateIndex(index)
}

// UpdateIndex ...
func (c *Client) UpdateIndex(project, logstore string, index Index) error {
	ls := convertLogstore(c, project, logstore)
	return ls.UpdateIndex(index)
}

// DeleteIndex ...
func (c *Client) DeleteIndex(project, logstore string) error {
	ls := convertLogstore(c, project, logstore)
	return ls.DeleteIndex()
}

// GetIndex ...
func (c *Client) GetIndex(project, logstore string) (*Index, error) {
	ls := convertLogstore(c, project, logstore)
	return ls.GetIndex()
}
