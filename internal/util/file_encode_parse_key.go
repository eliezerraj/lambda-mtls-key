package util

import(
	"context"
    "os"
	"errors"
    "fmt"
	"encoding/pem"
	"crypto/rsa"
    "crypto/x509"
	"encoding/base64"

	"github.com/rs/zerolog/log"
	"github.com/lambda-mtls-key/internal/lib"
)

var childLogger = log.With().Str("internal", "utils").Logger()

// Aux function
func readFile(filePath string) ([]byte, error) {
	file, err := os.Open(filePath)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	fileInfo, err := file.Stat()
	if err != nil {
		return nil, err
	}

	data := make([]byte, fileInfo.Size())
	_, err = file.Read(data)
	if err != nil {
		return nil, err
	}

	return data, err
}

func LoadPemCert(filePath string) (*[]byte, error) {
    childLogger.Debug().Msg("LoadPemCert")

    pemBytes, err := readFile(filePath)
	if err != nil {
		return nil, err
	}

	fmt.Println(string(pemBytes))

	block, _ := pem.Decode(pemBytes)
	if block == nil {
		return nil, errors.New("failed to decode PEM-encoded certificate")
	}

	return &block.Bytes, nil
}

func SaveKeyAsFile(	key string,
					filename 	string) error {
	childLogger.Debug().Msg("SaveKeyAsFile")

	data := []byte(key)
	err := os.WriteFile(filename, data, 0644)
	if err != nil {
		return err
	}

    return nil
}

func ParsePEMToPrivateKey(ctx context.Context, pemString string) (*rsa.PrivateKey, error) {
    childLogger.Debug().Msg("ParsePEMToPrivateKey")

	span := lib.Span(ctx, "util.parsePEMToPrivateKey")	
    defer span.End()
	
	fmt.Println(pemString)

	block, _ := pem.Decode([]byte(pemString))
	if block == nil {
		return nil, errors.New("Failed to decode PEM-encoded key")
	}

	privateKey, err := x509.ParsePKCS1PrivateKey(block.Bytes)
    if err != nil {
        return nil, err
    }

	return privateKey, nil
}

func ParsePemToCertx509(ctx context.Context, pemString string) (*x509.Certificate, error) {
    childLogger.Debug().Msg("ParsePemToCertx509")
	
	span := lib.Span(ctx, "util.parsePemToCertx509")	
    defer span.End()

	fmt.Println(pemString)

	block, _ := pem.Decode([]byte(pemString))
	if block == nil || block.Type != "CERTIFICATE" {
		return nil, errors.New("Failed to decode PEM-encoded cert")
	}

	cert, err := x509.ParseCertificate(block.Bytes)
    if err != nil {
        return nil, err
    }

	return cert, nil
}

func DecodeB64(ctx context.Context, base64String string) (string, error){
    childLogger.Debug().Msg("DecodeB64")
	
	span := lib.Span(ctx, "util.decodeB64")	
    defer span.End()

    decodedBytes, err := base64.StdEncoding.DecodeString(base64String)
    if err != nil {
		return "", err
    }

	decodedString := string(decodedBytes)

	return decodedString, nil
}