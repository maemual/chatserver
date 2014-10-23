package main

import (
	"crypto/md5"
	"crypto/sha1"
	"fmt"

	"code.google.com/p/go-uuid/uuid"
)

func NewUUID() string {
	return uuid.New()
}

func PasswordHash(password string) string {
	return fmt.Sprintf("%x", sha1.Sum([]byte(fmt.Sprintf("%x", md5.Sum([]byte(password))))))
}
