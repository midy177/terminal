package utils

import "testing"

func TestName(t *testing.T) {
	enc, err := AesEncryptByGCM([]byte("hello"))
	if err != nil {
		t.Fatal(err)
	}
	dec, err := AesDecryptByGCM(enc)
	if err != nil {
		t.Fatal(err)
	}
	t.Log(string(dec))
}
