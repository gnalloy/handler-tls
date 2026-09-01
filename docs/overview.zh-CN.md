# 概览

[English](overview.md) | [文档索引](README.zh-CN.md)

## 目标

基于 crypto/tls 的 Gnalloy TLS Pipeline Handler，覆盖 ALPN、SNI、StartTLS、Optional TLS、OCSP 与 cipher suite 辅助能力。

该模块提供 Pipeline handler。handler 在 Channel 已存在之后观察、转换、拒绝、延迟、记录或保护消息；除非模块名明确说明，否则不拥有监听 socket 或应用协议。

## 仓库身份

- 模块路径：`gnalloy.org/handler-tls`
- GitHub 仓库：`github.com/gnalloy/handler-tls`
- 默认分支：`dev`
- 许可证：Apache-2.0

## 包结构
- `gnalloy.org/handler-tls`（`tls`）
- `gnalloy.org/handler-tls/provider/standard`（`standard`）

## 直接 Gnalloy 依赖
- `gnalloy.org/gnalloy`

## 当前模块规划中的直接下游
- `gnalloy.org/benchmarks`

## 架构位置

Gnalloy 保持核心小而轻依赖。本仓库围绕单一职责形成可替换模块，通过显式 Go package 连接，而不是依靠运行时发现。

## 兼容性

公共导入路径是 `gnalloy.org/handler-tls`。首个稳定 tag 发布前，请按依赖策略使用 `@dev` 或明确的 pseudo-version。
