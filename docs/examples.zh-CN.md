# 案例

[English](examples.md) | [文档索引](README.zh-CN.md)

## 案例 1：将模块加入应用

```bash
mkdir gnalloy-app && cd gnalloy-app
go mod init example.com/gnalloy-app
go get gnalloy.org/handler-tls@dev
go doc gnalloy.org/handler-tls
```

## 案例 2：查看当前包

当前源码树暴露这些 package 导入路径：
- `gnalloy.org/handler-tls`
- `gnalloy.org/handler-tls/provider/standard`

按需要的行为对对应 package 执行 `go doc`：

```bash
go doc gnalloy.org/handler-tls
go doc gnalloy.org/handler-tls/provider/standard
```

精选当前导出入口：
- `var ErrNeedInput = errors.New("gnalloy/handler/tls: need more encrypted input") ...`
- `func CertificateWithOCSPStaple(cert cryptotls.Certificate, response []byte) (cryptotls.Certificate, error)`
- `func CipherSuiteName(id uint16) string`
- `func CipherSuiteNames(ids []uint16) []string`
- `func ConfigureCipherSuiteNames(cfg *cryptotls.Config, value string, options CipherSuiteOptions) error`
- `func ConfigureCipherSuites(cfg *cryptotls.Config, suites []uint16) error`
- `type Provider struct{ ... }`

## 案例 3：将可执行测试作为行为示例

仓库测试是受支持行为的可执行示例。先从下面的精选名称开始，再阅读对应 `_test.go` 文件中的完整 setup 和断言。完整发现列表见 [测试与性能](testing.zh-CN.md)。

```bash
GOWORK=off GOTOOLCHAIN=local go test ./... -run Test -count=1
```

精选当前 test、benchmark、fuzz 与 example 入口：
- `BenchmarkCopyReadableBytesComposite`
- `TestCipherSuiteCatalogIncludesRuntimeSuites`
- `TestCipherSuiteCatalogReportsCertificateAuth`
- `TestCipherSuiteCatalogRequiresInsecureOptIn`
- `TestCloseReleasesQueuedApplicationBuffers`
- `TestConfigureCipherSuitesCopiesTLS12Suites`
- `TestConfigureCipherSuitesRejectsTLS13Suites`
- `TestCopyReadableBytesCopiesCompositeOnce`
- `TestCopyReadableBytesUsesConfiguredPool`
- `TestEvaluateNativeProviderRequiresTLS13ALPNAndQUIC`
- `TestEvaluateProviderDoesNotRequireQUICPacketProtection`
- `TestEvaluateProviderRejectsTypedNil`
- `TestHandlerEmitsStapledOCSPResponse`
- `TestHandlerFlushesRepeatedApplicationWrites`
- `TestHandlerNegotiatesAndPassesPlaintext`
- `TestHandlerNegotiatesConfiguredCipherSuites`
- `TestHandlerNegotiatesConfiguredTLSVersions`
- `TestHandlerRejectsUnsupportedProvider`

## 案例 4：跨模块装配

本模块的直接 Gnalloy 依赖：
- `gnalloy.org/gnalloy`

装配说明：
- handler 在 Gnalloy pipeline 内提供策略、生命周期、防护、日志、指标或流量行为。
- 协议解析留给 codec 模块，socket 所有权留给 transport 模块。
- handler 的位置应匹配所选协议栈的 inbound/outbound 语义。

## 案例 5：压测 Harness

持续负载测试时，如果该模块参与网络流量路径，将它接入 `gnalloy.org/benchmarks` 的场景，或接入 `gnalloy.org/examples` 的可运行客户端。报告中记录 host、OS、CPU、Go version、protocol、payload、concurrency、warmup、repetitions、throughput 和 p99 latency。
