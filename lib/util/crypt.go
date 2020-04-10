package util

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/base64"
	"os"
)

func GenerateSecretHash(username string) string {
	tobehashed := username + os.Getenv("AWS_CLIENT_ID")
	secret := os.Getenv("AWS_CLIENT_SECRET")
	h := hmac.New(sha256.New, []byte(secret))
	h.Write([]byte(tobehashed))
	sha := base64.StdEncoding.EncodeToString(h.Sum(nil))
	return sha
}
