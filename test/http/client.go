/**
 * Tencent is pleased to support the open source community by making Polaris available.
 *
 * Copyright (C) 2019 THL A29 Limited, a Tencent company. All rights reserved.
 *
 * Licensed under the BSD 3-Clause License (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 * https://opensource.org/licenses/BSD-3-Clause
 *
 * Unless required by applicable law or agreed to in writing, software distributed
 * under the License is distributed on an "AS IS" BASIS, WITHOUT WARRANTIES OR
 * CONDITIONS OF ANY KIND, either express or implied. See the License for the
 * specific language governing permissions and limitations under the License.
 */

package http

import (
	"bytes"
	"errors"
	"fmt"
	"io"
	"net/http"

	"github.com/golang/protobuf/jsonpb"

	api "github.com/polarismesh/polaris-server/common/api/v1"
)

// NewClient 创建HTTP客户端
func NewClient(address, version string) *Client {
	return &Client{
		Address: address,
		Version: version,
		Worker:  &http.Client{},
	}
}

// Client HTTP客户端
type Client struct {
	Address string
	Version string
	Worker  *http.Client
}

// SendRequest 发送请求 HTTP Post/Put
func (c *Client) SendRequest(method string, url string, body *bytes.Buffer) (*http.Response, error) {
	var request *http.Request
	var err error

	if body == nil {
		request, err = http.NewRequest(method, url, nil)
	} else {
		request, err = http.NewRequest(method, url, body)
	}

	if err != nil {
		return nil, err
	}

	request.Header.Add("Content-Type", "application/json")
	request.Header.Add("Request-Id", "test")
	request.Header.Add("X-Polaris-Token", "nu/0WRA4EqSR1FagrjRj0fZwPXuGlMpX+zCuWu4uMqy8xr1vRjisSbA25aAC3mtU8MeeRsKhQiDAynUR09I=")

	response, err := c.Worker.Do(request)
	if err != nil {
		return nil, err
	}

	return response, nil
}

// CompleteURL 生成GET请求的完整URL
func (c *Client) CompleteURL(url string, params map[string][]interface{}) string {
	count := 1
	url += "?"

	num := 0
	for _, param := range params {
		num += len(param)
	}

	for index, param := range params {
		for _, item := range param {
			url += fmt.Sprintf("%v=%v", index, item)
			if count != num {
				url += "&"
			}
			count++
		}
	}
	return url
}

// GetBatchWriteResponse 获取BatchWriteResponse
func GetBatchWriteResponse(response *http.Response) (*api.BatchWriteResponse, error) {
	// 打印回复
	fmt.Printf("http code: %v\n", response.StatusCode)

	ret := &api.BatchWriteResponse{}
	checkErr := jsonpb.Unmarshal(response.Body, ret)
	if checkErr == nil {
		fmt.Printf("%+v\n", ret)
	} else {
		fmt.Printf("%v\n", checkErr)
	}

	// 检查回复
	if response.StatusCode != 200 {
		return nil, errors.New("invalid http code")
	}

	if checkErr == nil {
		return ret, nil
	} else if checkErr == io.EOF {
		return nil, io.EOF
	} else {
		return nil, errors.New("body decode failed")
	}
}

// GetBatchQueryResponse 获取BatchQueryResponse
func GetBatchQueryResponse(response *http.Response) (*api.BatchQueryResponse, error) {
	// 打印回复
	fmt.Printf("http code: %v\n", response.StatusCode)

	ret := &api.BatchQueryResponse{}
	checkErr := jsonpb.Unmarshal(response.Body, ret)
	if checkErr == nil {
		fmt.Printf("%+v\n", ret)
	} else {
		fmt.Printf("%v\n", checkErr)
	}

	// 检查回复
	if response.StatusCode != 200 {
		return nil, errors.New("invalid http code")
	}

	if checkErr == nil {
		return ret, nil
	} else if checkErr == io.EOF {
		return nil, io.EOF
	} else {
		return nil, errors.New("body decode failed")
	}
}

// GetSimpleResponse 获取SimpleResponse
func GetSimpleResponse(response *http.Response) (*api.Response, error) {
	// 打印回复
	fmt.Printf("http code: %v\n", response.StatusCode)

	ret := &api.Response{}
	checkErr := jsonpb.Unmarshal(response.Body, ret)
	if checkErr == nil {
		fmt.Printf("%+v\n", ret)
	} else {
		fmt.Printf("%v\n", checkErr)
	}

	// 检查回复
	if response.StatusCode != 200 {
		return nil, errors.New("invalid http code")
	}

	if checkErr == nil {
		return ret, nil
	} else if checkErr == io.EOF {
		return nil, io.EOF
	} else {
		return nil, errors.New("body decode failed")
	}
}
