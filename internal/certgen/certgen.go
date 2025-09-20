package certgen

import (
	"crypto/ed25519"
	"crypto/rand"
	"crypto/x509"
	"encoding/pem"
	"fmt"
	"math/big"
	"net"
	"os"
	"time"
)

func Generate(ip string, dns []string) error {
	if ip == "" {
		return fmt.Errorf("ip cannot be empty")
	}

	pubKey, pvtKey, err := ed25519.GenerateKey(rand.Reader)
	if err != nil {
		return err
	}

	cert := &x509.Certificate{
		SerialNumber: big.NewInt(1),

		PublicKey:          pubKey,
		PublicKeyAlgorithm: x509.Ed25519,

		DNSNames:    dns,
		IPAddresses: []net.IP{net.IP(ip)},

		NotBefore: time.Now(),
		NotAfter:  time.Now().Add(36 * time.Hour),

		KeyUsage:              x509.KeyUsageKeyEncipherment | x509.KeyUsageDigitalSignature,
		ExtKeyUsage:           []x509.ExtKeyUsage{x509.ExtKeyUsageServerAuth},
		BasicConstraintsValid: true,
	}

	certDER, err := x509.CreateCertificate(rand.Reader, cert, cert, pubKey, pvtKey)
	if err != nil {
		return err
	}

	// create cert file
	certF, err := os.Create("cert.pem")
	if err != nil {
		return err
	}

	certBlock := &pem.Block{Type: "CERTIFICATE", Bytes: certDER}
	pem.Encode(certF, certBlock)
	certF.Close()

	// create key file
	keyF, err := os.Create("key.pem")
	if err != nil {
		return err
	}

	privateKey, err := x509.MarshalPKCS8PrivateKey(pvtKey)
	if err != nil {
		return err
	}

	keyBlock := &pem.Block{Type: "PRIVATE KEY", Bytes: privateKey}
	pem.Encode(keyF, keyBlock)
	keyF.Close()

	return nil
}
