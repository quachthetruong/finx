package jwttoken

import (
	"crypto/rsa"
	"crypto/x509"
	"encoding/base64"
	"errors"
)

func GetPublicKey(base64Key string) (*rsa.PublicKey, error) {
	key, err := base64.StdEncoding.DecodeString(base64Key)
	if err != nil {
		return nil, err
	}

	var parsedKey interface{}
	if parsedKey, err = x509.ParsePKIXPublicKey(key); err != nil {
		if cert, err := x509.ParseCertificate(key); err == nil {
			parsedKey = cert.PublicKey
		} else {
			if parsedKey, err = x509.ParsePKCS1PublicKey(key); err != nil {
				return nil, err
			}
		}
	}

	var pkey *rsa.PublicKey
	var ok bool
	if pkey, ok = parsedKey.(*rsa.PublicKey); !ok {
		return nil, errors.New("key is not a valid RSA public key")
	}

	return pkey, nil
}
