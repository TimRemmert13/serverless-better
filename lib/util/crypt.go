package util

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/base64"
	"os"
)

func GenerateSecretHash(username string) string {
	tobehashed := username + os.Getenv("AWS_CLIENT_ID")
	h := hmac.New(sha256.New, []byte("AWS_CLIENT_SECRET"))
	h.Write([]byte(tobehashed))
	sha := base64.StdEncoding.EncodeToString(h.Sum(nil))
	return sha
}
