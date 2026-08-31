# handler-tls

TLS pipeline handlers for Gnalloy backed by crypto/tls, including ALPN, SNI, StartTLS, optional TLS, OCSP, and cipher-suite helpers.

This repository is part of the Gnalloy modular networking stack. The default branch is `dev`; no release tag is created during bootstrap.

## Install

```bash
go get gnalloy.org/handler-tls@dev
```

## Module Boundary

- Module path: `gnalloy.org/handler-tls`
- Responsibility: TLS pipeline handlers for Gnalloy backed by crypto/tls, including ALPN, SNI, StartTLS, optional TLS, OCSP, and cipher-suite helpers
- Core dependency: `gnalloy.org/gnalloy` when this module uses Gnalloy buffers, channels, event loops, or bootstrap contracts.

## Gnalloy Dependencies

- `gnalloy.org/gnalloy`

## Development

```bash
go test ./... -count=1
go vet ./...
go test ./... -run '^$' -bench . -benchmem -benchtime=100ms -count=1
```

For multi-repository development, use the workspace at `G:\opensource\gnalloy\go.work`. For standalone verification, set `GOWORK=off`.

## License

Apache-2.0.
