package vpc

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

// AddGlobalAccelerationInstanceIp invokes the vpc.AddGlobalAccelerationInstanceIp API synchronously
// api document: https://help.aliyun.com/api/vpc/addglobalaccelerationinstanceip.html
func (client *Client) AddGlobalAccelerationInstanceIp(request *AddGlobalAccelerationInstanceIpRequest) (response *AddGlobalAccelerationInstanceIpResponse, err error) {
	response = CreateAddGlobalAccelerationInstanceIpResponse()
	err = client.DoAction(request, response)
	return
}

// AddGlobalAccelerationInstanceIpWithChan invokes the vpc.AddGlobalAccelerationInstanceIp API asynchronously
// api document: https://help.aliyun.com/api/vpc/addglobalaccelerationinstanceip.html
// asynchronous document: https://help.aliyun.com/document_detail/66220.html
func (client *Client) AddGlobalAccelerationInstanceIpWithChan(request *AddGlobalAccelerationInstanceIpRequest) (<-chan *AddGlobalAccelerationInstanceIpResponse, <-chan error) {
	responseChan := make(chan *AddGlobalAccelerationInstanceIpResponse, 1)
	errChan := make(chan error, 1)
	err := client.AddAsyncTask(func() {
		defer close(responseChan)
		defer close(errChan)
		response, err := client.AddGlobalAccelerationInstanceIp(request)
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

// AddGlobalAccelerationInstanceIpWithCallback invokes the vpc.AddGlobalAccelerationInstanceIp API asynchronously
// api document: https://help.aliyun.com/api/vpc/addglobalaccelerationinstanceip.html
// asynchronous document: https://help.aliyun.com/document_detail/66220.html
func (client *Client) AddGlobalAccelerationInstanceIpWithCallback(request *AddGlobalAccelerationInstanceIpRequest, callback func(response *AddGlobalAccelerationInstanceIpResponse, err error)) <-chan int {
	result := make(chan int, 1)
	err := client.AddAsyncTask(func() {
		var response *AddGlobalAccelerationInstanceIpResponse
		var err error
		defer close(result)
		response, err = client.AddGlobalAccelerationInstanceIp(request)
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

// AddGlobalAccelerationInstanceIpRequest is the request struct for api AddGlobalAccelerationInstanceIp
type AddGlobalAccelerationInstanceIpRequest struct {
	*requests.RpcRequest
	ResourceOwnerId              requests.Integer `position:"Query" name:"ResourceOwnerId"`
	GlobalAccelerationInstanceId string           `position:"Query" name:"GlobalAccelerationInstanceId"`
	ResourceOwnerAccount         string           `position:"Query" name:"ResourceOwnerAccount"`
	OwnerAccount                 string           `position:"Query" name:"OwnerAccount"`
	OwnerId                      requests.Integer `position:"Query" name:"OwnerId"`
	IpInstanceId                 string           `position:"Query" name:"IpInstanceId"`
}

// AddGlobalAccelerationInstanceIpResponse is the response struct for api AddGlobalAccelerationInstanceIp
type AddGlobalAccelerationInstanceIpResponse struct {
	*responses.BaseResponse
	RequestId string `json:"RequestId" xml:"RequestId"`
}

// CreateAddGlobalAccelerationInstanceIpRequest creates a request to invoke AddGlobalAccelerationInstanceIp API
func CreateAddGlobalAccelerationInstanceIpRequest() (request *AddGlobalAccelerationInstanceIpRequest) {
	request = &AddGlobalAccelerationInstanceIpRequest{
		RpcRequest: &requests.RpcRequest{},
	}
	request.InitWithApiInfo("Vpc", "2016-04-28", "AddGlobalAccelerationInstanceIp", "vpc", "openAPI")
	request.Method = requests.POST
	return
}

// CreateAddGlobalAccelerationInstanceIpResponse creates a response to parse from AddGlobalAccelerationInstanceIp response
func CreateAddGlobalAccelerationInstanceIpResponse() (response *AddGlobalAccelerationInstanceIpResponse) {
	response = &AddGlobalAccelerationInstanceIpResponse{
		BaseResponse: &responses.BaseResponse{},
	}
	return
}
