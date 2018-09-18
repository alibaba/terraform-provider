package aliyunclient

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/aliyun/alibaba-cloud-sdk-go/sdk/requests"

	"net/url"

	"regexp"

	"sync"

	"github.com/aliyun/alibaba-cloud-sdk-go/sdk"
	"github.com/aliyun/alibaba-cloud-sdk-go/sdk/endpoints"
	"github.com/aliyun/alibaba-cloud-sdk-go/sdk/resource"
	"github.com/aliyun/alibaba-cloud-sdk-go/sdk/utils"
	"github.com/aliyun/alibaba-cloud-sdk-go/services/cbn"
	"github.com/aliyun/alibaba-cloud-sdk-go/services/cms"
	"github.com/aliyun/alibaba-cloud-sdk-go/services/dds"
	"github.com/aliyun/alibaba-cloud-sdk-go/services/ecs"
	"github.com/aliyun/alibaba-cloud-sdk-go/services/ess"
	"github.com/aliyun/alibaba-cloud-sdk-go/services/ots"
	"github.com/aliyun/alibaba-cloud-sdk-go/services/pvtz"
	"github.com/aliyun/alibaba-cloud-sdk-go/services/r-kvstore"
	"github.com/aliyun/alibaba-cloud-sdk-go/services/rds"
	"github.com/aliyun/alibaba-cloud-sdk-go/services/slb"
	"github.com/aliyun/alibaba-cloud-sdk-go/services/sts"
	"github.com/aliyun/alibaba-cloud-sdk-go/services/vpc"
	"github.com/aliyun/aliyun-log-go-sdk"
	"github.com/aliyun/aliyun-oss-go-sdk/oss"
	"github.com/aliyun/aliyun-tablestore-go-sdk/tablestore"
	"github.com/aliyun/fc-go-sdk"
	"github.com/denverdino/aliyungo/cdn"
	"github.com/denverdino/aliyungo/common"
	"github.com/denverdino/aliyungo/cs"
	"github.com/denverdino/aliyungo/dns"
	"github.com/denverdino/aliyungo/kms"
	"github.com/denverdino/aliyungo/location"
	"github.com/denverdino/aliyungo/ram"
	"github.com/hashicorp/terraform/terraform"
)

// AliyunClient of aliyun
type AliyunClient struct {
	Region   common.Region
	RegionId string
	//In order to build ots table client, add accesskey and secretkey in aliyunclient temporarily.
	AccessKey                    string
	SecretKey                    string
	SecurityToken                string
	OtsInstanceName              string
	accountIdMutex               sync.RWMutex
	config                       *Config
	accountId                    string
	ecsconn                      *ecs.Client
	essconn                      *ess.Client
	rdsconn                      *rds.Client
	vpcconn                      *vpc.Client
	slbconn                      *slb.Client
	ossconn                      *oss.Client
	dnsconn                      *dns.Client
	ramconn                      ram.RamClientInterface
	csconn                       *cs.Client
	cdnconn                      *cdn.CdnClient
	kmsconn                      *kms.Client
	otsconn                      *ots.Client
	cmsconn                      *cms.Client
	logconn                      *sls.Client
	fcconn                       *fc.Client
	cenconn                      *cbn.Client
	pvtzconn                     *pvtz.Client
	ddsconn                      *dds.Client
	stsconn                      *sts.Client
	rkvconn                      *r_kvstore.Client
	tablestoreconnByInstanceName map[string]*tablestore.TableStoreClient
	csprojectconnByKey           map[string]*cs.ProjectClient
}

type ApiVersion string

const (
	ApiVersion20140526 = ApiVersion("2014-05-26")
	ApiVersion20160815 = ApiVersion("2016-08-15")
	ApiVersion20140515 = ApiVersion("2014-05-15")
)

const businessInfoKey = "Terraform"

var goSdkMutex = sync.RWMutex{} // The Go SDK is not thread-safe

