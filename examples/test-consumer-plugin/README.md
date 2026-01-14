# Test Consumer Plugin

这是一个测试插件，用于验证 fork 的 wasm-go SDK 的 consumer 功能。

## 功能

- 检测请求中的 consumer 信息
- 如果存在 consumer，添加 `X-Consumer-Name` 和 `X-Consumer-Message` 请求头
- 如果不存在 consumer，添加 `X-No-Consumer` 请求头
- 记录配置信息到日志

## 编译

### 标准 Go 编译（用于测试）
```bash
go build -o test-plugin main.go
```

### WASM 编译（用于部署）
需要先安装 tinygo:
```bash
# macOS
brew install tinygo

# 或从官网下载: https://tinygo.org/getting-started/install/
```

然后编译为 WASM:
```bash
tinygo build -o main.wasm -scheduler=none -target=wasi -gc=custom -tags='custommalloc nottinygc_finalizer' main.go
```

## 配置示例

### 包含 consumer 规则的配置

```yaml
apiVersion: extensions.higress.io/v1alpha1
kind: WasmPlugin
metadata:
  name: test-consumer-plugin
spec:
  phase: AUTHN
  priority: 800
  matchRules:
  - _match_consumer_:
    - "premium-user"
    - "enterprise-user"
    config:
      message: "Welcome premium/enterprise user!"
  - _match_route_:
    - "default-route"
    config:
      message: "Welcome regular user!"
  config:
    message: "Welcome guest!"
```

### 测试 consumer 优先级

当请求同时匹配 consumer 和 route 规则时，consumer 规则优先：

1. 如果请求包含 `consumer_name=premium-user`，使用 consumer 配置（message: "Welcome premium/enterprise user!"）
2. 如果请求不包含 consumer 但匹配 route，使用 route 配置（message: "Welcome regular user!"）
3. 如果都不匹配，使用全局配置（message: "Welcome guest!"）

## 验证

编译成功表明：
- ✅ Fork 的 SDK 模块路径正确 (`github.com/ISADBA/wasm-go`)
- ✅ 所有依赖正确解析
- ✅ Consumer 功能代码可用
- ✅ 向后兼容性保持
