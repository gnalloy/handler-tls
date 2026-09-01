# 用法

[English](usage.md) | [文档索引](README.zh-CN.md)

## 要求

- Go 1.25 或更新版本，并与 module 的 `go` 指令一致。
- 由 Gnalloy 应用、recipe、example 或 benchmark harness 负责生命周期与部署配置。
- 独立模块复验应设置 `GOWORK=off`，确保通过已发布依赖图测试。

## 安装
```bash
go get gnalloy.org/handler-tls@dev
```

## 导入
```go
import "gnalloy.org/handler-tls"
```

## 集成模式
- handler 构造函数通常承载策略：limit、timeout、match rule、日志级别、recorder 或流量预算。
- handler 顺序很重要。检查协议对象的 handler 应放在协议 decoder 之后，出站 encoder 应放在写协议对象的 handler 之后。
- handler 必须显式处理 backpressure 和消息所有权；不要在没有生命周期控制的情况下保留 ByteBuf。
- TLS 策略属于调用方提供的 `crypto/tls.Config`，版本、cipher suite、证书、SNI 与 ALPN 都需要明确选择。

## API 选择

通过 API 清单选择当前协议路径需要的具体构造函数或 option 类型：

```bash
go doc gnalloy.org/handler-tls
```

当前常用入口：
- `var ErrNeedInput = errors.New("gnalloy/handler/tls: need more encrypted input") ...`
- `type BytePoolConfig struct{ ... }`
- `type CipherSuiteOptions struct{ ... }`
- `type ClientHelloProviderFunc func(hello ClientHello) (*cryptotls.Config, error)`
- `type Config struct{ ... }`
- `type OCSPConfig struct{ ... }`
- `type ServerNameSelector func(serverName string) (*cryptotls.Config, error)`

## 跨模块装配

多个 Gnalloy 仓库一起开发时，在自己选择的 workspace 中创建本地 `go.work` 文件。不要把应用本地 `replace` 指令提交到发布用 library module，除非它是明确的临时变更且不会进入提交。

## 错误处理

网络输入、对端行为、平台能力和超时失败都必须作为普通错误处理。不要用 panic 恢复协议正确性。返回或传播模块错误，并在所有权要求时关闭受影响的 Channel。
