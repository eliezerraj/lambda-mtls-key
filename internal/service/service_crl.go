package service

import(
	"fmt"
	"time"
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"crypto/x509/pkix"
	"encoding/pem"
)

func(w WorkerService) CreateCRL(privkey *rsa.PrivateKey, 
								cacert *x509.Certificate) (	*pkix.CertificateList, *[]byte, error){
	childLogger.Debug().Msg("CreateCRL")

	revokedCerts := []pkix.RevokedCertificate{
        {
            SerialNumber:   cacert.SerialNumber,
            RevocationTime: time.Now(),
        },
    }

	now := time.Now()
	crlBytes, err := cacert.CreateCRL(rand.Reader, privkey, revokedCerts, now, now.Add(365*24*time.Hour))

    if err != nil {
		return nil, nil, err
    }

	crl_pem := pem.EncodeToMemory(
		&pem.Block{
			Type: "X509 CRL", 
			Bytes: crlBytes,
		},
	)
	
	res, _ := x509.ParseDERCRL(crl_pem)
	return res, &crl_pem, nil
}

func(w WorkerService) VerifyCertCRL(crl []byte, 
									cacert *x509.Certificate) (bool, error){
	childLogger.Debug().Msg("VerifyCertCRL")

	certSerialNumber := cacert.SerialNumber
	fmt.Println(certSerialNumber)

	_crl, err := x509.ParseCRL(crl)
	if err != nil {
		return false, err
	}

	for _, revokedCert := range _crl.TBSCertList.RevokedCertificates {
		if revokedCert.SerialNumber.Cmp(certSerialNumber) == 0 {
			return true, nil
		}
	}

	fmt.Println(cacert.SerialNumber)
	return false, nil
}