package util

import (
	"crypto"
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

func ParsePrivateKeyFile(path string) (crypto.Signer, x509.SignatureAlgorithm, error) {
	keyBytes, err := ioutil.ReadFile(path)
	if err != nil {
		return nil, -1, err
	}

	return DecodePrivateKeyBytes(keyBytes)
}

// DecodePrivateKeyBytes will decode a PEM encoded private key into a crypto.Signer.
// It supports ECDSA and RSA private keys only. All other types will return err.
func DecodePrivateKeyBytes(keyBytes []byte) (crypto.Signer, x509.SignatureAlgorithm, error) {
	// decode the private key pem
	block, _ := pem.Decode(keyBytes)
	if block == nil {
		return nil, -1, errors.New("error decoding private key PEM block")
	}

	// If the file just has `PRIVATE KEY` then need to check the types (*rsa.PrivateKey)

	return nil, -1, nil

	//switch block.Type {
	//case "PRIVATE KEY":
	//	key, err := x509.ParsePKCS8PrivateKey(block.Bytes)
	//	if err != nil {
	//		return nil, -1, fmt.Errorf("error parsing pkcs#8 private key: %s", err)
	//	}

	//	signer, ok := key.(crypto.Signer)
	//	if !ok {
	//		return nil, -1, errors.New("error parsing pkcs#8 private key: invalid key type")
	//	}
	//	return signer, -1, nil
	//	//	case "EC PRIVATE KEY":
	//	//		key, err := x509.ParseECPrivateKey(block.Bytes)
	//	//		if err != nil {
	//	//			return nil, fmt.Errorf("error parsing ecdsa private key: %s", err)
	//	//		}
	//	//
	//	//		return key, nil
	//	//	case "RSA PRIVATE KEY":
	//	//		key, err := x509.ParsePKCS1PrivateKey(block.Bytes)
	//	//		if err != nil {
	//	//			return nil, fmt.Errorf("error parsing rsa private key: %s", err)
	//	//		}
	//	//
	//	//		err = key.Validate()
	//	//		if err != nil {
	//	//			return nil, fmt.Errorf("rsa private key failed validation: %s", err)
	//	//		}
	//	//		return key, nil
	//default:
	//	return nil, -1, fmt.Errorf("unknown private key type: %s", block.Type)
	//}
}
