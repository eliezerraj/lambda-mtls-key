package main

import(
	//"fmt"
	"context"

	"github.com/aws/aws-lambda-go/lambda"
	"github.com/aws/aws-lambda-go/events"

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
	workerHandler	*handler.WorkerHandler
	response		*events.APIGatewayProxyResponse
)

func init(){
	log.Debug().Msg("init")
	zerolog.SetGlobalLevel(logLevel)
	appServer = util.GetAppInfo()
}

func main(){
	log.Debug().Msg("--- main ---")

	ctx := context.Background()

	awsConfig, err := config.LoadDefaultConfig(ctx)
	if err != nil {
		panic("Error loading AWS configuration: " + err.Error())
	}

	workerService = service.NewWorkerService()
	workerHandler = handler.NewWorkerHandler(*workerService, appServer, awsConfig)

	// Start lambda handler
	lambda.Start(lambdaHandler)
}

func lambdaHandler(ctx context.Context, req events.APIGatewayProxyRequest) (*events.APIGatewayProxyResponse, error) {
	log.Debug().Msg("--- lambdaHandler 1.1---")

	// Check the http method and path
	switch req.HTTPMethod {
		case "GET":
			if (req.Resource == "/getKey/{id}"){  
				response, _ = workerHandler.UnhandledMethod()
			}else if (req.Resource == "/info"){
				response, _ = workerHandler.GetInfo()
			}else {
				response, _ = workerHandler.UnhandledMethod()
			}
		case "POST":
			if (req.Resource == "/createCRL"){  
				response, _ = workerHandler.CreateCRL(ctx, req)
			}else if (req.Resource == "/verifyCertCRL"){
				response, _ = workerHandler.VerifyCertCRL(ctx, req)
			}else {
				response, _ = workerHandler.UnhandledMethod()
			}
		case "DELETE":
				response, _ = workerHandler.UnhandledMethod()
		case "PUT":
				response, _ = workerHandler.UnhandledMethod()
		default:
				response, _ = workerHandler.UnhandledMethod()
	}
	
	return response, nil
}