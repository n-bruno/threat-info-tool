package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
	"github.com/google/logger"
)

// ServerResponse is used for server response payloads to present information to the enduser.
// It will present if there was an error with the request, a user message, and data of interest.
type ServerResponse struct {
	Success bool        `json:"success"`
	Message string      `json:"message"`
	Data    interface{} `json:"data"`
}

// FunctionReturn is a structure for function responses.
// It's very similar to ServerResponse with some extra information.
// Users don't need to see server error messages and HttpStatus is set in the Http header.
type FunctionReturn struct {
	IsError    bool        `json:"IsError"`
	HttpStatus int         `json:"httpStatus"`
	Error      error       `json:"error"`
	Message    string      `json:"message"`
	Data       interface{} `json:"data"`
}

func serveEndpoints() {
	router := chi.NewRouter()

	router.Use(middleware.Recoverer) // recovers from panics without crashing server
	router.Use(middleware.Logger)

	router.Group(func(r chi.Router) {
		r.Get("/", func(w http.ResponseWriter, r *http.Request) { mainWrapperFunc(w, r, homePage) })
		r.Get("/ipInfo", func(w http.ResponseWriter, r *http.Request) { mainWrapperFunc(w, r, ipData) })
	})

	logger.Info(fmt.Sprintf("Listening on %s:%s", config.Server.Host, config.Server.Port))

	if config.Server.SSL.Enable {
		log.Fatal(http.ListenAndServeTLS(
			config.Server.Host+":"+config.Server.Port,
			config.Server.SSL.CertFile, config.Server.SSL.KeyFile,
			router))
	} else {
		log.Fatal(http.ListenAndServe(
			config.Server.Host+":"+config.Server.Port, router))
	}
}

// mainHandler is a wrapper function for error handling and for server responses
func mainWrapperFunc(w http.ResponseWriter, req *http.Request, serverFunction func(req *http.Request) FunctionReturn) {
	funcResp := serverFunction(req)

	if funcResp.Error != nil {
		logErrorWithIp(funcResp.Error, req)
	}

	// Write status code
	w.WriteHeader(funcResp.HttpStatus)

	// Prepare response to client
	response := ServerResponse{
		Message: funcResp.Message,
		Success: !funcResp.IsError,
		Data:    funcResp.Data,
	}

	if err := json.NewEncoder(w).Encode(response); err != nil {
		logErrorWithIp(fmt.Errorf("writing response message failed: %v", err.Error()), req)
	}
}

// logErrorWithIp generates an error log with client's IP address
func logErrorWithIp(err error, req *http.Request) {
	logger.Error(fmt.Sprintf("Client Ip: %v; Error: %v", readUserIP(req), err.Error()))
}

// logErrorWithIp generates an info log with client's IP address
func logInfoWithIp(str string, req *http.Request) {
	logger.Info(fmt.Sprintf("Client Ip: %v; Message: %v", readUserIP(req), str))
}

// readUserIP extracts the source IP address from an http request packet
func readUserIP(r *http.Request) string {
	IPAddress := r.Header.Get("X-Real-Ip")
	if IPAddress == "" {
		IPAddress = r.Header.Get("X-Forwarded-For")
	}
	if IPAddress == "" {
		IPAddress = r.RemoteAddr
	}
	return IPAddress
}

func checkAuthorization(req *http.Request) error {
	providedApiKey := ""
	if _, ok := req.Header["Key"]; ok && len(req.Header["Key"]) == 1 {
		providedApiKey = req.Header["Key"][0]
	}

	if providedApiKey == "" {
		return fmt.Errorf("please provide your API key in \"Key\" in the request header")
	}

	authorized := false
	user := ""
	for key, val := range config.ApiKeys {
		if val == providedApiKey {
			user = key
			authorized = true
			break
		}
	}

	if !authorized {
		return fmt.Errorf("provided API key is invalid")
	}

	// Logs the API that's being used. Good for correlating where abusive behavior is coming from.
	// This is a bit too verbose for my liking because it will print for each request, but alas.
	logInfoWithIp("User authorized: "+user, req)

	return nil
}
