# 配置说明

[English](configuration.md) | [文档索引](README.zh-CN.md)

Gnalloy 模块的配置必须显式。优先使用构造参数、option struct 和应用自有配置文件，不使用包级可变状态。

## 主要配置点
- handler 构造函数通常承载策略：limit、timeout、match rule、日志级别、recorder 或流量预算。
- handler 顺序很重要。检查协议对象的 handler 应放在协议 decoder 之后，出站 encoder 应放在写协议对象的 handler 之后。
- handler 必须显式处理 backpressure 和消息所有权；不要在没有生命周期控制的情况下保留 ByteBuf。
- TLS 策略属于调用方提供的 `crypto/tls.Config`，版本、cipher suite、证书、SNI 与 ALPN 都需要明确选择。

## 推荐默认值

- 从有界大小和较短集成测试超时开始。
- 只有在测量真实 payload 与对端行为后才提高限制。
- 安全相关默认值保持关闭或保守。
- 每个生产覆盖项都应记录在拥有该配置的服务中，而不是写进 library module。

## 环境变量

普通单元测试不要求该 library module 提供仓库专属环境变量。应用、examples、benchmarks 和 CI job 可以围绕它定义自己的变量。

## 本地 Workspace 开发

本地 `go.work` 只作为开发便利。发布用 module metadata 必须保持可移植，不能包含机器相关路径。
