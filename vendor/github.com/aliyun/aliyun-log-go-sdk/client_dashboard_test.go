package sls

import (
	"os"
	"testing"

	"github.com/golang/glog"
	"github.com/stretchr/testify/suite"
)

func TestDashboard(t *testing.T) {
	suite.Run(t, new(DashboardTestSuite))
	glog.Flush()
}

type DashboardTestSuite struct {
	suite.Suite
	endpoint        string
	projectName     string
	logstoreName    string
	accessKeyID     string
	accessKeySecret string
	client          Client
}

func (s *DashboardTestSuite) SetupSuite() {
	s.endpoint = os.Getenv("LOG_TEST_ENDPOINT")
	s.projectName = os.Getenv("LOG_TEST_PROJECT")
	s.logstoreName = os.Getenv("LOG_TEST_LOGSTORE")
	s.accessKeyID = os.Getenv("LOG_TEST_ACCESS_KEY_ID")
	s.accessKeySecret = os.Getenv("LOG_TEST_ACCESS_KEY_SECRET")
	s.client.AccessKeyID = s.accessKeyID
	s.client.AccessKeySecret = s.accessKeySecret
	s.client.Endpoint = s.endpoint
	s.Nil(makeSureLogstoreExist(&s.client, s.projectName, s.logstoreName))
}

func (s *DashboardTestSuite) TearDownSuite() {
	//err := s.client.DeleteMachineGroup(s.projectName, s.machineGroupName)
	//s.Nil(err)
}

func (s *DashboardTestSuite) TestDashboard() {
	// @todo
}

func (s *DashboardTestSuite) TestChart() {
	// @todo
}
