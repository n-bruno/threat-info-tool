package main

import (
	"fmt"
	"net/http"

	"threatInfoTool/abuseIpDb"
	"threatInfoTool/validate"
)

type formattedIpInfo struct {
	IsPublic          bool          `json:"isPublic"`
	IPVersion         int           `json:"ipVersion"`
	CountryCode       string        `json:"countryCode"`
	UsageType         string        `json:"usageType"`
	Isp               string        `json:"isp"`
	Domain            string        `json:"domain"`
	PotentialForAbuse int           `json:"potentialForAbuse"`
	Hostnames         []interface{} `json:"hostnames"`
}

func homePage(_ *http.Request) FunctionReturn {
	return FunctionReturn{
		IsError:    false,
		HttpStatus: http.StatusOK,
		Message:    "Please make requests to the \"/ipInfo\" endpoint for ip info.",
	}
}

func ipData(req *http.Request) FunctionReturn {
	if err := checkAuthorization(req); err != nil {
		return FunctionReturn{
			IsError:    true,
			HttpStatus: http.StatusUnauthorized,
			Message:    err.Error(),
			Error:      err,
		}
	}

	ipAddrParam := req.URL.Query()["ipAddress"]

	if len(ipAddrParam) != 1 {
		return FunctionReturn{
			IsError:    true,
			HttpStatus: http.StatusInternalServerError,
			Message:    "Please specify one \"ipAddress\" query.",
		}
	}

	ipAddress := ipAddrParam[0]

	if !validate.IsValidIpV4Address(ipAddress) {
		return FunctionReturn{
			IsError:    true,
			HttpStatus: http.StatusBadRequest,
			Message:    "Supplied IPv4 / IPv6 address is invalid.",
		}
	}

	logInfoWithIp("User requested info on the following IP: "+ipAddress, req)

	aResp, err := abuseIpDb.GetIPInfo(config.AbuseIPDb, httpClient, ipAddress)

	// If the request fails, log the error and hide the technical details from the user.
	if err != nil {
		return FunctionReturn{
			IsError:    true,
			HttpStatus: http.StatusInternalServerError,
			Error:      fmt.Errorf("ipData: %v", err),
		}
	}

	ipInfo := formattedIpInfo{
		IsPublic:          aResp.Data.IsPublic,
		IPVersion:         aResp.Data.IPVersion,
		CountryCode:       aResp.Data.CountryCode,
		UsageType:         aResp.Data.UsageType,
		Isp:               aResp.Data.Isp,
		Domain:            aResp.Data.Domain,
		Hostnames:         aResp.Data.Hostnames,
		PotentialForAbuse: aResp.Data.AbuseConfidenceScore,
	}

	return FunctionReturn{
		IsError:    false,
		HttpStatus: http.StatusOK,
		Data:       ipInfo,
	}
}
