package emr

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

// ListFailureJobExecutionInstances invokes the emr.ListFailureJobExecutionInstances API synchronously
// api document: https://help.aliyun.com/api/emr/listfailurejobexecutioninstances.html
func (client *Client) ListFailureJobExecutionInstances(request *ListFailureJobExecutionInstancesRequest) (response *ListFailureJobExecutionInstancesResponse, err error) {
	response = CreateListFailureJobExecutionInstancesResponse()
	err = client.DoAction(request, response)
	return
}

// ListFailureJobExecutionInstancesWithChan invokes the emr.ListFailureJobExecutionInstances API asynchronously
// api document: https://help.aliyun.com/api/emr/listfailurejobexecutioninstances.html
// asynchronous document: https://help.aliyun.com/document_detail/66220.html
func (client *Client) ListFailureJobExecutionInstancesWithChan(request *ListFailureJobExecutionInstancesRequest) (<-chan *ListFailureJobExecutionInstancesResponse, <-chan error) {
	responseChan := make(chan *ListFailureJobExecutionInstancesResponse, 1)
	errChan := make(chan error, 1)
	err := client.AddAsyncTask(func() {
		defer close(responseChan)
		defer close(errChan)
		response, err := client.ListFailureJobExecutionInstances(request)
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

// ListFailureJobExecutionInstancesWithCallback invokes the emr.ListFailureJobExecutionInstances API asynchronously
// api document: https://help.aliyun.com/api/emr/listfailurejobexecutioninstances.html
// asynchronous document: https://help.aliyun.com/document_detail/66220.html
func (client *Client) ListFailureJobExecutionInstancesWithCallback(request *ListFailureJobExecutionInstancesRequest, callback func(response *ListFailureJobExecutionInstancesResponse, err error)) <-chan int {
	result := make(chan int, 1)
	err := client.AddAsyncTask(func() {
		var response *ListFailureJobExecutionInstancesResponse
		var err error
		defer close(result)
		response, err = client.ListFailureJobExecutionInstances(request)
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

// ListFailureJobExecutionInstancesRequest is the request struct for api ListFailureJobExecutionInstances
type ListFailureJobExecutionInstancesRequest struct {
	*requests.RpcRequest
	ResourceOwnerId requests.Integer `position:"Query" name:"ResourceOwnerId"`
	Count           requests.Integer `position:"Query" name:"Count"`
}

// ListFailureJobExecutionInstancesResponse is the response struct for api ListFailureJobExecutionInstances
type ListFailureJobExecutionInstancesResponse struct {
	*responses.BaseResponse
	RequestId    string                                         `json:"RequestId" xml:"RequestId"`
	JobInstances JobInstancesInListFailureJobExecutionInstances `json:"JobInstances" xml:"JobInstances"`
}

// CreateListFailureJobExecutionInstancesRequest creates a request to invoke ListFailureJobExecutionInstances API
func CreateListFailureJobExecutionInstancesRequest() (request *ListFailureJobExecutionInstancesRequest) {
	request = &ListFailureJobExecutionInstancesRequest{
		RpcRequest: &requests.RpcRequest{},
	}
	request.InitWithApiInfo("Emr", "2016-04-08", "ListFailureJobExecutionInstances", "emr", "openAPI")
	return
}

// CreateListFailureJobExecutionInstancesResponse creates a response to parse from ListFailureJobExecutionInstances response
func CreateListFailureJobExecutionInstancesResponse() (response *ListFailureJobExecutionInstancesResponse) {
	response = &ListFailureJobExecutionInstancesResponse{
		BaseResponse: &responses.BaseResponse{},
	}
	return
}