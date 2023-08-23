package tunnel

import (
	"fmt"
	"log"
	"net"

	"github.com/miekg/dns"
)

func SendDNSData() {
	c := dns.Client{}
	m := dns.Msg{}

	// Create a DNS message with a custom question
	m.SetQuestion("example.com.", dns.TypeA)

	// Encode the message into binary form
	b, err := m.Pack()
	if err != nil {
		log.Fatal(err)
	}

	// Encode the binary message into a DNS query
	qname := dns.Fqdn("dns.example.com")
	req := new(dns.Msg)
	req.SetQuestion(qname, dns.TypeA)
	req.RecursionDesired = true
	req.Id = dns.Id()

	// Send the query to the remote server
	r, _, err := c.Exchange(req, net.JoinHostPort("8.8.8.8", "53"))
	if err != nil {
		log.Fatal(err)
	}

	// Decode the DNS response and extract the message data
	var msg dns.Msg
	err = msg.Unpack(r.Answer[0].(*dns.A).A)
	if err != nil {
		log.Fatal(err)
	}

	// Print the extracted message data
	fmt.Printf("Received message: %s\n", msg.Question[0].Name)
}
