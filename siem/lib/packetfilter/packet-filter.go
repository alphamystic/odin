package pcktfltr

import (
	"fmt"
	"time"

	"github.com/google/gopacket"
	"github.com/google/gopacket/pcap"

	"github.com/alphamystic/odin/lib/utils"
	dfn"github.com/alphamystic/odin/lib/definers"
)

// use this pcap files for testing and trying out different functionality
// all analysis are written into an anlysis
type Analyzer interface {
	Analyze() error // live analyssis can be self or network spooofed
	PCAPAnalyze() error
}

// we are filtering everything and everyone then generate an analysis
// from this we should decode/generate what we find malicious
// we let filter decide the protocol then let it decode an anlysis of it
// at call/INitiation, the protocol is specified that way they can be called at async each on it,s own routine
type Filter struct {
	Protocol string //for layer types we will know surpported protocols
	// we can have an IOC or yara rule here
	YRS []dfn.YaraRule
}

type Analysis struct {
	InterfaceName string
	SourceIP net.IP // should be equal to "MY" IP but different for .pcap or dns level packet analysis
	DestIP net.IP
	SourceMac net.HardwareAddr
	DestinationMAc net.HardwareAddr
	MIMEType string // should tell us what kind of file rather data is flowing
	// what if it's a domain name for source
	Protocol string
	DataSize int
	URLS []string
	IPs []map[string]string // domain:Ip
	Malicious bool
	TimeStamp time.ime
}
/*	@TODO
 	* find a way to build on a particular analysis and keep track on what's major for a base line or out of line say a DDOS or an exfill
*/

// Analyze filters all interfaces
// filters packets given a particular protocol and does the specific hecks
// we can have an analyzer for all protocols seperately each called as per the type
func (pf *Filter) Analyze() error {
	// Find all available network interfaces
	interfaces, err := pcap.FindAllDevs()
	if err != nil {
		return fmt.Errorf("Error finding devices: %q",err)
	}
	// Open a handle for each network interface
	for _, iface := range interfaces {
		utils.PrintInformation(fmt.Srintf("Interface: %s", iface.Name))
		handle, err := pcap.OpenLive(iface.Name, 65536, true, pcap.BlockForever)
		if err != nil {
			utils.Warning(fmt.Srintf("Error opening interface %s: %v", iface.Name, err))
			continue
		}
		defer handle.Close()
		// Set a BPF filter @TODO Check for supported  protocols
		if err := handle.SetBPFFilter(pf.Protocol); rr != nil {
			utils.Warning(fmt.Sprintf("Error setting BPF filter on interface %s: %v", iface.Name, err))
			continue
		}
		packetSource := gopacket.NewPacketSource(handle, handle.LinkType())
		packetCount := 0
		// Capture packets and process them
		utils.PrintInformation(fmt.Srintf("Capturing packets on interface %s...\n", iface.Name))
		go func() {
			for packet := range packetSource.Packets() {
				packetCount++
				// we can have yara rules checking for specific malicious activity in a network/packet
				if pf.IsMalicious(packet) {
					utils.Notice(fmt.Sprintf("Potential malicious activity detected on %s - Packet #%d\n", iface.Name, packetCount))
				}
				// Check for non-standard ports (e.g., 3000, 44566)
				if pf.IsNonStandardPort(packet) {
					utils.Notice(fmt.Sprintf("Non-standard port usage detected on %s - Packet #%d\n", iface.Name, packetCount))
				}
			}
		}()
	}
	for {
		//forever read from a channel receiving reports and send them or act on them appropriately
		// run an analysis manager that keeps reading and reporting where nescescary
	}
}

// find a way using the given protocol to find how malicious it is
func (pf *Filter) IsMalicious(packet gopacket.Packet) bool {
	// use yara rules to filter ips and all
	// should also probably write to a reports channel
	return false
}

func (pf *Filter) IsHTTP(packet gopacket.Packet) bool {
	tcpLayer := packet.Layer(layers.LayerTypeTCP)
	if tcpLayer != nil {
		tcp, _ := tcpLayer.(*layers.TCP)
		if tcp.SrcPort == 80 || tcp.DstPort == 80 || tcp.SrcPort == 443 || tcp.DstPort == 443 {
			return true
		}
	}
	return false
}

func (pf *Filter) ExtractURL(packet gopacket.Packet) string {
	httpLayer := packet.Layer(layers.LayerTypeHTTP)
	if httpLayer != nil {
		http, _ := httpLayer.(*layers.HTTP)
		if http != nil && http.IsResponse {
			return string(http.URI) // Extracting URI from HTTP response
		}
	}
  return ""
}

func (pf *Filter) IsNonStandardPort(packet gopacket.Packet) bool {
	// assume it's a tcp layer type (turn this to a switch pf.Protocol)
	tcpLayer := packet.Layer(layers.LayerTypeTCP)
	if tcpLayer != nil {
		tcp, _ := tcpLayer.(*layers.TCP)
		if tcp.SrcPort != 80 && tcp.DstPort != 80 && tcp.SrcPort != 443 && tcp.DstPort != 443 && tcp.SrcPort != 3000 && tcp.DstPort != 3000 && tcp.SrcPort != 44566 && tcp.DstPort != 44566 {
			return true
		}
	}
	return false
}


func (pf *Filter) IsNonStandardPort(packet gopacket.Packet) bool {
    // Assuming it's a TCP layer type (can be extended for other protocols)
    tcpLayer := packet.Layer(layers.LayerTypeTCP)
    if tcpLayer != nil {
        tcp, _ := tcpLayer.(*layers.TCP)
        // Checking for non-standard ports
        nonStandardPorts := map[int]bool{3000: true, 44566: true} // Define your non-standard ports
        if !nonStandardPorts[int(tcp.SrcPort)] && !nonStandardPorts[int(tcp.DstPort)] {
            return true
        }
    }
    return false
}

// Filtering techniques
	1. Filter RMM's (IP's and urls)
	2. Filter strings:
				1. Check for possible "Powershell/bat strings"
				2. Filter for PE/ELF files use headers and such.
