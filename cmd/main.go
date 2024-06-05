package main

import(
	"fmt"
	"context"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"

	"github.com/lambda-mtls-key/internal/handler"
	"github.com/lambda-mtls-key/internal/util"
	"github.com/lambda-mtls-key/internal/service"
	"github.com/lambda-mtls-key/internal/core"
	"github.com/aws/aws-sdk-go-v2/config"
)

var (
	logLevel		=	zerolog.DebugLevel // InfoLevel DebugLevel
	appServer		core.AppServer
	workerService	*service.WorkerService
	lambdaHandler	*handler.LambdaHandler
)

func init(){
	log.Debug().Msg("init")
	zerolog.SetGlobalLevel(logLevel)
	
	appServer = util.GetAppInfo()
}

func main(){
	log.Debug().Msg("main - lambda-mtls-key")

	ctx := context.Background()

	awsConfig, err := config.LoadDefaultConfig(ctx)
	if err != nil {
		panic("Error loading AWS configuration: " + err.Error())
	}

	workerService = service.NewWorkerService()
	lambdaHandler = handler.NewLambdaHandler(*workerService)

	//-------------
	privateKey, _, _, _, err := workerService.GenerateRsaKeyPair()
	crt_serial := 12345
	x509cert, _, err := workerService.GenerateX509Cert(privateKey, crt_serial)
	_, crl_pem, err := workerService.CreateCRL(privateKey, x509cert)
	appServer.InfoApp.BucketNameKey = "eliezerraj-908671954593-mtls-truststore"
	appServer.InfoApp.FilePath = "/"
	appServer.InfoApp.FileKey = "crl_ca.pem"
	//-------------

	err = util.SaveKeyAsFileS3(ctx, 
								awsConfig, 
								appServer.InfoApp.BucketNameKey,
								appServer.InfoApp.FilePath,
								appServer.InfoApp.FileKey,
								*crl_pem	)
	
	bodyBytes, err := util.LoadKeyAsFileS3(ctx, 
								awsConfig, 
								appServer.InfoApp.BucketNameKey,
								appServer.InfoApp.FilePath,
								appServer.InfoApp.FileKey	)

	fmt.Println(string(*bodyBytes))
}

//------
var (
	version			=	"lambda-mtls-key"

	server_pubkey_name		= 	"./certs/mtls-api-go-global-pub.key"
	server_privkey_name		= 	"./certs/mtls-api-go-global.key"
	server_crt_name			= 	"./certs/mtls-api-go-global.crt"

	client_pubkey_name		= "./certs/client-01-pub.key"
	client_privkey_name		= "./certs/client-01-priv.key"
	client_csr_name			= "./certs/client-01.csr"
	client_ca_name			= "./certs/ca-client-01.crt"

	crl_ca_name				= "./certs/crl_ca.pem"
)

func test(){
	log.Debug().Msg("test - lambda-mtls-key")

	WorkerService := service.NewWorkerService()

	fmt.Println("------------- SERVER : KEY----------------")
	// Create the RSA keys
	privateKey, _, privatePem, publicPem, err := WorkerService.GenerateRsaKeyPair()
	if err != nil {
		panic(err)
	}
	// Encode the RSA keys
	//privkey_pem := util.ExportRsaPrivateKeyAsPemStr(privkey)
	//privkey_pem_encoded := WorkerService.EncodedRsaPrivateKeyAsPemStr(privkey_pem)

	//var rsa_key core.RSA_Key
	//rsa_key.RSAPrivateKey = privkey_pem
	//rsa_key.RSAPrivateKeyEncoded = privkey_pem_encoded

	//fmt.Println(rsa_key)

	// Save as a file
	//fmt.Println(string(*publicPem))
	//fmt.Println(string(*privatePem))

	err = util.SaveKeyAsFile(string(*publicPem), 
							server_pubkey_name)
		
	err = util.SaveKeyAsFile(string(*privatePem), 
							server_privkey_name)
	if err != nil {
		panic(err)
	}
	fmt.Println("----------- SERVER : KEY=> CRT  ---------------")
	// Create the CSR keys
	crt_serial := 12345
	x509cert, crt_pem, err := WorkerService.GenerateX509Cert(privateKey, crt_serial)
	if err != nil {
		panic(err)
	}
	// Save as a file
	err = util.SaveKeyAsFile(string(*crt_pem), 
							server_crt_name)
	if err != nil {
		panic(err)
	}
	fmt.Println("------------- CLIENT : KEY ----------------")
	privateKeyCli, _, privateCliPem, publicCliPem, err := WorkerService.GenerateRsaKeyPair()
	if err != nil {
		panic(err)
	}

	err = util.SaveKeyAsFile(string(*publicCliPem), 
							client_pubkey_name)
	err = util.SaveKeyAsFile(string(*privateCliPem), 
							client_privkey_name)
	if err != nil {
		panic(err)
	}

	fmt.Println("------------- CLIENT : CSR ----------------")
	// Create the CSR keys
	_, csr_pem, err := WorkerService.GenerateCSRKey(privateKeyCli)
	//fmt.Println(csr_pem)
	// Save as a file
	err = util.SaveKeyAsFile(string(*csr_pem), client_csr_name)
	if err != nil {
		panic(err)
	}
	fmt.Println("-------------CREATE CRL----------------")
	_, crl_pem, err := WorkerService.CreateCRL(privateKeyCli, x509cert)
	if err != nil {
		panic(err)
	}
	err = util.SaveKeyAsFile(string(*crl_pem), crl_ca_name)
	if err != nil {
		panic(err)
	}
	fmt.Println("-------------VERIFY CERT/CRL----------------")
	crl_ver, err := WorkerService.VerifyCertCRL(*crl_pem, x509cert)
	fmt.Println(crl_ver)

	fmt.Println("----------------------------------")
	fmt.Println("-------------LOAD----------------")
	fmt.Println("----------------------------------")

	fmt.Println("-------------LOAD CERT/CRL----------------")

	fmt.Println(crl_ca_name)
	load_crl_pem, err := util.LoadPemCert(crl_ca_name)
	if err != nil {
		panic(err)
	}

	load_crl_ver, err := WorkerService.VerifyCertCRL(*load_crl_pem, x509cert)
	fmt.Println(load_crl_ver)
}