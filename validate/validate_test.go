package validate

import (
	"testing"
)

func TestCheck(t *testing.T) {

	ipv4False1 := IsValidIpV4Address("127")
	ipv4False2 := IsValidIpV4Address("127.0")
	ipv4False3 := IsValidIpV4Address("127.0.0")
	ipv4False4 := IsValidIpV4Address("127.0.0.")
	ipv4False5 := IsValidIpV4Address("127.0.0.1.1")
	ipv4False6 := IsValidIpV4Address("@")

	ipv6False1 := IsValidIpV4Address("0000:0000")
	ipv6False2 := IsValidIpV4Address("::::0001:::")
	ipv6False3 := IsValidIpV4Address("0000:0000:0000:0000:0000:0000:0000:000x")

	ipv4True1 := IsValidIpV4Address("127.0.0.1")
	ipv4True2 := IsValidIpV4Address("188.68.50.25")
	ipv4True3 := IsValidIpV4Address("255.255.255.255")

	ipv6True1 := IsValidIpV4Address("0000:0000:0000:0000:0000:0000:0000:0001")
	ipv6True2 := IsValidIpV4Address("::ffff:bc44:3219")
	ipv6True3 := IsValidIpV4Address("0:0:0:0:0:ffff:bc44:3219")
	ipv6True4 := IsValidIpV4Address("0000:0000:0000:0000:0000:ffff:bc44:3219")

	if ipv4False1 {
		t.Error("Check was incorrect for ipv4False1. Got: true, Expected: false")
	}

	if ipv4False2 {
		t.Error("Check was incorrect for ipv4False2. Got: true, Expected: false")
	}

	if ipv4False3 {
		t.Error("Check was incorrect for ipv4False3. Got: true, Expected: false")
	}

	if ipv4False4 {
		t.Error("Check was incorrect for ipv4False4. Got: true, Expected: false")
	}

	if ipv4False5 {
		t.Error("Check was incorrect for ipv4False5. Got: true, Expected: false")
	}

	if ipv4False6 {
		t.Error("Check was incorrect for ipv4False6. Got: true, Expected: false")
	}

	if ipv6False1 {
		t.Error("Check was incorrect for ipv4False1. Got: true, Expected: false")
	}

	if ipv6False2 {
		t.Error("Check was incorrect for ipv4False2. Got: true, Expected: false")
	}

	if ipv6False3 {
		t.Error("Check was incorrect for ipv4False3. Got: true, Expected: false")
	}

	if !ipv4True1 {
		t.Error("Check was incorrect for ipv4True1. Got: false, Expected: true")
	}

	if !ipv4True2 {
		t.Error("Check was incorrect for ipv4True2. Got: false, Expected: true")
	}

	if !ipv4True3 {
		t.Error("Check was incorrect for ipv4True3. Got: false, Expected: true")
	}

	if !ipv6True1 {
		t.Error("Check was incorrect for ipv6True1. Got: false, Expected: true")
	}

	if !ipv6True2 {
		t.Error("Check was incorrect for ipv6True2. Got: false, Expected: true")
	}

	if !ipv6True3 {
		t.Error("Check was incorrect for ipv6True3. Got: false, Expected: true")
	}

	if !ipv6True4 {
		t.Error("Check was incorrect for ipv6True4. Got: false, Expected: true")
	}
}
