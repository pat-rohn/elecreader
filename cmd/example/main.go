package main

import (
	"fmt"
	"time"

	"github.com/pat-rohn/elecreader"
	log "github.com/sirupsen/logrus"
	"github.com/tarm/serial"
)

func main() {
	fmt.Print("test")
	c := serial.Config{
		Name:        "/dev/ttyUSB0",
		Baud:        300,
		Size:        7,
		Parity:      serial.ParityEven,
		StopBits:    serial.Stop1,
		ReadTimeout: time.Second * 2,
	}
	p, err := elecreader.OpenPort(&c)
	if err != nil {
		log.Fatalf("Failed to open port: '%s' ", err)
	}
	conn := elecreader.Connection{
		Port: p,
	}
	defer conn.Port.Close()
	answer, err := conn.Read()
	if err != nil {
		log.Fatalf("Failed to read: %v", err)
	}

	res, err := elecreader.Extract(answer)
	if err != nil {
		log.Fatalf("Failed to extract values: %v", err)
	}
	fmt.Printf("%v", res)
}
