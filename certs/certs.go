package main

import (
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"crypto/x509/pkix"
	"encoding/pem"
	"math/big"
	"net"
	"os"
	"time"
)

func main() {
	// TODO [grokrz]: refactor
	validFor := 365 * 24 * time.Hour
	generateDBCerts(validFor)
	generateClientCerts(validFor)
}

func generateDBCerts(validFor time.Duration) {
	// Generate a new RSA private key
	key, err := rsa.GenerateKey(rand.Reader, 2048)
	if err != nil {
		panic(err)
	}

	// Define the certificate template
	template := x509.Certificate{
		SerialNumber: big.NewInt(1),
		Subject: pkix.Name{
			CommonName: "db",
		},
		NotBefore: time.Now(),
		NotAfter:  time.Now().Add(validFor),
		KeyUsage:  x509.KeyUsageKeyEncipherment | x509.KeyUsageDigitalSignature,
		ExtKeyUsage: []x509.ExtKeyUsage{
			x509.ExtKeyUsageServerAuth,
		},
		BasicConstraintsValid: true,
	}

	template.IPAddresses = []net.IP{net.ParseIP("127.0.0.1")}
	template.DNSNames = []string{"localhost", "db"}

	// Create the certificate
	derBytes, err := x509.CreateCertificate(rand.Reader, &template, &template, &key.PublicKey, key)
	if err != nil {
		panic(err)
	}

	// Save the certificate to a file
	certOut, err := os.Create("server-ca.pem")
	if err != nil {
		panic(err)
	}
	defer certOut.Close()
	if err := pem.Encode(certOut, &pem.Block{Type: "CERTIFICATE", Bytes: derBytes}); err != nil {
		panic(err)
	}

	// Save the private key to a file
	keyOut, err := os.Create("server.key")
	if err != nil {
		panic(err)
	}
	defer keyOut.Close()
	if err := pem.Encode(keyOut, &pem.Block{Type: "RSA PRIVATE KEY", Bytes: x509.MarshalPKCS1PrivateKey(key)}); err != nil {
		panic(err)
	}

	// chmod 600 server.key
	if err := keyOut.Chmod(0600); err != nil {
		panic(err)
	}
}

func generateClientCerts(validFor time.Duration) {
	// Generate a new RSA private key
	key, err := rsa.GenerateKey(rand.Reader, 2048)
	if err != nil {
		panic(err)
	}

	// Encode the private key in PEM format and write to file
	keyOut, err := os.Create("client-key.pem")
	if err != nil {
		panic(err)
	}
	defer keyOut.Close()
	if err := pem.Encode(keyOut, &pem.Block{Type: "RSA PRIVATE KEY", Bytes: x509.MarshalPKCS1PrivateKey(key)}); err != nil {
		panic(err)
	}

	// Create a certificate template
	template := x509.Certificate{
		SerialNumber:          big.NewInt(1),
		Subject:               pkix.Name{CommonName: "flowdber"},
		NotBefore:             time.Now(),
		NotAfter:              time.Now().Add(validFor),
		KeyUsage:              x509.KeyUsageDigitalSignature,
		ExtKeyUsage:           []x509.ExtKeyUsage{x509.ExtKeyUsageClientAuth},
		BasicConstraintsValid: true,
	}

	// Generate a new self-signed certificate
	derBytes, err := x509.CreateCertificate(rand.Reader, &template, &template, &key.PublicKey, key)
	if err != nil {
		panic(err)
	}

	// Encode the certificate in PEM format and write to file
	certOut, err := os.Create("client-cert.pem")
	if err != nil {
		panic(err)
	}
	defer certOut.Close()
	if err := pem.Encode(certOut, &pem.Block{Type: "CERTIFICATE", Bytes: derBytes}); err != nil {
		panic(err)
	}
}
