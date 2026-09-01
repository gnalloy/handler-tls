# handler-tls

[English](README.md) | [文档](docs/README.zh-CN.md)

基于 crypto/tls 的 Gnalloy TLS Pipeline Handler，覆盖 ALPN、SNI、StartTLS、Optional TLS、OCSP 与 cipher suite 辅助能力。

该模块提供 Pipeline handler。handler 在 Channel 已存在之后观察、转换、拒绝、延迟、记录或保护消息；除非模块名明确说明，否则不拥有监听 socket 或应用协议。

## 状态

- 导入路径：`gnalloy.org/handler-tls`
- 仓库：`github.com/gnalloy/handler-tls`
- 默认分支：`dev`
- 预览安装：`go get gnalloy.org/handler-tls@dev`
- 许可证：Apache-2.0

## 安装
```bash
go get gnalloy.org/handler-tls@dev
go doc gnalloy.org/handler-tls
GOWORK=off GOTOOLCHAIN=local go test ./... -count=1
```

## 文档
- [概览](docs/overview.zh-CN.md) ([English](docs/overview.md))
- [用法](docs/usage.zh-CN.md) ([English](docs/usage.md))
- [案例](docs/examples.zh-CN.md) ([English](docs/examples.md))
- [配置说明](docs/configuration.zh-CN.md) ([English](docs/configuration.md))
- [测试与性能](docs/testing.zh-CN.md) ([English](docs/testing.md))
- [API 参考](docs/api.zh-CN.md) ([English](docs/api.md))
- [注意事项与备注](docs/notes.zh-CN.md) ([English](docs/notes.md))
- [ADR-001 模块边界](docs/decisions/0001-module-boundary.zh-CN.md) ([English](docs/decisions/0001-module-boundary.md))

## 模块边界

本仓库负责：基于 crypto/tls 的 Gnalloy TLS Pipeline Handler，覆盖 ALPN、SNI、StartTLS、Optional TLS、OCSP 与 cipher suite 辅助能力。

它不吸收相邻模块职责。核心基础能力保留在 `gnalloy.org/gnalloy`；协议 codec、transport、handler、resolver、examples 与 benchmarks 分别由独立仓库负责。

## 包结构
- `gnalloy.org/handler-tls`（`tls`）
- `gnalloy.org/handler-tls/provider/standard`（`standard`）

## Gnalloy 依赖
- `gnalloy.org/gnalloy`

## 常见集成方式
- handler 构造函数通常承载策略：limit、timeout、match rule、日志级别、recorder 或流量预算。
- handler 顺序很重要。检查协议对象的 handler 应放在协议 decoder 之后，出站 encoder 应放在写协议对象的 handler 之后。
- handler 必须显式处理 backpressure 和消息所有权；不要在没有生命周期控制的情况下保留 ByteBuf。
- TLS 策略属于调用方提供的 `crypto/tls.Config`，版本、cipher suite、证书、SNI 与 ALPN 都需要明确选择。

## 当前公共入口

生成的 API 参考列出了完整公共面。当前常用构造函数或 option 类型包括：
- `var ErrNeedInput = errors.New("gnalloy/handler/tls: need more encrypted input") ...`
- `type BytePoolConfig struct{ ... }`
- `type CipherSuiteOptions struct{ ... }`
- `type ClientHelloProviderFunc func(hello ClientHello) (*cryptotls.Config, error)`
- `type Config struct{ ... }`
- `type OCSPConfig struct{ ... }`
- `type ServerNameSelector func(serverName string) (*cryptotls.Config, error)`

## 验证

```bash
GOWORK=off GOTOOLCHAIN=local go test ./... -count=1
GOWORK=off GOTOOLCHAIN=local go vet ./...
GOWORK=off GOTOOLCHAIN=local go test ./... -run '^$' -bench . -benchmem -count=1
```

压测时，将该模块和相应 transport、codec、handler 栈装配后，使用 `gnalloy.org/benchmarks` 或 `gnalloy.org/examples` 中的场景运行。报告必须保留主机、操作系统、payload、并发度、warmup 和 repetition。

## 注意事项
- 本仓库保持窄边界。跨模块行为应在应用、recipes、examples 或 benchmark harness 中装配。
- 公共 API 必须保持 Go 原生和显式；热路径避免运行时扫描、隐藏全局注册表和重反射。
- 网络输入一律视为不可信。配置解析上限，返回类型化错误，不使用 panic 处理输入错误。
- 性能结论必须绑定具体主机、操作系统、协议、payload、并发度、warmup 和 repetition。
- handler 只有放在正确 pipeline 位置时才可复用，位置取决于它观察的消息类型。
