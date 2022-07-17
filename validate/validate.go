package validate

import (
	"net"
	"regexp"
)

var regexIsValidUsername = regexp.MustCompile("[A-Za-z]{4,16}")
var regexIsStrongAPIKey = regexp.MustCompile("[A-Za-z0-9]{64,}")

func IsValidIpV4Address(ipStr string) bool {
	return net.ParseIP(ipStr) != nil
}

func IsValidUserName(s string) bool {
	return regexIsValidUsername.MatchString(s)
}

func IsStrongAPIKey(s string) bool {
	return regexIsStrongAPIKey.MatchString(s)
}
