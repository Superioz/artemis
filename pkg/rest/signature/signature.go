package signature

type Info struct {
	client    string
	algorithm string
	signature string
}

func (i *Info) Client() string {
	return i.client
}

func (i *Info) Algorithm() string {
	return i.algorithm
}

func (i *Info) Signature() string {
	return i.signature
}

type Authenticator interface {
	Encode(body []byte) string
	Verify(body []byte, signature string) bool
}
