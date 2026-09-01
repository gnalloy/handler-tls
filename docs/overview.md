# Overview

[简体中文](overview.zh-CN.md) | [Docs Index](README.md)

## Purpose

TLS pipeline handlers for Gnalloy backed by crypto/tls, including ALPN, SNI, StartTLS, optional TLS, OCSP, and cipher-suite helpers.

This module provides Pipeline handlers. A handler observes, transforms, rejects, delays, records, or protects messages after a Channel already exists. It does not own listening sockets or application protocols unless explicitly named.

## Repository Identity

- Module path: `gnalloy.org/handler-tls`
- GitHub repository: `github.com/gnalloy/handler-tls`
- Default branch: `dev`
- License: Apache-2.0

## Package Map
- `gnalloy.org/handler-tls` (`tls`)
- `gnalloy.org/handler-tls/provider/standard` (`standard`)

## Direct Gnalloy Dependencies
- `gnalloy.org/gnalloy`

## Direct Dependents in the Current Module Plan
- `gnalloy.org/benchmarks`

## Architecture Position

Gnalloy keeps the core small and dependency-light. This repository is a replaceable module around one responsibility, connected through explicit Go packages instead of runtime discovery.

## Compatibility

The public import path is `gnalloy.org/handler-tls`. Until the first stable tag is published, use `@dev` or an explicit pseudo-version selected by your dependency policy.
