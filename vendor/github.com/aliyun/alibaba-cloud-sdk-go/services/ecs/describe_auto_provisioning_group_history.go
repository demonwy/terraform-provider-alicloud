package ecs

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

// DescribeAutoProvisioningGroupHistory invokes the ecs.DescribeAutoProvisioningGroupHistory API synchronously
// api document: https://help.aliyun.com/api/ecs/describeautoprovisioninggrouphistory.html
func (client *Client) DescribeAutoProvisioningGroupHistory(request *DescribeAutoProvisioningGroupHistoryRequest) (response *DescribeAutoProvisioningGroupHistoryResponse, err error) {
	response = CreateDescribeAutoProvisioningGroupHistoryResponse()
	err = client.DoAction(request, response)
	return
}

// DescribeAutoProvisioningGroupHistoryWithChan invokes the ecs.DescribeAutoProvisioningGroupHistory API asynchronously
// api document: https://help.aliyun.com/api/ecs/describeautoprovisioninggrouphistory.html
// asynchronous document: https://help.aliyun.com/document_detail/66220.html
func (client *Client) DescribeAutoProvisioningGroupHistoryWithChan(request *DescribeAutoProvisioningGroupHistoryRequest) (<-chan *DescribeAutoProvisioningGroupHistoryResponse, <-chan error) {
	responseChan := make(chan *DescribeAutoProvisioningGroupHistoryResponse, 1)
	errChan := make(chan error, 1)
	err := client.AddAsyncTask(func() {
		defer close(responseChan)
		defer close(errChan)
		response, err := client.DescribeAutoProvisioningGroupHistory(request)
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

// DescribeAutoProvisioningGroupHistoryWithCallback invokes the ecs.DescribeAutoProvisioningGroupHistory API asynchronously
// api document: https://help.aliyun.com/api/ecs/describeautoprovisioninggrouphistory.html
// asynchronous document: https://help.aliyun.com/document_detail/66220.html
func (client *Client) DescribeAutoProvisioningGroupHistoryWithCallback(request *DescribeAutoProvisioningGroupHistoryRequest, callback func(response *DescribeAutoProvisioningGroupHistoryResponse, err error)) <-chan int {
	result := make(chan int, 1)
	err := client.AddAsyncTask(func() {
		var response *DescribeAutoProvisioningGroupHistoryResponse
		var err error
		defer close(result)
		response, err = client.DescribeAutoProvisioningGroupHistory(request)
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

// DescribeAutoProvisioningGroupHistoryRequest is the request struct for api DescribeAutoProvisioningGroupHistory
type DescribeAutoProvisioningGroupHistoryRequest struct {
	*requests.RpcRequest
	ResourceOwnerId         requests.Integer `position:"Query" name:"ResourceOwnerId"`
	StartTime               string           `position:"Query" name:"StartTime"`
	PageNumber              requests.Integer `position:"Query" name:"PageNumber"`
	PageSize                requests.Integer `position:"Query" name:"PageSize"`
	ResourceOwnerAccount    string           `position:"Query" name:"ResourceOwnerAccount"`
	OwnerAccount            string           `position:"Query" name:"OwnerAccount"`
	EndTime                 string           `position:"Query" name:"EndTime"`
	OwnerId                 requests.Integer `position:"Query" name:"OwnerId"`
	AutoProvisioningGroupId string           `position:"Query" name:"AutoProvisioningGroupId"`
}

// DescribeAutoProvisioningGroupHistoryResponse is the response struct for api DescribeAutoProvisioningGroupHistory
type DescribeAutoProvisioningGroupHistoryResponse struct {
	*responses.BaseResponse
	RequestId                      string                         `json:"RequestId" xml:"RequestId"`
	TotalCount                     int                            `json:"TotalCount" xml:"TotalCount"`
	PageNumber                     int                            `json:"PageNumber" xml:"PageNumber"`
	PageSize                       int                            `json:"PageSize" xml:"PageSize"`
	AutoProvisioningGroupHistories AutoProvisioningGroupHistories `json:"AutoProvisioningGroupHistories" xml:"AutoProvisioningGroupHistories"`
}

// CreateDescribeAutoProvisioningGroupHistoryRequest creates a request to invoke DescribeAutoProvisioningGroupHistory API
func CreateDescribeAutoProvisioningGroupHistoryRequest() (request *DescribeAutoProvisioningGroupHistoryRequest) {
	request = &DescribeAutoProvisioningGroupHistoryRequest{
		RpcRequest: &requests.RpcRequest{},
	}
	request.InitWithApiInfo("Ecs", "2014-05-26", "DescribeAutoProvisioningGroupHistory", "ecs", "openAPI")
	request.Method = requests.POST
	return
}

// CreateDescribeAutoProvisioningGroupHistoryResponse creates a response to parse from DescribeAutoProvisioningGroupHistory response
func CreateDescribeAutoProvisioningGroupHistoryResponse() (response *DescribeAutoProvisioningGroupHistoryResponse) {
	response = &DescribeAutoProvisioningGroupHistoryResponse{
		BaseResponse: &responses.BaseResponse{},
	}
	return
}
