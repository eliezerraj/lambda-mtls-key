package core

type AppServer struct {
	InfoApp 		*InfoApp 		`json:"info_app"`
	ConfigOTEL		*ConfigOTEL		`json:"otel_config"`
}

type InfoApp struct {
	AppName				string `json:"app_name"`
	AWSRegion			string `json:"aws_region"`
	ApiVersion			string `json:"version"`
	AvailabilityZone 	string `json:"availabilityZone"`
	BucketNameKey		string `json:"bucket_name_key"`
	FilePath			string `json:"file_path"`
	FileKey				string `json:"file_key"`
	Env					string `json:"env,omitempty"`
	AccountID			string `json:"account,omitempty"`
}

type Key_Pem struct  {
	CertX509Pem			string 	`json:"cert_x509_pem,omitempty"`
	RSAPrivateKeyPem	string 	`json:"rsa_private_key_pem,omitempty"`
	CrlCaPem			string 	`json:"crl_ca_pem,omitempty"`
}

type ConfigOTEL struct {
	OtelExportEndpoint		string
	TimeInterval            int64    `mapstructure:"TimeInterval"`
	TimeAliveIncrementer    int64    `mapstructure:"RandomTimeAliveIncrementer"`
	TotalHeapSizeUpperBound int64    `mapstructure:"RandomTotalHeapSizeUpperBound"`
	ThreadsActiveUpperBound int64    `mapstructure:"RandomThreadsActiveUpperBound"`
	CpuUsageUpperBound      int64    `mapstructure:"RandomCpuUsageUpperBound"`
	SampleAppPorts          []string `mapstructure:"SampleAppPorts"`
}