// Client for AliyunClient
func (c *Config) Client() (*AliyunClient, error) {
	err := c.loadAndValidate()
	if err != nil {
		return nil, err
	}

	return &AliyunClient{
		config:                       c,
		Region:                       c.Region,
		RegionId:                     c.RegionId,
		AccessKey:                    c.AccessKey,
		SecretKey:                    c.SecretKey,
		SecurityToken:                c.SecurityToken,
		OtsInstanceName:              c.OtsInstanceName,
		accountId:                    c.AccountId,
		tablestoreconnByInstanceName: make(map[string]*tablestore.TableStoreClient),
		csprojectconnByKey:           make(map[string]*cs.ProjectClient),
	}, nil
}

func (client *AliyunClient) RunSafelyWithEcsClient(do func(*ecs.Client) (interface{}, error)) (interface{}, error) {
	goSdkMutex.Lock()
	defer goSdkMutex.Unlock()

	// Initialize the ECS client if necessary
	if client.ecsconn == nil {
		endpoint := loadEndpoint(client.config.RegionId, ECSCode)
		if endpoint != "" {
			endpoints.AddEndpointMapping(client.config.RegionId, string(ECSCode), endpoint)
		}
		ecsconn, err := ecs.NewClientWithOptions(client.config.RegionId, client.getSdkConfig().WithTimeout(60000000000), client.config.getAuthCredential(true))
		if err != nil {
			return nil, fmt.Errorf("unable to initialize the ECS client: %#v", err)
		}

		if _, err := ecsconn.DescribeRegions(ecs.CreateDescribeRegionsRequest()); err != nil {
			return nil, err
		}
		client.ecsconn = ecsconn
	}

	return do(client.ecsconn)
}

func (client *AliyunClient) RunSafelyWithRdsClient(do func(*rds.Client) (interface{}, error)) (interface{}, error) {
	goSdkMutex.Lock()
	defer goSdkMutex.Unlock()

	// Initialize the RDS client if necessary
	if client.rdsconn == nil {
		endpoint := loadEndpoint(client.config.RegionId, RDSCode)
		if endpoint != "" {
			endpoints.AddEndpointMapping(client.config.RegionId, string(RDSCode), endpoint)
		}
		rdsconn, err := rds.NewClientWithOptions(client.config.RegionId, client.getSdkConfig(), client.config.getAuthCredential(true))
		if err != nil {
			return nil, fmt.Errorf("unable to initialize the RDS client: %#v", err)
		}

		client.rdsconn = rdsconn
	}

	return do(client.rdsconn)
}

func (client *AliyunClient) RunSafelyWithSlbClient(do func(*slb.Client) (interface{}, error)) (interface{}, error) {
	goSdkMutex.Lock()
	defer goSdkMutex.Unlock()

	// Initialize the SLB client if necessary
	if client.slbconn == nil {
		endpoint := loadEndpoint(client.config.RegionId, SLBCode)
		if endpoint != "" {
			endpoints.AddEndpointMapping(client.config.RegionId, string(SLBCode), endpoint)
		}
		slbconn, err := slb.NewClientWithOptions(client.config.RegionId, client.getSdkConfig(), client.config.getAuthCredential(true))
		if err != nil {
			return nil, fmt.Errorf("unable to initialize the SLB client: %#v", err)
		}

		client.slbconn = slbconn
	}

	return do(client.slbconn)
}

func (client *AliyunClient) RunSafelyWithVpcClient(do func(*vpc.Client) (interface{}, error)) (interface{}, error) {
	goSdkMutex.Lock()
	defer goSdkMutex.Unlock()

	// Initialize the VPC client if necessary
	if client.vpcconn == nil {
		endpoint := loadEndpoint(client.config.RegionId, VPCCode)
		if endpoint != "" {
			endpoints.AddEndpointMapping(client.config.RegionId, string(VPCCode), endpoint)
		}
		vpcconn, err := vpc.NewClientWithOptions(client.config.RegionId, client.getSdkConfig(), client.config.getAuthCredential(true))
		if err != nil {
			return nil, fmt.Errorf("unable to initialize the VPC client: %#v", err)
		}

		client.vpcconn = vpcconn
	}

	return do(client.vpcconn)
}

