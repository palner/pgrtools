/*
 * Copyright (C) 2021	The Palner Group, Inc. (palner.com)
 *						Fred Posner (@fredposner)
 *
 * pgiptables is free software; you can redistribute it and/or modify
 * it under the terms of the GNU General Public License as published by
 * the Free Software Foundation; either version 2 of the License, or
 * (at your option) any later version
 *
 * pgiptables is distributed in the hope that it will be useful,
 * but WITHOUT ANY WARRANTY; without even the implied warranty of
 * MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
 * GNU General Public License for more details.
 *
 * You should have received a copy of the GNU General Public License
 * along with this program; if not, write to the Free Software
 * Foundation, Inc., 51 Franklin Street, Fifth Floor, Boston, MA  02110-1301  USA
 *
 */

package pgiptables

import (
	"errors"
	"fmt"
	"log"
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

	log.Print("IPTABLES doesn't contain APIBANLOCAL. Creating now...")

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
	log.Println("IPtableHandle:", proto, task, ipvar)

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
		log.Println("IPtableHandle:", err)
		return "", err
	}

	_, err = InitializeIPTables(ipt)
	if err != nil {
		log.Fatalln("IPtableHandler: failed to initialize IPTables:", err)
		return "", err
	}

	switch task {
	case "add":
		err = ipt.AppendUnique("filter", "APIBANLOCAL", "-s", ipvar, "-d", "0/0", "-j", targetChain)
		if err != nil {
			log.Println("IPtableHandler: error adding address", err)
			return "", err
		} else {
			return "added", nil
		}
	case "delete":
		err = ipt.DeleteIfExists("filter", "APIBANLOCAL", "-s", ipvar, "-d", "0/0", "-j", targetChain)
		if err != nil {
			log.Println("IPtableHandler: error removing address", err)
			return "", err
		} else {
			return "deleted", nil
		}
	case "flush":
		err = ipt.ClearChain("filter", "APIBANLOCAL")
		if err != nil {
			log.Println("IPtableHandler:", proto, err)
			return "", err
		} else {
			return "flushed", nil
		}
	default:
		log.Println("IPtableHandler: unknown task")
		return "", errors.New("unknown task")
	}
}
