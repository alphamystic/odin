package handlers
/*
	Classifies a given vulnerability when found
*/
import (
	"fmt"
	//"time"
	//"strings"
	"net/http"
	//"math/rand"
	"github.com/alphamystic/odin/lib/utils"
)


// checkPayload checks the given payload for the specified vulnerability.
// this will work so well as an interface for various checkers
func checkPayload(payload string, vulnerability Vulnerability) bool {
	// Send the payload in a request to the target URL and check the response for indicators of the vulnerability.
	// For example, you could check the response for a known error message or a page that
	// Send the payload in a request to the target URL and check the response for indicators of the vulnerability.
	// For example, you could check the response for a known error message or a page that only appears when the vulnerability is present.
	// This is just an example and you will need to implement this part of the function based on your specific use case.
	targetURL := "http://example.com/index.php?file=../../etc/passwd"
	req, err := http.NewRequest("GET", targetURL+payload, nil)
	if err != nil {
		return false
	}
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return false
	}
	defer resp.Body.Close()
	if resp.StatusCode == 200 {
		// Check the response for indicators of the vulnerability.
		return true
	}
	return false
}

func PredictEncodingMethod(payload string) (e EncodingMethod){
	return
}
/* place the vulnerabilities into a range then for each, mutate to find the right payload
check the payload if successful and write to a success channel
when done send a done signal
*/
// should the parameter be in the form of a vulnerability?
// probably take in a channel of vulnerabilitites from a target and chain them into an exploit
// if an interface then we can implement each for one and called when needed
func Classifier() ([]Vulnerabilities,error) {
	utils.PrintTextInASpecificColorInBold("white","[+]		Loading payloads!!!!!!")
	/*payloads := LoadPayloads("LFI")
	// create a function to load the payload as specified
	var data []DataPoint
	for _, p := range payloads {
		features := []float64{float64(len(p)), float64(strings.Count(p, ".")), float64(strings.Count(p, "/"))}
		label := encodePayload(p, URLEncoding)
		data = append(data, DataPoint{features, label})
	}*/
	//targetURL := "http://example.com/index.php?file=../../etc/passwd"
	payload := "../../etc/passwd"
	encodingMethod := PredictEncodingMethod(payload)
	encodedPayload := encodePayload(payload, encodingMethod)
	if checkPayload(encodedPayload, LFI) {
		fmt.Println("LFI vulnerability found!")
	}
	if checkPayload(encodedPayload, RFI) {
		fmt.Println("RFI vulnerability found!")
	}
	if checkPayload(encodedPayload, SQLI) {
		fmt.Println("SQLI vulnerability found!")
	}
	if checkPayload(encodedPayload, SSRF) {
		fmt.Println("SSRF vulnerability found!")
	}
	if checkPayload(encodedPayload, IDOR) {
		fmt.Println("IDOR vulnerability found!")
	}
	if checkPayload(encodedPayload, CSRF) {
		fmt.Println("CSRF vulnerability found!")
	}
	return nil,nil
}

func LoadPayloads(name string) {
	/*var payloads = make([]string)

	payload := ""
	switch name {
	case "lfi":
		//open payloads directory and load all lfi ayloads
	case "xss":
		// load all xss
	default:
		//try looking for a zero day
		payloads = append(payloads,payload)
	}
	payloads = append(payloads,payload)
	return payloads*/
}
