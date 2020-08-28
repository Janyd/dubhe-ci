package utils

import (
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"encoding/pem"
	"errors"
	"golang.org/x/crypto/ssh"
)

func GenRsaKey(bits int) (string, string, error) {
	// 生成私钥文件
	privateKey, err := rsa.GenerateKey(rand.Reader, bits)

	if err != nil {
		return "", "", err
	}
	derStream := x509.MarshalPKCS1PrivateKey(privateKey)
	priBlock := &pem.Block{
		Type:  "RSA PRIVATE KEY",
		Bytes: derStream,
	}

	// 生成公钥文件
	publicKey := &privateKey.PublicKey
	derPkix, err := x509.MarshalPKIXPublicKey(publicKey)
	if err != nil {
		return "", "", err
	}
	publicBlock := &pem.Block{
		Type:  "PUBLIC KEY",
		Bytes: derPkix,
	}

	block, _ := pem.Decode(pem.EncodeToMemory(publicBlock))
	keyInterface, err := x509.ParsePKIXPublicKey(block.Bytes)
	key, ok := keyInterface.(*rsa.PublicKey)
	if !ok {
		return "", "", errors.New("not RSA public key")
	}
	skey, err := ssh.NewPublicKey(key)
	b := ssh.MarshalAuthorizedKey(skey)
	return string(pem.EncodeToMemory(priBlock)), string(b), nil
}
