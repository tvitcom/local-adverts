package main

import (
	"os"
	"log"
	"errors"
	"os/signal"
	"time"

	"github.com/nats-io/nats.go"
)


func main() {
	err := listenReplyChannel("ImageService","imagechannel", "uid.1:1_123.jpg")
	if err != nil {
		log.Println(err.Error())
	}
}

func listenReplyChannel(qName, subj, reply string) error {
	// Options.
	opts := []nats.Option{nats.Name(qName)}
	opts, errConf := setupConnOptions(opts)
	if errConf != nil {
		return errConf
	}
	// Connect to NATS
	nc, err := nats.Connect(nats.DefaultURL, opts...)
	if err != nil {
		return err
	}

	i :=  0
	queueId := "NATS-RPLY-22"
	nc.QueueSubscribe(subj, queueId, func(msg *nats.Msg) {
		i++
		printMsg(msg, i)
		msg.Respond([]byte(reply))
	})
	nc.Flush()

	if err := nc.LastError(); err != nil {
		return err
	}

	log.Println(qName, "listening on", subj)

	// Setup the interrupt handler to drain so we don't miss
	// requests when scaling down.
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)
	<-c
	log.Println()
	log.Printf("Draining...")
	nc.Drain()
	
	return errors.New("Exiting")
}

func setupConnOptions(opts []nats.Option) ([]nats.Option, error) {
	totalWait := 10 * time.Minute
	reconnectDelay := time.Second

	opts = append(opts, nats.ReconnectWait(reconnectDelay))
	opts = append(opts, nats.MaxReconnects(int(totalWait/reconnectDelay)))
	opts = append(opts, nats.DisconnectErrHandler(func(nc *nats.Conn, err error) {
		log.Printf("Disconnected due to: %s, will attempt reconnects for %.0fm", err, totalWait.Minutes())
	}))
	opts = append(opts, nats.ReconnectHandler(func(nc *nats.Conn) {
		log.Printf("Reconnected [%s]", nc.ConnectedUrl())
	}))
	opts = append(opts, nats.ClosedHandler(func(nc *nats.Conn) {
		log.Fatalf("Exiting: %v", nc.LastError())
	}))
	return opts, nil
}

func printMsg(m *nats.Msg, i int) {
	log.Printf("[#%d] Received on [%s]: '%s'\n", i, m.Subject, string(m.Data))
}