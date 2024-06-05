package handler

import(
	"github.com/rs/zerolog/log"
	"github.com/lambda-mtls-key/internal/service"
)

var childLogger = log.With().Str("handler", "LambdaHandler").Logger()

type LambdaHandler struct {
	workerService service.WorkerService
}

type MessageBody struct {
	ErrorMsg 	*string `json:"error,omitempty"`
	Msg 		*string `json:"message,omitempty"`
}

func NewLambdaHandler(workerService service.WorkerService) *LambdaHandler{
	childLogger.Debug().Msg("NewLambdaHandler")
	return &LambdaHandler{
		workerService: workerService,
	}
}

/*func (h *LambdaHandler) UnhandledMethod() (*events.APIGatewayProxyResponse, error){
	return ApiHandlerResponse(http.StatusMethodNotAllowed, MessageBody{ErrorMsg: aws.String(erro.ErrMethodNotAllowed.Error())})
}*/