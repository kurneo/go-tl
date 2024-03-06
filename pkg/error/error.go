package error

type err struct {
	err     error
	errType int
}

func (e err) GetType() int {
	return e.errType
}

func (e err) GetError() error {
	return e.err
}

func (e err) GetMessage() string {
	return e.err.Error()
}

func (e err) IsTransportLevelErr() bool {
	return e.errType == TransportError
}

func (e err) IsDomainError() bool {
	return e.errType == DomainError
}

func (e err) IsDatasourceErr() bool {
	return e.errType == DatasourceError
}

func NewDatasource(e error) Contract {
	return &err{
		err:     e,
		errType: DatasourceError,
	}
}

func NewDomain(e error) Contract {
	return &err{
		err:     e,
		errType: DomainError,
	}
}

func NewTransport(e error) Contract {
	return &err{
		err:     e,
		errType: TransportError,
	}
}
