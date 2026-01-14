# Fork SDK 迁移指南

本文档描述如何将 Higress 插件从原始 `wasm-go` SDK 迁移到 fork 的版本，以使用 consumer 级别作用域功能。

## 概述

本 fork 版本的主要变更：
- ✅ 模块路径从 `github.com/higress-group/wasm-go` 更改为 `github.com/ISADBA/wasm-go`
- ✅ 新增 Consumer 级别作用域支持
- ✅ Consumer 规则优先级高于 Route/Domain/Global
- ✅ 完全向后兼容，不影响现有功能

## 前提条件

- Go 1.24.1 或更高版本
- Git 访问权限
- （可选）TinyGo 用于编译 WASM

## 迁移步骤

### 1. 更新插件的 go.mod

在你的插件项目中，更新 `go.mod` 文件：

```bash
# 进入插件目录
cd /path/to/your-plugin

# 移除旧依赖
go mod edit -droprequire=github.com/higress-group/wasm-go

# 添加新依赖
go get github.com/ISADBA/wasm-go@v1.0.0-consumer-support

# 整理依赖
go mod tidy
```

**更新后的 go.mod 示例：**
```go
module my-higress-plugin

go 1.24.1

require (
    github.com/ISADBA/wasm-go v1.0.0-consumer-support
    github.com/higress-group/proxy-wasm-go-sdk v0.0.0-20251103120604-77e9cce339d2
)
```

### 2. 更新代码中的 import 路径

批量替换所有 `.go` 文件中的 import 路径：

```bash
# 在插件项目根目录执行
find . -name "*.go" -type f -exec sed -i '' \
  's|github.com/higress-group/wasm-go|github.com/ISADBA/wasm-go|g' {} +
```

**手动更新示例：**

**之前：**
```go
import (
    "github.com/higress-group/wasm-go/pkg/log"
    "github.com/higress-group/wasm-go/pkg/wrapper"
)
```

**之后：**
```go
import (
    "github.com/ISADBA/wasm-go/pkg/log"
    "github.com/ISADBA/wasm-go/pkg/wrapper"
)
```

### 3. 验证编译

```bash
# 下载依赖
go mod download

# 编译验证
go build ./...

# 运行测试
go test ./...
```

### 4. 使用 Consumer 功能（可选）

如果你想使用新的 consumer 功能，可以在配置中添加 consumer 规则：

```yaml
apiVersion: extensions.higress.io/v1alpha1
kind: WasmPlugin
metadata:
  name: my-plugin
spec:
  matchRules:
  # Consumer 规则 - 最高优先级
  - _match_consumer_:
    - "premium-user"
    config:
      # consumer 特定配置
  
  # Route 规则 - 次优先级
  - _match_route_:
    - "my-route"
    config:
      # route 特定配置
  
  # 全局配置 - 最低优先级
  config:
    # 默认配置
```

在插件代码中获取 consumer 信息：

```go
func onHttpRequestHeaders(ctx wrapper.HttpContext, config PluginConfig) types.Action {
    // 获取 consumer 信息
    consumerName, err := proxywasm.GetProperty([]string{"consumer_name"})
    if err == nil && string(consumerName) != "" {
        proxywasm.LogInfof("Consumer: %s", string(consumerName))
        // 处理 consumer 特定逻辑
    }
    return types.ActionContinue
}
```

## 常见问题

### Q1: 模块路径冲突错误

**错误信息：**
```
module declares its path as: github.com/higress-group/wasm-go
but was required as: github.com/ISADBA/wasm-go
```

**解决方案：**
确保 `go.mod` 中正确引用了新的模块路径，并执行 `go clean -modcache` 清理缓存。

### Q2: 找不到版本标签

**错误信息：**
```
unknown revision v1.0.0-consumer-support
```

**解决方案：**
1. 检查标签是否已推送到远程：`git ls-remote --tags origin | grep consumer`
2. 清理模块缓存：`go clean -modcache`
3. 重新下载：`go get github.com/ISADBA/wasm-go@v1.0.0-consumer-support`

### Q3: Import 路径未更新

**错误信息：**
```
no required module provides package github.com/higress-group/wasm-go/pkg/wrapper
```

**解决方案：**
检查所有 `.go` 文件，确保 import 路径已更新：
```bash
grep -r "github.com/higress-group/wasm-go" --include="*.go" .
```

### Q4: Consumer 功能不生效

**可能原因：**
1. 认证插件未设置 `consumer_name` 属性
2. Consumer 规则配置错误
3. 插件优先级设置不当

**解决方案：**
1. 确保认证插件（如 basic-auth）正确设置 consumer 信息
2. 检查配置中的 `_match_consumer_` 字段
3. 设置合适的插件优先级（建议 800+）

## 回滚方案

如果迁移后遇到问题，可以回滚到原始 SDK：

```bash
# 1. 更新 go.mod
go mod edit -droprequire=github.com/ISADBA/wasm-go
go get github.com/higress-group/wasm-go@latest

# 2. 恢复 import 路径
find . -name "*.go" -type f -exec sed -i '' \
  's|github.com/ISADBA/wasm-go|github.com/higress-group/wasm-go|g' {} +

# 3. 清理和验证
go clean -modcache
go mod tidy
go build ./...
```

## 验证清单

迁移完成后，请验证以下项目：

- [ ] `go.mod` 中的模块路径已更新
- [ ] 所有 `.go` 文件的 import 路径已更新
- [ ] `go build ./...` 编译成功
- [ ] `go test ./...` 测试通过
- [ ] （如果使用 consumer 功能）配置文件包含 consumer 规则
- [ ] （如果使用 consumer 功能）插件能够正确识别 consumer

## 示例项目

参考 `examples/test-consumer-plugin/` 目录查看完整的示例项目。

## 支持

如有问题，请：
1. 查看本文档的常见问题部分
2. 检查 `examples/test-consumer-plugin/README.md`
3. 在 GitHub 仓库提交 Issue

## 版本历史

- **v1.0.0-consumer-support** (2024-01-14)
  - 初始 fork 版本
  - 新增 consumer 级别作用域支持
  - 完全向后兼容