func (client *AliyunClient) RunSafelyWithCenClient(do func(*cbn.Client) (interface{}, error)) (interface{}, error) {
	goSdkMutex.Lock()
	defer goSdkMutex.Unlock()

	// Initialize the CEN client if necessary
	if client.cenconn == nil {
		endpoint := loadEndpoint(client.config.RegionId, CENCode)
		if endpoint != "" {
			endpoints.AddEndpointMapping(client.config.RegionId, string(CENCode), endpoint)
		}
		cenconn, err := cbn.NewClientWithOptions(client.config.RegionId, client.getSdkConfig(), client.config.getAuthCredential(true))
		if err != nil {
			return nil, fmt.Errorf("unable to initialize the CEN client: %#v", err)
		}

		client.cenconn = cenconn
	}

	return do(client.cenconn)
}

func (client *AliyunClient) RunSafelyWithEssClient(do func(*ess.Client) (interface{}, error)) (interface{}, error) {
	goSdkMutex.Lock()
	defer goSdkMutex.Unlock()

	// Initialize the ESS client if necessary
	if client.essconn == nil {
		endpoint := loadEndpoint(client.config.RegionId, ESSCode)
		if endpoint != "" {
			endpoints.AddEndpointMapping(client.config.RegionId, string(ESSCode), endpoint)
		}
		essconn, err := ess.NewClientWithOptions(client.config.RegionId, client.getSdkConfig(), client.config.getAuthCredential(true))
		if err != nil {
			return nil, fmt.Errorf("unable to initialize the ESS client: %#v", err)
		}

		client.essconn = essconn
	}

	return do(client.essconn)
}

func (client *AliyunClient) RunSafelyWithOssClient(do func(*oss.Client) (interface{}, error)) (interface{}, error) {
	goSdkMutex.Lock()
	defer goSdkMutex.Unlock()

	// Initialize the OSS client if necessary
	if client.ossconn == nil {
		endpointClient := location.NewClient(client.config.AccessKey, client.config.SecretKey)
		endpointClient.SetSecurityToken(client.config.SecurityToken)
		endpoint := loadEndpoint(client.config.RegionId, OSSCode)
		if endpoint == "" {
			endpointItem := client.describeEndpointForService(OSSCode)
			if endpointItem != nil {
				endpoint = strings.ToLower(endpointItem.Protocols.Protocols[0]) + "://" + endpointItem.Endpoint
			} else {
				endpoint = fmt.Sprintf("http://oss-%s.aliyuncs.com", client.RegionId)
			}
		}

		log.Printf("[DEBUG] Instantiate OSS client using endpoint: %#v", endpoint)
		clientOptions := []oss.ClientOption{oss.UserAgent(client.getUserAgent())}
		proxyUrl := client.getHttpProxyUrl()
		if proxyUrl != nil {
			clientOptions = append(clientOptions, oss.Proxy(proxyUrl.String()))
		}

		ossconn, err := oss.New(endpoint, client.config.AccessKey, client.config.SecretKey, clientOptions...)
		if err != nil {
			return nil, fmt.Errorf("unable to initialize the OSS client: %#v", err)
		}

		client.ossconn = ossconn
	}

	return do(client.ossconn)
}

func (client *AliyunClient) RunSafelyWithDnsClient(do func(*dns.Client) (interface{}, error)) (interface{}, error) {
	goSdkMutex.Lock()
	defer goSdkMutex.Unlock()

	// Initialize the DNS client if necessary
	if client.dnsconn == nil {
		dnsconn := dns.NewClientNew(client.config.AccessKey, client.config.SecretKey)
		dnsconn.SetBusinessInfo(businessInfoKey)
		dnsconn.SetUserAgent(client.getUserAgent())
		dnsconn.SetSecurityToken(client.config.SecurityToken)

		client.dnsconn = dnsconn
	}

	return do(client.dnsconn)
}

func (client *AliyunClient) RunSafelyWithRamClient(do func(ram.RamClientInterface) (interface{}, error)) (interface{}, error) {
	goSdkMutex.Lock()
	defer goSdkMutex.Unlock()

	// Initialize the RAM client if necessary
	if client.ramconn == nil {
		ramconn := ram.NewClientWithSecurityToken(client.config.AccessKey, client.config.SecretKey, client.config.SecurityToken)
		client.ramconn = ramconn
	}

	return do(client.ramconn)
}

