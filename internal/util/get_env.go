package util

import(
	"os"

	"github.com/rs/zerolog/log"
	"github.com/lambda-mtls-key/internal/core"
)

func GetAppInfo() core.AppServer {
	log.Debug().Msg("getEnv")

	var appServer	core.AppServer
	var infoApp		core.InfoApp

	if os.Getenv("APP_NAME") !=  "" {
		infoApp.AppName = os.Getenv("APP_NAME")
	}

	if os.Getenv("AWSRegion") !=  "" {
		infoApp.AWSRegion = os.Getenv("AWSRegion")
	}

	if os.Getenv("VERSION") !=  "" {
		infoApp.ApiVersion = os.Getenv("VERSION")
	}

	if os.Getenv("BUCKET_NAME_KEY") !=  "" {
		infoApp.BucketNameKey = os.Getenv("BUCKET_NAME_KEY")
	}

	if os.Getenv("FILE_PATH") !=  "" {
		infoApp.FilePath = os.Getenv("FILE_PATH")
	}

	if os.Getenv("FILE_KEY") !=  "" {
		infoApp.FileKey = os.Getenv("FILE_KEY")
	}

	appServer.InfoApp = &infoApp

	return appServer
}