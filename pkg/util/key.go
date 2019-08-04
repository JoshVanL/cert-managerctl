package util

import (
	"crypto"
	"crypto/ecdsa"
	"crypto/rsa"
	"crypto/x509"
	"encoding/pem"
	"errors"
	"fmt"
	"io/ioutil"
)

type KeyBundle struct {
	PrivateKey         crypto.Signer
	SignatureAlgorithm x509.SignatureAlgorithm
	PublicKeyAlgorithm x509.PublicKeyAlgorithm
}

func ParsePrivateKeyFile(path string) (*KeyBundle, error) {
	keyBytes, err := ioutil.ReadFile(path)
	if err != nil {
		return nil, err
	}

	return DecodePrivateKeyBytes(keyBytes)
}

// DecodePrivateKeyBytes will decode a PEM encoded private key into a crypto.Signer.
// It supports ECDSA and RSA private keys only. All other types will return err.
func DecodePrivateKeyBytes(keyBytes []byte) (*KeyBundle, error) {
	// decode the private key pem
	block, _ := pem.Decode(keyBytes)
	if block == nil {
		return nil, errors.New("error decoding private key PEM block")
	}

	var err error
	var key interface{}
	var sigAlgo x509.SignatureAlgorithm
	var pubAlgo x509.PublicKeyAlgorithm

	switch block.Type {
	case "PRIVATE KEY":
		key, err = x509.ParsePKCS8PrivateKey(block.Bytes)
		if err != nil {
			return nil, fmt.Errorf("error parsing pkcs#8 private key: %s", err)
		}

		_, ok := key.(*rsa.PrivateKey)
		if !ok {
			_, ok = key.(*ecdsa.PrivateKey)
			if !ok {
				return nil, errors.New("error determining private key type")
			}

			sigAlgo = x509.ECDSAWithSHA256
			pubAlgo = x509.ECDSA

			break
		}

		sigAlgo = x509.SHA256WithRSA
		pubAlgo = x509.RSA

	case "EC PRIVATE KEY":
		key, err = x509.ParseECPrivateKey(block.Bytes)
		if err != nil {
			return nil, fmt.Errorf("error parsing ecdsa private key: %s", err)
		}

		sigAlgo = x509.ECDSAWithSHA256
		pubAlgo = x509.ECDSA

		break

	case "RSA PRIVATE KEY":
		key, err = x509.ParsePKCS1PrivateKey(block.Bytes)
		if err != nil {
			return nil, fmt.Errorf("error parsing rsa private key: %s", err)
		}

		sigAlgo = x509.SHA256WithRSA
		pubAlgo = x509.RSA

	default:
		return nil, fmt.Errorf("unknown private key type: %s", block.Type)
	}

	signer, ok := key.(crypto.Signer)
	if !ok {
		return nil, errors.New("error parsing pkcs#8 private key: invalid key type")
	}

	return &KeyBundle{
		PrivateKey:         signer,
		SignatureAlgorithm: sigAlgo,
		PublicKeyAlgorithm: pubAlgo,
	}, nil
}
