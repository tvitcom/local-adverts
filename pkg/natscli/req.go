package main

import (
	"log"
	"time"
	"github.com/nats-io/nats.go"
)

func main() {
	res, err := sendReq("Cli1", "imagechannel", "file:1")
	if err != nil {
		log.Println("Error:",err.Error())
		return
	}
	log.Println("MQ OMMIT:", res)
}

func sendReq(cliName, subj, payload string) (string, error) {
	// Connect Options.
	opts := []nats.Option{nats.Name(cliName)}

	// Connect to NATS
	nc, err := nats.Connect(nats.DefaultURL, opts...)
	if err != nil {
		return "", err
	}
	defer nc.Close()
	
	msg, err := nc.Request(subj, []byte(payload), 2*time.Second)
	if err != nil {
		return "", err
	}
	if nc.LastError() != nil {
		return "", nc.LastError()
	}
	res := string(msg.Data)
	log.Printf("Published [%s] : '%s'", subj, payload)
	log.Printf("Received  [%v] : '%s'", msg.Subject, res)
	return res, nil
}