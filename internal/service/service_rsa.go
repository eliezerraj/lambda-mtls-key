package service

import (
	"context"
	"crypto/rand"
	"crypto/rsa"
	"encoding/pem"
	"crypto/x509"

	"github.com/rs/zerolog/log"
	"github.com/lambda-mtls-key/internal/lib"
)

var childLogger = log.With().Str("service", "service").Logger()
var size = 4096

type WorkerService struct {
}

func NewWorkerService() *WorkerService{
	childLogger.Debug().Msg("NewWorkerService")

	return &WorkerService{
	}
}

func(w WorkerService) GenerateRsaKeyPair(ctx context.Context) (*rsa.PrivateKey, 
																*rsa.PublicKey,
																*[]byte,
																*[]byte,
																error){
	childLogger.Debug().Msg("GenerateRsaKeyPair (private + public)")
    
	span := lib.Span(ctx, "service.generateRsaKeyPair")	
    defer span.End()

	// Generate the private key
	privateKey, err := rsa.GenerateKey(rand.Reader, size)
	if err != nil {
		return nil, nil, nil, nil, err
	}

	// Encode the private
	privatePem:= pem.EncodeToMemory(
		&pem.Block{
			Type: "RSA PRIVATE KEY", 
			Bytes: x509.MarshalPKCS1PrivateKey(privateKey),
		},
	)

	// Extract the public Key
	publicKey := privateKey.PublicKey
	// Encode the private
	publicKey_bytes, err := x509.MarshalPKIXPublicKey(&publicKey)
	if err != nil {
		return nil, nil, nil, nil, err
	}
	publicPem := pem.EncodeToMemory(
		&pem.Block{
			Type: "RSA PUBLIC KEY", 
			Bytes: publicKey_bytes,
		},
	)

	return privateKey, &publicKey, &privatePem, &publicPem, nil
}
