package abuseIpDb

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"
	"time"
)

type Config struct {
	APIKey string `json:"apiKey"`
	URL    string `json:"url"`
}

type ApiResponse struct {
	Data struct {
		IPAddress            string        `json:"ipAddress"`
		IsPublic             bool          `json:"isPublic"`
		IPVersion            int           `json:"ipVersion"`
		IsWhitelisted        bool          `json:"isWhitelisted"`
		AbuseConfidenceScore int           `json:"abuseConfidenceScore"`
		CountryCode          string        `json:"countryCode"`
		UsageType            string        `json:"usageType"`
		Isp                  string        `json:"isp"`
		Domain               string        `json:"domain"`
		Hostnames            []interface{} `json:"hostnames"`
		TotalReports         int           `json:"totalReports"`
		NumDistinctUsers     int           `json:"numDistinctUsers"`
		LastReportedAt       time.Time     `json:"lastReportedAt"`
	} `json:"data"`
}

type ApiErrorResponse struct {
	Errors []struct {
		Detail string `json:"detail"`
		Status int    `json:"status"`
		Source struct {
			Parameter string `json:"parameter"`
		} `json:"source"`
	} `json:"errors"`
}

func GetIPInfo(config Config, httpClient *http.Client, ipAddr string) (ApiResponse, error) {

	// Set up request
	req, err := http.NewRequest(http.MethodGet, config.URL, nil)

	if err != nil {
		return ApiResponse{}, fmt.Errorf("getIPInfo: unable to create request: %v", err)
	}

	// Add url query with IP address
	q := req.URL.Query()
	q.Add("ipAddress", ipAddr)

	// Add header info
	req.Header.Add("Key", config.APIKey)

	req.URL.RawQuery = q.Encode()

	// Preform request
	resp, err := httpClient.Do(req)

	if err != nil {
		return ApiResponse{}, fmt.Errorf("getIPInfo: unable to do request: %v", err)
	}

	defer resp.Body.Close()

	// Read body
	respBytes, err := ioutil.ReadAll(resp.Body)

	if err != nil {
		return ApiResponse{}, fmt.Errorf("getIPInfo: unable to read response payload: %v", err)
	}

	// If the API returns an error, capture it
	if resp.StatusCode != http.StatusOK {
		var apiErrResp ApiErrorResponse

		err = json.Unmarshal(respBytes, &apiErrResp)

		if err != nil {
			return ApiResponse{}, fmt.Errorf("getIPInfo: non-200 response(%v). Unable to unmarshall error response: %v, ", resp.StatusCode, err)
		}

		var errStr []string
		for _, apiErr := range apiErrResp.Errors {
			errStr = append(errStr, apiErr.Detail)
		}

		return ApiResponse{}, fmt.Errorf("getIPInfo: non-200 response(%v). Error details: %v", resp.StatusCode, strings.Join(errStr, "; "))
	}

	var apiResp ApiResponse
	err = json.Unmarshal(respBytes, &apiResp)

	if err != nil {
		return ApiResponse{}, fmt.Errorf("getIPInfo: unable to unmarshall response: %v", err)
	}

	return apiResp, nil
}
