package net

import (
	"testing"
)

// constants & variables ******************************************************************************************************************
// ****************************************************************************************************************************************
// ****************************************************************************************************************************************

// public functions ***********************************************************************************************************************
// ****************************************************************************************************************************************
// ****************************************************************************************************************************************

// TestIsPrivate
// ****************************************************************************************************************************************
func TestIsPrivate(t *testing.T) {
	ips := []string{
		"192.168.0.32",
		"::1",
		"127.0.0.1",
		"10.37.25.44",
		"65.33.184.42",
	}

	for i, ip := range ips {
		t.Log(ip, " -> ", IsPrivate(ip))
		if (i == 4 && IsPrivate(ip)) ||
			(i != 4 && !IsPrivate(ip)) {
			t.Fail()
		}
	}
}

// type defineds **************************************************************************************************************************
// ****************************************************************************************************************************************
// ****************************************************************************************************************************************

// private functions **********************************************************************************************************************
// ****************************************************************************************************************************************
// ****************************************************************************************************************************************
