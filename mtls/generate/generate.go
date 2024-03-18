package generate

import (
	"bytes"
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"crypto/x509/pkix"
	"encoding/pem"
	"fmt"
	"math/big"
	"net"
	"net/netip"
	"os"
	"path/filepath"
	"time"
)

func MakeCA(subject *pkix.Name, path string) (*x509.Certificate, *rsa.PrivateKey, error) {
	// creating a CA which will be used to sign all of our certificates using the x509 package from the Go Standard Library
	caCert := &x509.Certificate{
		SerialNumber:          big.NewInt(2019),
		Subject:               *subject,
		NotBefore:             time.Now(),
		NotAfter:              time.Now().AddDate(10*365, 0, 0),
		IsCA:                  true, // <- indicating this certificate is a CA certificate.
		ExtKeyUsage:           []x509.ExtKeyUsage{x509.ExtKeyUsageClientAuth, x509.ExtKeyUsageServerAuth},
		KeyUsage:              x509.KeyUsageDigitalSignature | x509.KeyUsageCertSign,
		BasicConstraintsValid: true,
	}

	// generate a private key for the CA
	caKey, err := rsa.GenerateKey(rand.Reader, 2048)
	if err != nil {
		return nil, nil, fmt.Errorf("generating the CA Private key: %w", err)
	}

	// create the CA certificate
	caBytes, err := x509.CreateCertificate(rand.Reader, caCert, caCert, &caKey.PublicKey, caKey)
	if err != nil {
		return nil, nil, fmt.Errorf("creating the CA Certificate: %w", err)
	}

	// Create the CA PEM files
	caPEM := new(bytes.Buffer)
	pem.Encode(caPEM, &pem.Block{Type: "CERTIFICATE", Bytes: caBytes})

	p := filepath.Join(path, "ca.crt")
	fmt.Printf("dumping %q\n", p)

	if err := os.WriteFile(p, caPEM.Bytes(), 0644); err != nil {
		return nil, nil, fmt.Errorf("writing the CA certificate file: %w", err)
	}

	caPrivKeyPEM := new(bytes.Buffer)
	pem.Encode(caPrivKeyPEM, &pem.Block{
		Type: "RSA PRIVATE KEY", Bytes: x509.MarshalPKCS1PrivateKey(caKey),
	})

	p = filepath.Join(path, "ca.key")
	fmt.Printf("dumping %q\n", p)

	if err := os.WriteFile(p, caPrivKeyPEM.Bytes(), 0644); err != nil {
		return nil, nil, fmt.Errorf("writing the CA certificate file: %w", err)
	}

	return caCert, caKey, nil
}

func MakeCert(caCert *x509.Certificate, caKey *rsa.PrivateKey,
	subject *pkix.Name, name, ip, path string,
) error {
	ipAddr, err := netip.ParseAddr(ip)
	if err != nil {
		return fmt.Errorf("parsing ip address: %w", err)
	}

	cert := &x509.Certificate{
		SerialNumber: big.NewInt(1658),
		Subject:      *subject,
		IPAddresses:  []net.IP{net.IPv4(127, 0, 0, 1), net.IPv6loopback, ipAddr.AsSlice()},

		DNSNames:     []string{"localhost"},
		NotBefore:    time.Now(),
		NotAfter:     time.Now().AddDate(10, 0, 0),
		SubjectKeyId: []byte{1, 2, 3, 4, 6},
		ExtKeyUsage:  []x509.ExtKeyUsage{x509.ExtKeyUsageClientAuth, x509.ExtKeyUsageServerAuth},
		KeyUsage:     x509.KeyUsageDigitalSignature,
	}

	certKey, err := rsa.GenerateKey(rand.Reader, 4096)
	if err != nil {
		return fmt.Errorf("generating the Key: %w", err)
	}

	certBytes, err := x509.CreateCertificate(rand.Reader, cert, caCert, &certKey.PublicKey, caKey)
	if err != nil {
		return fmt.Errorf("generate the certificate: %w", err)
	}

	certPEM := new(bytes.Buffer)
	pem.Encode(certPEM, &pem.Block{Type: "CERTIFICATE", Bytes: certBytes})

	p := filepath.Join(path, name+".crt")
	fmt.Printf("dumping %q\n", p)

	if err := os.WriteFile(p, certPEM.Bytes(), 0644); err != nil {
		return fmt.Errorf("writing the CA certificate file: %w", err)
	}

	certKeyPEM := new(bytes.Buffer)
	pem.Encode(certKeyPEM, &pem.Block{
		Type: "RSA PRIVATE KEY", Bytes: x509.MarshalPKCS1PrivateKey(certKey),
	})

	p = filepath.Join(path, name+".key")
	fmt.Printf("dumping %q\n", p)

	if err := os.WriteFile(p, certKeyPEM.Bytes(), 0644); err != nil {
		return fmt.Errorf("writing the CA certificate file: %w", err)
	}

	return nil
}
