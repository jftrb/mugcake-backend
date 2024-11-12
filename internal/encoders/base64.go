package encoders

import (
	"bytes"
	"encoding/base64"
)

func EncodeToBase64(s string) string {
	byteVal := bytes.NewBufferString(s).Bytes()
	return base64.StdEncoding.EncodeToString(byteVal)
}

func DecodeBase64(s string) (string, error) {
	bytes, err := base64.StdEncoding.DecodeString(s)
	return string(bytes[:]), err
}