package main

import (
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strconv"
	"sync"
	"time"
)

//TODO replace it with configuration file
//const SAVEFILE string = "clients.json"
//const PORT string = "12468"
//const SLEEPTIME string = "15s"

var devicesMutex = &sync.Mutex{}

type Device struct {
	PhysicalAddress string
	Name            string
	LastResponse    int64
	IsActive        bool
}

type Server struct {
	IsInited          bool
	Port              int
	RegisteredDevices map[string]*Device
	HttpServer        *http.Server
	Handlers          *http.ServeMux
	TimeTicker        *time.Ticker
}

func (s *Server) init(port string) {
	if !s.IsInited {
		s.loadDevices(configuration.ClientsSaveFile)
		prt, _ := strconv.Atoi(port)
		s.Port = prt
		slpTime, _ := time.ParseDuration(configuration.SleepTime + "s")
		s.TimeTicker = time.NewTicker(slpTime)
		s.Handlers = http.NewServeMux()
		s.HttpServer = &http.Server{
			Addr:         ":" + port,
			ReadTimeout:  time.Second * 20,
			WriteTimeout: time.Second * 20,
			Handler:      s.Handlers,
		}
		s.HttpServer.SetKeepAlivesEnabled(false)
		s.IsInited = true
	}
}

func (s *Server) start() chan bool {
	if !s.IsInited {
		s.init(configuration.ServerPort)
	}
	var done = make(chan bool)
	startDeviceValidator(s.TimeTicker, s.RegisteredDevices)
	go func() {
		s.HttpServer.ListenAndServe()
		done <- true
		return
	}()
	return done
}

func (s *Server) stop() error {
	s.saveDevices(configuration.ClientsSaveFile)
	s.TimeTicker.Stop()
	timeout, _ := time.ParseDuration("30s")
	ctx, _ := context.WithTimeout(context.Background(), timeout)
	return s.HttpServer.Shutdown(ctx)
}

func (s *Server) saveDevices(filename string) {
	bytes, err := json.Marshal(&s.RegisteredDevices)
	if err != nil {
		log.Println("Error in marshaling saved devices")
		log.Println(err.Error())
		return
	}
	if err = ioutil.WriteFile(filename, bytes, 0600); err != nil {
		log.Println("Error in writing saved devices:")
		log.Println(err.Error())
	}
}

func (s *Server) loadDevices(filename string) {
	if _, existErr := os.Stat(filename); existErr != nil {
		//file do not exists
		if os.IsNotExist(existErr) {
			log.Printf("File %s can't be found. Initializing new storage of saved devices.", filename)
			s.RegisteredDevices = map[string]*Device{}
		}
	} else {
		bytes, rdErr := ioutil.ReadFile(filename)
		//if file read insuccessfully
		if rdErr != nil {
			log.Printf("Error ocured during reading saved devices in: %s:\n", filename)
			log.Println(rdErr.Error())
		}
		//initialize an empty map
		s.RegisteredDevices = map[string]*Device{}
		if marshallErr := json.Unmarshal(bytes, &s.RegisteredDevices); marshallErr != nil {
			log.Printf("Error in unmarshaling %s\n", filename)
			log.Println(marshallErr.Error())
		}
	}
}

func startDeviceValidator(ticker *time.Ticker, devices map[string]*Device) {
	go func() {
		log.Println("Device validator service started")
		for tm := range ticker.C {
			devicesMutex.Lock()
			waitTime, _ := time.ParseDuration(configuration.SleepTime + "s")
			for _, dev := range devices {
				lastResponseTime := time.Unix(dev.LastResponse, 0)
				if tm.Sub(lastResponseTime) > waitTime && dev.IsActive {
					dev.IsActive = false
					notification := fmt.Sprintf("Device %s with name %s go offline!!!\n", dev.PhysicalAddress, dev.Name)
					log.Printf(notification)
					botNotify(notification)
				}
			}
			devicesMutex.Unlock()
		}
		log.Println("Device validator service started")
	}()
}
