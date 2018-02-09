package main

import (
	"encoding/json"
	"net/http"
)

type StatusResponse struct {
	OK      bool
	Message string
}

// /api/devicestatus
// return list of registered device in JSON
func deviceStatusHandler(devices map[string]*Device) http.HandlerFunc {
	return func(respWriter http.ResponseWriter, request *http.Request) {
		respWriter.Header().Add("Content-Type", "application/json")
		bytes, err := json.Marshal(&devices)
		if err != nil {
			http.Error(respWriter, "Can't marshall devices", http.StatusInternalServerError)
			return
		}
		respWriter.Write(bytes)
	}
}

// /api/deletedevice?physaddr=MACAddr
// remove specific device from map
func deleteDeviceHandler(devices map[string]*Device) http.HandlerFunc {
	return func(respWriter http.ResponseWriter, request *http.Request) {
		mac := request.FormValue("physaddr")
		if mac == "" {
			http.Error(respWriter, "Specify Physical device's address", http.StatusBadRequest)
			return
		}
		delete(devices, mac)
		response := StatusResponse{
			OK:      true,
			Message: "Device successfully deleted.",
		}
		bytes, err := json.Marshal(&response)
		if err != nil {
			http.Error(respWriter, "Can't marshall response json: "+err.Error(), http.StatusInternalServerError)
		}
		respWriter.Header().Add("Content-Type", "application/json")
		respWriter.Write(bytes)
	}
}

// /web/*
// return fileserver for static pages
func getFileserverHandler(folder string) http.Handler {
	return http.StripPrefix("/"+folder+"/", http.FileServer(http.Dir("./"+folder)))
}
