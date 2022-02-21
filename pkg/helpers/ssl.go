package helpers

import (
	"crypto/rsa"
	"crypto/x509"
	"encoding/pem"
	"io/ioutil"
)

func LoadRSAKey(path string) *rsa.PrivateKey {
	data, err := ioutil.ReadFile(path)
	if err != nil {
		panic(err)
	}
	block, _ := pem.Decode(data)
	key, err := x509.ParsePKCS8PrivateKey(block.Bytes)
	if err != nil {
		panic(err)
	}
	return key.(*rsa.PrivateKey)
}

func LoadCertificate(path string) *x509.Certificate {
	data, err := ioutil.ReadFile(path)
	if err != nil {
		panic(err)
	}
	block, _ := pem.Decode(data)
	cert, err := x509.ParseCertificate(block.Bytes)
	if err != nil {
		panic(err)
	}
	return cert
}
