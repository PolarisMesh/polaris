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

package httpserver

import (
	"context"
	"fmt"
	api "github.com/polarismesh/polaris-server/common/api/v1"
	"github.com/polarismesh/polaris-server/common/log"
	"github.com/polarismesh/polaris-server/common/utils"
	"github.com/emicklei/go-restful"
	"github.com/golang/protobuf/jsonpb"
	"github.com/golang/protobuf/proto"
	"go.uber.org/zap"
	"net/http"
	"strings"
)

/**
 * @brief HTTP请求/回复处理器
 */
type Handler struct {
	*restful.Request
	*restful.Response
}

/**
 * @brief 解析请求
 */
func (h *Handler) Parse(message proto.Message) (context.Context, error) {
	requestID := h.Request.HeaderParameter("Request-Id")
	platformID := h.Request.HeaderParameter("Platform-Id")
	platformToken := h.Request.HeaderParameter("Platform-Token")
	token := h.Request.HeaderParameter("Polaris-Token")

	if err := jsonpb.Unmarshal(h.Request.Request.Body, message); err != nil {
		log.Error(err.Error(), zap.String("request-id", requestID))
		return nil, err
	}

	ctx := context.Background()
	ctx = context.WithValue(ctx, utils.StringContext("request-id"), requestID)
	ctx = context.WithValue(ctx, utils.StringContext("platform-id"), platformID)
	ctx = context.WithValue(ctx, utils.StringContext("platform-token"), platformToken)
	if token != "" {
		ctx = context.WithValue(ctx, utils.StringContext("polaris-token"), token)
	}

	var operator string
	addrSlice := strings.Split(h.Request.Request.RemoteAddr, ":")
	if len(addrSlice) == 2 {
		operator = "HTTP:" + addrSlice[0]
		if platformID != "" {
			operator += "(" + platformID + ")"
		}
	}
	if staffName := h.Request.HeaderParameter("Staffname"); staffName != "" {
		operator = staffName
	}
	ctx = context.WithValue(ctx, utils.StringContext("operator"), operator)

	return ctx, nil
}

/**
 * @brief 仅返回Code
 */
func (h *Handler) WriteHeader(polarisCode uint32, httpStatus int) {
	requestID := h.Request.HeaderParameter(utils.PolarisRequestID)
	h.Request.SetAttribute(utils.PolarisCode, polarisCode) // api统计的时候，用该code

	// 对于非200000的返回，补充实际的code到header中
	if polarisCode != api.ExecuteSuccess {
		h.Response.AddHeader(utils.PolarisCode, fmt.Sprintf("%d", polarisCode))
		h.Response.AddHeader(utils.PolarisMessage, api.Code2Info(polarisCode))
	}
	h.Response.AddHeader("Request-Id", requestID)
	h.Response.WriteHeader(httpStatus)
}

/**
 * @brief 返回Code和Proto
 */
func (h *Handler) WriteHeaderAndProto(obj api.ResponseMessage) {
	requestID := h.Request.HeaderParameter(utils.PolarisRequestID)
	h.Request.SetAttribute(utils.PolarisCode, obj.GetCode().GetValue())
	status := api.CalcCode(obj)

	if status != http.StatusOK {
		log.Error(obj.String(), zap.String("request-id", requestID))
	}
	if code := obj.GetCode().GetValue(); code != api.ExecuteSuccess {
		h.Response.AddHeader(utils.PolarisCode, fmt.Sprintf("%d", code))
		h.Response.AddHeader(utils.PolarisMessage, api.Code2Info(code))
	}

	h.Response.AddHeader(utils.PolarisRequestID, requestID)
	h.Response.WriteHeader(status)

	m := jsonpb.Marshaler{Indent: " "}
	err := m.Marshal(h.Response, obj)
	if err != nil {
		log.Error(err.Error(), zap.String("request-id", requestID))
	}
}

// http答复简单封装
func HTTPResponse(req *restful.Request, rsp *restful.Response, code uint32) {
	handler := &Handler{req, rsp}
	resp := api.NewResponse(code)
	handler.WriteHeaderAndProto(resp)
}
