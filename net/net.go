package net

import (
	base "net"
	"net/url"
	"sync"

	"github.com/huyungtang/go-lib/slices"
)

// constants & variables ******************************************************************************************************************
// ****************************************************************************************************************************************
// ****************************************************************************************************************************************

var (
	privateIPBlocks     []*base.IPNet
	privateIPBlocksOnce sync.Once
)

// public functions ***********************************************************************************************************************
// ****************************************************************************************************************************************
// ****************************************************************************************************************************************

// IsPrivate
// ****************************************************************************************************************************************
func IsPrivate(addr string) bool {
	privateIPBlocksOnce.Do(func() {
		ips := []string{
			"127.0.0.0/8",
			"10.0.0.0/8",
			"172.16.0.0/12",
			"192.168.0.0/16",
			"::1/128",
			"fe80::/10",
		}

		privateIPBlocks = make([]*base.IPNet, 0, len(ips))
		for _, cidr := range ips {
			_, block, _ := base.ParseCIDR(cidr)
			privateIPBlocks = append(privateIPBlocks, block)
		}
	})

	ip := base.ParseIP(addr)
	_, isExists := slices.IndexOf(len(privateIPBlocks), func(i int) bool {
		return privateIPBlocks[i].Contains(ip)
	})

	return isExists
}

// ParseURL
// ****************************************************************************************************************************************
func ParseURL(uri string, params map[string]string) string {
	req, _ := url.Parse(uri)
	qry := req.Query()
	for key, val := range params {
		if key != "" {
			qry.Add(key, val)
		}
	}
	req.RawQuery = qry.Encode()

	return req.String()
}

// type defineds **************************************************************************************************************************
// ****************************************************************************************************************************************
// ****************************************************************************************************************************************

// private functions **********************************************************************************************************************
// ****************************************************************************************************************************************
// ****************************************************************************************************************************************
