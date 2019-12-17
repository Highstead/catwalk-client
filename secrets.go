package client

import (
	"encoding/json"
	"io/ioutil"
	"os"
	"path/filepath"

	"github.com/pkg/errors"
)

type Secrets struct {
}

var (
	tsEndpoint   = ""
	authEndpoint = ""
	secretToken  = ""
)

func init() {
	tsEndpoint = os.Getenv("TS_ENDPOINT")
	authEndpoint = os.Getenv("AUTH_ENDPOINT")
	secretToken = os.Getenv("SECRET_TOKEN")
}

func ParseSecretsFile(installDir string) (*Secrets, error) {
	path := filepath.Join(installDir, "config", "secrets.json")
	data, err := ioutil.ReadFile(path)
	if err != nil {
		return nil, errors.Wrapf(err, "failed to read file: %s", path)
	}

	return ParseSecrets(path, data)
}

func ParseSecrets(path string, data []byte) (*Secrets, error) {
	s := &Secrets{}
	if err := json.Unmarshal(data, s); err != nil {
		return nil, errors.Wrapf(err, "failed to unmarshal json, path: %s", path)
	}

	return s, nil
}
