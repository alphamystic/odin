package tunnel

import (
	"bytes"
	"fmt"
	"log"
	"net"

	"github.com/miekg/dns"
)

const serverAddr = ":53"

func RunDNSChunkBuilder() {
	// Listen for incoming DNS requests
	server := &dns.Server{Addr: serverAddr, Net: "udp"}
	server.Handler = dns.HandlerFunc(handleDNSRequest)
	if err := server.ListenAndServe(); err != nil {
		log.Fatalf("Failed to start DNS server: %s\n", err.Error())
	}
}

// Handle incoming DNS requests and extract data from responses
func handleDNSRequest(w dns.ResponseWriter, r *dns.Msg) {
	// Check if the request is a valid tunneling request
	qname := r.Question[0].Name
	if !isTunnelingRequest(qname) {
		dns.HandleFailed(w, r)
		return
	}

	// Reassemble the data from the DNS responses
	var data bytes.Buffer
	for i := 0; ; i++ {
		answerName := fmt.Sprintf("%d.%s", i, qname)
		answer := getAnswerForName(r, answerName)
		if answer == nil {
			break
		}
		data.Write(answer.A)
	}

	fmt.Printf("Received data: %s\n", data.Bytes())
}

// Check if the specified domain name is a valid tunneling request
func isTunnelingRequest(qname string) bool {
	parts := dns.SplitDomainName(qname)
	if len(parts) < 3 {
		return false
	}
	if parts[len(parts)-2] != "data" {
		return false
	}
	_, err := strconv.Atoi(parts[len(parts)-3])
	return err == nil
}

// Get the DNS answer for the specified domain name
func getAnswerForName(r *dns.Msg, name string) *dns.A {
	for _, answer := range r.Answer {
		if answer.Header().Name == name && answer.Header().Rrtype == dns.TypeA {
			a, ok := answer.(*dns.A)
			if ok {
				return a
			}
		}
	}
	return nil
}
