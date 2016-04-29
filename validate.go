package myapp

import (
	"crypto/hmac"
	"crypto/sha256"
	"fmt"
	"io"
	"strings"
)

func getCode(data string) string {
	h := hmac.New(sha256.New, []byte("key"))
	io.WriteString(h, data)
	return fmt.Sprintf("%x", h.Sum(nil))
}

func tampered(s string) bool {
	xs := strings.Split(s, "|")
	usrData := xs[1]
	usrCode := xs[2]
	return usrCode != getCode(usrData)
}
