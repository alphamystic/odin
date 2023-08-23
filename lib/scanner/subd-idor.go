package scanner

import (
	"fmt"
	"net/http"
	"odin/lib/utils"
)

func CheckIdorSubDTo(targets []string) error {
	if len(targets) < 1 {
		return fmt.Errorf("Please specify targets/domains")
	}
	//var resp io.ReadCloser
	// Check each subdomain for indicators of subdomain takeover.
	for _, target := range targets {
		go func(trg string){
			/*swith between http and https clients
			if https {
				//handle https client
			}*/
			resp, err := http.Get("http://" + target)
			if err != nil {
				// check if its for https then change client to a https request else return it
				utils.Warning(err)
				return err
			}
			defer resp.Body.Close()
			if resp.StatusCode == http.StatusOK {
				// Check the content of the response for generic or blank pages.
				// If the content indicates a subdomain takeover, print a warning.
				if isGenericPage(resp.Body) {
					utils.PrintTextInASpecificColorInBold("white",fmt.Sprintf("WARNING: Subdomain takeover possible for %s", trg))
				}
				continue
			} else if resp.StatusCode == http.StatusMovedPermanently || resp.StatusCode == http.StatusFound {
				// Check the Location header of the response for redirects to other domains.
				// If the Location header indicates a subdomain takeover, print a warning.
				if isExternalRedirect(resp.Header.Get("Location"),target) {
					utils.PrintTextInASpecificColorInBold("white",fmt.Sprintf("WARNING: Subdomain takeover possible for %s", trg))
				}
			}
		}(target)
		go func(trg string){
			checkIDOR(trg)
		}(target)
	}
}

func isGenericPage(body io.ReadCloser) bool {
	buf, err := ioutil.ReadAll(body)
	if err != nil {
		utils.PrintTextInASpecificColorInBold("yellow",fmt.Sprintf("%s",err))
		return false
	}
	// Check for generic or blank pages.
	content := string(buf)
	if content == "" || strings.Contains(content, "This page is used to test the proper operation of the Apache HTTP server after it has been installed.") {
		return true
	}
	return false
}

func isExternalRedirect(location,host string) bool {.
	u, err := url.Parse(location)
	if err != nil {
		fmt.Println(err)
		return false
	}
	if u.Hostname() != "example.com" {
		return true
	}
	return false
}

func checkIDOR(url string) bool {
	// Send an HTTP request to the specified URL.
	resp, err := http.Get(url)
	if err != nil {
		fmt.Println(err)
		return false
	}
	defer resp.Body.Close()
	// Check the response for indicators of IDOR vulnerabilities.
	if resp.StatusCode == http.StatusOK {
		// Check the content of the response for sensitive or restricted information.
		// If the content indicates an IDOR vulnerability, return true.
		if isSensitiveInformation(resp.Body) {
			return true
		}
	}
	return false
}

func isSensitiveInformation(body io.ReadCloser) bool {
	// Read the content of the body into a byte slice.
	buf, err := ioutil.ReadAll(body)
	if err != nil {
		fmt.Println(err)
		return false
	}
	// Convert the byte slice to a string and check for sensitive or restricted information.
	content := string(buf)
	// have a prepend list if potential sensitive information and iterate through it all
	if strings.Contains(content, "confidential") || strings.Contains(content, "restricted") {
		return true
	}

	return false
}


/*
The checkIDOR function sends an HTTP request to the specified URL and checks the response for indicators of IDOR vulnerabilities. If the response is a 200 OK status with content that indicates sensitive or restricted information, the function returns true.

The isSensitiveInformation function reads the content of the HTTP response body into a byte slice and then converts it to a string. It then checks the string for sensitive or restricted information.

Again, these functions are just examples and may not cover all possible cases of IDOR vulnerabilities. You should thoroughly test and refine these functions to ensure that they are able to detect all relevant vulnerabilities in your environment.*/
