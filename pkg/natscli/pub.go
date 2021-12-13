package natscli

import (
	"log"
	"github.com/nats-io/nats.go"
)
// nats-pub -s demo.nats.io <subject> <text>

// func main() {
// 	log.Println("Published:", systemPubMessage("Frontend localadverts Publisher","image", "123_1.jpg" ))
// }

//!!! TODO: https://docs.nats.io/developing-with-nats/tutorials/custom_dialer
func PubMessage(owner, topic, content string) bool {
	// Connect Options.
	opts := []nats.Option{nats.Name("NATS Sample Requestor")}

	// Connect to NATS
	var err error
	nc, err := nats.Connect(nats.DefaultURL, opts...)
	if err != nil {
		// log.Fatal(err)
		return false
	}
	defer nc.Close()
	nc.Publish(topic, []byte(content))
	nc.Flush()
	if err := nc.LastError(); err != nil {
		// log.Fatal(err)
		return false
	}
	return  true
}
