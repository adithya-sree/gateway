package file

import (
	"encoding/json"
	"github.com/adithya-sree/gateway/config"
	"github.com/adithya-sree/logger"
	"io/ioutil"
	"os"
)

// File Logger
var out = logger.GetLogger(config.LogFile, "file-config")

// Service Definition
type Definition struct {
	ServiceName  string        `json:"service_name"`
	ServiceRoute string        `json:"service_route"`
	Secret       string        `json:"secret"`
	Enabled      bool          `json:"enabled"`
	Config       BackendConfig `json:"config"`
}

// Service Backend Configuration
type BackendConfig struct {
	BackendTimeout int    `json:"backend_timeout"`
	BackendFQDN    string `json:"backend_fqdn"`
}

// Reads all definitions from backend template file
func Read() ([]Definition, error) {
	// Get Home Directory
	home, err := os.UserHomeDir()
	if err != nil {
		return nil, err
	}
	// Open Configuration File
	file, err := os.Open(home + config.BackendTemplateFile)
	if err != nil {
		return nil, err
	}
	// Defer File Close
	defer func() {
		err = file.Close()
		if err != nil {
			out.Errorf("Unable to close file input stream [%v]", err)
		}
	}()
	// Read all file contents as bytes
	bytes, err := ioutil.ReadAll(file)
	if err != nil {
		return nil, err
	}
	// Marshal all definitions
	var definitions []Definition
	err = json.Unmarshal(bytes, &definitions)
	if err != nil {
		return nil, err
	}
	// Return Definitions
	return definitions, nil
}
