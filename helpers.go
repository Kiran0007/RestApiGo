package main

import(
  "crypto/sha1"
  "fmt"
  "crypto/rand"
)

const(
  PASSWORD_SALT = "12dfd6sf65ds4f65ds6f5d5f65ds"
)

func GetSHA1(s string) string {
	h := sha1.New()
	h.Write([]byte(s))
	bs := h.Sum(nil)
	return fmt.Sprintf("%x", bs)
}

func RandToken() string {
	b := make([]byte, 16)
	rand.Read(b)
	return fmt.Sprintf("%x", b)
}
