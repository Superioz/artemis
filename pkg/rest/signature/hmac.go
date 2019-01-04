package signature

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/base64"
)

type HmacAuthenticator struct {
	secret []byte
}

func (h *HmacAuthenticator) Encode(body []byte) string {
	return base64.StdEncoding.EncodeToString(h.encode(body))
}

func (h *HmacAuthenticator) Verify(body []byte, signature string) (bool, error) {
	hash := h.encode(body)

	dec, err := base64.StdEncoding.DecodeString(signature)
	if err != nil {
		return false, err
	}
	return hmac.Equal(hash, dec), nil
}

func (h *HmacAuthenticator) encode(body []byte) []byte {
	hash := hmac.New(sha256.New, h.secret)
	hash.Write(body)

	return hash.Sum(nil)
}
