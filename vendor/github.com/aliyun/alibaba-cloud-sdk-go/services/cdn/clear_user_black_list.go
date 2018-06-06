package cdn

//Licensed under the Apache License, Version 2.0 (the "License");
//you may not use this file except in compliance with the License.
//You may obtain a copy of the License at
//
//http://www.apache.org/licenses/LICENSE-2.0
//
//Unless required by applicable law or agreed to in writing, software
//distributed under the License is distributed on an "AS IS" BASIS,
//WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
//See the License for the specific language governing permissions and
//limitations under the License.
//
// Code generated by Alibaba Cloud SDK Code Generator.
// Changes may cause incorrect behavior and will be lost if the code is regenerated.

import (
	"github.com/aliyun/alibaba-cloud-sdk-go/sdk/requests"
	"github.com/aliyun/alibaba-cloud-sdk-go/sdk/responses"
)

// ClearUserBlackList invokes the cdn.ClearUserBlackList API synchronously
// api document: https://help.aliyun.com/api/cdn/clearuserblacklist.html
func (client *Client) ClearUserBlackList(request *ClearUserBlackListRequest) (response *ClearUserBlackListResponse, err error) {
	response = CreateClearUserBlackListResponse()
	err = client.DoAction(request, response)
	return
}

// ClearUserBlackListWithChan invokes the cdn.ClearUserBlackList API asynchronously
// api document: https://help.aliyun.com/api/cdn/clearuserblacklist.html
// asynchronous document: https://help.aliyun.com/document_detail/66220.html
func (client *Client) ClearUserBlackListWithChan(request *ClearUserBlackListRequest) (<-chan *ClearUserBlackListResponse, <-chan error) {
	responseChan := make(chan *ClearUserBlackListResponse, 1)
	errChan := make(chan error, 1)
	err := client.AddAsyncTask(func() {
		defer close(responseChan)
		defer close(errChan)
		response, err := client.ClearUserBlackList(request)
		if err != nil {
			errChan <- err
		} else {
			responseChan <- response
		}
	})
	if err != nil {
		errChan <- err
		close(responseChan)
		close(errChan)
	}
	return responseChan, errChan
}

// ClearUserBlackListWithCallback invokes the cdn.ClearUserBlackList API asynchronously
// api document: https://help.aliyun.com/api/cdn/clearuserblacklist.html
// asynchronous document: https://help.aliyun.com/document_detail/66220.html
func (client *Client) ClearUserBlackListWithCallback(request *ClearUserBlackListRequest, callback func(response *ClearUserBlackListResponse, err error)) <-chan int {
	result := make(chan int, 1)
	err := client.AddAsyncTask(func() {
		var response *ClearUserBlackListResponse
		var err error
		defer close(result)
		response, err = client.ClearUserBlackList(request)
		callback(response, err)
		result <- 1
	})
	if err != nil {
		defer close(result)
		callback(nil, err)
		result <- 0
	}
	return result
}

// ClearUserBlackListRequest is the request struct for api ClearUserBlackList
type ClearUserBlackListRequest struct {
	*requests.RpcRequest
	OwnerId       requests.Integer `position:"Query" name:"OwnerId"`
	OwnerAccount  string           `position:"Query" name:"OwnerAccount"`
	SecurityToken string           `position:"Query" name:"SecurityToken"`
}

// ClearUserBlackListResponse is the response struct for api ClearUserBlackList
type ClearUserBlackListResponse struct {
	*responses.BaseResponse
	RequestId string `json:"RequestId" xml:"RequestId"`
}

// CreateClearUserBlackListRequest creates a request to invoke ClearUserBlackList API
func CreateClearUserBlackListRequest() (request *ClearUserBlackListRequest) {
	request = &ClearUserBlackListRequest{
		RpcRequest: &requests.RpcRequest{},
	}
	request.InitWithApiInfo("Cdn", "2014-11-11", "ClearUserBlackList", "", "")
	return
}

// CreateClearUserBlackListResponse creates a response to parse from ClearUserBlackList response
func CreateClearUserBlackListResponse() (response *ClearUserBlackListResponse) {
	response = &ClearUserBlackListResponse{
		BaseResponse: &responses.BaseResponse{},
	}
	return
}
