package main

import (
	"bytes"
	"encoding/binary"
	"errors"
	"fmt"
	"io"
	"time"

	"github.com/jacobsa/go-serial/serial"
)

const BAUD uint = 9600

type HardwareResetter struct {
	SerialConfiguration serial.OpenOptions
	Port                io.ReadWriteCloser
	SleepTime           time.Duration
}

//functions for HardwareResetter
func initResetter(portName string, baud uint) (*HardwareResetter, error) {
	resetter := &HardwareResetter{}
	resetter.SerialConfiguration = serial.OpenOptions{
		PortName:   portName,
		BaudRate:   baud,
		ParityMode: serial.PARITY_NONE,
		StopBits:   1,
		DataBits:   8,
	}
	port, err := serial.Open(resetter.SerialConfiguration)
	if err != nil {
		return nil, err
	}
	resetter.Port = port
	return resetter, nil
}

func (rst *HardwareResetter) serialHandshake() bool {
	respBuff := make([]byte, 2)
	io.ReadFull(rst.Port, respBuff)
	buff := []byte("HI")
	if bytes.Equal(buff, respBuff) {
		rst.Port.Write(buff)
		return true
	}
	return false
}

func (rst *HardwareResetter) serialSendTiming(dur time.Duration) {
	var millis int32 = int32(dur.Nanoseconds() / 1000 / 1000)
	binary.Write(rst.Port, binary.LittleEndian, millis)
	rst.SleepTime = dur
}

func (rst *HardwareResetter) startPolling() {
	go func() {
		for {
			var a byte = 97 // 'a'
			binary.Write(rst.Port, binary.LittleEndian, a)
			time.Sleep(rst.SleepTime)
		}
	}()
}

func findDevice(from uint, to uint) (*HardwareResetter, error) {
	for i := from; i <= to; i++ {
		deviceName := fmt.Sprintf("COM%d", i)
		port, err := initResetter(deviceName, BAUD)
		if err == nil && port.serialHandshake() {
			return port, nil
		}
	}
	return nil, errors.New("ResetterDevice not connected or unaviable")
}
