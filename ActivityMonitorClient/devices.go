package main

import (
	"crypto/md5"
	"encoding/hex"
	"net"
)

func getMd5(in string) string {
	hashBytes := md5.Sum([]byte(in))
	return hex.EncodeToString(hashBytes[:])
}

func getHashedMacAddr() string {
	ifaces, _ := net.Interfaces()
	hash := ""
	for _, iface := range ifaces {
		if addr := iface.HardwareAddr.String(); addr != "" {
			hash = getMd5(addr)
			break
		}
	}
	return hash
}

func haveMacAddr(macHash string) bool {
	ifaces, _ := net.Interfaces()
	for _, iface := range ifaces {
		if addr := iface.HardwareAddr.String(); getMd5(addr) == macHash {
			return true
		}
	}
	return false
}
