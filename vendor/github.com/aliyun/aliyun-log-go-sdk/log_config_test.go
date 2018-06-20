package sls

import (
	"encoding/json"
	"fmt"
	"os"
	"testing"

	"github.com/golang/glog"
	"github.com/stretchr/testify/suite"
)

func TestConfig(t *testing.T) {
	suite.Run(t, new(ConfigTestSuite))
	glog.Flush()
}

type ConfigTestSuite struct {
	suite.Suite
	endpoint         string
	projectName      string
	logstoreName     string
	accessKeyID      string
	accessKeySecret  string
	client           Client
	machineGroupName string
}

func makeSureLogstoreExist(c *Client, project, logstore string) error {
	if ok, err := c.CheckProjectExist(project); err != nil {
		return err
	} else if !ok {
		_, err := c.CreateProject(project, "go sdk test")
		if err != nil {
			return err
		}
	}
	if ok, err := c.CheckLogstoreExist(project, logstore); err != nil {
		return err
	} else if !ok {
		err := c.CreateLogStore(project, logstore, 1, 2)
		if err != nil {
			return err
		}
	}
	return nil
}

func (s *ConfigTestSuite) SetupSuite() {
	s.endpoint = os.Getenv("LOG_TEST_ENDPOINT")
	s.projectName = os.Getenv("LOG_TEST_PROJECT")
	s.logstoreName = os.Getenv("LOG_TEST_LOGSTORE")
	s.accessKeyID = os.Getenv("LOG_TEST_ACCESS_KEY_ID")
	s.accessKeySecret = os.Getenv("LOG_TEST_ACCESS_KEY_SECRET")
	s.client.AccessKeyID = s.accessKeyID
	s.client.AccessKeySecret = s.accessKeySecret
	s.client.Endpoint = s.endpoint
	machineGroupName := "go-test-machine-group"
	s.Nil(makeSureLogstoreExist(&s.client, s.projectName, s.logstoreName))
	groups, _, err := s.client.ListMachineGroup(s.projectName, 0, 100)
	s.Nil(err)
	for _, name := range groups {
		m, err := s.client.GetMachineGroup(s.projectName, name)
		s.Nil(err)
		fmt.Println(*m)
		if name == machineGroupName {
			s.machineGroupName = machineGroupName
			break
		}
	}
	if len(s.machineGroupName) == 0 {
		m := &MachineGroup{
			Name:          machineGroupName,
			MachineIDType: MachineIDTypeUserDefined,
			MachineIDList: []string{"go-sdk-test-id", "k8s-demo"},
		}
		err := s.client.CreateMachineGroup(s.projectName, m)
		s.Nil(err)
		s.machineGroupName = machineGroupName
	}
	m, err := s.client.GetMachineGroup(s.projectName, s.machineGroupName)
	s.Nil(err)
	fmt.Println("machine group :", *m)
}

func (s *ConfigTestSuite) TearDownSuite() {
	//err := s.client.DeleteMachineGroup(s.projectName, s.machineGroupName)
	//s.Nil(err)
}

func (s *ConfigTestSuite) TestListConfig() {
	configNames, count, err := s.client.ListConfig(s.projectName, 0, 100)
	s.Nil(err)
	s.True(count >= 0, count)
	s.True(len(configNames) >= 0, configNames)
}

