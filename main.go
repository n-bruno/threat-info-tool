package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strings"
	"time"

	"threatInfoTool/abuseIpDb"
	"threatInfoTool/validate"

	"github.com/google/logger"
)

type Config struct {
	Enable  bool   `json:"enable"`
	LogFile string `json:"logFile"`
	Server  struct {
		Host string `json:"host"`
		Port string `json:"port"`
		SSL  struct {
			Enable   bool   `json:"enable"`
			CertFile string `json:"certFile"`
			KeyFile  string `json:"keyFile"`
		} `json:"ssl"`
	} `json:"server"`
	HTTP struct {
		Timeout int `json:"timeout"`
	} `json:"http"`
	ApiKeys   map[string]string `json:"apiKeys"`
	AbuseIPDb abuseIpDb.Config  `json:"abuseIpDb"`
}

var (
	configFilePath string
	verbose        bool
	config         Config
	httpClient     *http.Client
)

func main() {
	flag.StringVar(&configFilePath, "c", "", "Config file path")
	flag.BoolVar(&verbose, "v", false, "Print info level logs to stdout")
	flag.BoolVar(&verbose, "h", false, "Print usage info")
	flag.Parse()

	if len(os.Args) == 0 || configFilePath == "" {
		printUsage()
	}

	// Read config file
	var err error
	config, err = readConfig(configFilePath)

	if err != nil {
		log.Fatalf("Failed to open config file: %v", err)
	}

	if err := validateCredentials(config); err != nil {
		log.Fatalf(err.Error())
	}

	// Set up logging
	lf, err := os.OpenFile(config.LogFile, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0660)
	if err != nil {
		log.Fatalf("Failed to open log file: %v", err)
	}

	defer lf.Close()
	defer logger.Init("ThreatInfoTool", verbose, true, lf).Close()

	// Set up http client
	httpClient = &http.Client{Timeout: time.Duration(config.HTTP.Timeout) * time.Second}

	logger.Info("Program launched")

	// Set up endpoints
	serveEndpoints()
}

func readConfig(filePath string) (Config, error) {

	file, err := ioutil.ReadFile(filePath)

	if err != nil {
		return Config{}, fmt.Errorf("readConfig: unable to read config file: %v", err)
	}

	var c Config
	err = json.Unmarshal(file, &c)

	if err != nil {
		return Config{}, fmt.Errorf("readConfig: unable to unmarshall config file: %v", err)
	}

	return c, nil
}

// validateCredentials looks through all the credentials specified in the config file and
// verifies the items pass our validation checks.
// API keys are validated for strength and usernames are validates for looking normal.
func validateCredentials(config Config) error {
	var errStr []string

	for username, apiKey := range config.ApiKeys {
		if !validate.IsValidUserName(username) {
			errStr = append(errStr, fmt.Sprintf("Not a valid username: \"%v\". "+
				"Usernames must be alphabetical characters and be 4 to 16 characters.", username))
		}

		if !validate.IsValidUserName(apiKey) {
			errStr = append(errStr, fmt.Sprintf("Not a valid API key: \"%v\". "+
				"API keys must be at least 64 characters and contain alphabetical and numbers", apiKey))
		}
	}

	if len(config.ApiKeys) == 0 {
		errStr = append(errStr, "Please provide at least one API key in the config file under \"apiKeys\". "+
			"The key should be a username and the value should be a strong API key.")
	}

	if len(errStr) != 0 {
		return fmt.Errorf(strings.Join(errStr, "; "))
	}

	return nil
}

func printUsage() {
	fmt.Println("Threat Info Tool v1.0")
	fmt.Println("Usage: ./threatInfoTool -c /configFilePath.json -v")
	flag.PrintDefaults()
	os.Exit(0)
}