func (client *AliyunClient) RunSafelyWithCsClient(do func(*cs.Client) (interface{}, error)) (interface{}, error) {
	goSdkMutex.Lock()
	defer goSdkMutex.Unlock()

	// Initialize the CS client if necessary
	if client.csconn == nil {
		csconn := cs.NewClientForAussumeRole(client.config.AccessKey, client.config.SecretKey, client.config.SecurityToken)
		csconn.SetUserAgent(client.getUserAgent())
		client.csconn = csconn
	}

	return do(client.csconn)
}

func (client *AliyunClient) RunSafelyWithCdnClient(do func(*cdn.CdnClient) (interface{}, error)) (interface{}, error) {
	goSdkMutex.Lock()
	defer goSdkMutex.Unlock()

	// Initialize the CDN client if necessary
	if client.cdnconn == nil {
		cdnconn := cdn.NewClient(client.config.AccessKey, client.config.SecretKey)
		cdnconn.SetBusinessInfo(businessInfoKey)
		cdnconn.SetUserAgent(client.getUserAgent())
		cdnconn.SetSecurityToken(client.config.SecurityToken)
		client.cdnconn = cdnconn
	}

	return do(client.cdnconn)
}

func (client *AliyunClient) RunSafelyWithKmsClient(do func(*kms.Client) (interface{}, error)) (interface{}, error) {
	goSdkMutex.Lock()
	defer goSdkMutex.Unlock()

	// Initialize the KMS client if necessary
	if client.kmsconn == nil {
		kmsconn := kms.NewECSClientWithSecurityToken(client.config.AccessKey, client.config.SecretKey, client.config.SecurityToken, client.config.Region)
		kmsconn.SetBusinessInfo(businessInfoKey)
		kmsconn.SetUserAgent(client.getUserAgent())
		client.kmsconn = kmsconn
	}

	return do(client.kmsconn)
}

func (client *AliyunClient) RunSafelyWithOtsClient(do func(*ots.Client) (interface{}, error)) (interface{}, error) {
	goSdkMutex.Lock()
	defer goSdkMutex.Unlock()

	// Initialize the OTS client if necessary
	if client.otsconn == nil {
		endpoint := loadEndpoint(client.config.RegionId, OTSCode)
		if endpoint != "" {
			endpoints.AddEndpointMapping(client.config.RegionId, string(OTSCode), endpoint)
		}
		otsconn, err := ots.NewClientWithOptions(client.config.RegionId, client.getSdkConfig(), client.config.getAuthCredential(true))
		if err != nil {
			return nil, fmt.Errorf("unable to initialize the OTS client: %#v", err)
		}

		client.otsconn = otsconn
	}

	return do(client.otsconn)
}

func (client *AliyunClient) RunSafelyWithCmsClient(do func(*cms.Client) (interface{}, error)) (interface{}, error) {
	goSdkMutex.Lock()
	defer goSdkMutex.Unlock()

	// Initialize the CMS client if necessary
	if client.cmsconn == nil {
		cmsconn, err := cms.NewClientWithOptions(client.config.RegionId, client.getSdkConfig(), client.config.getAuthCredential(false))
		if err != nil {
			return nil, fmt.Errorf("unable to initialize the CMS client: %#v", err)
		}

		client.cmsconn = cmsconn
	}

	return do(client.cmsconn)
}

func (client *AliyunClient) RunSafelyWithPvtzClient(do func(*pvtz.Client) (interface{}, error)) (interface{}, error) {
	goSdkMutex.Lock()
	defer goSdkMutex.Unlock()

	// Initialize the PVTZ client if necessary
	if client.pvtzconn == nil {
		endpoint := loadEndpoint(client.config.RegionId, PVTZCode)
		if endpoint != "" {
			endpoints.AddEndpointMapping(client.config.RegionId, string(PVTZCode), endpoint)
		} else {
			endpoints.AddEndpointMapping(client.config.RegionId, string(PVTZCode), "pvtz.aliyuncs.com")
		}
		pvtzconn, err := pvtz.NewClientWithOptions(client.config.RegionId, client.getSdkConfig(), client.config.getAuthCredential(true))
		if err != nil {
			return nil, fmt.Errorf("unable to initialize the PVTZ client: %#v", err)
		}

		client.pvtzconn = pvtzconn
	}

	return do(client.pvtzconn)
}