func (s *ConfigTestSuite) TestNormalFileConfig() {
	configName := "go-sdk-simple-file-config"
	s.client.DeleteConfig(s.projectName, configName)
	regexConfig := &RegexConfigInputDetail{}
	InitRegexConfigInputDetail(regexConfig)
	config := &LogConfig{
		Name:        configName,
		InputDetail: regexConfig,
		InputType:   InputTypeFile,
		OutputType:  OutputTypeLogService,
		OutputDetail: OutputDetail{
			ProjectName:  s.projectName,
			LogStoreName: s.logstoreName,
		},
	}
	regexConfig.Key = []string{"content"}
	regexConfig.Regex = "(.*)"
	regexConfig.LogBeginRegex = ".*"
	regexConfig.LogPath = "/usr/local/ilogtail"
	regexConfig.FilePattern = "ilogtail.LOG"
	regexConfig.DiscardUnmatch = false
	regexConfig.IsDockerFile = true
	// 采集所有K8S logtail的日志，自循环
	regexConfig.DockerIncludeEnv = map[string]string{
		"ALIYUN_LOGTAIL_USER_DEFINED_ID": "",
	}
	s.Equal(regexConfig.LogType, LogFileTypeRegexLog)
	err := s.client.CreateConfig(s.projectName, config)
	s.Nil(err)
	destConfig, err := s.client.GetConfig(s.projectName, configName)
	s.Nil(err)
	s.Equal(destConfig.Name, configName)
	s.Equal(destConfig.InputType, InputTypeFile)
	s.Equal(destConfig.OutputDetail.ProjectName, s.projectName)
	s.Equal(destConfig.OutputDetail.LogStoreName, s.logstoreName)
	s.Equal(destConfig.OutputType, OutputTypeLogService)
	regexConfigDest, ok := ConvertToRegexConfigInputDetail(destConfig.InputDetail)
	s.True(ok)
	s.Equal(regexConfigDest.Key, regexConfig.Key)
	s.Equal(regexConfigDest.TimeFormat, regexConfig.TimeFormat)
	s.Equal(regexConfigDest.Regex, regexConfig.Regex)
	s.Equal(regexConfigDest.LogBeginRegex, regexConfig.LogBeginRegex)
	s.Equal(regexConfigDest.LogPath, regexConfig.LogPath)
	s.Equal(regexConfigDest.LogType, regexConfig.LogType)
	s.Equal(regexConfigDest.FilePattern, regexConfig.FilePattern)
	s.Nil(s.client.ApplyConfigToMachineGroup(s.projectName, configName, s.machineGroupName))
}

func (s *ConfigTestSuite) TestRegexFileConfig() {
	configName := "go-sdk-regex-file-config"
	s.client.DeleteConfig(s.projectName, configName)
	regexConfig := &RegexConfigInputDetail{}
	InitRegexConfigInputDetail(regexConfig)
	config := &LogConfig{
		Name:        configName,
		InputDetail: regexConfig,
		InputType:   InputTypeFile,
		OutputType:  OutputTypeLogService,
		OutputDetail: OutputDetail{
			ProjectName:  s.projectName,
			LogStoreName: s.logstoreName,
		},
	}
	regexConfig.DiscardUnmatch = false
	regexConfig.Key = []string{"logger", "time", "cluster", "hostname", "sr", "app", "workdir", "exe", "corepath", "signature", "backtrace"}
	regexConfig.Regex = "\\S*\\s+(\\S*)\\s+(\\S*\\s+\\S*)\\s+\\S*\\s+(\\S*)\\s+(\\S*)\\s+(\\S*)\\s+(\\S*)\\s+(\\S*)\\s+(\\S*)\\s+(\\S*)\\s+\\S*\\s+(\\S*)\\s*([^$]+)"
	regexConfig.TimeFormat = "%Y/%m/%d %H:%M:%S"
	regexConfig.LogBeginRegex = `INFO core_dump_info_data .*`
	regexConfig.LogPath = "/cloud/log/tianji/TianjiClient#/core_dump_manager"
	regexConfig.FilePattern = "core_dump_info_data.log*"
	regexConfig.MaxDepth = 0
	s.Equal(regexConfig.LogType, LogFileTypeRegexLog)
	err := s.client.CreateConfig(s.projectName, config)
	s.Nil(err)
	destConfig, err := s.client.GetConfig(s.projectName, configName)
	s.Nil(err)
	s.Equal(destConfig.Name, configName)
	s.Equal(destConfig.InputType, InputTypeFile)
	s.Equal(destConfig.OutputDetail.ProjectName, s.projectName)
	s.Equal(destConfig.OutputDetail.LogStoreName, s.logstoreName)
	s.Equal(destConfig.OutputType, OutputTypeLogService)
	regexConfigDest, ok := ConvertToRegexConfigInputDetail(destConfig.InputDetail)
	s.True(ok)
	s.Equal(regexConfigDest.Key, regexConfig.Key)
	s.Equal(regexConfigDest.TimeFormat, regexConfig.TimeFormat)
	s.Equal(regexConfigDest.Regex, regexConfig.Regex)
	s.Equal(regexConfigDest.LogBeginRegex, regexConfig.LogBeginRegex)
	s.Equal(regexConfigDest.LogPath, regexConfig.LogPath)
	s.Equal(regexConfigDest.LogType, regexConfig.LogType)
	s.Equal(regexConfigDest.FilePattern, regexConfig.FilePattern)
	s.Equal(regexConfigDest.MaxDepth, regexConfigDest.MaxDepth)
	s.Nil(s.client.ApplyConfigToMachineGroup(s.projectName, configName, s.machineGroupName))
}

