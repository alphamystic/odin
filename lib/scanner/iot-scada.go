package scanner

import (
	"fmt"
	"net"
	"strings"
)

func IotScadaScan(input net.Ip) {
	// Read user input
	fmt.Print("Enter a network range (e.g. 192.168.1.0/24): ")
	/*var input string
	fmt.Scanln(&input)*/
	// Parse the network range
	_, network, err := net.ParseCIDR(input)
	if err != nil {
		fmt.Println(err)
		return
	}
	// Scan the network for SCADA and IoT systems
	fmt.Println("Scanning network for SCADA and IoT systems...")
	for ip := network.IP.Mask(network.Mask); network.Contains(ip); inc(ip) {
		// Check if the current IP is a SCADA or IoT system
		if isSCADA(ip) || isIoT(ip) {
			fmt.Println("SCADA or IoT system detected:", ip)
		}
	}
	fmt.Println("Scan complete.")
}

// inc increments the IP address
func inc(ip net.IP) {
	for j := len(ip) - 1; j >= 0; j-- {
		ip[j]++
		if ip[j] > 0 {
			break
		}
	}
}

type SCADA struct {
	Protocol string
	Port int
	Vendor string
}
// isSCADA returns true if the IP address is a SCADA system
func isSCADA(ip net.IP) bool {
	// Check for common SCADA ports
	if isTCPPortOpen(ip, 502) || isTCPPortOpen(ip, 161) || isTCPPortOpen(ip, 502) {
		return true
	}
	// Check for common SCADA protocols
	if isModbusOpen(ip) || isDNP3Open(ip) || isIEC104Open(ip) {
		return true
	}
	// Check for common SCADA vendors
	if isSchneiderOpen(ip) || isSiemensOpen(ip) || isRockwellOpen(ip) {
		return true
	}
	return false
}

// isIoT returns true if the IP address is an IoT device
func isIoT(ip net.IP) bool {
	// Check for common IoT ports
	if isTCPPortOpen(ip, 80) || isTCPPortOpen(ip, 443) || isTCPPortOpen(ip, 8883) {
		return true
	}
	// Check for common IoT protocols
	if isMQTTOpen(ip) || isCOAPOpen(ip) || isHTTPOpen(ip) {
		return true
	}
	// Check for common IoT vendors
	if isNestOpen(ip) || isPhilipsOpen(ip) || isSamsungOpen(ip) {
		return true
	}
	return false
}

// isTCPPortOpen returns true if the specified TCP port is open on the given IP address
func isTCPPortOpen(ip net.IP, port int) bool {
	// Connect to the port
	conn, err := net.Dial("tcp", fmt.Sprintf("%s:%d", ip, port))
	if err != nil {
		return false
	}
	defer conn.Close()
	return true
}

// isModbusOpen returns true if the Modbus protocol is open on the given IP address
func isModbusOpen(ip net.IP) bool {
	// Connect to the Modbus port (502)
	conn, err := net.Dial("tcp", fmt.Sprintf("%s:%d", ip, 502))
	if err != nil {
		return false
	}
	defer conn.Close()
	// Send a Modbus request
	_, err = conn.Write([]byte{0x01, 0x03, 0x00, 0x00, 0x00, 0x02, 0xC4, 0x0B})
	if err != nil {
		return false
	}
	// Read the response
	buf := make([]byte, 1024)
	n, err := conn.Read(buf)
	if err != nil {
		return false
	}
	// Check if the response is a valid Modbus response
	return n >= 5 && buf[0] == 0x01 && buf[1] == 0x03 && buf[2] == 0x02 && buf[3] == 0xC4 && buf
}

// isDNP3Open returns true if the DNP3 protocol is open on the given IP address
func isDNP3Open(ip net.IP) bool {
	// Connect to the DNP3 port (20000)
	conn, err := net.Dial("tcp", fmt.Sprintf("%s:%d", ip, 20000))
	if err != nil {
		return false
	}
	defer conn.Close()
	// Send a DNP3 request
	_, err = conn.Write([]byte{0x05, 0x64, 0xC0, 0x01, 0x00, 0x00, 0x00, 0x14})
	if err != nil {
		return false
	}
	// Read the response
	buf := make([]byte, 1024)
	n, err := conn.Read(buf)
	if err != nil {
		return false
	}
	// Check if the response is a valid DNP3 response
	return n >= 5 && buf[0] == 0x05 && buf[1] == 0x64 && buf[2] == 0xC0 && buf[3] == 0x01 && buf[4] == 0x00
}

// isSchneiderOpen returns true if the Schneider Electric vendor is open on the given IP address
func isSchneiderOpen(ip net.IP) bool {
	// Check if the IP address responds to a request for the Schneider Electric website
	resp, err := http.Get(fmt.Sprintf("http://%s/schneider", ip))
	if err != nil {
		return false
	}
	defer resp.Body.Close()
	// Check if the response indicates support for Schneider Electric products
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return false
	}
	return strings.Contains(string(body), "Schneider Electric")
}

