package util

import (
	"bytes"
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"crypto/x509/pkix"
	"encoding/pem"
	"math/big"
	"time"

	jose "gopkg.in/square/go-jose.v2"
)

// KeyPair define private/public key
type KeyPair struct {
	PrivateKey jose.JSONWebKey
	Cert       string
}

// GenJwk generates Jwk with rsa
func GenJwk(kid, cn string) (*KeyPair, error) {
	rsaKey, err := rsa.GenerateKey(rand.Reader, 3072)
	if err != nil {
		return nil, err
	}

	serialNumber, err := rand.Int(rand.Reader, big.NewInt(time.Now().Unix()))
	if err != nil {
		return nil, err
	}

	template := x509.Certificate{
		SerialNumber: serialNumber,
		Subject: pkix.Name{
			Country:      []string{"CN"},
			Organization: []string{"www.example.com"},
			CommonName:   cn,
		},
		Issuer: pkix.Name{
			Country:      []string{"CN"},
			Organization: []string{"www.example.com"},
			CommonName:   cn,
		},
		NotBefore:             time.Now().Add(-24 * time.Hour),
		NotAfter:              time.Now().Add(10 * 365 * 24 * time.Hour),
		SignatureAlgorithm:    x509.SHA256WithRSA,
		KeyUsage:              x509.KeyUsageKeyEncipherment | x509.KeyUsageDigitalSignature,
		ExtKeyUsage:           []x509.ExtKeyUsage{x509.ExtKeyUsageServerAuth},
		BasicConstraintsValid: true,
	}

	derBytes, err := x509.CreateCertificate(rand.Reader, &template, &template, &rsaKey.PublicKey, rsaKey)
	if err != nil {
		return nil, err
	}

	certificate, err := x509.ParseCertificate(derBytes)
	if err != nil {
		return nil, err
	}

	block := &pem.Block{
		Type:  "CERTIFICATE",
		Bytes: derBytes,
	}

	certPem := new(bytes.Buffer)
	if err := pem.Encode(certPem, block); err != nil {
		return nil, err
	}

	priv := jose.JSONWebKey{
		Key:          rsaKey,
		KeyID:        kid,
		Algorithm:    "RS256",
		Use:          "sig",
		Certificates: []*x509.Certificate{certificate},
	}

	return &KeyPair{
		PrivateKey: priv,
		Cert:       certPem.String(),
	}, nil
}

// ParseCertificates parse a pem-like cert string to x509.Certificates
func ParseCertificates(certStr string) ([]*x509.Certificate, error) {
	certBytes := []byte(certStr)
	var (
		certs []*x509.Certificate
		block *pem.Block
	)
	for {
		block, certBytes = pem.Decode(certBytes)
		if block == nil {
			break
		}
		cert, err := x509.ParseCertificate(block.Bytes)
		if err != nil {
			return nil, err
		}
		certs = append(certs, cert)
	}

	return certs, nil
}
