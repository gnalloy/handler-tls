# Usage

[简体中文](usage.zh-CN.md) | [Docs Index](README.md)

## Requirements

- Go 1.25 or newer, matching the module `go` directive.
- A Gnalloy application, recipe, example, or benchmark harness that owns lifecycle and deployment configuration.
- Standalone module verification should set `GOWORK=off` so the module is tested through its published dependency graph.

## Install
```bash
go get gnalloy.org/handler-tls@dev
```

## Import
```go
import "gnalloy.org/handler-tls"
```

## Integration Pattern
- Handler constructors usually carry the policy: limits, timeouts, match rules, logging level, recorder, or traffic budget.
- Handler order matters. Place protocol decoders before handlers that inspect protocol objects, and place outbound encoders after handlers that write protocol objects.
- Handlers must keep backpressure and message ownership explicit; never retain ByteBuf values without clear lifetime control.
- TLS policy belongs in caller-supplied `crypto/tls.Config`; choose versions, cipher suites, certificates, SNI, and ALPN deliberately.

## API Selection

Use the API inventory to choose the exact constructor or option type for your protocol path:

```bash
go doc gnalloy.org/handler-tls
```

Common current entry points:
- `var ErrNeedInput = errors.New("gnalloy/handler/tls: need more encrypted input") ...`
- `type BytePoolConfig struct{ ... }`
- `type CipherSuiteOptions struct{ ... }`
- `type ClientHelloProviderFunc func(hello ClientHello) (*cryptotls.Config, error)`
- `type Config struct{ ... }`
- `type OCSPConfig struct{ ... }`
- `type ServerNameSelector func(serverName string) (*cryptotls.Config, error)`

## Cross-Module Assembly

When multiple Gnalloy repositories are developed together, create a local `go.work` file in your chosen workspace. Keep application-local `replace` directives out of published library modules unless the change is intentionally temporary and never committed.

## Error Handling

Network input, peer behavior, platform capability, and timeout failures must be handled as normal errors. Do not recover protocol correctness by panicking. Return or propagate the module error and close the affected Channel when ownership requires it.