func (client *AliyunClient) RunSafelyWithStsClient(do func(*sts.Client) (interface{}, error)) (interface{}, error) {
	goSdkMutex.Lock()
	defer goSdkMutex.Unlock()

	// Initialize the STS client if necessary
	if client.stsconn == nil {
		endpoint := loadEndpoint(client.config.RegionId, STSCode)
		if endpoint != "" {
			endpoints.AddEndpointMapping(client.config.RegionId, string(STSCode), endpoint)
		}
		stsconn, err := sts.NewClientWithOptions(client.config.RegionId, client.getSdkConfig(), client.config.getAuthCredential(true))
		if err != nil {
			return nil, fmt.Errorf("unable to initialize the STS client: %#v", err)
		}

		client.stsconn = stsconn
	}

	return do(client.stsconn)
}

func (client *AliyunClient) RunSafelyWithLogClient(do func(*sls.Client) (interface{}, error)) (interface{}, error) {
	goSdkMutex.Lock()
	defer goSdkMutex.Unlock()

	// Initialize the LOG client if necessary
	if client.logconn == nil {
		endpoint := client.config.LogEndpoint
		if endpoint == "" {
			endpoint = loadEndpoint(client.config.RegionId, LOGCode)
			if endpoint == "" {
				endpoint = fmt.Sprintf("%s.log.aliyuncs.com", client.config.RegionId)
			}
		}

		client.logconn = &sls.Client{
			AccessKeyID:     client.config.AccessKey,
			AccessKeySecret: client.config.SecretKey,
			Endpoint:        endpoint,
			SecurityToken:   client.config.SecurityToken,
			UserAgent:       client.getUserAgent(),
		}
	}

	return do(client.logconn)
}

func (client *AliyunClient) RunSafelyWithDdsClient(do func(*dds.Client) (interface{}, error)) (interface{}, error) {
	goSdkMutex.Lock()
	defer goSdkMutex.Unlock()

	// Initialize the DDS client if necessary
	if client.ddsconn == nil {
		endpoint := loadEndpoint(client.config.RegionId, DDSCode)
		if endpoint != "" {
			endpoints.AddEndpointMapping(client.config.RegionId, string(DDSCode), endpoint)
		}
		ddsconn, err := dds.NewClientWithOptions(client.config.RegionId, client.getSdkConfig(), client.config.getAuthCredential(true))
		if err != nil {
			return nil, fmt.Errorf("unable to initialize the DDS client: %#v", err)
		}

		client.ddsconn = ddsconn
	}

	return do(client.ddsconn)
}

func (client *AliyunClient) RunSafelyWithRkvClient(do func(*r_kvstore.Client) (interface{}, error)) (interface{}, error) {
	goSdkMutex.Lock()
	defer goSdkMutex.Unlock()

	// Initialize the RKV client if necessary
	if client.rkvconn == nil {
		endpoint := loadEndpoint(client.config.RegionId, KVSTORECode)
		if endpoint != "" {
			endpoints.AddEndpointMapping(client.config.RegionId, fmt.Sprintf("R-%s", string(KVSTORECode)), endpoint)
		}
		rkvconn, err := r_kvstore.NewClientWithOptions(client.config.RegionId, client.getSdkConfig(), client.config.getAuthCredential(true))
		if err != nil {
			return nil, fmt.Errorf("unable to initialize the RKV client: %#v", err)
		}

		client.rkvconn = rkvconn
	}

	return do(client.rkvconn)
}

