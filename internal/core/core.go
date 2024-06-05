package core

type AppServer struct {
	InfoApp 		*InfoApp 		`json:"info_app"`
}

type InfoApp struct {
	AppName				string `json:"app_name"`
	AWSRegion			string `json:"aws_region"`
	ApiVersion			string `json:"version"`
	AvailabilityZone 	string `json:"availabilityZone"`
	BucketNameKey		string `json:"bucket_name_key"`
	FilePath			string `json:"file_path"`
	FileKey				string `json:"file_key"`
}

type RSA_Key struct  {
	RSAPrivateKey			string 	`json:"rsa_private_key"`
	RSAPrivateKeyEncoded	string 	`json:"rsa_private_encoded"`
}