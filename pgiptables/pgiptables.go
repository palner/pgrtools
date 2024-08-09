/*

Copyright (C) 2021 Fred Posner. All Rights Reserved.
Copyright (C) 2021 The Palner Group, Inc. All Rights Reserved.

Permission is hereby granted, free of charge, to any person obtaining a copy
of this software and associated documentation files (the "Software"), to deal
in the Software without restriction, including without limitation the rights
to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
copies of the Software, and to permit persons to whom the Software is
furnished to do so, subject to the following conditions:

The above copyright notice and this permission notice shall be included in all
copies or substantial portions of the Software.

THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE
SOFTWARE.

*/

package pgiptables

import (
	"errors"
	"fmt"
	"net"

	"github.com/coreos/go-iptables/iptables"
)

const targetChain string = "REJECT" // REJECT or DROP

// Function to see if string within string
func Contains(list []string, value string) bool {
	for _, val := range list {
		if val == value {
			return true
		}
	}
	return false
}

func CheckIPAddress(ip string) bool {
	if net.ParseIP(ip) == nil {
		return false
	} else {
		return true
	}
}

func CheckIPAddressv4(ip string) (string, error) {
	if net.ParseIP(ip) == nil {
		return "", errors.New("Not an IP address")
	}
	for i := 0; i < len(ip); i++ {
		switch ip[i] {
		case '.':
			return "ipv4", nil
		case ':':
			return "ipv6", nil
		}
	}

	return "", errors.New("unknown error")
}

func InitializeIPTables(ipt *iptables.IPTables) (string, error) {
	// Get existing chains from IPTABLES
	originaListChain, err := ipt.ListChains("filter")
	if err != nil {
		return "error", fmt.Errorf("failed to read iptables: %w", err)
	}

	// Search for INPUT in IPTABLES
	chain := "INPUT"
	if !Contains(originaListChain, chain) {
		return "error", errors.New("iptables does not contain expected INPUT chain")
	}

	// Search for FORWARD in IPTABLES
	chain = "FORWARD"
	if !Contains(originaListChain, chain) {
		return "error", errors.New("iptables does not contain expected FORWARD chain")
	}

	// Search for APIBAN in IPTABLES
	chain = "APIBANLOCAL"
	if Contains(originaListChain, chain) {
		// APIBAN chain already exists
		return "chain exists", nil
	}

	// Add APIBAN chain
	err = ipt.ClearChain("filter", chain)
	if err != nil {
		return "error", fmt.Errorf("failed to clear APIBANLOCAL chain: %w", err)
	}

	// Add APIBAN chain to INPUT
	err = ipt.Insert("filter", "INPUT", 1, "-j", chain)
	if err != nil {
		return "error", fmt.Errorf("failed to add APIBANLOCAL chain to INPUT chain: %w", err)
	}

	// Add APIBAN chain to FORWARD
	err = ipt.Insert("filter", "FORWARD", 1, "-j", chain)
	if err != nil {
		return "error", fmt.Errorf("failed to add APIBANLOCAL chain to FORWARD chain: %w", err)
	}

	return "chain created", nil
}

func IPtableHandle(proto string, task string, ipvar string) (string, error) {
	var ipProto iptables.Protocol
	switch proto {
	case "ipv6":
		ipProto = iptables.ProtocolIPv6
	default:
		ipProto = iptables.ProtocolIPv4
	}

	// Go connect for IPTABLES
	ipt, err := iptables.NewWithProtocol(ipProto)
	if err != nil {
		return "", err
	}

	_, err = InitializeIPTables(ipt)
	if err != nil {
		return "", err
	}

	switch task {
	case "add":
		err = ipt.AppendUnique("filter", "APIBANLOCAL", "-s", ipvar, "-d", "0/0", "-j", targetChain)
		if err != nil {
			return "", err
		} else {
			return "added", nil
		}
	case "delete":
		err = ipt.DeleteIfExists("filter", "APIBANLOCAL", "-s", ipvar, "-d", "0/0", "-j", targetChain)
		if err != nil {
			return "", err
		} else {
			return "deleted", nil
		}
	case "flush":
		err = ipt.ClearChain("filter", "APIBANLOCAL")
		if err != nil {
			return "", err
		} else {
			return "flushed", nil
		}
	default:
		return "", errors.New("unknown task")
	}
}
