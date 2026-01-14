// Copyright (c) 2024 Alibaba Group Holding Ltd.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//      http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package main

import (
	"github.com/ISADBA/wasm-go/pkg/wrapper"
	"github.com/higress-group/proxy-wasm-go-sdk/proxywasm"
	"github.com/higress-group/proxy-wasm-go-sdk/proxywasm/types"
	"github.com/tidwall/gjson"
)

func main() {
	wrapper.SetCtx(
		"test-consumer-plugin",
		wrapper.ParseConfig(parseConfig),
		wrapper.ProcessRequestHeaders(onHttpRequestHeaders),
	)
}

// PluginConfig 插件配置
type PluginConfig struct {
	Message string `json:"message"`
}

// parseConfig 解析插件配置
func parseConfig(json gjson.Result, config *PluginConfig) error {
	config.Message = json.Get("message").String()
	if config.Message == "" {
		config.Message = "default message"
	}
	return nil
}

// onHttpRequestHeaders 处理请求头阶段
func onHttpRequestHeaders(ctx wrapper.HttpContext, config PluginConfig) types.Action {
	// 尝试获取 consumer 信息
	consumerName, err := proxywasm.GetProperty([]string{"consumer_name"})
	if err == nil && string(consumerName) != "" {
		// 如果存在 consumer 信息，记录日志并添加请求头
		proxywasm.LogInfof("[test-consumer-plugin] Consumer detected: %s", string(consumerName))
		proxywasm.AddHttpRequestHeader("X-Consumer-Name", string(consumerName))
		proxywasm.AddHttpRequestHeader("X-Consumer-Message", config.Message)
	} else {
		// 没有 consumer 信息
		proxywasm.LogInfo("[test-consumer-plugin] No consumer detected")
		proxywasm.AddHttpRequestHeader("X-No-Consumer", "true")
	}
	
	// 记录配置信息
	proxywasm.LogInfof("[test-consumer-plugin] Plugin message: %s", config.Message)
	
	return types.ActionContinue
}
