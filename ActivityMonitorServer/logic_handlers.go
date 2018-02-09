package main

import (
	"fmt"
	"log"
	"net/http"
	"time"
)

// /check?name=Name1&physaddr=MacAddr
// renew device access time
func checkCommunicationHandler(source map[string]*Device) http.HandlerFunc {
	return func(respWriter http.ResponseWriter, request *http.Request) {
		var name = request.FormValue("name")
		var mac = request.FormValue("physaddr")
		if mac == "" || name == "" {
			fmt.Fprint(respWriter, "There is missed arguments.")
			return
		}
		devicesMutex.Lock()
		dev, ok := source[mac]
		if ok {
			dev.LastResponse = time.Now().Unix()
			if !dev.IsActive {
				dev.IsActive = true
				notification := fmt.Sprintf("Device %s with name %s now online\n", dev.PhysicalAddress, dev.Name)
				log.Printf(notification)
				botNotify(notification)
			}
			fmt.Fprint(respWriter, "Device's access time updated.")
		} else {
			source[mac] = &Device{
				Name:            name,
				PhysicalAddress: mac,
				LastResponse:    time.Now().Unix(),
				IsActive:        true,
			}
			fmt.Fprint(respWriter, "New device added!")
		}
		devicesMutex.Unlock()
	}
}

// /ticktime
// return duration in which performs checking clients
func tickTimeHandler(duration string) http.HandlerFunc {
	return func(respWriter http.ResponseWriter, request *http.Request) {
		respWriter.Write([]byte(duration))
	}
}

// /
//default listeners
func notFoundHandler(respWriter http.ResponseWriter, request *http.Request) {
	respWriter.WriteHeader(404)
	fmt.Fprint(respWriter, "404 Page Not Found")
}