func (client *AliyunClient) RunSafelyWithFcClient(do func(*fc.Client) (interface{}, error)) (interface{}, error) {
	goSdkMutex.Lock()
	defer goSdkMutex.Unlock()

	// Initialize the FC client if necessary
	if client.fcconn == nil {
		endpoint := client.config.LogEndpoint
		if endpoint == "" {
			endpoint = loadEndpoint(client.config.RegionId, FCCode)
			if endpoint == "" {
				endpoint = fmt.Sprintf("%s.fc.aliyuncs.com", client.config.RegionId)
			}
		}

		accountId, err := client.AccountId()
		if err != nil {
			return nil, err
		}

		config := client.getSdkConfig()
		fcconn, err := fc.NewClient(
			fmt.Sprintf("%s.%s", accountId, endpoint),
			string(ApiVersion20160815),
			client.config.AccessKey,
			client.config.SecretKey,
			fc.WithTransport(config.HttpTransport))
		if err != nil {
			return nil, fmt.Errorf("unable to initialize the FC client: %#v", err)
		}

		fcconn.Config.UserAgent = client.getUserAgent()
		fcconn.Config.SecurityToken = client.config.SecurityToken
		client.fcconn = fcconn
	}

	return do(client.fcconn)
}

func (client *AliyunClient) RunSafelyWithTableStoreClient(instanceName string, do func(*tablestore.TableStoreClient) (interface{}, error)) (interface{}, error) {
	goSdkMutex.Lock()
	defer goSdkMutex.Unlock()

	// Initialize the TABLESTORE client if necessary
	tableStoreClient, ok := client.tablestoreconnByInstanceName[instanceName]
	if !ok {
		endpoint := loadEndpoint(client.RegionId, OTSCode)
		if endpoint == "" {
			endpoint = fmt.Sprintf("%s.%s.ots.aliyuncs.com", instanceName, client.RegionId)
		}
		if !strings.HasPrefix(endpoint, "https") && !strings.HasPrefix(endpoint, "http") {
			endpoint = "https://" + endpoint
		}
		tableStoreClient = tablestore.NewClient(endpoint, instanceName, client.AccessKey, client.SecretKey)
		client.tablestoreconnByInstanceName[instanceName] = tableStoreClient
	}

	return do(tableStoreClient)
}

func (client *AliyunClient) RunSafelyWithCsProjectClient(clusterId, endpoint string, clusterCerts cs.ClusterCerts, do func(*cs.ProjectClient) (interface{}, error)) (interface{}, error) {
	goSdkMutex.Lock()
	defer goSdkMutex.Unlock()

	// Initialize the PROJECT client if necessary
	key := fmt.Sprintf("%s|%s|%s|%s|%s", clusterId, endpoint, clusterCerts.CA, clusterCerts.Cert, clusterCerts.Key)
	csProjectClient, ok := client.csprojectconnByKey[key]
	if !ok {
		var err error
		csProjectClient, err = cs.NewProjectClient(clusterId, endpoint, clusterCerts)
		if err != nil {
			return nil, fmt.Errorf("Getting Application Client failed by cluster id %s: %#v.", clusterCerts, err)
		}
		csProjectClient.SetDebug(false)
		csProjectClient.SetUserAgent(client.getUserAgent())
		client.csprojectconnByKey[key] = csProjectClient
	}

	return do(csProjectClient)
}

func (client *AliyunClient) NewCommonRequest(serviceCode ServiceCode, apiVersion ApiVersion) *requests.CommonRequest {
	request := requests.NewCommonRequest()
	endpoint := loadEndpoint(client.RegionId, serviceCode)
	if endpoint == "" {
		endpointItem := client.describeEndpointForService(serviceCode)
		if endpointItem != nil {
			endpoint = endpointItem.Endpoint
		}
	}
	if endpoint == "" {
		switch serviceCode {
		case ECSCode:
			endpoint = "ecs.aliyuncs.com"
		case VPCCode:
			endpoint = fmt.Sprintf("vpc.%s.aliyuncs.com", client.RegionId)
		case SLBCode:
			endpoint = fmt.Sprintf("slb.%s.aliyuncs.com", client.RegionId)
		case ESSCode:
			endpoint = "ess.aliyuncs.com"
		}
	}
	request.Domain = endpoint
	request.Version = string(apiVersion)
	request.RegionId = client.RegionId
	return request
}