func (s *ConfigTestSuite) TestJSONFileConfig() {
	configName := "go-sdk-json-config"
	s.client.DeleteConfig(s.projectName, configName)
	jsonConfig := &JSONConfigInputDetail{}
	InitJSONConfigInputDetail(jsonConfig)
	config := &LogConfig{
		Name:        configName,
		InputDetail: jsonConfig,
		InputType:   InputTypeFile,
		OutputType:  OutputTypeLogService,
		OutputDetail: OutputDetail{
			ProjectName:  s.projectName,
			LogStoreName: s.logstoreName,
		},
	}
	jsonConfig.TimeKey = "key_time"
	jsonConfig.TimeFormat = "%Y/%m/%d %H:%M:%S"
	s.Equal(jsonConfig.LogType, LogFileTypeJSONLog)
	err := s.client.CreateConfig(s.projectName, config)
	s.Nil(err)
	destConfig, err := s.client.GetConfig(s.projectName, configName)
	s.Nil(err)
	s.Equal(destConfig.Name, configName)
	s.Equal(destConfig.InputType, InputTypeFile)
	s.Equal(destConfig.OutputDetail.ProjectName, s.projectName)
	s.Equal(destConfig.OutputDetail.LogStoreName, s.logstoreName)
	s.Equal(destConfig.OutputType, OutputTypeLogService)
	jsonConfigDest, ok := ConvertToJSONConfigInputDetail(destConfig.InputDetail)
	s.True(ok)
	s.Equal(jsonConfig.TimeKey, jsonConfigDest.TimeKey)
	s.Equal(jsonConfig.TimeFormat, jsonConfigDest.TimeFormat)
	s.Equal(jsonConfig.LogPath, jsonConfigDest.LogPath)
	s.Equal(jsonConfig.LogType, jsonConfigDest.LogType)
	s.Equal(jsonConfig.FilePattern, jsonConfigDest.FilePattern)

	// update config
	jsonConfig.MaxDepth = 88
	s.Nil(s.client.UpdateConfig(s.projectName, config))
	destConfig, err = s.client.GetConfig(s.projectName, configName)
	s.Nil(err)
	jsonConfigDest, ok = ConvertToJSONConfigInputDetail(destConfig.InputDetail)
	s.True(ok)
	s.Equal(jsonConfigDest.MaxDepth, jsonConfig.MaxDepth)
	s.Nil(s.client.ApplyConfigToMachineGroup(s.projectName, configName, s.machineGroupName))
}

