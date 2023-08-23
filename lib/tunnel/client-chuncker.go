package tunnel

import (
	"fmt"
	"log"
	"net"
	"strconv"
	"strings"
	"time"

	"github.com/miekg/dns"
)

const serverAddr = "8.8.8.8:53"
const maxChunkSize = 250 // maximum size of each chunk

func ClientChunk() {
	data := "Hello, world! This is some sample data to be tunneled over DNS using chunking."

	// Split the data into chunks
	chunks := splitIntoChunks([]byte(data), maxChunkSize)

	// Send each chunk as a separate DNS request
	for i, chunk := range chunks {
		err := sendChunk(fmt.Sprintf("%d.%s", i, serverAddr), chunk)
		if err != nil {
			log.Fatal(err)
		}
	}
}

// Split data into chunks of the specified maximum size
func splitIntoChunks(data []byte, chunkSize int) [][]byte {
	var chunks [][]byte
	for i := 0; i < len(data); i += chunkSize {
		end := i + chunkSize
		if end > len(data) {
			end = len(data)
		}
		chunks = append(chunks, data[i:end])
	}
	return chunks
}

// Send a single chunk as a DNS request to the specified server
func sendChunk(server string, chunk []byte) error {
	// Encode the chunk into a DNS query
	qname := fmt.Sprintf("%s.data.example.com", time.Now().UnixNano())
	req := new(dns.Msg)
	req.SetQuestion(dns.Fqdn(qname), dns.TypeA)
	req.RecursionDesired = true
	req.Id = dns.Id()

	reqLen := len(req.String())
	// Adjust reqLen to fit maximum query size if necessary
	if reqLen+len(chunk) > dns.MinMsgSize {
		reqLen = dns.MinMsgSize - len(chunk)
	}

	for i := 0; i < len(chunk); i += reqLen {
		end := i + reqLen
		if end > len(chunk) {
			end = len(chunk)
		}
		req.Question[0] = dns.Question{
			Name:   dns.Fqdn(qname),
			Qtype:  dns.TypeA,
			Qclass: dns.ClassINET,
		}
		req.Question[0].Name = fmt.Sprintf("%d.%s", i/reqLen, req.Question[0].Name)

		// Create a DNS client and send the query to the server
		c := dns.Client{}
		r, _, err := c.Exchange(req, server)
		if err != nil {
			return err
		}
		if len(r.Answer) == 0 {
			return fmt.Errorf("No DNS answer received for chunk %d", i/reqLen)
		}
	}
	return nil
}