func (client *AliyunClient) AccountId() (string, error) {
	client.accountIdMutex.Lock()
	defer client.accountIdMutex.Unlock()

	if client.accountId == "" {
		log.Printf("[DEBUG] account_id not provided, attempting to retrieve it automatically...")
		identity, err := client.getCallerIdentity()
		if err != nil {
			return "", err
		}
		if identity.AccountId == "" {
			return "", fmt.Errorf("caller identity doesn't contain any AccountId")
		}
		log.Printf("[DEBUG] account_id retrieved with success.")
		client.accountId = identity.AccountId
	}
	return client.accountId, nil
}

func (client *AliyunClient) getSdkConfig() *sdk.Config {
	// Fix bug "open /usr/local/go/lib/time/zoneinfo.zip: no such file or directory" which happened in windows.
	if data, ok := resource.GetTZData("GMT"); ok {
		utils.TZData = data
		utils.LoadLocationFromTZData = time.LoadLocationFromTZData
	}
	return sdk.NewConfig().
		WithMaxRetryTime(5).
		WithTimeout(time.Duration(30000000000)).
		WithUserAgent(client.getUserAgent()).
		WithGoRoutinePoolSize(10).
		WithDebug(false).
		WithHttpTransport(client.getTransport())
}

func (client *AliyunClient) getUserAgent() string {
	return fmt.Sprintf("HashiCorp-Terraform-v%s", strings.TrimSuffix(terraform.VersionString(), "-dev"))
}

func (client *AliyunClient) getTransport() *http.Transport {
	handshakeTimeout, err := strconv.Atoi(os.Getenv("TLSHandshakeTimeout"))
	if err != nil {
		handshakeTimeout = 120
	}
	transport := &http.Transport{}
	transport.TLSHandshakeTimeout = time.Duration(handshakeTimeout) * time.Second

	// After building a new transport and it need to set http proxy to support proxy.
	proxyUrl := client.getHttpProxyUrl()
	if proxyUrl != nil {
		transport.Proxy = http.ProxyURL(proxyUrl)
	}
	return transport
}

func (client *AliyunClient) getHttpProxyUrl() *url.URL {
	for _, v := range []string{"HTTPS_PROXY", "https_proxy", "HTTP_PROXY", "http_proxy"} {
		value := strings.Trim(os.Getenv(v), " ")
		if value != "" {
			if !regexp.MustCompile(`^http(s)?://`).MatchString(value) {
				value = fmt.Sprintf("http://%s", value)
			}
			proxyUrl, err := url.Parse(value)
			if err == nil {
				return proxyUrl
			}
			break
		}
	}
	return nil
}

func (client *AliyunClient) describeEndpointForService(serviceCode ServiceCode) *location.EndpointItem {
	args := &location.DescribeEndpointsArgs{
		Id:          common.Region(client.RegionId),
		ServiceCode: strings.ToLower(string(serviceCode)),
		Type:        "openAPI",
	}
	locationClient := location.NewClient(client.AccessKey, client.SecretKey)
	locationClient.SetSecurityToken(client.SecurityToken)
	endpointsResponse, err := locationClient.DescribeEndpoints(args)
	if err != nil {
		log.Printf("[DEBUG] Describe %s endpoint using region: %#v got an error: %#v.", serviceCode, client.RegionId, err)
	} else if endpointsResponse != nil && len(endpointsResponse.Endpoints.Endpoint) > 0 {
		endpointItem := endpointsResponse.Endpoints.Endpoint
		return &endpointItem[0]
	}
	return nil
}

func (client *AliyunClient) getCallerIdentity() (*sts.GetCallerIdentityResponse, error) {
	args := sts.CreateGetCallerIdentityRequest()
	args.Scheme = "https"

	endpoint := loadEndpoint(client.config.RegionId, STSCode)
	if endpoint != "" {
		endpoints.AddEndpointMapping(client.config.RegionId, string(STSCode), endpoint)
	}
	stsClient, err := sts.NewClientWithOptions(client.config.RegionId, client.getSdkConfig(), client.config.getAuthCredential(true))
	if err != nil {
		return nil, fmt.Errorf("unable to initialize the STS client: %#v", err)
	}

	identity, err := stsClient.GetCallerIdentity(args)
	if err != nil {
		return nil, err
	}
	if identity == nil {
		return nil, fmt.Errorf("caller identity not found")
	}
	return identity, err
}
