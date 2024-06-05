package service

import (
	"time"
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"crypto/x509/pkix"
	"encoding/pem"
	"net"
	"math/big"
)

func(w WorkerService) GenerateX509Cert(	privkey *rsa.PrivateKey, 
										crt_serial int) (*x509.Certificate, 
														*[]byte,
														error){
	childLogger.Debug().Msg("GenerateX509Cert")

	crt := &x509.Certificate{
		SerialNumber: big.NewInt(int64(crt_serial)),
		Subject: pkix.Name{
			Organization:  []string{"DOCk"},
			Country:       []string{"BR"},
			Province:      []string{"SP"},
			Locality:      []string{"SP"},
			CommonName:    "localhost.com.br",
		},
		NotBefore:             time.Now(),
		NotAfter:              time.Now().AddDate(10, 0, 0), // Add 10 years
		IsCA:                  true,
		ExtKeyUsage:           []x509.ExtKeyUsage{x509.ExtKeyUsageClientAuth, x509.ExtKeyUsageServerAuth},
		KeyUsage:              x509.KeyUsageDigitalSignature | x509.KeyUsageCertSign,
		BasicConstraintsValid: true,
	}

	// Add the Subject Alternative Names
	var AlternativeNames [2]string
	AlternativeNames[0] = "localhost"
	AlternativeNames[1] = "go-api-global-mtls.architecture.caradhras.io"

    for _, altName := range AlternativeNames {
        crt.DNSNames = append(crt.DNSNames, altName)
    }

	var IPAddresses  [2]string
	IPAddresses [0] = "127.0.0.1"
	IPAddresses [1] = "0.0.0.0"

    for _, ipAddress := range IPAddresses  {
		ip := net.ParseIP(ipAddress)
        crt.IPAddresses  = append(crt.IPAddresses, ip)
    }

	// create the CA
	crtBytes, err := x509.CreateCertificate(rand.Reader, crt, crt, &privkey.PublicKey, privkey)
	if err != nil {
		return nil, nil, err
	}

	// pem encode
	crt_pem := pem.EncodeToMemory(
		&pem.Block{
				Type:  "CERTIFICATE",
				Bytes: crtBytes,
		},
	)

	return crt, &crt_pem, nil
}
