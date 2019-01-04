package signature

import "testing"

// makes sure that text gets encoded properly.
// and that it can be verified at the other end.
func TestHmacAuthenticatorEncodeVerify(t *testing.T) {
	secret := []byte("greatsecret")

	auth := HmacAuthenticator{
		secret: secret,
	}

	text := "some text"
	hash := auth.Encode([]byte(text))

	if text == hash {
		t.Fatal("couldn't encode text properly")
	}

	res, err := auth.Verify([]byte("some text"), hash)
	if err != nil {
		t.Fatal(err)
	}
	if !res {
		t.Fatal("couldn't verify text")
	}
}
