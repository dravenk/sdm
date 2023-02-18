package main

import (
	"crypto/md5"
	"encoding/base64"
	"fmt"
	"github.com/google/uuid"
)

func hashSalt() string {
	Uuid, _ := uuid.NewUUID()
	input := []byte(Uuid.String())
	baseSalt := base64.StdEncoding.EncodeToString(input)
	md5Sum := fmt.Sprintf("%x", md5.Sum(input))
	sha2 := base64.StdEncoding.EncodeToString([]byte(baseSalt + md5Sum))
	return sha2
}

func hashPass() string {
	Uuid, _ := uuid.NewUUID()
	input := []byte(Uuid.String())
	baseSalt := base64.StdEncoding.EncodeToString(input)
	return baseSalt
}
