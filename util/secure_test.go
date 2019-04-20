package util_test

import (
	"encoding/json"
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/westwin/gocommon/util"
)

func TestGenJwk(t *testing.T) {
	jwk, err := util.GenJwk("clientid", "theoneid.io")
	assert.NoError(t, err)

	keyStr, err := json.Marshal(jwk)

	fmt.Printf("Json format of key: \n%s\n", keyStr)
	jwkPub := jwk.PrivateKey.Public()

	keyPubStr, err := json.Marshal(jwkPub)

	fmt.Printf("Json format of public key: \n%s\n", keyPubStr)
	fmt.Printf("Cert: %s\n", jwk.Cert)
}

func TestParseCertificates(t *testing.T) {
	cn := "example.com"
	jwk, _ := util.GenJwk("kid", cn)
	certStr := jwk.Cert
	certs, err := util.ParseCertificates(certStr)
	assert.NoError(t, err)
	assert.Equal(t, cn, certs[0].Subject.CommonName)
}
