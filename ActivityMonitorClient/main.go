package main

import (
	"fmt"
	"os"
	"time"
)

const CONFIGFILE string = "configuration.json"

var configuration *AppConfig = loadConfig(CONFIGFILE)

func main() {
	//loading configuration
	if configuration == nil {
		fmt.Print("Enter name: ")
		var name string
		fmt.Scanln(&name)
		fmt.Print("Enter server addr in format host:port ")
		var addr string
		fmt.Scanln(&addr)
		configuration = &AppConfig{
			DeviceHash: getHashedMacAddr(),
			Name:       name,
			ServerAddr: addr,
		}
		configuration.saveConfig(CONFIGFILE)
	} else {
		if !haveMacAddr(configuration.DeviceHash) {
			fmt.Println("There is no registered device. Delete configuration file and run program again.")
			os.Exit(1)
		}
	}
	fmt.Printf("Loaded device hash: %s\nUser name: %s\nServer addr: %s\n", configuration.DeviceHash, configuration.Name, configuration.ServerAddr)

	//setup some variables
	serverAddr := fmt.Sprintf("http://%s", configuration.ServerAddr)
	chekArgs := map[string]string{
		"physaddr": configuration.DeviceHash,
		"name":     configuration.Name,
	}

	//main part
	initClient()
	//load ticktime from server
	ticktimestr := MakeGet(serverAddr + "/ticktime")
	ticktime, errDurr := time.ParseDuration(ticktimestr)
	if errDurr != nil {
		fmt.Println("Error in parsing timetick. Set to 30 seconds")
		fmt.Println(errDurr.Error())
		ticktime = time.Second * 30
	} else {
		ticktime = time.Duration(ticktime.Nanoseconds() / 2)
		fmt.Printf("According to server timetick set to %.0f seconds\n", ticktime.Seconds())
	}

	//initializing device if enabled
	if configuration.UseDevice {
		resetter, err := findDevice(1, 20)
		if err == nil {
			fmt.Printf("ResetterDevice found at: %s\n", resetter.SerialConfiguration.PortName)
			resetter.serialSendTiming(time.Second)
			resetter.startPolling()
		} else {
			fmt.Println(err.Error())
		}
	}

	//start main cycle
	ticker := time.NewTicker(ticktime)
	for _ = range ticker.C {
		response := MakeGetWithArgs(serverAddr+"/check", chekArgs)
		fmt.Printf("Server response: %s\n", response)
	}
}
