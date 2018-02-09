package main

import (
	"fmt"
	"os"
	"os/signal"
)

const CONFIGFILE string = "configuration.json"

var serverInstance *Server = nil
var configuration *AppConfig = loadConfig(CONFIGFILE)

func main() {
	if configuration == nil {
		fmt.Println("No configuration fing. Exiting...")
		os.Exit(1)
	}
	//init server
	serverInstance = &Server{}
	serverInstance.init(configuration.ServerPort)
	//setup server stoping
	quit := make(chan os.Signal)
	signal.Notify(quit, os.Interrupt, os.Kill)
	go func() {
		sign := <-quit
		fmt.Printf("Catched signal: %v\n", sign)
		if err := serverInstance.stop(); err != nil {
			fmt.Println(err.Error())
		}
	}()
	//register handlers
	//client handlers
	serverInstance.Handlers.HandleFunc("/", notFoundHandler)
	serverInstance.Handlers.HandleFunc("/ticktime", tickTimeHandler(configuration.SleepTime+"s"))
	serverInstance.Handlers.HandleFunc("/check", checkCommunicationHandler(serverInstance.RegisteredDevices))
	//web handlers
	serverInstance.Handlers.Handle("/web/", getFileserverHandler("web"))
	serverInstance.Handlers.HandleFunc("/api/devicestatus", deviceStatusHandler(serverInstance.RegisteredDevices))
	serverInstance.Handlers.HandleFunc("/api/deletedevice", deleteDeviceHandler(serverInstance.RegisteredDevices))
	//start server
	var isFinish = serverInstance.start()
	fmt.Printf("Go to online console at http://localhost:%d/web\n", serverInstance.Port)
	defer fmt.Println("Server will stop maximum in 30 seconds...")
	<-isFinish //wait until program do not closed
}
