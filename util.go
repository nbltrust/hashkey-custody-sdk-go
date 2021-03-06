package jadepoolsaas

const (
	defaultAddr = "https://openapi.jadepool.io"
)

type client interface {
	getKey() string
	getKeyHeaderName() string
	getSecret() string
	getAddr() string
}
