package main

import (
	"crypto"
	"crypto/rand"
	"crypto/rsa"
	"crypto/sha256"
	"crypto/x509"
	"encoding/base64"
	"encoding/pem"
	"fmt"
	"io/ioutil"
	"os"
)

func main() {
	args := os.Args[1:]
	if len(args) != 1 {
		fmt.Println("Usage: signtool filename")
		os.Exit(1)
	}
	key := os.Getenv("KEY")
	if key == "" {
		panic("Env KEY is not set")
	}
	rng := rand.Reader
	binary, err := ioutil.ReadFile(args[0])
	if err != nil {
		panic(err)
	}
	keyData, _ := base64.URLEncoding.DecodeString(key)
	block, _ := pem.Decode(keyData)
	rsaPrivateKey, err := x509.ParsePKCS8PrivateKey(block.Bytes)
	if err != nil {
		panic(err)
	}
	hashed := sha256.Sum256(binary)
	signature, err := rsa.SignPKCS1v15(rng, rsaPrivateKey.(*rsa.PrivateKey), crypto.SHA256, hashed[:])
	if err != nil {
		panic(err)
	}
	sha256 := fmt.Sprintf("%x", hashed)
	sigfile := fmt.Sprintf("%s.sig", args[0])
	sha256file := fmt.Sprintf("%s.sha256", args[0])
	err = ioutil.WriteFile(sigfile, signature, 0644)
	if err != nil {
		panic(err)
	}
	err = ioutil.WriteFile(sha256file, []byte(sha256), 0644)
	if err != nil {
		panic(err)
	}
}
