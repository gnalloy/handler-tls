# API Reference

[简体中文](api.zh-CN.md) | [Docs Index](README.md)

This inventory is generated from `go doc -short` for the packages in this repository. It is a quick public-surface map; source files and tests remain the authority for exact semantics.

## Packages

### `gnalloy.org/handler-tls`

Package name: `tls`

```text
var ErrNeedInput = errors.New("gnalloy/handler/tls: need more encrypted input") ...
func CertificateWithOCSPStaple(cert cryptotls.Certificate, response []byte) (cryptotls.Certificate, error)
func CipherSuiteName(id uint16) string
func CipherSuiteNames(ids []uint16) []string
func ConfigureCipherSuiteNames(cfg *cryptotls.Config, value string, options CipherSuiteOptions) error
func ConfigureCipherSuites(cfg *cryptotls.Config, suites []uint16) error
func OpenSSLCipherSuiteName(id uint16) (string, bool)
func ParseCipherSuites(value string, options CipherSuiteOptions) ([]uint16, error)
func ServerConfigWithClientHelloProvider(base *cryptotls.Config, provider ClientHelloProvider) *cryptotls.Config
func ServerConfigWithSNI(base *cryptotls.Config, selector ServerNameSelector) *cryptotls.Config
func ValidateConfigurableCipherSuites(suites []uint16) error
type BytePool interface{ ... }
type BytePoolConfig struct{ ... }
type CipherSuiteCertificateAuth uint8
    const CipherSuiteCertificateAny CipherSuiteCertificateAuth = iota ...
type CipherSuiteInfo struct{ ... }
    func CipherSuiteCatalog(options CipherSuiteOptions) []CipherSuiteInfo
    func LookupCipherSuite(name string, options CipherSuiteOptions) (CipherSuiteInfo, error)
    func LookupCipherSuiteID(id uint16, options CipherSuiteOptions) (CipherSuiteInfo, error)
type CipherSuiteOptions struct{ ... }
type ClientHello struct{ ... }
    func InspectClientHello(raw []byte) (ClientHello, bool, error)
type ClientHelloProvider interface{ ... }
type ClientHelloProviderFunc func(hello ClientHello) (*cryptotls.Config, error)
type Config struct{ ... }
type Conn interface{ ... }
type CryptoProvider struct{}
type Handler struct{ ... }
    func Client(cfg Config) *Handler
    func Server(cfg Config) *Handler
type HandshakeEvent struct{ ... }
type Mode uint8
    const ModeClient Mode = iota + 1 ...
type NativeCapabilities struct{ ... }
type NativeEvaluation struct{ ... }
    func EvaluateNativeProvider(provider NativeProvider) NativeEvaluation
    func EvaluateProvider(provider Provider) NativeEvaluation
type NativeProvider interface{ ... }
type OCSPConfig struct{ ... }
type OCSPEvent struct{ ... }
type OCSPValidator interface{ ... }
type OCSPValidatorFunc func(state cryptotls.ConnectionState, response []byte) error
type OptionalEvent struct{ ... }
type OptionalHandler struct{ ... }
    func NewOptionalHandler() *OptionalHandler
type PooledBytePool struct{ ... }
    func NewPooledBytePool(cfg BytePoolConfig) *PooledBytePool
type Provider interface{ ... }
type ServerNameSelector func(serverName string) (*cryptotls.Config, error)
    func ServerConfigMap(configs map[string]*cryptotls.Config) ServerNameSelector
type StartEvent struct{}
type UnsupportedNativeProvider struct{}
```

### `gnalloy.org/handler-tls/provider/standard`

Package name: `standard`

```text
type Provider struct{ ... }
    func Default() Provider
    func New(name ...string) Provider
```