// isSiemensOpen returns true if the Siemens vendor is open on the given IP address
func isSiemensOpen(ip net.IP) bool {
	// Check if the IP address responds to a request for the Siemens website
	resp, err := http.Get(fmt.Sprintf("http://%s/siemens", ip))
	//@TODO Add a https possibility though minimal
	if err != nil {
		return false
	}
	defer resp.Body.Close()
	// Check if the response indicates support for Siemens products
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return false
	}
	return strings.Contains(string(body), "Siemens")
}

// isRockwellOpen returns true if the Rockwell Automation vendor is open on the given IP address
func isRockwellOpen(ip net.IP) bool {
	// Check if the IP address responds to a request for the Rockwell Automation website
	resp, err := http.Get(fmt.Sprintf("http://%s/rockwell", ip))
	if err != nil {
		return false
	}
	defer resp.Body.Close()
	// Check if the response indicates support for Rockwell Automation products
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return false
	}
	return strings.Contains(string(body), "Rockwell Automation")
}

// isHTTPOpen returns true if the given IP address is hosting an HTTP server
func isHTTPOpen(ip net.IP) bool {
	// Check if the IP address responds to a request for the HTTP protocol
	_, err := http.Get(fmt.Sprintf("http://%s/", ip))
	if err != nil {
		return false
	}
	return true
}

// isMQTTOpen returns true if the given IP address is hosting an MQTT server
func isMQTTOpen(ip net.IP) bool {
	// Connect to the MQTT port (1883)
	conn, err := net.Dial("tcp", fmt.Sprintf("%s:%d", ip, 1883))
	if err != nil {
		return false
	}
	defer conn.Close()
	// Send an MQTT CONNECT request
	_, err = conn.Write([]byte{0x10, 0x0e, 0x00, 0x04, 0x4d, 0x51, 0x54, 0x54, 0x04, 0xc2, 0x00, 0x3c, 0x00})
	if err != nil {
		return false
	}
	// Read the response
	buf := make([]byte, 1024)
	n, err := conn.Read(buf)
	if err != nil {
		return false
	}
	// Check if the response is a valid MQTT CONNACK response
	return n >= 2 && buf[0] == 0x20 && buf[1] == 0x02
}

// isCOAPOpen returns true if the COAP protocol is open on the given IP address
func isCOAPOpen(ip net.IP) bool {
	// Connect to the COAP port (5683)
	conn, err := net.Dial("udp", fmt.Sprintf("%s:%d", ip, 5683))
	if err != nil {
		return false
	}
	defer conn.Close()
	// Send a COAP GET request
	_, err = conn.Write([]byte{0x41, 0x01, 0x61, 0x61, 0x61, 0x61, 0x61, 0x61})
	if err != nil {
		return false
	}
	// Read the response
	buf := make([]byte, 1024)
	n, err := conn.Read(buf)
	if err != nil {
		return false
	}
	// Check if the response is a valid COAP response
	return n >= 2 && buf[0] == 0x51 && buf[1] == 0x01
}

// isNestOpen returns true if the Nest vendor is open on the given IP address
func isNestOpen(ip net.IP) bool {
	// Check if the IP address responds to a request for the Nest website
	resp, err := http.Get(fmt.Sprintf("http://%s/nest", ip))
	if err != nil {
		return false
	}
	defer resp.Body.Close()
	// Check if the response indicates support for Nest products
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return false
	}
	return strings.Contains(string(body), "Nest")
}

// isPhilipsOpen returns true if the Philips vendor is open on the given IP address
func isPhilipsOpen(ip net.IP) bool {
	// Check if the IP address responds to a request for the Philips website
	resp, err := http.Get(fmt.Sprintf("http://%s/philips", ip))
	if err != nil {
		return false
	}
	defer resp.Body.Close()
	// Check if the response indicates support for Philips products
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return false
	}
	return strings.Contains(string(body), "Philips")
}

// isSamsungOpen returns true if the Samsung vendor is open on the given IP address
func isSamsungOpen(ip net.IP) bool {
	// Check if the IP address responds to a request for the Samsung website
	resp, err := http.Get(fmt.Sprintf("http://%s/samsung", ip))
	if err != nil {
		return false
	}
	defer resp.Body.Close()
	// Check if the response indicates support for Samsung products
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return false
	}
	return strings.Contains(string(body), "Samsung")
}

// isIEC104Open returns true if the IEC 104 protocol is open on the given IP address
func isIEC104Open(ip net.IP) bool {
	// Connect to the IEC 104 port (2404)
	conn, err := net.Dial("tcp", fmt.Sprintf("%s:%d", ip, 2404))
	if err != nil {
		return false
	}
	defer conn.Close()
	// Send an IEC 104 request
	//_, err = conn.Write([]byte{0x68, 0x04, 0x7B, 0x12, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x68, 0x11, 0x04, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0}
	return true
}
