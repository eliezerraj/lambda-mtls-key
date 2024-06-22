package main

import(
	"context"

	"github.com/aws/aws-lambda-go/lambda"
	"github.com/aws/aws-lambda-go/events"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"

	"github.com/lambda-mtls-key/internal/handler"
	"github.com/lambda-mtls-key/internal/util"
	"github.com/lambda-mtls-key/internal/service"
	"github.com/lambda-mtls-key/internal/core"
	"github.com/lambda-mtls-key/internal/lib"
	"github.com/aws/aws-sdk-go-v2/config"

	"go.opentelemetry.io/contrib/propagators/aws/xray"
	"go.opentelemetry.io/otel"
 	"go.opentelemetry.io/contrib/instrumentation/github.com/aws/aws-lambda-go/otellambda"
	"go.opentelemetry.io/contrib/instrumentation/github.com/aws/aws-sdk-go-v2/otelaws"
	"go.opentelemetry.io/contrib/instrumentation/github.com/aws/aws-lambda-go/otellambda/xrayconfig"
	"go.opentelemetry.io/otel/trace"
)

var (
	logLevel		=	zerolog.DebugLevel // InfoLevel DebugLevel
	appServer		core.AppServer
	workerService	*service.WorkerService
	workerHandler	*handler.WorkerHandler
	response		*events.APIGatewayProxyResponse
	tracer 			trace.Tracer
)

func init(){
	log.Debug().Msg("init")
	zerolog.SetGlobalLevel(logLevel)
	appServer = util.GetAppInfo()
	configOTEL := util.GetOtelEnv()
	appServer.ConfigOTEL = &configOTEL
}

func main(){
	log.Debug().Msg("--- main ---")
	log.Debug().Interface("appServer :",appServer).Msg("")
	
	ctx := context.Background()
	awsConfig, err := config.LoadDefaultConfig(ctx)
	if err != nil {
		panic("Error loading AWS configuration: " + err.Error())
	}

	// Instrument all AWS clients.
	otelaws.AppendMiddlewares(&awsConfig.APIOptions)

	workerService = service.NewWorkerService()
	workerHandler = handler.NewWorkerHandler(*workerService, appServer, awsConfig)

		//----- OTEL ----//
	tp := lib.NewTracerProvider(ctx, appServer.ConfigOTEL, appServer.InfoApp)
	defer func(ctx context.Context) {
			err := tp.Shutdown(ctx)
			if err != nil {
				log.Error().Err(err).Msg("Error shutting down tracer provider")
			}
	}(ctx)
	
	otel.SetTextMapPropagator(xray.Propagator{})
	otel.SetTracerProvider(tp)
	
	tracer = tp.Tracer("lambda-mtls-key-tracer")
	lambda.Start(otellambda.InstrumentHandler(lambdaHandler, xrayconfig.WithRecommendedOptions(tp)... ))

	// Start lambda handler
	//lambda.Start(lambdaHandler)
}

func lambdaHandler(ctx context.Context, req events.APIGatewayProxyRequest) (*events.APIGatewayProxyResponse, error) {
	log.Debug().Msg("--- lambdaHandler ---")

	ctx, span := tracer.Start(ctx, "lambdaHandler_otel_v1.2")
    defer span.End()
	
	// Check the http method and path
	switch req.HTTPMethod {
		case "GET":
			if (req.Resource == "/info"){  
				response, _ = workerHandler.GetInfo(ctx)
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