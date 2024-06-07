package handler

import(
	"net/http"
	"encoding/json"
	"context"

	"github.com/rs/zerolog/log"
	"github.com/lambda-mtls-key/internal/core"
	"github.com/lambda-mtls-key/internal/erro"
	"github.com/lambda-mtls-key/internal/service"
	"github.com/lambda-mtls-key/internal/util"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-lambda-go/events"
)

var (
	childLogger = log.With().Str("handler", "LambdaHandler").Logger()
	crl_ca_name = "crl_ca.pem"
)

type WorkerHandler struct {
	workerService 	service.WorkerService
	appServer		core.AppServer
	awsConfig		aws.Config
}

type MessageBody struct {
	ErrorMsg 	*string `json:"error,omitempty"`
	Msg 		*string `json:"message,omitempty"`
}

func NewWorkerHandler(	workerService service.WorkerService, 
						appServer	core.AppServer,
						awsConfig	aws.Config) *WorkerHandler{
	childLogger.Debug().Msg("NewWorkerHandler")
	return &WorkerHandler{
		workerService: workerService,
		appServer:		appServer,
		awsConfig: 		awsConfig,
	}
}

func ApiHandlerResponse(statusCode int, body interface{}) (*events.APIGatewayProxyResponse, error){
	stringBody, err := json.Marshal(&body)
	if err != nil {
		return nil, erro.ErrUnmarshal
	}

	return &events.APIGatewayProxyResponse{
		StatusCode: statusCode,
		Headers: map[string]string{
			"Content-Type": "application/json",
		},
		Body: string(stringBody),
	}, nil
}

func (h *WorkerHandler) UnhandledMethod() (*events.APIGatewayProxyResponse, error){
	return ApiHandlerResponse(http.StatusMethodNotAllowed, MessageBody{ErrorMsg: aws.String(erro.ErrMethodNotAllowed.Error())})
}

func (h *WorkerHandler) GetInfo() (*events.APIGatewayProxyResponse, error) {
	childLogger.Debug().Msg("GetInfo")

	handlerResponse, err := ApiHandlerResponse(http.StatusOK, h.appServer)
	if err != nil {
		return ApiHandlerResponse(http.StatusInternalServerError, MessageBody{ErrorMsg: aws.String(err.Error())})
	}

	return handlerResponse, nil
}

func (h *WorkerHandler) CreateCRL(ctx context.Context, req events.APIGatewayProxyRequest) (*events.APIGatewayProxyResponse, error) {
	childLogger.Debug().Msg("CreateCRL")

	// Retrieve data from request
	var key_pem core.Key_Pem
    if err := json.Unmarshal([]byte(req.Body), &key_pem); err != nil {
        return ApiHandlerResponse(http.StatusBadRequest, MessageBody{ErrorMsg: aws.String(err.Error())})
    }

	// Decode to Base64
	RSAPrivateKeyPemDecoded, err := util.DecodeB64(key_pem.RSAPrivateKeyPem)
	if err != nil {
		return ApiHandlerResponse(http.StatusBadRequest, MessageBody{ErrorMsg: aws.String(err.Error())})
	}
	CertX509PemDecoded, err := util.DecodeB64(key_pem.CertX509Pem)
	if err != nil {
		return ApiHandlerResponse(http.StatusBadRequest, MessageBody{ErrorMsg: aws.String(err.Error())})
	}

	// Parse the key pem
	privateKey, err := util.ParsePEMToPrivateKey(RSAPrivateKeyPemDecoded)
	if err != nil {
		return ApiHandlerResponse(http.StatusBadRequest, MessageBody{ErrorMsg: aws.String(err.Error())})
	}
	certX509, err := util.ParsePemToCertx509(CertX509PemDecoded)
	if err != nil {
		return ApiHandlerResponse(http.StatusBadRequest, MessageBody{ErrorMsg: aws.String(err.Error())})
	}

	// Create a CRL
	_, crl_pem ,err := h.workerService.CreateCRL(privateKey,certX509)
	if err != nil {
		return ApiHandlerResponse(http.StatusBadRequest, MessageBody{ErrorMsg: aws.String(err.Error())})
	}

	// Save the CRL
	err = util.SaveKeyAsFileS3(	ctx, 
								h.awsConfig, 
								h.appServer.InfoApp.BucketNameKey,
								h.appServer.InfoApp.FilePath,
								crl_ca_name,
								*crl_pem)
	if err != nil {
		return ApiHandlerResponse(http.StatusBadRequest, MessageBody{ErrorMsg: aws.String(err.Error())})
	}

	msg_response := "CRL created Successful !!!"
	response := MessageBody{ Msg: &msg_response }

	handlerResponse, err := ApiHandlerResponse(http.StatusOK, response)
	if err != nil {
		return ApiHandlerResponse(http.StatusInternalServerError, MessageBody{ErrorMsg: aws.String(err.Error())})
	}

	return handlerResponse, nil
}

func (h *WorkerHandler) VerifyCertCRL(ctx context.Context, req events.APIGatewayProxyRequest) (*events.APIGatewayProxyResponse, error) {
	childLogger.Debug().Msg("VerifyCertCRL")

	// Retrieve data from request
	var key_pem core.Key_Pem
    if err := json.Unmarshal([]byte(req.Body), &key_pem); err != nil {
        return ApiHandlerResponse(http.StatusBadRequest, MessageBody{ErrorMsg: aws.String(err.Error())})
    }

	// Decode to Base64
	CertX509PemDecoded, err := util.DecodeB64(key_pem.CertX509Pem)
	if err != nil {
		return ApiHandlerResponse(http.StatusBadRequest, MessageBody{ErrorMsg: aws.String(err.Error())})
	}

	// Parse the key pem
	certX509, err := util.ParsePemToCertx509(CertX509PemDecoded)
	if err != nil {
		return ApiHandlerResponse(http.StatusBadRequest, MessageBody{ErrorMsg: aws.String(err.Error())})
	}

	// Load the CRL
	load_crl_pem, err := util.LoadKeyAsFileS3(ctx, 
											h.awsConfig, 
											h.appServer.InfoApp.BucketNameKey,
											h.appServer.InfoApp.FilePath,
											crl_ca_name)
	if err != nil {
		return ApiHandlerResponse(http.StatusBadRequest, MessageBody{ErrorMsg: aws.String(err.Error())})
	}

	response, err := h.workerService.VerifyCertCRL(*load_crl_pem, certX509)
	if err != nil {
		return ApiHandlerResponse(http.StatusBadRequest, MessageBody{ErrorMsg: aws.String(err.Error())})
	}

	handlerResponse, err := ApiHandlerResponse(http.StatusOK, response)
	if err != nil {
		return ApiHandlerResponse(http.StatusInternalServerError, MessageBody{ErrorMsg: aws.String(err.Error())})
	}

	return handlerResponse, nil
}