package util

import(
    "os"
	"errors"
    "fmt"
	"encoding/pem"
	"github.com/rs/zerolog/log"
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
