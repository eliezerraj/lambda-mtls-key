package service

import (
    "context"
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"crypto/x509/pkix"
	"encoding/pem"
	"encoding/asn1"

    "github.com/lambda-mtls-key/internal/lib"
)

var oidEmailAddress = asn1.ObjectIdentifier{1, 2, 840, 113549, 1, 9, 1}

func(w WorkerService) GenerateCSRKey(ctx context.Context,
                                    privkey *rsa.PrivateKey) ( *x509.CertificateRequest, *[]byte, error){
	childLogger.Debug().Msg("GenerateCSRKey")
	
    span := lib.Span(ctx, "service.generateCSRKey")	
    defer span.End()

	emailAddress := "eliezer.junior@dock.tech"
    subj := pkix.Name{
        CommonName:         "localhost.com",
        Country:            []string{"BR"},
        Province:           []string{"BR"},
        Locality:           []string{"SP"},
        Organization:       []string{"DOCK"},
        OrganizationalUnit: []string{"ARCH"},
        ExtraNames: []pkix.AttributeTypeAndValue{
            {
                Type:  oidEmailAddress, 
                Value: asn1.RawValue{
                    Tag:   asn1.TagIA5String, 
                    Bytes: []byte(emailAddress),
                },
            },
        },
    }

	csr := x509.CertificateRequest{
        Subject:            subj,
        SignatureAlgorithm: x509.SHA256WithRSA,
        EmailAddresses: []string{emailAddress},
    }

	csrBytes, err := x509.CreateCertificateRequest(rand.Reader, &csr, privkey)
    if err != nil {
		return nil, nil, err
	}

	csr_pem := pem.EncodeToMemory(
            &pem.Block{
                    Type:  "CERTIFICATE REQUEST",
                    Bytes: csrBytes,
            },
    )

    return &csr, &csr_pem, nil
}
