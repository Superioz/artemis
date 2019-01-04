package rest

import "github.com/superioz/artemis/pkg/rest/signature"

type HttpRequest struct {
	route       string
	signature   signature.Info
	contentType string
	body        []byte
}

func (r *HttpRequest) Route() string {
	return r.route
}

func (r *HttpRequest) Signature() signature.Info {
	return r.signature
}

func (r *HttpRequest) ContentType() string {
	return r.contentType
}

func (r *HttpRequest) Body() []byte {
	return r.body
}