func (s *ConfigTestSuite) TestDelimiterFileConfig() {
	configName := "go-sdk-delimiter-config"
	s.client.DeleteConfig(s.projectName, configName)
	delimiterConfig := &DelimiterConfigInputDetail{}
	InitDelimiterConfigInputDetail(delimiterConfig)
	config := &LogConfig{
		Name:        configName,
		InputDetail: delimiterConfig,
		InputType:   InputTypeFile,
		OutputType:  OutputTypeLogService,
		OutputDetail: OutputDetail{
			ProjectName:  s.projectName,
			LogStoreName: s.logstoreName,
		},
	}
	delimiterConfig.Quote = "\u0001"
	delimiterConfig.Key = []string{"1", "2", "3", "4", "5"}
	delimiterConfig.Separator = "\""
	delimiterConfig.TimeKey = "1"
	delimiterConfig.TimeFormat = "xxxx"
	delimiterConfig.LogPath = "/var/log/log"
	delimiterConfig.FilePattern = "xxxx.log"
	s.Equal(delimiterConfig.LogType, LogFileTypeDelimiterLog)
	err := s.client.CreateConfig(s.projectName, config)
	s.Nil(err)
	destConfig, err := s.client.GetConfig(s.projectName, configName)
	s.Nil(err)
	s.Equal(destConfig.Name, configName)
	s.Equal(destConfig.InputType, InputTypeFile)
	s.Equal(destConfig.OutputDetail.ProjectName, s.projectName)
	s.Equal(destConfig.OutputDetail.LogStoreName, s.logstoreName)
	s.Equal(destConfig.OutputType, OutputTypeLogService)
	delimiterConfigDest, ok := ConvertToDelimiterConfigInputDetail(destConfig.InputDetail)
	s.True(ok)
	s.Equal(delimiterConfigDest.Quote, delimiterConfig.Quote)
	s.Equal(delimiterConfigDest.Separator, delimiterConfig.Separator)
	s.Equal(delimiterConfigDest.Key, delimiterConfig.Key)
	s.Equal(delimiterConfigDest.Quote, delimiterConfig.Quote)
	s.Equal(delimiterConfigDest.TimeKey, delimiterConfig.TimeKey)
	s.Equal(delimiterConfigDest.TimeFormat, delimiterConfig.TimeFormat)
	s.Equal(delimiterConfigDest.LogPath, delimiterConfig.LogPath)
	s.Equal(delimiterConfigDest.LogType, delimiterConfig.LogType)
	s.Equal(delimiterConfigDest.FilePattern, delimiterConfig.FilePattern)
	s.Nil(s.client.ApplyConfigToMachineGroup(s.projectName, configName, s.machineGroupName))
}

func (s *ConfigTestSuite) TestPluginConfig() {
	configName := "go-sdk-plugin-config"
	s.client.DeleteConfig(s.projectName, configName)
	pluginConfig := &PluginLogConfigInputDetail{}
	InitPluginLogConfigInputDetail(pluginConfig)
	config := &LogConfig{
		Name:        configName,
		InputDetail: pluginConfig,
		InputType:   InputTypePlugin,
		OutputType:  OutputTypeLogService,
		OutputDetail: OutputDetail{
			ProjectName:  s.projectName,
			LogStoreName: s.logstoreName,
		},
	}
	dockerStdoutPlugin := LogConfigPluginInput{}
	dockerStdoutPluginDetail := CreateConfigPluginDockerStdout()
	dockerStdoutPluginDetail.IncludeEnv = map[string]string{
		"x":    "y",
		"dddd": "",
	}
	dockerStdoutPluginDetail.ExcludeEnv = map[string]string{
		"no_this_env": "",
	}
	dockerStdoutPlugin.Inputs = append(dockerStdoutPlugin.Inputs, CreatePluginInputItem(PluginInputTypeDockerStdout, dockerStdoutPluginDetail))

	pluginConfig.PluginDetail = dockerStdoutPlugin

	//b, err := json.Marshal(config)
	//fmt.Println("config:", string(b))
	s.Nil(s.client.CreateConfig(s.projectName, config))

	// check get config
	destConfig, err := s.client.GetConfig(s.projectName, configName)
	s.Nil(err)
	s.Equal(destConfig.Name, configName)
	s.Equal(destConfig.InputType, InputTypePlugin)
	s.Equal(destConfig.OutputDetail.ProjectName, s.projectName)
	s.Equal(destConfig.OutputDetail.LogStoreName, s.logstoreName)
	s.Equal(destConfig.OutputType, OutputTypeLogService)
	pluginConfigDest, ok := ConvertToPluginLogConfigInputDetail(destConfig.InputDetail)
	s.True(ok)
	destBytes, _ := json.Marshal(&pluginConfigDest.PluginDetail)
	srcBytes, _ := json.Marshal(&pluginConfig.PluginDetail)
	s.JSONEq(string(destBytes), string(srcBytes))
	s.Nil(s.client.ApplyConfigToMachineGroup(s.projectName, configName, s.machineGroupName))
}
