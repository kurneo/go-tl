package error

const (
	DatasourceError = iota
	DomainError
	TransportError
)

type Contract interface {
	GetType() int
	GetError() error
	GetMessage() string
	IsTransportLevelErr() bool
	IsDomainError() bool
	IsDatasourceErr() bool
}
