# Configuration

[简体中文](configuration.zh-CN.md) | [Docs Index](README.md)

Configuration in Gnalloy modules is explicit. Prefer constructor arguments, option structs, and application-owned configuration files over package-level mutable state.

## Primary Configuration Points
- Handler constructors usually carry the policy: limits, timeouts, match rules, logging level, recorder, or traffic budget.
- Handler order matters. Place protocol decoders before handlers that inspect protocol objects, and place outbound encoders after handlers that write protocol objects.
- Handlers must keep backpressure and message ownership explicit; never retain ByteBuf values without clear lifetime control.
- TLS policy belongs in caller-supplied `crypto/tls.Config`; choose versions, cipher suites, certificates, SNI, and ALPN deliberately.

## Recommended Defaults

- Start with bounded sizes and short integration-test timeouts.
- Increase limits only after measuring realistic payloads and peer behavior.
- Keep security-sensitive defaults closed or conservative.
- Document every production override in the owning service, not in this library module.

## Environment Variables

This library module does not require repository-specific environment variables for normal unit tests. Applications, examples, benchmarks, and CI jobs may define their own variables around it.

## Local Workspace Development

Use a local `go.work` file only as a developer convenience. Published module metadata should remain portable and must not contain machine-specific paths.
